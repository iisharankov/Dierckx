package dierckx

// #cgo FFLAGS: -fdefault-real-8
// void splev (double*, int*, double*, int*, double*, double*, int*, int*);
// void curfit (int*, int*, double*, double*, double*, double*, double*, int*, double*, int*, int*, double*, double*, double*, double*, int*, int*, int*);
import "C"

import (
	"errors"
	"fmt"
)

// Spline1D takes in x and y arrays, and ports the data to Direckx to return a spline
func Spline1D(x, y []float64, k int) ([]float64, []float64, error) {

	/*
		Spline1D ports the given data so it can be ran by the included FORTRAN dierckx library. To do this, we heavily use
		the C library to convert all types to C types and pass all the references to the FORTRAN code. This function is truly
		jut a port from Golang to FORTRAN to use the required FORTRAN function.

		If you want detailed descriptions of the variables ported to the FORTRAN library, many of them contain the same names,
		and can be expanded by reading the very well documented FORTRAN docstring on the input types.
	*/

	// Used for returning arrays durring errors
	var emptyArray []float64

	/*
		The x and y parameters are arrays of float64, while the C compiler that ports to FORTRAN is an array of C.double.
		Obviously, this should not be an issue, since other than the name, these two types are identical in every way. But
		Golang is a bit too picky and even this needs to be converted. This is done by making an array of C.doubles and
		manually appending each element in x and y.
	*/
	copyOfX := make([]C.double, len(x))
	copyOfY := make([]C.double, len(y))

	// w is an array of ones of same length as the input arrays. This is done in the for loop
	w := make([]C.double, len(x))

	for i := range x {
		copyOfX[i] = C.double(x[i])
		copyOfY[i] = C.double(y[i])
		w[i] = 1
	}

	// Depreciated parameters left for possible later use
	// var periodic bool = false
	// var bc string = "nearest"

	var s float64 = 0.0
	var m int = len(x)

	if len(y) != m {
		return emptyArray, emptyArray, errors.New("length of x and y must match")
	}
	if len(w) != m {
		return emptyArray, emptyArray, errors.New("length of x and w must match")
	}
	if !(m > k) {
		return emptyArray, emptyArray, errors.New("k must be less than length(x)")
	}
	if !(1 <= k) || !(k <= 5) {
		return emptyArray, emptyArray, errors.New("1 <= k = $k <= 5 must hold")
	}

	var nest int = MaxOf(m+k+1, (2*k)+3)
	var lwrk int = m*(k+1) + nest*(7+(3*k))

	var wrk = make([]C.double, lwrk)
	var iwrk = make([]C.int, nest)
	var varOfZero C.int = 0

	// outputs
	var n int32
	var t = make([]C.double, nest)
	var c = make([]C.double, nest)
	var fp float64
	var ier int32

	// Cast all the above variables created to their C counterpart

	castOfm := C.int(m)
	castOfk := C.int(k)
	castOfs := C.double(s)
	castOfnest := C.int(nest)
	castOfn := C.int(n)
	castOffp := C.double(fp)
	castOflwrk := C.int(lwrk)
	castOfier := C.int(ier)

	// C.curfit is the name call given to the FORTRAN function that calculates the spline knots and coefficients
	// a := []C.double{0, m, x, y, w, x[0], x[len(x)-1], k, s, nest, n, t, c, fp, wrk, lwrk, iwrk, ier}
	C.curfit(&varOfZero, &castOfm, &copyOfX[0], &copyOfY[0], &w[0], &copyOfX[0], &copyOfX[len(x)-1], &castOfk, &castOfs, &castOfnest, &castOfn, &t[0], &c[0], &castOffp, &wrk[0], &castOflwrk, &iwrk[0], &castOfier) // pass addresses
	/*
		The above call populated the t and c variables, which now contain the location of all the knots in the spline, and the
		coefficients (respectively). These two arrays are of type []_Ctype_double, so naturally we must convert them back to
		[]float64 as the inverse of what we did at the beginning of the function.
	*/
	// Change back from list of []_Ctype_double to []float64

	if castOfier != -1 {
		fmt.Println("Fortran compiled with error of:", castOfier)
	}

	locOfKnots := make([]float64, len(t))
	coefficients := make([]float64, len(c))

	for i := range t {
		locOfKnots[i] = float64(t[i])
		coefficients[i] = float64(c[i])
	}

	// Return the arrays
	return locOfKnots, coefficients, nil
}

//////////////////////////////////////////////////////////////////
func Evaluate(spline, coefficients, xValues []float64, k int) []float64 {

	CCompatableSpline := make([]C.double, len(spline))
	CCompatableXValues := make([]C.double, len(xValues))
	CCompatableCoefficients := make([]C.double, len(coefficients))

	for i := range xValues {
		CCompatableXValues[i] = C.double(xValues[i])
	}

	for i := range spline {
		CCompatableSpline[i] = C.double(spline[i])
		CCompatableCoefficients[i] = C.double(coefficients[i])
	}

	numOfKnots := C.int(len(spline))

	degree := C.int(k)
	var m C.int = C.int(len(xValues))

	// // outputs
	var y = make([]C.double, len(xValues))
	var ier int32
	castOfier := C.int(ier)

	// Change back from list of []_Ctype_double to []float64
	listOfYValues := make([]float64, len(y))

	C.splev(&CCompatableSpline[0], &numOfKnots, &CCompatableCoefficients[0], &degree, &CCompatableXValues[0], &y[0], &m, &castOfier) // pass addresses
	// for i := range xValues {
	// 	// TODO If I use package unsaffe to somehow pass the full array, I don't need to call fortran len(xValues) times!
	// 	C.splev(&CCompatableSpline[0], &numOfKnots, &CCompatableCoefficients[0], &degree, &CCompatableXValues[i], &y[i], &m, &castOfier) // pass addresses
	// }
	for i := range listOfYValues {
		// Convert _Ctype_double to float64  and append to new list that can store it
		listOfYValues[i] = float64(y[i])
	}

	return listOfYValues
}

// MaxOf returns the max parameter passed in
func MaxOf(vars ...int) int {

	max := vars[0]
	for _, i := range vars {
		if max < i {
			max = i
		}
	}
	return max
}
