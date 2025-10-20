package util

import "fmt"

type Rect struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

func (r *Rect) Width() int {
	return r.Right - r.Left
}

func (r *Rect) Height() int {
	return r.Bottom - r.Top
}

func (r *Rect) ToString() string {
	return fmt.Sprintf("{left: %d, top: %d, right: %d, bottom: %d}", r.Left, r.Top, r.Right, r.Bottom)
}

func (r *Rect) Contains(x, y int) bool {
	return r.Left <= x && x < r.Right && r.Top <= y && y < r.Bottom
}

func (r *Rect) SplitGrid(x, y int) []Rect {
	width := r.Width() / x
	height := r.Height() / y
	rects := []Rect{}
	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			rects = append(
				rects,
				Rect{
					Left:   r.Left + width*i,
					Top:    r.Top + height*j,
					Right:  r.Left + width*(i+1),
					Bottom: r.Top + height*(j+1),
				},
			)
		}
	}
	return rects
}
