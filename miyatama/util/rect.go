package util

import "fmt"

type Rect[T Numeric] struct {
	Left   T
	Top    T
	Right  T
	Bottom T
}

func (r *Rect[T]) Width() T {
	return r.Right - r.Left
}

func (r *Rect[T]) Height() T {
	return r.Bottom - r.Top
}

func (r *Rect[T]) ToString() string {
	return fmt.Sprintf("{left: %d, top: %d, right: %d, bottom: %d}", r.Left, r.Top, r.Right, r.Bottom)
}

func (r *Rect[T]) Contains(x, y T) bool {
	return r.Left <= x && x < r.Right && r.Top <= y && y < r.Bottom
}

func (r *Rect[T]) SplitGrid(x, y T) []Rect[T] {
	width := r.Width() / x
	height := r.Height() / y
	rects := []Rect[T]{}
	for j := 0; j < int(y); j++ {
		for i := 0; i < int(x); i++ {
			rects = append(
				rects,
				Rect[T]{
					Left:   r.Left + width*T(i),
					Top:    r.Top + height*T(j),
					Right:  r.Left + width*(T(i)+1),
					Bottom: r.Top + height*(T(j)+1),
				},
			)
		}
	}
	return rects
}
