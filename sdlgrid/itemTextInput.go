package sdlgrid

import (
	"fmt"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

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

func (i *ItemTextInput) oNotifyMouseButton(x, y int32, button uint8, state uint8) {
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
