package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	BatchSize       = 250000                              // Number of IPs processed in one batch. Optimized for memory and performance.
	WorkerCount     = 8                                   // Number of goroutines (workers) for parallel processing.
	MaxIP           = 1 << 32                             // Total possible IPv4 addresses (2^32).
	ReadBufferSize  = 2 * 1024 * 1024                     // Buffer size for file reading in bytes (2 MB).
	BitmapSegments  = 256                                 // Number of segments for bitmap parallelism.
	BitsInByte      = 8                                   // Number of bits in one byte.
	BytesPerSegment = MaxIP / BitmapSegments / BitsInByte // Bytes per segment in the segmented bitmap.
)

// SegmentedBitmap tracks unique IPs using a segmented bitmap for thread-safe parallel processing.
type SegmentedBitmap struct {
	segments [BitmapSegments][]byte     // Segmented storage for reduced contention.
	mutexes  [BitmapSegments]sync.Mutex // Mutexes for thread-safe access to segments.
}

// NewSegmentedBitmap initializes a SegmentedBitmap and allocates memory for each segment.
func NewSegmentedBitmap() *SegmentedBitmap {
	sb := &SegmentedBitmap{}
	for i := range sb.segments {
		sb.segments[i] = make([]byte, BytesPerSegment) // Allocate memory for each segment.
	}
	return sb
}

// MarkIP marks an IP as seen in the segmented bitmap.
// Returns true if the IP was not previously marked, false otherwise.
func (sb *SegmentedBitmap) MarkIP(ipInt uint32) bool {
	// Determine the segment and position within the segment for the IP.
	segmentIndex := ipInt / (MaxIP / BitmapSegments)
	localIP := ipInt % (MaxIP / BitmapSegments)
	byteIndex := localIP / BitsInByte
	bitOffset := localIP % BitsInByte

	sb.mutexes[segmentIndex].Lock() // Lock the segment for safe access.
	defer sb.mutexes[segmentIndex].Unlock()

	// Check if the bit corresponding to the IP is already set.
	if sb.segments[segmentIndex][byteIndex]&(1<<bitOffset) != 0 {
		return false // IP was already marked.
	}

	// Mark the IP in the bitmap.
	sb.segments[segmentIndex][byteIndex] |= 1 << bitOffset
	return true
}

// Cleanup releases the memory allocated for the bitmap segments to free resources.
func (sb *SegmentedBitmap) Cleanup() {
	for i := range sb.segments {
		sb.mutexes[i].Lock()
		sb.segments[i] = nil // Clear segment to allow garbage collection.
		sb.mutexes[i].Unlock()
	}
}

// ipToInt converts an IP address string (e.g., "192.168.1.1") to a unique 32-bit integer representation.
// Each segment of the IP is converted to an integer and shifted to its position.
func ipToInt(ip string) uint32 {
	var result uint32
	segments := strings.Split(ip, ".") // Split IP into four segments.
	for i, segment := range segments {
		value := uint32(0)
		for _, c := range segment {
			value = value*10 + uint32(c-'0') // Convert string segment to integer.
		}
		result += value << (BitsInByte * uint(3-i)) // Shift each segment based on position.
	}
	return result
}

// processBatch processes a batch of IPs, marking unique ones in the bitmap.
// Returns the count of new unique IPs found in the batch.
func processBatch(batch []string, sb *SegmentedBitmap) (int, error) {
	uniqueCount := 0
	for _, ip := range batch {
		if ip == "" {
			continue // Skip empty lines.
		}
		ipInt := ipToInt(ip) // Convert IP to integer representation.
		if sb.MarkIP(ipInt) {
			uniqueCount++ // Increment count if IP was newly marked.
		}
	}
	return uniqueCount, nil
}

// worker processes batches from a channel and sends results to a result channel.
func worker(batchChan <-chan []string, sb *SegmentedBitmap, resultChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done() // Notify WaitGroup when the worker is done.

	for batch := range batchChan {
		uniqueCount, err := processBatch(batch, sb)
		if err != nil {
			log.Printf("Error processing batch: %v", err)
			continue
		}
		resultChan <- uniqueCount // Send the count of unique IPs to the result channel.
	}
}

// readBatch reads a batch of lines (IPs) from the file.
func readBatch(reader *bufio.Reader) ([]string, error) {
	batch := make([]string, 0, BatchSize)
	for len(batch) < BatchSize {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF && len(batch) > 0 {
				return batch, nil // Return the last batch if EOF is reached.
			}
			return nil, err // Return error if encountered.
		}
		batch = append(batch, strings.TrimSpace(line)) // Add trimmed line to the batch.
	}
	return batch, nil
}

// countUniqueIPs orchestrates the process of counting unique IPs in a file.
func countUniqueIPs(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err) // Standardized error handling.
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, ReadBufferSize) // Efficient buffered file reading.
	sb := NewSegmentedBitmap()                          // Initialize segmented bitmap.
	defer sb.Cleanup()                                  // Ensure memory is released.

	batchChan := make(chan []string, WorkerCount) // Channel for passing batches to workers.
	resultChan := make(chan int, WorkerCount)     // Channel for collecting results from workers.

	var wg sync.WaitGroup // Synchronizes worker goroutines.

	// Start worker goroutines.
	for i := 0; i < WorkerCount; i++ {
		wg.Add(1)
		go worker(batchChan, sb, resultChan, &wg)
	}

	// Goroutine to read the file and send batches to workers.
	go func() {
		defer close(batchChan) // Close the batch channel after all batches are sent.
		for {
			batch, err := readBatch(reader)
			if err == io.EOF {
				break // End of file reached.
			}
			if err != nil {
				log.Printf("Error reading batch: %v", err)
				break
			}
			batchChan <- batch // Send batch to workers.
		}
	}()

	// Goroutine to close the result channel when all workers are done.
	go func() {
		wg.Wait()         // Wait for all workers to complete.
		close(resultChan) // Close the result channel.
	}()

	totalUnique := 0
	for uniqueCount := range resultChan {
		totalUnique += uniqueCount // Accumulate unique counts from all batches.
	}

	return totalUnique, nil
}

func main() {
	start := time.Now()                                           // Start timing the execution.
	filePath := "/Users/kirillmalovicko/go/src/test/ip_addresses" // Path to the file.

	uniqueIPs, err := countUniqueIPs(filePath)
	if err != nil {
		log.Fatalf("Error: %v", err) // Log and exit on error.
	}

	duration := time.Since(start) // Calculate execution time.
	log.Printf("Unique IP addresses: %d", uniqueIPs)
	log.Printf("Execution time: %v", duration)
}
