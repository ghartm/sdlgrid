// sdlgrid2 project sdlgrid2.go
//package sdlgrid2

// TODO: listitem

// TODO: neues item viewport mit slidern

// TODO: anpassen von size() für alle ecken
// TODO: resize handles auch in allen ecken

// TODO: getminsize prüfen. bei absoluter positionierung auch die Position der subitems einbeziehen

// TODO: TextInput textmarkierung mit tastatur und maus

// TODO: Standard Icon set für menuepfeile und schliessen minimieren vergroessern

package main

import (
	"fmt"
	"reflect"
	"stopwatch"

	//	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

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

// ----------------------------------------------------------------

var sw *stopwatch.Stopwatch = stopwatch.NewStopwatch()

func newTestWindow(win *RootWindow) *ItemWindow {
	// Window
	win1 := NewItemWindow(win)
	win1.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_COLLAPSE, 20000, 50000, 50000, 0)

	rg := NewItemGrid(win, 0, 0)
	win1.AddSubItem(rg)

	sl := NewItemSlider(win, ITEM_SLIDER_STYLE_VERTICAL)
	rg.SetSubItem(rg.AppendColumn(), 0, sl)

	frm11 := NewItemFrame(win)
	rg.SetSubItem(rg.AppendColumn(), 0, frm11)

	rg.SetColSpec(0, layoutParam{LS_SIZE_ABS, win.GetStyle().decoUnit})

	frm11.SetName("wFrame")
	frm11.SetText("wFrame")
	frm11.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_COLLAPSE, LS_SIZE_COLLAPSE, 50000, 50000, 0, 0)

	gr3 := NewItemGrid(win, 0, 0)
	gr3.SetName("wGrid")
	gr3.SetSpacing(win.style.spacing)
	frm11.SetSubItem(gr3)

	ti1 := NewItemTextInput(win, win.style.averageRuneWidth*20)
	ti1.SetName("textInput1")

	ti2 := NewItemTextInput(win, 80)
	ti2.SetName("textInput2")
	gr3.SetSubItem(0, gr3.AppendRow(), ti1)
	gr3.SetSubItem(0, gr3.AppendRow(), ti2)

	b1 := NewItemButton(win)
	gr3.SetSubItem(0, gr3.AppendRow(), b1)
	b1.SetName("firstButton")
	b1.SetText("Change Alignment")
	b1.SetCallback(func() {
		fmt.Printf("Button Callback: %s\n", b1.GetName())
		ti1.SetTextAlignment(ITEM_TEXTINPUT_ALIGN_CENTER)
		ti2.SetTextAlignment(ITEM_TEXTINPUT_ALIGN_RIGHT)
	})

	b1 = NewItemButton(win)
	gr3.SetSubItem(0, gr3.AppendRow(), b1)
	b1.SetText("toggle Input")
	b1.SetCallback(func() {
		ti1.SetHidden(!ti1.hidden)
		ti2.SetHidden(!ti2.hidden)
	})

	gr4 := NewItemGrid(win, 0, 0)
	gr3.SetSubItem(0, gr3.AppendRow(), gr4)
	gr4.SetSpacing(win.style.spacing)

	lr := NewItemTextInput(win, 50)
	gr4.SetSubItem(gr4.AppendColumn(), 0, lr)
	lr.SetTextAlignment(ITEM_TEXTINPUT_ALIGN_RIGHT)
	lr.SetText("50")

	lg := NewItemTextInput(win, 50)
	gr4.SetSubItem(gr4.AppendColumn(), 0, lg)
	lg.SetTextAlignment(ITEM_TEXTINPUT_ALIGN_RIGHT)
	lg.SetText("50")

	lb := NewItemTextInput(win, 50)
	gr4.SetSubItem(gr4.AppendColumn(), 0, lb)
	lb.SetTextAlignment(ITEM_TEXTINPUT_ALIGN_RIGHT)
	lb.SetText("50")

	b1 = NewItemButton(win)
	gr3.SetSubItem(0, gr3.AppendRow(), b1)
	b1.SetText("set color")
	b1.SetCallback(func() {
		ncs := new(ColorScheme)
		ncs.SetBaseColor(&sdl.Color{R: uint8(lr.GetNumber()), G: uint8(lg.GetNumber()), B: uint8(lb.GetNumber()), A: 255})
		win1.oWithItems(func(i Item) {
			i.SetColorScheme(ncs)
			i.SetChanged(true)
		})
	})
	return win1
}

func main() {

	fmt.Println("start")
	sw.Start("main")

	wc := NewWindowController()
	defer wc.Destroy()

	s := NewStyle("default")
	defer s.Destroy()

	win := NewRootWindow(s, "RootWindow", 500, 300)
	defer win.Destroy()

	wc.AddRootWindow(win)

	gr1 := NewItemGrid(win, 1, 2)
	gr1.SetName("rootGrid")
	//gr1.SetRowSpec(0, 50)
	//gr1.SetRowSpec(1, 100)
	//gr1.SetSpacing(2)
	win.AddRootItem(gr1)

	// menue entries
	mu1 := NewItemMenue(win)
	gr1.SetSubItem(0, 0, mu1)
	mu1.SetName("mu1")
	mu1.SetOrientation(MENU_ORIENTATION_HORIZONTAL)

	wm := NewItemWindowManager(win)
	gr1.SetSubItem(0, 1, wm)

	mu2 := NewItemMenue(win)
	mu1.AddNewMenuEntry("File", mu2, nil)
	mu2.SetName("mu2")
	mu2.AddNewMenuEntry("New", nil, nil)
	mu2.AddNewMenuEntry("Open", nil, nil)
	mu2.AddNewMenuEntry("Save", nil, nil)
	mu2.AddNewMenuEntry("Close", nil, nil)

	mu3 := NewItemMenue(win)
	mu1.AddNewMenuEntry("Window", mu3, nil)
	mu3.AddNewMenuEntry("New", nil, func() {
		wm.AddSubItem(newTestWindow(win))

	})

	w1 := newTestWindow(win)
	wm.AddSubItem(w1)

	win.ReportItems()
	wc.Start()

	sw.Stop("main")
	sw.Report()

}

func eventSenderCustom(wg *sync.WaitGroup, c chan interface{}, ctrl chan bool) {
	wg.Add(1)
	var e EventCustom
	var running bool = true
	for running {
		select {
		case _, running = <-ctrl:
			fmt.Println("quitting eventSenderCustom")
			wg.Done()
		case c <- &e:
			time.Sleep(time.Duration(1000) * time.Millisecond)
		}
	}
}

func eventSenderTick(wg *sync.WaitGroup, c chan interface{}, ctrl chan bool, msec int) {
	wg.Add(1)
	e := EventRenderTick{msec: msec}
	var running bool = true
	for running {
		select {
		case _, running = <-ctrl:
			fmt.Println("quitting eventSenderTick")
			wg.Done()
		default:
			c <- &e
			time.Sleep(time.Duration(msec) * time.Millisecond)
		}
	}
}

//--------------------------------------------------------------------

type EventRenderTick struct {
	msec int
}

//--------------------------------------------------------------------

type EventCustom struct {
	round int
}

// --------------------------------------------------------------------
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
	oNotifyMouseButton(x int32, y int32, button sdl.Button, state sdl.ButtonState) // state:0=released 1=pressed
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
func (i *ItemBase) oNotifyMouseButton(x, y int32, button sdl.Button, state sdl.ButtonState) {
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

// --------------------------------------------------------------------
// Container item that manages/arranges its subitems in a grid. It has no decoration
type ItemGrid struct {
	ItemBase
	cols                int
	rows                int
	spacing             int32
	grid                [][]Item // addressed by (column;row)
	rowSpec             []layoutParam
	colSpec             []layoutParam
	lastFoundSubItemCol int
	lastFoundSubItemRow int
	buttonCallback      func(int32, int32, sdl.Button, sdl.ButtonState) //x, y int32, button, state uint8
}

func NewItemGrid(win *RootWindow, cols int, rows int) *ItemGrid {
	i := new(ItemGrid)
	i.o = Item(i)
	i.setRootWindow(win)
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_PCT, 0, 0, 100000, 100000)

	i.cols = cols
	i.rows = rows

	i.rowSpec = make([]layoutParam, i.rows)
	i.colSpec = make([]layoutParam, i.cols)
	i.grid = make([][]Item, i.cols)

	var c, r int

	for c = 0; c < i.cols; c++ {
		i.grid[c] = make([]Item, i.rows)
	}

	// default is to collapse all but the last
	if i.cols > 0 {
		for c = 0; c < i.cols-1; c++ {
			i.colSpec[c] = layoutParam{LS_SIZE_COLLAPSE, 0}
		}
		i.colSpec[c] = layoutParam{LS_SIZE_PCT, 100000}
	}
	if i.rows > 0 {
		for r = 0; r < i.rows-1; r++ {
			i.rowSpec[r] = layoutParam{LS_SIZE_COLLAPSE, 0}

		}
		i.rowSpec[r] = layoutParam{LS_SIZE_PCT, 100000}

	}

	return i
}

func (i *ItemGrid) SetButtonCallback(cb func(int32, int32, sdl.Button, sdl.ButtonState)) {
	i.buttonCallback = cb
}
func (i *ItemGrid) oNotifyMouseButton(x, y int32, button sdl.Button, state sdl.ButtonState) {
	fmt.Printf("%s: ItemGrid.oNotifyMouseButton()\n", i.GetName())
	if i.buttonCallback != nil {
		i.buttonCallback(x, y, button, state)
	}
}

// extend grid by one column. return index of new column
func (i *ItemGrid) AppendColumn() int {
	i.grid = append(i.grid, make([]Item, i.rows))

	i.colSpec = append(i.colSpec, layoutParam{LS_SIZE_PCT, 100000})
	i.cols = len(i.colSpec)
	if i.cols > 1 {
		i.colSpec[i.cols-1] = layoutParam{LS_SIZE_COLLAPSE, 0}
	}

	// if it was the first column add a row as well
	if i.cols == 1 {
		i.grid[i.cols-1] = append(i.grid[i.cols-1], nil)
		i.rowSpec = append(i.rowSpec, layoutParam{LS_SIZE_PCT, 100000})
		i.rows = len(i.rowSpec)
	}

	return i.cols - 1
}

// extend grid by one row. return index of new row
func (i *ItemGrid) AppendRow() int {

	// if there is no column jet, add a column first
	if i.cols == 0 {
		i.AppendColumn()
		//if it was the first column a row was added as well
	} else {

		for n := range i.grid {
			i.grid[n] = append(i.grid[n], nil)
		}
		i.rowSpec = append(i.rowSpec, layoutParam{LS_SIZE_PCT, 100000})
	}

	i.rows = len(i.rowSpec)
	if i.rows > 1 {
		i.rowSpec[i.rows-1] = layoutParam{LS_SIZE_COLLAPSE, 0}
	}

	return i.rows - 1
}

// set the inner spacing between the cells
func (i *ItemGrid) SetSpacing(spc int32) {
	i.spacing = spc
	i.useTmpMinSize = false
}

// set the horizontal layout spec for a column
func (i *ItemGrid) SetColSpec(col int, p layoutParam) {
	i.colSpec[col] = p
	i.useTmpMinSize = false
}

// set the vertical layout spec for a row
func (i *ItemGrid) SetRowSpec(row int, p layoutParam) {
	i.rowSpec[row] = p
	i.useTmpMinSize = false
}

// places a sub-item into a cell
func (i *ItemGrid) SetSubItem(col int, row int, si Item) bool {
	if col < i.cols && row < i.rows {
		i.grid[col][row] = si
		si.SetParent(i)
		i.useTmpMinSize = false
		return true
	}
	return false
}

func (i *ItemGrid) oGetSubFrame() *sdl.Rect { return &i.iframe }

