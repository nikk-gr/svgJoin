// Package svgJoin Copyright 2023 Gryaznov Nikita
// Licensed under the Apache License, Version 2.0
package svgJoin

import (
	"errors"
	"fmt"
)

func (s *xy) add(x xy) {
	s.x += x.x
	s.y += x.y
}
func (s *xy) sub(x xy) {
	s.x -= x.x
	s.y -= x.y
}

// print convert Chunk to svg group (<g>) that located in pos coordinates
// pos - coordinates of result <g>
// id - uniq number of clipPath generator
// id should be the same for all print functions that draw a resulting picture
func (s Chunk) print(pos xy, id *clipId) (result string, err error) {
	// check input
	if s.viewBox.x == 0 || s.viewBox.y == 0 || s.viewport.x == 0 || s.viewport.y == 0 {
		if s.body == "" && s.viewBox.x == 0 && s.viewBox.y == 0 && s.viewport.x == 0 && s.viewport.y == 0 {
			return
		}
		var errStr string
		if s.viewBox.x == 0 {
			errStr += "viewbox width"
		}
		if s.viewBox.y == 0 {
			if errStr != "" {
				errStr += ", "
			}
			errStr += "viewbox height"
		}
		if s.viewport.x == 0 {
			if errStr != "" {
				errStr += ", "
			}
			errStr += "width"
		}
		if s.viewport.y == 0 {
			if errStr != "" {
				errStr += ", "
			}
			errStr += "height"
		}
		errStr += " is empty"
		err = errors.New(errStr)
		return
	}

	// set the zero position
	s.position.sub(pos)
	// check if the viewport has non-zero coordinates or if its size differs from the viewport
	var (
		isTranslate, isScale bool
	)
	if s.viewBox.y != s.viewport.y || s.viewBox.x != s.viewport.x {
		isScale = true
	}
	if s.position.x != 0 || s.position.y != 0 {
		isTranslate = true
	}
	// cut off the parts that go beyond the viewbox
	clpId := id.get()
	result += fmt.Sprintf("<g clip-path=\"url(#clp%d)\">\n", clpId)
	result += fmt.Sprintf("<clipPath id=\"clp%d\">\n<rect x=\"%f\" y=\"%f\" width=\"%f\" height=\"%f\" />\n</clipPath>\n", clpId, pos.x, pos.y, s.viewport.x, s.viewport.y)
	// scale and translate
	if isTranslate {
		result += fmt.Sprintf("<g transform=\"translate(%-1g, %-1g)\">\n", -s.position.x, -s.position.y)
	}
	if isScale {
		result += fmt.Sprintf("<g transform=\"scale(%-1g, %-1g)\">\n", s.viewport.x/s.viewBox.x, s.viewport.y/s.viewBox.y)
	}
	// add body and close <g> tags
	result += s.body
	result += "\n</g>"
	if isScale {
		result += "\n</g>"
	}
	if isTranslate {
		result += "\n</g>"
	}
	return
}

// size returns the size of the Chunk
func (s Chunk) size() xy {
	return s.viewport
}

// Possible directions and alignment
// rightward, leftward, upward, downward
// 0 - top, left, 1 - mid, 2 = bottom, right

// print convert Group to svg group (<g>) that located in pos coordinates
// pos - coordinates of result <g>
// id - uniq number of clipPath generator
// id should be the same for all print functions that draw a resulting picture
func (s Group) print(pos xy, id *clipId) (result string, err error) {
	size := s.size()
	var (
		localZero, resultZero xy
		tmp                   string
		from, stp             int
		check                 func(int) bool
	)

	// set the direction of the bodies array iteration
	if s.toForward {
		from = 0
		check = func(i int) bool {
			return i < len(s.body)
		}
		stp = 1
	} else {
		from = len(s.body) - 1
		check = func(i int) bool {
			return i >= 0
		}
		stp = -1
	}
	// cut off the parts that go beyond the viewbox
	// ToDo possibly not required
	clpId := id.get()
	result += fmt.Sprintf("<g clip-path=\"url(#clp%d)\">\n", clpId)
	result += fmt.Sprintf("<clipPath id=\"clp%d\">\n<rect x=\"%f\" y=\"%f\" width=\"%f\" height=\"%f\" />\n</clipPath>\n", clpId, pos.x, pos.y, size.x, size.y)
	// iterate bodies
	for i := from; check(i); i += stp {
		localZero, err = getCoordinates(localZero, s.body[i].size(), size, s.offset, s.align, s.isVertical)
		if err != nil {
			return "", err
		}
		if s.isVertical {
			resultZero.x = localZero.x
		} else {
			resultZero.y = localZero.y
		}

		resultZero.add(pos)
		if result != "" {
			result += "\n"
		}
		tmp, err = s.body[i].print(resultZero, id)
		if err != nil {
			return "", err
		}
		resultZero = localZero
		result += tmp
	}
	result += "\n</g>\n"
	return
}

// getCoordinates calculate the coordinates of the upper right corner of the part
func getCoordinates(prev, partSize, groupSize xy, offset float64, align uint8, isVertical bool) (new xy, err error) {
	if !isVertical {
		new.x = prev.x + offset + partSize.x
		switch align {
		case 0:
			new.y = prev.y
		case 1:
			new.y = (groupSize.y - partSize.y) / 2
		case 2:
			new.y = groupSize.y - partSize.y
		default:
			err = errors.New("invalid align code")
			return
		}
	} else {
		new.y = prev.y + offset + partSize.y
		switch align {
		case 0:
			new.x = prev.x
		case 1:
			new.x = (groupSize.x - partSize.x) / 2
		case 2:
			new.x = groupSize.x - partSize.x
		default:
			err = errors.New("invalid align code")
			return
		}
	}
	return
}

// size returns the size of the Group
func (s Group) size() (o xy) {
	var tmp xy
	if !s.isVertical {
		for key, val := range s.body {
			tmp = val.size()
			if o.y < tmp.y {
				o.y = tmp.y
			}
			if key > 0 {
				o.x += s.offset
			}
			o.x += tmp.x
		}

	} else {
		for key, val := range s.body {
			tmp = val.size()
			if key > 0 {
				o.y += s.offset
			}
			o.y += tmp.y
			if o.x < tmp.x {
				o.x = tmp.x
			}
		}
	}
	return
}

// Draw returns a svg string of Group
func (s Group) Draw() (pic string, err error) {
	size := s.size()
	var id clipId
	pic, err = s.print(xy{}, &id)
	pic = fmt.Sprintf("<svg width=\"%g\" height=\"%g\" viewBox=\"0 0 %g %g\" xmlns=\"http://www.w3.org/2000/svg\">\n%s\n</svg>", size.x, size.y, size.x, size.y, pic)
	return
}

// Draw returns a svg string of Chunk
func (s Chunk) Draw() (pic string, err error) {
	size := s.size()
	var id clipId
	pic, err = s.print(xy{}, &id)
	pic = fmt.Sprintf("<svg width=\"%g\" height=\"%g\" viewBox=\"0 0 %g %g\" xmlns=\"http://www.w3.org/2000/svg\">\n%s\n</svg>", size.x, size.y, size.x, size.y, pic)
	return
}
