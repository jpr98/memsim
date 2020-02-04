package mem

import "errors"

type page struct {
	pid            string
	virtualAddress int
}

// Memory ...
type Memory struct {
	freeList []int
	pages    []page
	PageSize int
}

// New creates a new Memory
func New(size, pageSize int) (*Memory, error) {
	if pageSize == 0 {
		return nil, errors.New("PageSize should not be zero")
	}
	numOfPages := size / pageSize
	return &Memory{
		freeList: createFreeList(numOfPages),
		pages:    make([]page, numOfPages),
		PageSize: pageSize,
	}, nil
}

func createFreeList(size int) []int {
	list := make([]int, size)
	for i := 0; i < size; i++ {
		list[i] = i
	}
	return list
}

func (m *Memory) getNextFreeAddress() (int, bool) {
	if len(m.freeList) < 1 {
		return -1, false
	}
	addr := m.freeList[0]
	m.freeList = m.freeList[1:]
	return addr, true
}

// AllocatePage ...
func (m *Memory) AllocatePage(pid string, processPage int) bool {
	addr, ok := m.getNextFreeAddress()
	if !ok {
		return false
	}

	if m.pages[addr].pid != "" {
		return false
	}

	m.pages[addr] = page{pid, processPage}
	return true
}

// AccessPage ...
func (m *Memory) AccessPage(pid string, address int) (int, bool) {
	displacedAddress := address / m.PageSize
	for realAddress, page := range m.pages {
		if page.pid == pid && page.virtualAddress == displacedAddress {
			return realAddress, true
		}
	}
	return -1, false
}

// RemovePages ...
func (m *Memory) RemovePages(pid string) bool {
	found := false
	for i, p := range m.pages {
		if p.pid == pid {
			m.pages[i] = page{}
			found = true
			m.freeList = append(m.freeList, i)
		}
	}
	return found
}
