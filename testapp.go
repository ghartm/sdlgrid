package main

import (
	"fmt"
	"sdlgrid/sdlgrid"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func newTestWindow(win *sdlgrid.RootWindow) *sdlgrid.ItemWindow {
	// Window
	win1 := sdlgrid.NewItemWindow(win)
	win1.SetSpec(sdlgrid.LS_POS_PCT, sdlgrid.LS_POS_PCT, sdlgrid.LS_SIZE_PCT, sdlgrid.LS_SIZE_COLLAPSE, 20000, 50000, 50000, 0)

	rg := sdlgrid.NewItemGrid(win, 0, 0)
	win1.AddSubItem(rg)

	sl := sdlgrid.NewItemSlider(win, sdlgrid.ITEM_SLIDER_STYLE_VERTICAL)
	rg.SetSubItem(rg.AppendColumn(), 0, sl)

	frm11 := sdlgrid.NewItemFrame(win)
	rg.SetSubItem(rg.AppendColumn(), 0, frm11)

	rg.SetColSpec(0, layoutParam{sdlgrid.LS_SIZE_ABS, win.GetStyle().decoUnit})

	frm11.SetName("wFrame")
	frm11.SetText("wFrame")
	frm11.SetSpec(sdlgrid.LS_POS_PCT, sdlgrid.LS_POS_PCT, sdlgrid.LS_SIZE_COLLAPSE, sdlgrid.LS_SIZE_COLLAPSE, 50000, 50000, 0, 0)

	gr3 := sdlgrid.NewItemGrid(win, 0, 0)
	gr3.SetName("wGrid")
	gr3.SetSpacing(win.style.spacing)
	frm11.SetSubItem(gr3)

	ti1 := sdlgrid.NewItemTextInput(win, win.style.averageRuneWidth*20)
	ti1.SetName("textInput1")

	ti2 := sdlgrid.NewItemTextInput(win, 80)
	ti2.SetName("textInput2")
	gr3.SetSubItem(0, gr3.AppendRow(), ti1)
	gr3.SetSubItem(0, gr3.AppendRow(), ti2)

	b1 := sdlgrid.NewItemButton(win)
	gr3.SetSubItem(0, gr3.AppendRow(), b1)
	b1.SetName("firstButton")
	b1.SetText("Change Alignment")
	b1.SetCallback(func() {
		fmt.Printf("Button Callback: %s\n", b1.GetName())
		ti1.SetTextAlignment(sdlgrid.ITEM_TEXTINPUT_ALIGN_CENTER)
		ti2.SetTextAlignment(sdlgrid.ITEM_TEXTINPUT_ALIGN_RIGHT)
	})

	b1 = sdlgrid.NewItemButton(win)
	gr3.SetSubItem(0, gr3.AppendRow(), b1)
	b1.SetText("toggle Input")
	b1.SetCallback(func() {
		ti1.SetHidden(!ti1.hidden)
		ti2.SetHidden(!ti2.hidden)
	})

	gr4 := sdlgrid.NewItemGrid(win, 0, 0)
	gr3.SetSubItem(0, gr3.AppendRow(), gr4)
	gr4.SetSpacing(win.style.spacing)

	lr := sdlgrid.NewItemTextInput(win, 50)
	gr4.SetSubItem(gr4.AppendColumn(), 0, lr)
	lr.SetTextAlignment(sdlgrid.ITEM_TEXTINPUT_ALIGN_RIGHT)
	lr.SetText("50")

	lg := sdlgrid.NewItemTextInput(win, 50)
	gr4.SetSubItem(gr4.AppendColumn(), 0, lg)
	lg.SetTextAlignment(sdlgrid.ITEM_TEXTINPUT_ALIGN_RIGHT)
	lg.SetText("50")

	lb := sdlgrid.NewItemTextInput(win, 50)
	gr4.SetSubItem(gr4.AppendColumn(), 0, lb)
	lb.SetTextAlignment(sdlgrid.ITEM_TEXTINPUT_ALIGN_RIGHT)
	lb.SetText("50")

	b1 = sdlgrid.NewItemButton(win)
	gr3.SetSubItem(0, gr3.AppendRow(), b1)
	b1.SetText("set color")
	b1.SetCallback(func() {
		ncs := new(sdlgrid.ColorScheme)
		ncs.SetBaseColor(&sdl.Color{R: uint8(lr.GetNumber()), G: uint8(lg.GetNumber()), B: uint8(lb.GetNumber()), A: 255})
		win1.oWithItems(func(i sdlgrid.Item) {
			i.SetColorScheme(ncs)
			i.SetChanged(true)
		})
	})
	return win1
}

func main() {

	fmt.Println("start")

	wc := sdlgrid.NewWindowController()
	defer wc.Destroy()

	s := sdlgrid.NewStyle("default")
	defer s.Destroy()

	win := sdlgrid.NewRootWindow(s, "RootWindow", 500, 300)
	defer win.Destroy()

	wc.AddRootWindow(win)

	gr1 := sdlgrid.NewItemGrid(win, 1, 2)
	gr1.SetName("rootGrid")
	//gr1.SetRowSpec(0, 50)
	//gr1.SetRowSpec(1, 100)
	//gr1.SetSpacing(2)
	win.AddRootItem(gr1)

	// menue entries
	mu1 := sdlgrid.NewItemMenue(win)
	gr1.SetSubItem(0, 0, mu1)
	mu1.SetName("mu1")
	mu1.SetOrientation(sdlgrid.MENU_ORIENTATION_HORIZONTAL)

	wm := sdlgrid.NewItemWindowManager(win)
	gr1.SetSubItem(0, 1, wm)

	mu2 := sdlgrid.NewItemMenue(win)
	mu1.AddNewMenuEntry("File", mu2, nil)
	mu2.SetName("mu2")
	mu2.AddNewMenuEntry("New", nil, nil)
	mu2.AddNewMenuEntry("Open", nil, nil)
	mu2.AddNewMenuEntry("Save", nil, nil)
	mu2.AddNewMenuEntry("Close", nil, nil)

	mu3 := sdlgrid.NewItemMenue(win)
	mu1.AddNewMenuEntry("Window", mu3, nil)
	mu3.AddNewMenuEntry("New", nil, func() {
		wm.AddSubItem(newTestWindow(win))

	})

	w1 := newTestWindow(win)
	wm.AddSubItem(w1)

	win.ReportItems()
	wc.Start()

}

func eventSenderCustom(wg *sync.WaitGroup, c chan interface{}, ctrl chan bool) {
	wg.Add(1)
	var e sdlgrid.EventCustom
	var running bool = true
	for running {
		select {
		case _, running = <-ctrl:
			fmt.Println("quitting eventSenderCustom")
			wg.Done()
		case c <- &e:
			time.Sleep(time.Duration(1000) * time.Millisecond)
		}
	}
}

func eventSenderTick(wg *sync.WaitGroup, c chan interface{}, ctrl chan bool, msec int) {
	wg.Add(1)
	e := sdlgrid.EventRenderTick{msec: msec}
	var running bool = true
	for running {
		select {
		case _, running = <-ctrl:
			fmt.Println("quitting eventSenderTick")
			wg.Done()
		default:
			c <- &e
			time.Sleep(time.Duration(msec) * time.Millisecond)
		}
	}
}

//--------------------------------------------------------------------

//--------------------------------------------------------------------
