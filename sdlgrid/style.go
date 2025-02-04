package sdlgrid

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Style struct {
	Font *ttf.Font
	name string

	cursorBlinkRate int

	decoUnit         int32
	spacing          int32
	fontSize         int32
	averageFontHight int32
	averageRuneWidth int32

	colorRed    *sdl.Color
	colorGreen  *sdl.Color
	colorBlue   *sdl.Color
	colorPurple *sdl.Color
	colorBlack  *sdl.Color
	colorWhite  *sdl.Color

	//colorBg   *sdl.Color
	csDefault *ColorScheme
	csWindow  *ColorScheme
}

func NewStyle(n string) *Style {
	s := new(Style)

	if !ttf.WasInit() {
		ttf.Init()
	}

	s.InitDefault()
	s.name = n
	return s
}

func (s *Style) InitDefault() {
	s.cursorBlinkRate = 300
	//red := &sdl.Color{255, 0, 0, 255}
	// R, G, B, A
	s.colorRed = &sdl.Color{R: 255, G: 50, B: 50, A: 255}
	s.colorGreen = &sdl.Color{R: 5, G: 200, B: 5, A: 255}
	s.colorBlue = &sdl.Color{R: 50, G: 50, B: 255, A: 255}
	s.colorPurple = &sdl.Color{R: 255, G: 50, B: 255, A: 255}
	s.colorBlack = &sdl.Color{R: 0, G: 0, B: 0, A: 255}
	s.colorWhite = &sdl.Color{R: 255, G: 255, B: 255, A: 255}

	//standard item
	s.csDefault = new(ColorScheme)
	//s.csDefault.SetBaseColor(&sdl.Color{70, 100, 100, 255}, 50)
	s.csDefault.SetBaseColor(&sdl.Color{R: 190, G: 190, B: 190, A: 255})

	// Windows colors
	s.csWindow = new(ColorScheme)
	s.csWindow.SetBaseColor(&sdl.Color{R: 150, G: 200, B: 150, A: 255})

	s.decoUnit = 16
	// spacing shall not drop below 2 !
	s.spacing = 2
	s.fontSize = 13

	var err error
	// find /usr/share/fonts/truetype -name '*.ttf'
	// /usr/share/fonts/truetype/dejavu/DejaVuSans.ttf
	// /usr/share/fonts/truetype/ubuntu/Ubuntu-RI.ttf

	//fmt.Printf("Available Fonts: in /usr/share/fonts/truetype/\n")
	//files, err := filepath.Glob("/usr/share/fonts/truetype/*/*.ttf")
	//for _, v := range files {
	//	fmt.Printf("%s\n", v)
	//}

	//s.Font, err = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", int(s.fontSize))
	//s.Font, err = ttf.OpenFont("/usr/share/fonts/truetype/ubuntu/Ubuntu-RI.ttf", int(s.fontSize))
	//s.Font, err = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", int(s.fontSize))
	s.Font, err = ttf.OpenFont("/usr/share/fonts/truetype/ubuntu/Ubuntu-R.ttf", int(s.fontSize))
	if err != nil {
		panic(err)
	}
	s.Font.SetHinting(ttf.HINTING_NORMAL)
	//s.Font.SetKerning(false)

	s.averageFontHight = s.GetTextHight("QqJjPpGg")
	s.averageRuneWidth = s.GetAverageRuneWidth("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQTSTUVWXYZ")
}

func (s *Style) Destroy() {
	s.Font.Close()
}
func (s *Style) GetAverageRuneWidth(txt string) int32 {
	var sum, n int
	for _, r := range txt {
		x, _, _ := s.Font.SizeUTF8(string(r))
		n++
		sum += x
	}
	return int32(sum / n)

}

func (s *Style) GetTextLen(txt string) int32 {
	x, _, _ := s.Font.SizeUTF8(txt)
	return int32(x)
}

func (s *Style) GetTextHight(txt string) int32 {
	_, x, _ := s.Font.SizeUTF8(txt)
	return int32(x)
}
