package sdlgrid

import (
	"fmt"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

type RootWindow struct {
	id             uint32
	title          string
	sizeX          int32
	sizeY          int32
	Window         *sdl.Window
	Renderer       *sdl.Renderer
	Style          *Style
	RootItems      infraItemList
	mouseFocusItem Item
	mouseFocusLock bool
	kbFocusItem    Item
	changed        bool
	stateButton1   bool
	stateButton2   bool
	stateButton3   bool
	posMouseX      int32
	posMouseY      int32
	relMouseX      int32
	relMouseY      int32
}

func NewRootWindow(s *Style, title string, x int32, y int32) *RootWindow {
	w := new(RootWindow)
	w.sizeX = x
	w.sizeY = y
	w.title = title
	w.Style = s

	var err error

	w.Window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, w.sizeX, w.sizeY, sdl.WINDOW_RESIZABLE|sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	w.Renderer, err = sdl.CreateRenderer(w.Window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_TARGETTEXTURE)
	if err != nil {
		panic(err)
	}

	w.Renderer.SetDrawColor(0, 0, 0, 255)
	w.Renderer.Clear()

	w.id, _ = w.Window.GetID()
	w.mouseFocusLock = false

	return w

}

func (w *RootWindow) SetMouseFocusLock(b bool) {
	w.mouseFocusLock = b
	if b {
		fmt.Printf("Mouse focus locked\n")
	} else {
		fmt.Printf("Mouse focus unlocked\n")
		// if it was unlocked outside the focus item notify the new item at the last known position for focus change
		if found, i := w.GetItemByPos(w.posMouseX, w.posMouseY, nil); found {
			w.SetMouseFocusItem(i)
		}

	}
}
func (w *RootWindow) GetChanged() bool  { return w.changed }
func (w *RootWindow) SetChanged(c bool) { w.changed = c }

func (w *RootWindow) GetStyle() *Style  { return w.Style }
func (w *RootWindow) SetStyle(s *Style) { w.Style = s }

func (w *RootWindow) AddRootItem(i Item) Item {
	w.RootItems.AddItem(i) // add to content
	return i
}
func (w *RootWindow) RemoveRootItem(i Item) bool {
	return w.RootItems.RemoveItem(i)
}

// returns the first item that is hit top down by a position
func (w *RootWindow) GetItemByPos(x, y int32, e sdl.Event) (bool, Item) {

	if i, found := w.RootItems.CheckTopDown(x, y); found {
		return i.GetItemByPos(x, y, e)
	}
	return false, nil
}

/*
delete a position i from array
copy(a[i:], a[i+1:])
a[len(a)-1] = nil // or the zero value of T
a = a[:len(a)-1]
*/

/*

item lokalisierung über position:

	items layered auf dem mainwindow
	items nested als content.

	- suche zuerst in den layered items des main windows. von oben nach unten. erster fund
	- in diesem Fund suche das nested item.


Item
	Ein Item muss auf Anfrage seine Minimalgröße bestimmen können


ItemGrid
	Für das Layout steht ein Grid-Item zur verfügung.
	Das Grid-Item ist ein Container für andere Items und hat keine sichtbaren Elemente.
	Es kann Einzelzelle, horizontale Zeile von zellen, vertikale spalte von zellen , oder zweidimensionales gitter von zellen sein.

	Eine Zelle kann nur ein Item enthalten.
	Zellengröße:
		- mindestens Minimalgröße
		- sollgröße in x oder y (mindestgröße überbietet maximalgröße)
		- ohne maximum -> so groß wie möglich expandiert


	Konkurierende Zellen bei Expansion
		- die Mindestgröße nach Minimalgröße des Inhalts
		- konkurierende Expansionen werden in den verbleibenden rest expandert
		- verbleibender rest wird anteilig in gewichtetn und und ungewichteten bereich geteilt.
		- ungewichteter Bereich wird an zellen gleichverteilt gewichteter bereich wird an gewichete zellen verteilt.

*/

// Layout all the items on the window starting from the root item.
// All root-items will be layed out in reference to the same root-window and may overlay each other
// The items of each root item will be layed out in reference to its root item.
func (w *RootWindow) Layout() {
	printSeparator()
	fmt.Printf("RootWindow.Layout() \n")
	wf := w.GetFrame()
	for _, ri := range w.RootItems.GetList() {
		ri.Layout(wf, true)
	}
}
func (w *RootWindow) ReportItems() {
	printSeparator()
	fmt.Printf("RootWindow.ReportItems() \n")
	for _, ri := range w.RootItems.GetList() {
		ri.Report(0)
	}
}

func (w *RootWindow) Destroy() {

	w.Renderer.Destroy()
	w.Renderer = nil

	w.Window.Destroy()
	w.Window = nil
}
func (w *RootWindow) GetFrame() *sdl.Rect {
	b, h := w.Window.GetSize()
	return &sdl.Rect{X: 0, Y: 0, W: b, H: h}
}

func (w *RootWindow) RenderAll() {
	//fmt.Printf("RootWindow.RenderAll() \n")

	utilRendererSetDrawColor(w.Renderer, &w.Style.csDefault.baseDark)
	w.Renderer.Clear()

	// render from bootom to top

	for _, ri := range w.RootItems.GetList() {
		//utilRendererSetDrawColor(w.Renderer, w.Style.colorBlack)
		//w.Renderer.FillRect(ri.GetFrame())
		ri.Render()
	}
	w.SetChanged(false)

}

func (w *RootWindow) Present() {
	// render texture to rootWindow
	//w.Renderer.SetRenderTarget(nil)
	//w.Renderer.Copy(w.winTex, nil, nil)
	//w.Renderer.SetRenderTarget(w.winTex)
	w.Renderer.Present()
}

func (w *RootWindow) handleMouseMotionEvent(e *sdl.MouseMotionEvent) {
	//Xrel and Yrel is not relieable due to fake moves sent by sdl when Mouse leaves window border with pressed buttons
	w.relMouseX = e.X - w.posMouseX
	w.relMouseY = e.Y - w.posMouseY
	w.posMouseX = e.X
	w.posMouseY = e.Y

	if w.mouseFocusLock {
		// allways send the motion events to the locked focus item
		if w.mouseFocusItem != nil {
			w.mouseFocusItem.oNotifyMouseMotion(e.X, e.Y, w.relMouseX, w.relMouseY)
		}
	} else {
		if found, i := w.GetItemByPos(e.X, e.Y, e); found {
			w.SetMouseFocusItem(i)
			i.oNotifyMouseMotion(e.X, e.Y, w.relMouseX, w.relMouseY)
		}
	}
}

func (w *RootWindow) handleMouseButtonEvent(e *sdl.MouseButtonEvent) {

	switch e.Button {
	case 1:
		if e.State == 0 {
			w.stateButton1 = false
		} else {
			w.stateButton1 = true
		}
	case 2:
		if e.State == 0 {
			w.stateButton2 = false
		} else {
			w.stateButton2 = true
		}
	case 3:
		if e.State == 0 {
			w.stateButton3 = false
		} else {
			w.stateButton3 = true
		}
	}

	if w.mouseFocusLock {
		// allways send the button events to the locked focus item
		if w.mouseFocusItem != nil {
			w.mouseFocusItem.oNotifyMouseButton(e.X, e.Y, e.Button, e.State)
		}
	} else {
		if found, i := w.GetItemByPos(e.X, e.Y, e); found {
			// send notification to the item
			i.oNotifyMouseButton(e.X, e.Y, e.Button, e.State)
			if w.kbFocusItem != i {
				// if the button event does not refer to the kb focus item, release the focus
				w.SetKbFocusItem(nil)
			}
		}
	}
}

func (w *RootWindow) handleWindowEvent(t *sdl.WindowEvent) {

	switch t.Event {
	case sdl.WINDOWEVENT_SHOWN:
		fmt.Printf("Window %d shown\n", t.WindowID)
	case sdl.WINDOWEVENT_HIDDEN:
		fmt.Printf("Window %d hidden\n", t.WindowID)
	case sdl.WINDOWEVENT_EXPOSED: //window has been exposed and should be redrawn
		fmt.Printf("Window %d exposed\n", t.WindowID)
		w.RenderAll()
		w.Present()
	case sdl.WINDOWEVENT_MOVED:
		fmt.Printf("Window %d moved to %d,%d\n", t.WindowID, t.Data1, t.Data2)
	case sdl.WINDOWEVENT_RESIZED: //SDL_WINDOWEVENT_RESIZED if the size was changed by an external event, i.e. the user or the window manager
		fmt.Printf("Window %d resized to %dx%d\n", t.WindowID, t.Data1, t.Data2)
		w.sizeX = t.Data1
		w.sizeY = t.Data2
		w.Layout()
	case sdl.WINDOWEVENT_SIZE_CHANGED: //window size has changed, either as a result of an API call or through the system or user changing the window size
		fmt.Printf("Window %d size changed to %dx%d\n", t.WindowID, t.Data1, t.Data2)
	case sdl.WINDOWEVENT_MINIMIZED:
		fmt.Printf("Window %d minimized\n", t.WindowID)
	case sdl.WINDOWEVENT_MAXIMIZED:
		fmt.Printf("Window %d maximized\n", t.WindowID)
	case sdl.WINDOWEVENT_RESTORED:
		fmt.Printf("Window %d restored\n", t.WindowID)
	case sdl.WINDOWEVENT_ENTER:
		fmt.Printf("Mouse entered window %d\n", t.WindowID)
	case sdl.WINDOWEVENT_LEAVE:
		fmt.Printf("Mouse left window %d\n", t.WindowID)
		w.SetMouseFocusItem(nil)
	case sdl.WINDOWEVENT_FOCUS_GAINED:
		fmt.Printf("Window %d gained keyboard focus\n", t.WindowID)
	case sdl.WINDOWEVENT_FOCUS_LOST:
		fmt.Printf("Window %d lost keyboard focus\n", t.WindowID)
	//case sdl.WINDOWEVENT_CLOSE:	 // CLOSE is handled by the window controller
	case sdl.WINDOWEVENT_TAKE_FOCUS:
		fmt.Printf("Window %d is offered a focus\n", t.WindowID)
	case sdl.WINDOWEVENT_HIT_TEST:
		fmt.Printf("Window %d has a special hit test\n", t.WindowID)
	default:
		fmt.Printf("[%d ms] unknown Window event\twid:%d\tevent:%d\n", t.Timestamp, t.WindowID, t.Event)
	}
}

func (w *RootWindow) handleKeyboardEvent(e *sdl.KeyboardEvent) {
	if w.kbFocusItem != nil {
		w.kbFocusItem.oNotifyKbEvent(e)
	}
}

func (w *RootWindow) handleTextInputEvent(e *sdl.TextInputEvent) {
	if w.kbFocusItem != nil {
		//e.Text : the null-terminated input text in UTF-8 encoding
		// find the null termination
		nullIndex := strings.Index(string(e.Text[:]), "\x00")
		// get a rune slice from the string until null index
		r := []rune(string(e.Text[:nullIndex]))
		//TODO handle input text longer than one char
		w.kbFocusItem.oNotifyTextInput(r[0])
	}
}

func (w *RootWindow) SetMouseFocusItem(i Item) {
	if !w.mouseFocusLock {
		if w.mouseFocusItem != i {
			if w.mouseFocusItem != nil {
				w.mouseFocusItem.oNotifyMouseFocusLost()
			}
			w.mouseFocusItem = i
			if i != nil {
				w.mouseFocusItem.oNotifyMouseFocusGained()
			}
		}
	}
}
func (w *RootWindow) SetKbFocusItem(i Item) {

	if w.kbFocusItem != i {
		if w.kbFocusItem != nil {
			w.kbFocusItem.oNotifyKbFocusLost()
		}
		w.kbFocusItem = i
		if i != nil {
			i.oNotifyKbFocusGained()
		}
	}
}
