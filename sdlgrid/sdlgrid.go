// TODO: listitem
// TODO: neues item viewport mit slidern
// TODO: anpassen von size() für alle ecken
// TODO: resize handles auch in allen ecken
// TODO: getminsize prüfen. bei absoluter positionierung auch die Position der subitems einbeziehen
// TODO: TextInput textmarkierung mit tastatur und maus
// TODO: Standard Icon set für menuepfeile und schliessen minimieren vergroessern

package sdlgrid

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//	"path/filepath"

//
// apt install libsdl2{,-image,-mixer,-ttf,-gfx}-dev
// go get -v github.com/veandco/go-sdl2/sdl@master
// go get -v github.com/veandco/go-sdl2/{sdl,img,mix,ttf}

// ----------------------------------------------------------------

const (
	// Layout Spec.
	//LS_POS_CENTER = -1 //center the item in x or y dimension LS_POS_PCT 50
	//LS_POS_LEFT   = -2 // align item left LS_POS_PCT 0
	//LS_POS_RIGHT  = -3 // align item right LS_POS_PCT 100
	//LS_POS_TOP    = -2 // align item top LS_POS_PCT 0
	//LS_POS_BOTTOM = -3 // align item bottom LS_POS_PCT 100

	LS_POS_PCT = -4 // position is set by percentage of parent frame
	LS_POS_ABS = 0  // position  in pixel

	//LS_SIZE_EXPAND   = -1 // width/height shall be expanded to parent item LS_SIZE_PCT 100
	LS_SIZE_COLLAPSE = -2 // width/height shall be collapsed to minimal extent of content
	LS_SIZE_PCT      = -4 // size is set by percentage of parent frame
	LS_SIZE_ABS      = 0  // absolute size in pixel

	ITEM_TEXTINPUT_ALIGN_CENTER = -1
	ITEM_TEXTINPUT_ALIGN_LEFT   = -2
	ITEM_TEXTINPUT_ALIGN_RIGHT  = -3

	ITEM_MOVE_MODE_X  = -1 // item can be moved horizontally
	ITEM_MOVE_MODE_Y  = -2 // item can be moved vertically
	ITEM_MOVE_MODE_XY = -3 // item can be moved in any direction

	ITEM_BUTTON_STYLE_EDGE = -1
	ITEM_BUTTON_STYLE_FLAT = -2

	ITEM_HANDLE_STYLE_HIDDEN   = -1
	ITEM_HANDLE_STYLE_DARKEDGE = -2
	ITEM_HANDLE_STYLE_FLAT     = -3

	ITEM_SLIDER_STYLE_VERTICAL   = -1
	ITEM_SLIDER_STYLE_HORIZONTAL = -2

	MENU_ORIENTATION_VERTICAL   = -1
	MENU_ORIENTATION_HORIZONTAL = -2

	ITEM_STATE_ACTIVE     = 0
	ITEM_STATE_BACKGROUND = -1
)

type LayoutParam struct {
	S int32
	V int32
}

type LayoutSpec struct {
	X LayoutParam
	Y LayoutParam
	W LayoutParam
	H LayoutParam
}

type infraItemList struct {
	items []Item
}

func (i *infraItemList) GetList() []Item { return i.items }

func (i *infraItemList) CheckTopDown(x, y int32) (Item, bool) {

	for n := len(i.items) - 1; n >= 0; n-- {
		if i.items[n].CheckPos(x, y) {
			return i.items[n], true
		}
	}
	return nil, false
}
func (i *infraItemList) ClearList(s Item) {
	for n := range i.items {
		i.items[n] = nil
	}
	i.items = i.items[:0] // cutt off
}

func (i *infraItemList) AddItem(s Item) {
	i.items = append(i.items, s) // add to top
}

func (i *infraItemList) getItemIndex(s Item) (int, bool) {

	for n, ri := range i.items {
		if ri == s {
			return n, true
		}
	}
	return 0, false
}

func (i *infraItemList) RemoveItem(s Item) bool {

	if n, found := i.getItemIndex(s); found {
		al := len(i.items)
		if n < (al - 1) {
			//if its not the last entry to remove
			copy(i.items[n:], i.items[n+1:]) // left shift tail over found entry
		}
		i.items[al-1] = nil      // dont let references remain in unused part of array
		i.items = i.items[:al-1] // cutt off last
		return true
	}
	return false
}
func (i *infraItemList) swapItem(n1, n2 int) {
	if n1 != n2 {
		ti := i.items[n2]
		i.items[n2] = i.items[n1]
		i.items[n1] = ti
	}
}
func (i *infraItemList) GetTop() Item {
	l := len(i.items)
	if l > 0 {
		// if it is not allready the top one
		return i.items[l-1]
	}
	return nil
}
func (i *infraItemList) ShiftTop(s Item) {
	// if there is more than one entry
	l := len(i.items)
	if l > 1 {
		// and if it is not allready the top one
		if i.items[l-1] != s {
			if i.RemoveItem(s) {
				i.AddItem(s)
			}
		}
	}
}

