// Package svgJoin Copyright 2023 Gryaznov Nikita
// Licensed under the Apache License, Version 2.0
package svgJoin

import "strconv"

type myErr struct {
	code uint16
	msg  string
}

func newErr(code uint16) myErr {
	return myErr{
		code: code,
	}
}

func wrapErr(code uint16, errs ...error) (err myErr) {
	err.code = code
	for _, val := range errs {
		if val != nil {
			if err.msg != "" {
				err.msg += " "
			}
			err.msg += val.Error()
		}
	}
	return
}

func (err myErr) Error() string {
	part1 := err.readCodeMsg()
	part2 := err.msg
	if part1 != "" && part2 != "" {
		return part1 + " " + part2
	} else {
		return part1 + part2
	}
}

func (err myErr) readCodeMsg() string {
	switch err.code {
	case 10:
		return "invalid viewbox format"
	case 11:
		return "viewbox not found"
	case 20:
		return "no viewport data"
	case 21:
		return "no sizeWidth data"
	case 22:
		return "no sizeHeight data"
	default:
		return "invalid error code" + strconv.FormatUint(uint64(err.code), 16)
	}
}
