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

	memtotal        = "memtotal"
	memfree         = "memfree"
	memused         = "memused"
	swapcached      = "swapcached"
	active          = "active"
	inactive        = "inactive"
	unevictable     = "unevictable"
	mlocked         = "mlocked"
	dirty           = "dirty"
	writeback       = "writeback"
	filepages       = "filepages"
	mapped          = "mapped"
	anonpages       = "anonpages"
	shmem           = "shmem"
	kernelstack     = "kernelstack"
	pagetables      = "pagetables"
	secpagetables   = "secpagetables"
	nfs_unstable    = "nfs_unstable"
	bounce          = "bounce"
	writebacktmp    = "writebacktmp"
	kreclaimable    = "kreclaimable"
	slab            = "slab"
	sreclaimable    = "sreclaimable"
	sunreclaim      = "sunreclaim"
	anonhugepages   = "anonhugepages"
	shmemhugepages  = "shmemhugepages"
	shmempmdmapped  = "shmempmdmapped"
	filehugepages   = "filehugepages"
	filepmdmapped   = "filepmdmapped"
	hugepages_total = "hugepages_total"
	hugepages_free  = "hugepages_free"
	hugepages_surp  = "hugepages_surp"
)

var (
	numaHitDesc       = prometheus.NewDesc("numa_hit", "Number of pages successfully allocated to this node.", []string{"node"}, nil)
	numaMissDesc      = prometheus.NewDesc("numa_miss", "Number of pages allocated to this node due to insufficient memory on the intended node.", []string{"node"}, nil)
	numaForeignDesc   = prometheus.NewDesc("numa_foreign", "This is the number of pages first used for this node that were instead allocated to another node.", []string{"node"}, nil)
	interleaveHitDesc = prometheus.NewDesc("interleave_hit", "Number of temporary policy pages successfully allocated to this node.", []string{"node"}, nil)
	localNodeDesc     = prometheus.NewDesc("local_node", "The number of pages successfully allocated to this node by processes on this node.", []string{"node"}, nil)
	otherNodeDesc     = prometheus.NewDesc("other_node", "Number of pages allocated to this node by processes on other nodes.", []string{"node"}, nil)

	memtotalDesc       = prometheus.NewDesc("memtotal", "Total usable RAM (i.e. physical RAM minus a few reserved bits and the kernel binary code)", []string{"node"}, nil)
	memfreeDesc        = prometheus.NewDesc("memfree", "Total free RAM. On highmem systems, the sum of LowFree+HighFree.", []string{"node"}, nil)
	memusedDesc        = prometheus.NewDesc("memused", "memused", []string{"node"}, nil)
	swapcachedDesc     = prometheus.NewDesc("swapcached", "Memory that once was swapped out, is swapped back in but still also is in the swapfile (if memory is needed it doesn't need to be swapped out AGAIN because it is already in the swapfile. This saves I/O)", []string{"node"}, nil)
	activeDesc         = prometheus.NewDesc("active", "Memory that has been used more recently and usually not reclaimed unless absolutely necessary.", []string{"node", "type"}, nil)
	inactiveDesc       = prometheus.NewDesc("inactive", "Memory which has been less recently used. It is more eligible to be reclaimed for other purposes", []string{"node", "type"}, nil)
	unevictableDesc    = prometheus.NewDesc("unevictable", "Memory allocated for userspace which cannot be reclaimed, such as mlocked pages, ramfs backing pages, secret memfd pages etc.", []string{"node"}, nil)
	mlockedDesc        = prometheus.NewDesc("mlocked", "Memory locked with mlock().", []string{"node"}, nil)
	dirtyDesc          = prometheus.NewDesc("dirty", "Memory which is waiting to get written back to the disk", []string{"node"}, nil)
	writebackDesc      = prometheus.NewDesc("writeback", "Memory which is actively being written back to the disk", []string{"node"}, nil)
	filepagesDesc      = prometheus.NewDesc("filepages", "file backed pages mapped into userspace page tables", []string{"node"}, nil)
	mappedDesc         = prometheus.NewDesc("mapped", "files which have been mmapped, such as libraries", []string{"node"}, nil)
	anonpagesDesc      = prometheus.NewDesc("anonpages", "Non-file backed pages mapped into userspace page tables", []string{"node"}, nil)
	shmemDesc          = prometheus.NewDesc("shmem", "Total memory used by shared memory (shmem) and tmpfs", []string{"node"}, nil)
	kernelstackDesc    = prometheus.NewDesc("kernelstack", "Memory consumed by the kernel stacks of all tasks", []string{"node"}, nil)
	pagetablesDesc     = prometheus.NewDesc("pagetables", "Memory consumed by userspace page tables", []string{"node"}, nil)
	secpagetablesDesc  = prometheus.NewDesc("secpagetables", "Memory consumed by secondary page tables, this currently currently includes KVM mmu allocations on x86 and arm64.", []string{"node"}, nil)
	nfsUnstableDesc    = prometheus.NewDesc("nfs_unstable", "Always zero. Previous counted pages which had been written to the server, but has not been committed to stable storage.", []string{"node"}, nil)
	bounceDesc         = prometheus.NewDesc("bounce", "Memory used for block device 'bounce buffers'", []string{"node"}, nil)
	writebacktmpDesc   = prometheus.NewDesc("writebacktmp", "Memory used by FUSE for temporary writeback buffers", []string{"node"}, nil)
	kreclaimableDesc   = prometheus.NewDesc("kreclaimable", "Kernel allocations that the kernel will attempt to reclaim under memory pressure. ", []string{"node"}, nil)
	slabDesc           = prometheus.NewDesc("slab", "in-kernel data structures cache", []string{"node"}, nil)
	sreclaimableDesc   = prometheus.NewDesc("sreclaimable", "Part of Slab, that might be reclaimed, such as caches", []string{"node"}, nil)
	sunreclaimDesc     = prometheus.NewDesc("sunreclaim", "Part of Slab, that cannot be reclaimed on memory pressure", []string{"node"}, nil)
	anonhugepagesDesc  = prometheus.NewDesc("anonhugepages", "Non-file backed huge pages mapped into userspace page tables", []string{"node"}, nil)
	shmemhugepagesDesc = prometheus.NewDesc("shmemhugepages", "Memory used by shared memory (shmem) and tmpfs allocated with huge pages", []string{"node"}, nil)
	shmempmdmappedDesc = prometheus.NewDesc("shmempmdmapped", "Shared memory mapped into userspace with huge pages", []string{"node"}, nil)
	filehugepagesDesc  = prometheus.NewDesc("filehugepages", "Memory used for filesystem data (page cache) allocated with huge pages", []string{"node"}, nil)
	filepmdmappedDesc  = prometheus.NewDesc("filepmdmapped", "Page cache mapped into userspace with huge pages", []string{"node"}, nil)
	hugepagesTotalDesc = prometheus.NewDesc("hugepages_total", "is the size of the pool of huge pages", []string{"node"}, nil)
	hugepagesFreeDesc  = prometheus.NewDesc("hugepages_free", "is the number of huge pages in the pool that are not yet allocated.", []string{"node"}, nil)
	hugepagesSurpDesc  = prometheus.NewDesc("hugepages_surp", "is short for 'surplus,' and is the number of huge pages in the pool above the value in '/proc/sys/vm/nr_hugepages'", []string{"node"}, nil)
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

		memtotal:        memtotalDesc,
		memfree:         memfreeDesc,
		memused:         memusedDesc,
		swapcached:      swapcachedDesc,
		active:          activeDesc,
		inactive:        inactiveDesc,
		unevictable:     unevictableDesc,
		mlocked:         mlockedDesc,
		dirty:           dirtyDesc,
		writeback:       writebackDesc,
		filepages:       filepagesDesc,
		mapped:          mappedDesc,
		anonpages:       anonpagesDesc,
		shmem:           shmemDesc,
		kernelstack:     kernelstackDesc,
		pagetables:      pagetablesDesc,
		secpagetables:   secpagetablesDesc,
		nfs_unstable:    nfsUnstableDesc,
		bounce:          bounceDesc,
		writebacktmp:    writebacktmpDesc,
		kreclaimable:    kreclaimableDesc,
		slab:            slabDesc,
		sreclaimable:    sreclaimableDesc,
		sunreclaim:      sunreclaimDesc,
		anonhugepages:   anonhugepagesDesc,
		shmemhugepages:  shmemhugepagesDesc,
		shmempmdmapped:  shmempmdmappedDesc,
		filehugepages:   filehugepagesDesc,
		filepmdmapped:   filepmdmappedDesc,
		hugepages_total: hugepagesTotalDesc,
		hugepages_free:  hugepagesFreeDesc,
		hugepages_surp:  hugepagesSurpDesc,
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
		if stat.Type != nil {
			ch <- prometheus.MustNewConstMetric(sourcesDesc[stat.Name],
				prometheus.CounterValue,
				stat.Value,
				strconv.Itoa(stat.Node),
				*stat.Type,
			)
		} else {
			ch <- prometheus.MustNewConstMetric(sourcesDesc[stat.Name],
				prometheus.CounterValue,
				stat.Value,
				strconv.Itoa(stat.Node),
			)
		}
	}
}
