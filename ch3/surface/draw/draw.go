package draw

import (
	"fmt"
	"math"
    "strings"
)

const (
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func Surface(width int, height int) (string) {
	s := fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
    polygonPoints := make([]string, cells * cells)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j, width, height)
			bx, by := corner(i, j, width, height)
			cx, cy := corner(i, j+1, width, height)
			dx, dy := corner(i+1, j+1, width, height)
			polygonPoints = append(polygonPoints,fmt.Sprintf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>",
                ax, ay, bx, by, cx, cy, dx, dy))
		}
	}
    return s + strings.Join(polygonPoints,"\n") + "</svg>"
}

func corner(i, j, width, height int) (float64, float64) {
	xyscale := float64(width) / 2 / xyrange
	zscale := float64(height) * 0.4

	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := float64(width)/2 + (x-y)*cos30*xyscale
	sy := float64(height)/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

//!-
