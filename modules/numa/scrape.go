package numa

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
	0, 0-1
*/

/*
   numa_hit 6422219
   numa_miss 0
   numa_foreign 0
   interleave_hit 6460
   local_node 6417749
   other_node 4470
*/

/*
Node 0 MemTotal:        7856504 kB
Node 0 MemFree:         1400616 kB
Node 0 MemUsed:         6455888 kB
Node 0 SwapCached:            0 kB
Node 0 Active:          1238936 kB
Node 0 Inactive:        4570064 kB
Node 0 Active(anon):      11688 kB
Node 0 Inactive(anon):  1871708 kB
Node 0 Active(file):    1227248 kB
Node 0 Inactive(file):  2698356 kB
Node 0 Unevictable:          64 kB
Node 0 Mlocked:              64 kB
Node 0 Dirty:                 0 kB
Node 0 Writeback:             0 kB
Node 0 FilePages:       3958580 kB
Node 0 Mapped:           724772 kB
Node 0 AnonPages:       1765740 kB
Node 0 Shmem:             32976 kB
Node 0 KernelStack:       10640 kB
Node 0 PageTables:        19588 kB
Node 0 SecPageTables:         0 kB
Node 0 NFS_Unstable:          0 kB
Node 0 Bounce:                0 kB
Node 0 WritebackTmp:          0 kB
Node 0 KReclaimable:     251456 kB
Node 0 Slab:             398876 kB
Node 0 SReclaimable:     251456 kB
Node 0 SUnreclaim:       147420 kB
Node 0 AnonHugePages:   1292288 kB
Node 0 ShmemHugePages:        0 kB
Node 0 ShmemPmdMapped:        0 kB
Node 0 FileHugePages:        0 kB
Node 0 FilePmdMapped:        0 kB
Node 0 HugePages_Total:     0
Node 0 HugePages_Free:      0
Node 0 HugePages_Surp:      0
*/
const (
	NODE_ONLINE     = "/sys/devices/system/node/online"
	NUMASTAT_SCRAPE = "/sys/devices/system/node/%s/numastat"
	MEMINFO_SCRAPE  = "/sys/devices/system/node/%s/meminfo"
)

type Numastats []Numastat

type Numastat struct {
	Node  int
	Name  string
	Value float64
	Type  *string
}

func Scrape() (Numastats, error) {
	var items Numastats
	if err := foundFile(NODE_ONLINE); err != nil {
		return nil, err
	}

	cnt, err := os.ReadFile(NODE_ONLINE)
	if err != nil {
		return nil, err
	}

	values := strings.Split(strings.ReplaceAll(string(cnt), "\n", ""), "-")
	lastNum, err := strconv.Atoi(values[len(values)-1])
	if err != nil {
		return nil, err
	}
	for i := 0; i <= lastNum; i++ {
		nodeName := fmt.Sprintf("node%d", i)

		if contents, err := fileContent(fmt.Sprintf(NUMASTAT_SCRAPE, nodeName), func(line string) *Numastat {
			lineValue := strings.Split(line, " ")
			val, err := strconv.ParseFloat(lineValue[1], 64)
			if err != nil {
				return nil
			}
			return &Numastat{
				Node:  i,
				Name:  lineValue[0],
				Value: val,
			}
		}); err != nil {
			return nil, err
		} else {
			items = append(items, contents...)
		}

		if contents, err := fileContent(fmt.Sprintf(MEMINFO_SCRAPE, nodeName), func(line string) *Numastat {
			lines := strings.Split(line, ":")

			valuekey := strings.ToLower(strings.Split(lines[0], " ")[2])

			numastat := Numastat{
				Node: i,
				Name: valuekey,
			}

			if strings.HasPrefix(valuekey, "active") {
				var typeValue string
				if strings.HasSuffix(valuekey, "(anon)") {
					typeValue = "anon"
				} else if strings.HasSuffix(valuekey, "(file)") {
					typeValue = "file"
				} else {
					typeValue = "default"
				}
				numastat.Name = "active"
				numastat.Type = &typeValue
			} else if strings.HasPrefix(valuekey, "inactive") {
				var typeValue string
				if strings.HasSuffix(valuekey, "(anon)") {
					typeValue = "anon"
				} else if strings.HasSuffix(valuekey, "(file)") {
					typeValue = "file"
				} else {
					typeValue = "default"
				}
				numastat.Name = "inactive"
				numastat.Type = &typeValue
			}

			val, err := strconv.ParseFloat(strings.Split(lines[1], " ")[0], 64)
			if err == nil {
				numastat.Value = val
			}

			return &numastat

		}); err != nil {
			return nil, err
		} else {
			items = append(items, contents...)
		}

	}
	return items, nil
}

func fileContent(filepath string, parse func(string) *Numastat) (Numastats, error) {
	if err := foundFile(filepath); err != nil {
		return nil, err
	}
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var items Numastats
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		item := parse(scanner.Text())
		if item != nil {
			items = append(items, *item)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func foundFile(filepath string) error {
	fsInfo, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	if fsInfo.IsDir() {
		return errors.New("file not found, found directory")
	}
	return nil
}
