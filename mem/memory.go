package mem

// Memory represents the real memory
type Memory interface {
	AllocatePage(string, int) bool
	AccessPage(string, int) (int, bool)
	RemovePages(string)
	NextSwappingCandidate() (page, bool)
	GetPages() []page
}

// NewMemory creates a new Memory interface
func NewMemory(size, pageSize int) (Memory, error) {
	return new(size, pageSize)
}

// AllocatePage ...
func (m *memory) AllocatePage(pid string, processPage int) bool {
	addr, ok := m.getNextFreeAddress()
	if !ok {
		return false
	}

	if m.pages[addr].pid != "" {
		return false
	}

	m.pages[addr] = page{pid, processPage}
	m.queue = append(m.queue, addr)
	return true
}

// AccessPage ...
func (m *memory) AccessPage(pid string, address int) (int, bool) {
	displacedAddress := address / m.PageSize
	for realAddress, page := range m.pages {
		if page.pid == pid && page.virtualAddress == displacedAddress {
			return realAddress, true
		}
	}
	return -1, false
}

// RemovePages ...
func (m *memory) RemovePages(pid string) {
	for i, p := range m.pages {
		if p.pid == pid {
			m.pages[i] = page{}
			m.freeAddress(i)
		}
	}
}

// NextSwappingCandidate ...
func (m *memory) NextSwappingCandidate() (page, bool) {
	// FIFO
	if len(m.queue) > 0 {
		// get next from queue
		candidateAddress := m.queue[0]
		m.queue = m.queue[1:]

		// retrieve page for candidate address
		p := m.pages[candidateAddress]
		m.pages[candidateAddress] = page{}
		m.freeAddress(candidateAddress)

		return p, true
	}
	return page{}, false
}
