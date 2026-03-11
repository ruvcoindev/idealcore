package analysis

import (
	"math"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/core"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/data"
)

// extractDigits извлекает три цифры числа (локальная версия)
func extractDigits(n int32) []int32 {
	n = int32(math.Abs(float64(n)))
	return []int32{
		(n / 100) % 10,
		(n / 10) % 10,
		n % 10,
	}
}

// ShameAnalysisResult содержит результаты анализа стыда
type ShameAnalysisResult struct {
	TotalScore      float64
	PrimaryPattern  string
	BodyZones       map[string]float64
	Triggers        []string
	Recommendations []string
}

// CalculateBodyShame вычисляет уровень телесного стыда
func CalculateBodyShame(coords core.HypercubeCoords, traumaRole core.TraumaRole, ageAtEvent int, member *data.ExtendedFamilyMember) float64 {
	// Используем XOR для более стабильной "суммы"
	sum := coords.X ^ coords.Y ^ coords.Z
	digits := extractDigits(sum)
	baseShame := float64(digits[0] * 10)

	// Коэффициенты для ролей
	roleMultiplier := map[core.TraumaRole]float64{
		core.TraumaRoleScapegoat:       1.8,
		core.TraumaRoleGoldenChild:     1.3,
		core.TraumaRoleLostChild:       1.6,
		core.TraumaRoleInvisible:       1.5,
		core.TraumaRoleGlassChild:      1.4,
		core.TraumaRoleParentified:     1.7,
		core.TraumaRoleMascot:          1.2,
		core.TraumaRoleEmotionalSpouse: 1.9,
		core.TraumaRoleTruthTeller:     1.5,
		core.TraumaRoleShadow:          1.8,
		core.TraumaRoleHealer:          1.3,
		core.TraumaRoleReplacement:     2.2,
		core.TraumaRoleGhost:           0.5, // Призрак чувствует меньше
		core.TraumaRoleAncestor:        1.5,
	}[traumaRole]

	if roleMultiplier == 0 {
		roleMultiplier = 1.0
	}

	// Базовый расчёт
	accumulated := baseShame * roleMultiplier * math.Exp(float64(ageAtEvent)/20.0)

	// Мультипликативные бонусы за события
	if member != nil {
		eventBonus := 0.0

		if member.IsReplacement && member.ReplacementData != nil {
			eventBonus += member.ReplacementData.BurdenLevel
		}

		eventBonus += member.SurvivorGuilt * 0.4

		// Учитываем количество абортов как фактор
if member.AbortionCount > 0 {
    eventBonus += float64(member.AbortionCount) * 0.15
}

		for _, i := range member.Infidelities {
			eventBonus += i.Impact * 0.2
		}

		for _, d := range member.Diseases {
			eventBonus += d.Severity * 0.25
		}

		for _, debt := range member.KarmicDebts {
			if !debt.IsResolved {
				eventBonus += debt.Weight * 0.5
			}
		}

		// Множитель вместо сложения
		accumulated *= (1 + eventBonus)
	}

	return math.Min(accumulated, 250.0)
}

// AnalyzeShamePattern определяет паттерн стыда
func AnalyzeShamePattern(vectors core.PersonVectors, traumaRole core.TraumaRole) *data.ShamePattern {
	amp := core.VectorAmplitude(vectors)
	
	pattern := &data.ShamePattern{
		Score:    0,
		Triggers: make([]string, 0),
	}
	
	if amp > 6.0 {
		pattern.Pattern = "burning"
		pattern.BodyZone = "whole"
		pattern.Score = amp * 10
	} else if amp < 3.0 {
		pattern.Pattern = "freezing"
		pattern.BodyZone = "gut"
		pattern.Score = 100 - amp*10
	} else {
		pattern.Pattern = "fleeing"
		pattern.BodyZone = "heart"
		pattern.Score = 50
	}
	
	switch traumaRole {
	case core.TraumaRoleScapegoat:
		pattern.Triggers = []string{"критика", "обвинения", "конфликты", "публичные ситуации"}
	case core.TraumaRoleGoldenChild:
		pattern.Triggers = []string{"неудача", "критика", "сравнение с другими", "несовершенство"}
	case core.TraumaRoleLostChild:
		pattern.Triggers = []string{"внимание к себе", "публичность", "выражение потребностей"}
	case core.TraumaRoleReplacement:
		pattern.Triggers = []string{"сравнение с умершим", "дни рождения", "семейные праздники"}
	default:
		pattern.Triggers = []string{"близкие отношения", "конфликты", "семейные события"}
	}
	
	return pattern
}

// GenerateShameReport генерирует отчет по стыду
func GenerateShameReport(member *data.ExtendedFamilyMember, traumaRole core.TraumaRole, vectors core.PersonVectors) ShameAnalysisResult {
	coords := core.ParseDateToCoords(member.BirthDate)
	age := member.GetAge()
	
	totalShame := CalculateBodyShame(coords, traumaRole, age, member)
	pattern := AnalyzeShamePattern(vectors, traumaRole)
	
	bodyZones := map[string]float64{
		"голова": totalShame * 0.2,
		"горло":  totalShame * 0.15,
		"грудь":  totalShame * 0.25,
		"живот":  totalShame * 0.3,
		"таз":    totalShame * 0.1,
	}
	
	recommendations := make([]string, 0)
	
	if totalShame > 150 {
		recommendations = append(recommendations,
			"🔴 Критический уровень стыда - требуется срочная терапия",
			"   • Телесно-ориентированная терапия",
			"   • Работа с внутренним критиком",
			"   • Практики заземления")
	} else if totalShame > 80 {
		recommendations = append(recommendations,
			"🟠 Высокий уровень стыда - рекомендуется терапия",
			"   • Психотерапия, фокусированная на стыде",
			"   • Йога/цигун для работы с телом",
			"   • Группы поддержки")
	} else if totalShame > 40 {
		recommendations = append(recommendations,
			"🟡 Средний уровень стыда - возможна самостоятельная работа",
			"   • Дневник чувств",
			"   • Медитация",
			"   • Работа с уязвимостью")
	} else {
		recommendations = append(recommendations,
			"🟢 Низкий уровень стыда - профилактика",
			"   • Поддерживающие практики",
			"   • Осознанность")
	}
	
	switch pattern.Pattern {
	case "burning":
		recommendations = append(recommendations,
			"🔥 Паттерн 'сжигание' - стыд проявляется как жар, гнев",
			"   • Работа с гневом",
			"   • Охлаждающие практики")
	case "freezing":
		recommendations = append(recommendations,
			"❄️ Паттерн 'заморозка' - стыд вызывает оцепенение",
			"   • Разогревающие практики",
			"   • Постепенное возвращение чувствительности")
	case "fleeing":
		recommendations = append(recommendations,
			"🏃 Паттерн 'бегство' - стыд вызывает избегание",
			"   • Практики заземления",
			"   • Постепенное приближение к триггерам")
	}
	
	return ShameAnalysisResult{
		TotalScore:      totalShame,
		PrimaryPattern:  pattern.Pattern,
		BodyZones:       bodyZones,
		Triggers:        pattern.Triggers,
		Recommendations: recommendations,
	}
}
