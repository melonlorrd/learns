package main

import (
	"log"
	"math"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	aspectRatio = 16.0 / 9.0
	winWidth    = 1600
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

func hitSphere(cx, cy, cz, radius float64, ox, oy, oz, dx, dy, dz float64) float64 {
	// c -> center, o -> ray origin, d -> direction

	ocx := ox - cx
	ocy := oy - cy
	ocz := oz - cz

	a := dx*dx + dy*dy + dz*dz
	b := 2.0 * (ocx*dx + ocy*dy + ocz*dz)
	c := ocx*ocx + ocy*ocy + ocz*ocz - radius*radius

	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return -1.0
	} else {
		return (-b - math.Sqrt(discriminant)) / (2.0 * a)
	}
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

	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * (float64(winWidth) / float64(winHeight))

	pixelDeltaUx := viewportWidth / float64(winWidth)
	pixelDeltaUy := 0.0
	pixelDeltaUz := 0.0

	pixelDeltaVx := 0.0
	pixelDeltaVy := -viewportHeight / float64(winHeight)
	pixelDeltaVz := 0.0

	viewportUpperLeftX := -viewportWidth / 2.0
	viewportUpperLeftY := viewportHeight / 2.0
	viewportUpperLeftZ := -focalLength

	pixel00X := viewportUpperLeftX + 0.5*(pixelDeltaUx+pixelDeltaVx)
	pixel00Y := viewportUpperLeftY + 0.5*(pixelDeltaUy+pixelDeltaVy)
	pixel00Z := viewportUpperLeftZ + 0.5*(pixelDeltaUz+pixelDeltaVz)

	for j := range winHeight {
		for i := range winWidth {
			// incremental pixel position
			px := pixel00X + float64(i)*pixelDeltaUx + float64(j)*pixelDeltaVx
			py := pixel00Y + float64(i)*pixelDeltaUy + float64(j)*pixelDeltaVy
			pz := pixel00Z + float64(i)*pixelDeltaUz + float64(j)*pixelDeltaVz

			// ray direction = pixel_center - camera_center (camera at origin => direction = pixel center)
			rx := px
			ry := py
			rz := pz

			// unit vector of direction
			len := math.Sqrt(rx*rx + ry*ry + rz*rz) // magnitude of ray direction vector

			// normalization so we get values between -1 and 1, so it becomes a unit vector, i.e. root(ux2 + uy2 + uz2) = 1
			// with normalization, we don't care about distance, only care about the direction

			ux := rx / len
			uy := ry / len
			uz := rz / len

			// sky gradient: a = 0.5*(unit_direction.y + 1.0), only using y because top to bottom
			a := 0.5 * (uy + 1.0)

			// linear rgb (doing linear interpolation)
			rr := (1.0-a)*1.0 + a*0.5 // white -> light blue.r
			gg := (1.0-a)*1.0 + a*0.7 // white -> light blue.g
			bb := (1.0-a)*1.0 + a*1.0 // white -> light blue.b

			sphereCenterX := 0.0
			sphereCenterY := 0.0
			sphereCenterZ := -1.0

			// ray's origin is camera, so 0.0, 0.0, 0.0
			ox := 0.0
			oy := 0.0
			oz := 0.0

			t := hitSphere(sphereCenterX, sphereCenterY, sphereCenterZ, 0.5, ox, oy, oz, ux, uy, uz)

			if t > 0.0 {
				// Compute hit point: P = O + t*D, parametric equation of ray
				px := ox + t*ux
				py := oy + t*uy
				pz := oz + t*uz

				// Normal vector: N = unit_vector(P - center)
				nx := px - sphereCenterX
				ny := py - sphereCenterY
				nz := pz - sphereCenterZ
				nlen := math.Sqrt(nx*nx + ny*ny + nz*nz)
				nx /= nlen
				ny /= nlen
				nz /= nlen

				// Map from [-1,1] to [0,1]
				r := 0.5 * (nx + 1.0)
				g := 0.5 * (ny + 1.0)
				b := 0.5 * (nz + 1.0)

				setPixel(i, j, color{
					byte(255 * r),
					byte(255 * g),
					byte(255 * b),
				}, pixels)
				continue
			}

			c := color{byte(int(255 * rr)), byte(int(255 * gg)), byte(int(255 * bb))} // rr, gg, bb is between 0 and 1

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
