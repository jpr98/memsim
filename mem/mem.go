package mem

import (
	"errors"
)

type page struct {
	pid            string
	virtualAddress int
}

type memory struct {
	freeList []int
	pages    []page
	PageSize int
	queue    []int
	policy   *policy
}

// new creates a new Memory
func new(size, pageSize int, policy *policy) (*memory, error) {
	if pageSize == 0 {
		return nil, errors.New("PageSize should not be zero")
	}
	numOfPages := size / pageSize
	return &memory{
		freeList: createFreeList(numOfPages),
		pages:    make([]page, numOfPages),
		PageSize: pageSize,
		policy:   policy,
	}, nil
}

func createFreeList(size int) []int {
	list := make([]int, size)
	for i := 0; i < size; i++ {
		list[i] = i
	}
	return list
}

func (m *memory) GetPages() []page {
	return m.pages
}

func (m *memory) getNextFreeAddress() (int, bool) {
	if len(m.freeList) < 1 {
		return -1, false
	}
	addr := m.freeList[0]
	m.freeList = m.freeList[1:]
	return addr, true
}

func (m *memory) freeAddress(i int) {
	m.freeList = append(m.freeList, i)
}
