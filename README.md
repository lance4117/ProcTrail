# ProcTrail

ProcTrail is a cross-platform process metrics collection tool that supports collecting various performance metrics from target processes through PID or process name, with multiple data export formats.

[English](README.md) | [中文](README_zh.md)

## Features

- Cross-platform support (Windows, Linux, macOS)
- Multiple process selectors
  - By Process ID (PID)
  - By Process Name
- Rich metrics collection
  - CPU Usage (%)
  - Memory Usage (MB)
  - I/O Read/Write Bytes
  - Thread Count
  - File Descriptor Count
- Flexible sampling options
  - One-time sampling
  - Timed sampling (configurable interval)
  - Sampling duration control
- Multiple export formats
  - CSV
  - JSON
  - JSONL (JSON Lines)
  - XML
- Pretty print support (JSON/XML)

## Installation

### From Source

Ensure you have Go 1.16 or higher installed, then run:

```bash
git clone https://github.com/lance4117/ProcTrail.git
cd ProcTrail
go install ./cmd/proctrail
```

## Usage

### Basic Command Format

```bash
proctrail [flags]
```

### Required Parameters (Choose One)

- `--pid`: Process ID to monitor
- `--name`: Process name to monitor

### Optional Parameters

- `--interval`: Sampling interval in seconds (default: 0, one-time sampling)
- `--duration`: Total monitoring duration in seconds (default: 0, one-time sampling)
- `--format`: Output format (csv/json/jsonl/xml) (default: json)
- `--output`: Output file path (default: stdout)
- `--pretty`: Pretty print output (for json and xml formats)

### Examples

1. One-time process metrics collection by PID:
```bash
proctrail --pid 1234 --format json
```

2. Collect metrics by process name and save as CSV:
```bash
proctrail --name chrome --format csv --output chrome_metrics.csv
```

3. Sample every 5 seconds for 1 minute with pretty JSON output:
```bash
proctrail --name firefox --interval 5 --duration 60 --format json --pretty
```

4. Collect process metrics and output as XML:
```bash
proctrail --pid 1234 --format xml --pretty --output process_metrics.xml
```

5. Continuous monitoring with JSONL format:
```bash
proctrail --name nginx --interval 10 --format jsonl --output nginx_metrics.jsonl
```

### Output Format Examples

#### JSON
```json
{
  "timestamp": "2025-11-07T15:04:05Z",
  "pid": 1234,
  "name": "example",
  "cpu_percent": 2.5,
  "memory_mb": 150.6,
  "io_read_bytes": 1024,
  "io_write_bytes": 2048,
  "thread_count": 4,
  "fd_count": 8
}
```

#### CSV
```csv
timestamp,pid,name,cpu_percent,memory_mb,io_read_bytes,io_write_bytes,thread_count,fd_count
2025-11-07T15:04:05Z,1234,example,2.5,150.6,1024,2048,4,8
```

#### XML
```xml
<?xml version="1.0" encoding="UTF-8"?>
<metrics>
  <timestamp>2025-11-07T15:04:05Z</timestamp>
  <pid>1234</pid>
  <name>example</name>
  <cpu_percent>2.5</cpu_percent>
  <memory_mb>150.6</memory_mb>
  <io_read_bytes>1024</io_read_bytes>
  <io_write_bytes>2048</io_write_bytes>
  <thread_count>4</thread_count>
  <fd_count>8</fd_count>
</metrics>
```

## Notes

1. Some metrics collection may require specific permissions:
   - Root privileges on Linux/macOS
   - Administrator privileges on Windows

2. Available metrics may vary across operating systems:
   - File descriptor count may not be available on Windows
   - I/O statistics may require special permissions on some systems

3. Performance considerations:
   - Smaller sampling intervals may increase system load
   - Adjust sampling interval based on actual requirements

## Troubleshooting

1. Insufficient Permissions
```bash
Error: failed to collect metrics: permission denied
```
Solution: Run the program with appropriate privileges (sudo or administrator)

2. Process Not Found
```bash
Error: process not found
```
Solution: Verify the PID or process name is correct

3. Output File Access Failed
```bash
Error: failed to create output file
```
Solution: Ensure write permissions for the output path

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