// get minimum sizes of items in cells ( column-width; row-heigt)
func (i *ItemGrid) getMinSizeFields() (cw, rh []int32) {

	cw = make([]int32, i.cols)
	rh = make([]int32, i.rows)
	for c := 0; c < i.cols; c++ {
		for r := 0; r < i.rows; r++ {
			if i.grid[c][r] != nil {
				x, y, mw, mh := i.grid[c][r].GetCollapsedSpec()
				if i.grid[c][r].IsAutoSize() {
					//if it is not a fixed size - Minimum size must be computed
					mw, mh = i.grid[c][r].oGetMinSize()
				}

				mw += x
				mh += y

				// if absolute size
				if i.colSpec[c].S == LS_SIZE_ABS {
					mw = i.colSpec[c].V
				}

				if i.rowSpec[r].S == LS_SIZE_ABS {
					mh = i.rowSpec[r].V
				}

				if cw[c] < mw {
					cw[c] = mw
				}
				if rh[r] < mh {
					rh[r] = mh
				}
			}
		}
	}

	return cw, rh

}
func (i *ItemGrid) oNotifyPostLayout(sizeChanged bool) {
	//fmt.Printf("%s: ItemGrid.oNotifyPostLayout()\n", i.GetName())
	// compute frame-size for every cell and layout the subitem in each frame
	// get required minimum size from every cell. cw:column-width rh:row-hight min for row and col
	//column-width and row-height

	cw, rh := i.getMinSizeFields()

	//TODO  make grid layout aware to sizeChanged

	// get minimum needed space for grid Frame
	// and count the number of expand-specs if it has content

	var sum int32
	var npct int32
	var basepct int32

	// distribute free space between required size and size of grid-item frame across the expandable cells according to the percentages
	for n := 0; n < i.cols; n++ {
		sum += cw[n]
		if i.colSpec[n].S == LS_SIZE_PCT {
			npct++
			basepct += i.colSpec[n].V
		}
	}
	// if there is space for distribution and percentage layout exists
	if free := (i.iframe.W - (sum + (i.spacing * int32(i.cols-1)))); free > 0 && npct > 0 {
		for n := 0; n < i.cols; n++ {
			if i.colSpec[n].S == LS_SIZE_PCT {
				// normalize percentages
				cw[n] += utilPct(free, utilNormPct(basepct, i.colSpec[n].V))
			}
		}
	}
	//-----------

	npct = 0
	basepct = 0
	sum = 0
	for n := 0; n < i.rows; n++ {
		sum += rh[n]
		if i.rowSpec[n].S == LS_SIZE_PCT {
			npct++
			basepct += i.rowSpec[n].V
		}
	}

	// if there is space for distribution and percentage layout exists
	if free := (i.iframe.H - (sum + (i.spacing * int32(i.rows-1)))); free > 0 && npct > 0 {
		for n := 0; n < i.rows; n++ {
			if i.rowSpec[n].S == LS_SIZE_PCT {
				// normalize percentages
				rh[n] += utilPct(free, utilNormPct(basepct, i.rowSpec[n].V))
			}
		}
	}

	// layout subitems in the cells according to cw and rh
	var pf sdl.Rect //parent frame of sub-item
	pf.X = i.iframe.X
	for c := 0; c < i.cols; c++ {
		pf.Y = i.iframe.Y
		pf.W = cw[c]
		for r := 0; r < i.rows; r++ {
			pf.H = rh[r]
			if i.grid[c][r] != nil {
				i.grid[c][r].Layout(&pf, sizeChanged)
			}
			pf.Y += rh[r] + i.spacing
		}
		pf.X += cw[c] + i.spacing
	}
}

func (i *ItemGrid) oReportSubitems(lvl int) {

	for c := 0; c < i.cols; c++ {
		for r := 0; r < i.rows; r++ {
			for n := 0; n < lvl; n++ {
				fmt.Print("--")
			}

			if i.grid[c][r] != nil {
				fmt.Printf("%s: c:%d r:%d\n", i.GetName(), c, r)
				i.grid[c][r].Report(lvl)
			} else {
				fmt.Printf("%s: c:%d r:%d empty\n", i.GetName(), c, r)
			}
		}
	}
}
func (i *ItemGrid) oWithItems(fn func(Item)) {
	fmt.Printf("%s: ItemBase.oWithItems() no subitems\n", i.GetName())
	// handle the call for myself
	fn(i)
	// forward to all subitems
	for c := 0; c < i.cols; c++ {
		for r := 0; r < i.rows; r++ {
			if i.grid[c][r] != nil {
				i.grid[c][r].oWithItems(fn)
			}
		}
	}
}
func (i *ItemGrid) oGetMinSize() (int32, int32) {
	//fmt.Printf("%s: ItemGrid.oGetMinSize()\n", i.GetName())
	if !i.useTmpMinSize {
		//i.minw, i.minh

		//column-width and row-height
		if i.cols < 1 || i.rows < 1 {
			return 0, 0
		} else {

			cw, rh := i.getMinSizeFields()

			// get minimum needed space for grid Frame
			var cwsum, rhsum int32
			for c := 0; c < i.cols; c++ {
				cwsum += cw[c]
			}
			for r := 0; r < i.rows; r++ {
				rhsum += rh[r]
			}
			i.minw = cwsum + (i.spacing * int32(i.cols-1))
			i.minh = rhsum + (i.spacing * int32(i.rows-1))
			fmt.Printf("%s: ItemGrid.oGetMinSize() w:%d h:%d\n", i.GetName(), i.minw, i.minh)
		}
		i.useTmpMinSize = true
	}
	return i.minw, i.minh
}

func (i *ItemGrid) oRender() {
	//fmt.Printf("%s: ItemGrid.oRender()\n", i.GetName())
	// ItemGrid does not have a decoration.

	/*// debug frame
	rd := i.GetRenderer()
	s := i.GetStyle()
	rd.SetDrawBlendMode(sdl.BLENDMODE_NONE)
	utilRenderSolidBorder(rd, &i.iframe, s.colorRed)
	*/

	// forward to content
	for c := 0; c < i.cols; c++ {
		for r := 0; r < i.rows; r++ {
			if i.grid[c][r] != nil {
				i.grid[c][r].Render()
			}
		}
	}

}

func (i *ItemGrid) oFindSubItem(x, y int32, e sdl.Event) (found bool, item Item) {
	if i.cols < 1 || i.rows < 1 {
		return false, nil
	}
	// first check last hit
	if (i.lastFoundSubItemCol >= 0) && (i.lastFoundSubItemCol < i.cols) && (i.lastFoundSubItemRow < i.rows) {
		if si := i.grid[i.lastFoundSubItemCol][i.lastFoundSubItemRow]; si != nil {
			if si.CheckPos(x, y) {
				return true, si
			}
		}
	}

	// check sequential if not found
	for c := 0; c < i.cols; c++ {
		for r := 0; r < i.rows; r++ {
			if si := i.grid[c][r]; si != nil {
				if si.CheckPos(x, y) {
					i.lastFoundSubItemCol = c
					i.lastFoundSubItemRow = r
					return true, si
				}
			}
		}
	}
	return false, nil
}

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
		i.textBox.SetSpecX(LS_POS_ABS, i.textBox.GetStyle().spacing*2)
		i.textBox.UseBaseColor(true)
	}
	i.textBox.SetText(t)
	return i
}
func (i *ItemFrame) getBorderValues() (t, r, b, l int32) {
	// frame line
	s := i.GetStyle().spacing
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

	//utilRenderSolidBorder(r, &i.outerFrame, s.colorGreen)
	//utilRenderSolidBorder(r, &i.outerFrame, i.GetStyle().colorGreen)

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

// -----------------------------------------------------------------------------
type ItemSlider struct {
	ItemBase
	upLeftCb        func()
	downRightCb     func()
	valueSetCb      func()
	rootGrid        *ItemGrid
	upLeftButton    *ItemButton
	handle          *ItemHandle
	downRightButton *ItemButton
	sliderStyle     int

	userUnit  int32 // units to use for object movements
	userView  int32 // viewable size of object
	userTotal int32 // total size of object to slide

	userTotalPerPx float32

	handleSize int32 // size of handle accoring viewable size of object
}

func NewItemSlider(win *RootWindow, sl int) *ItemSlider {
	i := new(ItemSlider)
	i.o = Item(i)
	i.setRootWindow(win)
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_PCT, 0, 0, 100000, 100000)
	i.sliderStyle = sl
	i.SetName("Slider")

	i.rootGrid = NewItemGrid(win, 0, 0)
	i.rootGrid.SetParent(i)

	i.upLeftButton = NewItemButton(win)
	i.upLeftButton.SetItemStyle(ITEM_BUTTON_STYLE_FLAT)
	i.upLeftButton.SetActOnClick(true)
	i.upLeftButton.SetCallback(i.upLeftButtonCb)

	i.downRightButton = NewItemButton(win)
	i.downRightButton.SetItemStyle(ITEM_BUTTON_STYLE_FLAT)
	i.downRightButton.SetActOnClick(true)
	i.downRightButton.SetCallback(i.downRightButtonCb)

	i.handle = NewItemHandle(win)
	i.handle.SetTargetCb(i.handleMoveCb)
	i.handle.SetItemStyle(ITEM_HANDLE_STYLE_FLAT)

	size := win.GetStyle().decoUnit
	i.userUnit = size

	if i.sliderStyle == ITEM_SLIDER_STYLE_HORIZONTAL {
		i.upLeftButton.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_ABS, LS_SIZE_PCT, 0, 0, size, 0)
		i.downRightButton.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_ABS, LS_SIZE_PCT, 100000, 0, size, 100000)
		i.handle.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_ABS, LS_SIZE_PCT, 0, 0, size, 100000)

		n := i.rootGrid.AppendColumn()
		i.rootGrid.SetSubItem(n, 0, i.downRightButton)
		i.rootGrid.SetColSpec(n, layoutParam{LS_SIZE_COLLAPSE, 0})
		n = i.rootGrid.AppendColumn()
		i.rootGrid.SetSubItem(n, 0, i.handle)
		i.rootGrid.SetColSpec(n, layoutParam{LS_SIZE_PCT, 100000})
		n = i.rootGrid.AppendColumn()
		i.rootGrid.SetSubItem(n, 0, i.downRightButton)
		i.rootGrid.SetColSpec(n, layoutParam{LS_SIZE_COLLAPSE, 0})
	} else {
		i.upLeftButton.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_ABS, 0, 0, 100000, size)
		i.downRightButton.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_ABS, 0, 100000, 100000, size)
		i.handle.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_ABS, 0, 0, 100000, size)

		n := i.rootGrid.AppendRow()
		i.rootGrid.SetSubItem(0, n, i.upLeftButton)
		i.rootGrid.SetRowSpec(n, layoutParam{LS_SIZE_COLLAPSE, 0})
		n = i.rootGrid.AppendRow()
		i.rootGrid.SetSubItem(0, n, i.handle)
		i.rootGrid.SetRowSpec(n, layoutParam{LS_SIZE_PCT, 100000})
		n = i.rootGrid.AppendRow()
		i.rootGrid.SetSubItem(0, n, i.downRightButton)
		i.rootGrid.SetRowSpec(n, layoutParam{LS_SIZE_COLLAPSE, 0})
	}
	return i
}

func (i *ItemSlider) computeUnitSize() {
	// according handles ability to move, handle size and total and part values
	// compute the jump unit in pixel tha handle will jump on a decrease or increase
	var pix int32
	if i.sliderStyle == ITEM_SLIDER_STYLE_HORIZONTAL {
		pix = i.handle.pframe.W
	} else {
		pix = i.handle.pframe.H
	}
	i.userTotalPerPx = float32(i.userTotal) / float32(pix)

	s := i.GetStyle().decoUnit
	if i.handleSize = int32(float32(i.userView) / i.userTotalPerPx); i.handleSize < s {
		i.handleSize = s
	}

	if i.sliderStyle == ITEM_SLIDER_STYLE_HORIZONTAL {
		pix = i.handle.pframe.W
	} else {
		pix = i.handle.pframe.H
	}

}

func (i *ItemSlider) SetUnitRange(unit, total, view int32) {
	i.userView = view
	i.userTotal = total

}

func (i *ItemSlider) downRightButtonCb() {
	i.moveHandle(i.userUnit, i.userUnit)
	if i.downRightCb != nil {
		i.downRightCb()
	}
}

func (i *ItemSlider) upLeftButtonCb() {
	i.moveHandle(-i.userUnit, -i.userUnit)
	if i.upLeftCb != nil {
		i.upLeftCb()
	}
}

func (i *ItemSlider) handleMoveCb(dx, dy int32) {
	i.moveHandle(dx, dy)
}

func (i *ItemSlider) moveHandle(dx, dy int32) {
	if i.sliderStyle == ITEM_SLIDER_STYLE_HORIZONTAL {
		i.handle.Move(dx, 0)
	} else {
		i.handle.Move(0, dy)
	}
	//i.handle.Layout(&i.handle.pframe, false)
	i.SetChanged(true)
	if i.valueSetCb != nil {
		i.valueSetCb()
	}
}

func (i *ItemSlider) oRender() {
	r := i.GetRenderer()
	c := i.GetColorScheme()

	// background of slider
	utilRenderFillRect(r, &i.iframe, &c.baseDark)
	// content
	i.rootGrid.Render()
}
func (i *ItemSlider) oReportSubitems(lvl int) {
	i.rootGrid.Report(lvl)
}
func (i *ItemSlider) oNotifyPostLayout(sizeChanged bool) {

	i.rootGrid.Layout(i.oGetSubFrame(), sizeChanged)

	if sizeChanged {
		i.computeUnitSize()
		//find the new position for the handle
	}
}
func (i *ItemSlider) oGetSubFrame() *sdl.Rect {
	return &i.iframe
}
func (i *ItemSlider) oGetMinSize() (int32, int32) {
	return i.rootGrid.oGetMinSize()
}
func (i *ItemSlider) oWithItems(fn func(Item)) {
	fn(i)
	i.rootGrid.oWithItems(fn)
}
func (i *ItemSlider) oFindSubItem(x, y int32, e sdl.Event) (bool, Item) {
	if i.rootGrid.CheckPos(x, y) {
		return true, i.rootGrid
	}
	return false, nil
}

//-----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
type ItemBackground struct {
	ItemBase
	subItem Item
	margin  int32
}

