// Package svgJoin Copyright 2023 Gryaznov Nikita
// Licensed under the Apache License, Version 2.0
package svgJoin

import (
	"errors"
	"fmt"
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
		} else {
			t.Logf("%sTest %d  success%s\t%s\n", green, i, normal, "")
		}
	}
	testArray := []testCase{
		{
			res: `<svg width="200" height="200" viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
<rect width="100" height="100" x="50" y="50" />
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
<g transform="scale(2, 1)">
<g transform="translate(2, 0)">
<rect width="100" height="100" x="50" y="50" />
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
			err:  errors.New("viewbox width, viewbox height, width, height is empty"),
		},
	}

	for k, v := range testArray {
		doTest(k, v)
	}
}
