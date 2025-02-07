package sdlgrid

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

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
		mgn := i.GetStyle().Spacing

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

func (i *ItemButton) oNotifyMouseButton(x, y int32, button uint8, state uint8) {
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
