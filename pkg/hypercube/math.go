// Package hypercube реализует математику из фильма "Куб/Гиперкуб"
// с добавлением 4-го измерения (W) для моделирования жизненных траекторий
package hypercube

import (
	"github.com/ruvcoindev/idealcore/internal/domain"
	"math"
        "fmt"
)

// CubeSize — размер куба (26 комнат по каждой оси, как в фильме)
const CubeSize = 26

// ParseDateToCoords преобразует дату рождения в координаты HypercubeCoords
// Формат: DD.MM.YYYY → [0DD, 0MM, YYYY, checksum]
func ParseDateToCoords(date string) domain.HypercubeCoords {
	// Упрощенный парсер: ожидаем "YYYY-MM-DD" или "DD.MM.YYYY"
	var day, month, year int
	_, err := fmt.Sscanf(date, "%d.%d.%d", &day, &month, &year)
	if err != nil {
		_, err = fmt.Sscanf(date, "%d-%d-%d", &year, &month, &day)
		if err != nil {
			return domain.HypercubeCoords{X: 0, Y: 0, Z: 0, W: 0}
		}
	}

	// Преобразуем в трехзначные числа
	x := toThreeDigit(day)    // день
	y := toThreeDigit(month)  // месяц
	z := toThreeDigit(year)   // год

	// W — контрольная сумма (сумма цифр всех чисел)
	w := sumDigits(x) + sumDigits(y) + sumDigits(z)

	return domain.HypercubeCoords{
		X: int32(x),
		Y: int32(y),
		Z: int32(z),
		W: int32(w),
	}
}

// toThreeDigit дополняет число до трех знаков (15 → 015)
func toThreeDigit(n int) int {
	if n < 10 {
		return n // 3 → 003
	}
	if n < 100 {
		return n // 15 → 015
	}
	return n % 1000 // 1974 → 974
}

// sumDigits считает сумму цифр числа (974 → 9+7+4=20)
func sumDigits(n int) int {
	sum := 0
	for n > 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

// CalculateVectors вычисляет векторы движения комнаты по алгоритму из фильма
// Для каждого числа ABC: V1=A-B, V2=B-C, V3=C-A
func CalculateVectors(coords domain.HypercubeCoords) [3][3]int32 {
	numbers := [3]int32{coords.X, coords.Y, coords.Z}
	var vectors [3][3]int32

	for i, num := range numbers {
		digits := extractDigits(int(num)) // [A, B, C]
		if len(digits) < 3 {
			continue
		}
		vectors[i][0] = int32(digits[0] - digits[1]) // A-B
		vectors[i][1] = int32(digits[1] - digits[2]) // B-C
		vectors[i][2] = int32(digits[2] - digits[0]) // C-A
	}
	return vectors
}

// extractDigits извлекает цифры из числа (974 → [9,7,4])
func extractDigits(n int) []int {
	var digits []int
	for n > 0 {
		digits = append([]int{n % 10}, digits...)
		n /= 10
	}
	// Дополняем нулями слева до 3 цифр
	for len(digits) < 3 {
		digits = append([]int{0}, digits...)
	}
	return digits
}

// MoveRoom применяет векторы к координатам, возвращая новые координаты
func MoveRoom(coords domain.HypercubeCoords, vectors [3][3]int32, step int) domain.HypercubeCoords {
	newCoords := coords
	for i := 0; i < 3; i++ {
		// Применяем i-й вектор на шаге step (по модулю 3, т.к. цикл из 3 векторов)
		vectorIdx := step % 3
		newCoords.X += vectors[i][vectorIdx]
		newCoords.Y += vectors[i][vectorIdx]
		newCoords.Z += vectors[i][vectorIdx]
		// W не меняется (четвертое измерение — "мост")
	}
	// Ограничиваем координаты размером куба
	newCoords.X = modulo(newCoords.X, CubeSize)
	newCoords.Y = modulo(newCoords.Y, CubeSize)
	newCoords.Z = modulo(newCoords.Z, CubeSize)
	return newCoords
}

// modulo — корректный модуль для отрицательных чисел
func modulo(a, m int32) int32 {
	return ((a % m) + m) % m
}

// IsBridgeRoom проверяет, является ли комната "мостом" (как 999 в фильме)
// В 4D: если W == 27 (сумма цифр 999) или одна из координат == 999
func IsBridgeRoom(coords domain.HypercubeCoords) bool {
	return coords.W == 27 || coords.X == 999 || coords.Y == 999 || coords.Z == 999
}

// CalculateCycle определяет, на каком шаге цикла находится комната
// Цикл замыкается, когда координаты возвращаются к стартовым
func CalculateCycle(start, current domain.HypercubeCoords, vectors [3][3]int32) int {
	if start == current {
		return 0
	}
	test := start
	for step := 1; step <= 9; step++ { // макс 9 шагов (3 вектора × 3 оси)
		test = MoveRoom(test, vectors, step-1)
		if test == start {
			return step
		}
	}
	return -1 // цикл не найден
}

// AnalyzeCompatibility сравнивает две "комнаты" (людей) на совместимость траекторий
func AnalyzeCompatibility(person1, person2 string) (compatibility float64, syncSteps []int) {
	coords1 := ParseDateToCoords(person1)
	coords2 := ParseDateToCoords(person2)
	vectors1 := CalculateVectors(coords1)
	vectors2 := CalculateVectors(coords2)

	// 1. Проверка общей оси (как Y=1 у Виталия и Дины)
	commonAxes := 0
	if coords1.X == coords2.X { commonAxes++ }
	if coords1.Y == coords2.Y { commonAxes++ }
	if coords1.Z == coords2.Z { commonAxes++ }
	if coords1.W == coords2.W { commonAxes++ }

	// 2. Сравнение амплитуд векторов (стабильность vs динамика)
	amp1 := vectorAmplitude(vectors1)
	amp2 := vectorAmplitude(vectors2)
	ampDiff := math.Abs(float64(amp1 - amp2))

	// 3. Поиск синхронных шагов (когда комнаты в "безопасной" близости)
	for step := 0; step < 12; step++ {
		pos1 := coords1
		pos2 := coords2
		for s := 0; s < step; s++ {
			pos1 = MoveRoom(pos1, vectors1, s)
			pos2 = MoveRoom(pos2, vectors2, s)
		}
		distance := euclideanDistance(pos1, pos2)
		if distance < 5.0 { // порог "близости"
			syncSteps = append(syncSteps, step)
		}
	}

	// Формула совместимости (эвристическая)
	compatibility = float64(commonAxes)*0.3 + 
	                (1.0 - math.Min(ampDiff/10.0, 1.0))*0.4 + 
	                float64(len(syncSteps))*0.1

	return compatibility, syncSteps
}

func vectorAmplitude(vectors [3][3]int32) float64 {
	sum := 0.0
	for _, v := range vectors {
		for _, comp := range v {
			sum += math.Abs(float64(comp))
		}
	}
	return sum / 9.0 // среднее абсолютное значение
}

func euclideanDistance(a, b domain.HypercubeCoords) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	dz := float64(a.Z - b.Z)
	dw := float64(a.W - b.W)
	return math.Sqrt(dx*dx + dy*dy + dz*dz + dw*dw)
}
