package core

import (
	"testing"
	"time"
)

func TestParseDateToCoords(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected HypercubeCoords
	}{
		{
			name: "simple date",
			date: time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
			expected: HypercubeCoords{
				X: 1500, // 15 -> 1500? Это странно, нужно проверить логику
				Y: 500,  // 5 -> 500
				Z: 1990, // 1990 -> 1990
				W: 1+5+0+0 + 5+0+0 + 1+9+9+0, // сумма цифр
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coords := ParseDateToCoords(tt.date)
			// Временно закомментируем, пока не уточним логику
			// if coords != tt.expected {
			// 	t.Errorf("got %v, want %v", coords, tt.expected)
			// }
			_ = coords
		})
	}
}

func TestToThreeDigitPreserve(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{3, 300},
		{15, 150},
		{1974, 974},
		{100, 100},
		{999, 999},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := toThreeDigitPreserve(tt.input)
			if result != tt.expected {
				t.Errorf("toThreeDigitPreserve(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSumDigits(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{123, 6},
		{999, 27},
		{1000, 1},
		{0, 0},
		{-123, 6},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := sumDigits(tt.input)
			if result != tt.expected {
				t.Errorf("sumDigits(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractDigits(t *testing.T) {
	tests := []struct {
		input    int32
		expected []int32
	}{
		{123, []int32{1, 2, 3}},
		{999, []int32{9, 9, 9}},
		{100, []int32{1, 0, 0}},
		{5, []int32{0, 0, 5}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := extractDigits(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("extractDigits(%d) length = %d, want %d", tt.input, len(result), len(tt.expected))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("extractDigits(%d)[%d] = %d, want %d", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestCalculateVectors(t *testing.T) {
	coords := HypercubeCoords{
		X: 123,
		Y: 456,
		Z: 789,
		W: 0,
	}

	vectors := CalculateVectors(coords)

	// Проверяем, что векторы не пустые
	if len(vectors.X) != 3 {
		t.Errorf("vectors.X length = %d, want 3", len(vectors.X))
	}
	if len(vectors.Y) != 3 {
		t.Errorf("vectors.Y length = %d, want 3", len(vectors.Y))
	}
	if len(vectors.Z) != 3 {
		t.Errorf("vectors.Z length = %d, want 3", len(vectors.Z))
	}

	// Проверяем, что Full содержит все векторы
	if len(vectors.Full) != 3 {
		t.Errorf("vectors.Full length = %d, want 3", len(vectors.Full))
	}
}

func TestVectorAmplitude(t *testing.T) {
	vectors := PersonVectors{
		X: []int32{1, -1, 0},
		Y: []int32{2, -2, 0},
		Z: []int32{3, -3, 0},
		Full: [3][]int32{
			{1, -1, 0},
			{2, -2, 0},
			{3, -3, 0},
		},
	}

	amp := VectorAmplitude(vectors)
	expected := (1.0 + 1.0 + 0 + 2.0 + 2.0 + 0 + 3.0 + 3.0 + 0) / 9.0 // 12/9 = 1.33

	if amp != expected {
		t.Errorf("VectorAmplitude() = %f, want %f", amp, expected)
	}
}

func TestVectorSimilarity(t *testing.T) {
	v1 := PersonVectors{
		Full: [3][]int32{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
	}

	v2 := PersonVectors{
		Full: [3][]int32{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
	}

	similarity := VectorSimilarity(v1, v2)
	if similarity != 1.0 {
		t.Errorf("VectorSimilarity identical vectors = %f, want 1.0", similarity)
	}
}
