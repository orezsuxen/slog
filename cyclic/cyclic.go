package cyclic

type Cyclic struct {
	pos   int
	count int
	Data  []string
}

func New(size int) Cyclic {
	var retval Cyclic
	retval.Data = make([]string, size)

	return retval
}

func (c *Cyclic) Store(msg string) {
	c.Data[c.pos] = msg
	c.pos += 1
	if c.pos >= len(c.Data) {
		c.pos = 0
	}
	if c.count < len(c.Data) {
		c.count += 1
	}
}

func (c *Cyclic) Get() []string {
	retval := make([]string, 0) //v1

	if c.count >= len(c.Data) {
		retval = append(retval, c.Data[c.pos:]...)
	}

	retval = append(retval, c.Data[:c.pos]...)
	return retval
}

func (c *Cyclic) GetNum(req int) []string {
	retval := make([]string, 0) //v1
	if req >= c.count {
		if c.count >= len(c.Data) {
			retval = append(retval, c.Data[c.pos:]...)
		}

		retval = append(retval, c.Data[:c.pos]...)
		return retval
	} else { // req smaler than count
		n := req
		if c.count >= len(c.Data) {
			back := len(c.Data) - c.pos
			if req >= back {
				retval = append(retval, c.Data[c.pos:]...)
				n = req - back
			} else {
				retval = append(retval, c.Data[c.pos:c.pos+req]...)
				return retval
			}
		}

		retval = append(retval, c.Data[:n]...)
		return retval

	}

}

// Interfaces
func (c *Cyclic) Write(p []byte) (n int, err error) {
	c.Store(string(p))

	return len(p), nil
}
