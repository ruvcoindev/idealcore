package analysis

import (
	"fmt"
	"strings"
	"time"  // <--- ЭТО НУЖНО ДОБАВИТЬ
	"github.com/ruvcoindev/idealcore/pkg/hypercube/core"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/data"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/model"
)

// AddictionType определяет тип зависимости
type AddictionType string

const (
	AddictionTypeAlcohol      AddictionType = "alcohol"       // Алкогольная зависимость
	AddictionTypeDrugs        AddictionType = "drugs"         // Наркотическая зависимость
	AddictionTypeGambling     AddictionType = "gambling"      // Игровая зависимость (лудомания)
	AddictionTypeFood         AddictionType = "food"          // Пищевая зависимость
	AddictionTypeSex          AddictionType = "sex"           // Сексуальная зависимость
	AddictionTypeWork         AddictionType = "work"          // Трудоголизм
	AddictionTypeInternet     AddictionType = "internet"      // Интернет-зависимость
	AddictionTypeGaming       AddictionType = "gaming"        // Игровая (компьютерные игры)
	AddictionTypeShopping     AddictionType = "shopping"      // Шопоголизм
	AddictionTypeRelationships AddictionType = "relationships" // Зависимость от отношений (любовная)
	AddictionTypeCodependency AddictionType = "codependency"   // Созависимость
	AddictionTypePrescription AddictionType = "prescription"   // Зависимость от лекарств
)

// AddictionSeverity определяет тяжесть зависимости
type AddictionSeverity string

const (
	AddictionSeverityMild      AddictionSeverity = "mild"      // Легкая
	AddictionSeverityModerate  AddictionSeverity = "moderate"  // Умеренная
	AddictionSeveritySevere    AddictionSeverity = "severe"    // Тяжелая
	AddictionSeverityCritical  AddictionSeverity = "critical"  // Критическая (угроза жизни)
)

// AddictionStage определяет стадию зависимости
type AddictionStage string

const (
	AddictionStageExperimental  AddictionStage = "experimental"   // Экспериментальная
	AddictionStageRegular       AddictionStage = "regular"        // Регулярное употребление
	AddictionStageRisky         AddictionStage = "risky"          // Рискованное употребление
	AddictionStageDependent     AddictionStage = "dependent"      // Зависимость
	AddictionStageCrisis        AddictionStage = "crisis"         // Кризис
	AddictionStageRecovery      AddictionStage = "recovery"       // Выздоровление
	AddictionStageRelapse       AddictionStage = "relapse"        // Срыв/рецидив
)

// AddictionPattern определяет паттерн зависимости в семье
type AddictionPattern string

const (
	AddictionPatternSingle      AddictionPattern = "single"       // Один зависимый
	AddictionPatternMultiple    AddictionPattern = "multiple"     // Несколько зависимых
	AddictionPatternGenerational AddictionPattern = "generational" // Межпоколенческая
	AddictionPatternEnabling    AddictionPattern = "enabling"     // С системой энейблеров
	AddictionPatternCoAddiction AddictionPattern = "co-addiction" // Со-зависимости (разные типы у разных членов)
)

// AddictionIncident представляет эпизод, связанный с зависимостью
type AddictionIncident struct {
	ID              string
	Date            *time.Time
	AddictionType   AddictionType
	Severity        AddictionSeverity
	Stage           AddictionStage
	PersonID        string
	Substance       string              // для веществ - название
	Amount          string              // количество
	Context         string              // контекст (один, в компании, в семье)
	Consequences    []string            // последствия
	Treatment       bool                // было ли лечение
	TreatmentType   string              // тип лечения
	Relapse         bool                // был ли это срыв
	Impact          float64             // влияние на жизнь (0-1)
}

