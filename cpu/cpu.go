package cpu

import (
	"fmt"
	"io"
	"math"
	"time"

	"github.com/jpr98/memsim/mem"
)

// CPU ...
type CPU struct {
	mmu  mem.MMU
	pids pidMap
}

type process struct {
	live       bool
	start      time.Time
	end        time.Time
	pageFaults int
}

type pidMap map[string]*process

// New creates a new CPU
func New(mmu mem.MMU) CPU {
	return CPU{
		mmu:  mmu,
		pids: make(pidMap),
	}
}

// CreateProcess ...
func (c *CPU) CreateProcess(pid string, size int) error {
	if size == 0 {
		return fmt.Errorf("size should not be 0")
	}

	if _, present := c.pids[pid]; present {
		return fmt.Errorf("PID %s is already in cpu", pid)
	}

	fRequiredPages := float64(size) / float64(c.mmu.PageSize)
	requiredPages := int(math.Ceil(fRequiredPages))
	for i := 0; i < requiredPages; i++ {
		err := c.mmu.AllocatePage(pid, i)
		if err != nil {
			return err
		}
	}

	c.pids[pid] = &process{true, time.Now(), time.Time{}, 0}

	return nil
}

// AccessProcess ...
func (c *CPU) AccessProcess(pid string, addr int) (int, error) {
	if _, present := c.pids[pid]; !present {
		return -1, fmt.Errorf("PID %s is not present in cpu", pid)
	}

	rAddr, pageFault, err := c.mmu.AccessPage(pid, addr)
	if err != nil {
		return -1, err
	}

	if pageFault {
		c.pids[pid].pageFaults++
	}

	return rAddr, nil
}

// DeleteProcess ...
func (c *CPU) DeleteProcess(pid string) error {
	if _, present := c.pids[pid]; !present {
		return fmt.Errorf("process with PID %s not found in memory", pid)
	}

	c.mmu.RemovePages(pid)

	c.pids[pid].end = time.Now()
	c.pids[pid].live = false

	return nil
}

// ReportStats prints to std output the stats of the cpu in the simulation (times in µs)
func (c *CPU) ReportStats() {
	var accumTurnaround time.Duration

	fmt.Println("\n** Stats **")
	fmt.Println("(times in ms)")
	// turnaround time per process
	fmt.Println("\nTurnaround time")
	for k, v := range c.pids {
		fmt.Printf("PID: %s\t-  ", k)
		if v.live {
			v.end = time.Now()
		}
		turnaround := v.end.Sub(v.start)
		accumTurnaround += turnaround
		fmt.Printf("%d\n", turnaround.Microseconds())
	}

	// turnaround promedio
	fmt.Println("\nAverage turnaround time")
	processCount := len(c.pids)
	avgTime := int(accumTurnaround) / processCount
	fmt.Println(time.Duration(avgTime).Microseconds())

	// page faults por proceso
	fmt.Println("\nPage Faults per process")
	for k, v := range c.pids {
		fmt.Printf("PID: %s\t-  ", k)
		fmt.Printf("%d\n", v.pageFaults)
	}

	// numero total de swaps
	fmt.Println("\nTotal swaps")
	fmt.Printf("%d\n", c.mmu.SCount)
	fmt.Print("\n** End of Stats ** \n\n")
}

// Print prints the state of memory if in debug mode
func (c *CPU) Print(debug bool, w io.Writer) {
	if debug {
		c.mmu.Print(w)
	}
}
