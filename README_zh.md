# ProcTrail

ProcTrail 是一个跨平台的进程指标采集工具，支持通过 PID 或进程名称来采集目标进程的各项性能指标，并支持多种数据导出格式。

## 功能特点

- 跨平台支持 (Windows, Linux, macOS)
- 多种进程选择器
  - 通过进程 ID (PID)
  - 通过进程名称
- 丰富的指标采集
  - CPU 使用率 (%)
  - 内存使用量 (MB)
  - I/O 读写字节数
  - 线程数量
  - 文件描述符数量
- 灵活的采样选项
  - 一次性采样
  - 定时采样（可配置间隔时间）
  - 采样持续时间控制
- 多种导出格式
  - CSV
  - JSON
  - JSONL (JSON Lines)
  - XML
- 支持美化输出 (JSON/XML)

## 安装

### 从源码安装

确保你已安装 Go 1.16 或更高版本，然后执行：

```bash
git clone https://github.com/lance4117/ProcTrail.git
cd ProcTrail
go install ./cmd/proctrail
```

## 使用方法

### 基本命令格式

```bash
proctrail [flags]
```

### 必需参数（二选一）

- `--pid`: 指定进程 ID
- `--name`: 指定进程名称

### 可选参数

- `--interval`: 采样间隔（秒），默认为 0（一次性采样）
- `--duration`: 总采样时长（秒），默认为 0（一次性采样）
- `--format`: 输出格式（csv/json/jsonl/xml），默认为 json
- `--output`: 输出文件路径，默认输出到标准输出
- `--pretty`: 美化输出（适用于 json 和 xml 格式）

### 使用示例

1. 通过 PID 一次性采集进程指标：
```bash
proctrail --pid 1234 --format json
```

2. 通过进程名称采集并保存为 CSV：
```bash
proctrail --name chrome --format csv --output chrome_metrics.csv
```

3. 每 5 秒采样一次，持续 1 分钟，输出美化的 JSON：
```bash
proctrail --name firefox --interval 5 --duration 60 --format json --pretty
```

4. 采集指定进程的指标并输出为 XML：
```bash
proctrail --pid 1234 --format xml --pretty --output process_metrics.xml
```

5. 使用 JSONL 格式进行连续监控：
```bash
proctrail --name nginx --interval 10 --format jsonl --output nginx_metrics.jsonl
```

### 输出格式示例

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

## 注意事项

1. 某些指标的采集可能需要特定的权限：
   - Linux/macOS 下可能需要 root 权限
   - Windows 下可能需要管理员权限

2. 不同操作系统下的可用指标可能有所不同：
   - 文件描述符计数在 Windows 下可能不可用
   - I/O 统计在某些系统上可能需要特殊权限

3. 性能影响：
   - 较小的采样间隔可能会增加系统负载
   - 建议根据实际需求调整采样间隔

## 故障排除

1. 权限不足
```bash
Error: failed to collect metrics: permission denied
```
解决方案：使用适当的权限运行程序（sudo 或管理员权限）

2. 进程不存在
```bash
Error: process not found
```
解决方案：确认 PID 或进程名称是否正确

3. 输出文件访问失败
```bash
Error: failed to create output file
```
解决方案：确认输出路径是否有写入权限

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

[English](README.md) | 中文
