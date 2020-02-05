package mem

import (
	"fmt"
	"strings"
)

type policy struct {
	FIFO bool
	LRU  bool
}

func createPolicy(s string) (policy, error) {
	s = strings.ToLower(s)
	if s == "fifo" {
		return policy{
			FIFO: true,
			LRU:  false,
		}, nil
	} else if s == "lru" {
		return policy{
			FIFO: false,
			LRU:  true,
		}, nil
	} else {
		return policy{}, fmt.Errorf("%s is not a known policy", s)
	}
}
