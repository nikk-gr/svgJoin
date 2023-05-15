// Package svgJoin Copyright 2023 Gryaznov Nikita
// Licensed under the Apache License, Version 2.0
package svgJoin

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Parse prepare a svg string to be joined with others
// Error returns if <svg…> tag has a wrong format
// Deeper structure is not checked
func Parse(svg string) (part Chunk, err error) {
	var (
		noViewbox bool
	)
	firstLine, err := getFirstLine(svg)
	if err != nil {
		return
	}
	part.position.x, part.position.y, part.viewBox.x, part.viewBox.y, err = getViewBoxData(firstLine)
	if err != nil {
		if err.Error() == "viewbox not found" {
			noViewbox = true
			err = nil
		} else {
			return
		}
	}
	part.viewport.x, part.viewport.y, err = getSize(firstLine)
	if err != nil && noViewbox {
		err = fmt.Errorf("no svg size data. Viewbox not found and %w", err)
		return
	}
	if err == nil {
		if noViewbox {
			part.viewBox = part.viewport
		}
	} else {
		v := strings.TrimSpace(err.Error())
		switch v {
		default:
			return
		case "no viewport data %!w(<nil>) %!w(<nil>)":
			part.viewport = part.viewBox
			err = nil
		case "no viewport width data %!w(<nil>)":
			part.viewport.x = part.viewBox.x * part.viewport.y / part.viewBox.y
			err = nil
		case "no viewport height data %!w(<nil>)":
			part.viewport.y = part.viewBox.y * part.viewport.x / part.viewBox.x
			err = nil
		}
	}
	svg = regexp.MustCompile("<\\?xml.*\\?>").ReplaceAllString(svg, "")
	svg = regexp.MustCompile("<svg.*?>").ReplaceAllString(svg, "")
	svg = regexp.MustCompile("</svg>").ReplaceAllString(svg, "")
	svg = strings.TrimSpace(svg)
	part.body = svg
	return
}

// getViewBoxData returns four numbers from the viewbox
func getViewBoxData(firstLine string) (x0, y0, w, h float64, err error) {
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
				tmp, err = strconv.ParseFloat(values[i], 64)
				if err != nil {
					err = fmt.Errorf("invalid viewbox format %w", err)
					return
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
			err = errors.New("invalid viewbox format")
		}

		return
	} else {
		err = errors.New("viewbox not found")
		return
	}
}

// getViewBoxData returns two numbers of the viewport
func getSize(firstLine string) (w, h float64, err error) {
	var (
		isW, isH   bool
		errW, errH error
	)
	wString := regexp.MustCompile("width=\".*?\"").FindAllString(firstLine, 1)
	hString := regexp.MustCompile("height=\".*?\"").FindAllString(firstLine, 1)

	if len(wString) > 0 {
		wString[0] = regexp.MustCompile("width=\"").ReplaceAllString(wString[0], "")
		wString[0] = regexp.MustCompile("\"").ReplaceAllString(wString[0], "")
		w, errW = parseSize(wString[0])
		if errW == nil {
			isW = true
		}
	}
	if len(hString) > 0 {
		hString[0] = regexp.MustCompile("height=\"").ReplaceAllString(hString[0], "")
		hString[0] = regexp.MustCompile("\"").ReplaceAllString(hString[0], "")
		h, errH = parseSize(hString[0])
		if errH == nil {
			isH = true
		}
	}
	switch {
	case !isW && !isH:
		err = fmt.Errorf("no viewport data %w %w", errW, errH)
	case !isW && isH:
		err = fmt.Errorf("no viewport width data %w", errW)
	case isW && !isH:
		err = fmt.Errorf("no viewport height data %w", errH)
	}
	return
}

// parseSize parse value with units of measurement
// Return result in px
func parseSize(s string) (f float64, e error) {
	var k float64 = 1
	pt := regexp.MustCompile("pt")
	px := regexp.MustCompile("px")
	in := regexp.MustCompile("in")
	cm := regexp.MustCompile("cm")
	mm := regexp.MustCompile("mm")
	pc := regexp.MustCompile("pc")
	em := regexp.MustCompile("em")
	ex := regexp.MustCompile("ex")
	switch {
	case pt.MatchString(s):
		s = pt.ReplaceAllString(s, "")
		k = 48.0 / 36.0
	case px.MatchString(s):
		s = px.ReplaceAllString(s, "")
	case in.MatchString(s):
		s = in.ReplaceAllString(s, "")
		k = 96
	case cm.MatchString(s):
		s = cm.ReplaceAllString(s, "")
		k = 37.7952755906
	case mm.MatchString(s):
		s = mm.ReplaceAllString(s, "")
		k = 3.7795275591
	case pc.MatchString(s):
		s = pc.ReplaceAllString(s, "")
		k = 16
	case em.MatchString(s):
		e = errors.New("em not supported. Use pt")
		return
	case ex.MatchString(s):
		e = errors.New("ex not supported. Use px")
		return

	}
	f, e = strconv.ParseFloat(s, 64)
	f *= k
	return
}

// selects the <svg… > tag from the body
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
