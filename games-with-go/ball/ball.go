package main

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth  = 800
	winHeight = 600
	radius    = 160
	numVerts  = 32
	springK   = 0.1
	damping   = 0.85
)

type vertex struct {
	x, y   float64
	vx, vy float64
	angle  float64
	radius float64
}

type ball struct {
	x, y     int
	vertices []vertex
	dragged  int // index of dragged vertex (in terms of vertices), -1 if none
}

func newBall(x, y int, r int) *ball {
	verts := make([]vertex, numVerts)
	for i := 0; i < numVerts; i++ {
		angleRad := 2 * math.Pi * float64(i) / float64(numVerts)

		x := float64(x) + float64(r)*math.Cos(angleRad)
		y := float64(y) + float64(r)*math.Sin(angleRad)

		verts[i] = vertex{x, y, 0, 0, angleRad, float64(r)}
	}

	return &ball{x, y, verts, -1}
}

func (b *ball) update(mouseX, mouseY int32, dragging bool) {
	x := float64(b.x)
	y := float64(b.y)

	for i := 0; i < numVerts; i++ {
		v := &b.vertices[i]

		// follow when dragged
		if dragging && i == b.dragged {
			v.x = float64(mouseX)
			v.y = float64(mouseY)
			v.vx = 0
			v.vy = 0
			continue
		}

		// apply spring

		// calculate original pos
		restX := x + v.radius*math.Cos(v.angle)
		restY := y + v.radius*math.Sin(v.angle)

		// hooke's law
		forceX := springK * (restX - v.x)
		forceY := springK * (restY - v.y)

		// Pull toward neighbors
		for nd := 1; nd < 8; nd++ {
			weight := springK / float64(8-nd)

			next := &b.vertices[(i+nd)%numVerts]
			prev := &b.vertices[(i-nd+numVerts)%numVerts]

			midX := (next.x + prev.x) / 2
			midY := (next.y + prev.y) / 2

			neighborFX := weight * (midX - v.x)
			neighborFY := weight * (midY - v.y)

			v.vx += neighborFX
			v.vy += neighborFY
		}

		v.vx += forceX
		v.vy += forceY

		// damping
		v.vx *= damping
		v.vy *= damping

		// update new pos
		v.x += v.vx
		v.y += v.vy
	}
}

func (b *ball) draw(r *sdl.Renderer) {
	for i := 0; i < numVerts; i++ {
		// this gives next vertex index in polygon
		// i.e. 0 -> 1, 1 -> 2, 2 -> 3, 3 -> 0
		j := (i + 1) % numVerts
		r.DrawLine(int32(b.vertices[i].x), int32(b.vertices[i].y),
			int32(b.vertices[j].x), int32(b.vertices[j].y))
	}
}

func (b *ball) tryStartDrag(mouseX, mouseY int32) bool {
	// logic of minDist is here so we don't drag two vertices together
	minDist := math.MaxFloat64

	// none of vertices are being dragged
	idx := -1

	for k, v := range b.vertices {
		distMouseVertexX := float64(mouseX) - v.x
		distMouseVertexY := float64(mouseY) - v.y
		dist := distMouseVertexX*distMouseVertexX + distMouseVertexY*distMouseVertexY

		if dist < 100 && dist < minDist { // 10*10, i..e 10 radius
			minDist = dist
			idx = k // index of vertex
		}
	}

	if idx != -1 {
		b.dragged = idx
		return true
	}
	return false
}

func (b *ball) stopDrag() {
	b.dragged = -1
}

func main() {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Bawls", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)
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

	ball := newBall(winWidth/2, winHeight/2, radius)
	var mouseX, mouseY int32
	var dragging bool

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			// type switch concept bhul gya tha
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.MouseButtonEvent:
				mouseX, mouseY = e.X, e.Y
				if e.Type == sdl.MOUSEBUTTONDOWN && e.Button == sdl.BUTTON_LEFT {
					dragging = ball.tryStartDrag(mouseX, mouseY)
				}
				if e.Type == sdl.MOUSEBUTTONUP && e.Button == sdl.BUTTON_LEFT {
					dragging = false
					ball.stopDrag()
				}
			case *sdl.MouseMotionEvent:
				mouseX, mouseY = e.X, e.Y
			}
		}

		ball.update(mouseX, mouseY, dragging)

		// because sdl2 needs a color to clear
		renderer.SetDrawColor(20, 20, 20, 255)
		renderer.Clear()

		renderer.SetDrawColor(200, 255, 200, 255)
		ball.draw(renderer)

		renderer.Present()

		sdl.Delay(16)
	}
}
