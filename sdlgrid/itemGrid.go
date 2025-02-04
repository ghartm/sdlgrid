package sdlgrid

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// --------------------------------------------------------------------
// Container item that manages/arranges its subitems in a grid. It has no decoration
type ItemGrid struct {
	ItemBase
	cols                int
	rows                int
	spacing             int32
	grid                [][]Item // addressed by (column;row)
	rowSpec             []layoutParam
	colSpec             []layoutParam
	lastFoundSubItemCol int
	lastFoundSubItemRow int
	buttonCallback      func(int32, int32, uint8, uint8) //x, y int32, button, state uint8
}

func NewItemGrid(win *RootWindow, cols int, rows int) *ItemGrid {
	i := new(ItemGrid)
	i.o = Item(i)
	i.setRootWindow(win)
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_PCT, LS_SIZE_PCT, 0, 0, 100000, 100000)

	i.cols = cols
	i.rows = rows

	i.rowSpec = make([]layoutParam, i.rows)
	i.colSpec = make([]layoutParam, i.cols)
	i.grid = make([][]Item, i.cols)

	var c, r int

	for c = 0; c < i.cols; c++ {
		i.grid[c] = make([]Item, i.rows)
	}

	// default is to collapse all but the last
	if i.cols > 0 {
		for c = 0; c < i.cols-1; c++ {
			i.colSpec[c] = layoutParam{LS_SIZE_COLLAPSE, 0}
		}
		i.colSpec[c] = layoutParam{LS_SIZE_PCT, 100000}
	}
	if i.rows > 0 {
		for r = 0; r < i.rows-1; r++ {
			i.rowSpec[r] = layoutParam{LS_SIZE_COLLAPSE, 0}

		}
		i.rowSpec[r] = layoutParam{LS_SIZE_PCT, 100000}

	}

	return i
}

func (i *ItemGrid) SetButtonCallback(cb func(int32, int32, uint8, uint8)) {
	i.buttonCallback = cb
}
func (i *ItemGrid) oNotifyMouseButton(x, y int32, button uint8, state uint8) {
	fmt.Printf("%s: ItemGrid.oNotifyMouseButton()\n", i.GetName())
	if i.buttonCallback != nil {
		i.buttonCallback(x, y, button, state)
	}
}

// extend grid by one column. return index of new column
func (i *ItemGrid) AppendColumn() int {
	i.grid = append(i.grid, make([]Item, i.rows))

	i.colSpec = append(i.colSpec, layoutParam{LS_SIZE_PCT, 100000})
	i.cols = len(i.colSpec)
	if i.cols > 1 {
		i.colSpec[i.cols-1] = layoutParam{LS_SIZE_COLLAPSE, 0}
	}

	// if it was the first column add a row as well
	if i.cols == 1 {
		i.grid[i.cols-1] = append(i.grid[i.cols-1], nil)
		i.rowSpec = append(i.rowSpec, layoutParam{LS_SIZE_PCT, 100000})
		i.rows = len(i.rowSpec)
	}

	return i.cols - 1
}

// extend grid by one row. return index of new row
func (i *ItemGrid) AppendRow() int {

	// if there is no column jet, add a column first
	if i.cols == 0 {
		i.AppendColumn()
		//if it was the first column a row was added as well
	} else {

		for n := range i.grid {
			i.grid[n] = append(i.grid[n], nil)
		}
		i.rowSpec = append(i.rowSpec, layoutParam{LS_SIZE_PCT, 100000})
	}

	i.rows = len(i.rowSpec)
	if i.rows > 1 {
		i.rowSpec[i.rows-1] = layoutParam{LS_SIZE_COLLAPSE, 0}
	}

	return i.rows - 1
}

// set the inner spacing between the cells
func (i *ItemGrid) SetSpacing(spc int32) {
	i.spacing = spc
	i.useTmpMinSize = false
}

// set the horizontal layout spec for a column
func (i *ItemGrid) SetColSpec(col int, p layoutParam) {
	i.colSpec[col] = p
	i.useTmpMinSize = false
}

