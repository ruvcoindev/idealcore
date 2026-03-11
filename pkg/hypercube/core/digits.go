package core

import "math"

// extractDigits извлекает три цифры числа
// Размещено в отдельном файле для гарантии видимости во всём пакете
func extractDigits(n int32) []int32 {
	n = int32(math.Abs(float64(n)))
	return []int32{
		(n / 100) % 10,
		(n / 10) % 10,
		n % 10,
	}
}
