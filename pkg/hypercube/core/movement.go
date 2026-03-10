package core

// MoveRoom перемещает комнату по векторам на заданное количество шагов
// В фильме "Куб" комнаты двигались по циклическому алгоритму
// Здесь мы моделируем жизненный путь как последовательность перемещений
// step - номер шага (0, 1, 2,...), определяет, какой вектор применяется
func MoveRoom(coords HypercubeCoords, vectors PersonVectors, step int) HypercubeCoords {
	newCoords := coords
	vectorIdx := step % 3 // Цикл из трех векторов, как в фильме

	// Применяем векторы к каждой координате
	newCoords.X += vectors.Full[0][vectorIdx]
	newCoords.Y += vectors.Full[1][vectorIdx]
	newCoords.Z += vectors.Full[2][vectorIdx]

	// Пересчитываем W (контрольную сумму) после перемещения
	newCoords.W = int32(sumDigits(int(newCoords.X)) +
		sumDigits(int(newCoords.Y)) +
		sumDigits(int(newCoords.Z)))

	// Нормализуем координаты до размера куба (0-25)
	// В фильме куб имел конечный размер, и комнаты не могли выйти за его пределы
	newCoords.X = modulo(newCoords.X, CubeSize)
	newCoords.Y = modulo(newCoords.Y, CubeSize)
	newCoords.Z = modulo(newCoords.Z, CubeSize)

	return newCoords
}

// modulo — корректный модуль для отрицательных чисел
// В Go оператор % может давать отрицательный результат,
// а нам нужно всегда получать положительные координаты
func modulo(a, m int32) int32 {
	return ((a % m) + m) % m
}

// IsBridgeRoom проверяет, является ли комната "мостом" (как комната 999 в фильме)
// В фильме комната 999 была выходом из куба
// Здесь мост - это возможность выйти из травматического цикла
func IsBridgeRoom(coords HypercubeCoords) bool {
	// W=27 получается, когда сумма цифр всех координат = 27 (как у 999)
	// Или любая координата равна 999
	return coords.W == 27 || coords.X == 999 || coords.Y == 999 || coords.Z == 999
}

// IsTrapRoom проверяет, является ли комната "ловушкой"
// В фильме некоторые комнаты содержали ловушки
// Здесь ловушка - это повторяющийся паттерн (все цифры одинаковые: 111, 222, 333...)
// Такие комнаты символизируют "застревание" в повторяющихся сценариях
func IsTrapRoom(coords HypercubeCoords) bool {
	digitsX := extractDigits(coords.X)
	digitsY := extractDigits(coords.Y)
	digitsZ := extractDigits(coords.Z)

	// Проверяем, все ли цифры одинаковы
	allEqual := func(d []int32) bool {
		if len(d) == 0 {
			return false
		}
		first := d[0]
		for _, val := range d {
			if val != first {
				return false
			}
		}
		return true
	}
	
	return allEqual(digitsX) || allEqual(digitsY) || allEqual(digitsZ)
}

// FindCycle определяет, через сколько шагов комната возвращается в исходное положение
// В фильме это был ключевой механизм - комнаты двигались по циклам
// Здесь цикл может показать, как быстро человек "возвращается" к своим паттернам
func FindCycle(start HypercubeCoords, vectors PersonVectors) int {
	current := start
	for step := 1; step <= 27; step++ { // Максимум 27 шагов (3^3)
		current = MoveRoom(current, vectors, step-1)
		if current == start {
			return step
		}
	}
	return -1 // Цикл не найден (в реальном кубе такого быть не может)
}

// EuclideanDistance вычисляет евклидово расстояние между двумя точками в 4D
// Используется для определения "близости" жизненных траекторий
func EuclideanDistance(a, b HypercubeCoords) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	dz := float64(a.Z - b.Z)
	dw := float64(a.W - b.W)
	return math.Sqrt(dx*dx + dy*dy + dz*dz + dw*dw)
}

// FindSafeRooms находит все комнаты-мосты на траектории за заданное количество шагов
// Помогает определить "точки выхода" из травматических циклов
func FindSafeRooms(coords HypercubeCoords, steps int) []HypercubeCoords {
	var safeRooms []HypercubeCoords
	vectors := CalculateVectors(coords)
	current := coords
	
	for i := 0; i < steps; i++ {
		current = MoveRoom(current, vectors, i)
		if IsBridgeRoom(current) {
			safeRooms = append(safeRooms, current)
		}
	}
	return safeRooms
}
