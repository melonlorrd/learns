// main.go
package main

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth  = 800
	winHeight = 800
)

type color struct {
	r, g, b, a byte
}

func setPixel(x, y int, c color, pixels []byte) {
	if x < 0 || x >= winWidth || y < 0 || y >= winHeight {
		return
	}

	index := (y*winWidth + x) * 4

	if (index + 3) > len(pixels) {
		return
	}

	pixels[index] = c.r
	pixels[index+1] = c.g
	pixels[index+2] = c.b
	pixels[index+3] = c.a
}

func main() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatalf("could not initialize sdl: %v", err)
	}
	defer sdl.Quit()

	win, err := sdl.CreateWindow(
		"Gradient",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		winWidth,
		winHeight,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		log.Fatalf("could not create window: %v", err)
	}
	defer win.Destroy()

	renderer, err := sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("could not create renderer: %v", err)
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	for y := range winHeight {
		for x := range winWidth {
			r := uint8(float64(y) / float64(winWidth) * 256)
			g := uint8(float64(x) / float64(winHeight) * 256)
			b := uint8(0)

			setPixel(x, y, color{r, g, b, 255}, pixels)
		}
	}

	tex.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*4)

	renderer.Copy(tex, nil, nil)
	renderer.Present()

	running := true
	for running {
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch e := ev.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN && e.Keysym.Sym == sdl.K_ESCAPE {
					running = false
				}
			}
		}
		sdl.Delay(16)
	}
}
