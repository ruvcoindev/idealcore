package core

import (
	"math"
        "fmt"
)

// MoveRoom перемещает комнату по векторам на заданное количество шагов
func MoveRoom(coords HypercubeCoords, vectors PersonVectors, step int) HypercubeCoords {
	newCoords := coords
	vectorIdx := step % 3

	newCoords.X += vectors.Full[0][vectorIdx]
	newCoords.Y += vectors.Full[1][vectorIdx]
	newCoords.Z += vectors.Full[2][vectorIdx]

	// ИСПРАВЛЕНИЕ: используем sumDigits из coordinates.go (она уже импортирована через пакет)
	newCoords.W = int32(sumDigits(int(newCoords.X)) +
		sumDigits(int(newCoords.Y)) +
		sumDigits(int(newCoords.Z)))

	newCoords.X = modulo(newCoords.X, CubeSize)
	newCoords.Y = modulo(newCoords.Y, CubeSize)
	newCoords.Z = modulo(newCoords.Z, CubeSize)

	return newCoords
}

// modulo — корректный модуль для отрицательных чисел
func modulo(a, m int32) int32 {
	return ((a % m) + m) % m
}

// IsBridgeRoom проверяет, является ли комната "мостом" (как комната 999 в фильме)
func IsBridgeRoom(coords HypercubeCoords) bool {
	return coords.W == 27 || coords.X == 999 || coords.Y == 999 || coords.Z == 999
}

// IsTrapRoom проверяет, является ли комната "ловушкой"
func IsTrapRoom(coords HypercubeCoords) bool {
	// Получаем цифры
	digitsX := extractDigits(coords.X)
	digitsY := extractDigits(coords.Y)
	digitsZ := extractDigits(coords.Z)
	
	// Для отладки - распечатаем, что получаем
	fmt.Printf("DEBUG - coords: X=%d, Y=%d, Z=%d\n", coords.X, coords.Y, coords.Z)
	fmt.Printf("DEBUG - digitsX: %v\n", digitsX)
	fmt.Printf("DEBUG - digitsY: %v\n", digitsY)
	fmt.Printf("DEBUG - digitsZ: %v\n", digitsZ)
	
	isTrap := func(d []int32) bool {
		return len(d) == 3 && d[0] == d[1] && d[1] == d[2] && d[0] > 0
	}
	
	result := isTrap(digitsX) || isTrap(digitsY) || isTrap(digitsZ)
	fmt.Printf("DEBUG - result: %v\n", result)
	fmt.Println("---")
	
	return result
}

// FindCycle определяет, через сколько шагов комната возвращается в исходное положение
func FindCycle(start HypercubeCoords, vectors PersonVectors) int {
	current := start
	for step := 1; step <= 27; step++ {
		current = MoveRoom(current, vectors, step-1)
		if current == start {
			return step
		}
	}
	return -1
}

// EuclideanDistance вычисляет евклидово расстояние между двумя точками в 4D
func EuclideanDistance(a, b HypercubeCoords) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	dz := float64(a.Z - b.Z)
	dw := float64(a.W - b.W)
	return math.Sqrt(dx*dx + dy*dy + dz*dz + dw*dw)
}

// FindSafeRooms находит все комнаты-мосты на траектории
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

// ВНИМАНИЕ: функции sumDigits и extractDigits НЕ НУЖНО определять здесь!
// Они уже определены в coordinates.go и доступны в этом пакете
