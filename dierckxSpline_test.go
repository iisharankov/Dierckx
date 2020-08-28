package dierckx

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"testing"
)

func outputToTxt(name string, value []float64) {
	// Open Azimuth file to write to
	AzFile, err := os.Create(name + ".txt")
	if err != nil {
		log.Println(err)
	}
	defer AzFile.Close()

	for i := range value {
		AzData := fmt.Sprintf("%v\n", value[i])
		_, err = io.WriteString(AzFile, AzData)
		if err != nil {
			log.Println(err)
		}
	}
}

type Spline1DData struct {
	filename string
	azPoints []float64
	timeVals []float64
}

var splineTestTable = []Spline1DData{
	// {"TestData1", []float64{80, 82, 85, 84, 82, 80}, []float64{1592771864, 1592771874, 1592771884, 1592771889, 1592771894, 1592771904}},
	{"TestData2", []float64{40, 70, 90, 150, 160, 200, 80, 30, 50, 80}, []float64{0, 50, 75, 140, 150, 200, 300, 350, 450, 500}},
	{"TestData3", []float64{80, 40, 80, 120, 80, 40, 80, 120, 80}, []float64{0, 50, 100, 150, 200, 250, 300, 350, 400}},
	// {"TestData4", []float64{80, 40, 80}, []float64{0, 100, 200}},
}

func TestSpline1D(t *testing.T) {

	// actual := Temp(tt.n)
	// if actual != tt.expected {
	// 	t.Errorf("Fib(%d): expected %d, actual %d", tt.n, tt.expected, actual)
	// }

	/*
		The below azPoints and timeVals was sent to the Julia package and is saved as JuliaOutput.txt.
		Send the same command to the two functions necessary to get the same spline and see if all the
		values of the two text files are identical
	*/

	for _, tt := range splineTestTable {
		// azPoints := []float64{80, 82, 85, 84, 82, 80}
		// timeVals := []float64{1592771864, 1592771874, 1592771884, 1592771889, 1592771894, 1592771904}

		azPoints := tt.azPoints
		timeVals := tt.timeVals
		fmt.Println(azPoints)

		var k int
		if len(azPoints) >= 4 {
			k = 3
		} else if len(azPoints) == 3 {
			k = 2
		} else if len(azPoints) == 2 {
			k = 1
		}
		fmt.Println(k)
		locOfKnots, coefficients, errr := Spline1D(timeVals, azPoints, k)
		if errr != 0 {
			// ErrorMsg := fmt.Sprintf("Spline1D returned an unexpected error: %v", err)
			// t.Errorf(ErrorMsg)
			return
		}

		// Need to Evaluate() over the locOfKnots at the same step size as Julia
		start := timeVals[0]
		end := timeVals[len(timeVals)-1]
		step := 1.0 / 20.0
		xValues := []float64{}
		for i := start; i < end; i += step {
			xValues = append(xValues, i)
		}

		// Send the xValues over to evaluate at each x point point
		yVals, _ := Evaluate(locOfKnots, coefficients, xValues, 3)

		// Print out to text file, then load it along side the Julia text file
		outputToTxt("TestOutput", yVals)

		file, err := os.Open("TestOutput.txt")
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		}

		juliaFile, err := os.Open(tt.filename + ".txt")
		defer juliaFile.Close()
		if err != nil {
			fmt.Println(err)
		}

		// Read each line of the two text files and make sure they're identical
		i := 0
		for {
			i++
			var julia float64
			var test float64

			var n int
			n, err = fmt.Fscanln(juliaFile, &julia)
			if n == 0 || err != nil {
				break
			}

			n, err = fmt.Fscanln(file, &test)
			if n == 0 || err != nil {
				break
			}

			// Floating point rounding, so make sure abs difference is very low
			if math.Abs(julia-test) > 0.00000001 {
				ErrorMsg := fmt.Sprintf("%v: The created spline was not identical to the Julia spline: %v", i, math.Abs(julia-test))
				t.Errorf(ErrorMsg)
				// continue
			}
		}
	}
}

// func TestSpline1DONLYONE(t *testing.T) {

// 	/*
// 		The below azPoints and timeVals was sent to the Julia package and is saved as JuliaOutput.txt.
// 		Send the same command to the two functions necessary to get the same spline and see if all the
// 		values of the two text files are identical
// 	*/
// 	azPoints := []float64{80, 82, 85, 84, 82, 80}
// 	timeVals := []float64{1592771864, 1592771874, 1592771884, 1592771889, 1592771894, 1592771904}

// 	locOfKnots, coefficients, err := Spline1D(timeVals, azPoints, 3)
// 	if err != nil {
// 		ErrorMsg := fmt.Sprintf("Spline1D returned an unexpected error: %v", err)
// 		t.Errorf(ErrorMsg)
// 	}

// 	// Need to Evaluate() over the locOfKnots at the same step size as Julia
// 	start := timeVals[0]
// 	end := timeVals[len(timeVals)-1]
// 	step := 1.0 / 20.0
// 	xValues := []float64{}
// 	for i := start; i < end; i += step {
// 		xValues = append(xValues, i)
// 	}

// 	// Send the xValues over to evaluate at each x point point
// 	yVals := Evaluate(locOfKnots, coefficients, xValues, 3)

// 	// Print out to text file, then load it along side the Julia text file
// 	outputToTxt("TestOutput", yVals)

// 	file, err := os.Open("TestOutput.txt")
// 	defer file.Close()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	juliaFile, err := os.Open("JuliaOutput_1.txt")
// 	defer juliaFile.Close()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Read each line of the two text files and make sure they're identical
// 	i := 0
// 	for {
// 		i++
// 		var julia float64
// 		var test float64

// 		var n int
// 		n, err = fmt.Fscanln(juliaFile, &julia)
// 		if n == 0 || err != nil {
// 			break
// 		}

// 		n, err = fmt.Fscanln(file, &test)
// 		if n == 0 || err != nil {
// 			break
// 		}
// 		fmt.Println(math.Abs(julia-test))
// 		// fmt.Println(math.Abs(julia-test))
// 		// Floating point rounding, so make sure abs difference is very low
// 		if math.Abs(julia-test) > 0.000001 {
// 			ErrorMsg := fmt.Sprintf("%v: The created spline was not identical to the Julia spline: %v", i, err)
// 			t.Errorf(ErrorMsg)
// 			// break
// 		}

// 	}
// }

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}

	}
}