func NewItemBackground(win *RootWindow) *ItemBackground {
	i := new(ItemBackground)
	i.o = Item(i)
	i.setRootWindow(win)
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_PCT, 0, 0, 100000, 100000)
	//i.margin=win.style.spacing
	return i
}

func (i *ItemBackground) SetMargin(m int32) { i.margin = m }

func (i *ItemBackground) SetSubItem(s Item) {
	i.subItem = s
	if i.subItem != nil {
		s.SetParent(i)
	}
}

func (i *ItemBackground) oRender() {
	r := i.GetRenderer()
	c := i.GetColorScheme()

	// background of window
	utilRenderFillRect(r, &i.iframe, &c.baseDark)

	// content
	if i.subItem != nil {
		i.subItem.Render()
	}

}
func (i *ItemBackground) oReportSubitems(lvl int) {
	i.subItem.Report(lvl)

}
func (i *ItemBackground) oNotifyPostLayout(sizeChanged bool) {
	i.subItem.Layout(i.oGetSubFrame(), sizeChanged)
}

func (i *ItemBackground) oGetSubFrame() *sdl.Rect {
	return i.MakeInnerFrame(i.margin)
}

func (i *ItemBackground) oGetMinSize() (int32, int32) {
	w, h := i.subItem.oGetMinSize()
	spc := i.margin * 2
	return w + spc, h + spc
}
func (i *ItemBackground) oWithItems(fn func(Item)) {
	fmt.Printf("%s: ItemBackground.oWithItems()\n", i.GetName())
	// handle the call for myself
	fn(i)
	// forward to all subitems
	i.subItem.oWithItems(fn)
}
func (i *ItemBackground) oFindSubItem(x, y int32, e sdl.Event) (bool, Item) {
	if i.subItem.CheckPos(x, y) {
		return true, i.subItem
	}
	return false, nil
}

// -----------------------------------------------------------------------------
type ItemButton struct {
	ItemBase
	buttonState bool   // false=unpressed true=pressed
	buttonStyle int    // flat/ edge ...
	userCb      func() // callback function if button is pressed
	textBox     *ItemText
	mouseFocus  bool
	actOnClick  bool
}

func NewItemButton(win *RootWindow) *ItemButton {
	i := new(ItemButton)
	i.o = Item(i)
	i.setRootWindow(win)
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_COLLAPSE, LS_SIZE_COLLAPSE, 0, 0, 0, 0)
	i.buttonStyle = ITEM_BUTTON_STYLE_EDGE
	return i
}

func (i *ItemButton) oReportSubitems(lvl int) {

	if i.textBox != nil {
		i.textBox.Report(lvl)
	}
}

func (i *ItemButton) SetItemStyle(s int)   { i.buttonStyle = s }
func (i *ItemButton) SetActOnClick(b bool) { i.actOnClick = b }

func (i *ItemButton) SetText(t string) *ItemButton {
	if i.textBox == nil {
		i.textBox = NewItemText(i.win)
		i.textBox.SetParent(i)
		i.textBox.SetName(i.name + ".textBox")
		i.textBox.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_COLLAPSE, LS_SIZE_COLLAPSE, 50000, 50000, 0, 0)
	}
	i.textBox.SetText(t)
	return i
}
func (i *ItemButton) SetCallback(f func()) {
	i.userCb = f
}
func (i *ItemButton) oWithItems(fn func(Item)) {
	fmt.Printf("%s: ItemBase.oWithItems() no subitems\n", i.GetName())
	// handle the call for myself
	fn(i)
	// forward to all subitems
	if i.textBox != nil {
		i.textBox.oWithItems(fn)
	}
}
func (i *ItemButton) oNotifyPostLayout(sizeChanged bool) {
	//fmt.Printf("%s: ItemButton.oNotifyPostLayout()\n", i.GetName())

	// layout decoration
	if i.textBox != nil {
		i.textBox.Layout(i.GetFrame(), sizeChanged)
	}
}

func (i *ItemButton) oGetMinSize() (int32, int32) {
	// get the minimum bounding box for myself
	if !i.useTmpMinSize { //i.minw, i.minh
		var w, h, mw, mh int32

		//at least decoration and margin =border
		mgn := i.GetStyle().spacing

		if i.textBox != nil {
			w, h = i.textBox.oGetMinSize()
			mw, mh = utilSizeMax(mw, mh, w, h)
		}

		_, _, w, h = i.GetCollapsedSpec()
		w -= mgn * 2
		h -= mgn * 2
		mw, mh = utilSizeMax(mw, mh, w, h)

		mw += mgn * 2
		mh += mgn * 2

		i.minw = mw
		i.minh = mh
		i.useTmpMinSize = true
	}
	return i.minw, i.minh
}

func (i *ItemButton) oRender() {
	//fmt.Printf("%s: ItemButton.oRender()\n", i.GetName())

	r := i.GetRenderer()
	c := i.GetColorScheme()
	r.SetDrawBlendMode(sdl.BLENDMODE_NONE)
	f := &i.iframe

	// background

	switch i.buttonStyle {
	case ITEM_BUTTON_STYLE_EDGE:
		// shadow border
		utilRenderFillRect(r, f, &c.base)
		utilRenderShadowBorder(r, f, c, i.buttonState)
	case ITEM_BUTTON_STYLE_FLAT:
		if i.buttonState {
			utilRenderFillRect(r, f, &c.baseBright)
		} else {
			utilRenderFillRect(r, f, &c.base)
		}

		if i.mouseFocus {
			utilRenderSolidBorder(r, f, &c.baseReverse)
		}
	}

	// render the textbox
	if i.textBox != nil {
		i.textBox.oRender()
	}
}

func (i *ItemButton) oNotifyMouseFocusGained() {
	fmt.Printf("%s: ItemButton.oNotifyMouseFocusGained()\n", i.GetName())
	i.mouseFocus = true
	i.win.SetChanged(true)
}

func (i *ItemButton) oNotifyMouseFocusLost() {
	fmt.Printf("%s: ItemButton.oNotifyMouseFocusLost()\n", i.GetName())
	i.SetButtonState(false)
	i.mouseFocus = false
	i.win.SetChanged(true)
}

func (i *ItemButton) SetButtonState(s bool) {
	if i.buttonState != s {
		i.buttonState = s
		switch i.buttonStyle {
		case ITEM_BUTTON_STYLE_EDGE:

			if i.textBox != nil {
				// if there is text move it
				fx, fy := i.textBox.GetFramePos()
				if s {
					// sunken text
					i.textBox.iframe.X = fx + 1
					i.textBox.iframe.Y = fy + 1

				} else {
					//normal text
					i.textBox.iframe.X = fx - 1
					i.textBox.iframe.Y = fy - 1
				}
			}

			i.SetChanged(true)
		case ITEM_BUTTON_STYLE_FLAT:
		}

	}

}

func (i *ItemButton) oNotifyMouseButton(x, y int32, button sdl.Button, state sdl.ButtonState) {
	fmt.Printf("%s: ItemButton.oNotifyMouseButton()\n", i.GetName())
	if button == 1 {
		switch state {
		case 1:
			// button 1 pressed

			if i.actOnClick {

				if i.userCb != nil {
					i.userCb()
				}
			} else {
				i.SetButtonState(true)
			}
		case 0:
			// button 1 released

			if i.buttonState {
				// call usercb if button was in state pressed before only
				if !i.actOnClick {
					if i.userCb != nil {
						i.userCb()
					}
				}
			}
			i.SetButtonState(false)
		}
	}
	i.SetChanged(true)
}

// -----------------------------------------------------------------------------
// area that can be picked and dragged by mouse to move a specific Target item. Initially itself.
type ItemHandle struct {
	ItemBase
	picked       bool
	pickedx      int32 //relative to item outerframe. Where in item was the pick.
	pickedy      int32 //relative to item outerframe
	targetCb     func(int32, int32)
	tmpCursor    *sdl.Cursor
	cursor       *sdl.Cursor
	cursorPicked *sdl.Cursor
	handleStyle  int
}

func NewItemHandle(win *RootWindow) *ItemHandle {
	i := new(ItemHandle)
	i.o = Item(i)
	i.setRootWindow(win)
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_PCT, 0, 0, 100000, 100000)
	i.handleStyle = ITEM_HANDLE_STYLE_HIDDEN

	return i
}

func (i *ItemHandle) SetItemStyle(s int) {
	i.handleStyle = s
}
func (i *ItemHandle) SetCursor(c sdl.SystemCursor) {
	i.cursor = sdl.CreateSystemCursor(c)
	if i.cursorPicked == nil {
		i.cursorPicked = i.cursor
	}
}
func (i *ItemHandle) SetCursorPicked(c sdl.SystemCursor) {
	i.cursorPicked = sdl.CreateSystemCursor(c)
}

func (i *ItemHandle) SetTargetCb(cb func(int32, int32)) {
	i.targetCb = cb
}

func (i *ItemHandle) oNotifyMouseMotion(x, y, dx, dy int32) {
	fmt.Printf("ItemHandle: %s.oNotifyMouseMotion(%d,%d,%d,%d)\n", i.name, x, y, dx, dy)
	if i.picked {
		if i.targetCb != nil {
			// send dx,dy to come back to picked position
			i.targetCb((x-i.iframe.X)-i.pickedx, (y-i.iframe.Y)-i.pickedy)
		}
	}
}
func (i *ItemHandle) oNotifyMouseFocusLost() {
	fmt.Printf("ItemHandle: %s.oNotifyMouseFocusLost()\n", i.name)
	i.pick(false)
	if i.tmpCursor != nil {
		sdl.SetCursor(i.tmpCursor)
	}
	i.tmpCursor = nil
}
func (i *ItemHandle) oNotifyMouseFocusGained() {
	fmt.Printf("ItemHandle: %s.oNotifyMouseFocusGained()\n", i.name)
	i.pick(false)
	i.tmpCursor = sdl.GetCursor()
	if i.cursor != nil {
		sdl.SetCursor(i.cursor)
	}
}

func (i *ItemHandle) pick(b bool) {
	fmt.Printf("ItemHandle: %s.pick(%t)\n", i.name, b)
	if i.picked != b {
		i.picked = b
		if i.picked {
			if i.cursorPicked != nil {
				sdl.SetCursor(i.cursorPicked)
			}
		} else {
			if i.cursor != nil {
				sdl.SetCursor(i.cursor)
			} else {
				sdl.SetCursor(i.tmpCursor)
			}
		}

		i.win.SetMouseFocusLock(b)
	}
}

func (i *ItemHandle) oNotifyMouseButton(x, y int32, button sdl.Button, state sdl.ButtonState) {
	fmt.Printf("%s: ItemHandle.oNotifyMouseButton()\n", i.GetName())
	if button == 1 {
		switch state {
		case 1:
			// button 1 pressed
			i.pick(true)
			i.pickedx = x - i.iframe.X
			i.pickedy = y - i.iframe.Y
		case 0:
			// button 1 released
			i.pick(false)
		}
	}
}
func (i *ItemHandle) oRender() {
	//fmt.Printf("%s: ItemHandle.oRender()\n", i.GetName())
	r := i.GetRenderer()

	//utilRenderSolidBorder(r, &i.iframe, i.GetStyle().colorPurple)

	switch i.handleStyle {
	case ITEM_HANDLE_STYLE_DARKEDGE:
		r.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
		utilRendererSetDrawColor(r, &sdl.Color{R: 0, G: 0, B: 0, A: 30})
		//r.FillRect(&i.iframe)
		r.DrawRect(&i.iframe)
	case ITEM_HANDLE_STYLE_FLAT:
		c := i.GetColorScheme()
		utilRendererSetDrawColor(r, &c.base)
		r.FillRect(&i.iframe)
		utilRendererSetDrawColor(r, &c.lowEdge)
		r.DrawRect(&i.iframe)
	}

}

// -----------------------------------------------------------------------------
// Display a text
type ItemText struct {
	ItemBase
	text         string // Buttons text
	textTexture  *sdl.Texture
	textureColor sdl.Color
	useBaseColor bool
}

func NewItemText(win *RootWindow) *ItemText {
	i := new(ItemText)
	i.o = Item(i)
	i.setRootWindow(win)
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_COLLAPSE, LS_SIZE_COLLAPSE, 0, 0, 0, 0)
	return i
}
func (i *ItemText) UseBaseColor(b bool) { i.useBaseColor = b }

func (i *ItemText) oGetMinSize() (int32, int32) {
	i.prepareTexture()
	return i.GetSize()
}

func (i *ItemText) oRender() {

	//fmt.Printf("%s: ItemText.oRender()\n", i.GetName())
	r := i.GetRenderer()

	//utilRenderSolidBorder(r, &i.iframe, i.style.colorPurple)

	if i.textTexture != nil {
		if i.textureColor.Uint32() != i.GetColorScheme().text.Uint32() {
			i.textTexture.Destroy()
			i.textTexture = nil
			i.prepareTexture()
		}
		r.Copy(i.textTexture, nil, &i.iframe)
	}

}

func (i *ItemText) SetText(s string) {
	i.text = s
	if i.textTexture != nil {
		i.textTexture.Destroy()
		i.textTexture = nil
	}
	i.prepareTexture()
}

