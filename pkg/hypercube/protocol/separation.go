package protocol

import (
	"fmt"
	"time"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/core"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/data"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/analysis"
)

// SeparationStage описывает этап разделения
type SeparationStage struct {
	Year        int      // Год разделения
	Name        string   // Название этапа
	Symptoms    []string // Симптомы, которые могут возникнуть
	SurvivalRate float64 // Вероятность успешного прохождения (0-1)
}

// SeparationProtocol содержит план безопасного разделения для деструктивных отношений
type SeparationProtocol struct {
	RequiredYears       int                // Минимальное количество лет разделения
	Stages              []SeparationStage  // Этапы разделения
	RelapseProbability  float64            // Вероятность рецидива (0-1)
	SuccessProbability  float64            // Вероятность успешного исцеления (0-1)
	ContactRules        []string           // Правила контакта
	TherapyRecommendations []string        // Рекомендации по терапии
	Warning             string             // Предупреждение (если есть)
}

// GenerateSeparationProtocol создает протокол разделения для двух людей
func GenerateSeparationProtocol(p1, p2 *data.ExtendedFamilyMember, relType core.RelationshipType, 
	compat analysis.CompatibilityResult) SeparationProtocol {
	
	// Базовые параметры
	requiredYears := 3
	if relType == core.RelationshipRomantic || relType == core.RelationshipHusbandWife {
		requiredYears = 5 // Романтические отношения требуют больше времени
	}
	if relType == core.RelationshipMotherSon || relType == core.RelationshipMotherDaughter {
		requiredYears = 4 // Родительско-детские отношения
	}
	
	// Корректировка на основе совместимости
	if compat.Score < 0.3 {
		requiredYears += 1 // Очень низкая совместимость - сложнее разрыв
	}
	
	// Учет травматических ролей
	if p1.IsReplacement || p2.IsReplacement {
		requiredYears += 1 // Замещающие дети требуют больше времени
	}
	
	// Вероятность рецидива
	relapseProb := 0.5 // Базовая
	if compat.Score < 0.4 {
		relapseProb += 0.2
	}
	if relType == core.RelationshipMotherSon {
		relapseProb += 0.2 // Эмоциональный инцест особенно трудно разорвать
	}
	if p1.IsReplacement && p2.IsReplacement {
		relapseProb += 0.2 // Два замещающих - взаимная травма
	}
	if relapseProb > 0.95 {
		relapseProb = 0.95
	}
	
	successProb := 1.0 - relapseProb
	
	// Правила контакта
	contactRules := []string{
		"❌ Полное отсутствие контакта (включая соцсети и мессенджеры)",
		"❌ Не обсуждать с общими знакомыми",
		"❌ Избегать мест возможных встреч",
		"✅ При случайной встрече - минимальное общение и уход",
		fmt.Sprintf("⏱ Минимальный срок без контакта: %d года", requiredYears),
	}
	
	// Рекомендации по терапии
	therapyRecs := []string{
		"🧑‍⚕️ Индивидуальная терапия для каждого участника",
		"👥 Группы поддержки",
		"📝 Дневник чувств и прогресса",
	}
	
	if p1.IsReplacement || p2.IsReplacement {
		therapyRecs = append(therapyRecs,
			"🔄 Терапия идентичности для замещающих детей",
			"💔 Работа с виной выжившего")
	}
	
	// Стадии разделения
	stages := []SeparationStage{
		{
			Year: 1,
			Name: "Телесная ломка и острая боль",
			Symptoms: []string{
				"Острое чувство стыда и вины",
				"Желание восстановить контакт любой ценой",
				"Физическая боль в теле, психосоматика",
				"Кошмары с участием другого",
				"Активация детских травм",
				"Бессонница, потеря аппетита",
			},
			SurvivalRate: 0.4,
		},
		{
			Year: 2,
			Name: "Осознание и сепарация",
			Symptoms: []string{
				"Понимание роли другого в своей жизни",
				"Осознание собственных паттернов",
				"Злость и гнев вместо боли",
				"Попытки найти замену",
				"Укрепление границ",
			},
			SurvivalRate: 0.6,
		},
		{
			Year: 3,
			Name: "Интеграция и исцеление",
			Symptoms: []string{
				"Способность быть одному без боли",
				"Принятие опыта",
				"Благодарность за уроки",
				"Готовность к новым здоровым отношениям",
				"Восстановление целостности",
			},
			SurvivalRate: 0.8,
		},
	}
	
	// Добавляем дополнительные годы при необходимости
	if requiredYears > 3 {
		stages = append(stages, SeparationStage{
			Year: 4,
			Name: "Углубленная проработка",
			Symptoms: []string{
				"Работа с родовыми сценариями",
				"Интеграция теневых аспектов",
				"Формирование новой идентичности",
			},
			SurvivalRate: 0.9,
		})
	}
	
	if requiredYears > 4 {
		stages = append(stages, SeparationStage{
			Year: 5,
			Name: "Завершение и трансформация",
			Symptoms: []string{
				"Полное отделение",
				"Новая жизнь",
				"Возможность безопасного контакта (если нужно)",
			},
			SurvivalRate: 0.95,
		})
	}
	
	// Предупреждение для критических случаев
	warning := ""
	if successProb < 0.4 {
		warning = "⚠️ КРИТИЧЕСКИ ВЫСОКИЙ РИСК РЕЦИДИВА. Требуется интенсивная терапия и полная изоляция."
	}
	
	return SeparationProtocol{
		RequiredYears:       requiredYears,
		Stages:              stages,
		RelapseProbability:  relapseProb,
		SuccessProbability:  successProb,
		ContactRules:        contactRules,
		TherapyRecommendations: therapyRecs,
		Warning:             warning,
	}
}

// GenerateHealingProtocol создает протокол исцеления для одного человека
func GenerateHealingProtocol(member *data.ExtendedFamilyMember, traumaRole core.TraumaRole) []string {
	recommendations := []string{
		fmt.Sprintf("🧑‍⚕️ Протокол исцеления для %s", member.Name),
		"",
		"🎯 Основные направления:",
	}
	
	switch traumaRole {
	case core.TraumaRoleScapegoat:
		recommendations = append(recommendations,
			"   • Работа с чувством стыда и вины",
			"   • Отделение от родительских проекций",
			"   • Восстановление самоценности",
			"   • Терапия гнева")
	case core.TraumaRoleGoldenChild:
		recommendations = append(recommendations,
			"   • Работа с перфекционизмом",
			"   • Принятие права на ошибку",
			"   • Отделение от ожиданий",
			"   • Поиск собственных желаний")
	case core.TraumaRoleLostChild:
		recommendations = append(recommendations,
			"   • Активация эмоциональной жизни",
			"   • Поиск голоса и места",
			"   • Работа с телом",
			"   • Постепенное расширение контактов")
	case core.TraumaRoleReplacement:
		recommendations = append(recommendations,
			"   • Терапия идентичности",
			"   • Работа с виной выжившего",
			"   • Ритуал прощания с потерянным",
			"   • Присвоение права на собственную жизнь")
	default:
		recommendations = append(recommendations,
			"   • Индивидуальная терапия",
			"   • Работа с телом",
			"   • Группы поддержки")
	}
	
	recommendations = append(recommendations,
		"",
		"📚 Рекомендуемые практики:",
		"   • Дневник чувств",
		"   • Медитация и майндфулнес",
		"   • Йога или цигун",
		"   • Арт-терапия",
		"",
		"⏱ Ожидаемая длительность: 2-5 лет")
	
	return recommendations
}
