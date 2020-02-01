package cpu

import (
	"fmt"

	"github.com/jpr98/memsim/mem"
)

// CPU ...
type CPU struct {
	realMemory mem.Memory
	swapMemory mem.Memory
}

// New creates a new CPU
func New(real, swap mem.Memory) CPU {
	return CPU{
		realMemory: real,
		swapMemory: swap,
	}
}

// CreateProcess ...
func (c *CPU) CreateProcess(pid string, size int) error {
	requiredPages := size / c.realMemory.PageSize
	for i := 0; i < requiredPages; i++ {
		if ok := c.realMemory.AllocatePage(pid, i); !ok {
			return fmt.Errorf("not enough space in memory for PID: %s", pid)
		}
	}
	return nil
}

// AccessProcess ...
func (c *CPU) AccessProcess(pid string, addr int) (*int, error) {
	rAddr, ok := c.realMemory.AccessPage(pid, addr)
	if !ok {
		return nil, fmt.Errorf("address %d for PID %s not found", addr, pid)
	}

	return &rAddr, nil
}
