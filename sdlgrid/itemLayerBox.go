package sdlgrid

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// ItemLayerBox is designed manage layered items
// all items share the same frame
type ItemLayerBox struct {
	ItemBase
	subItems       infraItemList
	buttonCallback func(int32, int32, uint8, uint8) //x, y int32, button uint8, state uint8
	//spacing        int32
}

func NewItemLayerBox(win *RootWindow) *ItemLayerBox {
	i := new(ItemLayerBox)
	i.o = Item(i)
	i.setRootWindow(win)
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_PCT, 0, 0, 100000, 100000)
	return i
}

func (i *ItemLayerBox) SetButtonCallback(cb func(int32, int32, uint8, uint8)) {
	i.buttonCallback = cb
}
func (i *ItemLayerBox) oNotifyMouseButton(x, y int32, button uint8, state uint8) {
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
				x, y, _, _ := si.GetCollapsedSpec()
				w, h := si.oGetMinSize()
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
