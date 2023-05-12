// Package svgJoin Copyright 2023 Gryaznov Nikita
// Licensed under the Apache License, Version 2.0
package svgJoin

type (
	direction rune
	align     rune
	xy        struct {
		x, y float64
	}
	Chunk struct {
		viewport, viewBox, position xy
		body                        string
	}
	Group struct {
		body       []part
		isVertical bool
		toForward  bool
		offset     float64
		align      uint
	}
	part interface {
		print(xy) (string, error)
		size() xy
	}
)
