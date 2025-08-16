package main

import (
	"fmt"
//	"math/rand"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth = 800
const winHeight = 600

type color struct {
//	r, g, b, a byte
	r, g, b byte
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
//	pixels[index+3] = c.a
}

func main() {
	window, err := sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*3) // buffer = width * height * 4 as each pixel has 4 values R G B A

	for y := range winHeight {
		for x := range winWidth {
			setPixel(x, y, color{255, 0, 0}, pixels)
			//setPixel(x, y, color{byte(rand.Intn((256))), byte(rand.Intn((256))), byte(rand.Intn((256))), byte(rand.Intn((256)))}, pixels)
		}
	}

	tex.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*3)
	renderer.Copy(tex, nil, nil)
	renderer.Present()
	sdl.Delay(10000)
}