// prepares and renders the Texture of the Text
// Textures need to be prepared for layout.
func (i *ItemText) prepareTexture() {
	//fmt.Printf("ItemText.prepareTexture()\n")
	if i.textTexture == nil {
		if i.text != "" {
			if i.useBaseColor {
				i.textureColor = i.GetColorScheme().base
			} else {
				i.textureColor = i.GetColorScheme().text
			}
			sfc, err := i.style.Font.RenderUTF8Blended(i.text, i.textureColor)
			//surface, _ := f.RenderUTF8Shaded(t, *c)
			if err != nil {
				panic(err)
			}
			i.textTexture, _ = i.GetRenderer().CreateTextureFromSurface(sfc)
			i.SetSpecSize(LS_POS_ABS, LS_POS_ABS, sfc.W, sfc.H)
			// layout
			i.iframe.W = sfc.W
			i.iframe.H = sfc.H
			sfc.Free()
		} else {
			i.SetSpecSize(LS_POS_ABS, LS_POS_ABS, 0, 0)
			//layout
			i.iframe.W = 0
			i.iframe.H = 0

		}
	}
}

//--------------------------------------------------------------------

type ItemTextInput struct {
	ItemBase
	textBox         *ItemText
	text            []rune   // the text itself as a slice of runes; r := []rune{'\u0061', '\u0062', '\u0063', '\u65E5', -1}; s := string(r)
	cursorPosition  []int32  // cursor pixel-position for each gap between runes in texture. referres to x of uter-frame of textbox textBox
	cursorLocation  int      // location of cursor in text. 0=pos1; 1= after first rune ...
	cursorVisible   bool     // used for blinking curser switched on and off
	cursorRect      sdl.Rect // the cursor itself
	active          bool     // input is able to receive kb focus
	textWindowStart int      // start rune position of visible text in textBox
	textWindowEnd   int      // end rune position of visible text in textBox
	textAlign       int32    // alignment of the text LS_POS_PCT/CENTER/RIGHT
	picked          bool
	pickedx         int32 //relative to item outerframe. Where in item was the pick.
	pickedy         int32 //relative to item outerframe
	tmpCursor       *sdl.Cursor
	cursor          *sdl.Cursor
}

func NewItemTextInput(win *RootWindow, width int32) *ItemTextInput {
	i := new(ItemTextInput)
	i.o = Item(i)
	i.setRootWindow(win)

	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_ABS, LS_SIZE_ABS, 0, 0, width, (i.style.spacing*2 + i.style.averageFontHight))
	i.active = false
	i.textAlign = ITEM_TEXTINPUT_ALIGN_LEFT

	i.textBox = NewItemText(win)
	i.textBox.SetName(i.name + ".textBox")

	i.cursorRect.Y = i.style.spacing
	i.cursorRect.W = 1
	i.cursorRect.H = i.style.averageFontHight

	i.SetCursor(sdl.SYSTEM_CURSOR_IBEAM)

	return i
}

func (i *ItemTextInput) SetCursor(c sdl.SystemCursor) {
	i.cursor = sdl.CreateSystemCursor(c)
}

func (i *ItemTextInput) pick(b bool) {
	fmt.Printf("ItemHandle: %s.pick(%t)\n", i.name, b)
	if i.picked != b {
		i.picked = b
		if i.picked {

		} else {

		}

		i.win.SetMouseFocusLock(b)
	}
}

func (i *ItemTextInput) SetTextAlignment(a int32) {
	i.textAlign = a
	switch i.textAlign {
	case ITEM_TEXTINPUT_ALIGN_LEFT:
		i.textBox.SetSpecX(LS_POS_PCT, 0)
	case ITEM_TEXTINPUT_ALIGN_RIGHT:
		i.textBox.SetSpecX(LS_POS_PCT, 100000)
	case ITEM_TEXTINPUT_ALIGN_CENTER:
		i.textBox.SetSpecX(LS_POS_PCT, 50000)
	}

	i.oNotifyPostLayout(false)
	i.SetChanged(true)
}

// finds a cursor position by a physical position
func (i *ItemTextInput) findCursorLoc(x int32) int {
	base := i.cursorPosition[i.textWindowStart]
	for n := i.textWindowStart; n < (len(i.cursorPosition) - 1); n++ {
		if x < (((i.cursorPosition[n] + i.cursorPosition[n+1]) / 2) - base) {
			//fmt.Printf("findCursorLoc(%d)=%d\n", x, n)
			return n
		}
	}
	//fmt.Printf("findCursorLoc(%d)=%d\n", x, len(i.Text))
	return len(i.text)
}

func (i *ItemTextInput) insertRune(pos int, r rune) {
	// insert into runes and update cursor positions
	i.text = append(i.text, r)
	copy(i.text[pos+1:], i.text[pos:])
	i.text[pos] = r

	// extend the position array by one
	i.cursorPosition = append(i.cursorPosition, 0)
	// shift tail by one position to the right  copy(dst,src)
	copy(i.cursorPosition[pos+2:], i.cursorPosition[pos+1:])
	// insert the new position right of the inserted rune
	i.cursorPosition[pos+1] = i.style.GetTextLen(string(i.text[:pos+1]))

	// if it was an insert, not an append
	if pos < len(i.cursorPosition)-2 {
		//compute positional change of all right shifted
		dif := i.style.GetTextLen(string(i.text[:pos+2])) - i.cursorPosition[pos+2]
		// correct all right shifted by the difference that was caused by the insert
		l := len(i.cursorPosition)
		for n := pos + 2; n < l; n++ {
			i.cursorPosition[n] += dif
		}
	}
}

func (i *ItemTextInput) removeRune(pos int) {

	// delete from runes and update cursor positions
	copy(i.text[pos:], i.text[pos+1:])
	i.text = i.text[:len(i.text)-1]

	// shift tail by one position to the left - copy(dst,src)
	copy(i.cursorPosition[pos+1:], i.cursorPosition[pos+2:])
	// reduce the position array by one
	i.cursorPosition = i.cursorPosition[:len(i.cursorPosition)-1]

	// if it wasnt the last rune removed
	if pos < len(i.cursorPosition)-1 {
		oldpos := i.cursorPosition[pos+1]
		// update the new position right of the removed one
		i.cursorPosition[pos+1] = i.style.GetTextLen(string(i.text[:pos+1]))
		//compute positional change of all shifted
		dif := i.cursorPosition[pos+1] - oldpos
		// correct all left shifted by the difference that was caused by the insert
		l := len(i.cursorPosition)
		for n := pos + 2; n < l; n++ {
			i.cursorPosition[n] += dif
		}
	}

}

// computes a visible text window based on the current text and cursor position.
// text and cursor positions need to be consistent before calling this function
func (i *ItemTextInput) computeTextWindow() {

	// show as much text as possible
	// cursor moves inside text window do not move text window

	// available space
	spc := i.oGetSubFrame().W
	plen := len(i.cursorPosition)

	getend := func(begin int) (int, bool) {
		var max bool = false
		if begin < plen {
			base := i.cursorPosition[begin]
			var n int
			for n = begin + 1; n < plen; n++ {
				if i.cursorPosition[n]-base > spc {
					max = true
					break
				}
			}
			fmt.Printf("getend: begin:%d plen:%d loc:%d\n", begin, plen, n-1)
			return n - 1, max

		} else {
			fmt.Printf("getend: begin:%d plen:%d loc:%d\n", begin, plen, plen-1)
			return plen - 1, max
		}
	}

	getstart := func(begin int) (int, bool) {
		var max bool = false
		if begin > 0 {
			base := i.cursorPosition[begin]
			var n int
			for n = begin - 1; n >= 0; n-- {
				if base-i.cursorPosition[n] > spc {
					max = true
					break
				}
			}
			fmt.Printf("getstart: begin:%d plen:%d loc:%d\n", begin, plen, n+1)
			return n + 1, max
		} else {
			fmt.Printf("getstart: begin:%d plen:%d loc:%d\n", begin, plen, 0)
			return 0, max
		}
	}

	// text may have changed so current start and end positions may not be correct
	// first check if cursor moved out of current range
	if i.cursorLocation < i.textWindowStart {
		//compute window end
		fmt.Printf("i.cursorLoc < i.textWindowStart\n")

		i.textWindowStart = i.cursorLocation
		i.textWindowEnd, _ = getend(i.cursorLocation)

	} else if i.textWindowEnd < i.cursorLocation {
		//compute window start
		fmt.Printf("i.textWindowEnd < i.cursorLoc\n")
		i.textWindowStart, _ = getstart(i.cursorLocation)
		i.textWindowEnd = i.cursorLocation
	} else {
		fmt.Printf("cursor in window\n")
		// cursor is in current window

		// check bounds of window are stil valid
		if i.textWindowEnd > plen-1 {
			// end is beyond text-end
			i.textWindowStart, _ = getstart(plen - 1)
			i.textWindowEnd = plen - 1
		} else {
			i.textWindowEnd, _ = getend(i.textWindowStart)
			// if cursor location is outside end - compute window start
			if i.textWindowEnd < i.cursorLocation {
				fmt.Printf("i.textWindowEnd < i.cursorLoc\n")
				//compute window start
				i.textWindowStart, _ = getstart(i.cursorLocation)
				i.textWindowEnd = i.cursorLocation
			}
		}
	}

	s := string(i.text[i.textWindowStart:i.textWindowEnd])
	i.textBox.SetText(s)
	if i.textAlign != ITEM_TEXTINPUT_ALIGN_LEFT {
		// layout textbox
		i.textBox.Layout(i.oGetSubFrame(), true)
	}
	fmt.Printf("computeTextWindow start:%d end:%d loc:%d text:%s\n", i.textWindowStart, i.textWindowEnd, i.cursorLocation, s)

}

func (i *ItemTextInput) centerTextWindow() {

	plen := len(i.cursorPosition)
	startPosition := plen / 2
	i.cursorLocation = startPosition
	base := i.cursorPosition[startPosition]
	spc := i.oGetSubFrame().W
	var sumr int32
	var suml int32

	i.textWindowStart = startPosition
	i.textWindowEnd = startPosition

	for n := 0; n < plen; n++ {
		if startPosition+n < plen {
			sumr = i.cursorPosition[startPosition+n] - base
			//fmt.Printf("r: start:%d end:%d sumr:%d\n", i.textWindowStart, i.textWindowEnd, sumr)
			if sumr+suml > spc {
				break
			} else {
				i.textWindowEnd = startPosition + n
			}
		}

		if startPosition-n >= 0 {
			suml = base - i.cursorPosition[startPosition-n]
			//fmt.Printf("l: start:%d end%d suml:%d\n", i.textWindowStart, i.textWindowEnd, suml)
			if sumr+suml > spc {
				break
			} else {
				i.textWindowStart = startPosition - n
			}
		}
	}
	s := string(i.text[i.textWindowStart:i.textWindowEnd])
	i.textBox.SetText(s)
	if i.textAlign != ITEM_TEXTINPUT_ALIGN_LEFT {
		// layout textbox
		i.textBox.Layout(i.oGetSubFrame(), true)
	}
}

// For every rune in text it computes its pixel end-position in the texture. So the cursor knows where to print itself in the texture.
func (i *ItemTextInput) computeCursorPositions() {
	fmt.Printf("computeCursorPositions:\n")
	// adapt to font kerning
	l := len(i.text)

	if l >= (cap(i.cursorPosition)) {
		i.cursorPosition = make([]int32, l+1, l+8)
	}

	//record positions between each rune. 0=pos1 l+1=end
	for n := range i.text {
		i.cursorPosition[n] = i.style.GetTextLen(string(i.text[:n]))
	}
	i.cursorPosition[l] = i.style.GetTextLen(string(i.text[:]))

	// make cursor position array same size as Text +1
	i.cursorPosition = i.cursorPosition[:l+1]
}

func (i *ItemTextInput) SetText(s string) {

	i.text = []rune(s)
	i.computeCursorPositions()
	i.cursorLocation = 0
	i.computeTextWindow()

}

func (i *ItemTextInput) GetText() string {
	return string(i.text)
}
func (i *ItemTextInput) GetNumber() float64 {
	n, err := strconv.ParseFloat(string(i.text), 64)
	if err != nil {
		return 0
	}
	return n
}

func (i *ItemTextInput) oRender() {
	//fmt.Printf("%s: ItemTextInput.oRender()\n", i.GetName())
	r := i.GetRenderer()
	//s := i.GetStyle()
	c := i.GetColorScheme()
	r.SetDrawBlendMode(sdl.BLENDMODE_NONE)

	// render background

	if i.active {
		utilRendererSetDrawColor(r, &c.baseBright)
	} else {
		utilRendererSetDrawColor(r, &c.base)
	}
	r.FillRect(&i.iframe)
	utilRenderShadowBorder(r, &i.iframe, c, true)

	// render text
	i.textBox.oRender()

	// render active box
	if i.active {
		if i.cursorVisible {
			// render cursor
			i.cursorRect.X = i.textBox.iframe.X + (i.cursorPosition[i.cursorLocation] - i.cursorPosition[i.textWindowStart])
			i.cursorRect.Y = i.textBox.iframe.Y
			//r.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
			utilRendererSetDrawColor(r, &c.text)
			r.DrawRect(&i.cursorRect)
		}
	}

}