// set the vertical layout spec for a row
func (i *ItemGrid) SetRowSpec(row int, p layoutParam) {
	i.rowSpec[row] = p
	i.useTmpMinSize = false
}

// places a sub-item into a cell
func (i *ItemGrid) SetSubItem(col int, row int, si Item) bool {
	if col < i.cols && row < i.rows {
		i.grid[col][row] = si
		si.SetParent(i)
		i.useTmpMinSize = false
		return true
	}
	return false
}

func (i *ItemGrid) oGetSubFrame() *sdl.Rect { return &i.iframe }

// get minimum sizes of items in cells ( column-width; row-heigt)
func (i *ItemGrid) getMinSizeFields() (cw, rh []int32) {

	cw = make([]int32, i.cols)
	rh = make([]int32, i.rows)
	for c := 0; c < i.cols; c++ {
		for r := 0; r < i.rows; r++ {
			if i.grid[c][r] != nil {
				x, y, mw, mh := i.grid[c][r].GetCollapsedSpec()
				if i.grid[c][r].IsAutoSize() {
					//if it is not a fixed size - Minimum size must be computed
					mw, mh = i.grid[c][r].oGetMinSize()
				}

				mw += x
				mh += y

				// if absolute size
				if i.colSpec[c].S == LS_SIZE_ABS {
					mw = i.colSpec[c].V
				}

				if i.rowSpec[r].S == LS_SIZE_ABS {
					mh = i.rowSpec[r].V
				}

				if cw[c] < mw {
					cw[c] = mw
				}
				if rh[r] < mh {
					rh[r] = mh
				}
			}
		}
	}

	return cw, rh

}
func (i *ItemGrid) oNotifyPostLayout(sizeChanged bool) {
	//fmt.Printf("%s: ItemGrid.oNotifyPostLayout()\n", i.GetName())
	// compute frame-size for every cell and layout the subitem in each frame
	// get required minimum size from every cell. cw:column-width rh:row-hight min for row and col
	//column-width and row-height

	cw, rh := i.getMinSizeFields()

	//TODO  make grid layout aware to sizeChanged

	// get minimum needed space for grid Frame
	// and count the number of expand-specs if it has content

	var sum int32
	var npct int32
	var basepct int32

	// distribute free space between required size and size of grid-item frame across the expandable cells according to the percentages
	for n := 0; n < i.cols; n++ {
		sum += cw[n]
		if i.colSpec[n].S == LS_SIZE_PCT {
			npct++
			basepct += i.colSpec[n].V
		}
	}
	// if there is space for distribution and percentage layout exists
	if free := (i.iframe.W - (sum + (i.spacing * int32(i.cols-1)))); free > 0 && npct > 0 {
		for n := 0; n < i.cols; n++ {
			if i.colSpec[n].S == LS_SIZE_PCT {
				// normalize percentages
				cw[n] += utilPct(free, utilNormPct(basepct, i.colSpec[n].V))
			}
		}
	}
	//-----------

	npct = 0
	basepct = 0
	sum = 0
	for n := 0; n < i.rows; n++ {
		sum += rh[n]
		if i.rowSpec[n].S == LS_SIZE_PCT {
			npct++
			basepct += i.rowSpec[n].V
		}
	}

	// if there is space for distribution and percentage layout exists
	if free := (i.iframe.H - (sum + (i.spacing * int32(i.rows-1)))); free > 0 && npct > 0 {
		for n := 0; n < i.rows; n++ {
			if i.rowSpec[n].S == LS_SIZE_PCT {
				// normalize percentages
				rh[n] += utilPct(free, utilNormPct(basepct, i.rowSpec[n].V))
			}
		}
	}

	// layout subitems in the cells according to cw and rh
	var pf sdl.Rect //parent frame of sub-item
	pf.X = i.iframe.X
	for c := 0; c < i.cols; c++ {
		pf.Y = i.iframe.Y
		pf.W = cw[c]
		for r := 0; r < i.rows; r++ {
			pf.H = rh[r]
			if i.grid[c][r] != nil {
				i.grid[c][r].Layout(&pf, sizeChanged)
			}
			pf.Y += rh[r] + i.spacing
		}
		pf.X += cw[c] + i.spacing
	}
}