func (i *infraItemList) ShiftUp(s Item) {
	if n, found := i.getItemIndex(s); found {
		if n < len(i.items)-1 {
			i.swapItem(n, n+1)
		}
	}
}

func (i *infraItemList) ShiftDown(s Item) {
	if n, found := i.getItemIndex(s); found {
		if n > 0 {
			i.swapItem(n, n-1)
		}
	}
}

func (i *infraItemList) ShiftBottom(s Item) {
	if len(i.items) > 1 {
		// if there is more than one entry
		if n, found := i.getItemIndex(s); found {
			// shift the beginning items to the right over the found one
			copy(i.items[1:n+1], i.items[0:n]) // left shift tail over found entry
			i.items[0] = s
		}
	}
}

//--------------------------------------------------------------------

// gets the normalized percentage
func utilNormPct(base, pct int32) int32 {
	r := int32((100000.0 / float32(base)) * float32(pct))
	if r > 100000 {
		r = 100000
	}
	return r
}

// gets the percentage of fraction from base
func utilGetPct(base, fraction int32) int32 {
	r := int32(((float32(fraction) / float32(base)) * 100000.0))
	if r > 100000 {
		r = 100000
	}
	return r
}

// gets pct percent of base
func utilPct(base, pct int32) int32 {
	r := int32((float32(base) * (float32(pct) / 100000.0)) + 0.5)
	if r > base {
		r = base
	}
	return r
}

func utilFrameReduce(f *sdl.Rect, mgn int32) {
	f.X += mgn
	f.Y += mgn
	f.W -= (mgn * 2)
	f.H -= (mgn * 2)
}
func utilFrameGetData(f *sdl.Rect) (x, y, w, h int32) {
	return f.X, f.Y, f.W, f.H
}

func utilCollapseLayoutSpec(ls *sdl.Rect) (x, y, w, h int32) {
	// translate the spec into a minimum fixed size spec
	x, y, w, h = utilFrameGetData(ls)

	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}

	return x, y, w, h
}

func printSeparator() {
	fmt.Println("-------------------------------------------------------------------------------------------------")
}

func utilSizeMax(w1, h1, w2, h2 int32) (w, h int32) {
	if w = w2; w1 > w2 {
		w = w1
	}
	if h = h2; h1 > h2 {
		h = h1
	}
	return w, h
}

func utilFrameCpyMoved(src *sdl.Rect, x, y int32) (dst *sdl.Rect) {
	return &sdl.Rect{X: src.X + x, Y: src.Y + y, W: src.W, H: src.H}
}

func utilRendererSetDrawColor(r *sdl.Renderer, c *sdl.Color) {
	r.SetDrawColor(c.R, c.G, c.B, c.A)
}

func utilRenderSolidBorder(rd *sdl.Renderer, r *sdl.Rect, c *sdl.Color) {
	rd.SetDrawBlendMode(sdl.BLENDMODE_NONE)
	utilRendererSetDrawColor(rd, c)
	rd.DrawRect(r)
}

func utilRenderFillRect(rd *sdl.Renderer, r *sdl.Rect, c *sdl.Color) {
	rd.SetDrawBlendMode(sdl.BLENDMODE_NONE)
	utilRendererSetDrawColor(rd, c)
	rd.FillRect(r)
}

