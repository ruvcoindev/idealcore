package core

import (
	"math"
	"time"
)

// HypercubeCoords представляет координаты в 4-мерном гиперкубе
type HypercubeCoords struct {
	X int32
	Y int32
	Z int32
	W int32
}

// PersonVectors содержит векторы движения для человека
type PersonVectors struct {
	X    []int32
	Y    []int32
	Z    []int32
	Full [3][]int32
}

// ParseDateToCoords преобразует дату рождения в координаты гиперкуба
func ParseDateToCoords(date time.Time) HypercubeCoords {
	day := date.Day()
	month := int(date.Month())
	year := date.Year()

	x := toThreeDigitPreserve(day)
	y := toThreeDigitPreserve(month)
	z := toThreeDigitPreserve(year)
	
	w := sumDigits(x) + sumDigits(y) + sumDigits(z)

	return HypercubeCoords{
		X: int32(x),
		Y: int32(y),
		Z: int32(z),
		W: int32(w),
	}
}

// toThreeDigitPreserve дополняет число до трех знаков
func toThreeDigitPreserve(n int) int {
	if n < 10 {
		return n * 100 // 3 -> 300
	}
	if n < 100 {
		return n * 10 // 15 -> 150
	}
	return n % 1000 // 1974 -> 974
}

// sumDigits считает сумму цифр числа
func sumDigits(n int) int {
	sum := 0
	n = int(math.Abs(float64(n)))
	for n > 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

// extractDigits извлекает три цифры из числа
// ИСПРАВЛЕНО: всегда возвращает 3 цифры, дополняя нулями слева
func extractDigits(n int32) []int32 {
	n = int32(math.Abs(float64(n)))
	
	// Разбиваем на сотни, десятки, единицы
	hundreds := (n / 100) % 10
	tens := (n / 10) % 10
	units := n % 10
	
	return []int32{hundreds, tens, units}
}

// CalculateVectors вычисляет векторы для координат
func CalculateVectors(coords HypercubeCoords) PersonVectors {
	digitsX := extractDigits(coords.X)
	digitsY := extractDigits(coords.Y)
	digitsZ := extractDigits(coords.Z)
	
	return PersonVectors{
		X: []int32{
			digitsX[0] - digitsX[1],
			digitsX[1] - digitsX[2],
			digitsX[2] - digitsX[0],
		},
		Y: []int32{
			digitsY[0] - digitsY[1],
			digitsY[1] - digitsY[2],
			digitsY[2] - digitsY[0],
		},
		Z: []int32{
			digitsZ[0] - digitsZ[1],
			digitsZ[1] - digitsZ[2],
			digitsZ[2] - digitsZ[0],
		},
		Full: [3][]int32{
			{
				digitsX[0] - digitsX[1],
				digitsX[1] - digitsX[2],
				digitsX[2] - digitsX[0],
			},
			{
				digitsY[0] - digitsY[1],
				digitsY[1] - digitsY[2],
				digitsY[2] - digitsY[0],
			},
			{
				digitsZ[0] - digitsZ[1],
				digitsZ[1] - digitsZ[2],
				digitsZ[2] - digitsZ[0],
			},
		},
	}
}

// VectorAmplitude вычисляет среднюю амплитуду векторов
func VectorAmplitude(vectors PersonVectors) float64 {
	sum := 0.0
	count := 0
	for _, vec := range vectors.Full {
		for _, v := range vec {
			sum += math.Abs(float64(v))
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

// VectorSimilarity вычисляет сходство векторов двух людей
func VectorSimilarity(v1, v2 PersonVectors) float64 {
	similarity := 0.0
	count := 0
	
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			diff := math.Abs(float64(v1.Full[i][j] - v2.Full[i][j]))
			similarity += 1.0 - math.Min(diff/20.0, 1.0)
			count++
		}
	}
	
	if count == 0 {
		return 0
	}
	return similarity / float64(count)
}
