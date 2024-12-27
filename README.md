# IP Counter

**IP Counter** is a highly efficient Go-based solution for counting unique IPv4 addresses from extremely large files. It is optimized for memory usage and execution speed, handling files up to 115GB or more.

## Features

- **Scalability**: Handles massive files without loading them entirely into memory.
- **Efficiency**: Uses segmented bitmaps to minimize memory usage.
- **Concurrency**: Employs a worker pool for parallel batch processing.
- **Accuracy**: Guarantees precise counting of unique IPs.

---

## Installation

1. Ensure you have **Go 1.19+** installed on your system.
2. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/ip-counter.git
   cd ip-counter
   ```
3. Build the application:
    ```bash
   go build -o ip_counter
    ```

## Usage
1. Prepare a file containing one IPv4 address per line.
2. Run the program:
    ```bash
    /ip_counter /path/to/ip_addresses.txt
    ```

## Example Output
```bash
$ ./ip_counter /path/to/ip_addresses.txt
2024/12/24 10:00:00 Unique IP addresses: 14562387
2024/12/24 10:00:00 Execution time: 12m48.562s
```
## Configuration Parameters

| Parameter        | Default Value | Description                               |
|------------------|---------------|-------------------------------------------|
| `BatchSize`      | 250,000       | Number of IPs processed per batch.        |
| `WorkerCount`    | 8             | Number of concurrent worker goroutines.   |
| `ReadBufferSize` | 2MB           | Buffer size for file reading.             |
| `BitmapSegments` | 256           | Number of bitmap segments for parallelism.|

---


## Technical Highlights
**Memory Optimization**:
- Uses segmented bitmaps to track unique IPs efficiently.
**Concurrency**:
- Employs worker goroutines for parallel processing.
**Error Handling**:
- Logs issues such as invalid input or file errors.
**Resource Management**:
- Ensures cleanup of memory and resources upon completion.

## Known Limitations
- Only supports IPv4 addresses.
- Does not handle corrupted or non-standard file formats.

## Potential Improvements
- Add support for IPv6 addresses.
- Implement dynamic configuration for batch sizes based on system memory.
- Enhance the user interface with a progress bar or more detailed metrics.

## License
This project is licensed under the MIT License.
