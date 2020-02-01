package mem

type page struct {
	pid            string
	virtualAddress int
}

// Memory ...
type Memory struct {
	pages    []page
	PageSize int
}

// New creates a new Memory
func New(size, pageSize int) Memory {
	numOfPages := size / pageSize
	return Memory{
		pages:    make([]page, numOfPages),
		PageSize: pageSize,
	}
}

// AllocatePage ...
func (m *Memory) AllocatePage(pid string, processPage int) bool {
	for i, p := range m.pages {
		if p.pid == "" {
			m.pages[i] = page{pid, processPage}
			return true
		}
	}
	return false
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
