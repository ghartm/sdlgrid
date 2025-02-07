package sdlgrid

import "github.com/veandco/go-sdl2/sdl"

// -----------------------------------------------------------------------------
// Display a text
type ItemText struct {
	ItemBase
	text         string // Buttons text
	textTexture  *sdl.Texture
	textureColor sdl.Color
	useBaseColor bool
}

func NewItemText(win *RootWindow) *ItemText {
	i := new(ItemText)
	i.o = Item(i)
	i.setRootWindow(win)
	i.SetSpec(LS_POS_PCT, LS_POS_PCT, LS_SIZE_COLLAPSE, LS_SIZE_COLLAPSE, 0, 0, 0, 0)
	return i
}
func (i *ItemText) UseBaseColor(b bool) { i.useBaseColor = b }

func (i *ItemText) oGetMinSize() (int32, int32) {
	i.prepareTexture()
	return i.GetSize()
}

func (i *ItemText) oRender() {

	//fmt.Printf("%s: ItemText.oRender()\n", i.GetName())
	r := i.GetRenderer()

	//utilRenderSolidBorder(r, &i.iframe, i.Style.ColorPurple)

	if i.textTexture != nil {
		if i.textureColor.Uint32() != i.GetColorScheme().text.Uint32() {
			i.textTexture.Destroy()
			i.textTexture = nil
			i.prepareTexture()
		}
		r.Copy(i.textTexture, nil, &i.iframe)
	}

}

func (i *ItemText) SetText(s string) {
	i.text = s
	if i.textTexture != nil {
		i.textTexture.Destroy()
		i.textTexture = nil
	}
	i.prepareTexture()
}

// prepares and renders the Texture of the Text
// Textures need to be prepared for layout.
func (i *ItemText) prepareTexture() {
	//fmt.Printf("ItemText.prepareTexture()\n")
	if i.textTexture == nil {
		if i.text != "" {
			if i.useBaseColor {
				i.textureColor = i.GetColorScheme().base
			} else {
				i.textureColor = i.GetColorScheme().text
			}
			sfc, err := i.Style.Font.RenderUTF8Blended(i.text, i.textureColor)
			//surface, _ := f.RenderUTF8Shaded(t, *c)
			if err != nil {
				panic(err)
			}
			i.textTexture, _ = i.GetRenderer().CreateTextureFromSurface(sfc)
			i.SetSpecSize(LS_POS_ABS, LS_POS_ABS, sfc.W, sfc.H)
			// layout
			i.iframe.W = sfc.W
			i.iframe.H = sfc.H
			sfc.Free()
		} else {
			i.SetSpecSize(LS_POS_ABS, LS_POS_ABS, 0, 0)
			//layout
			i.iframe.W = 0
			i.iframe.H = 0

		}
	}
}
