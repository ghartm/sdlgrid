package sdlgrid

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

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

func (i *ItemMenu) layerButtonCb(x, y int32, button uint8, state uint8) {
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
