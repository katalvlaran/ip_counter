package main

//import (
//	"bufio"
//	"fmt"
//	"os"
//	"strconv"
//	"strings"
//	"sync"
//)
//
//// Adjust batch size as needed
//const BatchSize = 10000000
//
//// Convert IP address to an integer representation
//func ipToInt(ip string) uint {
//	var result uint = 0
//	parts := strings.Split(ip, ".")
//	for i, segment := range parts {
//		num, _ := strconv.Atoi(segment)
//		result += uint(num) << (8 * uint(3-i))
//	}
//	return result
//}
//
//// Worker function to process a batch of IPs
//func processBatch(batch []string, uniqueIPs map[uint]struct{}, mutex *sync.Mutex, wg *sync.WaitGroup) {
//	fmt.Println("another batch started")
//	defer wg.Done()
//	batchMap := make(map[uint]struct{})
//
//	// Process each IP in the batch
//	for _, ip := range batch {
//		ipInt := ipToInt(ip)
//		batchMap[ipInt] = struct{}{}
//	}
//
//	// Write batchMap into the shared uniqueIPs map with mutex
//	mutex.Lock()
//	for ip := range batchMap {
//		uniqueIPs[ip] = struct{}{}
//	}
//	mutex.Unlock()
//	fmt.Printf("batch finished with %d uniqIP\n", len(batchMap))
//}
//
//func countUniqueIPs(filePath string) (int, error) {
//	uniqueIPs := make(map[uint]struct{}) // Shared map for storing unique IPs
//	var mutex sync.Mutex                 // Mutex for synchronizing access to uniqueIPs
//	var wg sync.WaitGroup                // WaitGroup for tracking goroutines
//
//	file, err := os.Open(filePath)
//	if err != nil {
//		return 0, err
//	}
//	defer file.Close()
//
//	scanner := bufio.NewScanner(file)
//	batch := make([]string, 0, BatchSize)
//
//	for scanner.Scan() {
//		batch = append(batch, scanner.Text())
//
//		// If batch is full, process it in a goroutine
//		if len(batch) == BatchSize {
//			wg.Add(1)
//			go processBatch(batch, uniqueIPs, &mutex, &wg)
//			batch = make([]string, 0, BatchSize) // Reset the batch
//		}
//	}
//
//	// Process remaining IPs in the final batch
//	if len(batch) > 0 {
//		wg.Add(1)
//		go processBatch(batch, uniqueIPs, &mutex, &wg)
//	}
//
//	wg.Wait() // Wait for all goroutines to finish
//
//	if err := scanner.Err(); err != nil {
//		return 0, err
//	}
//
//	return len(uniqueIPs), nil
//}
//
//func main() {
//	filePath := "/path/to/ip_addresses" // Replace with your file path
//
//	count, err := countUniqueIPs(filePath)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	fmt.Println("Unique IPs:", count)
//}
