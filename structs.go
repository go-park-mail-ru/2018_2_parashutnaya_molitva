package chess

type Coord struct {
	r int
	c int
}

func (c *Coord) add(o *Coord) Coord {
	return Coord{c.r + o.r, c.c + o.c}
}
