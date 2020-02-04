package cpu

import (
	"fmt"

	"github.com/jpr98/memsim/mem"
)

// CPU ...
type CPU struct {
	realMemory mem.Memory
	swapMemory mem.Memory
	pids       set
}

type set map[string]bool

// New creates a new CPU
func New(real, swap mem.Memory) CPU {
	return CPU{
		realMemory: real,
		swapMemory: swap,
		pids:       make(map[string]bool),
	}
}

// CreateProcess ...
func (c *CPU) CreateProcess(pid string, size int) error {
	if size == 0 {
		return fmt.Errorf("size should not be 0")
	}

	present := c.pids[pid]
	if present {
		return fmt.Errorf("PID %s is already in cpu", pid)
	}

	requiredPages := size / c.realMemory.PageSize
	for i := 0; i < requiredPages; i++ {
		if ok := c.realMemory.AllocatePage(pid, i); !ok {
			return fmt.Errorf("not enough space in memory for PID: %s", pid)
		}
	}
	c.pids[pid] = true
	// TODO: Check swap
	return nil
}

// AccessProcess ...
func (c *CPU) AccessProcess(pid string, addr int) (int, error) {
	if present := c.pids[pid]; !present {
		return -1, fmt.Errorf("PID %s is not present in cpu", pid)
	}

	rAddr, ok := c.realMemory.AccessPage(pid, addr)
	if !ok {
		return -1, fmt.Errorf("address %d for PID %s not found", addr, pid)
	}
	// TODO: look for process in swap
	return rAddr, nil
}

// DeleteProcess ...
func (c *CPU) DeleteProcess(pid string) error {
	foundInRealMem := c.realMemory.RemovePages(pid)

	// TODO: delte process' pages from swap

	if !foundInRealMem {
		return fmt.Errorf("PID %s not found in memory", pid)
	}
	c.pids[pid] = false
	return nil
}
