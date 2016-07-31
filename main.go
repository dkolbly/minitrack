// +build darwin linux windows

// An app that allows a kitty to capture territory while being harassed
// by lines
//
package main

import (
	"log"
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	//"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	//"golang.org/x/mobile/exp/sprite"
	//"golang.org/x/mobile/exp/sprite/clock"
	//"golang.org/x/mobile/exp/sprite/glsprite"
	"golang.org/x/mobile/gl"
	"github.com/golang/freetype/truetype"
)

type AppContext struct {
	images *glutil.Images
	startTime time.Time
	glctx       gl.Context
	currentSize size.Event
	appFont *truetype.Font
}

func (ac *AppContext) start(glctx gl.Context) {
	ac.glctx = glctx

	if ac.startTime.IsZero() {
		log.Printf("first start")
		// first start
		ac.startTime = time.Now()
	} else {
		log.Printf("subsequent start")
	}

	ac.images = glutil.NewImages(glctx)
	
	ac.initFonts()
}

func (ac *AppContext) stop() {
	ac.glctx = nil
}

func (ac *AppContext) resize(e size.Event) {
	ac.currentSize = e
}

func (ac *AppContext) paint(e paint.Event) {
	glctx := ac.glctx
	if glctx == nil {
		// we are not yet initialized with a context to paint with
		return
	}
	// if we were painting as fast as we could, we could ignore
	// external paint events
	//if e.External { return }

	//onPaint(glctx, currentSize)

	glctx.ClearColor(153.0/255.0, 204.0/255.0, 1, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	
	ac.renderText()
}

func (ac *AppContext) run(a app.App) {

	for e := range a.Events() {
		switch e := a.Filter(e).(type) {
		case lifecycle.Event:
			switch e.Crosses(lifecycle.StageVisible) {
			case lifecycle.CrossOn:
				glctx, _ := e.DrawContext.(gl.Context)
				ac.start(glctx)
				a.Send(paint.Event{})
			case lifecycle.CrossOff:
				ac.stop()
			}
		case size.Event:
			ac.resize(e)
			// TODO: Screen reorientation?
			//touchX = float32(sz.WidthPx / 2)
			//touchY = float32(sz.HeightPx / 2)
		case paint.Event:
			ac.paint(e)
			a.Publish()
			// Drive the animation by preparing to
			// paint the next frame after this one
			// is shown.
			//a.Send(paint.Event{})
		case touch.Event:
			/*
			gg.lastTouch = e
			gg.touch.X = int(e.X)
			gg.touch.Y = int(e.Y)
*/
			log.Printf("Touch (%.1f,%.1f)", e.X, e.Y)
		}
	}
}

func main() {
	ac := &AppContext{}
	app.Main(ac.run)
}

