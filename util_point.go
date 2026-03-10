package gointrum

import (
	"fmt"
	"strconv"
)

// Point - точка на карте.
type Point struct {
	Lat float64 // Широта
	Lon float64 // Долгота
}

// NewPoint конвертирует массив point в стурктуру.
func NewPoint(point [2]float64) Point {
	return Point{
		Lat: point[0],
		Lon: point[1],
	}
}

// NewPointFromStrings парсит массив строк point в стурктуру.
func NewPointFromStrings(point [2]string) (Point, error) {
	latF, err := strconv.ParseFloat(point[0], 64)
	if err != nil {
		return Point{}, fmt.Errorf("parse lat error: %w", err)
	}

	lonF, err := strconv.ParseFloat(point[1], 64)
	if err != nil {
		return Point{}, fmt.Errorf("parse lon error: %w", err)
	}

	return Point{
		Lat: latF,
		Lon: lonF,
	}, nil
}
