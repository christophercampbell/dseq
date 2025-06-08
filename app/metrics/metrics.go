package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Transaction metrics
	TransactionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dseq_transactions_total",
			Help: "Total number of transactions processed",
		},
		[]string{"status"},
	)

	TransactionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dseq_transaction_duration_seconds",
			Help:    "Transaction processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	// Block metrics
	BlockHeight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "dseq_block_height",
			Help: "Current block height",
		},
	)

	BlockProcessingDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "dseq_block_processing_duration_seconds",
			Help:    "Block processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	// Stream metrics
	StreamEntriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dseq_stream_entries_total",
			Help: "Total number of stream entries processed",
		},
		[]string{"type"},
	)

	StreamLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dseq_stream_latency_seconds",
			Help:    "Stream entry processing latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"type"},
	)

	// System metrics
	MemoryUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "dseq_memory_usage_bytes",
			Help: "Current memory usage in bytes",
		},
	)

	GoroutineCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "dseq_goroutine_count",
			Help: "Current number of goroutines",
		},
	)
)

// RecordTransaction records transaction metrics
func RecordTransaction(status string, duration float64) {
	TransactionTotal.WithLabelValues(status).Inc()
	TransactionDuration.WithLabelValues(status).Observe(duration)
}

// RecordBlock records block metrics
func RecordBlock(height int64, duration float64) {
	BlockHeight.Set(float64(height))
	BlockProcessingDuration.Observe(duration)
}

// RecordStreamEntry records stream entry metrics
func RecordStreamEntry(entryType string, latency float64) {
	StreamEntriesTotal.WithLabelValues(entryType).Inc()
	StreamLatency.WithLabelValues(entryType).Observe(latency)
}

// RecordSystemMetrics records system metrics
func RecordSystemMetrics(memoryBytes int64, goroutineCount int) {
	MemoryUsage.Set(float64(memoryBytes))
	GoroutineCount.Set(float64(goroutineCount))
}
