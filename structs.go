package chess

type Coord struct {
	r int
	c int
}

func (c *Coord) add(o *Coord) Coord {
	return Coord{c.r + o.r, c.c + o.c}
}

func (c *Coord) subtract(o *Coord) Coord {
	return Coord{c.r - o.r, c.c - o.c}
}

func (c *Coord) multiply(o *Coord) Coord {
	return Coord{c.r * o.r, c.c * o.c}
}
