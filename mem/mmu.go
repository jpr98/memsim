package mem

import (
	"errors"
	"fmt"
	"io"
)

// MMU ...
type MMU struct {
	real     Memory
	swap     Swap
	PageSize int
}

// NewMMU creates a new Memory Management Unit
func NewMMU(memSize, swapSize, pageSize int, policyStr string) (*MMU, error) {
	p, err := createPolicy(policyStr)
	if err != nil {
		return nil, err
	}

	real, err := NewMemory(memSize, pageSize, p)
	if err != nil {
		return nil, err
	}

	swap, err := NewSwap(swapSize, pageSize)
	if err != nil {
		return nil, err
	}

	return &MMU{
		real:     real,
		swap:     swap,
		PageSize: pageSize,
	}, nil
}

// AllocatePage ...
func (m *MMU) AllocatePage(pid string, processPage int) error {
	if ok := m.real.AllocatePage(pid, processPage); ok {
		return nil // page allocated correctly in real memory
	}

	// get page that can be swapped
	p, ok := m.real.NextSwappingCandidate()
	if !ok {
		return errors.New("No swapping candidate found")
	}

	// move page to disk
	err := m.swap.StorePage(p)
	if err != nil {
		return err // error allocating swapped-out page
	}

	// insert process on empty page
	ok = m.real.AllocatePage(pid, processPage)
	if !ok {
		return errors.New("Error allocating swapped-in page")
	}
	return nil
}

// AccessPage ...
func (m *MMU) AccessPage(pid string, address int) (int, error) {
	if addr, ok := m.real.AccessPage(pid, address); ok {
		return addr, nil
	}

	// search for page in swap
	swapInPage, err := m.swap.RetrievePage(pid, address)
	if err != nil {
		return -1, err
	}

	// if present, get pages that can be swapped
	swapOutPage, ok := m.real.NextSwappingCandidate()
	if !ok {
		return -1, errors.New("No swapping candidate found")
	}

	// exchange pages
	err = m.swap.StorePage(swapOutPage)
	if err != nil {
		return -1, err
	}
	if ok := m.real.AllocatePage(swapInPage.pid, swapInPage.virtualAddress); !ok {
		return -1, fmt.Errorf("Error allocating swapped-in page with PID %s", swapInPage.pid)
	}

	// get new address of desired page
	addr, ok := m.real.AccessPage(pid, address)
	if !ok {
		return -1, fmt.Errorf("Error accessing page with PID %s", pid)
	}
	return addr, nil
}

// RemovePages ...
func (m *MMU) RemovePages(pid string) {
	m.real.RemovePages(pid)
	m.swap.RemovePages(pid)
}

// Print ...
func (m *MMU) Print(w io.Writer) {
	pages := m.real.GetPages()
	swappedPages := m.swap.GetPages()
	fmt.Fprintln(w, "RAM\t\tSWAP")
	for i := 0; i < len(swappedPages); i++ {
		if len(pages) > i {
			rp := pages[i]
			fmt.Fprintf(w, "%s\t%d\t", rp.pid, rp.virtualAddress)
		} else {
			fmt.Fprint(w, "\t\t")
		}
		sp := swappedPages[i]
		fmt.Fprintf(w, "%s\t%d\n", sp.pid, sp.virtualAddress)
	}
}
