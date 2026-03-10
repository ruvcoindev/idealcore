package analysis

import (
	"math"
	"github.com/ruvcoindev/idealcore/hypercube/core"
	"github.com/ruvcoindev/idealcore/hypercube/data"
)

// CompatibilityResult содержит результат анализа совместимости двух людей
type CompatibilityResult struct {
	Score           float64  // Общий балл совместимости (0-1)
	CommonAxes      int      // Количество общих осей координат (0-4)
	AmpDiff         float64  // Разница в амплитуде векторов
	SyncSteps       int      // Количество синхронных шагов (0-12)
	Level           string   // Уровень: "Низкая", "Средняя", "Высокая"
	KarmicFactor    float64  // Кармический фактор (>1 если есть кармическая связь)
	Recommendations []string // Рекомендации
}

// GroupCompatibilityResult содержит результат анализа совместимости группы
type GroupCompatibilityResult struct {
	AverageScore    float64   // Средняя совместимость в группе
	ConflictScore   float64   // Уровень конфликтности (0-1)
	StabilityScore  float64   // Уровень стабильности (0-1)
	GroupDynamics   []string  // Описание групповой динамики
	RequiredActions []string  // Необходимые действия
}

// CalculateCompatibility вычисляет совместимость двух людей на основе их координат
// Комбинация трех факторов:
// 1. Общие оси (похожесть базовых паттернов) - 25% веса
// 2. Разница амплитуд (совместимость темпов жизни) - 35% веса
// 3. Синхронные шаги (совпадение важных моментов) - 40% веса
func CalculateCompatibility(coords1, coords2 core.HypercubeCoords, vectors1, vectors2 core.PersonVectors) CompatibilityResult {
	// 1. Общие оси - насколько похожи базовые координаты
	commonAxes := 0
	if coords1.X == coords2.X {
		commonAxes++
	}
	if coords1.Y == coords2.Y {
		commonAxes++
	}
	if coords1.Z == coords2.Z {
		commonAxes++
	}
	if coords1.W == coords2.W {
		commonAxes++
	}

	// 2. Разница в амплитуде векторов - насколько похожи темпы жизни
	amp1 := core.VectorAmplitude(vectors1)
	amp2 := core.VectorAmplitude(vectors2)
	ampDiff := math.Abs(amp1 - amp2)
	ampScore := 1.0 - math.Min(ampDiff/10.0, 1.0)

	// 3. Синхронность движения - сколько раз траектории сближаются
	syncSteps := 0
	pos1, pos2 := coords1, coords2
	for step := 0; step < 12; step++ {
		pos1 = core.MoveRoom(pos1, vectors1, step)
		pos2 = core.MoveRoom(pos2, vectors2, step)
		distance := core.EuclideanDistance(pos1, pos2)
		if distance < 5.0 {
			syncSteps++
		}
	}
	syncScore := float64(syncSteps) / 12.0

	// Итоговый балл (веса: оси - 25%, амплитуда - 35%, синхронность - 40%)
	score := float64(commonAxes)*0.25/4.0 + ampScore*0.35 + syncScore*0.4

	// Определяем уровень совместимости
	level := "Низкая"
	recommendations := make([]string, 0)
	
	if score >= 0.7 {
		level = "Высокая"
		recommendations = append(recommendations, 
			"✓ Высокая совместимость: хороший потенциал для здоровых отношений",
			"✓ Рекомендуется осознанная работа над индивидуальными травмами",
			"✓ Важно сохранять личные границы")
	} else if score >= 0.4 {
		level = "Средняя"
		recommendations = append(recommendations,
			"➤ Средняя совместимость: требуется осознанная работа над отношениями",
			"➤ Важно выстраивать четкие границы",
			"➤ Рекомендуется парная терапия для проработки сложностей")
	} else {
		level = "Низкая"
		recommendations = append(recommendations,
			"✗ Низкая совместимость: высокий риск деструктивных паттернов",
			"✗ Рекомендуется осторожность в сближении",
			"✗ Важно осознавать разницу в жизненных ритмах",
			"✗ Рекомендуется индивидуальная терапия перед построением отношений")
	}

	return CompatibilityResult{
		Score:           score,
		CommonAxes:      commonAxes,
		AmpDiff:         ampDiff,
		SyncSteps:       syncSteps,
		Level:           level,
		KarmicFactor:    1.0,
		Recommendations: recommendations,
	}
}

// CalculateGroupCompatibility анализирует совместимость в группе (3+ человек)
func CalculateGroupCompatibility(members []*data.ExtendedFamilyMember) GroupCompatibilityResult {
	if len(members) < 2 {
		return GroupCompatibilityResult{
			AverageScore:    0,
			ConflictScore:   0,
			StabilityScore:  0,
			GroupDynamics:   []string{"Недостаточно участников для анализа"},
			RequiredActions: []string{},
		}
	}

	totalScore := 0.0
	conflicts := 0
	stablePairs := 0
	pairs := 0

	// Анализируем все пары в группе
	for i, m1 := range members {
		for j, m2 := range members {
			if i >= j {
				continue
			}
			pairs++
			
			coords1 := core.ParseDateToCoords(m1.BirthDate)
			coords2 := core.ParseDateToCoords(m2.BirthDate)
			vectors1 := core.CalculateVectors(coords1)
			vectors2 := core.CalculateVectors(coords2)
			
			compat := CalculateCompatibility(coords1, coords2, vectors1, vectors2)
			totalScore += compat.Score
			
			if compat.Level == "Низкая" {
				conflicts++
			}
			if compat.Level == "Высокая" {
				stablePairs++
			}
		}
	}

	avgScore := totalScore / float64(pairs)
	conflictScore := float64(conflicts) / float64(pairs)
	stabilityScore := float64(stablePairs) / float64(pairs)

	// Анализ групповой динамики
	dynamics := make([]string, 0)
	actions := make([]string, 0)

	if conflictScore > 0.5 {
		dynamics = append(dynamics, "⚠️ Высокий уровень конфликтности в группе")
		actions = append(actions, "Групповая терапия", "Медиация", "Четкие границы")
	}
	if stabilityScore < 0.3 {
		dynamics = append(dynamics, "⚠️ Низкая стабильность связей")
		actions = append(actions, "Формирование новых связей", "Поиск поддерживающих фигур")
	}
	if avgScore > 0.6 {
		dynamics = append(dynamics, "✓ Группа в целом совместима")
	}
	if avgScore < 0.4 {
		dynamics = append(dynamics, "⚠️ Группа в целом несовместима - высокий риск распада")
		actions = append(actions, "Пересмотр состава группы", "Индивидуальная терапия для участников")
	}

	// Гендерный анализ
	maleCount := 0
	femaleCount := 0
	for _, m := range members {
		if m.Gender == core.GenderMale {
			maleCount++
		} else {
			femaleCount++
		}
	}
	
	if maleCount > femaleCount*2 {
		dynamics = append(dynamics, "➤ Преобладание мужчин - маскулинная динамика")
	} else if femaleCount > maleCount*2 {
		dynamics = append(dynamics, "➤ Преобладание женщин - феминная динамика")
	}

	return GroupCompatibilityResult{
		AverageScore:    avgScore,
		ConflictScore:   conflictScore,
		StabilityScore:  stabilityScore,
		GroupDynamics:   dynamics,
		RequiredActions: actions,
	}
}
