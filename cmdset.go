package go_hotstuff

import (
	"container/list"
	"sync"
)

type CmdSet interface {
	Add(cmds ...string)
	Remove(cmds ...string)
	GetFirst(n int) []string
	IsProposed(cmd string) bool
	MarkProposed(cmds ...string)
}

type cmdElement struct {
	cmd      string
	proposed bool
}

type cmdSetImpl struct {
	lock  sync.Mutex
	order list.List
	set   map[string]*list.Element
}

func NewCmdSet() *cmdSetImpl {
	c := &cmdSetImpl{
		set: make(map[string]*list.Element),
	}
	c.order.Init()
	return c
}

// Add add cmds to the list and set, duplicate one will be ignored
func (c *cmdSetImpl) Add(cmds ...string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, cmd := range cmds {
		// avoid duplication
		if _, ok := c.set[cmd]; ok {
			continue
		}
		e := c.order.PushBack(&cmdElement{
			cmd:           cmd,
			proposed:      false,
		})
		c.set[cmd] = e
	}
}

// Remove remove commands from set and list
func (c *cmdSetImpl) Remove(cmds ...string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, cmd := range cmds {
		if e, ok := c.set[cmd]; ok {
			c.order.Remove(e)
			delete(c.set, cmd)
		}
	}
}

// GetFirst return the top n unused commands from the list
func (c *cmdSetImpl) GetFirst(n int) []string {
	c.lock.Lock()
	defer c.lock.Unlock()

	if len(c.set) == 0 {
		return nil
	}

	cmds := make([]string, 0, n)
	i := 0
	// get the first element of list
	e := c.order.Front()

	for i < n {
		if e == nil {
			break
		}
		if cmd := e.Value.(*cmdElement); !cmd.proposed {
			cmds = append(cmds, cmd.cmd)
			i++
		}
		e = e.Next()
	}
	return cmds
}

func (c *cmdSetImpl) IsProposed(cmd string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if e, ok := c.set[cmd]; ok {
		return e.Value.(*cmdElement).proposed
	}
	return false
}

// MarkProposed will mark the given commands as proposed and move them to the back of the queue
func (c *cmdSetImpl) MarkProposed(cmds ...string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, cmd := range cmds {
		if e, ok := c.set[cmd]; ok {
			e.Value.(*cmdElement).proposed = true
			// Move to back so that it's not immediately deleted by a call to TrimToLen()
			c.order.MoveToBack(e)
		} else {
			// new cmd, store it to back
			e := c.order.PushBack(&cmdElement{cmd: cmd, proposed: true})
			c.set[cmd] = e
		}
	}
}