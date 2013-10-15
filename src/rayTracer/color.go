package rayTracer

const (
	OFFSET_RED   = iota
	OFFSET_GREEN = iota
	OFFSET_BLUE  = iota
	OFFSET_MAX   = iota
)

type ColorV struct {
	Red, Green, Blue float64
}

func (c *ColorV) AddEq(c2 ColorV) ColorV {
	c.Red += c2.Red
	c.Green += c2.Green
	c.Blue += c2.Blue
	return *c
}

func (c *ColorV) Mul(c2 ColorV) ColorV {
	return ColorV{c.Red * c2.Red, c.Green * c2.Green, c.Blue * c2.Blue}
}

func (c *ColorV) Add(c2 ColorV) ColorV {
	return ColorV{c.Red + c2.Red, c.Green + c2.Green, c.Blue + c2.Blue}
}

func (c *ColorV) MulC(coef float64) ColorV {
	return ColorV{c.Red * coef, c.Green * coef, c.Blue * coef}
}
