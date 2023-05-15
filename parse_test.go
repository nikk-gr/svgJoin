// Package svgJoin Copyright 2023 Gryaznov Nikita
// Licensed under the Apache License, Version 2.0
package svgJoin

import (
	"errors"
	"fmt"
	"testing"
)

const (
	red    string = "\x1b[31m"
	green  string = "\x1b[32m"
	normal string = "\x1b[0m"
)

func TestParse(t *testing.T) {
	type testCase struct {
		pic string
		res Chunk
		err error
	}
	doTest := func(i int, c testCase) {
		t.Logf("Test %d\tstart", i)
		res, err := Parse(c.pic)

		if fmt.Sprint(err) != fmt.Sprint(c.err) {
			t.Errorf("%sTest %d failed%s\twant: %s, got: %s\n", red, i, normal, c.err, err)
		} else if err == nil && res != c.res {
			t.Errorf("%sTest %d failed%s\twant: %s, got: %s\n", red, i, normal, fmt.Sprint(c.res), fmt.Sprint(res))
		} else {
			t.Logf("%sTest %d  success%s\t%s\n", green, i, normal, "")
		}
	}
	testArray := []testCase{
		{
			pic: `<svg width="200" height="200" viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
<rect width="100" height="100" x="50" y="50" />
</svg>`,
			res: Chunk{
				viewport: xy{200, 200},
				viewBox:  xy{200, 200},
				position: xy{0, 0},
				body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
			},
		},
		{
			pic: `<svg height="200" viewBox="0, 0, 200, 200" xmlns="http://www.w3.org/2000/svg">
<rect width="100" height="100" x="50" y="50" />
</svg>`,
			res: Chunk{
				viewport: xy{200, 200},
				viewBox:  xy{200, 200},
				position: xy{0, 0},
				body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
			},
		},
		{
			pic: `<svg width="200" viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
<rect width="100" height="100" x="50" y="50" />
</svg>`,
			res: Chunk{
				viewport: xy{200, 200},
				viewBox:  xy{200, 200},
				position: xy{0, 0},
				body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
			},
		},
		{
			pic: `<svg viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
<rect width="100" height="100" x="50" y="50" />
</svg>`,
			res: Chunk{
				viewport: xy{200, 200},
				viewBox:  xy{200, 200},
				position: xy{0, 0},
				body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
			},
		},
		{
			pic: `<svg width="200"  xmlns="http://www.w3.org/2000/svg">
<rect width="100" height="100" x="50" y="50" />
</svg>`,
			err: errors.New("no svg size data. Viewbox not found and no viewport height data %!w(<nil>)"),
		},
		{
			pic: `<svg width="200" height="200" xmlns="http://www.w3.org/2000/svg">
<rect width="100" height="100" x="50" y="50" />
</svg>`,
			res: Chunk{
				viewport: xy{200, 200},
				viewBox:  xy{200, 200},
				position: xy{0, 0},
				body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
			},
		},
		{
			pic: `<svg width="e3" height="200" xmlns="http://www.w3.org/2000/svg">
<rect width="100" height="100" x="50" y="50" />
</svg>`,
			err: errors.New("no svg size data. Viewbox not found and no viewport width data strconv.ParseFloat: parsing \"e3\": invalid syntax"),
		},
		{
			pic: `<sg width="e3" height="200" xmlns="http://www.w3.org/2000/svg">
<rect width="100" height="100" x="50" y="50" />
</svg>`,
			err: errors.New("no svg first line match"),
		},
		{
			pic: `<svg width="e3" height="200" viewBox="0 200 200" xmlns="http://www.w3.org/2000/svg">
<rect width="100" height="100" x="50" y="50" />
</svg>`,
			err: errors.New("invalid viewbox format"),
		},
		{
			pic: `<svg width="e3" height="200" viewBox="0 200 200 323 243" xmlns="http://www.w3.org/2000/svg">
<rect width="100" height="100" x="50" y="50" />
</svg>`,
			err: errors.New("invalid viewbox format"),
		},
		{
			pic: `<svg width="200" height="200" viewBox="23rd 200 200 323" xmlns="http://www.w3.org/2000/svg">
<rect width="100" height="100" x="50" y="50" />
</svg>`,
			err: errors.New("invalid viewbox format strconv.ParseFloat: parsing \"23rd\": invalid syntax"),
		},
		{
			pic: `<svg width="200" height="200" viewBox="0 0 ere 200" xmlns="http://www.w3.org/2000/svg">
<rect width="100" height="100" x="50" y="50" />
</svg>`,
			err: errors.New("invalid viewbox format strconv.ParseFloat: parsing \"ere\": invalid syntax"),
		},
		{
			pic: `<svg width="200" height="2egf" viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
<rect width="100" height="100" x="50" y="50" />
</svg>`,
			err: errors.New("no viewport height data strconv.ParseFloat: parsing \"2egf\": invalid syntax"),
		},
	}

	for k, v := range testArray {
		doTest(k, v)
	}
}

func TestParseSize(t *testing.T) {
	type testCase struct {
		in  string
		res float64
		err error
	}
	doTest := func(i int, c testCase) {
		t.Logf("Test %d\tstart", i)
		res, err := parseSize(c.in)

		if fmt.Sprint(err) != fmt.Sprint(c.err) {
			t.Errorf("%sTest %d failed%s\twant: %s, got: %s\n", red, i, normal, c.err, err)
		} else if err == nil && res != c.res {
			t.Errorf("%sTest %d failed%s\twant: %s, got: %s\n", red, i, normal, fmt.Sprint(c.res), fmt.Sprint(res))
		} else {
			t.Logf("%sTest %d  success%s\t%s\n", green, i, normal, "")
		}
	}
	testArray := []testCase{
		{
			in:  "10px",
			res: 10,
		},
		{
			in:  "10pt",
			res: 13.333333333333332,
		},
		{
			in:  "10pc",
			res: 160,
		},
		{
			in:  "10cm",
			res: 377.952755906,
		},
		{
			in:  "10mm",
			res: 37.795275591,
		},
		{
			in:  "10in",
			res: 960,
		},
		{
			in:  "10",
			res: 10,
		},
		{
			in:  "10em",
			err: errors.New("em not supported. Use pt"),
		},
		{
			in:  "10ex",
			err: errors.New("ex not supported. Use px"),
		},
		{
			in:  "abc",
			err: errors.New("strconv.ParseFloat: parsing \"abc\": invalid syntax"),
		},
	}

	for k, v := range testArray {
		doTest(k, v)
	}
}
