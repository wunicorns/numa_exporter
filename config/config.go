package config

var (
	PROCESS_NAME = "numa_exporter"

	// 포트
	HOST      = ""
	URL       = "/metrics"
	POOL_SIZE = 100
	PORT      = 9999

	LOG_DIR         = "/tmp/"
	LOG_MAX_SIZE    = 500
	LOG_MAX_BACKUPS = 50
	LOG_MAX_AGE     = 30
	LOG_COMPRESS    = false
	DEBUG_LEVEL     = "debug"

	TYPE        = ""
	MEASUREMENT = ""

	METRICS_PATH = "/metrics"
	USE_METRICS  = true
)
