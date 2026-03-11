package core

import (
	"math"
	"time"
)

// HypercubeCoords представляет координаты в 4-мерном гиперкубе
// X, Y, Z - три пространственных измерения (как в фильме "Куб")
// W - четвертое измерение, контрольная сумма, символизирует "мост" между измерениями
// Каждая координата - трехзначное число (от 0 до 999), но после перемещения
// нормализуется до размера куба (0-25) через модуль
type HypercubeCoords struct {
	X int32 // Первая координата (день рождения в преобразовании)
	Y int32 // Вторая координата (месяц рождения)
	Z int32 // Третья координата (год рождения)
	W int32 // Четвертая координата (контрольная сумма цифр)
}

// PersonVectors содержит векторы движения для человека
// Векторы вычисляются из цифр координат и показывают направление
// и амплитуду жизненных изменений
type PersonVectors struct {
	X    []int32      // Векторы для X-координаты: [A-B, B-C, C-A]
	Y    []int32      // Векторы для Y-координаты: [A-B, B-C, C-A]
	Z    []int32      // Векторы для Z-координаты: [A-B, B-C, C-A]
	Full [3][]int32   // Все векторы в одном массиве для удобства вычислений
}

// ParseDateToCoords преобразует дату рождения в координаты гиперкуба
// Философский смысл: дата рождения - это "входная точка" в гиперкуб жизни
// День -> X, Месяц -> Y, Год -> Z, сумма цифр -> W (ключ к выходу)
func ParseDateToCoords(date time.Time) HypercubeCoords {
	day := date.Day()
	month := int(date.Month())
	year := date.Year()

	// Преобразуем в трехзначные числа с сохранением позиций
	// Например: 5 мая -> 500 (день), 050 (месяц)
	x := toThreeDigitPreserve(day)
	y := toThreeDigitPreserve(month)
	z := toThreeDigitPreserve(year)
	
	// W - контрольная сумма, в фильме комната 999 была мостом
	// Здесь W = сумма всех цифр - "ключ" к пониманию системы
	w := sumDigits(x) + sumDigits(y) + sumDigits(z)

	return HypercubeCoords{
		X: int32(x),
		Y: int32(y),
		Z: int32(z),
		W: int32(w),
	}
}

// toThreeDigitPreserve дополняет число до трех знаков, сохраняя значимость позиций
// 3 -> 300 (3 сотни)
// 15 -> 150 (1 сотня, 5 десятков)
// 1974 -> 974 (последние три цифры)
func toThreeDigitPreserve(n int) int {
	if n < 10 {
		return n * 100 // 3 становится 300
	}
	if n < 100 {
		return n * 10 // 15 становится 150
	}
	return n % 1000 // 1974 становится 974
}

// sumDigits считает сумму цифр числа (нужно для вычисления W)
func sumDigits(n int) int {
	sum := 0
	n = int(math.Abs(float64(n)))
	for n > 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

// extractDigits извлекает три цифры из числа (сотни, десятки, единицы)
// Используется для вычисления векторов
func ExtractDigits(n int32) []int32 {
	n = int32(math.Abs(float64(n)))
	return []int32{
		(n / 100) % 10, // Сотни
		(n / 10) % 10,  // Десятки
		n % 10,         // Единицы
	}
}

// CalculateVectors вычисляет векторы движения для координат
// Из каждой трехзначной координаты ABC получаем три вектора:
// V1 = A-B (разница между первой и второй цифрой)
// V2 = B-C (разница между второй и третьей)
// V3 = C-A (разница между третьей и первой, замыкающая)
// Это создает циклическую систему, как в фильме
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
// Чем выше амплитуда, тем более "турбулентная" жизнь человека
// Низкая амплитуда может указывать на застой или депрессию
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
// Используется для определения "резонанса" или "синхронности" жизненных путей
func VectorSimilarity(v1, v2 PersonVectors) float64 {
	similarity := 0.0
	count := 0
	
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			// Чем меньше разница, тем выше сходство
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
