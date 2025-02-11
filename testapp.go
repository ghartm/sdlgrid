package main

import (
	"fmt"
	"sdlgrid/sdlgrid"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func newTestWindow2(win *sdlgrid.RootWindow) *sdlgrid.ItemWindow {
	win1 := sdlgrid.NewItemWindow(win)
	win1.SetSpec(sdlgrid.LS_POS_PCT, sdlgrid.LS_POS_PCT, sdlgrid.LS_SIZE_COLLAPSE, sdlgrid.LS_SIZE_COLLAPSE, 20000, 20000, 0, 0)
	rg := sdlgrid.NewItemGrid(win, 0, 0)
	rg.SetSpacing(win.Style.Spacing)

	win1.AddSubItem(rg)

	//r1 := rg.AppendRow()
	//r2 := rg.AppendRow()

	c := rg.AppendColumn()
	rg.SetColSpec(c, sdlgrid.LayoutParam{S: sdlgrid.LS_SIZE_PCT, V: 100000})
	rg.SetSubItem(c, 0, sdlgrid.NewItemFrame(win).SetText("frame0"))
	r1 := rg.AppendRow()
	rg.SetRowSpec(r1, sdlgrid.LayoutParam{S: sdlgrid.LS_SIZE_PCT, V: 100000})
	rg.SetSubItem(c, r1, sdlgrid.NewItemFrame(win).SetText("frame1"))
	r2 := rg.AppendRow()
	rg.SetRowSpec(r2, sdlgrid.LayoutParam{S: sdlgrid.LS_SIZE_PCT, V: 100000})
	rg.SetSubItem(c, r2, sdlgrid.NewItemFrame(win).SetText("frame2"))

	c = rg.AppendColumn()
	rg.SetColSpec(c, sdlgrid.LayoutParam{S: sdlgrid.LS_SIZE_PCT, V: 100000})
	rg.SetSubItem(c, 0, sdlgrid.NewItemFrame(win).SetText("frame3"))
	rg.SetSubItem(c, r1, sdlgrid.NewItemFrame(win).SetText("frame4"))
	rg.SetSubItem(c, r2, sdlgrid.NewItemFrame(win).SetText("frame5"))

	c = rg.AppendColumn()
	rg.SetColSpec(c, sdlgrid.LayoutParam{S: sdlgrid.LS_SIZE_PCT, V: 100000})
	rg.SetSubItem(c, 0, sdlgrid.NewItemFrame(win).SetText("frame6"))
	rg.SetSubItem(c, r1, sdlgrid.NewItemFrame(win).SetText("frame7"))
	rg.SetSubItem(c, r2, sdlgrid.NewItemFrame(win).SetText("frame8"))

	return win1
}

func newTestWindow1(win *sdlgrid.RootWindow) *sdlgrid.ItemWindow {
	// Window
	win1 := sdlgrid.NewItemWindow(win)
	win1.SetSpec(sdlgrid.LS_POS_PCT, sdlgrid.LS_POS_PCT, sdlgrid.LS_SIZE_PCT, sdlgrid.LS_SIZE_COLLAPSE, 20000, 50000, 50000, 0)

	rg := sdlgrid.NewItemGrid(win, 0, 0)
	win1.AddSubItem(rg)

	sl := sdlgrid.NewItemSlider(win, sdlgrid.ITEM_SLIDER_STYLE_VERTICAL)
	rg.SetSubItem(rg.AppendColumn(), 0, sl)

	frm11 := sdlgrid.NewItemFrame(win)
	rg.SetSubItem(rg.AppendColumn(), 0, frm11)

	rg.SetColSpec(0, sdlgrid.LayoutParam{S: sdlgrid.LS_SIZE_ABS, V: win.GetStyle().DecoUnit})

	frm11.SetName("wFrame")
	frm11.SetText("wFrame")
	frm11.SetSpec(sdlgrid.LS_POS_PCT, sdlgrid.LS_POS_PCT, sdlgrid.LS_SIZE_COLLAPSE, sdlgrid.LS_SIZE_COLLAPSE, 50000, 50000, 0, 0)

	gr3 := sdlgrid.NewItemGrid(win, 0, 0)
	gr3.SetName("wGrid")
	gr3.SetSpacing(win.Style.Spacing)
	frm11.SetSubItem(gr3)

	ti1 := sdlgrid.NewItemTextInput(win, win.Style.EstimateTextLength(20))
	ti1.SetName("textInput1")

	ti2 := sdlgrid.NewItemTextInput(win, 80)
	ti2.SetName("textInput2")
	gr3.SetSubItem(0, gr3.AppendRow(), ti1)
	gr3.SetSubItem(0, gr3.AppendRow(), ti2)

	b1 := sdlgrid.NewItemButton(win)
	gr3.SetSubItem(0, gr3.AppendRow(), b1)
	b1.SetName("firstButton")
	b1.SetText("Change Alignment")
	var aligntoggle int = 1
	b1.SetCallback(func() {
		fmt.Printf("Button Callback: %s\n", b1.GetName())
		if aligntoggle == 0 {
			ti1.SetTextAlignment(sdlgrid.ITEM_TEXTINPUT_ALIGN_LEFT)
			ti2.SetTextAlignment(sdlgrid.ITEM_TEXTINPUT_ALIGN_LEFT)
		} else if aligntoggle == 1 {
			ti1.SetTextAlignment(sdlgrid.ITEM_TEXTINPUT_ALIGN_CENTER)
			ti2.SetTextAlignment(sdlgrid.ITEM_TEXTINPUT_ALIGN_CENTER)
		} else if aligntoggle == 2 {
			ti1.SetTextAlignment(sdlgrid.ITEM_TEXTINPUT_ALIGN_RIGHT)
			ti2.SetTextAlignment(sdlgrid.ITEM_TEXTINPUT_ALIGN_RIGHT)
		}

		aligntoggle = (aligntoggle + 1) % 3
	})

	b1 = sdlgrid.NewItemButton(win)
	gr3.SetSubItem(0, gr3.AppendRow(), b1)
	b1.SetText("toggle Input")
	b1.SetCallback(func() {
		ti1.SetHidden(!ti1.GetHidden())
		ti2.SetHidden(!ti2.GetHidden())
	})

	gr4 := sdlgrid.NewItemGrid(win, 0, 0)
	gr3.SetSubItem(0, gr3.AppendRow(), gr4)
	gr4.SetSpacing(win.Style.Spacing)

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
		win1.ChangeAllColourSchemes(ncs)
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
	mu3.AddNewMenuEntry("New 1", nil, func() {
		wm.AddSubItem(newTestWindow1(win))

	})
	mu3.AddNewMenuEntry("New 2", nil, func() {
		wm.AddSubItem(newTestWindow2(win))

	})
	mu3.AddNewMenuEntry("SubMenue File >", mu2, nil)

	w1 := newTestWindow1(win)
	wm.AddSubItem(w1)

	win.ReportItems()

	wc.RegisterCustomEventSender(eventSenderCustom)
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

//--------------------------------------------------------------------

//--------------------------------------------------------------------
