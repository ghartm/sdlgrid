package sdlgrid

import "github.com/veandco/go-sdl2/sdl"

type ColorScheme struct {
	base        sdl.Color
	baseDark    sdl.Color
	baseBright  sdl.Color
	baseReverse sdl.Color
	baseDeco    sdl.Color
	text        sdl.Color
	lowEdge     sdl.Color
	midEdge     sdl.Color
	highEdge    sdl.Color
}

func (c *ColorScheme) SetBaseColor(bc *sdl.Color) {

	c.base = *bc
	c.baseReverse = *utilColorReverse(bc)

	cl := utilColorLevel(bc)
	//fmt.Printf("SetBaseColor r:%d g:%d b:%d level: %d\n", bc.R, bc.G, bc.B, cl)

	c.baseDark = *utilColorDim(bc, -40)
	c.baseBright = *utilColorDim(bc, 40)

	if cl > 130 {
		//c.text = sdl.Color{0, 0, 0, 255}
		c.text = *utilColorDim(bc, -70)
		c.baseDeco = c.baseDark
	} else {
		//c.text = sdl.Color{255, 255, 255, 255}
		c.text = *utilColorDim(bc, 70)
		c.baseDeco = c.baseBright
	}
	c.baseDeco = c.baseBright

	c.baseReverse = *utilColorReverse(bc)
	c.highEdge = *utilColorDim(bc, 15)
	c.lowEdge = *utilColorDim(bc, -15)
	c.midEdge = *utilColorMix(&c.lowEdge, &c.highEdge)
}
