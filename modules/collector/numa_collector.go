package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	promVersion "github.com/prometheus/common/version"
	"github.com/wunicorns/numa_exporter/modules/numa"
)

const (
	NUMA_Hit_Measurement       = "numa_hit"
	NUMA_Miss_Measurement      = "numa_miss"
	NUMA_Foreign_Measurement   = "numa_foreign"
	Interleave_Hit_Measurement = "interleave_hit"
	Local_Node_Measurement     = "local_node"
	Other_Node_Measurement     = "other_node"
)

var (
	numaHitDesc       = prometheus.NewDesc("numa_hit", "Number of pages successfully allocated to this node.", []string{"node"}, nil)
	numaMissDesc      = prometheus.NewDesc("numa_miss", "Number of pages allocated to this node due to insufficient memory on the intended node.", []string{"node"}, nil)
	numaForeignDesc   = prometheus.NewDesc("numa_foreign", "This is the number of pages first used for this node that were instead allocated to another node.", []string{"node"}, nil)
	interleaveHitDesc = prometheus.NewDesc("interleave_hit", "Number of temporary policy pages successfully allocated to this node.", []string{"node"}, nil)
	localNodeDesc     = prometheus.NewDesc("local_node", "The number of pages successfully allocated to this node by processes on this node.", []string{"node"}, nil)
	otherNodeDesc     = prometheus.NewDesc("other_node", "Number of pages allocated to this node by processes on other nodes.", []string{"node"}, nil)
)

type (
	scrapeRequest struct {
		results chan<- prometheus.Metric
		done    chan struct{}
	}

	ScrapeSourceDesc map[string]*prometheus.Desc
	ScrapeSource     func()

	NumastatCollector struct {
		scrapeChan  chan scrapeRequest
		Source      func() interface{}
		SourceDesc  ScrapeSourceDesc
		LabelValues []string
	}
)

var version string

func init() {
	promVersion.Version = version
	prometheus.MustRegister(promVersion.NewCollector("numa_exporter"))
}

func NewNumaCollector(source func() interface{}) (*NumastatCollector, error) {

	var sourceDesc = ScrapeSourceDesc{
		NUMA_Hit_Measurement:       numaHitDesc,
		NUMA_Miss_Measurement:      numaMissDesc,
		NUMA_Foreign_Measurement:   numaForeignDesc,
		Interleave_Hit_Measurement: interleaveHitDesc,
		Local_Node_Measurement:     localNodeDesc,
		Other_Node_Measurement:     otherNodeDesc,
	}

	p := &NumastatCollector{
		scrapeChan: make(chan scrapeRequest),
		Source:     source,
		SourceDesc: sourceDesc,
	}

	go p.start()

	return p, nil
}

func (p *NumastatCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range p.SourceDesc {
		ch <- desc
	}
}

func (p *NumastatCollector) Collect(ch chan<- prometheus.Metric) {
	req := scrapeRequest{results: ch, done: make(chan struct{})}
	p.scrapeChan <- req
	<-req.done
}

func (p *NumastatCollector) start() {
	for req := range p.scrapeChan {
		ch := req.results
		p.scrape(ch)
		req.done <- struct{}{}
	}
}

func (p *NumastatCollector) scrape(ch chan<- prometheus.Metric) {
	sources := p.Source().(*numa.Numastats)
	sourcesDesc := p.SourceDesc
	for _, stat := range *sources {
		ch <- prometheus.MustNewConstMetric(sourcesDesc[stat.Name],
			prometheus.CounterValue,
			stat.Value,
			strconv.Itoa(stat.Node),
		)
	}
}
