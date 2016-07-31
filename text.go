package main

import (
	//"log"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/size"

	//"golang.org/x/mobile/exp/gl/glutil"
	//"golang.org/x/mobile/gl"
	"golang.org/x/mobile/geom"
)

func F(x float32) fixed.Int26_6 {
	return fixed.Int26_6((1<<6) * x)
}

func (ac *AppContext) initFonts() {
	f, err := asset.Open("basic.ttf")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	raw, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	font, err := freetype.ParseFont(raw)
	if err != nil {
		panic(err)
	}
	ac.appFont = font
}

func (ac *AppContext) textArea() image.Rectangle {
	gap := ac.currentSize.WidthPx / 24
	if ac.currentSize.Orientation == size.OrientationLandscape {
		gap = ac.currentSize.HeightPx / 24
	}
	//log.Printf("gap=%d", gap)
	androidTopBar := 64
	return image.Rectangle{
		Min: image.Point{
			X: gap,
			Y: gap + androidTopBar,
		},
		Max: image.Point{
			X: ac.currentSize.WidthPx - gap,
			Y: gap + 6*gap + androidTopBar,
		},
	}
}

func (ac *AppContext) renderText() {

	sz := ac.currentSize
	images := ac.images
	rect := ac.textArea()

	textArea := images.NewImage(rect.Dx(), rect.Dy())
	draw.Draw(textArea.RGBA, textArea.RGBA.Bounds(),
		image.NewUniform(color.RGBA{0xf5, 0x73, 0x99, 0xff}),
		image.ZP,
		draw.Src)

	fg := image.NewUniform(color.RGBA{0, 0, 0, 0xff})
	face := truetype.NewFace(
		ac.appFont,
		&truetype.Options{
			Size:    12,
			DPI:     72 * float64(sz.PixelsPerPt),
			Hinting: font.HintingNone,
		})
	d := &font.Drawer{
		Dst:  textArea.RGBA,
		Src:  fg,
		Face: face,
		Dot: fixed.Point26_6{
			X: F(10 * sz.PixelsPerPt),
			Y: F(14 * sz.PixelsPerPt),
		},
	}
	msg := bytes.Buffer{}
	fmt.Fprintf(&msg, "Hi: %d x %d ",
		ac.currentSize.WidthPx, ac.currentSize.HeightPx)
	switch ac.currentSize.Orientation {
	case size.OrientationUnknown:
		msg.Write([]byte("unk."))
	case size.OrientationPortrait:
		msg.Write([]byte("port."))
	case size.OrientationLandscape:
		msg.Write([]byte("land."))
	}
	fmt.Fprintf(&msg, " %.1f px/pt", sz.PixelsPerPt)
	d.DrawString(msg.String())

	if true {
		// do some color coding
		d := textArea.RGBA
		for y := 0; y < 100; y++ {
			c := color.RGBA{0, 0, 0, 0xff}
			switch y / 10 {
			case 0:
				c.R = 0xff
			case 1:
				c.G = 0xff
			case 2:
				c.B = 0xff
			case 3:
				c.R = 0xff
				c.G = 0xff
			case 4:
				c.G = 0xff
				c.B = 0xff
			case 5:
				c.R = 0xff
				c.B = 0xff
			}
			d.Set(y%10, y, c)
		}
	}

	textArea.Upload()
	tl := geom.Point{
		X: geom.Pt(float32(rect.Min.X) / sz.PixelsPerPt),
		Y: geom.Pt(float32(rect.Min.Y) / sz.PixelsPerPt),
	}
	tr := geom.Point{
		X: geom.Pt(float32(rect.Max.X) / sz.PixelsPerPt),
		Y: geom.Pt(float32(rect.Min.Y) / sz.PixelsPerPt),
	}
	bl := geom.Point{
		X: geom.Pt(float32(rect.Min.X) / sz.PixelsPerPt),
		Y: geom.Pt(float32(rect.Max.Y) / sz.PixelsPerPt),
	}

	textArea.Draw(
		sz,
		tl,
		tr,
		bl,
		textArea.RGBA.Bounds()) //rect) //sz.Bounds())
	textArea.Release()
}
