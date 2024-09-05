package collect

import ()

type msgInfo struct {
	count   int
	arrival int
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
	c.data[msg] = msgInfo{v.count + 1, c.msgCount}
	c.msgCount += 1
}
