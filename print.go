package svgJoin

import (
	"errors"
	"fmt"
)

func (s *xy) add(x xy) {
	s.x += x.x
	s.y += x.y
}

func (s Chunk) print(pos xy) (result string, err error) {
	if s.body == "" || s.viewBox.x == 0 || s.viewBox.y == 0 || s.viewport.x == 0 || s.viewport.y == 0 {
		if s.body == "" && s.viewBox.x == 0 && s.viewBox.y == 0 && s.viewport.x == 0 && s.viewport.y == 0 {
			return
		}
		var errStr string
		if s.body == "" {
			errStr += "body"
		}
		if s.viewBox.x == 0 {
			if errStr != "" {
				errStr += ", "
			}
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
	s.position.add(pos)
	var (
		isTranslate, isScale bool
	)
	if s.viewBox.y != s.viewport.y || s.viewBox.x != s.viewport.x {
		isScale = true
	}
	if s.position.x != 0 || s.position.y != 0 {
		isTranslate = true
	}

	if isScale {
		result += fmt.Sprintf("<g transform=\"scale(%f, %f)\">\n", s.viewport.x/s.viewBox.x, s.viewport.y/s.viewBox.y)
	}
	if isTranslate {
		result += fmt.Sprintf("<g transform=\"translate(%f, %f)\">\n", s.position.x, s.position.y)
	}
	result += s.body
	if isScale {
		result += "\n</g>"
	}
	if isTranslate {
		result += "\n</g>"
	}
	return
}
func (s Chunk) size() xy {
	return s.viewport
}

// rightward, leftward, upward, downward
// 0 - top, left, 1 - mid, 2 = bottom, right

func (s Group) print(zero xy) (result string, err error) {
	size := s.size()

	var (
		localZero, resultZero xy
		tmp                   string
		from, stp             int
		check                 func(int) bool
	)
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
	for i := from; check(i); i += stp {
		localZero, err = getCoordinates(localZero, s.body[i].size(), size, s.offset, s.align, s.isVertical)
		resultZero = localZero
		resultZero.add(zero)
		if result != "" {
			result += "\n"
		}
		tmp, err = s.body[i].print(resultZero)
		if err != nil {
			return "", err
		}
		result += tmp
	}
	return
}

func getCoordinates(prev, partSize, groupSize xy, offset float64, align uint, isVertical bool) (new xy, err error) {
	if isVertical {
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

func (s Group) Draw() (pic string, err error) {
	size := s.size()
	pic, err = s.print(xy{})
	pic = fmt.Sprintf("<svg width=\"%g\" height=\"%g\" viewBox=\"0 0 %g %g\" version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\">\n%s\n</svg>", size.x, size.y, size.x, size.y, pic)
	return
}
func (s Chunk) Draw() (pic string, err error) {
	size := s.size()
	pic, err = s.print(xy{})
	pic = fmt.Sprintf("<svg width=\"%g\" height=\"%g\" viewBox=\"0 0 %g %g\" version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\">\n%s\n</svg>", size.x, size.y, size.x, size.y, pic)
	return
}
