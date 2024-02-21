package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wunicorns/numa_exporter/config"
	"github.com/wunicorns/numa_exporter/modules/collector"
	"github.com/wunicorns/numa_exporter/modules/log"
	"github.com/wunicorns/numa_exporter/modules/numa"
)

type Scraper struct {
	*numa.Numastats
}

func (s *Scraper) collect() error {
	stats, err := numa.Scrape()
	if err != nil {
		return err
	}
	s.Numastats = &stats
	return nil
}

func (s *Scraper) source() interface{} {
	s.collect()
	return s.Numastats
}

func Serve() {

	scraper := Scraper{}

	nc, err := collector.NewNumaCollector(scraper.source)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	prometheus.MustRegister(nc)

	http.Handle(config.METRICS_PATH, promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>NUMA Exporter</title></head>
			<body>
			<h1>NUMA Exporter</h1>
			<p><a href="` + config.METRICS_PATH + `">Metrics</a></p>
			</body>
			</html>`))
	})

	serve := fmt.Sprintf("%s:%d", config.HOST, config.PORT)

	log.Debug("Running ...", serve)

	if err := http.ListenAndServe(serve, nil); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
