package main

import (
	"log"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	aspectRatio = 16.0 / 9.0
	winWidth    = 800
	winHeight   = int(float64(winWidth) / aspectRatio)
)

type color struct {
	r, g, b byte
}

func setPixel(x, y int, c color, pixels []byte) {
	if x < 0 || x >= winWidth || y < 0 || y >= winHeight {
		return
	}
	idx := (y*winWidth + x) * 3
	if idx+2 >= len(pixels) {
		return
	}
	pixels[idx+0] = c.r
	pixels[idx+1] = c.g
	pixels[idx+2] = c.b
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
		int32(winWidth),
		int32(winHeight),
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
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		log.Fatalf("could not create texture: %v", err)
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*3)

	for j := range winHeight {
		for i := range winWidth {
			c := color{255, 0, 0}
			setPixel(i, j, c, pixels)
		}
	}

	if err := tex.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*3); err != nil {
		log.Fatalf("texture update failed: %v", err)
	}
	if err := renderer.Copy(tex, nil, nil); err != nil {
		log.Fatalf("renderer copy failed: %v", err)
	}
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