// AddictionDynamics описывает динамику зависимостей в семье
type AddictionDynamics struct {
	FamilyID                string
	HasAddiction            bool
	AddictionTypes          map[AddictionType][]string // тип -> список ID
	PrimaryAddictID         string
	SecondaryAddicts        []string
	Enablers                []string                    // те, кто способствует/покрывает
	CoAddicts               []string                    // со-зависимые
	Pattern                 AddictionPattern
	Generational            bool                        // передается через поколения
	NormalizationLevel      float64                     // уровень нормализации (0-1)
	SecrecyLevel            float64                     // уровень секретности (0-1)
	FirstIncidentDate       *time.Time
	LastIncidentDate        *time.Time
	IncidentCount           int
	TreatmentAttempts       int
	RelapseCount            int
	FinancialImpact         string                      // финансовые последствия
	HealthImpact            string                      // последствия для здоровья
	SocialImpact            string                      // социальные последствия
}

// AddictionAnalysisResult содержит результаты анализа зависимостей
type AddictionAnalysisResult struct {
	FamilyID            string
	Dynamics            *AddictionDynamics
	Incidents           []AddictionIncident
	AddictsByType       map[AddictionType][]string
	Timeline            []AddictionIncident
	RiskLevel           string              // "low", "medium", "high", "critical"
	Recommendations     []string
	Warnings            []string
	Resources           []string
}

// AddictionHistory хранит историю зависимостей для человека
type AddictionHistory struct {
	PersonID            string
	Addictions          []AddictionIncident
	PrimaryAddiction    AddictionType
	AgeOfOnset          int                   // возраст начала
	Duration            int                   // длительность в годах
	TreatmentHistory    []string              // история лечения
	RelapseHistory      []int                 // годы срывов
	RecoveryPeriods     []string              // периоды ремиссии
	TriggerPatterns     []string              // паттерны триггеров
	CopingMechanisms    []string              // механизмы совладания
	SupportSystems      []string              // системы поддержки
	FamilyHistory       []AddictionType       // типы зависимостей в семье
}

// ExtendFamilyMemberWithAddiction расширяет члена семьи данными о зависимостях
type ExtendFamilyMemberWithAddiction struct {
	*data.ExtendedFamilyMember
	AddictionHistory    *AddictionHistory
	AddictionRiskScore  float64               // риск развития зависимости (0-1)
	ProtectiveFactors   []string              // защитные факторы
	VulnerabilityFactors []string             // факторы уязвимости
	RecoveryPotential   float64               // потенциал выздоровления (0-1)
}