// Layout has happened for the Item - so every item gets asked to layout its subitems
func (i *ItemTextInput) oNotifyPostLayout(sizeChanged bool) {
	//fmt.Printf("%s: ItemTextInput.oNotifyPostLayout()\n", i.GetName())
	// layout decoration
	i.textBox.Layout(i.oGetSubFrame(), sizeChanged)

	if sizeChanged {
		// after a layout has happened the text window may be resized
		// so let the text window adjust
		i.computeCursorPositions()
		i.cursorLocation = 0
		i.computeTextWindow()
	}

}
func (i *ItemTextInput) oWithItems(fn func(Item)) {
	fmt.Printf("%s: ItemBase.oWithItems() no subitems\n", i.GetName())
	// handle the call for myself
	fn(i)
	// forward to all subitems
	i.textBox.oWithItems(fn)
}
func (i *ItemTextInput) oNotifyTimer() {
	// toggle cursor
	//fmt.Printf("ItemTextInput: %s.oNotifyTimer()\n", i.name)
	if i.cursorVisible {
		i.cursorVisible = false
	} else {
		i.cursorVisible = true
	}
	if i.active {
		i.SetTimer(i.style.cursorBlinkRate)
	}
	i.SetChanged(true)
}

func (i *ItemTextInput) oNotifyMouseMotion(x, y, dx, dy int32) {
	fmt.Printf("ItemTextInput: %s.NotifyMouseMotion(%d,%d,%d,%d)\n", i.name, x, y, dx, dy)
}

func (i *ItemTextInput) oNotifyMouseButton(x, y int32, button sdl.Button, state sdl.ButtonState) {
	fmt.Printf("ItemTextInput: %s.NotifyMouseButton(%d,%d,%d,%d)\n", i.name, x, y, button, state)

	if button == 1 {
		switch state {
		case 1:
			// button 1 pressed
			i.pick(true)
			i.pickedx = x - i.iframe.X
			i.pickedy = y - i.iframe.Y

			if i.active {
				i.cursorLocation = i.findCursorLoc(x - i.textBox.iframe.X)
				i.computeTextWindow()

			} else {
				i.win.SetKbFocusItem(i)
				i.cursorLocation = i.findCursorLoc(x - i.textBox.iframe.X)
				i.computeTextWindow()
			}
		case 0:
			// button 1 released
			// button 1 released
			i.pick(false)
		}
	}
}

func (i *ItemTextInput) cursorOn() {
	if i.active {
		i.tmpCursor = sdl.GetCursor()
		if i.cursor != nil {
			sdl.SetCursor(i.cursor)
		}
	}

}
func (i *ItemTextInput) cursorOff() {
	if i.active {
		if i.tmpCursor != nil {
			sdl.SetCursor(i.tmpCursor)
		}
		i.tmpCursor = nil
	}
}

func (i *ItemTextInput) oNotifyMouseFocusLost() {
	i.pick(false)
	i.cursorOff()
}

func (i *ItemTextInput) oNotifyMouseFocusGained() {
	i.pick(false)
	i.cursorOn()
}

func (i *ItemTextInput) oNotifyKbFocusGained() {
	fmt.Printf("ItemTextInput: %s.NotifyKbFocusGained()\n", i.name)

	i.active = true
	i.cursorOn()
	sdl.StartTextInput()
	// start cursor timer
	i.cursorVisible = true
	i.SetTimer(i.style.cursorBlinkRate)
	i.SetChanged(true)
}
func (i *ItemTextInput) oNotifyKbFocusLost() {
	fmt.Printf("ItemTextInput: %s.NotifyKbFocusLost()\n", i.name)
	i.cursorOff()
	i.active = false
	sdl.StopTextInput()

	switch i.textAlign {
	case ITEM_TEXTINPUT_ALIGN_LEFT:
		i.cursorLocation = 0
		i.computeTextWindow()
	case ITEM_TEXTINPUT_ALIGN_RIGHT:
		i.cursorLocation = len(i.cursorPosition) - 1
		i.computeTextWindow()
	case ITEM_TEXTINPUT_ALIGN_CENTER:
		i.centerTextWindow()
	}

	i.SetChanged(true)
}

func (i *ItemTextInput) oNotifyTextInput(r rune) {
	fmt.Printf("ItemTextInput: %s.NotifyTextInput() %c\n", i.name, r)
	i.insertRune(i.cursorLocation, r)
	i.cursorLocation++
	i.computeTextWindow()
	i.SetChanged(true)
}

func (i *ItemTextInput) oNotifyKbEvent(e *sdl.KeyboardEvent) {
	if e.State == 1 {
		// only key down events
		fmt.Printf("ItemTextInput: %s.NotifyKeyboardEvent() mod:%d scancode:%d sym:%d\n", i.name, e.Keysym.Mod, e.Keysym.Scancode, e.Keysym.Sym)
		switch e.Keysym.Scancode {
		case sdl.SCANCODE_BACKSPACE:
			if i.cursorLocation > 0 {
				i.cursorLocation--
				i.removeRune(i.cursorLocation)
				i.computeTextWindow()
			}

		case sdl.SCANCODE_DELETE:
			if i.cursorLocation < len(i.text) {
				i.removeRune(i.cursorLocation)
				i.computeTextWindow()
			}

		case sdl.SCANCODE_RIGHT:
			if i.cursorLocation < len(i.text) {
				i.cursorLocation++
				i.computeTextWindow()
			}
		case sdl.SCANCODE_LEFT:
			if i.cursorLocation > 0 {
				i.cursorLocation--
				i.computeTextWindow()
			}
		case sdl.SCANCODE_END:
			i.cursorLocation = len(i.text)
			i.computeTextWindow()

		case sdl.SCANCODE_HOME:
			i.cursorLocation = 0
			i.computeTextWindow()
		}
		i.SetChanged(true)
	}

}

//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

type WindowController struct {
	windows    []*RootWindow
	ctrl       chan bool
	queue      chan interface{}
	waitgroup  *sync.WaitGroup
	lastId     uint32
	lastWindow *RootWindow
}

func NewWindowController() *WindowController {
	w := new(WindowController)

	if sdl.WasInit(sdl.INIT_VIDEO) != sdl.INIT_VIDEO {
		sdl.Init(sdl.INIT_VIDEO)
	}

	return w
}

func (wc *WindowController) Destroy() {

	// check if control channel ist closed by reading from it. if not closed close it
	select {
	case _, running := <-wc.ctrl:
		if running {
			close(wc.ctrl)
		}
	}

	if ttf.WasInit() {
		ttf.Quit()
	}

	if sdl.WasInit(sdl.INIT_VIDEO) == sdl.INIT_VIDEO {
		sdl.Quit()
	}
}

func (wc *WindowController) AddRootWindow(rw *RootWindow) {
	wc.windows = append(wc.windows, rw)
}

func (wc *WindowController) getRootWindowById(id uint32) (bool, *RootWindow) {
	if id == wc.lastId {
		return true, wc.lastWindow
	} else {
		for _, w := range wc.windows {
			if i, _ := w.Window.GetID(); i == id {
				wc.lastId = i
				wc.lastWindow = w
				return true, w
			}
		}
	}
	return false, nil
}

func (wc *WindowController) eventSenderRenderTick(msec int) {
	wc.waitgroup.Add(1)
	e := EventRenderTick{msec: msec}
	var running bool = true
runloop:
	for {
		select {
		case _, running = <-wc.ctrl:
			if !running {
				fmt.Println("quitting eventSenderTick")
				wc.waitgroup.Done()
				break runloop
			}
		default:
			wc.queue <- &e
			time.Sleep(time.Duration(msec) * time.Millisecond)
		}
	}
}

func (wc *WindowController) eventSenderSDL() {
	wc.waitgroup.Add(1)
	var running bool = true

runloop:
	for {
		// A select blocks until one of its cases can run, then it executes that case. It chooses one at random if multiple are ready.
		select {
		case _, running = <-wc.ctrl:
			if !running {
				fmt.Println("quitting eventSenderSDL")
				wc.waitgroup.Done()
				break runloop
			}
		case wc.queue <- sdl.WaitEventTimeout(10): // wait here until an sdl-event arives. pick one.
		}
	}
}

func (wc *WindowController) Start() {

	// layout the windows
	for _, w := range wc.windows {
		w.Layout()
	}

	// event sender/receiver control channel and common event queue
	wc.ctrl = make(chan bool, 2)
	wc.queue = make(chan interface{}, 2)
	wc.waitgroup = new(sync.WaitGroup)

	// starting the event senders
	go wc.eventSenderSDL()
	go wc.eventSenderRenderTick(5) // 5 msec = 20frames/sec

	//go wc.eventSenderCustom(&wg, wc.queue, wc.ctrl)

	//starting the main event processor
	var event interface{}
queuereader:
	for event = range wc.queue {
		switch t := event.(type) {
		case *EventRenderTick:
			// render windows that have changed
			for _, w := range wc.windows {
				if w.GetChanged() {
					sw.Start("render")
					w.RenderAll()
					w.Present()
					sw.Stop("render")
				}
			}
		case *EventCustom:
			fmt.Println("Custom event!")
		case *sdl.QuitEvent:
			close(wc.ctrl)
			break queuereader
		case *sdl.WindowEvent:
			if t.Event == sdl.WINDOWEVENT_CLOSE {
				close(wc.ctrl)
				break queuereader
			} else {
				if found, rw := wc.getRootWindowById(t.WindowID); found {
					rw.handleWindowEvent(t)
				}
			}
		case *sdl.MouseMotionEvent:
			//fmt.Printf("[%d ms] MouseMotion\ttype:%d\twhich:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\twindow:%d\tstate:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel, t.WindowID, t.State)
			if found, rw := wc.getRootWindowById(t.WindowID); found {
				rw.handleMouseMotionEvent(t)
			}
		case *sdl.MouseButtonEvent:
			//fmt.Printf("[%d ms] MouseButton\ttype:%d\twhich:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\twindow:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State, t.WindowID)
			if found, rw := wc.getRootWindowById(t.WindowID); found {
				rw.handleMouseButtonEvent(t)
			} else if t.WindowID == 0 {
				// button was released outside a window
				// notify last known window
				if wc.lastWindow != nil {
					wc.lastWindow.handleMouseButtonEvent(t)
				}
			}
		case *sdl.MouseWheelEvent:
			fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y)
		case *sdl.KeyboardEvent:
			if found, rw := wc.getRootWindowById(t.WindowID); found {
				rw.handleKeyboardEvent(t)
			}
			fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n", t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
		case *sdl.TextInputEvent:
			//t.Text : the null-terminated input text in UTF-8 encoding
			//fmt.Printf("[%d ms] TextInput\ttype:%d\ttext:%s\n", t.Timestamp, t.Type, string(t.Text[:]))
			if found, rw := wc.getRootWindowById(t.WindowID); found {
				rw.handleTextInputEvent(t)
			}
		case *sdl.SysWMEvent:
			fmt.Printf("[%d ms] SysWMEvent\ttype:%d\n", t.Timestamp, t.Type)
		case *sdl.TouchFingerEvent:
		case *sdl.UserEvent:
			fmt.Printf("[%d ms] UserEvent\ttype:%d\n", t.Timestamp, t.Type)
		}
	}

	fmt.Println("waiting for wait group to finish")
	wc.waitgroup.Wait()
	fmt.Println("wait group is done")

	// Clean up the sdl resources befor exit
	for _, w := range wc.windows {
		w.Destroy()
	}

	sdl.Quit()
}

// --------------------------------------------------------------------
type layoutParam struct {
	S int32
	V int32
}
type layoutSpec struct {
	X layoutParam
	Y layoutParam
	W layoutParam
	H layoutParam
}

// --------------------------------------------------------------------
type ItemWindowManager struct {
	ItemLayerBox
}

func NewItemWindowManager(win *RootWindow) *ItemWindowManager {
	i := new(ItemWindowManager)
	i.o = Item(i)
	i.setRootWindow(win)
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_PCT, 0, 0, 100000, 100000)
	return i
}

func (i *ItemWindowManager) oNotifyPostLayout(sizeChanged bool) {

	for _, si := range i.subItems.GetList() {
		//fmt.Printf("--> subitem:%s\n", si.GetName())
		sf := i.oGetSubFrame()
		i.subItems.GetTop().SetState(ITEM_STATE_BACKGROUND)
		si.Layout(sf, sizeChanged)
	}
	i.subItems.GetTop().SetState(ITEM_STATE_ACTIVE)
}
func (i *ItemWindowManager) AddSubItem(w *ItemWindow) {
	// can only have ItemWindow as Subitems
	i.subItems.AddItem(Item(w))
	w.SetState(ITEM_STATE_BACKGROUND)
	w.SetParent(i)
	w.manager = i
	i.oNotifyPostLayout(true)
}

