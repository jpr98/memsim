package mem_test

import (
	"testing"

	"github.com/jpr98/memsim/mem"
)

func TestNew(t *testing.T) {
	_, err := mem.New(0, 0)
	if err == nil {
		t.Error("New should not accept 0 as the page size")
	}

	memory, err := mem.New(100, 10)
	if err != nil {
		t.Error("New should not return error when proper values passed")
	}
	if memory == nil {
		t.Error("New shouldn't return a nil Memory")
	}
	if memory.PageSize != 10 {
		t.Error("PageSize should be value passed to New")
	}
}

func TestAllocatePage(t *testing.T) {
	memory, _ := mem.New(0, 10)
	ok := memory.AllocatePage("1", 1)
	if ok {
		t.Error("AllocatePage should return false if there's no space to allocate new pages")
	}

	memory, _ = mem.New(100, 10)
	ok = memory.AllocatePage("1", 1)
	if !ok {
		t.Error("AllocatePage should return true when it allocates a new page successfully")
	}
}

func TestAccessPage(t *testing.T) {
	memory, _ := mem.New(100, 10)
	_, ok := memory.AccessPage("1", 1)
	if ok {
		t.Error("AccessPage should return false if PID not found")
	}

	_ = memory.AllocatePage("1", 0)
	_ = memory.AllocatePage("1", 1)
	addr, okay := memory.AccessPage("1", 1)
	if !okay {
		t.Error("AccessPage should return true if PID and virtual address exists in memory")
	}
	if addr != 0 {
		t.Error("AccessPage should return correct real address")
	}
}

func TestRemovePages(t *testing.T) {
	memory, _ := mem.New(100, 10)
	ok := memory.RemovePages("1")
	if ok {
		t.Error("RemovePages should return false if PID not present in memory")
	}

	_ = memory.AllocatePage("1", 1)
	_ = memory.AllocatePage("1", 2)
	ok = memory.RemovePages("1")
	_, ok = memory.AccessPage("1", 2)
	if ok {
		t.Error("RemovePages should delete all pages with a given PID")
	}
}
