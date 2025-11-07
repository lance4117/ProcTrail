package exporter

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"

	"github.com/lance4117/ProcTrail/internal/collector"
)

// Exporter defines the interface for different export formats
type Exporter interface {
	Export(metrics []*collector.ProcessMetrics, w io.Writer) error
}

// CSVExporter implements CSV export format
type CSVExporter struct{}

func NewCSVExporter() *CSVExporter {
	return &CSVExporter{}
}

func (e *CSVExporter) WriteHeader(w io.Writer) error {
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	header := []string{"timestamp", "pid", "name", "cpu_percent", "memory_mb", "io_read_bytes", "io_write_bytes", "thread_count", "fd_count"}
	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}
	return nil
}

func (e *CSVExporter) Export(metrics []*collector.ProcessMetrics, w io.Writer) error {
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Write data rows
	for _, m := range metrics {
		row := []string{
			m.Timestamp.Format("2006-01-02 15:04:05"),
			strconv.Itoa(int(m.PID)),
			m.Name,
			strconv.FormatFloat(m.CPUPercent, 'f', 2, 64),
			strconv.FormatFloat(m.MemoryMB, 'f', 2, 64),
			strconv.FormatUint(m.IOReadB, 10),
			strconv.FormatUint(m.IOWriteB, 10),
			strconv.Itoa(int(m.ThreadCount)),
			strconv.Itoa(int(m.FDCount)),
		}
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}
	return nil
}

// JSONExporter implements JSON export format
type JSONExporter struct {
	Pretty bool
}

func NewJSONExporter(pretty bool) *JSONExporter {
	return &JSONExporter{Pretty: pretty}
}

func (e *JSONExporter) Export(metrics []*collector.ProcessMetrics, w io.Writer) error {
	encoder := json.NewEncoder(w)
	if e.Pretty {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(metrics)
}

// JSONLExporter implements JSONL export format
type JSONLExporter struct{}

func NewJSONLExporter() *JSONLExporter {
	return &JSONLExporter{}
}

func (e *JSONLExporter) Export(metrics []*collector.ProcessMetrics, w io.Writer) error {
	encoder := json.NewEncoder(w)
	for _, m := range metrics {
		if err := encoder.Encode(m); err != nil {
			return fmt.Errorf("failed to encode JSONL record: %w", err)
		}
	}
	return nil
}

// XMLExporter implements XML export format
type XMLExporter struct {
	Pretty bool
}

func NewXMLExporter(pretty bool) *XMLExporter {
	return &XMLExporter{Pretty: pretty}
}

type xmlWrapper struct {
	Metrics []*collector.ProcessMetrics `xml:"metrics"`
}

func (e *XMLExporter) Export(metrics []*collector.ProcessMetrics, w io.Writer) error {
	wrapper := xmlWrapper{Metrics: metrics}

	if _, err := io.WriteString(w, xml.Header); err != nil {
		return fmt.Errorf("failed to write XML header: %w", err)
	}

	encoder := xml.NewEncoder(w)
	if e.Pretty {
		encoder.Indent("", "  ")
	}

	return encoder.Encode(wrapper)
}
