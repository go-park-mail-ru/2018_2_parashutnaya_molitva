package chess

type coord struct {
	r int
	c int
}

func (c *coord) add(o *coord) coord {
	return coord{c.r + o.r, c.c + o.c}
}

func (c *coord) subtract(o *coord) coord {
	return coord{c.r - o.r, c.c - o.c}
}

func (c *coord) multiply(o *coord) coord {
	return coord{c.r * o.r, c.c * o.c}
}
