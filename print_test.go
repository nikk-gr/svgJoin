// Package svgJoin Copyright 2023 Gryaznov Nikita
// Licensed under the Apache License, Version 2.0
package svgJoin

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestChunk_Draw(t *testing.T) {
	type testCase struct {
		in  Chunk
		res string
		err error
	}
	doTest := func(i int, c testCase) {
		t.Logf("Test %d\tstart", i)
		res, err := c.in.Draw()

		if fmt.Sprint(err) != fmt.Sprint(c.err) {
			t.Errorf("%sTest %d failed%s\twant: %s, got: %s\n", red, i, normal, c.err, err)
		} else if err == nil && res != c.res {
			t.Errorf("%sTest %d failed%s\twant: %s, got: %s\n", red, i, normal, c.res, res)
			_ = os.WriteFile(fmt.Sprintf("TestChunk_Draw%d.svg", i), []byte(res), 777)

		} else {
			t.Logf("%sTest %d  success%s\t%s\n", green, i, normal, "")
		}
	}
	testArray := []testCase{
		{
			res: `<svg width="200" height="200" viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
<g clip-path="url(#clp0)">
<clipPath id="clp0">
<rect x="0.000000" y="0.000000" width="200.000000" height="200.000000" />
</clipPath>
<rect width="100" height="100" x="50" y="50" />
</g>
</svg>`,
			in: Chunk{
				viewport: xy{200, 200},
				viewBox:  xy{200, 200},
				position: xy{0, 0},
				body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
			},
		},
		{
			res: `<svg width="200" height="200" viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
<g clip-path="url(#clp0)">
<clipPath id="clp0">
<rect x="0.000000" y="0.000000" width="200.000000" height="200.000000" />
</clipPath>
<g transform="translate(-2, -0)">
<g transform="scale(2, 1)">
<rect width="100" height="100" x="50" y="50" />
</g>
</g>
</g>
</svg>`,
			in: Chunk{
				viewport: xy{200, 200},
				viewBox:  xy{100, 200},
				position: xy{2, 0},
				body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
			},
		},
		{
			res: "<svg width=\"0\" height=\"0\" viewBox=\"0 0 0 0\" xmlns=\"http://www.w3.org/2000/svg\">\n\n</svg>",
		},
		{
			in: Chunk{
				body: "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
			},
			err: errors.New("viewbox width, viewbox height, width, height is empty"),
		},
	}

	for k, v := range testArray {
		doTest(k, v)
	}
}
func TestGroup_Draw(t *testing.T) {
	type testCase struct {
		in  Group
		res string
		err error
	}
	doTest := func(i int, c testCase) {
		t.Logf("Test %d\tstart", i)
		res, err := c.in.Draw()

		if fmt.Sprint(err) != fmt.Sprint(c.err) {
			t.Errorf("%sTest %d failed%s\twant: %s, got: %s\n", red, i, normal, c.err, err)
		} else if err == nil && res != c.res {
			t.Errorf("%sTest %d failed%s\twant: %s, got: %s\n", red, i, normal, c.res, res)
			_ = os.WriteFile(fmt.Sprintf("TestGroup_Draw%d.svg", i), []byte(res), 777)
		} else {
			t.Logf("%sTest %d  success%s\t%s\n", green, i, normal, "")
		}
	}
	testArr := []testCase{
		{
			in: Group{
				toForward: true,
				body: []Part{
					Chunk{
						viewport: xy{200, 200},
						viewBox:  xy{200, 200},
						position: xy{0, 0},
						body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
					},
					Group{
						body: []Part{
							Chunk{
								viewport: xy{200, 200},
								viewBox:  xy{200, 200},
								position: xy{0, 0},
								body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
							},
							Chunk{
								viewport: xy{200, 200},
								viewBox:  xy{100, 200},
								position: xy{4, 0},
								body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
							},
						},
					},
				},
			},
			res: `<svg width="600" height="200" viewBox="0 0 600 200" xmlns="http://www.w3.org/2000/svg">
<g clip-path="url(#clp0)">
<clipPath id="clp0">
<rect x="0.000000" y="0.000000" width="600.000000" height="200.000000" />
</clipPath>

<g clip-path="url(#clp1)">
<clipPath id="clp1">
<rect x="0.000000" y="0.000000" width="200.000000" height="200.000000" />
</clipPath>
<rect width="100" height="100" x="50" y="50" />
</g>
<g clip-path="url(#clp2)">
<clipPath id="clp2">
<rect x="200.000000" y="0.000000" width="400.000000" height="200.000000" />
</clipPath>

<g clip-path="url(#clp3)">
<clipPath id="clp3">
<rect x="200.000000" y="0.000000" width="200.000000" height="200.000000" />
</clipPath>
<g transform="translate(196, -0)">
<g transform="scale(2, 1)">
<rect width="100" height="100" x="50" y="50" />
</g>
</g>
</g>
<g clip-path="url(#clp4)">
<clipPath id="clp4">
<rect x="400.000000" y="0.000000" width="200.000000" height="200.000000" />
</clipPath>
<g transform="translate(400, -0)">
<rect width="100" height="100" x="50" y="50" />
</g>
</g>
</g>

</g>

</svg>`,
		},
		{
			in: Group{
				toForward: true,
				body: []Part{
					Chunk{
						viewport: xy{200, 200},
						viewBox:  xy{0, 0},
						position: xy{0, 0},
						body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
					},
				},
			},
			err: errors.New("viewbox width, viewbox height is empty"),
		},
		{
			in: Group{
				toForward: false,
				align:     5,
				body: []Part{
					Chunk{
						viewport: xy{200, 200},
						viewBox:  xy{200, 200},
						position: xy{0, 0},
						body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
					},
				},
			},
			err: errors.New("invalid align code"),
		},
		{
			in: Group{
				isVertical: true,
				toForward:  false,
				align:      5,
				body: []Part{
					Chunk{
						viewport: xy{200, 200},
						viewBox:  xy{200, 200},
						position: xy{0, 0},
						body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
					},
				},
			},
			err: errors.New("invalid align code"),
		},
		{
			in: Group{
				toForward: false,
				align:     1,
				body: []Part{
					Chunk{
						viewport: xy{100, 100},
						viewBox:  xy{100, 100},
						position: xy{0, 0},
						body:     "<rect width=\"100\" height=\"100\" x=\"0\" y=\"0\" />",
					},
					Chunk{
						viewport: xy{50, 50},
						viewBox:  xy{50, 50},
						position: xy{0, 0},
						body:     "<rect width=\"50\" height=\"50\" x=\"0\" y=\"0\" fill=\"red\" />",
					},
				},
			},
			res: `<svg width="150" height="100" viewBox="0 0 150 100" xmlns="http://www.w3.org/2000/svg">
<g clip-path="url(#clp0)">
<clipPath id="clp0">
<rect x="0.000000" y="0.000000" width="150.000000" height="100.000000" />
</clipPath>

<g clip-path="url(#clp1)">
<clipPath id="clp1">
<rect x="0.000000" y="25.000000" width="50.000000" height="50.000000" />
</clipPath>
<g transform="translate(-0, 25)">
<rect width="50" height="50" x="0" y="0" fill="red" />
</g>
</g>
<g clip-path="url(#clp2)">
<clipPath id="clp2">
<rect x="50.000000" y="0.000000" width="100.000000" height="100.000000" />
</clipPath>
<g transform="translate(50, -0)">
<rect width="100" height="100" x="0" y="0" />
</g>
</g>
</g>

</svg>`,
		},
		{
			in: Group{
				offset:     5,
				isVertical: false,
				toForward:  false,
				align:      2,
				body: []Part{
					Chunk{
						viewport: xy{100, 100},
						viewBox:  xy{100, 100},
						position: xy{0, 0},
						body:     "<rect width=\"100\" height=\"100\" x=\"0\" y=\"0\" />",
					},
					Chunk{
						viewport: xy{50, 50},
						viewBox:  xy{50, 50},
						position: xy{0, 0},
						body:     "<rect width=\"50\" height=\"50\" x=\"0\" y=\"0\" fill=\"red\" />",
					},
				},
			},
			res: `<svg width="155" height="100" viewBox="0 0 155 100" xmlns="http://www.w3.org/2000/svg">
<g clip-path="url(#clp0)">
<clipPath id="clp0">
<rect x="0.000000" y="0.000000" width="155.000000" height="100.000000" />
</clipPath>

<g clip-path="url(#clp1)">
<clipPath id="clp1">
<rect x="0.000000" y="50.000000" width="50.000000" height="50.000000" />
</clipPath>
<g transform="translate(-0, 50)">
<rect width="50" height="50" x="0" y="0" fill="red" />
</g>
</g>
<g clip-path="url(#clp2)">
<clipPath id="clp2">
<rect x="55.000000" y="0.000000" width="100.000000" height="100.000000" />
</clipPath>
<g transform="translate(55, -0)">
<rect width="100" height="100" x="0" y="0" />
</g>
</g>
</g>

</svg>`,
		},
		{
			in: Group{
				isVertical: true,
				toForward:  false,
				align:      0,
				body: []Part{
					Chunk{
						viewport: xy{100, 100},
						viewBox:  xy{100, 100},
						position: xy{0, 0},
						body:     "<rect width=\"100\" height=\"100\" x=\"0\" y=\"0\" />",
					},
					Chunk{
						viewport: xy{50, 50},
						viewBox:  xy{50, 50},
						position: xy{0, 0},
						body:     "<rect width=\"50\" height=\"50\" x=\"0\" y=\"0\" fill=\"red\" />",
					},
				},
			},
			res: `<svg width="100" height="150" viewBox="0 0 100 150" xmlns="http://www.w3.org/2000/svg">
<g clip-path="url(#clp0)">
<clipPath id="clp0">
<rect x="0.000000" y="0.000000" width="100.000000" height="150.000000" />
</clipPath>

<g clip-path="url(#clp1)">
<clipPath id="clp1">
<rect x="0.000000" y="0.000000" width="50.000000" height="50.000000" />
</clipPath>
<rect width="50" height="50" x="0" y="0" fill="red" />
</g>
<g clip-path="url(#clp2)">
<clipPath id="clp2">
<rect x="0.000000" y="50.000000" width="100.000000" height="100.000000" />
</clipPath>
<g transform="translate(-0, 50)">
<rect width="100" height="100" x="0" y="0" />
</g>
</g>
</g>

</svg>`,
		},
		{
			in: Group{
				isVertical: true,
				toForward:  false,
				align:      1,
				body: []Part{
					Chunk{
						viewport: xy{100, 100},
						viewBox:  xy{100, 100},
						position: xy{0, 0},
						body:     "<rect width=\"100\" height=\"100\" x=\"0\" y=\"0\" />",
					},
					Chunk{
						viewport: xy{50, 50},
						viewBox:  xy{50, 50},
						position: xy{0, 0},
						body:     "<rect width=\"50\" height=\"50\" x=\"0\" y=\"0\" fill=\"red\" />",
					},
				},
			},
			res: `<svg width="100" height="150" viewBox="0 0 100 150" xmlns="http://www.w3.org/2000/svg">
<g clip-path="url(#clp0)">
<clipPath id="clp0">
<rect x="0.000000" y="0.000000" width="100.000000" height="150.000000" />
</clipPath>

<g clip-path="url(#clp1)">
<clipPath id="clp1">
<rect x="25.000000" y="0.000000" width="50.000000" height="50.000000" />
</clipPath>
<g transform="translate(25, -0)">
<rect width="50" height="50" x="0" y="0" fill="red" />
</g>
</g>
<g clip-path="url(#clp2)">
<clipPath id="clp2">
<rect x="0.000000" y="50.000000" width="100.000000" height="100.000000" />
</clipPath>
<g transform="translate(-0, 50)">
<rect width="100" height="100" x="0" y="0" />
</g>
</g>
</g>

</svg>`,
		},
		{
			in: Group{
				offset:     10,
				isVertical: true,
				toForward:  true,
				align:      2,
				body: []Part{
					Chunk{
						viewport: xy{100, 100},
						viewBox:  xy{100, 100},
						position: xy{0, 0},
						body:     "<rect width=\"100\" height=\"100\" x=\"0\" y=\"0\" />",
					},
					Chunk{
						viewport: xy{50, 50},
						viewBox:  xy{50, 50},
						position: xy{0, 0},
						body:     "<rect width=\"50\" height=\"50\" x=\"0\" y=\"0\" fill=\"red\" />",
					},
				},
			},
			res: `<svg width="100" height="160" viewBox="0 0 100 160" xmlns="http://www.w3.org/2000/svg">
<g clip-path="url(#clp0)">
<clipPath id="clp0">
<rect x="0.000000" y="0.000000" width="100.000000" height="160.000000" />
</clipPath>

<g clip-path="url(#clp1)">
<clipPath id="clp1">
<rect x="0.000000" y="0.000000" width="100.000000" height="100.000000" />
</clipPath>
<rect width="100" height="100" x="0" y="0" />
</g>
<g clip-path="url(#clp2)">
<clipPath id="clp2">
<rect x="50.000000" y="110.000000" width="50.000000" height="50.000000" />
</clipPath>
<g transform="translate(50, 110)">
<rect width="50" height="50" x="0" y="0" fill="red" />
</g>
</g>
</g>

</svg>`,
		},
	}
	for k, v := range testArr {
		doTest(k, v)
	}
}
