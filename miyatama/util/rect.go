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
