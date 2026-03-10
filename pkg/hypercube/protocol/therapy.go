package protocol

import (
	"fmt"
	"github.com/ruvcoindev/idealcore/hypercube/core"
	"github.com/ruvcoindev/idealcore/hypercube/data"
)

// TherapyRecommendation содержит рекомендацию по терапии
type TherapyRecommendation struct {
	Type        string   // Тип: "individual", "group", "body", "art", etc.
	Description string   // Описание
	Priority    int      // Приоритет (1-5, 1 - наивысший)
	ForRoles    []core.TraumaRole // Для каких ролей подходит
	Duration    string   // Рекомендуемая длительность
}

// TherapyPlan содержит полный план терапии для человека
type TherapyPlan struct {
	PersonID        string
	PersonName      string
	PrimaryRole     core.TraumaRole
	Recommendations []TherapyRecommendation
	UrgentActions   []string
	LongTermGoals   []string
	Contraindications []string
}

// GenerateTherapyPlan создает план терапии на основе травматической роли
func GenerateTherapyPlan(member *data.ExtendedFamilyMember, role *core.IdentifiedTraumaRoles) TherapyPlan {
	plan := TherapyPlan{
		PersonID:      member.ID,
		PersonName:    member.Name,
		PrimaryRole:   role.PrimaryRole,
		Recommendations: make([]TherapyRecommendation, 0),
		UrgentActions:   make([]string, 0),
		LongTermGoals:   make([]string, 0),
		Contraindications: make([]string, 0),
	}
	
	// Общие рекомендации для всех
	plan.Recommendations = append(plan.Recommendations,
		TherapyRecommendation{
			Type:        "individual",
			Description: "Индивидуальная психотерапия",
			Priority:    1,
			ForRoles:    []core.TraumaRole{},
			Duration:    "2-5 лет",
		},
		TherapyRecommendation{
			Type:        "body",
			Description: "Телесно-ориентированная терапия",
			Priority:    2,
			ForRoles:    []core.TraumaRole{},
			Duration:    "1-3 года",
		})
	
	// Специфические рекомендации в зависимости от роли
	switch role.PrimaryRole {
	case core.TraumaRoleScapegoat:
		plan.Recommendations = append(plan.Recommendations,
			TherapyRecommendation{
				Type:        "group",
				Description: "Группы поддержки для взрослых детей из дисфункциональных семей",
				Priority:    2,
				ForRoles:    []core.TraumaRole{core.TraumaRoleScapegoat},
				Duration:    "1-2 года",
			},
			TherapyRecommendation{
				Type:        "anger",
				Description: "Терапия гнева и работа с обидой",
				Priority:    3,
				ForRoles:    []core.TraumaRole{core.TraumaRoleScapegoat},
				Duration:    "6-12 месяцев",
			})
		plan.UrgentActions = append(plan.UrgentActions,
			"Прекратить контакты с обесценивающими людьми",
			"Начать вести дневник самоценности")
		plan.LongTermGoals = append(plan.LongTermGoals,
			"Восстановление чувства собственной ценности",
			"Умение выстраивать здоровые границы")
		
	case core.TraumaRoleGoldenChild:
		plan.Recommendations = append(plan.Recommendations,
			TherapyRecommendation{
				Type:        "perfectionism",
				Description: "Работа с перфекционизмом и синдромом самозванца",
				Priority:    2,
				ForRoles:    []core.TraumaRole{core.TraumaRoleGoldenChild},
				Duration:    "1-2 года",
			})
		plan.UrgentActions = append(plan.UrgentActions,
			"Разрешить себе ошибаться",
			"Начать говорить 'нет' ожиданиям других")
		plan.LongTermGoals = append(plan.LongTermGoals,
			"Принятие своего несовершенства",
			"Поиск собственных желаний, отделенных от родительских")
		
	case core.TraumaRoleLostChild:
		plan.Recommendations = append(plan.Recommendations,
			TherapyRecommendation{
				Type:        "activation",
				Description: "Активирующая терапия, работа с телом и голосом",
				Priority:    1,
				ForRoles:    []core.TraumaRole{core.TraumaRoleLostChild},
				Duration:    "2-3 года",
			},
			TherapyRecommendation{
				Type:        "social",
				Description: "Постепенная социальная активация, группы по интересам",
				Priority:    3,
				ForRoles:    []core.TraumaRole{core.TraumaRoleLostChild},
				Duration:    "постоянно",
			})
		plan.UrgentActions = append(plan.UrgentActions,
			"Начать замечать свои потребности",
			"Практиковать маленькие шаги к контакту")
		plan.LongTermGoals = append(plan.LongTermGoals,
			"Способность заявлять о себе",
			"Умение получать удовольствие от жизни")
		
	case core.TraumaRoleInvisible:
		plan.Recommendations = append(plan.Recommendations,
			TherapyRecommendation{
				Type:        "visibility",
				Description: "Терапия, фокусированная на видимости и праве занимать место",
				Priority:    1,
				ForRoles:    []core.TraumaRole{core.TraumaRoleInvisible},
				Duration:    "2-3 года",
			})
		plan.UrgentActions = append(plan.UrgentActions,
			"Начать занимать больше места физически",
			"Практиковать выражение мнения в безопасной среде")
		plan.LongTermGoals = append(plan.LongTermGoals,
			"Способность быть видимым и слышимым",
			"Комфорт в центре внимания")
		
	case core.TraumaRoleParentified:
		plan.Recommendations = append(plan.Recommendations,
			TherapyRecommendation{
				Type:        "childhood",
				Description: "Терапия для возвращения себе детства и права на беззаботность",
				Priority:    1,
				ForRoles:    []core.TraumaRole{core.TraumaRoleParentified},
				Duration:    "2-4 года",
			})
		plan.UrgentActions = append(plan.UrgentActions,
			"Перестать быть 'спасателем' для других",
			"Начать получать заботу, а не только давать")
		plan.LongTermGoals = append(plan.LongTermGoals,
			"Способность принимать помощь",
			"Умение быть уязвимым")
		
	case core.TraumaRoleEmotionalSpouse:
		plan.Recommendations = append(plan.Recommendations,
			TherapyRecommendation{
				Type:        "separation",
				Description: "Срочная терапия для разрыва симбиотической связи с родителем",
				Priority:    1,
				ForRoles:    []core.TraumaRole{core.TraumaRoleEmotionalSpouse},
				Duration:    "3-5 лет",
			},
			TherapyRecommendation{
				Type:        "couple",
				Description: "Терапия для построения здоровых партнерских отношений",
				Priority:    2,
				ForRoles:    []core.TraumaRole{core.TraumaRoleEmotionalSpouse},
				Duration:    "2-3 года",
			})
		plan.UrgentActions = append(plan.UrgentActions,
			"Минимизировать контакты с родителем, с которым симбиоз",
			"Начать строить отношения вне семьи")
		plan.LongTermGoals = append(plan.LongTermGoals,
			"Здоровые романтические отношения",
			"Отделение своей идентичности от родительской")
		
	case core.TraumaRoleReplacement:
		plan.Recommendations = append(plan.Recommendations,
			TherapyRecommendation{
				Type:        "identity",
				Description: "Терапия идентичности - отделение от образа потерянного",
				Priority:    1,
				ForRoles:    []core.TraumaRole{core.TraumaRoleReplacement},
				Duration:    "3-5 лет",
			},
			TherapyRecommendation{
				Type:        "grief",
				Description: "Работа с виной выжившего и правом на жизнь",
				Priority:    2,
				ForRoles:    []core.TraumaRole{core.TraumaRoleReplacement},
				Duration:    "2-3 года",
			})
		plan.UrgentActions = append(plan.UrgentActions,
			"Ритуал прощания с потерянным сиблингом",
			"Начать исследовать свою уникальность")
		plan.LongTermGoals = append(plan.LongTermGoals,
			"Присвоение права на собственную жизнь",
			"Освобождение от роли 'замены'")
	}
	
	// Противопоказания
	if role.PrimaryRole == core.TraumaRoleLostChild || role.PrimaryRole == core.TraumaRoleInvisible {
		plan.Contraindications = append(plan.Contraindications,
			"Интенсивные групповые терапии без предварительной подготовки",
			"Конфронтационные методы")
	}
	
	return plan
}

