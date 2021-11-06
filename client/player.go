package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

type player struct {
	tex *sdl.Texture
}

func newPlayer(renderer *sdl.Renderer, path string) (p player, err error) {
	img, err := sdl.LoadBMP(path)
	if err != nil {
		return player{}, fmt.Errorf("loading player sprite: %v", err)
	}
	defer img.Free()
	p.tex, err = renderer.CreateTextureFromSurface(img)
	if err != nil {
		return player{}, fmt.Errorf("creating player texture: %v", err)
	}

	return p, nil
}

func (p *player) draw(renderer *sdl.Renderer) {
	renderer.Copy(p.tex,
		&sdl.Rect{X: 0, Y: 0, W: 100, H: 162},
		&sdl.Rect{X: 40, Y: 20, W: 105, H: 120})
}
