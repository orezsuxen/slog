package collect

import (
	"fmt"
	"slices"
)

type msgInfo struct {
	count   int
	arrival int
	msg     string
}

type Collector struct {
	data     map[string]msgInfo
	msgCount int
}

func New(size int) Collector {
	return Collector{
		data:     make(map[string]msgInfo),
		msgCount: 1,
	}
}

func (c *Collector) Store(msg string) {
	v := c.data[msg]
	c.data[msg] = msgInfo{v.count + 1, c.msgCount, msg}
	c.msgCount += 1
}

func (c *Collector) Display(amount int) string {
	//get messages from map into slice
	msgs := make([]msgInfo, 0)
	for _, v := range c.data {
		msgs = append(msgs, v)
	}

	//sort slice by arrival of msg
	slices.SortFunc(msgs, func(a msgInfo, b msgInfo) int {
		return a.arrival - b.arrival
	})

	retval := ""
	n := 0
	if amount > len(msgs) {
		n = len(msgs)
	} else {
		n = amount
	}
	msgs = msgs[len(msgs)-n:]
	slices.SortFunc(msgs, func(a msgInfo, b msgInfo) int {
		if a.msg > b.msg {
			return 1
		} else if a.msg < b.msg {
			return -1
		} else {
			return 0
		}
	})

	//get 10 last arrived msgs
	for _, m := range msgs[:] {
		retval += fmt.Sprint(m.msg, " ::: ", m.count, "\n")
	}
	return retval
}
