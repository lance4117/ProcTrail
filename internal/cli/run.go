package cli

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/lance4117/ProcTrail/internal/collector"
	"github.com/lance4117/ProcTrail/internal/exporter"
	"github.com/spf13/cobra"
)

// SpinnerAnimation 定义动画帧
var SpinnerAnimation = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner 结构体用于控制动画
type Spinner struct {
	stop chan struct{}
}

// NewSpinner 创建新的进度指示器
func NewSpinner() *Spinner {
	return &Spinner{
		stop: make(chan struct{}),
	}
}

// Start 启动进度指示器
func (s *Spinner) Start() {
	green := color.New(color.FgGreen).SprintFunc()
	i := 0
	go func() {
		for {
			select {
			case <-s.stop:
				fmt.Print("\r\033[K") // 清除当前行
				return
			default:
				fmt.Printf("\r%s Monitoring process... ", green(SpinnerAnimation[i]))
				i = (i + 1) % len(SpinnerAnimation)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// Stop 停止进度指示器
func (s *Spinner) Stop() {
	close(s.stop)
}

func runCommand(_ *cobra.Command, _ []string) error {
	if err := validateFlags(); err != nil {
		return err
	}

	// Create collector
	c := collector.NewPSUtilCollector()

	// Create file exporter
	fileExporter, err := createExporter()
	if err != nil {
		return err
	}
	defer func() {
		if err := fileExporter.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing output file: %v\n", err)
		}
	}()

	// 创建并启动进度指示器
	spinner := NewSpinner()
	spinner.Start()
	defer spinner.Stop()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create ticker if interval is specified
	var ticker *time.Ticker
	if interval > 0 {
		ticker = time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()
	}

	// Set up duration timer if specified
	var timer *time.Timer
	if duration > 0 {
		timer = time.NewTimer(getDuration())
		defer timer.Stop()
	}

	// Collect and export metrics
	for {
		metrics, err := collectMetrics(c)
		if err != nil {
			return err
		}

		if err := fileExporter.Export(metrics); err != nil {
			return fmt.Errorf("failed to export metrics: %w", err)
		}

		if interval == 0 {
			return nil
		}

		select {
		case <-sigChan:
			fmt.Println("\nMonitoring interrupted")
			return nil
		case <-ticker.C:
			continue
		case <-timer.C:
			fmt.Println("\nMonitoring completed")
			return nil
		}
	}
}

func collectMetrics(c collector.Collector) ([]*collector.ProcessMetrics, error) {
	if pid != 0 {
		metric, err := c.CollectByPID(pid)
		if err != nil {
			return nil, fmt.Errorf("failed to collect metrics for PID %d: %w", pid, err)
		}
		return []*collector.ProcessMetrics{metric}, nil
	}

	metrics, err := c.CollectByName(processName)
	if err != nil {
		return nil, fmt.Errorf("failed to collect metrics for process '%s': %w", processName, err)
	}
	return metrics, nil
}

func createExporter() (*exporter.SingleFileExporter, error) {
	var exp exporter.Exporter
	switch strings.ToLower(outputFormat) {
	case "csv":
		exp = exporter.NewCSVExporter()
	case "json":
		exp = exporter.NewJSONExporter(prettyPrint)
	case "jsonl":
		exp = exporter.NewJSONLExporter()
	case "xml":
		exp = exporter.NewXMLExporter(prettyPrint)
	default:
		return nil, fmt.Errorf("unsupported output format: %s", outputFormat)
	}

	// 创建单文件导出器
	fileExporter, err := exporter.NewSingleFileExporter(outputFile, exp)
	if err != nil {
		return nil, fmt.Errorf("failed to create file exporter: %w", err)
	}

	return fileExporter, nil
}
