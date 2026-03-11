package analysis

import (
	"math"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/core"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/data"
)

// ShameAnalysisResult содержит результаты анализа стыда
type ShameAnalysisResult struct {
	TotalScore      float64              // Общий уровень стыда (0-100)
	PrimaryPattern  string               // Основной паттерн: "burning", "freezing", "fleeing"
	BodyZones       map[string]float64   // Зоны тела и уровень напряжения
	Triggers        []string             // Основные триггеры
	Recommendations []string             // Рекомендации по работе со стыдом
}

// CalculateBodyShame вычисляет уровень телесного стыда на основе координат и роли
// Стыд - ключевой компонент травмы, особенно в нарциссических семьях
// Параметры:
// - coords: координаты в гиперкубе
// - traumaRole: выявленная травматическая роль
// - ageAtEvent: возраст для учета накопления
// - member: расширенные данные (для учета замещения, потерь и т.д.)
func CalculateBodyShame(coords core.HypercubeCoords, traumaRole core.TraumaRole, ageAtEvent int, member *data.ExtendedFamilyMember) float64 {
	// Базовый стыд из цифр координат (0-90)
	digits := core.ExtractDigits(coords.X + coords.Y + coords.Z)
	baseShame := float64(digits[0] * 10)
	
	// Коэффициенты для разных травматических ролей
	// Основано на клинических наблюдениях: разные роли несут разный уровень стыда
	roleMultiplier := map[core.TraumaRole]float64{
		core.TraumaRoleScapegoat:      1.8, // Козел отпущения - высокий стыд
		core.TraumaRoleGoldenChild:    1.3, // Любимчик - стыд за несовершенство
		core.TraumaRoleLostChild:      1.6, // Потерянный - стыд за само существование
		core.TraumaRoleInvisible:      1.5, // Невидимка - стыд за потребности
		core.TraumaRoleGlassChild:     1.4, // Стеклянный - стыд за то, что "нормальный"
		core.TraumaRoleParentified:    1.7, // Родифицированный - стыд за неспособность "спасти"
		core.TraumaRoleMascot:         1.2, // Шут - стыд, скрытый юмором
		core.TraumaRoleEmotionalSpouse: 1.9, // Эмоциональный супруг - глубокий стыд за инцестуозную связь
		core.TraumaRoleTruthTeller:    1.5, // Говорящий правду - стыд за "предательство" семьи
		core.TraumaRoleShadow:         1.8, // Тень - стыд за подавленные эмоции
		core.TraumaRoleHealer:         1.3, // Целитель - стыд за неспособность исцелить
		core.TraumaRoleReplacement:    2.2, // Замещающий ребенок - максимальный стыд
		core.TraumaRoleGhost:          1.0, // Призрак - не чувствует стыда (не жив)
		core.TraumaRoleAncestor:       1.5, // Предок - родовой стыд
	}
	
	multiplier := roleMultiplier[traumaRole]
	if multiplier == 0 {
		multiplier = 1.0
	}

	// Накопление стыда с возрастом (экспоненциально, если не прорабатывать)
	accumulated := baseShame * multiplier * math.Exp(float64(ageAtEvent)/20.0)
	
	// Дополнительные факторы из расширенных данных
	if member != nil {
		// Замещающие дети несут дополнительный стыд
		if member.IsReplacement && member.ReplacementData != nil {
			accumulated *= (1 + member.ReplacementData.BurdenLevel)
		}
		
		// Вина выжившего добавляет стыд
		accumulated += member.SurvivorGuilt * 50
		
		// Аборты и потери
		for _, abortion := range member.Abortions {
			accumulated += abortion.EmotionalImpact * 30
		}
		
		// Измены
		for _, infidelity := range member.Infidelities {
			accumulated += infidelity.Impact * 20
		}
		
		// Заболевания
		for _, disease := range member.Diseases {
			accumulated += disease.Severity * 25
		}
		
		// Кармические долги
		for _, debt := range member.KarmicDebts {
			if !debt.IsResolved {
				accumulated += debt.Weight * 40
			}
		}
	}
	
	// Ограничиваем до 250 (максимальный уровень)
	return math.Min(accumulated, 250.0)
}

// AnalyzeShamePattern определяет паттерн стыда по векторам и роли
func AnalyzeShamePattern(vectors core.PersonVectors, traumaRole core.TraumaRole) *data.ShamePattern {
	amp := core.VectorAmplitude(vectors)
	
	pattern := &data.ShamePattern{
		Score:    0,
		Triggers: make([]string, 0),
	}
	
	// Определяем паттерн по амплитуде
	if amp > 6.0 {
		// Высокая амплитуда - активный стыд, "сжигание"
		pattern.Pattern = "burning"
		pattern.BodyZone = "whole"
		pattern.Score = amp * 10
	} else if amp < 3.0 {
		// Низкая амплитуда - заморозка, оцепенение
		pattern.Pattern = "freezing"
		pattern.BodyZone = "gut"
		pattern.Score = 100 - amp*10
	} else {
		// Средняя амплитуда - бегство, избегание
		pattern.Pattern = "fleeing"
		pattern.BodyZone = "heart"
		pattern.Score = 50
	}
	
	// Добавляем триггеры в зависимости от роли
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

// GenerateShameReport генерирует отчет по стыду для человека
func GenerateShameReport(member *data.ExtendedFamilyMember, traumaRole core.TraumaRole, vectors core.PersonVectors) ShameAnalysisResult {
	coords := core.ParseDateToCoords(member.BirthDate)
	age := member.GetAge()
	
	totalShame := CalculateBodyShame(coords, traumaRole, age, member)
	pattern := AnalyzeShamePattern(vectors, traumaRole)
	
	// Анализ зон тела (упрощенно)
	bodyZones := map[string]float64{
		"голова":    totalShame * 0.2,
		"горло":     totalShame * 0.15,
		"грудь":     totalShame * 0.25,
		"живот":     totalShame * 0.3,
		"таз":       totalShame * 0.1,
	}
	
	// Рекомендации в зависимости от уровня стыда
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
	
	// Добавляем рекомендации по паттерну
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