func (i *ItemWindowManager) oFindSubItem(x, y int32, e sdl.Event) (found bool, item Item) {
	// browse thru all subitems top down. first hit wins
	if si, found := i.subItems.CheckTopDown(x, y); found {
		if si.CheckPos(x, y) {
			// if a Mouse button event was the reason
			if e != nil {
				switch e.(type) {
				case *sdl.MouseButtonEvent:
					ti := i.subItems.GetTop()
					if ti != nil && ti != si {
						ti.SetState(ITEM_STATE_BACKGROUND)
						i.subItems.ShiftTop(si)
						si.SetState(ITEM_STATE_ACTIVE)
						i.SetChanged(true)
					}
				}
			}
			return true, si
		}
	}

	return false, nil
}

// --------------------------------------------------------------------
// ItemBackground is a background whith a subitem
type ItemWindow struct {
	ItemBase
	handleGrid *ItemGrid
	rootGrid   *ItemGrid
	titleGrid  *ItemGrid
	titleText  *ItemText
	panelLayer *ItemLayerBox
	lastX      int32
	lastY      int32
	manager    *ItemWindowManager
	//wmCbButton
}

func NewItemWindow(win *RootWindow) *ItemWindow {
	i := new(ItemWindow)
	i.o = Item(i)
	i.setRootWindow(win)
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_PCT, 0, 0, 100000, 100000)
	i.colorScheme = win.style.csWindow
	i.SetName("NewWindow")

	i.handleGrid = NewItemGrid(win, 3, 3)
	i.handleGrid.SetParent(i)
	i.handleGrid.SetName("handleGrid")
	i.handleGrid.SetSpec(LS_POS_ABS, LS_POS_ABS, LS_SIZE_PCT, LS_SIZE_PCT, 0, 0, 100000, 100000)

	var hdim int32 = 10

	setHandleParam := func(h *ItemHandle, name string, p Item, cb func(int32, int32)) *ItemHandle {
		h.SetName(name)
		//h.SetItemStyle(ITEM_HANDLE_STYLE_DARKEDGE)
		h.SetItemStyle(ITEM_HANDLE_STYLE_DARKEDGE)
		h.SetCursor(sdl.SYSTEM_CURSOR_SIZENWSE)
		h.SetParent(p)
		h.SetTargetCb(cb)
		return h
	}

	sh := setHandleParam(NewItemHandle(win), "sizeHandleNW", i.handleGrid, i.sizeCbNW)
	sh.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_ABS, LS_SIZE_ABS, 0, 0, hdim, hdim)
	i.handleGrid.SetSubItem(0, 0, sh)

	sh = setHandleParam(NewItemHandle(win), "sizeHandleN", i.handleGrid, i.sizeCbN)
	sh.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_ABS, 0, 0, 100000, hdim)
	i.handleGrid.SetSubItem(1, 0, sh)

	sh = setHandleParam(NewItemHandle(win), "sizeHandleNE", i.handleGrid, i.sizeCbNE)
	sh.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_ABS, LS_SIZE_ABS, 0, 0, hdim, hdim)
	i.handleGrid.SetSubItem(2, 0, sh)

	sh = setHandleParam(NewItemHandle(win), "sizeHandleW", i.handleGrid, i.sizeCbW)
	sh.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_ABS, LS_SIZE_PCT, 0, 0, hdim, 100000)
	i.handleGrid.SetSubItem(0, 1, sh)

	sh = setHandleParam(NewItemHandle(win), "sizeHandleE", i.handleGrid, i.sizeCbE)
	sh.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_ABS, LS_SIZE_PCT, 0, 0, hdim, 100000)
	i.handleGrid.SetSubItem(2, 1, sh)

	sh = setHandleParam(NewItemHandle(win), "sizeHandleSW", i.handleGrid, i.sizeCbSW)
	sh.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_ABS, LS_SIZE_ABS, 0, 0, hdim, hdim)
	i.handleGrid.SetSubItem(0, 2, sh)

	sh = setHandleParam(NewItemHandle(win), "sizeHandleS", i.handleGrid, i.sizeCbS)
	sh.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_ABS, 0, 0, 100000, hdim)
	i.handleGrid.SetSubItem(1, 2, sh)

	sh = setHandleParam(NewItemHandle(win), "sizeHandleSE", i.handleGrid, i.sizeCbSE)
	sh.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_ABS, LS_SIZE_ABS, 0, 0, hdim, hdim)
	i.handleGrid.SetSubItem(2, 2, sh)

	i.handleGrid.SetColSpec(0, layoutParam{LS_SIZE_COLLAPSE, 0})
	i.handleGrid.SetColSpec(1, layoutParam{LS_SIZE_PCT, 100000})
	i.handleGrid.SetColSpec(2, layoutParam{LS_SIZE_COLLAPSE, 0})
	i.handleGrid.SetRowSpec(0, layoutParam{LS_SIZE_COLLAPSE, 0})
	i.handleGrid.SetRowSpec(1, layoutParam{LS_SIZE_PCT, 100000})
	i.handleGrid.SetRowSpec(2, layoutParam{LS_SIZE_COLLAPSE, 0})

	// rootGrid separates Title and panel of window
	i.rootGrid = NewItemGrid(win, 1, 2)
	i.rootGrid.SetParent(i)
	i.rootGrid.SetName("windowGrid")
	i.rootGrid.SetSpacing(i.GetStyle().spacing)

	// titleGrid separates a laxerbox with text and handle from further buttons at the rigth
	i.titleGrid = NewItemGrid(win, 4, 1)
	i.rootGrid.SetSubItem(0, 0, i.titleGrid)
	i.titleGrid.SetColSpec(0, layoutParam{LS_SIZE_PCT, 100000})
	i.titleGrid.SetColSpec(1, layoutParam{LS_SIZE_COLLAPSE, 0})
	i.titleGrid.SetColSpec(2, layoutParam{LS_SIZE_COLLAPSE, 0})
	i.titleGrid.SetColSpec(3, layoutParam{LS_SIZE_COLLAPSE, 0})
	i.titleGrid.SetSpacing(2)
	i.titleGrid.SetName("titleGrid")

	lb := NewItemLayerBox(win)
	i.titleGrid.SetSubItem(0, 0, lb)
	lb.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_PCT, 0, 0, 100000, 100000)
	lb.SetName("titleLayer")

	btn := NewItemButton(win)
	i.titleGrid.SetSubItem(1, 0, btn)
	btn.SetName("titleButtonMin")
	btn.SetText("m")
	btn.SetCallback(i.minCb)
	btn.SetSpecSize(LS_POS_ABS, LS_POS_ABS, i.GetStyle().averageFontHight, i.GetStyle().averageFontHight)

	btn = NewItemButton(win)
	i.titleGrid.SetSubItem(2, 0, btn)
	btn.SetName("titleButtonMax")
	btn.SetText("M")
	btn.SetCallback(i.maxCb)
	btn.SetSpecSize(LS_POS_ABS, LS_POS_ABS, i.GetStyle().averageFontHight, i.GetStyle().averageFontHight)

	btn = NewItemButton(win)
	i.titleGrid.SetSubItem(3, 0, btn)
	btn.SetName("titleButtonX")
	btn.SetText("x")
	btn.SetCallback(i.closeCb)
	btn.SetSpecSize(LS_POS_ABS, LS_POS_ABS, i.GetStyle().averageFontHight, i.GetStyle().averageFontHight)

	hdl := NewItemHandle(win)
	lb.AddSubItem(hdl)
	hdl.SetTargetCb(i.moveCb)
	hdl.SetName("titleHandle")
	hdl.SetCursorPicked(sdl.SYSTEM_CURSOR_SIZEALL)

	i.titleText = NewItemText(win)
	lb.AddSubItem(i.titleText)
	i.titleText.SetText("Title")
	i.titleText.SetName("titleText")

	bg := NewItemBackground(win)
	i.rootGrid.SetSubItem(0, 1, bg)
	//bg.SetMargin(5)

	// the panel of the window
	i.panelLayer = NewItemLayerBox(win)
	bg.SetSubItem(i.panelLayer)
	i.panelLayer.SetName("panelLayer")

	return i
}

func (i *ItemWindow) moveCb(dx, dy int32) {
	i.Move(dx, dy)
	_, _, _, _, i.lastX, i.lastY, _, _ = i.GetSpec()
	i.SetChanged(true)
}

func (i *ItemWindow) sizeCbCollector(dx, dy int32, bw, bh bool) {
	i.Size(dx, dy, bw, bh)
	_, _, _, _, i.lastX, i.lastY, _, _ = i.GetSpec()
	i.SetChanged(true)
}

func (i *ItemWindow) sizeCbNW(dx, dy int32) { i.sizeCbCollector(dx, dy, true, true) }
func (i *ItemWindow) sizeCbN(dx, dy int32)  { i.sizeCbCollector(0, dy, false, true) }
func (i *ItemWindow) sizeCbNE(dx, dy int32) { i.sizeCbCollector(dx, dy, false, true) }

func (i *ItemWindow) sizeCbW(dx, dy int32) { i.sizeCbCollector(dx, 0, true, false) }
func (i *ItemWindow) sizeCbE(dx, dy int32) { i.sizeCbCollector(dx, 0, false, false) }

func (i *ItemWindow) sizeCbSW(dx, dy int32) { i.sizeCbCollector(dx, dy, true, false) }
func (i *ItemWindow) sizeCbS(dx, dy int32)  { i.sizeCbCollector(0, dy, false, false) }
func (i *ItemWindow) sizeCbSE(dx, dy int32) { i.sizeCbCollector(dx, dy, false, false) }

func (i *ItemWindow) closeCb() {
	if i.manager != nil {
		i.manager.subItems.RemoveItem(i)
		i.manager.SetChanged(true)
	}
}
func (i *ItemWindow) minCb() {

	i.SetSpec(LS_POS_ABS, LS_POS_ABS, LS_SIZE_COLLAPSE, LS_SIZE_COLLAPSE, i.lastX, i.lastY, 0, 0)
	i.Layout(&i.pframe, true)
	i.SetChanged(true)
}

func (i *ItemWindow) maxCb() {
	_, _, _, _, i.lastX, i.lastY, _, _ = i.GetSpec()
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_PCT, 0, 0, 100000, 100000)
	i.Layout(&i.pframe, true)
	i.SetChanged(true)
}

func (i *ItemWindow) AddSubItem(s Item) {
	i.panelLayer.AddSubItem(s)
}

func (i *ItemWindow) oRender() {
	r := i.GetRenderer()
	c := i.GetColorScheme()

	// background of window
	utilRenderFillRect(r, &i.iframe, &c.base)

	// outer border
	utilRenderShadowBorder(r, &i.iframe, c, false)

	// size handles
	//i.handleGrid.Render()

	// background of rootGrid
	utilRenderFillRect(r, &i.rootGrid.pframe, &c.base)

	//inner Border around panelLayer
	utilRenderShadowBorder(r, i.panelLayer.MakeInnerFrame(-1), c, true)

	//background of title grid
	//utilRenderFillRect(r, &i.titleGrid.pframe, &c.base)

	// panelLayer background
	//utilRenderFillRect(r, &i.panelLayer.iframe, &c.baseDark)

	// content
	i.rootGrid.Render()

}
func (i *ItemWindow) oReportSubitems(lvl int) {
	i.handleGrid.Report(lvl)
	i.rootGrid.Report(lvl)
}

func (i *ItemWindow) oNotifyPostLayout(sizeChanged bool) {
	i.handleGrid.Layout(&i.iframe, sizeChanged)
	i.rootGrid.Layout(i.oGetSubFrame(), sizeChanged)
}

func (i *ItemWindow) oGetSubFrame() *sdl.Rect {
	return i.MakeInnerFrame(i.GetStyle().spacing + 1)
}

func (i *ItemWindow) oGetMinSize() (int32, int32) {
	w, h := i.rootGrid.oGetMinSize()
	spc := (i.GetStyle().spacing + 1) * 2
	return w + spc, h + spc
}
func (i *ItemWindow) oWithItems(fn func(Item)) {
	fmt.Printf("%s: ItemBase.oWithItems() no subitems\n", i.GetName())
	// handle the call for myself
	fn(i)
	// forward to all subitems
	i.handleGrid.oWithItems(fn)
	i.rootGrid.oWithItems(fn)
}
func (i *ItemWindow) oFindSubItem(x, y int32, e sdl.Event) (bool, Item) {

	if i.rootGrid.CheckPos(x, y) {
		return true, i.rootGrid
	}

	if i.handleGrid.CheckPos(x, y) {
		return true, i.handleGrid
	}

	return false, nil

}

// --------------------------------------------------------------------
// ItemLayerBox is designed manage layered items
// all items share the same frame
type ItemLayerBox struct {
	ItemBase
	subItems       infraItemList
	buttonCallback func(int32, int32, sdl.Button, sdl.ButtonState) //x, y int32, button sdl.Button, state sdl.ButtonState
	//spacing        int32
}

func NewItemLayerBox(win *RootWindow) *ItemLayerBox {
	i := new(ItemLayerBox)
	i.o = Item(i)
	i.setRootWindow(win)
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_PCT, 0, 0, 100000, 100000)
	return i
}