func utilRenderText(rd *sdl.Renderer, f *ttf.Font, t string, r *sdl.Rect, c *sdl.Color) {
	// as TTF_RenderText_Solid could only be used on SDL_Surface then you have to create the surface first

	//surface, _ := f.RenderUTF8Solid(t, *c)
	sfc, _ := f.RenderUTF8Blended(t, *c)
	//surface, _ := f.RenderUTF8Shaded(t, *c)

	//now you can convert it into a texture
	//SDL_Texture* Message = SDL_CreateTextureFromSurface(renderer, surfaceMessage);
	tx, _ := rd.CreateTextureFromSurface(sfc)

	//you put the renderer's name first, the Message, the crop size(you can ignore this if you don't want to dabble with cropping), and the rect which is the size and coordinate of your texture
	//SDL_RenderCopy(renderer, Message, NULL, &Message_rect);

	// center target rect
	var tr sdl.Rect = *r
	tr.W = sfc.W
	tr.H = sfc.H
	tr.X += (r.W - sfc.W) / 2
	tr.Y += (r.H - sfc.H) / 2

	rd.Copy(tx, nil, &tr)

	sfc.Free()
	tx.Destroy()
}
func utilRenderSymbolX(rd *sdl.Renderer, r *sdl.Rect, cl *ColorScheme) {
	p1x := r.X + 1
	p1y := r.Y + 1
	p2x := p1x + r.W - 2
	p2y := p1y
	p3x := p1x
	p3y := p1y + r.H - 2
	p4x := p2x
	p4y := p3y

	if (p2x-p1x)%2 == 1 {
		p2x--
		p4x--
	}
	if (p3y-p1y)%2 == 1 {
		p3y--
		p4y--
	}

	rd.SetDrawBlendMode(sdl.BLENDMODE_NONE)
	utilRendererSetDrawColor(rd, &cl.baseDark)
	rd.DrawLine(p1x, p1y, p4x, p4y)
	rd.DrawLine(p3x, p3y, p2x, p2y)

	c := cl.baseDark
	//c := sdl.Color{255, 50, 50, 255}

	c.A = 100
	rd.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	utilRendererSetDrawColor(rd, &c)
	p := make([]sdl.Point, 13)
	p[0].X = p1x - 1
	p[0].Y = p1y
	p[1].X = p1x
	p[1].Y = p1y - 1

	p[2].X = (p1x + (p2x-p1x)/2)
	p[2].Y = (p1y + (p3y-p1y)/2) - 1

	p[3].X = p2x
	p[3].Y = p2y - 1
	p[4].X = p2x + 1
	p[4].Y = p2y

	p[5].X = (p1x + (p2x-p1x)/2) + 1
	p[5].Y = (p1y + (p3y-p1y)/2)

	p[6].X = p4x + 1
	p[6].Y = p4y
	p[7].X = p4x
	p[7].Y = p4y + 1

	p[8].X = (p1x + (p2x-p1x)/2)
	p[8].Y = (p1y + (p3y-p1y)/2) + 1

	p[9].X = p3x
	p[9].Y = p3y + 1

	p[10].X = p3x - 1
	p[10].Y = p3y

	p[11].X = (p1x + (p2x-p1x)/2) - 1
	p[11].Y = (p1y + (p3y-p1y)/2)

	p[12].X = p1x - 1
	p[12].Y = p1y

	rd.DrawLines(p)

}
func utilRenderShadowBorder(rd *sdl.Renderer, r *sdl.Rect, bc *ColorScheme, sunken bool) {

	rd.SetDrawBlendMode(sdl.BLENDMODE_NONE)
	//top and left edge
	if sunken {
		utilRendererSetDrawColor(rd, &bc.lowEdge)
	} else {
		utilRendererSetDrawColor(rd, &bc.highEdge)
	}
	rd.DrawLine(r.X, r.Y, r.X+r.W-2, r.Y)
	rd.DrawLine(r.X, r.Y, r.X, r.Y+r.H-2)

	//blend corners
	utilRendererSetDrawColor(rd, &bc.midEdge)
	rd.DrawPoint(r.X+r.W-1, r.Y)
	rd.DrawPoint(r.X, r.Y+r.H-1)

	//bottom and right
	if sunken {
		utilRendererSetDrawColor(rd, &bc.highEdge)
	} else {
		utilRendererSetDrawColor(rd, &bc.lowEdge)
	}
	rd.DrawLine(r.X+1, r.Y+r.H-1, r.X+r.W-1, r.Y+r.H-1)
	rd.DrawLine(r.X+r.W-1, r.Y+1, r.X+r.W-1, r.Y+r.H-1)

}

func utilColorMix(a *sdl.Color, b *sdl.Color) *sdl.Color {
	//average of every color channel
	return &sdl.Color{R: uint8((uint32(a.R) + uint32(b.R)) >> 1), G: uint8((uint32(a.G) + uint32(b.G)) >> 1), B: uint8((uint32(a.B) + uint32(b.B)) >> 1), A: uint8((uint32(a.A) + uint32(b.A)) >> 1)}
}

func utilColorDim(c *sdl.Color, pct int) *sdl.Color {

	if pct > 100 || pct < -100 || pct == 0 {
		return c
	}
	r := c.R
	g := c.G
	b := c.B

	if pct < 0 {
		//dim
		p := 1.0 - (float32(pct) / -100)
		r = (uint8(float32(r) * p))
		g = (uint8(float32(g) * p))
		b = (uint8(float32(b) * p))
	} else {
		//brighten
		p := 1.0 - (float32(pct) / 100)
		r = (255 - uint8(float32(255-r)*p)) & 0xff
		g = (255 - uint8(float32(255-g)*p)) & 0xff
		b = (255 - uint8(float32(255-b)*p)) & 0xff
	}

	return &sdl.Color{R: r, G: g, B: b, A: c.A}
}

// Gets the visual brightness of a color. For color levels until about 130 use bright text - otherwise dark text
func utilColorLevel(c *sdl.Color) int {
	return ((3 * int(c.R)) + (6 * int(c.G)) + int(c.B)) / 10

}

func utilColorReverse(c *sdl.Color) *sdl.Color {
	var rc sdl.Color
	rc.R = 255 - c.R
	rc.G = 255 - c.G
	rc.B = 255 - c.B
	rc.A = c.A
	return &rc
}
