package sdlgrid

import (
	"fmt"
	"reflect"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// --------------------------------------------------------------------
// TODO: "register Callback" functionality (button, motion)
type Item interface {
	Report(int)
	oReportSubitems(int)

	oGetMinSize() (int32, int32)
	oNotifyPostLayout(bool)
	IsAutoSize() bool
	IsAutoPos() bool
	SetSpec(sx, sy, sw, sh, x, y, w, h int32)
	SetSpecX(sv, v int32)
	SetSpecY(sv, v int32)
	SetSpecSize(int32, int32, int32, int32)
	SetSpecPos(int32, int32, int32, int32)
	GetSpec() (sx, sy, sw, sh, x, y, w, h int32)
	GetSize() (int32, int32)
	GetPos() (int32, int32)
	GetFrame() *sdl.Rect
	GetFramePos() (x, y int32)
	GetFrameSize() (w, h int32)
	MakeInnerFrame(mgn int32) *sdl.Rect
	Layout(pf *sdl.Rect, sizeChanged bool)
	GetCollapsedSpec() (x, y, w, h int32)
	SetState(s int)
	SetHidden(bool)

	// implemented by ItemBase and may be by specific Item

	oNotifyMouseMotion(x, y, dx, dy int32)
	oNotifyMouseButton(x int32, y int32, button uint8, state uint8) // state:0=released 1=pressed
	oNotifyMouseFocusGained()
	oNotifyMouseFocusLost()
	oNotifyKbFocusGained()
	oNotifyKbFocusLost()
	oNotifyKbEvent(e *sdl.KeyboardEvent)
	oNotifyTextInput(r rune)
	oNotifyTimer()
	oGetSubFrame() *sdl.Rect
	oWithItems(func(Item))

	Render()
	oRender()
	oFindSubItem(x, y int32, e sdl.Event) (found bool, item Item)

	GetItemByPos(x, y int32, e sdl.Event) (bool, Item)

	GetRenderer() *sdl.Renderer //standard way of getting a renderer if needed
	GetStyle() *Style           //standard way of getting style if needed
	SetStyle(s *Style)
	GetColorScheme() *ColorScheme
	SetColorScheme(*ColorScheme)

	GetLayoutParentFrame() *sdl.Rect

	//GetFrame() *sdl.Rect
	//SetFramePos(x, y int32)
	//GetFramePos() (x, y int32)
	//SetFrameSize(w, h int32)
	//GetFrameSize() (w, h int32)

	setRootWindow(*RootWindow)

	GetName() string
	SetName(string)

	CheckPos(x, y int32) bool

	SetChanged(bool)
	GetChanged() bool

	SetParent(pi Item)
	GetParent() Item
}

// --------------------------------------------------------------------
// the commons of every item
type ItemBase struct {
	name        string
	o           Item         // self reference as Item interface - used as a workaround for "late bound fn". see the oXXXX stub functions
	win         *RootWindow  // the root window the item is connected to
	colorScheme *ColorScheme // colors used for rendering
	style       *Style       // The style to be used on rendering if it is different from the window style
	changed     bool         // true if the item needs to be refreshed on screen
	parent      Item         // Parent item
	state       int          // item State hidden/active/
	hidden      bool

	spec layoutSpec
	//specx         sdl.Rect // spec relative to parent frame
	//specval       sdl.Rect // spec values relative to parent frame
	iframe        sdl.Rect // itemframe absolute positioning (computed by layout according to spec)
	pframe        sdl.Rect // parent frame absolute positioning (parent frame used in layout computation)
	minw          int32    //temporary minimum size for layout speedup
	minh          int32    //temporary minimum size for layout speedup
	useTmpMinSize bool     //true uses the temp value instead of recomputing it
}

//--------------------------------------------------------------------

func (i *ItemBase) oWithItems(fn func(Item)) {
	fmt.Printf("%s: ItemBase.oWithItems() no subitems\n", i.GetName())
	// handle the call for myself
	fn(i)
	// forward to all subitems
}

func (i *ItemBase) SetState(s int) { i.state = s }
func (i *ItemBase) SetHidden(b bool) {
	if i.hidden != b {
		i.hidden = b
	}
}

func (i *ItemBase) Report(lvl int) {
	for n := 0; n < lvl; n++ {
		fmt.Print("--")
	}
	fmt.Printf("%s:\n", i.name)
	lvl++
	i.o.oReportSubitems(lvl)
}
func (i *ItemBase) oReportSubitems(lvl int) {
	for n := 0; n < lvl; n++ {
		fmt.Print("--")
	}

	fmt.Printf("%s: no subitems\n", i.name)
}

// --------------------------------------------------------------------
func (i *ItemBase) GetColorScheme() *ColorScheme {
	// fine if it is set
	var cs *ColorScheme

	if i.colorScheme != nil {
		cs = i.colorScheme
	} else {

		// try to get it from the parent
		if i.parent != nil {
			cs = i.parent.GetColorScheme()
		} else {
			// no parent so take the default from the style
			cs = i.GetStyle().csDefault
		}
		i.colorScheme = cs
	}

	return cs
}

func (i *ItemBase) SetColorScheme(cs *ColorScheme) {
	i.colorScheme = cs
}
func (i *ItemBase) GetRenderer() *sdl.Renderer {
	if i.win == nil {
		panic("ItemBase.GetRenderer(): win is nil")
	}
	if i.win.Renderer == nil {
		panic("ItemBase.GetRenderer(): renderer is nil")
	}
	return i.win.Renderer
}

func (i *ItemBase) GetStyle() *Style {
	if i.style == nil {
		return i.win.GetStyle()
	} else {
		return i.style
	}
}

func (i *ItemBase) GetLayoutParentFrame() *sdl.Rect { return &i.pframe }

func (i *ItemBase) SetParent(pi Item) { i.parent = pi }
func (i *ItemBase) GetParent() Item   { return i.parent }
func (i *ItemBase) SetStyle(s *Style) { i.style = s }
func (i *ItemBase) SetName(n string)  { i.name = n }
func (i *ItemBase) GetName() string   { return i.name }

func (i *ItemBase) GetChanged() bool {
	return i.changed
}

func (i *ItemBase) SetChanged(c bool) {
	i.win.SetChanged(true)
	i.changed = c
}

// returnes the "inner" frame of the item where subitems possibly reside.
func (i *ItemBase) oGetSubFrame() *sdl.Rect {
	//fmt.Printf("%s: ItemBase.oGetSubFrame()\n", i.GetName())
	// default is outer frame substracted by defaul spacing
	return i.MakeInnerFrame(i.GetStyle().spacing)
}

func (i *ItemBase) setRootWindow(rw *RootWindow) {
	if rw == nil {
		panic("nil RootWindow was assigned to item")
	}
	i.win = rw

	// get style of root window if not set
	if (i.win.style != nil) && (i.style == nil) {
		i.style = i.win.style
	}
}

// checks if the position xy is located at the item
// func (i CommonItem) checkPos(r *sdl.Rect) bool {	return i.frame.HasIntersection(r)}
func (i *ItemBase) CheckPos(x, y int32) bool {
	return (i.iframe.X <= x) && (x < i.iframe.X+i.iframe.W) && (i.iframe.Y <= y) && (y < i.iframe.Y+i.iframe.H)
}

// tries to find a subitem at a specific absolute position, if no subitem is found return false. do not return the item
// do not forward to subitems.
// do not call GetSubItems()
func (i *ItemBase) oFindSubItem(x, y int32, e sdl.Event) (found bool, item Item) {
	// for all items that have no subitems and thus do not define fn:oFindSubItem it can not find anything.
	return false, nil
}

func (i *ItemBase) Render() {
	if !i.hidden {
		i.o.oRender()
	}
}

// Render() renders the item in its outerFrame. It must not change the outerFrame
func (i *ItemBase) oRender() {
	// reder a red box if nothing else is set for the item
	utilRenderSolidBorder(i.GetRenderer(), i.GetFrame(), i.GetStyle().colorRed)
}

// Waits for a certain time and calls the items oNotifyTimer function
func (i *ItemBase) SetTimer(msec int) {
	go func() {
		time.Sleep(time.Duration(msec) * time.Millisecond)
		i.o.oNotifyTimer()
	}()
}

// The item will be notified when it gets the Keyboard Focus
func (i *ItemBase) oNotifyKbFocusGained() {} //stub
// The item will be notified when it loses the Keyboard Focus
func (i *ItemBase) oNotifyKbFocusLost() {} //stub
// The item will be notified about key events when it has keyboard focus
func (i *ItemBase) oNotifyKbEvent(e *sdl.KeyboardEvent) {} //stub
// The item will be notified about text input events when it has keyboard focus
func (i *ItemBase) oNotifyTextInput(r rune) {} //stub

// Notification about the Timer has quit in response of call to SetTimer(). called as a seperate go thread.
func (i *ItemBase) oNotifyTimer() {
	fmt.Printf("ItemBase: %s.NotifyTimer()\n", i.name)
}

// Mouse move events will be reported via this function. if they are relevan tfor this item
func (i *ItemBase) oNotifyMouseMotion(x, y, dx, dy int32) {
	//fmt.Printf("ItemBase: %s.oNotifyMouseMotion(%d,%d,%d,%d)\n", i.name, x, y, dx, dy)
}

// Button events will be reported via this function. if they are relevan tfor this item
func (i *ItemBase) oNotifyMouseButton(x, y int32, button uint8, state uint8) {
	fmt.Printf("ItemBase: %s.oNotifyMouseButton(%d,%d,%d,%d)\n", i.name, x, y, button, state)
}

// The item will be notified when the Mouse enters
func (i *ItemBase) oNotifyMouseFocusGained() {
	fmt.Printf("ItemBase: %s.oNotifyMouseFocusGained()\n", i.name)
}

// The item will be notified when the Mouse leaves
func (i *ItemBase) oNotifyMouseFocusLost() {
	fmt.Printf("ItemBase: %s.oNotifyMouseFocusLost()\n", i.name)
}

func (i *ItemBase) GetItemByPos(x, y int32, e sdl.Event) (bool, Item) {
	// do not find if hidden
	if i.hidden {
		return false, nil
	}
	// try to find a SubItem
	if found, ci := i.o.oFindSubItem(x, y, e); found {
		// forward to subitem
		return ci.GetItemByPos(x, y, e)
	} else {
		// else handle it myself
		return true, i.o
	}
}

// item needs to handle the layout for all its subitems
//func (i *ItemBase) cbLayoutSubItems() {}

func (i *ItemBase) oGetMinSize() (int32, int32) {
	//report the collapsed spec
	_, _, w, h := i.GetCollapsedSpec()
	return w, h

}
func (i *ItemBase) oNotifyPostLayout(sizeChanged bool) {
	//fmt.Printf("ItemBase.oNotifyPostLayout: no subitems!\n")
}
func (i *ItemBase) IsAutoSize() bool {
	if (i.spec.W.S != LS_SIZE_ABS) || (i.spec.H.S != LS_SIZE_ABS) {
		return true
	}
	return false
}
func (i *ItemBase) IsAutoPos() bool {
	if (i.spec.X.S != LS_POS_ABS) || (i.spec.Y.S != LS_POS_ABS) {
		return true
	}
	return false
}
func (i *ItemBase) SetSpec(sx, sy, sw, sh, x, y, w, h int32) {
	i.spec.X.S = sx
	i.spec.Y.S = sy
	i.spec.W.S = sw
	i.spec.H.S = sh
	i.spec.X.V = x
	i.spec.Y.V = y
	i.spec.W.V = w
	i.spec.H.V = h
	i.useTmpMinSize = false
}
func (i *ItemBase) SetSpecPos(sx, sy, x, y int32) {
	i.spec.X.S = sx
	i.spec.Y.S = sy
	i.spec.X.V = x
	i.spec.Y.V = y
}
func (i *ItemBase) SetSpecSize(sw, sh, w, h int32) {
	i.spec.W.S = sw
	i.spec.H.S = sh
	i.spec.W.V = w
	i.spec.H.V = h
	i.useTmpMinSize = false
}
func (i *ItemBase) SetSpecX(sv, v int32) {
	i.spec.X.S = sv
	i.spec.X.V = v
}
func (i *ItemBase) SetSpecY(sv, v int32) {
	i.spec.Y.S = v
	i.spec.Y.V = v
}
func (i *ItemBase) GetSpec() (sx, sy, sw, sh, x, y, w, h int32) {
	return i.spec.X.S, i.spec.Y.S, i.spec.W.S, i.spec.H.S, i.spec.X.V, i.spec.Y.V, i.spec.W.V, i.spec.H.V
}
func (i *ItemBase) GetSize() (int32, int32) {
	return i.iframe.W, i.iframe.H
}
func (i *ItemBase) GetPos() (int32, int32) {
	return i.iframe.X, i.iframe.Y
}
func (i *ItemBase) GetFrame() *sdl.Rect        { return &i.iframe }
func (i *ItemBase) GetFramePos() (x, y int32)  { return i.iframe.X, i.iframe.Y }
func (i *ItemBase) GetFrameSize() (w, h int32) { return i.iframe.W, i.iframe.H }
func (i *ItemBase) MakeInnerFrame(mgn int32) *sdl.Rect {

	return &sdl.Rect{X: i.iframe.X + mgn, Y: i.iframe.Y + mgn, W: i.iframe.W - (mgn * 2), H: i.iframe.H - (mgn * 2)}
}

func (i *ItemBase) Layout(pf *sdl.Rect, sizeChanged bool) {
	var x, y, w, h, minw, minh int32
	i.pframe = *pf
	clipmin := false
	minsize := false

	if sizeChanged {
		switch i.spec.H.S {
		case LS_SIZE_COLLAPSE:
			// let item compute its minimum size, remember it for later use
			minw, minh = i.o.oGetMinSize()
			minsize = true
			h = minh // size collapse
		case LS_SIZE_PCT:
			h = utilPct(pf.H, i.spec.H.V)
			clipmin = true
		case LS_SIZE_ABS:
			h = i.spec.H.V
		}
	} else {
		w = i.iframe.W
		h = i.iframe.H
	}

	switch i.spec.Y.S {
	case LS_POS_ABS:
		y = i.spec.Y.V
	case LS_POS_PCT:
		y = utilPct(pf.H-h, i.spec.Y.V)
	}

	//---------------
	if sizeChanged {
		switch i.spec.W.S {
		case LS_SIZE_COLLAPSE:
			if !minsize {
				minw, minh = i.o.oGetMinSize()
				minsize = true
			}
			w = minw
		case LS_SIZE_PCT:
			w = utilPct(pf.W, i.spec.W.V)
			clipmin = true
		case LS_SIZE_ABS:
			w = i.spec.W.V
		}
	}

	switch i.spec.X.S {
	case LS_POS_ABS:
		x = i.spec.X.V
	case LS_POS_PCT:
		x = utilPct(pf.W-w, i.spec.X.V)
	}

	// clip position to parent frame - no negatove position allowed
	if x < 0 {
		x = 0
	}

	if y < 0 {
		y = 0
	}

	// if parent frame is to narrow -  first adjust position
	if x+w > pf.W {
		x = pf.W - w
	}
	if y+h > pf.H {
		y = pf.H - h
	}

	// if the position adjustment resulted in negative positions - expand item to parent frame
	if x < 0 {
		x = 0
		w = pf.W //become parents width
		clipmin = true
	}

	if y < 0 {
		y = 0
		h = pf.H
		clipmin = true
	}

	// dont allow w,h to come below minimum size
	if sizeChanged {
		if clipmin {
			if !minsize {
				minw, minh = i.o.oGetMinSize()
			}

			if w < minw {
				w = minw
			}

			if h < minh {
				h = minh
			}
		}
	}

	i.iframe.X = x + i.pframe.X
	i.iframe.Y = y + i.pframe.Y

	newSize := false

	if i.iframe.W != w {
		newSize = true
		i.iframe.W = w
	}

	if i.iframe.H != h {
		newSize = true
		i.iframe.H = h
	}

	fmt.Printf("%s:%s:\tItemBase.Layout: ParentFrame(x:%d,y:%d,w:%d,h:%d) \tLayout(sx:%d,sy:%d,sw:%d,sh:%d,x:%d,y:%d,w:%d,h:%d) \tResult(x:%d,y:%d,w:%d,h:%d)\n", reflect.TypeOf(i.o), i.name, pf.X, pf.Y, pf.W, pf.H, i.spec.X.S, i.spec.Y.S, i.spec.W.S, i.spec.H.S, i.spec.X.V, i.spec.Y.V, i.spec.W.V, i.spec.H.V, x, y, w, h)

	// let the item layout the subitems as well
	i.o.oNotifyPostLayout(newSize)
}

func (i *ItemBase) GetCollapsedSpec() (x, y, w, h int32) {
	// translate the spec into a minimum fixed size spec

	if i.spec.X.S == LS_POS_ABS {
		x = i.spec.X.V
	}
	if i.spec.Y.S == LS_POS_ABS {
		y = i.spec.Y.V
	}
	if i.spec.W.S == LS_SIZE_ABS {
		w = i.spec.W.V
	}
	if i.spec.H.S == LS_SIZE_ABS {
		h = i.spec.H.V
	}

	return x, y, w, h
}

/*
func (i *ItemBase) convertSpec2Fixed() {

	if i.spec.X.S != LS_POS_ABS {
		i.spec.X.V = i.iframe.X - i.pframe.X
		i.spec.X.S = LS_POS_ABS
	}
	if i.spec.Y.S != LS_POS_ABS {
		i.spec.Y.V = i.iframe.Y - i.pframe.Y
		i.spec.Y.S = LS_POS_ABS
	}
	if i.spec.W.S != LS_SIZE_ABS {
		i.spec.W.V = i.iframe.W
		i.spec.W.S = LS_SIZE_ABS
	}
	if i.spec.H.S != LS_SIZE_ABS {
		i.spec.H.V = i.iframe.H
		i.spec.H.S = LS_SIZE_ABS
	}
}
*/

// Size that item by desired dx,dy and limit to layautParenFrame. Sets layout-specification to positioning for expand and collapse spec. Returns the resulting change of size.
func (i *ItemBase) Size(dx, dy int32, sizeLeft bool, sizeTop bool) (int32, int32) {

	var neww, newh, oldw, oldh, mvx, mvy int32 //absolute positions
	var domove bool

	//store the old hight and with to compute the real/effective resize in the end
	oldw = i.iframe.W
	oldh = i.iframe.H
	neww = oldw
	newh = oldh
	mvx = 0
	mvy = 0

	// get minimum for bounds check
	mw, mh := i.o.oGetMinSize()

	if dx != 0 {
		// do resizes in absolute positioning. remember the spech and convert to absolute
		// later on it will be converted back if needed
		sw := i.spec.W.S
		if i.spec.W.S != LS_SIZE_ABS {
			i.spec.W.V = i.iframe.W
			i.spec.W.S = LS_SIZE_ABS
		}

		// modify Layout Specification by the possible amount in that direction
		// do clipping. position must not leave parent frame

		if sizeLeft {
			neww = oldw - dx
		} else {
			neww = oldw + dx
		}

		if neww < 0 {
			neww = 0
		} else if i.iframe.X+i.iframe.W-neww < i.pframe.X {
			neww = i.iframe.X + i.iframe.W - i.pframe.X
		}

		//respect minsize of content or at least no negative spec
		if neww < mw {
			neww = mw
		}
		// convert back to percentage if layout was percentage before
		if sw == LS_SIZE_PCT {
			i.spec.W.S = LS_SIZE_PCT
			i.spec.W.V = utilGetPct(i.pframe.W, neww)
		} else {
			i.spec.W.V = neww
		}
		// prevent position changes
		if i.spec.X.S == LS_POS_PCT {
			i.spec.X.V = utilGetPct(i.pframe.W-neww, i.iframe.X-i.pframe.X)
		}
		if sizeLeft {
			// position correction is needed
			dx = oldw - neww
			domove = true
			mvx = dx
		}
	}
	if dy != 0 {
		// do resizes in absolute positioning. remember the spec and convert to absolute
		// later on it will be converted back if needed
		sh := i.spec.H.S
		if i.spec.H.S != LS_SIZE_ABS {
			i.spec.H.V = i.iframe.H
			i.spec.H.S = LS_SIZE_ABS
		}

		if sizeTop {
			newh = oldh - dy
		} else {
			newh = oldh + dy
		}

		if newh < 0 { // right border
			newh = 0
		} else if i.iframe.Y+i.iframe.H-newh < i.pframe.Y { // left border
			newh = i.iframe.Y + i.iframe.H - i.pframe.Y
		}
		//respect minsize of content or at least no negative spec
		if newh < mh {
			newh = mh
		}
		if sh == LS_SIZE_PCT {
			// convert back to percentage
			i.spec.H.S = LS_SIZE_PCT
			i.spec.H.V = utilGetPct(i.pframe.H, newh)
		} else {
			i.spec.H.V = newh
		}
		// prevent position changes
		if i.spec.Y.S == LS_POS_PCT {
			i.spec.Y.V = utilGetPct(i.pframe.H-newh, i.iframe.Y-i.pframe.Y)

		}
		if sizeTop {
			// position correction is needed
			dy = oldh - newh
			domove = true
			mvy = dy
		}
	}

	i.Layout(&i.pframe, true)
	// correct the position if it was the left or top border
	if domove {
		i.Move(mvx, mvy)
	}
	return neww - oldw, newh - oldh

}

// Moves that item by desired dx,dy and limits to layautParenFrame if set. Sets layout-specification to positioning. Returns the resulting movement.
func (i *ItemBase) Move(dx, dy int32) (int32, int32) {

	var newx, newy, oldx, oldy int32 //absolute positions

	// modify Layaut Specification by the possible amount in that direction
	// check if new absolute positions leave the bounding box
	if dx != 0 {

		//remember spec
		sx := i.spec.X.S

		// convert ls to positioning if it was not abs
		if i.spec.X.S != LS_POS_ABS {
			i.spec.X.V = i.iframe.X - i.pframe.X
			i.spec.X.S = LS_POS_ABS
		}

		oldx = i.pframe.X + i.spec.X.V
		newx = oldx + dx

		if newx < i.pframe.X { // left border
			newx = i.pframe.X
		} else if newx+i.iframe.W > i.pframe.X+i.pframe.W { // right border
			newx -= (newx + i.iframe.W) - (i.pframe.X + i.pframe.W)
		}

		// convert back to percentage
		if sx == LS_POS_PCT {
			i.spec.X.S = LS_POS_PCT
			i.spec.X.V = utilGetPct(i.pframe.W-i.iframe.W, newx-i.pframe.X)
		} else {
			i.spec.X.V = newx - i.pframe.X
		}

	}
	if dy != 0 {

		//remember spec
		sy := i.spec.Y.S

		// convert ls to positioning if it was not abs
		if i.spec.Y.S != LS_POS_ABS {
			i.spec.Y.V = i.iframe.Y - i.pframe.Y
			i.spec.Y.S = LS_POS_ABS
		}

		oldy = i.pframe.Y + i.spec.Y.V
		newy = oldy + dy

		if newy < i.pframe.Y { // top border
			newy = i.pframe.Y
		} else if newy+i.iframe.H > i.pframe.Y+i.pframe.H { // bottom border
			newy -= (newy + i.iframe.H) - (i.pframe.Y + i.pframe.H)
		}

		// convert back to percentage
		if sy == LS_POS_PCT {
			i.spec.Y.S = LS_POS_PCT
			i.spec.Y.V = utilGetPct(i.pframe.H-i.iframe.H, newy-i.pframe.Y)
		} else {
			i.spec.Y.V = newy - i.pframe.Y
		}
	}

	i.Layout(&i.pframe, false)

	return newx - oldx, newy - oldy

}
