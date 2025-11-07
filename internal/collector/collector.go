package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

// ProcessMetrics represents the collected metrics for a single process
type ProcessMetrics struct {
	Timestamp   time.Time `json:"timestamp" xml:"timestamp"`
	PID         int32     `json:"pid" xml:"pid"`
	Name        string    `json:"name" xml:"name"`
	CPUPercent  float64   `json:"cpu_percent" xml:"cpu_percent"`
	MemoryMB    float64   `json:"memory_mb" xml:"memory_mb"`
	IOReadB     uint64    `json:"io_read_bytes" xml:"io_read_bytes"`
	IOWriteB    uint64    `json:"io_write_bytes" xml:"io_write_bytes"`
	ThreadCount int32     `json:"thread_count" xml:"thread_count"`
	FDCount     int32     `json:"fd_count" xml:"fd_count"`
}

// Collector interface defines methods for collecting process metrics
type Collector interface {
	CollectByPID(pid int32) (*ProcessMetrics, error)
	CollectByName(name string) ([]*ProcessMetrics, error)
}

// PSUtilCollector implements Collector using gopsutil
type PSUtilCollector struct{}

// NewPSUtilCollector creates a new PSUtilCollector instance
func NewPSUtilCollector() *PSUtilCollector {
	return &PSUtilCollector{}
}

func (c *PSUtilCollector) CollectByPID(pid int32) (*ProcessMetrics, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("process not found: %w", err)
	}
	return c.collectMetrics(p)
}

func (c *PSUtilCollector) CollectByName(name string) ([]*ProcessMetrics, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to list processes: %w", err)
	}

	var metrics []*ProcessMetrics
	for _, p := range processes {
		procName, err := p.Name()
		if err != nil {
			continue
		}
		if procName == name {
			if m, err := c.collectMetrics(p); err == nil {
				metrics = append(metrics, m)
			}
		}
	}
	return metrics, nil
}

func (c *PSUtilCollector) collectMetrics(p *process.Process) (*ProcessMetrics, error) {
	name, err := p.Name()
	if err != nil {
		return nil, fmt.Errorf("failed to get process name: %w", err)
	}

	cpu, err := p.CPUPercent()
	if err != nil {
		cpu = -1
	}

	memInfo, err := p.MemoryInfo()
	var memMB float64
	if err != nil {
		memMB = -1
	} else {
		memMB = float64(memInfo.RSS) / 1024 / 1024
	}

	ioCounters, err := p.IOCounters()
	var readB, writeB uint64
	if err == nil && ioCounters != nil {
		readB = ioCounters.ReadBytes
		writeB = ioCounters.WriteBytes
	}

	numThreads, err := p.NumThreads()
	if err != nil {
		numThreads = -1
	}

	numFDs, err := p.NumFDs()
	if err != nil {
		numFDs = -1
	}

	return &ProcessMetrics{
		Timestamp:   time.Now(),
		PID:         p.Pid,
		Name:        name,
		CPUPercent:  cpu,
		MemoryMB:    memMB,
		IOReadB:     readB,
		IOWriteB:    writeB,
		ThreadCount: numThreads,
		FDCount:     numFDs,
	}, nil
}
