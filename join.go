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
	Midle     align     = 1
	Bottom    align     = 2
	Right     align     = 2
)

func Join(direction direction, align align, offset float64, parts ...part) (res Group, err error) {
	if align > 2 {
		err = errors.New("invalid align code")
		return
	} else {
		res.align = uint8(align)
	}
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
	res.offset = offset
	res.body = parts
	return
}
