package mem

import (
	"container/list"
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

type set map[int]list.Element

type lruCache struct {
	list *list.List
	s    set
}

// Refer inserts a pid to the LRU Cache
func (c *lruCache) Refer(addr int) {
	el, ok := c.s[addr]
	if !ok {
		e := c.list.PushFront(addr)
		c.s[addr] = *e
	}
	c.list.Remove(&el)
	e := c.list.PushFront(addr)
	c.s[addr] = *e
}

func (c *lruCache) GetNext() int {
	e := c.list.Front()
	c.list.Remove(e)

	for k, v := range c.s {
		if v == *e {
			delete(c.s, k)
			return k
		}
	}

	return 0
}

// Remove deletes a pid from the LRU Cache
func (c *lruCache) Remove(addr int) {
	el, _ := c.s[addr]
	c.list.Remove(&el)
}
