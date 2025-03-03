package sdlgrid

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type EventCustom interface{}

type EventRenderTick struct {
	Msec int
}

type WindowController struct {
	windows     []*RootWindow
	ctrl        chan bool
	queue       chan interface{}
	waitgroup   *sync.WaitGroup
	lastId      uint32
	lastWindow  *RootWindow
	eventFnList []func(*sync.WaitGroup, chan interface{}, chan bool)
}

func NewWindowController() *WindowController {
	wc := new(WindowController)
	wc.eventFnList = make([]func(*sync.WaitGroup, chan interface{}, chan bool), 0)

	if sdl.WasInit(sdl.INIT_VIDEO) != sdl.INIT_VIDEO {
		sdl.Init(sdl.INIT_VIDEO)
	}

	return wc
}

func (wc *WindowController) Destroy() {

	// check if control channel ist closed by reading from it. if not closed close it
	_, running := <-wc.ctrl
	if running {
		close(wc.ctrl)
	}

	if ttf.WasInit() {
		ttf.Quit()
	}

	if sdl.WasInit(sdl.INIT_VIDEO) == sdl.INIT_VIDEO {
		sdl.Quit()
	}
}

func (wc *WindowController) AddRootWindow(rw *RootWindow) {
	wc.windows = append(wc.windows, rw)
}

func (wc *WindowController) getRootWindowById(id uint32) (bool, *RootWindow) {
	if id == wc.lastId {
		return true, wc.lastWindow
	} else {
		for _, w := range wc.windows {
			if i, _ := w.Window.GetID(); i == id {
				wc.lastId = i
				wc.lastWindow = w
				return true, w
			}
		}
	}
	return false, nil
}

func (wc *WindowController) eventSenderRenderTick(ms int) {
	wc.waitgroup.Add(1)
	e := EventRenderTick{Msec: ms}
	var running bool = true
runloop:
	for {
		select {
		case _, running = <-wc.ctrl:
			if !running {
				fmt.Println("quitting eventSenderTick")
				wc.waitgroup.Done()
				break runloop
			}
		default:
			wc.queue <- &e
			time.Sleep(time.Duration(ms) * time.Millisecond)
		}
	}
}

func (wc *WindowController) eventSenderSDL() {
	wc.waitgroup.Add(1)
	var running bool = true

runloop:
	for {
		// A select blocks until one of its cases can run, then it executes that case. It chooses one at random if multiple are ready.
		select {
		case _, running = <-wc.ctrl:
			if !running {
				fmt.Println("quitting eventSenderSDL")
				wc.waitgroup.Done()
				break runloop
			}
		case wc.queue <- sdl.WaitEventTimeout(200): // wait here until an sdl-event arives. pick one.
			// if timeout occurs no event will be returned but "nil" - the main loop needs to compensate for nil events.
			// As long as WaitEventTimeout blocks, the wc.ctrl channel can not be read and will block exits.
		}
	}
}

func (wc *WindowController) RegisterCustomEventSender(eventfn func(*sync.WaitGroup, chan interface{}, chan bool)) {
	wc.eventFnList = append(wc.eventFnList, eventfn)
}

func (wc *WindowController) Start() {

	// layout the windows
	for _, w := range wc.windows {
		w.Layout()
	}

	// event sender/receiver control channel and common event queue
	wc.ctrl = make(chan bool, 2)
	wc.queue = make(chan interface{}, 2)
	wc.waitgroup = new(sync.WaitGroup)

	// starting the event senders
	go wc.eventSenderSDL()
	go wc.eventSenderRenderTick(50) // 50 msec = 20frames/sec

	for _, efn := range wc.eventFnList {
		//go wc.eventSenderCustom(&wg, wc.queue, wc.ctrl)
		go efn(wc.waitgroup, wc.queue, wc.ctrl)
	}

	//starting the main event loop
	var event interface{}
	var channelreads int = 0
	var defaultreads int = 0
	var nilreads int = 0

queuereader:
	for event = range wc.queue {
		channelreads++
		//fmt.Printf("event: Type: %T \n", event)

		switch t := event.(type) {
		case nil:
			nilreads++
		case *EventRenderTick:
			// render windows that have changed
			for _, w := range wc.windows {
				if w.GetChanged() {
					fmt.Printf("EventRenderTick + Window changed. window:%d\n", w.id)
					w.RenderAll()
					w.Present()
				}
			}
		case *EventCustom:
			fmt.Printf("Custom event! channelreads:%d/s, defaultreads: %d/s nilreads: %d/s\n", channelreads, defaultreads, nilreads)
			channelreads = 0
			defaultreads = 0
			nilreads = 0
		case *sdl.QuitEvent:
			close(wc.ctrl)
			break queuereader
		case *sdl.WindowEvent:
			if t.Event == sdl.WINDOWEVENT_CLOSE {
				close(wc.ctrl)
				break queuereader
			} else {
				if found, rw := wc.getRootWindowById(t.WindowID); found {
					rw.handleWindowEvent(t)
				}
			}
		case *sdl.MouseMotionEvent:
			//fmt.Printf("[%d ms] MouseMotion\ttype:%d\twhich:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\twindow:%d\tstate:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel, t.WindowID, t.State)
			if found, rw := wc.getRootWindowById(t.WindowID); found {
				rw.handleMouseMotionEvent(t)
			}
		case *sdl.MouseButtonEvent:
			//fmt.Printf("[%d ms] MouseButton\ttype:%d\twhich:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\twindow:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State, t.WindowID)
			if found, rw := wc.getRootWindowById(t.WindowID); found {
				rw.handleMouseButtonEvent(t)
			} else if t.WindowID == 0 {
				// button was released outside a window
				// notify last known window
				if wc.lastWindow != nil {
					wc.lastWindow.handleMouseButtonEvent(t)
				}
			}
		case *sdl.MouseWheelEvent:
			fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y)
		case *sdl.KeyboardEvent:
			if found, rw := wc.getRootWindowById(t.WindowID); found {
				rw.handleKeyboardEvent(t)
			}
			fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n", t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
		case *sdl.TextInputEvent:
			//t.Text : the null-terminated input text in UTF-8 encoding
			//fmt.Printf("[%d ms] TextInput\ttype:%d\ttext:%s\n", t.Timestamp, t.Type, string(t.Text[:]))
			if found, rw := wc.getRootWindowById(t.WindowID); found {
				rw.handleTextInputEvent(t)
			}
		case *sdl.SysWMEvent:
			fmt.Printf("[%d ms] SysWMEvent\ttype:%d\n", t.Timestamp, t.Type)
		case *sdl.TouchFingerEvent:
		case *sdl.UserEvent:
			fmt.Printf("[%d ms] UserEvent\ttype:%d\n", t.Timestamp, t.Type)
		default:
			defaultreads++
			fmt.Println("default event: Type:", reflect.TypeOf(event).String())
		}
	}

	fmt.Println("waiting for wait group to finish")
	wc.waitgroup.Wait()
	fmt.Println("wait group is done")

	// Clean up the sdl resources befor exit
	for _, w := range wc.windows {
		w.Destroy()
	}

	sdl.Quit()
}