// AnalyzeAddictionDynamics анализирует динамику зависимостей в семье
func AnalyzeAddictionDynamics(fs *model.FamilySystem) *AddictionAnalysisResult {
	result := &AddictionAnalysisResult{
		AddictsByType:    make(map[AddictionType][]string),
		Incidents:        make([]AddictionIncident, 0),
		Timeline:         make([]AddictionIncident, 0),
		Recommendations:  make([]string, 0),
		Warnings:         make([]string, 0),
		Resources:        make([]string, 0),
	}

	dynamics := &AddictionDynamics{
		AddictionTypes:   make(map[AddictionType][]string),
		SecondaryAddicts: make([]string, 0),
		Enablers:         make([]string, 0),
		CoAddicts:        make([]string, 0),
		HasAddiction:     false,
	}

	// Собираем все инциденты из событий
	for id, member := range fs.Members {
		for _, event := range member.Events {
			if strings.Contains(event.EventType, "addiction") || 
			   strings.Contains(event.EventType, "alcohol") ||
			   strings.Contains(event.EventType, "drug") ||
			   strings.Contains(event.EventType, "gambling") ||
			   strings.Contains(event.EventType, "relapse") ||
			   strings.Contains(event.EventType, "treatment") {
				
				addictionType := determineAddictionType(event.EventType)
				severity := determineAddictionSeverity(event.Description)
				stage := determineAddictionStage(event.EventType)
				
				incident := AddictionIncident{
					ID:            fmt.Sprintf("addict_%s_%d", id, event.Date.Unix()),
					Date:          &event.Date,
					AddictionType: addictionType,
					Severity:      severity,
					Stage:         stage,
					PersonID:      id,
					Context:       event.Description,
					Treatment:     strings.Contains(event.EventType, "treatment"),
					Relapse:       strings.Contains(event.EventType, "relapse"),
					Impact:        event.Impact,
				}
				
				result.Incidents = append(result.Incidents, incident)
				result.Timeline = append(result.Timeline, incident)
				
				// Обновляем статистику
				dynamics.IncidentCount++
				dynamics.HasAddiction = true
				
				if dynamics.FirstIncidentDate == nil || event.Date.Before(*dynamics.FirstIncidentDate) {
					dynamics.FirstIncidentDate = &event.Date
				}
				if dynamics.LastIncidentDate == nil || event.Date.After(*dynamics.LastIncidentDate) {
					dynamics.LastIncidentDate = &event.Date
				}
				
				if incident.Treatment {
					dynamics.TreatmentAttempts++
				}
				if incident.Relapse {
					dynamics.RelapseCount++
				}
				
				// Добавляем в карты по типам
				dynamics.AddictionTypes[addictionType] = append(dynamics.AddictionTypes[addictionType], id)
				result.AddictsByType[addictionType] = append(result.AddictsByType[addictionType], id)
			}
		}
	}

	if !dynamics.HasAddiction {
		return result
	}

	// Определяем основного зависимого
	addictCount := make(map[string]int)
	for _, inc := range result.Incidents {
		addictCount[inc.PersonID]++
	}
	
	maxCount := 0
	for id, count := range addictCount {
		if count > maxCount {
			maxCount = count
			dynamics.PrimaryAddictID = id
		}
	}

	// Находим вторичных зависимых
	for id, count := range addictCount {
		if id != dynamics.PrimaryAddictID && count > 1 {
			dynamics.SecondaryAddicts = append(dynamics.SecondaryAddicts, id)
		}
	}

// Ищем энейблеров (тех, кто способствует/покрывает)
for id, member := range fs.Members {
    if id == dynamics.PrimaryAddictID {
        continue
    }

    // Проверяем, есть ли признаки созависимости
    coords := core.ParseDateToCoords(member.BirthDate)
    vectors := core.CalculateVectors(coords)
    amp := core.VectorAmplitude(vectors)

    // Низкая амплитуда может указывать на подавление себя (созависимость)
    if amp < 2.5 {
        // Проверяем, связан ли с основным зависимым
        if member.HasPartner(dynamics.PrimaryAddictID) ||
           ContainsString(member.Children, dynamics.PrimaryAddictID) ||
           ContainsString(member.Parents, dynamics.PrimaryAddictID) {
            dynamics.Enablers = append(dynamics.Enablers, id)
        }
    }
}

// Определяем паттерн (ЭТОТ БЛОК ДОЛЖЕН БЫТЬ ПОСЛЕ ЦИКЛА)
if len(dynamics.AddictionTypes) > 2 {
    dynamics.Pattern = AddictionPatternMultiple
} else if checkGenerationalAddiction(fs, result.Incidents) {
    dynamics.Pattern = AddictionPatternGenerational
    dynamics.Generational = true
} else if len(dynamics.Enablers) > 0 {
    dynamics.Pattern = AddictionPatternEnabling
} else {
    dynamics.Pattern = AddictionPatternSingle
}

// Оцениваем уровень нормализации
dynamics.NormalizationLevel = calculateNormalizationLevel(fs, dynamics.PrimaryAddictID)

// Оцениваем уровень секретности
dynamics.SecrecyLevel = calculateSecrecyLevel(result.Incidents)

// Оцениваем уровень риска
riskScore := 0.0
if dynamics.IncidentCount > 20 {
    riskScore += 0.4
}
if dynamics.RelapseCount > 3 {
    riskScore += 0.3

	}
	if dynamics.Pattern == AddictionPatternGenerational {
		riskScore += 0.4
	}
	if dynamics.NormalizationLevel > 0.7 {
		riskScore += 0.3
	}
	if dynamics.SecrecyLevel > 0.8 {
		riskScore += 0.2
	}
	
	switch {
	case riskScore > 1.2:
		result.RiskLevel = "critical"
	case riskScore > 0.8:
		result.RiskLevel = "high"
	case riskScore > 0.4:
		result.RiskLevel = "medium"
	default:
		result.RiskLevel = "low"
	}

	// Формируем предупреждения
	if dynamics.IncidentCount > 30 {
		result.Warnings = append(result.Warnings,
			"⚠️ Хроническая зависимость - требуется интенсивное лечение")
	}
	if dynamics.RelapseCount > 5 {
		result.Warnings = append(result.Warnings,
			"⚠️ Множественные срывы - необходима смена подхода к лечению")
	}
	if dynamics.Pattern == AddictionPatternGenerational {
		result.Warnings = append(result.Warnings,
			"⚠️ Межпоколенческая зависимость - требуется разрыв цикла")
	}
	if len(dynamics.Enablers) > 2 {
		result.Warnings = append(result.Warnings,
			"⚠️ Сильная система энейблеров - семейная система поддерживает зависимость")
	}

	// Рекомендации
	result.Recommendations = generateAddictionRecommendations(dynamics, result.RiskLevel)
	
		// Ресурсы помощи
	result.Resources = []string{
		"📞 В России: бесплатная горячая линия помощи зависимым 8-800-700-50-50",
		"📞 Анонимные Алкоголики: 8-800-555-55-55",
		"📞 Нарколог: 8-800-200-0-200",
		"🌐 Сайт: https://www.help-addiction.ru",
	}

	return result
}

