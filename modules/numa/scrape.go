package numa

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

			l := strings.TrimLeft(lines[1], " ")
			vals := strings.Split(l, " ")
			numastat.Value, err = strconv.ParseFloat(vals[0], 64)
			if len(vals) > 2 {
				switch vals[1] {
				case "kb":
					numastat.Value = numastat.Value * 1024
				case "mb":
					numastat.Value = numastat.Value * 1024 * 1024
				case "gb":
					numastat.Value = numastat.Value * 1024 * 1024 * 1024
				case "tb":
					numastat.Value = numastat.Value * 1024 * 1024 * 1024 * 1024
				}
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