func (i *ItemLayerBox) SetButtonCallback(cb func(int32, int32, sdl.Button, sdl.ButtonState)) {
	i.buttonCallback = cb
}
func (i *ItemLayerBox) oNotifyMouseButton(x, y int32, button sdl.Button, state sdl.ButtonState) {
	if i.buttonCallback != nil {
		i.buttonCallback(x, y, button, state)
	}
}

func (i *ItemLayerBox) oRender() {

	/* // debug
	r := i.GetRenderer()
	utilRenderSolidBorder( Border(r, &i.iframe, i.GetStyle().colorPurple)
	*/

	for _, si := range i.subItems.GetList() {
		si.Render()
	}
}

func (i *ItemLayerBox) AddSubItem(s Item) {
	i.subItems.AddItem(s)
	s.SetParent(i)
}

func (i *ItemLayerBox) RemoveSubItem(s Item) bool {
	return i.subItems.RemoveItem(i)
}
func (i *ItemLayerBox) oReportSubitems(lvl int) {
	for _, si := range i.subItems.GetList() {
		si.Report(lvl)
	}
}

func (i *ItemLayerBox) oNotifyPostLayout(sizeChanged bool) {
	//fmt.Printf("%s: ItemLayer.oNotifyPostLayout()\n", i.GetName())
	for _, si := range i.subItems.GetList() {
		//fmt.Printf("--> subitem:%s\n", si.GetName())
		sf := i.oGetSubFrame()
		si.Layout(sf, sizeChanged)
	}
}

func (i *ItemLayerBox) oFindSubItem(x, y int32, e sdl.Event) (found bool, item Item) {
	// browse thru all subitems top down. first hit wins

	if si, found := i.subItems.CheckTopDown(x, y); found {
		if si.CheckPos(x, y) {
			return true, si
		}
	}

	return false, nil

}
func (i *ItemLayerBox) oGetSubFrame() *sdl.Rect {
	return &i.iframe
}

func (i *ItemLayerBox) oWithItems(fn func(Item)) {
	fmt.Printf("%s: ItemBase.oWithItems() no subitems\n", i.GetName())
	// handle the call for myself
	fn(i)
	// forward to all subitems
	for _, si := range i.subItems.GetList() {
		si.oWithItems(fn)
	}
}
func (i *ItemLayerBox) oGetMinSize() (int32, int32) {
	if !i.useTmpMinSize {
		if len(i.subItems.GetList()) == 0 {
			i.minw = 0
			i.minh = 0
		} else {
			var maxx, maxy, minx, miny int32
			minx = 2147483647
			miny = 2147483647

			for _, si := range i.subItems.GetList() {
				x, y, w, h := si.GetCollapsedSpec()
				w, h = si.oGetMinSize()
				//fmt.Printf("%s: ItemLayerBox.oGetMinSize() w:%d h:%d\n", si.GetName(), w, h)

				if x < minx {
					minx = x
				}
				if y < miny {
					miny = y
				}
				if x+w > maxx {
					maxx = x + w
				}
				if y+h > maxy {
					maxy = y + h
				}
			}
			if minx < 0 {
				minx = 0
			}
			if miny < 0 {
				miny = 0
			}

			i.minw = maxx - minx
			i.minh = maxy - miny
			fmt.Printf("%s: ItemLayerBox.oGetMinSize() w:%d h:%d\n", i.GetName(), i.minw, i.minh)
		}
		i.useTmpMinSize = true
	}

	return i.minw, i.minh
}

// --------------------------------------------------------------------
// menue entry can be a submenu or action
type ItemMenuEntry struct {
	ItemButton
	subMenue    *ItemMenu // if set open a submenue
	parentMenue *ItemMenu // the entries parent-menu is the menue it is residing in
	menuAction  func()
}

func NewItemMenuEntry(win *RootWindow) *ItemMenuEntry {
	i := new(ItemMenuEntry)
	i.o = Item(i)
	i.setRootWindow(win)
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_PCT, 0, 0, 100000, 100000)
	i.buttonStyle = ITEM_BUTTON_STYLE_FLAT
	i.SetCallback(i.menuTrigger)
	return i
}

func (i *ItemMenuEntry) openSubMenu() {
	if i.subMenue != nil {
		fmt.Printf("openSubMenu\n")
		// tell the Parent Menue to open a sub-menue
		i.parentMenue.OpenSubMenu(i.subMenue, i.GetFrame())
		i.subMenue.SetParentMenu(i.parentMenue)
	}
}
func (i *ItemMenuEntry) SetMenuAction(f func()) {
	i.menuAction = f
}

func (i *ItemMenuEntry) menuTrigger() {
	fmt.Printf("%s: ItemMenuEntry.menuTrigger()\n", i.GetName())

	if i.subMenue != nil {
		fmt.Printf("%s: ItemMenuEntry: Submenue needs to open\n", i.GetName())
		i.parentMenue.OpenSubMenu(i.subMenue, i.GetFrame())
	} else {
		fmt.Printf("%s: ItemMenuEntry: menuAction!\n", i.GetName())
		// notify my Menu that a menue action is pending so it can close beforehand
		i.parentMenue.NotifyMenuAction()
		if i.menuAction != nil {
			i.menuAction()
		}
	}
}
func (i *ItemMenuEntry) SetSubMenu(s *ItemMenu) {
	i.subMenue = s
	if i.subMenue != nil {
		i.SetActOnClick(true)
	}
}

func (i *ItemMenuEntry) SetParentMenu(s *ItemMenu) {
	i.parentMenue = s
}

// --------------------------------------------------------------------
// Container holding MenueEntries vertical or horizontal composed of Grid
type ItemMenu struct {
	ItemGrid
	orientation    int           //MENU_ORIENTATION_VERTICAL/MENU_ORIENTATION_HORIZONTAL
	activeSubMenue *ItemMenu     // the currently activated subMenue
	parentMenue    *ItemMenu     // the currently activated subMenue
	layer          *ItemLayerBox // the layer item on which i get activated and shown
}

func NewItemMenue(win *RootWindow) *ItemMenu {
	i := new(ItemMenu)
	i.ItemGrid = *NewItemGrid(win, 0, 0)
	i.o = Item(i)
	i.orientation = MENU_ORIENTATION_VERTICAL
	return i
}

// when a menue action desition is made, this is called to notify the menue
func (i *ItemMenu) NotifyMenuAction() {
	fmt.Printf("%s: ItemMenu.NotifyMenuAction()\n", i.GetName())
	// hand it to the parent if there is any
	if i.parentMenue != nil {
		i.parentMenue.NotifyMenuAction()
	} else {
		// close submenu if there is any
		i.CloseSubMenu()
	}
}

func (i *ItemMenu) CloseSubMenu() {
	fmt.Printf("%s: ItemMenu.CloseSubMenu()\n", i.GetName())
	if i.activeSubMenue != nil {
		fmt.Printf("%s: ItemMenu: have active submenue ()\n", i.GetName())
		// forward the call first
		i.activeSubMenue.CloseSubMenu()
		i.layer.RemoveSubItem(i.activeSubMenue)
		i.activeSubMenue = nil
		if i.parentMenue != nil {
			i.layer.Layout(i.win.GetFrame(), true)
		}
	}
	// if i am the root Menue remove the layer as well
	if i.parentMenue == nil {
		fmt.Printf("%s: ItemMenu: removing Layer()\n", i.GetName())
		i.win.RemoveRootItem(Item(i.layer))
		i.layer = nil
	}
	i.win.SetChanged(true)
}

func (i *ItemMenu) layerButtonCb(x, y int32, button sdl.Button, state sdl.ButtonState) {
	fmt.Printf("%s: ItemMenu.layerButtonCb()\n", i.GetName())
	if state == 1 {
		// only on button down events
		i.CloseSubMenu()
	}
}

func (i *ItemMenu) OpenSubMenu(sm *ItemMenu, f *sdl.Rect) {
	fmt.Printf("%s: ItemMenu.OpenSubMenu()\n", i.GetName())
	// if  layer is nil this is the first menu to open a sub menue.
	// create one
	if i.layer == nil {
		fmt.Printf("%s: ItemMenu: creating layer ()\n", i.GetName())
		i.layer = NewItemLayerBox(i.win)
		i.layer.SetName(fmt.Sprintf("%s-layer", i.name))
		i.layer.SetButtonCallback(i.layerButtonCb)
		i.win.AddRootItem(Item(i.layer))
	}

	if i.activeSubMenue != nil {
		i.CloseSubMenu()
	}

	sm.layer = i.layer
	sm.SetParentMenu(i)
	i.activeSubMenue = sm

	if i.orientation == MENU_ORIENTATION_HORIZONTAL {
		sm.SetSpec(LS_POS_ABS, LS_POS_ABS, LS_SIZE_COLLAPSE, LS_SIZE_COLLAPSE, f.X, f.Y+f.H, 0, 0)
	} else if i.orientation == MENU_ORIENTATION_VERTICAL {
		sm.SetSpec(LS_POS_ABS, LS_POS_ABS, LS_SIZE_COLLAPSE, LS_SIZE_COLLAPSE, f.X+f.W, f.Y, 0, 0)
	}

	i.layer.AddSubItem(sm)
	i.layer.Layout(i.win.GetFrame(), true)
	i.win.SetChanged(true)

}
func (i *ItemMenu) SetParentMenu(pm *ItemMenu) {
	i.parentMenue = pm
}

func (i *ItemMenu) SetOrientation(o int) {
	i.orientation = o
}

func (i *ItemMenu) AddEntry(e *ItemMenuEntry) {
	if i.orientation == MENU_ORIENTATION_HORIZONTAL {
		n := i.AppendColumn()
		i.SetSubItem(n, 0, Item(e))
		i.SetColSpec(n, layoutParam{LS_SIZE_COLLAPSE, 0})
		e.SetParentMenu(i)
	} else if i.orientation == MENU_ORIENTATION_VERTICAL {
		n := i.AppendRow()
		i.SetSubItem(0, n, Item(e))
		i.SetRowSpec(n, layoutParam{LS_SIZE_COLLAPSE, 0})
		e.SetParentMenu(i)
	}
	e.SetParent(i)
}

func (i *ItemMenu) AddNewMenuEntry(text string, sub *ItemMenu, cb func()) {
	me := NewItemMenuEntry(i.win)
	me.SetText(text)
	me.SetName(text)
	if sub != nil {
		me.SetSubMenu(sub)
	} else {
		me.SetMenuAction(cb)
	}
	i.AddEntry(me)

}

//--------------------------------------------------------------------

// --------------------------------------------------------------------
type RootWindow struct {
	id             uint32
	title          string
	sizeX          int32
	sizeY          int32
	Window         *sdl.Window
	Renderer       *sdl.Renderer
	style          *Style
	RootItems      infraItemList
	mouseFocusItem Item
	mouseFocusLock bool
	kbFocusItem    Item
	winTex         *sdl.Texture
	changed        bool
	stateButton1   bool
	stateButton2   bool
	stateButton3   bool
	posMouseX      int32
	posMouseY      int32
	relMouseX      int32
	relMouseY      int32
}

func NewRootWindow(s *Style, title string, x int32, y int32) *RootWindow {
	w := new(RootWindow)
	w.sizeX = x
	w.sizeY = y
	w.title = title
	w.style = s

	var err error

	w.Window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, w.sizeX, w.sizeY, sdl.WINDOW_RESIZABLE|sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	w.Renderer, err = sdl.CreateRenderer(w.Window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_TARGETTEXTURE)
	if err != nil {
		panic(err)
	}

	w.Renderer.SetDrawColor(0, 0, 0, 255)
	w.Renderer.Clear()

	w.id, _ = w.Window.GetID()
	w.mouseFocusLock = false

	return w

}

func (w *RootWindow) SetMouseFocusLock(b bool) {
	w.mouseFocusLock = b
	if b {
		fmt.Printf("Mouse focus locked\n")
	} else {
		fmt.Printf("Mouse focus unlocked\n")
		// if it was unlocked outside the focus item notify the new item at the last known position for focus change
		if found, i := w.GetItemByPos(w.posMouseX, w.posMouseY, nil); found {
			w.SetMouseFocusItem(i)
		}

	}
}
func (w *RootWindow) GetChanged() bool  { return w.changed }
func (w *RootWindow) SetChanged(c bool) { w.changed = c }

func (w *RootWindow) GetStyle() *Style  { return w.style }
func (w *RootWindow) SetStyle(s *Style) { w.style = s }

func (w *RootWindow) AddRootItem(i Item) Item {
	w.RootItems.AddItem(i) // add to content
	return i
}
func (w *RootWindow) RemoveRootItem(i Item) bool {
	return w.RootItems.RemoveItem(i)
}

// returns the first item that is hit top down by a position
func (w *RootWindow) GetItemByPos(x, y int32, e sdl.Event) (bool, Item) {

	if i, found := w.RootItems.CheckTopDown(x, y); found {
		return i.GetItemByPos(x, y, e)
	}
	return false, nil
}

/*
delete a position i from array
copy(a[i:], a[i+1:])
a[len(a)-1] = nil // or the zero value of T
a = a[:len(a)-1]
*/

