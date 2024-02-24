MemTotal
              Total usable RAM (i.e. physical RAM minus a few reserved bits and the kernel binary code)
MemFree
              Total free RAM. On highmem systems, the sum of LowFree+HighFree
MemAvailable
              An estimate of how much memory is available for starting new applications, without swapping. 
			  Calculated from MemFree, SReclaimable, the size of the file LRU lists, and the low watermarks in each zone.
              The estimate takes into account that the system needs some page cache to function well, and that not all reclaimable slab will be reclaimable, due to items being in use. 
			  The impact of those factors will vary from system to system.	  
Buffers
              Relatively temporary storage for raw disk blocks shouldn't get tremendously large (20MB or so)
Cached
              In-memory cache for files read from the disk (the pagecache) as well as tmpfs & shmem. Doesn't include SwapCached. 
SwapCached
              Memory that once was swapped out, is swapped back in but still also is in the swapfile (if memory is needed it doesn't need to be swapped out AGAIN because it is already in the swapfile. This saves I/O)
Active
              Memory that has been used more recently and usually not reclaimed unless absolutely necessary.
Inactive
              Memory which has been less recently used. It is more eligible to be reclaimed for other purposes
Unevictable
              Memory allocated for userspace which cannot be reclaimed, such as mlocked pages, ramfs backing pages, secret memfd pages etc.
Mlocked
              Memory locked with mlock().
HighTotal, HighFree
              Highmem is all memory above ~860MB of physical memory.
              Highmem areas are for use by userspace programs, or for the pagecache.  The kernel must use tricks to access this memory, making it slower to access than lowmem.
LowTotal, LowFree
              Lowmem is memory which can be used for everything that highmem can be used for, but it is also available for the kernel's use for its own data structures.  
			  Among many other things, it is where everything from the Slab is allocated.  Bad things happen when you're out of lowmem.
SwapTotal
              total amount of swap space available
SwapFree
              Memory which has been evicted from RAM, and is temporarily on the disk
Zswap
              Memory consumed by the zswap backend (compressed size)
Zswapped
              Amount of anonymous memory stored in zswap (original size)
Dirty
              Memory which is waiting to get written back to the disk
Writeback
              Memory which is actively being written back to the disk
AnonPages
              Non-file backed pages mapped into userspace page tables
Mapped
              files which have been mmapped, such as libraries
Shmem
              Total memory used by shared memory (shmem) and tmpfs
KReclaimable
              Kernel allocations that the kernel will attempt to reclaim under memory pressure. 
			  Includes SReclaimable (below), and other direct allocations with a shrinker.
Slab
              in-kernel data structures cache
SReclaimable
              Part of Slab, that might be reclaimed, such as caches
SUnreclaim
              Part of Slab, that cannot be reclaimed on memory pressure
KernelStack
              Memory consumed by the kernel stacks of all tasks
PageTables
              Memory consumed by userspace page tables
SecPageTables
              Memory consumed by secondary page tables, this currently currently includes KVM mmu allocations on x86 and arm64.
NFS_Unstable
              Always zero. Previous counted pages which had been written to the server, but has not been committed to stable storage.
Bounce
              Memory used for block device "bounce buffers"
WritebackTmp
              Memory used by FUSE for temporary writeback buffers
CommitLimit
              Based on the overcommit ratio ('vm.overcommit_ratio'), this is the total amount of  memory currently available to be allocated on the system. 
			  This limit is only adhered to if strict overcommit accounting is enabled (mode 2 in 'vm.overcommit_memory').

              The CommitLimit is calculated with the following formula::

                CommitLimit = ([total RAM pages] - [total huge TLB pages]) *
                               overcommit_ratio / 100 + [total swap pages]

              For example, on a system with 1G of physical RAM and 7G
              of swap with a `vm.overcommit_ratio` of 30 it would
              yield a CommitLimit of 7.3G.

              For more details, see the memory overcommit documentation in mm/overcommit-accounting.
Committed_AS
              The amount of memory presently allocated on the system.
              The committed memory is a sum of all of the memory which has been allocated by processes, even if it has not been "used" by them as of yet. 
			  A process which malloc()'s 1G of memory, but only touches 300M of it will show up as using 1G. 
			  This 1G is memory which has been "committed" to by the VM and can be used at any time by the allocating application. 
			  With strict overcommit enabled on the system (mode 2 in 'vm.overcommit_memory'), allocations which would exceed the CommitLimit (detailed above) will not be permitted.
              This is useful if one needs to guarantee that processes will not fail due to lack of memory once that memory has been successfully allocated.
VmallocTotal
              total size of vmalloc virtual address space
VmallocUsed
              amount of vmalloc area which is used
VmallocChunk
              largest contiguous block of vmalloc area which is free
Percpu
              Memory allocated to the percpu allocator used to back percpu allocations. This stat excludes the cost of metadata.
EarlyMemtestBad
              The amount of RAM/memory in kB, that was identified as corrupted by early memtest. 
			  If memtest was not run, this field will not be displayed at all. Size is never rounded down to 0 kB.
              That means if 0 kB is reported, you can safely assume there was at least one pass of memtest and none of the passes found a single faulty byte of RAM.
HardwareCorrupted
              The amount of RAM/memory in KB, the kernel identifies as corrupted.
AnonHugePages
              Non-file backed huge pages mapped into userspace page tables
ShmemHugePages
              Memory used by shared memory (shmem) and tmpfs allocated with huge pages
ShmemPmdMapped
              Shared memory mapped into userspace with huge pages
FileHugePages
              Memory used for filesystem data (page cache) allocated with huge pages
FilePmdMapped
              Page cache mapped into userspace with huge pages
CmaTotal
              Memory reserved for the Contiguous Memory Allocator (CMA)
CmaFree
              Free remaining memory in the CMA reserves
HugePages_Total, HugePages_Free, HugePages_Rsvd, HugePages_Surp, Hugepagesize, Hugetlb
              See Documentation/admin-guide/mm/hugetlbpage.rst.
DirectMap4k, DirectMap2M, DirectMap1G
              Breakdown of page table sizes used in the kernel's identity mapping of RAM