// FormatTherapyPlan форматирует план терапии для вывода
func FormatTherapyPlan(plan TherapyPlan) string {
	result := fmt.Sprintf("🧑‍⚕️ ПЛАН ТЕРАПИИ ДЛЯ %s\n", plan.PersonName)
	result += fmt.Sprintf("Основная роль: %s\n\n", plan.PrimaryRole)
	
	result += "СРОЧНЫЕ ДЕЙСТВИЯ:\n"
	for i, action := range plan.UrgentActions {
		result += fmt.Sprintf("%d. %s\n", i+1, action)
	}
	
	result += "\nРЕКОМЕНДАЦИИ ПО ТЕРАПИИ:\n"
	for i, rec := range plan.Recommendations {
		result += fmt.Sprintf("%d. [%s] %s (приоритет %d, %s)\n",
			i+1, rec.Type, rec.Description, rec.Priority, rec.Duration)
	}
	
	result += "\nДОЛГОСРОЧНЫЕ ЦЕЛИ:\n"
	for i, goal := range plan.LongTermGoals {
		result += fmt.Sprintf("%d. %s\n", i+1, goal)
	}
	
	if len(plan.Contraindications) > 0 {
		result += "\nПРОТИВОПОКАЗАНИЯ:\n"
		for i, contra := range plan.Contraindications {
			result += fmt.Sprintf("%d. %s\n", i+1, contra)
		}
	}
	
	return result
}