/*

item lokalisierung über position:

	items layered auf dem mainwindow
	items nested als content.

	- suche zuerst in den layered items des main windows. von oben nach unten. erster fund
	- in diesem Fund suche das nested item.


Item
	Ein Item muss auf Anfrage seine Minimalgröße bestimmen können


ItemGrid
	Für das Layout steht ein Grid-Item zur verfügung.
	Das Grid-Item ist ein Container für andere Items und hat keine sichtbaren Elemente.
	Es kann Einzelzelle, horizontale Zeile von zellen, vertikale spalte von zellen , oder zweidimensionales gitter von zellen sein.

	Eine Zelle kann nur ein Item enthalten.
	Zellengröße:
		- mindestens Minimalgröße
		- sollgröße in x oder y (mindestgröße überbietet maximalgröße)
		- ohne maximum -> so groß wie möglich expandiert


	Konkurierende Zellen bei Expansion
		- die Mindestgröße nach Minimalgröße des Inhalts
		- konkurierende Expansionen werden in den verbleibenden rest expandert
		- verbleibender rest wird anteilig in gewichtetn und und ungewichteten bereich geteilt.
		- ungewichteter Bereich wird an zellen gleichverteilt gewichteter bereich wird an gewichete zellen verteilt.

*/

// Layout all the items on the window starting from the root item.
// All root-items will be layed out in reference to the same root-window and may overlay each other
// The items of each root item will be layed out in reference to its root item.
func (w *RootWindow) Layout() {
	printSeparator()
	fmt.Printf("RootWindow.Layout() \n")
	wf := w.GetFrame()
	for _, ri := range w.RootItems.GetList() {
		ri.Layout(wf, true)
	}
}
func (w *RootWindow) ReportItems() {
	printSeparator()
	fmt.Printf("RootWindow.ReportItems() \n")
	for _, ri := range w.RootItems.GetList() {
		ri.Report(0)
	}
}

func (w *RootWindow) Destroy() {

	w.Renderer.Destroy()
	w.Renderer = nil

	w.Window.Destroy()
	w.Window = nil
}
func (w *RootWindow) GetFrame() *sdl.Rect {
	b, h := w.Window.GetSize()
	return &sdl.Rect{X: 0, Y: 0, W: b, H: h}
}

func (w *RootWindow) RenderAll() {
	//fmt.Printf("RootWindow.RenderAll() \n")

	utilRendererSetDrawColor(w.Renderer, &w.style.csDefault.baseDark)
	w.Renderer.Clear()

	// render from bootom to top

	for _, ri := range w.RootItems.GetList() {
		//utilRendererSetDrawColor(w.Renderer, w.style.colorBlack)
		//w.Renderer.FillRect(ri.GetFrame())
		ri.Render()
	}
	w.SetChanged(false)

}

func (w *RootWindow) Present() {
	// render texture to rootWindow
	//w.Renderer.SetRenderTarget(nil)
	//w.Renderer.Copy(w.winTex, nil, nil)
	//w.Renderer.SetRenderTarget(w.winTex)
	w.Renderer.Present()
}

func (w *RootWindow) handleMouseMotionEvent(e *sdl.MouseMotionEvent) {
	//Xrel and Yrel is not relieable due to fake moves sent by sdl when Mouse leaves window border with pressed buttons
	w.relMouseX = e.X - w.posMouseX
	w.relMouseY = e.Y - w.posMouseY
	w.posMouseX = e.X
	w.posMouseY = e.Y

	if w.mouseFocusLock {
		// allways send the motion events to the locked focus item
		if w.mouseFocusItem != nil {
			w.mouseFocusItem.oNotifyMouseMotion(e.X, e.Y, w.relMouseX, w.relMouseY)
		}
	} else {
		if found, i := w.GetItemByPos(e.X, e.Y, e); found {
			w.SetMouseFocusItem(i)
			i.oNotifyMouseMotion(e.X, e.Y, w.relMouseX, w.relMouseY)
		}
	}
}

func (w *RootWindow) handleMouseButtonEvent(e *sdl.MouseButtonEvent) {

	switch e.Button {
	case 1:
		if e.State == 0 {
			w.stateButton1 = false
		} else {
			w.stateButton1 = true
		}
	case 2:
		if e.State == 0 {
			w.stateButton2 = false
		} else {
			w.stateButton2 = true
		}
	case 3:
		if e.State == 0 {
			w.stateButton3 = false
		} else {
			w.stateButton3 = true
		}
	}

	if w.mouseFocusLock {
		// allways send the button events to the locked focus item
		if w.mouseFocusItem != nil {
			w.mouseFocusItem.oNotifyMouseButton(e.X, e.Y, e.Button, e.State)
		}
	} else {
		if found, i := w.GetItemByPos(e.X, e.Y, e); found {
			// send notification to the item
			i.oNotifyMouseButton(e.X, e.Y, e.Button, e.State)
			if w.kbFocusItem != i {
				// if the button event does not refer to the kb focus item, release the focus
				w.SetKbFocusItem(nil)
			}
		}
	}
}

func (w *RootWindow) handleWindowEvent(t *sdl.WindowEvent) {

	switch t.Event {
	case sdl.WINDOWEVENT_SHOWN:
		fmt.Printf("Window %d shown\n", t.WindowID)
	case sdl.WINDOWEVENT_HIDDEN:
		fmt.Printf("Window %d hidden\n", t.WindowID)
	case sdl.WINDOWEVENT_EXPOSED: //window has been exposed and should be redrawn
		fmt.Printf("Window %d exposed\n", t.WindowID)
		w.RenderAll()
		w.Present()
	case sdl.WINDOWEVENT_MOVED:
		fmt.Printf("Window %d moved to %d,%d\n", t.WindowID, t.Data1, t.Data2)
	case sdl.WINDOWEVENT_RESIZED: //SDL_WINDOWEVENT_RESIZED if the size was changed by an external event, i.e. the user or the window manager
		fmt.Printf("Window %d resized to %dx%d\n", t.WindowID, t.Data1, t.Data2)
		w.sizeX = t.Data1
		w.sizeY = t.Data2
		w.Layout()
	case sdl.WINDOWEVENT_SIZE_CHANGED: //window size has changed, either as a result of an API call or through the system or user changing the window size
		fmt.Printf("Window %d size changed to %dx%d\n", t.WindowID, t.Data1, t.Data2)
	case sdl.WINDOWEVENT_MINIMIZED:
		fmt.Printf("Window %d minimized\n", t.WindowID)
	case sdl.WINDOWEVENT_MAXIMIZED:
		fmt.Printf("Window %d maximized\n", t.WindowID)
	case sdl.WINDOWEVENT_RESTORED:
		fmt.Printf("Window %d restored\n", t.WindowID)
	case sdl.WINDOWEVENT_ENTER:
		fmt.Printf("Mouse entered window %d\n", t.WindowID)
	case sdl.WINDOWEVENT_LEAVE:
		fmt.Printf("Mouse left window %d\n", t.WindowID)
		w.SetMouseFocusItem(nil)
	case sdl.WINDOWEVENT_FOCUS_GAINED:
		fmt.Printf("Window %d gained keyboard focus\n", t.WindowID)
	case sdl.WINDOWEVENT_FOCUS_LOST:
		fmt.Printf("Window %d lost keyboard focus\n", t.WindowID)
	//case sdl.WINDOWEVENT_CLOSE:	 // CLOSE is handled by the window controller
	//fmt.Printf("Window %d closed\n", t.WindowID)
	case sdl.WINDOWEVENT_TAKE_FOCUS:
		fmt.Printf("Window %d is offered a focus\n", t.WindowID)
	case sdl.WINDOWEVENT_HIT_TEST:
		fmt.Printf("Window %d has a special hit test\n", t.WindowID)
	default:
		fmt.Printf("[%d ms] unknown Window event\twid:%d\tevent:%d\n", t.Timestamp, t.WindowID, t.Event)
	}
}

func (w *RootWindow) handleKeyboardEvent(e *sdl.KeyboardEvent) {
	if w.kbFocusItem != nil {
		w.kbFocusItem.oNotifyKbEvent(e)
	}
}

func (w *RootWindow) handleTextInputEvent(e *sdl.TextInputEvent) {
	if w.kbFocusItem != nil {
		//e.Text : the null-terminated input text in UTF-8 encoding
		// find the null termination
		nullIndex := strings.Index(string(e.Text[:]), "\x00")
		// get a rune slice from the string until null index
		r := []rune(string(e.Text[:nullIndex]))
		//TODO handle input text longer than one char
		w.kbFocusItem.oNotifyTextInput(r[0])
	}
}

func (w *RootWindow) SetMouseFocusItem(i Item) {
	if !w.mouseFocusLock {
		if w.mouseFocusItem != i {
			if w.mouseFocusItem != nil {
				w.mouseFocusItem.oNotifyMouseFocusLost()
			}
			w.mouseFocusItem = i
			if i != nil {
				w.mouseFocusItem.oNotifyMouseFocusGained()
			}
		}
	}
}
func (w *RootWindow) SetKbFocusItem(i Item) {

	if w.kbFocusItem != i {
		if w.kbFocusItem != nil {
			w.kbFocusItem.oNotifyKbFocusLost()
		}
		w.kbFocusItem = i
		if i != nil {
			i.oNotifyKbFocusGained()
		}
	}
}

//-----------------------------------------------------------

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

// -----------------------------------------------------------
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

type Style struct {
	Font *ttf.Font
	name string

	cursorBlinkRate int

	decoUnit         int32
	spacing          int32
	fontSize         int32
	averageFontHight int32
	averageRuneWidth int32

	colorRed    *sdl.Color
	colorGreen  *sdl.Color
	colorBlue   *sdl.Color
	colorPurple *sdl.Color
	colorBlack  *sdl.Color
	colorWhite  *sdl.Color

	//colorBg   *sdl.Color
	csDefault *ColorScheme
	csWindow  *ColorScheme
}

func NewStyle(n string) *Style {
	s := new(Style)

	if !ttf.WasInit() {
		ttf.Init()
	}

	s.InitDefault()
	s.name = n
	return s
}

func (s *Style) InitDefault() {
	s.cursorBlinkRate = 300
	//red := &sdl.Color{255, 0, 0, 255}
	// R, G, B, A
	s.colorRed = &sdl.Color{R: 255, G: 50, B: 50, A: 255}
	s.colorGreen = &sdl.Color{R: 5, G: 200, B: 5, A: 255}
	s.colorBlue = &sdl.Color{R: 50, G: 50, B: 255, A: 255}
	s.colorPurple = &sdl.Color{R: 255, G: 50, B: 255, A: 255}
	s.colorBlack = &sdl.Color{R: 0, G: 0, B: 0, A: 255}
	s.colorWhite = &sdl.Color{R: 255, G: 255, B: 255, A: 255}

	//standard item
	s.csDefault = new(ColorScheme)
	//s.csDefault.SetBaseColor(&sdl.Color{70, 100, 100, 255}, 50)
	s.csDefault.SetBaseColor(&sdl.Color{R: 190, G: 190, B: 190, A: 255})

	// Windows colors
	s.csWindow = new(ColorScheme)
	s.csWindow.SetBaseColor(&sdl.Color{R: 150, G: 200, B: 150, A: 255})

	s.decoUnit = 16
	// spacing shall not drop below 2 !
	s.spacing = 2
	s.fontSize = 13

	var err error
	// find /usr/share/fonts/truetype -name '*.ttf'
	// /usr/share/fonts/truetype/dejavu/DejaVuSans.ttf
	// /usr/share/fonts/truetype/ubuntu/Ubuntu-RI.ttf

	//fmt.Printf("Available Fonts: in /usr/share/fonts/truetype/\n")
	//files, err := filepath.Glob("/usr/share/fonts/truetype/*/*.ttf")
	//for _, v := range files {
	//	fmt.Printf("%s\n", v)
	//}

	//s.Font, err = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", int(s.fontSize))
	//s.Font, err = ttf.OpenFont("/usr/share/fonts/truetype/ubuntu/Ubuntu-RI.ttf", int(s.fontSize))
	//s.Font, err = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", int(s.fontSize))
	s.Font, err = ttf.OpenFont("/usr/share/fonts/truetype/ubuntu/Ubuntu-R.ttf", int(s.fontSize))
	if err != nil {
		panic(err)
	}
	s.Font.SetHinting(ttf.HINTING_NORMAL)
	//s.Font.SetKerning(false)

	s.averageFontHight = s.GetTextHight("QqJjPpGg")
	s.averageRuneWidth = s.GetAverageRuneWidth("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQTSTUVWXYZ")
}

func (s *Style) Destroy() {
	s.Font.Close()
}
func (s *Style) GetAverageRuneWidth(txt string) int32 {
	var sum, n int
	for _, r := range txt {
		x, _, _ := s.Font.SizeUTF8(string(r))
		n++
		sum += x
	}
	return int32(sum / n)

}

func (s *Style) GetTextLen(txt string) int32 {
	x, _, _ := s.Font.SizeUTF8(txt)
	return int32(x)
}

func (s *Style) GetTextHight(txt string) int32 {
	_, x, _ := s.Font.SizeUTF8(txt)
	return int32(x)
}

// -----------------------------------------------
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
