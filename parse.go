// Package svgJoin Copyright 2023 Gryaznov Nikita
// Licensed under the Apache License, Version 2.0
package svgJoin

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func Parse(svg string) (part Chunk, err error) {
	var (
		tmpErr    myErr
		noViewbox bool
	)
	firstLine, err := getFirstLine(svg)
	if err != nil {
		return
	}
	part.position.x, part.position.y, part.viewBox.x, part.viewBox.y, tmpErr = getViewBoxData(firstLine)
	switch tmpErr.code {
	case 0:
	case 11:
		noViewbox = true
	default:
		err = tmpErr
		return
	}
	part.viewport.x, part.viewport.y, tmpErr = getSize(firstLine)
	if tmpErr.code != 0 && noViewbox {
		err = errors.New("no svg size data")
		return
	}
	switch tmpErr.code {
	case 0:
		if noViewbox {
			part.viewBox = part.viewport
		}
	default:
		part.viewport = part.viewBox
	case 21:
		part.viewport.x = part.viewBox.x * part.viewport.y / part.viewBox.y
	case 22:
		part.viewport.y = part.viewBox.y * part.viewport.x / part.viewBox.x
	}
	svg = regexp.MustCompile("<\\?xml.*\\?>").ReplaceAllString(svg, "")
	svg = regexp.MustCompile("<svg.*?>").ReplaceAllString(svg, "")
	svg = regexp.MustCompile("</svg>").ReplaceAllString(svg, "")
	svg = strings.TrimSpace(svg)
	part.body = svg
	return
}

func getViewBoxData(firstLine string) (x0, y0, w, h float64, err myErr) {
	firstLine = regexp.MustCompile(",").ReplaceAllString(firstLine, " ")
	viewBoxFind := regexp.MustCompile("viewBox=\".+?\"") // regular expression for viewbox

	if viewbox := viewBoxFind.FindAllString(firstLine, 1); len(viewbox) > 0 {
		var tmp float64
		viewbox[0] = regexp.MustCompile("viewBox=\"").ReplaceAllString(viewbox[0], "")
		viewbox[0] = regexp.MustCompile("\"").ReplaceAllString(viewbox[0], "")

		values := regexp.MustCompile(" ").Split(viewbox[0], -1)
		var counter uint8
		for i := 0; i < len(values); i++ {
			if values[i] != "" {
				var err1 error
				tmp, err1 = strconv.ParseFloat(values[i], 64)
				if err1 != nil {
					err = wrapErr(10, err1)
				}
				switch counter {
				case 0:
					x0 = tmp
				case 1:
					y0 = tmp
				case 2:
					w = tmp
				case 3:
					h = tmp
				}
				counter++
			}
		}
		if counter != 4 {
			err = newErr(10)
		}

		return
	} else {
		err = newErr(11)
		return
	}
}

func getSize(firstLine string) (w, h float64, err myErr) {
	var (
		isW, isH   bool
		errW, errH error
	)
	wString := regexp.MustCompile("width=\"\\d+\\.?\\d*\"").FindAllString(firstLine, 1)
	hString := regexp.MustCompile("height=\"\\d+\\.?\\d*\"").FindAllString(firstLine, 1)
	if len(wString) > 0 {
		wString[0] = regexp.MustCompile("width=\"").ReplaceAllString(wString[0], "")
		wString[0] = regexp.MustCompile("\"").ReplaceAllString(wString[0], "")
		w, errW = strconv.ParseFloat(wString[0], 64)
		if errW == nil {
			isW = true
		}
	}
	if len(hString) > 0 {
		hString[0] = regexp.MustCompile("height=\"").ReplaceAllString(hString[0], "")
		hString[0] = regexp.MustCompile("\"").ReplaceAllString(hString[0], "")
		h, errH = strconv.ParseFloat(hString[0], 64)
		if errH == nil {
			isH = true
		}
	}
	switch {
	case !isW && !isH:
		err = wrapErr(20, errW, errH)
	case !isW && isH:
		err = wrapErr(21, errW)
	case isW && !isH:
		err = wrapErr(22, errH)
	}
	return
}

func getFirstLine(body string) (firstLine string, err error) {
	firstLineArr := regexp.MustCompile("<svg.*?>").FindAllString(body, 1)
	if len(firstLineArr) == 0 {
		err = errors.New("no svg first line match")
		return
	} else {
		firstLine = firstLineArr[0]
		return
	}
}
