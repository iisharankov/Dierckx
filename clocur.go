package dierckx

// #cgo FFLAGS: -fdefault-real-8
// #cgo LDFLAGS: -lm
// void splev (double*, int*, double*, int*, double*, double*, int*, int*);
// void curfit (int*, int*, double*, double*, double*, double*, double*, int*, double*, int*, int*, double*, double*, double*, double*, int*, int*, int*);
import "C"

// subroutine curfit(iopt,m,x,y,w,xb,xe,k,s,nest,n,t,c,fp, wrk,lwrk,iwrk,ier)
// subroutine clocur(iopt,ipar,idim,m,u,mx,x,w,k,s,nest,n,t,nc,c,fp, wrk,lwrk,iwrk,ier)
// void clocur (int*, int*, int*, int*, double*, int*, double*, double*, int*, double*, int*, int*, double*, int*, double*, double* double*, int*, int*, int*);







// Spline1D takes in x and y arrays, and ports the data to Direckx to return a spline
// func clocur(x, y []float64, k int) ([]float64, []float64, int) {

// 	return 
// }

// subroutine curfit(iopt,m,x,y,w,xb,xe,k,s,nest,n,t,c,fp,
// 	* wrk,lwrk,iwrk,ier)

// c   iopt  : int flag. 	varOfZero := C.int(0) 	-> &varOfZero
// c   m     : integer.  	castOfm := C.int(m) 	-> &castOfm
// c   s     : real. 		castOfs := C.double(s)	-> &castOfs
// c   x     : real array of dimension at least (m).
// 		copyOfX := make([]C.double, len(x))
// 		// If need to populate...
// 		for i := range x {
// 			copyOfX[i] = C.double(x[i])
// 		}
// 		usage: &copyOfX[0]

// c   iwrk  : integer array of dimension at least (nest):
// 		var iwrk = make([]C.int, nest)
// 		&iwrk[0]

