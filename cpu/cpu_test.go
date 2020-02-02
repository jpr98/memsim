package cpu_test

import (
	"testing"

	"github.com/jpr98/memsim/cpu"
	"github.com/jpr98/memsim/mem"
)

func createEmptyCPU() cpu.CPU {
	real, _ := mem.New(0, 10)
	swap, _ := mem.New(0, 10)
	return cpu.New(*real, *swap)
}

func createCPU() cpu.CPU {
	real, _ := mem.New(2048, 16)
	swap, _ := mem.New(2048, 16)
	return cpu.New(*real, *swap)
}

func TestCreateProcess(t *testing.T) {
	c := createEmptyCPU()
	err := c.CreateProcess("1", 10)
	if err == nil {
		t.Error("CreateProcess should return an error if there is no space in memory")
	}

	c = createCPU()
	err = c.CreateProcess("1", 30)
	if err != nil {
		t.Error("CreateProcess should not return an error if there is space in memory")
	}
}

func TestAccessProcess(t *testing.T) {
	c := createEmptyCPU()
	addr, err := c.AccessProcess("1", 20)
	if err == nil {
		t.Error("AccessProcess should return an error if process not in memory")
	}

	c = createCPU()
	c.CreateProcess("1", 200)
	addr, err = c.AccessProcess("1", 20)
	if err != nil {
		t.Error("AccessProcess shouldn't return an error if process in memory")
	}
	if addr != (20 / 16) {
		t.Errorf("Address should be %d", (20 / 16))
	}
}

func TestDeleteProcess(t *testing.T) {
	c := createEmptyCPU()
	err := c.DeleteProcess("1")
	if err == nil {
		t.Error("DeleteProcess should return an error if process not found in memory")
	}

	c = createCPU()
	_ = c.CreateProcess("1", 1000)
	err = c.DeleteProcess("1")
	if err != nil {
		t.Error("DeleteProcess shouldn't return an error if process is in memory")
	}
	_, err = c.AccessProcess("1", 10)
	if err == nil {
		t.Error("DelteProcess is not deleting the process from memory")
	}

}
