package main

import (
	"fmt"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth = 800
const winHeight = 600

type color struct {
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
}

type pos struct {
	x, y int
}

type ball struct {
	pos    // composition, something like inheritance
	radius int
	// xvel   float32
	// yvel   float32
	// color color
}

// this parametric angle way to find dots in a circle is computationally expensive
// func (ball *ball) draw(pixels []byte) {
// 	for r := 0; r <= ball.radius; r++ {
// 		for angle := 0; angle < 360; angle++ {
// 			angleRad := float64(angle) * (math.Pi / 180)

// 			curX := ball.x + int((float64(r) * math.Cos(angleRad)))
// 			curY := ball.y + int((float64(r) * math.Sin(angleRad)))

// 			setPixel(curX, curY, color{255, 255, 255}, pixels)
// 		}
// 	}
// }

// this is simple, bound check way. efficient because normal calculations
func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y <= ball.radius; y++ {
		for x := -ball.radius; x <= ball.radius; x++ {
			if x*x+y*y <= ball.radius*ball.radius {
				setPixel(ball.x+x, ball.y+y, color{255, 255, 255}, pixels)
			}
		}
	}
}

type paddle struct {
	pos
	w     int
	h     int
	color color
}

func (paddle *paddle) draw(pixels []byte) {
	startX := paddle.x - paddle.w/2
	startY := paddle.y - paddle.h/2

	// y is outer, x is inner
	// because the 2d matrix is being stored as array
	// this makes iteration linear
	for y := 0; y < paddle.h; y++ {
		for x := 0; x < paddle.w; x++ {
			setPixel(startX+x, startY+y, color{255, 255, 255}, pixels)
		}
	}
}

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

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

	pixels := make([]byte, winWidth*winHeight*4) // buffer = width * height * 4 as each pixel has 4 values R G B A

	// p := &paddle{}

	// p.x = winWidth / 2
	// p.y = winHeight / 2
	// p.w = 20
	// p.h = 100
	// p.color = color{}

	// p.draw(pixels)

	b := &ball{}
	b.x = winWidth / 2
	b.y = winHeight / 2
	b.radius = 200

	b.draw(pixels)

	tex.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*4)
	renderer.Copy(tex, nil, nil)
	renderer.Present()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		sdl.Delay(16)
	}
}
