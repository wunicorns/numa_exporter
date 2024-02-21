package main

import (
	"runtime"

	"github.com/wunicorns/numa_exporter/middleware"
	_ "github.com/wunicorns/numa_exporter/middleware"
	"github.com/wunicorns/numa_exporter/modules/log"

	"github.com/wunicorns/numa_exporter/config"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := log.Init(log.Config{
		LogDir:        config.LOG_DIR,
		LogMaxSize:    config.LOG_MAX_SIZE,
		LogMaxBackups: config.LOG_MAX_BACKUPS,
		LogMaxAge:     config.LOG_MAX_AGE,
		LogCompress:   config.LOG_COMPRESS,
		DebugLevel:    config.DEBUG_LEVEL,
	}); err != nil {
		panic(err)
	}

	middleware.Serve()

}
