package cpu

import (
	"fmt"
	"io"
	"math"

	"github.com/jpr98/memsim/mem"
)

// CPU ...
type CPU struct {
	mmu  mem.MMU
	pids set
}

type set map[string]bool

// New creates a new CPU
func New(mmu mem.MMU) CPU {
	return CPU{
		mmu:  mmu,
		pids: make(set),
	}
}

// CreateProcess ...
func (c *CPU) CreateProcess(pid string, size int) error {
	if size == 0 {
		return fmt.Errorf("size should not be 0")
	}

	if c.pids[pid] {
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

	c.pids[pid] = true
	return nil
}

// AccessProcess ...
func (c *CPU) AccessProcess(pid string, addr int) (int, error) {
	if !c.pids[pid] {
		return -1, fmt.Errorf("PID %s is not present in cpu", pid)
	}

	rAddr, err := c.mmu.AccessPage(pid, addr)
	if err != nil {
		return -1, err
	}

	return rAddr, nil
}

// DeleteProcess ...
func (c *CPU) DeleteProcess(pid string) error {
	if !c.pids[pid] {
		return fmt.Errorf("process with PID %s not found in memory", pid)
	}
	c.mmu.RemovePages(pid)
	c.pids[pid] = false
	return nil
}

// Print prints the state of memory if in debug mode
func (c *CPU) Print(debug bool, w io.Writer) {
	if debug {
		c.mmu.Print(w)
	}
}
