// Package svgJoin Copyright 2023 Gryaznov Nikita
// Licensed under the Apache License, Version 2.0
package svgJoin

import "errors"

type (
	direction rune
	align     uint8
)

const (
	Rightward direction = 'R'
	Leftward  direction = 'L'
	Upward    direction = 'U'
	Downward  direction = 'D'
	Top       align     = 0
	Left      align     = 0
	Middle    align     = 1
	Bottom    align     = 2
	Right     align     = 2
)

// Join returns images sequentially connected in the selected direction with selected offset aligned along the selected edge.
// Directions: svgJoin.Rightward, svgJoin.Leftward, svgJoin.Upward, svgJoin.Downward
// Alignment: svgJoin.Top, svgJoin.Left, svgJoin.Middle, svgJoin.Bottom, svgJoin.Right
// Part could be Group or Chunk
func Join(direction direction, align align, offset float64, parts ...Part) (res Group, err error) {
	// Check align input
	if align > 2 {
		err = errors.New("invalid align code")
		return
	} else {
		res.align = uint8(align)
	}
	// Check direction input
	switch direction {
	case 'r', 'R':
		res.toForward = true
	case 'l', 'L':
	case 'u', 'U':
		res.isVertical = true
	case 'd', 'D':
		res.isVertical = true
		res.toForward = true
	default:
		err = errors.New("invalid direction code")
	}
	// Put offset and parts into res
	res.offset = offset
	res.body = parts
	return
}
