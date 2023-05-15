package svgJoin_test

import (
	"github.com/nikk-gr/svgJoin"
	"os"
	"testing"
)

func TestJoin2(t *testing.T) {
	// Make arrays of pictures paths
	upperNames := []string{
		"./testdata/STS1.svg",
		"./testdata/FC1.svg",
	}
	lowerNames := []string{
		"./testdata/EH.svg",
		"./testdata/CW.svg",
		"./testdata/WMH.svg",
		"./testdata/STS2.svg",
		"./testdata/FC2.svg",
	}
	leftNames := []string{
		"./testdata/FC3.svg",
		"./testdata/FVS1.svg",
		"./testdata/RR0.svg",
		"./testdata/FVS2.svg",
	}

	// Read pictures from paths and parse it
	upperParts := make([]svgJoin.Part, len(upperNames))
	lowerParts := make([]svgJoin.Part, len(lowerNames))
	leftParts := make([]svgJoin.Part, len(leftNames))

	for i := 0; i < len(upperNames); i++ {
		file, err := os.ReadFile(upperNames[i])
		if err != nil {
			t.Fatal(err.Error())
		}
		upperParts[i], err = svgJoin.Parse(string(file))
		if err != nil {
			t.Fatal(err.Error())
		}
	}
	for i := 0; i < len(lowerNames); i++ {
		file, err := os.ReadFile(lowerNames[i])
		if err != nil {
			t.Fatal(err.Error())
		}
		lowerParts[i], err = svgJoin.Parse(string(file))
		if err != nil {
			t.Fatal(err.Error())
		}
	}
	for i := 0; i < len(leftNames); i++ {
		file, err := os.ReadFile(leftNames[i])
		if err != nil {
			t.Fatal(err.Error())
		}
		leftParts[i], err = svgJoin.Parse(string(file))
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	// Join pictures
	var (
		UpperRight, LowerRight, Left, Right, Result svgJoin.Group
		err                                         error
	)
	UpperRight, err = svgJoin.Join(svgJoin.Rightward, svgJoin.Top, 10, upperParts...)
	if err != nil {
		t.Fatal(err.Error())
	}
	LowerRight, err = svgJoin.Join(svgJoin.Rightward, svgJoin.Bottom, 10, lowerParts...)
	if err != nil {
		t.Fatal(err.Error())
	}
	Right, err = svgJoin.Join(svgJoin.Downward, svgJoin.Left, 0, UpperRight, LowerRight)
	if err != nil {
		t.Fatal(err.Error())
	}
	Left, err = svgJoin.Join(svgJoin.Rightward, svgJoin.Bottom, 10, leftParts...)
	if err != nil {
		t.Fatal(err.Error())
	}
	Result, err = svgJoin.Join(svgJoin.Rightward, svgJoin.Bottom, 10, Left, Right)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Draw result svg
	var svg string
	svg, err = Result.Draw()

	// Save to file
	_ = os.WriteFile("Result.svg", []byte(svg), 770)
}

func TestJoin3(t *testing.T) {
	Names := []string{
		"./testdata/EH.svg",
		"./testdata/CW.svg",
		"./testdata/WMH.svg",
		"./testdata/STS2.svg",
		"./testdata/FC2.svg",
	}

	// Read pictures from paths and parse it
	Parts := make([]svgJoin.Part, len(Names))

	for i := 0; i < len(Names); i++ {
		file, err := os.ReadFile(Names[i])
		if err != nil {
			t.Fatal(err.Error())
		}
		Parts[i], err = svgJoin.Parse(string(file))
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	// Join pictures
	Result, err := svgJoin.Join(svgJoin.Rightward, svgJoin.Top, 10, Parts...)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Draw result svg
	var svg string
	svg, err = Result.Draw()

	// Save to file
	_ = os.WriteFile("Result2.svg", []byte(svg), 770)
}
