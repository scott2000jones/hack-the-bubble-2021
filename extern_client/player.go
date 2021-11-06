package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"net"
)

type player struct {
	tex *sdl.Texture
	x, y float64
}

const (
	playerHeight = 162
	playerWidth = 100
)

func newPlayer(renderer *sdl.Renderer, path string, startx int, starty int) (p player, err error) {
	img, err := sdl.LoadBMP(path)
	if err != nil {
		return player{}, fmt.Errorf("loading player sprite: %v", err)
	}
	defer img.Free()
	p.tex, err = renderer.CreateTextureFromSurface(img)
	if err != nil {
		return player{}, fmt.Errorf("creating player texture: %v", err)
	}

	p.x = float64(startx)
	p.y = float64(starty)

	return p, nil
}

func (p *player) draw(renderer *sdl.Renderer) {

	x := p.x - playerWidth/2.0
	y := p.y - playerHeight/2.0

	renderer.Copy(p.tex,
		&sdl.Rect{X: 0, Y: 0, W: playerWidth, H: playerHeight},
		&sdl.Rect{X: int32(x), Y: int32(y), W: playerWidth, H: playerHeight})
}

func (p *player) update(conn net.Conn) {
	keys := sdl.GetKeyboardState()

	if keys[sdl.SCANCODE_LEFT] == 1 {
		// send LEFT to UDP
		go SendUDP("left", conn)
	} else if keys[sdl.SCANCODE_RIGHT] == 1 {
		// send RIGHT to UDP
		go SendUDP("right", conn)
	} else if keys[sdl.SCANCODE_UP] == 1 {
		// send UP to UDP
		go SendUDP("up", conn)
	} else if keys[sdl.SCANCODE_DOWN] == 1 {
		// send DOWN to UDP
		go SendUDP("down", conn)
	}
}
