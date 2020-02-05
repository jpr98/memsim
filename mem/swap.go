package mem

import (
	"fmt"
)

// Swap represents the functionality of swap space in disk
type Swap interface {
	RetrievePage(string, int) (page, error)
	StorePage(page) error
	RemovePages(string)
	GetPages() []page
}

// NewSwap creates a new Swap interface
func NewSwap(size, pageSize int) (Swap, error) {
	return new(size, pageSize, nil)
}

func (m *memory) RetrievePage(pid string, address int) (page, error) {
	displacedAddress := address / m.PageSize
	for i, p := range m.pages {
		if p.pid == pid && p.virtualAddress == displacedAddress {
			m.freeAddress(i)
			return p, nil
		}
	}
	return page{}, fmt.Errorf("page with virtual address %d for PID %s not found", address, pid)
}

func (m *memory) StorePage(p page) error {
	addr, ok := m.getNextFreeAddress()
	if !ok {
		return fmt.Errorf("no space available in swap")
	}

	if m.pages[addr].pid != "" {
		return fmt.Errorf("error mapping swap address")
	}

	m.pages[addr] = p
	return nil
}
