package sdlgrid

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

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
	//i.margin=win.Style.Spacing
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
