package model

// WindowInfo represents a system window for task binding
type WindowInfo struct {
	Handle string `json:"handle"`
	Title  string `json:"title"`
	Process string `json:"process"`
}

// Rect represents a screen region
type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Point represents a screen coordinate
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}
