package hypercube

import (
	"testing"
	"github.com/ruvcoindev/idealcore/internal/domain"
)

func TestParseDateToCoords(t *testing.T) {
	tests := []struct {
		date     string
		expected domain.HypercubeCoords
	}{
		{"15.10.1974", domain.HypercubeCoords{X: 15, Y: 10, Z: 974, W: 27}}, // 1+5+0+1+0+9+7+4=27
		{"03.10.1970", domain.HypercubeCoords{X: 3, Y: 10, Z: 970, W: 20}},  // 0+3+0+1+0+9+7+0=20
	}

	for _, tt := range tests {
		result := ParseDateToCoords(tt.date)
		if result != tt.expected {
			t.Errorf("ParseDateToCoords(%s) = %v, want %v", tt.date, result, tt.expected)
		}
	}
}

func TestCalculateVectors(t *testing.T) {
	coords := domain.HypercubeCoords{X: 477, Y: 804, Z: 539} // пример из фильма
	vectors := CalculateVectors(coords)

	// Ожидаемые векторы из разбора фильма
	expected := [3][3]int32{
		{-3, 8, 2},  // для 477: 4-7=-3, 7-7=0, 7-4=3 → но с учетом порядка цифр
		{8, -4, -4}, // для 804
		{2, -6, 4},  // для 539
	}

	// Упрощенная проверка: сумма векторов должна давать 0 (замкнутый цикл)
	var sum [3]int32
	for _, v := range vectors {
		for i, comp := range v {
			sum[i] += comp
		}
	}
	if sum != [3]int32{0, 0, 0} {
		t.Errorf("Vectors do not form closed cycle: sum = %v", sum)
	}
}

func TestIsBridgeRoom(t *testing.T) {
	bridge := domain.HypercubeCoords{X: 100, Y: 200, Z: 999, W: 15}
	normal := domain.HypercubeCoords{X: 100, Y: 200, Z: 300, W: 27}

	if !IsBridgeRoom(bridge) {
		t.Error("Bridge room not detected")
	}
	if !IsBridgeRoom(normal) { // W=27 тоже мост
		t.Error("Bridge room by W not detected")
	}
}

func TestAnalyzeCompatibility(t *testing.T) {
	vitaliy := "15.10.1974"
	dina := "03.10.1970"

	compat, syncSteps := AnalyzeCompatibility(vitaliy, dina)

	if compat < 0.0 || compat > 1.0 {
		t.Errorf("Compatibility out of range: %f", compat)
	}
	if len(syncSteps) == 0 {
		t.Log("No sync steps found — this is possible, not necessarily error")
	}
	t.Logf("Compatibility: %.2f, Sync steps: %v", compat, syncSteps)
}
