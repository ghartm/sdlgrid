package sdlgrid

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

//--------------------------------------------------------------------

// Single Container Item with a border and an optional Headline
type ItemFrame struct {
	ItemBase
	textBox *ItemText
	subItem Item
}

func NewItemFrame(win *RootWindow) *ItemFrame {
	i := new(ItemFrame)
	i.o = Item(i)
	i.setRootWindow(win)
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_PCT, 0, 0, 100000, 100000)
	return i
}
func (i *ItemFrame) oReportSubitems(lvl int) {

	if i.textBox != nil {
		i.subItem.Report(lvl)
	}
}
func (i *ItemFrame) oFindSubItem(x, y int32, e sdl.Event) (found bool, item Item) {
	//there is only one subitem to check
	if i.subItem != nil {
		if i.subItem.CheckPos(x, y) {
			return true, i.subItem
		}
	}
	return false, nil
}
func (i *ItemFrame) oWithItems(fn func(Item)) {
	fmt.Printf("%s: ItemBase.oWithItems() no subitems\n", i.GetName())
	// handle the call for myself
	fn(i)
	// forward to all subitems
	if i.textBox != nil {
		i.textBox.oWithItems(fn)
	}
	if i.subItem != nil {
		i.subItem.oWithItems(fn)
	}

}

// report the items minimal size
func (i *ItemFrame) oGetMinSize() (int32, int32) {
	// get the minimum bounding box for myself
	if !i.useTmpMinSize {
		var w, h, mw, mh int32

		// get the size diff between iframe and subframe
		t, r, b, l := i.getBorderValues()

		// text
		if i.textBox != nil {
			tx, _, tw, _ := i.textBox.GetCollapsedSpec()
			mw = (tx - l) + tw
		}

		// add subitems
		if i.subItem != nil {
			w, h = i.subItem.oGetMinSize()
			mw, mh = utilSizeMax(mw, mh, w, h)
		}

		_, _, w, h = i.GetCollapsedSpec()
		w -= r + l
		h -= t + b
		mw, mh = utilSizeMax(mw, mh, w, h)

		i.minw = mw + r + l
		i.minh = mh + t + b

		fmt.Printf("%s: ItemFrame.oGetMinSize()= %d,%d\n", i.GetName(), i.minw, i.minh)
		i.useTmpMinSize = true
	}
	return i.minw, i.minh
}

// every item gets asked to layout its subitems
func (i *ItemFrame) oNotifyPostLayout(sizeChanged bool) {
	//fmt.Printf("%s: ItemFrame.oNotifyPostLayout()\n", i.GetName())

	// layout decoration
	if i.textBox != nil {
		//fmt.Printf("--> textBox:%s\n", i.textBox.GetName())
		i.textBox.Layout(i.GetFrame(), sizeChanged)
	}

	// layout content
	if i.subItem != nil {
		//fmt.Printf("--> subitem:%s\n", i.subItem.GetName())
		i.subItem.Layout(i.oGetSubFrame(), sizeChanged)
	}
}

func (i *ItemFrame) SetSubItem(ni Item) {
	i.subItem = ni
	ni.SetParent(i)
}

func (i *ItemFrame) SetText(t string) *ItemFrame {
	if i.textBox == nil {
		i.textBox = NewItemText(i.win)
		i.textBox.SetParent(i)
		i.textBox.SetName(i.name + ".textBox")
		i.textBox.SetSpecX(LS_POS_ABS, i.textBox.GetStyle().Spacing*2)
		i.textBox.UseBaseColor(true)
	}
	i.textBox.SetText(t)
	return i
}
func (i *ItemFrame) getBorderValues() (t, r, b, l int32) {
	// frame line
	s := i.GetStyle().Spacing
	r = s + 1
	b = s + 1
	l = s + 1
	if i.textBox != nil {
		_, _, _, _, _, _, _, t = i.textBox.GetSpec()
		t += s
	} else {
		t = s + 1
	}
	return t, r, b, l
}

func (i *ItemFrame) oGetSubFrame() *sdl.Rect {
	//
	// copy the item frame to a subframe
	var sf sdl.Rect
	t, r, b, l := i.getBorderValues()

	sf.X = i.iframe.X + l
	sf.Y = i.iframe.Y + t
	sf.W = i.iframe.W - r - l
	sf.H = i.iframe.H - t - b

	// clip it to minimum
	if sf.W < i.iframe.W-sf.W {
		sf.W = i.iframe.W - sf.W
	}
	if sf.H < i.iframe.H-sf.H {
		sf.H = i.iframe.H - sf.H
	}

	return &sf
}

func (i *ItemFrame) oRender() {
	//fmt.Printf("%s: ItemFrame.oRender()\n", i.GetName())

	r := i.GetRenderer()
	c := i.GetColorScheme()
	r.SetDrawBlendMode(sdl.BLENDMODE_NONE)

	//utilRenderSolidBorder(r, &i.outerFrame, s.ColorGreen)
	//utilRenderSolidBorder(r, &i.outerFrame, i.GetStyle().ColorGreen)

	if i.textBox != nil {
		//if there is a text, let a gap in the frame
		_, _, _, _, tx, _, tw, th := i.textBox.GetSpec()

		th /= 2

		rx, ry := i.GetFramePos()
		rw, rh := i.GetFrameSize()

		utilRendererSetDrawColor(r, &c.base)

		r.DrawLine(rx, ry+th, rx+tx-1, ry+th)           //top left
		r.DrawLine(rx+tx+tw+1, ry+th, rx+rw-1, ry+th)   // top right
		r.DrawLine(rx, ry+th, rx, ry+rh-1)              // left
		r.DrawLine(rx+rw-1, ry+th, rx+rw-1, ry+rh-1)    // right
		r.DrawLine(rx, ry+rh-1, (rx + rw - 1), ry+rh-1) // bottom

		// render the textbox
		i.textBox.oRender()

	} else {
		utilRenderSolidBorder(r, &i.iframe, &c.baseDeco)
	}

	// forward to content
	if i.subItem != nil {
		i.subItem.Render()
	}
}
