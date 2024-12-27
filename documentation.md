# IP Counter - Documentation

## Overview

The IP Counter is designed to efficiently count unique IPv4 addresses in large text files, exceeding 115GB. The solution leverages Go’s concurrency model and memory-efficient data structures for high performance.

---

## Design Architecture

---
### Fundamental Principles
1. **Scalability**:
   - Designed to handle files larger than system memory. 
   - Uses a segmented bitmap to minimize contention during concurrent writes.
2. **Efficiency**:
   - Memory usage is optimized using batch processing and a fixed bitmap size.
3. **Robustness**:
   - Comprehensive error handling and resource cleanup ensure stability.

---
### Components
1. **Segmented Bitmap**:
   - Tracks unique IPs using a memory-efficient bit representation.
   - Divides the bitmap into multiple segments to reduce contention in concurrent environments.

2. **Batch Processing**:
   - Processes the file in manageable chunks (batches) to avoid memory overflow.

3. **Worker Pool**:
   - Distributes workload among multiple goroutines for parallel processing.

4. **Pipeline**:
   - Uses Go channels for efficient communication between file reader, workers, and result aggregator.

---
## Workflow

1. **File Reading**:
   - The file is read line-by-line into memory batches.
   - Each batch is sent to a worker pool for processing.

2. **Batch Processing**:
   - Each IP is converted into a unique integer using a specialized conversion function.
   - Unique IPs are tracked using a segmented bitmap.

3. **Aggregation**:
   - Results from all workers are aggregated to produce the final count of unique IPs.

4. **Cleanup**:
   - Memory used by the bitmap is released after processing.

---
## Function Details

### ipToInt(ip string) uint32

- **Purpose**: Converts an IPv4 address into a unique 32-bit integer.
- **Logic**:
  - Splits the IP into four segments.
  - Converts each segment to an integer and shifts it to its respective byte position.
- **Example**:
  - Input: `"192.168.0.1"`
  - Output: `3232235521`

---
### processBatch(batch []string, sb *SegmentedBitmap) (int, error)

- **Purpose**: Processes a batch of IP addresses to identify unique entries.
- **Steps**:
  - Convert each IP to an integer.
  - Mark the IP in the segmented bitmap.
  - Return the count of new unique IPs.

---
### worker(batchChan <-chan []string, sb *SegmentedBitmap, resultChan chan<- int, wg *sync.WaitGroup)

- **Purpose**: Processes batches from the input channel and sends results to the output channel.
- **Steps**:
  - Reads batches from the channel.
  - Processes each batch using `processBatch`.
  - Sends the result (number of unique IPs) to the result channel.

---
### countUniqueIPs(filePath string) (int, error)

- **Purpose**: Orchestrates the file processing workflow.
- **Steps**:
  - Opens the file and initializes channels, bitmap, and worker pool.
  - Reads the file in batches and sends them to workers.
  - Aggregates results from the worker pool.

---
## Configuration Parameters

| Parameter        | Default Value | Description                               |
|------------------|---------------|-------------------------------------------|
| `BatchSize`      | 250,000       | Number of IPs processed per batch.        |
| `WorkerCount`    | 8             | Number of concurrent worker goroutines.   |
| `ReadBufferSize` | 2MB           | Buffer size for file reading.             |
| `BitmapSegments` | 256           | Number of bitmap segments for parallelism.|

---
## Error Handling

1. **File Errors**:
   - Handles missing or unreadable files gracefully.
2. **Input Validation**:
   - Skips invalid or empty lines.
3. **Memory Issues**:
   - Ensures memory usage is optimized and cleans up resources.

---
## Performance Metrics

- **Memory Usage**:
  - Approximately 1 byte per unique IP.
- **Execution Time**:
  - Dependent on file size and system resources.
  - File size: 115GB → ~30 minutes on an 8-core CPU with 16GB RAM.

---
## Limitations

- Supports only IPv4 addresses.
- Requires the file to have one IP per line.

---
## Future Enhancements

1. Add support for IPv6.
2. Implement dynamic worker pool adjustment based on system load.
3. Provide a progress bar for better user experience.

## FAQ
### Q1: Why use a segmented bitmap?
**A**: It allows concurrent marking by multiple workers without significant contention.
### Q2: How does BatchSize affect performance?
**A**: Larger batches reduce I/O operations but increase memory usage. The default size balances performance and resource usage.