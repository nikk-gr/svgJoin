package svgJoin

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestJoin(t *testing.T) {
	type testCase struct {
		direction direction
		align     align
		offset    float64
		in        []part
		res       Group
		err       error
	}
	doTest := func(i int, c testCase) {
		t.Logf("Test %d\tstart", i)
		res, err := Join(c.direction, c.align, c.offset, c.in...)

		if fmt.Sprint(err) != fmt.Sprint(c.err) {
			t.Errorf("%sTest %d failed%s\twant: %s, got: %s\n", red, i, normal, c.err, err)
		} else if err == nil && !reflect.DeepEqual(res, c.res) {
			t.Errorf("%sTest %d failed%s\twant: %s, got: %s\n", red, i, normal, fmt.Sprint(c.res), fmt.Sprint(res))
		} else {
			t.Logf("%sTest %d  success%s\t%s\n", green, i, normal, "")
		}
	}
	testArray := []testCase{
		{
			direction: Rightward,
			align:     Top,
			offset:    0,
			in: []part{
				Chunk{
					viewport: xy{200, 200},
					viewBox:  xy{200, 200},
					position: xy{0, 0},
					body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
				},
				Group{
					body: []part{
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
			res: Group{
				toForward: true,
				body: []part{
					Chunk{
						viewport: xy{200, 200},
						viewBox:  xy{200, 200},
						position: xy{0, 0},
						body:     "<rect width=\"100\" height=\"100\" x=\"50\" y=\"50\" />",
					},
					Group{
						body: []part{
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
		},
		{
			direction: Upward,
			align:     Midle,
			offset:    0,

			res: Group{
				toForward:  false,
				isVertical: true,
				align:      1,
			},
		},
		{
			direction: Leftward,
			align:     Bottom,
			offset:    0,

			res: Group{
				toForward:  false,
				isVertical: false,
				align:      2,
			},
		},
		{
			direction: Downward,
			align:     Midle,
			offset:    0,

			res: Group{
				toForward:  true,
				isVertical: true,
				align:      1,
			},
		},
		{
			direction: 5,
			align:     Midle,
			offset:    0,

			err: errors.New("invalid direction code"),
		},
		{
			direction: Downward,
			align:     4,
			offset:    0,

			err: errors.New("invalid align code"),
		},
	}

	for k, v := range testArray {
		doTest(k, v)
	}
}
