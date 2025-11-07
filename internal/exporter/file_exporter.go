package exporter

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/lance4117/ProcTrail/internal/collector"
)

// SingleFileExporter handles file export operations with a single file for all records
type SingleFileExporter struct {
	file       *os.File
	exporter   Exporter
	mu         sync.Mutex
	headerSent bool // For CSV format to track if header has been written
}

// NewSingleFileExporter creates a new file exporter that writes all records to a single file
func NewSingleFileExporter(outputPath string, exp Exporter) (*SingleFileExporter, error) {
	var filePath string
	if outputPath == "" || filepath.Ext(outputPath) == "" {
		// If no path specified or it's a directory, create a file with timestamp
		timestamp := time.Now().Format("20060102_150405")
		fileName := fmt.Sprintf("proctrail_%s.%s", timestamp, getFormatExtension(exp))
		if outputPath == "" {
			outputPath = "."
		}
		filePath = filepath.Join(outputPath, fileName)
	} else {
		filePath = outputPath
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Open file in append mode
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	return &SingleFileExporter{
		file:     file,
		exporter: exp,
	}, nil
}

// Export writes metrics to the file
func (e *SingleFileExporter) Export(metrics []*collector.ProcessMetrics) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// For CSV format, we need special handling of headers
	if csvExp, ok := e.exporter.(*CSVExporter); ok && !e.headerSent {
		// Write header for CSV only once
		if err := csvExp.WriteHeader(e.file); err != nil {
			return fmt.Errorf("failed to write CSV header: %w", err)
		}
		e.headerSent = true
	}

	if err := e.exporter.Export(metrics, e.file); err != nil {
		return fmt.Errorf("failed to export metrics: %w", err)
	}

	return nil
}

// Close closes the file
func (e *SingleFileExporter) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.file != nil {
		if err := e.file.Close(); err != nil {
			return fmt.Errorf("failed to close output file: %w", err)
		}
		e.file = nil
	}
	return nil
}

// Helper function to get file extension based on exporter type
func getFormatExtension(e Exporter) string {
	switch e.(type) {
	case *CSVExporter:
		return "csv"
	case *JSONExporter:
		return "json"
	case *JSONLExporter:
		return "jsonl"
	case *XMLExporter:
		return "xml"
	default:
		return "txt"
	}
}
