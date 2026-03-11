package core

import (
	"testing"
)

func TestMoveRoom(t *testing.T) {
	coords := HypercubeCoords{
		X: 10,
		Y: 10,
		Z: 10,
		W: 30,
	}

	vectors := PersonVectors{
		Full: [3][]int32{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
	}

	// Шаг 0
	newCoords := MoveRoom(coords, vectors, 0)
	if newCoords.X != (10+1)%CubeSize {
		t.Errorf("MoveRoom step 0 X = %d, want %d", newCoords.X, (10+1)%CubeSize)
	}

	// Шаг 1
	newCoords = MoveRoom(coords, vectors, 1)
	if newCoords.Y != (10+5)%CubeSize {
		t.Errorf("MoveRoom step 1 Y = %d, want %d", newCoords.Y, (10+5)%CubeSize)
	}

	// Шаг 2
	newCoords = MoveRoom(coords, vectors, 2)
	if newCoords.Z != (10+9)%CubeSize {
		t.Errorf("MoveRoom step 2 Z = %d, want %d", newCoords.Z, (10+9)%CubeSize)
	}

	// Шаг 3 (должен быть как шаг 0)
	newCoords = MoveRoom(coords, vectors, 3)
	if newCoords.X != (10+1)%CubeSize {
		t.Errorf("MoveRoom step 3 X = %d, want %d", newCoords.X, (10+1)%CubeSize)
	}
}

func TestModulo(t *testing.T) {
	tests := []struct {
		a, m, expected int32
	}{
		{5, 26, 5},
		{-5, 26, 21},
		{30, 26, 4},
		{-30, 26, 22},
		{0, 26, 0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := modulo(tt.a, tt.m)
			if result != tt.expected {
				t.Errorf("modulo(%d, %d) = %d, want %d", tt.a, tt.m, result, tt.expected)
			}
		})
	}
}

func TestIsBridgeRoom(t *testing.T) {
	tests := []struct {
		name     string
		coords   HypercubeCoords
		expected bool
	}{
		{"W=27 - мост", HypercubeCoords{W: 27}, true},
		{"X=999 - мост", HypercubeCoords{X: 999}, true},
		{"Y=999 - мост", HypercubeCoords{Y: 999}, true},
		{"Z=999 - мост", HypercubeCoords{Z: 999}, true},
		{"W=26 - не мост", HypercubeCoords{W: 26}, false},
		{"X=998 - не мост", HypercubeCoords{X: 998}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBridgeRoom(tt.coords)
			if result != tt.expected {
				t.Errorf("IsBridgeRoom(%v) = %v, want %v", tt.coords, result, tt.expected)
			}
		})
	}
}

func TestIsTrapRoom(t *testing.T) {
	tests := []struct {
		name     string
		coords   HypercubeCoords
		expected bool
	}{
		{"111 - ловушка (X)", HypercubeCoords{X: 111}, true},
		{"222 - ловушка (Y)", HypercubeCoords{Y: 222}, true},
		{"333 - ловушка (Z)", HypercubeCoords{Z: 333}, true},
		{"123 - не ловушка", HypercubeCoords{X: 123}, false},
		{"456 - не ловушка", HypercubeCoords{Y: 456}, false},
		{"789 - не ловушка", HypercubeCoords{Z: 789}, false},
		{"100 - не ловушка", HypercubeCoords{X: 100}, false},
		{"101 - не ловушка", HypercubeCoords{X: 101}, false},
		{"110 - не ловушка", HypercubeCoords{X: 110}, false},
		{"0 - не ловушка", HypercubeCoords{X: 0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTrapRoom(tt.coords)
			if result != tt.expected {
				t.Errorf("IsTrapRoom(%v) = %v, want %v", tt.coords, result, tt.expected)
			}
		})
	}
}

func TestFindCycle(t *testing.T) {
	coords := HypercubeCoords{X: 1, Y: 2, Z: 3, W: 6}
	vectors := CalculateVectors(coords)

	cycle := FindCycle(coords, vectors)
	if cycle < 1 || cycle > 27 {
		t.Errorf("FindCycle() returned %d, should be between 1 and 27", cycle)
	}
}

func TestEuclideanDistance(t *testing.T) {
	a := HypercubeCoords{X: 0, Y: 0, Z: 0, W: 0}
	b := HypercubeCoords{X: 3, Y: 4, Z: 0, W: 0}

	dist := EuclideanDistance(a, b)
	if dist != 5.0 {
		t.Errorf("EuclideanDistance = %f, want 5.0", dist)
	}
}

func TestFindSafeRooms(t *testing.T) {
	coords := HypercubeCoords{X: 1, Y: 2, Z: 3, W: 6}
	steps := 10

	rooms := FindSafeRooms(coords, steps)
	if len(rooms) > steps {
		t.Errorf("FindSafeRooms returned %d rooms, should be <= %d", len(rooms), steps)
	}
}