func determineAddictionType(eventType string) AddictionType {
	eventType = strings.ToLower(eventType)
	
	if strings.Contains(eventType, "alcohol") {
		return AddictionTypeAlcohol
	}
	if strings.Contains(eventType, "drug") {
		return AddictionTypeDrugs
	}
	if strings.Contains(eventType, "gambling") {
		return AddictionTypeGambling
	}
	return AddictionTypeAlcohol // по умолчанию
}

func determineAddictionSeverity(description string) AddictionSeverity {
	desc := strings.ToLower(description)
	
	if strings.Contains(desc, "critical") || strings.Contains(desc, "life-threatening") {
		return AddictionSeverityCritical
	}
	if strings.Contains(desc, "severe") {
		return AddictionSeveritySevere
	}
	if strings.Contains(desc, "moderate") {
		return AddictionSeverityModerate
	}
	return AddictionSeverityMild
}

func determineAddictionStage(eventType string) AddictionStage {
	eventType = strings.ToLower(eventType)
	
	if strings.Contains(eventType, "relapse") {
		return AddictionStageRelapse
	}
	if strings.Contains(eventType, "treatment") || strings.Contains(eventType, "recovery") {
		return AddictionStageRecovery
	}
	if strings.Contains(eventType, "crisis") {
		return AddictionStageCrisis
	}
	return AddictionStageRegular
}

func checkGenerationalAddiction(fs *model.FamilySystem, incidents []AddictionIncident) bool {
	// Упрощенная реализация - в реальном коде здесь будет анализ поколений
	return false
}

func calculateNormalizationLevel(fs *model.FamilySystem, primaryAddictID string) float64 {
	// Упрощенная реализация
	return 0.5
}

func calculateSecrecyLevel(incidents []AddictionIncident) float64 {
	// Упрощенная реализация
	return 0.5
}

func generateAddictionRecommendations(dynamics *AddictionDynamics, riskLevel string) []string {
	recs := []string{
		"Рекомендации по работе с зависимостью:",
	}
	
	switch riskLevel {
	case "critical":
		recs = append(recs,
			"• Немедленное обращение к наркологу",
			"• Стационарное лечение",
			"• Детоксикация")
	case "high":
		recs = append(recs,
			"• Консультация нарколога",
			"• Регулярная терапия",
			"• Группы поддержки")
	default:
		recs = append(recs,
			"• Профилактика",
			"• Образовательные программы")
	}
	
	return recs
}
// ContainsString проверяет наличие строки в слайсе
func ContainsString(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

