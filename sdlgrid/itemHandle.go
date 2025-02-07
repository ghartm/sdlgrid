package sdlgrid

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

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

func (i *ItemHandle) oNotifyMouseButton(x, y int32, button uint8, state uint8) {
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

	//utilRenderSolidBorder(r, &i.iframe, i.GetStyle().ColorPurple)

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
