# numa_exporter

## run 

```
go run main.go
```

```
# cat /sys/devices/system/node/node0/numastat
numa_hit 11522122
numa_miss 0
numa_foreign 0
interleave_hit 6460
local_node 11514479
other_node 7643
```

```
# cat /sys/devices/system/node/node0/meminfo
Node 0 MemTotal:       65165920 kB
Node 0 MemFree:        61730760 kB
Node 0 MemUsed:         3435160 kB
Node 0 SwapCached:            0 kB
Node 0 Active:           950520 kB
Node 0 Inactive:        1683088 kB
Node 0 Active(anon):       4020 kB
Node 0 Inactive(anon):   572268 kB
Node 0 Active(file):     946500 kB
Node 0 Inactive(file):  1110820 kB
Node 0 Unevictable:          16 kB
Node 0 Mlocked:               0 kB
Node 0 Dirty:                 0 kB
Node 0 Writeback:             0 kB
Node 0 FilePages:       2090864 kB
Node 0 Mapped:           396496 kB
Node 0 AnonPages:        491680 kB
Node 0 Shmem:             34508 kB
Node 0 KernelStack:       11240 kB
Node 0 PageTables:         8032 kB
Node 0 SecPageTables:         0 kB
Node 0 NFS_Unstable:          0 kB
Node 0 Bounce:                0 kB
Node 0 WritebackTmp:          0 kB
Node 0 KReclaimable:     323940 kB
Node 0 Slab:             526856 kB
Node 0 SReclaimable:     323940 kB
Node 0 SUnreclaim:       202916 kB
Node 0 AnonHugePages:    296960 kB
Node 0 ShmemHugePages:        0 kB
Node 0 ShmemPmdMapped:        0 kB
Node 0 FileHugePages:        0 kB
Node 0 FilePmdMapped:        0 kB
Node 0 HugePages_Total:     0
Node 0 HugePages_Free:      0
Node 0 HugePages_Surp:      0
```