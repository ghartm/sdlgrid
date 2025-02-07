package sdlgrid

import "github.com/veandco/go-sdl2/sdl"

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

	size := win.GetStyle().DecoUnit
	i.userUnit = size

	if i.sliderStyle == ITEM_SLIDER_STYLE_HORIZONTAL {
		i.upLeftButton.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_ABS, LS_SIZE_PCT, 0, 0, size, 0)
		i.downRightButton.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_ABS, LS_SIZE_PCT, 100000, 0, size, 100000)
		i.handle.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_ABS, LS_SIZE_PCT, 0, 0, size, 100000)

		n := i.rootGrid.AppendColumn()
		i.rootGrid.SetSubItem(n, 0, i.downRightButton)
		i.rootGrid.SetColSpec(n, LayoutParam{LS_SIZE_COLLAPSE, 0})
		n = i.rootGrid.AppendColumn()
		i.rootGrid.SetSubItem(n, 0, i.handle)
		i.rootGrid.SetColSpec(n, LayoutParam{LS_SIZE_PCT, 100000})
		n = i.rootGrid.AppendColumn()
		i.rootGrid.SetSubItem(n, 0, i.downRightButton)
		i.rootGrid.SetColSpec(n, LayoutParam{LS_SIZE_COLLAPSE, 0})
	} else {
		i.upLeftButton.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_ABS, 0, 0, 100000, size)
		i.downRightButton.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_ABS, 0, 100000, 100000, size)
		i.handle.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_ABS, 0, 0, 100000, size)

		n := i.rootGrid.AppendRow()
		i.rootGrid.SetSubItem(0, n, i.upLeftButton)
		i.rootGrid.SetRowSpec(n, LayoutParam{LS_SIZE_COLLAPSE, 0})
		n = i.rootGrid.AppendRow()
		i.rootGrid.SetSubItem(0, n, i.handle)
		i.rootGrid.SetRowSpec(n, LayoutParam{LS_SIZE_PCT, 100000})
		n = i.rootGrid.AppendRow()
		i.rootGrid.SetSubItem(0, n, i.downRightButton)
		i.rootGrid.SetRowSpec(n, LayoutParam{LS_SIZE_COLLAPSE, 0})
	}
	return i
}

func (i *ItemSlider) computeUnitSize() {
	// according handles ability to move, handle size and total and part values
	// compute the jump unit in pixel that handle will jump on a decrease or increase
	var pix int32
	if i.sliderStyle == ITEM_SLIDER_STYLE_HORIZONTAL {
		pix = i.handle.pframe.W
	} else {
		pix = i.handle.pframe.H
	}
	i.userTotalPerPx = float32(i.userTotal) / float32(pix)

	s := i.GetStyle().DecoUnit
	if i.handleSize = int32(float32(i.userView) / i.userTotalPerPx); i.handleSize < s {
		i.handleSize = s
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
