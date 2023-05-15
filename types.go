// Package svgJoin Copyright 2023 Gryaznov Nikita
// Licensed under the Apache License, Version 2.0
package svgJoin

type (
	xy struct {
		x, y float64
	}
	Chunk struct {
		viewport, viewBox, position xy
		body                        string
	}
	Group struct {
		body       []Part
		isVertical bool
		toForward  bool
		offset     float64
		align      uint8
	}
	Part interface {
		print(xy, *clipId) (string, error)
		size() xy
	}
)