func (i *ItemGrid) oReportSubitems(lvl int) {

	for c := 0; c < i.cols; c++ {
		for r := 0; r < i.rows; r++ {
			for n := 0; n < lvl; n++ {
				fmt.Print("--")
			}

			if i.grid[c][r] != nil {
				fmt.Printf("%s: c:%d r:%d\n", i.GetName(), c, r)
				i.grid[c][r].Report(lvl)
			} else {
				fmt.Printf("%s: c:%d r:%d empty\n", i.GetName(), c, r)
			}
		}
	}
}
func (i *ItemGrid) oWithItems(fn func(Item)) {
	fmt.Printf("%s: ItemBase.oWithItems() no subitems\n", i.GetName())
	// handle the call for myself
	fn(i)
	// forward to all subitems
	for c := 0; c < i.cols; c++ {
		for r := 0; r < i.rows; r++ {
			if i.grid[c][r] != nil {
				i.grid[c][r].oWithItems(fn)
			}
		}
	}
}
func (i *ItemGrid) oGetMinSize() (int32, int32) {
	//fmt.Printf("%s: ItemGrid.oGetMinSize()\n", i.GetName())
	if !i.useTmpMinSize {
		//i.minw, i.minh

		//column-width and row-height
		if i.cols < 1 || i.rows < 1 {
			return 0, 0
		} else {

			cw, rh := i.getMinSizeFields()

			// get minimum needed space for grid Frame
			var cwsum, rhsum int32
			for c := 0; c < i.cols; c++ {
				cwsum += cw[c]
			}
			for r := 0; r < i.rows; r++ {
				rhsum += rh[r]
			}
			i.minw = cwsum + (i.spacing * int32(i.cols-1))
			i.minh = rhsum + (i.spacing * int32(i.rows-1))
			fmt.Printf("%s: ItemGrid.oGetMinSize() w:%d h:%d\n", i.GetName(), i.minw, i.minh)
		}
		i.useTmpMinSize = true
	}
	return i.minw, i.minh
}

func (i *ItemGrid) oRender() {
	//fmt.Printf("%s: ItemGrid.oRender()\n", i.GetName())
	// ItemGrid does not have a decoration.

	/*// debug frame
	rd := i.GetRenderer()
	s := i.GetStyle()
	rd.SetDrawBlendMode(sdl.BLENDMODE_NONE)
	utilRenderSolidBorder(rd, &i.iframe, s.colorRed)
	*/

	// forward to content
	for c := 0; c < i.cols; c++ {
		for r := 0; r < i.rows; r++ {
			if i.grid[c][r] != nil {
				i.grid[c][r].Render()
			}
		}
	}

}

func (i *ItemGrid) oFindSubItem(x, y int32, e sdl.Event) (found bool, item Item) {
	if i.cols < 1 || i.rows < 1 {
		return false, nil
	}
	// first check last hit
	if (i.lastFoundSubItemCol >= 0) && (i.lastFoundSubItemCol < i.cols) && (i.lastFoundSubItemRow < i.rows) {
		if si := i.grid[i.lastFoundSubItemCol][i.lastFoundSubItemRow]; si != nil {
			if si.CheckPos(x, y) {
				return true, si
			}
		}
	}

	// check sequential if not found
	for c := 0; c < i.cols; c++ {
		for r := 0; r < i.rows; r++ {
			if si := i.grid[c][r]; si != nil {
				if si.CheckPos(x, y) {
					i.lastFoundSubItemCol = c
					i.lastFoundSubItemRow = r
					return true, si
				}
			}
		}
	}
	return false, nil
}
