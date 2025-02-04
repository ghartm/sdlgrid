package sdlgrid

import "fmt"

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
