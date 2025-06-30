package model

type Edge struct {
	From, To Point
	Weight   float32
}

type Point struct {
	X, Y float32
}
