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
		filepath := fmt.Sprintf(NUMASTAT_SCRAPE, nodeName)
		contents, err := fileContent(filepath)
		if err != nil {
			continue
		}
		for _, line := range contents {
			lineValue := strings.Split(line, " ")
			if val, err := strconv.ParseFloat(lineValue[1], 64); err != nil {
				return nil, err
			} else {
				items = append(items, Numastat{
					Node:  i,
					Name:  lineValue[0],
					Value: val,
				})
			}
		}

	}
	return items, nil
}

func fileContent(filepath string) ([]string, error) {
	if err := foundFile(filepath); err != nil {
		return nil, err
	}
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
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
