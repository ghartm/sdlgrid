package sdlgrid

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

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
