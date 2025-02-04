package sdlgrid

import "github.com/veandco/go-sdl2/sdl"

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
