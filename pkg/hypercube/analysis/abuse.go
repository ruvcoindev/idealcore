package analysis

import (
	"fmt"
	"strings"
	"time"  // <--- ЭТО НУЖНО ДОБАВИТЬ
	"github.com/ruvcoindev/idealcore/pkg/hypercube/core"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/data"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/model"
)

// AbuseType определяет тип насилия
type AbuseType string

const (
	AbuseTypePhysical      AbuseType = "physical"       // Физическое насилие
	AbuseTypeEmotional     AbuseType = "emotional"      // Эмоциональное насилие
	AbuseTypePsychological AbuseType = "psychological"  // Психологическое насилие
	AbuseTypeSexual        AbuseType = "sexual"         // Сексуальное насилие
	AbuseTypeFinancial     AbuseType = "financial"      // Финансовое насилие
	AbuseTypeSpiritual     AbuseType = "spiritual"      // Духовное насилие
	AbuseTypeMedical       AbuseType = "medical"        // Медицинское насилие (отказ в лечении)
	AbuseTypeNeglect       AbuseType = "neglect"        // Пренебрежение, игнорирование
)

// AbuseSeverity определяет тяжесть насилия
type AbuseSeverity string

const (
	AbuseSeverityMild      AbuseSeverity = "mild"       // Легкое
	AbuseSeverityModerate  AbuseSeverity = "moderate"   // Умеренное
	AbuseSeveritySevere    AbuseSeverity = "severe"     // Тяжелое
	AbuseSeverityExtreme   AbuseSeverity = "extreme"    // Крайне тяжелое
)

// AbusePattern определяет паттерн насилия
type AbusePattern string

const (
	AbusePatternCycle       AbusePattern = "cycle"        // Цикл насилия (напряжение-взрыв-примирение)
	AbusePatternEscalating  AbusePattern = "escalating"   // Эскалирующее (усиливается со временем)
	AbusePatternChronic     AbusePattern = "chronic"      // Хроническое (постоянное)
	AbusePatternSituational AbusePattern = "situational"  // Ситуативное (в определенных обстоятельствах)
	AbusePatternInherited   AbusePattern = "inherited"    // Унаследованное (из поколения в поколение)
)

// AbuseIncident представляет случай насилия
type AbuseIncident struct {
	ID              string
	Date            *time.Time
	AbuseType       AbuseType
	Severity        AbuseSeverity
	PerpetratorID   string      // ID того, кто совершил насилие
	VictimID        string      // ID жертвы
	Description     string
	WitnessedBy     []string    // ID свидетелей
	Reported        bool
	ReportedTo      string      // кому сообщили (полиция, соцслужба и т.д.)
	Intervention    bool        // было ли вмешательство
	Aftermath       string      // последствия
	Impact          float64     // влияние на жертву (0-1)
	CyclePhase      string      // фаза цикла насилия: "tension", "explosion", "honeymoon"
}

// AbuseDynamics описывает динамику насилия в семье
type AbuseDynamics struct {
	FamilyID                string
	HasPhysicalAbuse        bool
	HasEmotionalAbuse       bool
	HasSexualAbuse          bool
	HasFinancialAbuse       bool
	HasSpiritualAbuse       bool
	HasNeglect              bool
	PrimaryAbuserID         string
	PrimaryVictimID         string
	OtherVictims            []string
	OtherAbusers            []string
	Enablers                []string            // те, кто допускает/покрывает насилие
	Intergenerational       bool                // передается ли через поколения
	Pattern                 AbusePattern
	CycleLength             int                 // длительность цикла в днях (средняя)
	NormalizationLevel      float64             // уровень нормализации насилия в семье (0-1)
	SecrecyLevel            float64             // уровень секретности/сокрытия (0-1)
	FirstIncidentDate       *time.Time
	LastIncidentDate        *time.Time
	IncidentCount           int
	ReportedCount           int
}

// ExtendFamilyMemberWithAbuse расширяет члена семьи данными о насилии
type ExtendFamilyMemberWithAbuse struct {
	*data.ExtendedFamilyMember
	AbuseHistory        *AbuseHistory
	AbuseRiskScore      float64               // риск стать жертвой/абьюзером (0-1)
	ProtectiveFactors   []string              // защитные факторы
	VulnerabilityFactors []string             // факторы уязвимости
	ChildhoodAbuse      bool                  // было ли насилие в детстве
	SupportSystems      []string              // системы поддержки
	TherapyHistory      []string              // история терапии
}

// AbuseAnalysisResult содержит результаты анализа насилия
type AbuseAnalysisResult struct {
	FamilyID            string
	Dynamics            *AbuseDynamics
	Incidents           []AbuseIncident
	VictimsByType       map[AbuseType][]string
	PerpetratorsByType  map[AbuseType][]string
	Timeline            []AbuseIncident
	RiskLevel           string              // "low", "medium", "high", "critical"
	Recommendations     []string
	Warnings            []string
	Resources           []string            // ссылки на ресурсы помощи
}

// AbuseHistory хранит историю насилия для человека
type AbuseHistory struct {
	PersonID            string
	AsVictim            []AbuseIncident
	AsPerpetrator       []AbuseIncident
	AsWitness           []AbuseIncident
	ChildhoodAbuse      bool                  // было ли насилие в детстве
	TraumaSymptoms      []string              // симптомы ПТСР/КПТСР
	CopingMechanisms    []string              // механизмы совладания
	SupportSystems      []string              // системы поддержки
	TherapyHistory      []string              // история терапии
}

// AddAbuseIncident добавляет случай насилия в историю
func (ah *AbuseHistory) AddAbuseIncident(incident AbuseIncident, role string) {
	switch role {
	case "victim":
		ah.AsVictim = append(ah.AsVictim, incident)
	case "perpetrator":
		ah.AsPerpetrator = append(ah.AsPerpetrator, incident)
	case "witness":
		ah.AsWitness = append(ah.AsWitness, incident)
	}
}


// AnalyzeAbuseDynamics анализирует динамику насилия в семье
func AnalyzeAbuseDynamics(fs *model.FamilySystem) *AbuseAnalysisResult {
	result := &AbuseAnalysisResult{
		VictimsByType:      make(map[AbuseType][]string),
		PerpetratorsByType: make(map[AbuseType][]string),
		Incidents:          make([]AbuseIncident, 0),
		Timeline:           make([]AbuseIncident, 0),
		Recommendations:    make([]string, 0),
		Warnings:           make([]string, 0),
		Resources:          make([]string, 0),
	}

	dynamics := &AbuseDynamics{
		OtherVictims:   make([]string, 0),
		OtherAbusers:   make([]string, 0),
		Enablers:       make([]string, 0),
	}

	// Собираем все инциденты из событий
	for id, member := range fs.Members {
		for _, event := range member.Events {
			if strings.Contains(event.EventType, "abuse") || 
			   strings.Contains(event.EventType, "violence") ||
			   strings.Contains(event.EventType, "assault") {
				
				abuseType := determineAbuseType(event.EventType)
				severity := determineAbuseSeverity(event.Description)
				
				incident := AbuseIncident{
					ID:            fmt.Sprintf("abuse_%s_%d", id, event.Date.Unix()),
					Date:          &event.Date,
					AbuseType:     abuseType,
					Severity:      severity,
					PerpetratorID: event.WithPerson,
					VictimID:      id,
					Description:   event.Description,
					WitnessedBy:   make([]string, 0),
					Impact:        event.Impact,
				}
				
				result.Incidents = append(result.Incidents, incident)
				result.Timeline = append(result.Timeline, incident)
				
				// Обновляем статистику
				dynamics.IncidentCount++
				if dynamics.FirstIncidentDate == nil || event.Date.Before(*dynamics.FirstIncidentDate) {
					dynamics.FirstIncidentDate = &event.Date
				}
				if dynamics.LastIncidentDate == nil || event.Date.After(*dynamics.LastIncidentDate) {
					dynamics.LastIncidentDate = &event.Date
				}
				
				// Отмечаем типы насилия
				switch abuseType {
				case AbuseTypePhysical:
					dynamics.HasPhysicalAbuse = true
				case AbuseTypeEmotional, AbuseTypePsychological:
					dynamics.HasEmotionalAbuse = true
				case AbuseTypeSexual:
					dynamics.HasSexualAbuse = true
				case AbuseTypeFinancial:
					dynamics.HasFinancialAbuse = true
				case AbuseTypeSpiritual:
					dynamics.HasSpiritualAbuse = true
				case AbuseTypeNeglect:
					dynamics.HasNeglect = true
				}
				
				// Добавляем в карты по типам
				result.VictimsByType[abuseType] = append(result.VictimsByType[abuseType], id)
				if event.WithPerson != "" {
					result.PerpetratorsByType[abuseType] = append(result.PerpetratorsByType[abuseType], event.WithPerson)
				}
			}
		}
	}

	// Анализируем цикличность
	if len(result.Incidents) > 1 && dynamics.FirstIncidentDate != nil && dynamics.LastIncidentDate != nil {
		days := int(dynamics.LastIncidentDate.Sub(*dynamics.FirstIncidentDate).Hours() / 24)
		if days > 0 && dynamics.IncidentCount > 1 {
			dynamics.CycleLength = days / (dynamics.IncidentCount - 1)
		}
		
		// Определяем паттерн
		if dynamics.CycleLength < 30 && dynamics.IncidentCount > 5 {
			dynamics.Pattern = AbusePatternChronic
		} else if isEscalating(result.Incidents) {
			dynamics.Pattern = AbusePatternEscalating
		} else {
			dynamics.Pattern = AbusePatternCycle
		}
	}

	// Ищем основного абьюзера и жертву
	abuserCount := make(map[string]int)
	victimCount := make(map[string]int)
	
	for _, inc := range result.Incidents {
		abuserCount[inc.PerpetratorID]++
		victimCount[inc.VictimID]++
	}
	
	maxAbuser := 0
	maxVictim := 0
	
	for id, count := range abuserCount {
		if count > maxAbuser {
			maxAbuser = count
			dynamics.PrimaryAbuserID = id
		}
	}
	
	for id, count := range victimCount {
		if count > maxVictim {
			maxVictim = count
			dynamics.PrimaryVictimID = id
		}
	}

		// Ищем энейблеров (тех, кто знал и не вмешивался)
	for id, member := range fs.Members {
		if id == dynamics.PrimaryAbuserID || id == dynamics.PrimaryVictimID {
			continue
		}
		// Проверяем, связан ли с абьюзером или жертвой
		if member.HasPartner(dynamics.PrimaryAbuserID) || 
		   member.HasPartner(dynamics.PrimaryVictimID) ||
		   containsString(member.Children, dynamics.PrimaryVictimID) ||
		   containsString(member.Parents, dynamics.PrimaryVictimID) {
			dynamics.Enablers = append(dynamics.Enablers, id)
		}
	}

	// Проверяем межпоколенческую передачу
	dynamics.Intergenerational = checkIntergenerationalAbuse(fs, result.Incidents)

	// Оцениваем уровень риска
	riskScore := 0.0
	if dynamics.HasPhysicalAbuse {
		riskScore += 0.3
	}
	if dynamics.HasSexualAbuse {
		riskScore += 0.5
	}
	if dynamics.Pattern == AbusePatternEscalating {
		riskScore += 0.4
	}
	if dynamics.IncidentCount > 10 {
		riskScore += 0.3
	}
	if dynamics.CycleLength < 7 { // еженедельно
		riskScore += 0.4
	}
	
	switch {
	case riskScore > 1.5:
		result.RiskLevel = "critical"
	case riskScore > 1.0:
		result.RiskLevel = "high"
	case riskScore > 0.5:
		result.RiskLevel = "medium"
	default:
		result.RiskLevel = "low"
	}

	// Формируем предупреждения
	if dynamics.HasPhysicalAbuse {
		result.Warnings = append(result.Warnings, 
			"⚠️ В семье присутствует физическое насилие - требуется немедленное вмешательство")
	}
	if dynamics.HasSexualAbuse {
		result.Warnings = append(result.Warnings,
			"⚠️⚠️⚠️ СЕКСУАЛЬНОЕ НАСИЛИЕ - требуется экстренное вмешательство служб защиты")
	}
	if dynamics.Pattern == AbusePatternEscalating {
		result.Warnings = append(result.Warnings,
			"⚠️ Эскалирующий паттерн насилия - ситуация ухудшается")
	}
	if len(dynamics.Enablers) > 2 {
		result.Warnings = append(result.Warnings,
			"⚠️ Множество энейблеров - семейная система поддерживает насилие")
	}
	if dynamics.Intergenerational {
		result.Warnings = append(result.Warnings,
			"⚠️ Межпоколенческая передача насилия - требуется разрыв цикла")
	}

	// Рекомендации
	result.Recommendations = generateAbuseRecommendations(dynamics, result.RiskLevel)
	
	// Ресурсы помощи
	result.Resources = []string{
		"📞 В России: телефон доверия 8-800-2000-122",
		"📞 Кризисный центр для женщин: 8-800-7000-600",
		"📞 Центр помощи пережившим насилие: 8-800-333-33-33",
		"🌐 Сайт: http://www.anna-center.ru",
		"🌐 Международные ресурсы: https://www.hotpeachpages.net",
	}

	dynamics.FamilyID = "family_main" // в реальном коде нужно брать ID семьи
	result.Dynamics = dynamics

	return result
}

// determineAbuseType определяет тип насилия по описанию события
func determineAbuseType(eventType string) AbuseType {
	eventType = strings.ToLower(eventType)
	
	if strings.Contains(eventType, "physical") || 
	   strings.Contains(eventType, "beat") ||
	   strings.Contains(eventType, "hit") {
		return AbuseTypePhysical
	}
	if strings.Contains(eventType, "emotional") ||
	   strings.Contains(eventType, "psych") ||
	   strings.Contains(eventType, "verbal") {
		return AbuseTypeEmotional
	}
	if strings.Contains(eventType, "sexual") ||
	   strings.Contains(eventType, "rape") ||
	   strings.Contains(eventType, "molest") {
		return AbuseTypeSexual
	}
	if strings.Contains(eventType, "financial") ||
	   strings.Contains(eventType, "money") {
		return AbuseTypeFinancial
	}
	if strings.Contains(eventType, "spiritual") ||
	   strings.Contains(eventType, "relig") {
		return AbuseTypeSpiritual
	}
	if strings.Contains(eventType, "neglect") ||
	   strings.Contains(eventType, "ignore") {
		return AbuseTypeNeglect
	}
	
	return AbuseTypePsychological
}

// determineAbuseSeverity определяет тяжесть по описанию
func determineAbuseSeverity(description string) AbuseSeverity {
	desc := strings.ToLower(description)
	
	if strings.Contains(desc, "critical") || 
	   strings.Contains(desc, "life-threatening") ||
	   strings.Contains(desc, "hospital") {
		return AbuseSeverityExtreme
	}
	if strings.Contains(desc, "severe") || 
	   strings.Contains(desc, "serious") ||
	   strings.Contains(desc, "injury") {
		return AbuseSeveritySevere
	}
	if strings.Contains(desc, "moderate") {
		return AbuseSeverityModerate
	}
	
	return AbuseSeverityMild
}

// isEscalating проверяет, усиливается ли насилие со временем
func isEscalating(incidents []AbuseIncident) bool {
	if len(incidents) < 3 {
		return false
	}
	
	// Упрощенная проверка: растет ли тяжесть
	severityValues := make([]int, len(incidents))
	for i, inc := range incidents {
		switch inc.Severity {
		case AbuseSeverityMild:
			severityValues[i] = 1
		case AbuseSeverityModerate:
			severityValues[i] = 2
		case AbuseSeveritySevere:
			severityValues[i] = 3
		case AbuseSeverityExtreme:
			severityValues[i] = 4
		}
	}
	
	// Проверяем тренд
	for i := 1; i < len(severityValues); i++ {
		if severityValues[i] < severityValues[i-1] {
			return false
		}
	}
	return true
}

// checkIntergenerationalAbuse проверяет передачу насилия между поколениями
func checkIntergenerationalAbuse(fs *model.FamilySystem, incidents []AbuseIncident) bool {
	// Группируем по поколениям
	abusersByGen := make(map[int]map[string]bool)
	victimsByGen := make(map[int]map[string]bool)
	
	for _, inc := range incidents {
		if inc.PerpetratorID != "" {
			if perp, ok := fs.Members[inc.PerpetratorID]; ok {
				if _, exists := abusersByGen[perp.Generation]; !exists {
					abusersByGen[perp.Generation] = make(map[string]bool)
				}
				abusersByGen[perp.Generation][inc.PerpetratorID] = true
			}
		}
		
		if vic, ok := fs.Members[inc.VictimID]; ok {
			if _, exists := victimsByGen[vic.Generation]; !exists {
				victimsByGen[vic.Generation] = make(map[string]bool)
			}
			victimsByGen[vic.Generation][inc.VictimID] = true
		}
	}
	
	// Если есть абьюзеры в нескольких поколениях
	return len(abusersByGen) >= 2
}

// generateAbuseRecommendations генерирует рекомендации на основе анализа
func generateAbuseRecommendations(dynamics *AbuseDynamics, riskLevel string) []string {
	recs := []string{
		"🛡️ ПЛАН БЕЗОПАСНОСТИ:",
	}
	
	switch riskLevel {
	case "critical":
		recs = append(recs,
			"   • НЕМЕДЛЕННО связаться со службами экстренной помощи (112)",
			"   • Обеспечить безопасное место для жертв насилия",
			"   • Юридическая консультация для подачи заявления",
			"   • Временное изъятие абьюзера из семьи")
	case "high":
		recs = append(recs,
			"   • Срочная консультация с психологом, специализирующимся на насилии",
			"   • Разработать план безопасности на случай эскалации",
			"   • Собрать документы и важные вещи на случай ухода",
			"   • Информировать близких о ситуации")
	case "medium":
		recs = append(recs,
			"   • Начать терапию для всех членов семьи",
			"   • Программы по управлению гневом для абьюзера",
			"   • Группы поддержки для жертв",
			"   • Четкие границы и последствия")
	case "low":
		recs = append(recs,
			"   • Превентивная терапия",
			"   • Обучение здоровым способам выражения гнева",
			"   • Укрепление семейных связей")
	}
	
	if dynamics.HasPhysicalAbuse || dynamics.HasSexualAbuse {
		recs = append(recs,
			"⚖️ ЮРИДИЧЕСКАЯ ПОМОЩЬ:",
			"   • Заявление в полицию",
			"   • Защитный ордер",
			"   • Консультация с юристом по семейному праву")
	}
	
	if dynamics.Intergenerational {
		recs = append(recs,
			"🔄 РАЗРЫВ ЦИКЛА НАСИЛИЯ:",
			"   • Терапия для всех поколений",
			"   • Работа с родовыми сценариями",
			"   • Обучение здоровым моделям отношений")
	}
	
	if len(dynamics.Enablers) > 0 {
		recs = append(recs,
			"👥 РАБОТА С ЭНЕЙБЛЕРАМИ:",
			"   • Осознание роли в поддержании насилия",
			"   • Терапия созависимости",
			"   • Обучение правильному реагированию")
	}
	
	return recs
}

// AnalyzePersonAbuseHistory анализирует историю насилия конкретного человека
func AnalyzePersonAbuseHistory(member *ExtendFamilyMemberWithAbuse, fs *model.FamilySystem) *AbuseHistory {
	history := &AbuseHistory{
		PersonID:         member.ID,
		AsVictim:         make([]AbuseIncident, 0),
		AsPerpetrator:    make([]AbuseIncident, 0),
		AsWitness:        make([]AbuseIncident, 0),
		TraumaSymptoms:   make([]string, 0),
		CopingMechanisms: make([]string, 0),
		SupportSystems:   make([]string, 0),
		TherapyHistory:   make([]string, 0),
	}
	
	// Анализируем события
	for _, event := range member.Events {
		if strings.Contains(event.EventType, "abuse") || 
		   strings.Contains(event.EventType, "violence") {
			
			abuseType := determineAbuseType(event.EventType)
			severity := determineAbuseSeverity(event.Description)
			
			incident := AbuseIncident{
				ID:            fmt.Sprintf("abuse_%s_%d", member.ID, event.Date.Unix()),
				Date:          &event.Date,
				AbuseType:     abuseType,
				Severity:      severity,
				PerpetratorID: event.WithPerson,
				VictimID:      member.ID,
				Description:   event.Description,
				Impact:        event.Impact,
			}
			
			history.AsVictim = append(history.AsVictim, incident)
			
			// Отмечаем детское насилие
			ageAtIncident := int(event.Date.Sub(member.BirthDate).Hours() / 24 / 365)
			if ageAtIncident < 18 {
				history.ChildhoodAbuse = true
			}
		}
	}
	
	// Анализируем симптомы ПТСР на основе векторов
	coords := core.ParseDateToCoords(member.BirthDate)
	vectors := core.CalculateVectors(coords)
	amp := core.VectorAmplitude(vectors)
	
	if amp < 2.0 {
		history.TraumaSymptoms = append(history.TraumaSymptoms, 
			"Диссоциация, эмоциональное оцепенение")
	} else if amp > 8.0 {
		history.TraumaSymptoms = append(history.TraumaSymptoms,
			"Гипервозбуждение, тревога, вспышки гнева")
	}
	
	// Анализируем механизмы совладания по роли
	if role, ok := fs.TraumaRoles[member.ID]; ok {
		switch role.PrimaryRole {
		case core.TraumaRoleLostChild:
			history.CopingMechanisms = append(history.CopingMechanisms,
				"Избегание, уход в себя, минимизация контактов")
		case core.TraumaRoleMascot:
			history.CopingMechanisms = append(history.CopingMechanisms,
				"Юмор как защита, минимизация серьезности")
		case core.TraumaRoleParentified:
			history.CopingMechanisms = append(history.CopingMechanisms,
				"Контроль, гиперответственность, забота о других")
		}
	}
	
	return history
}

// CalculateAbuseRisk вычисляет риск стать жертвой или абьюзером
func CalculateAbuseRisk(member *ExtendFamilyMemberWithAbuse, fs *model.FamilySystem) float64 {
	risk := 0.0
	
	// Факторы риска
	if member.ChildhoodAbuse {
		risk += 0.4 // люди, пережившие насилие в детстве, чаще становятся жертвами или абьюзерами
	}
	
	// Роль в семье
	if role, ok := fs.TraumaRoles[member.ID]; ok {
		switch role.PrimaryRole {
		case core.TraumaRoleScapegoat:
			risk += 0.3 // козлы отпущения часто становятся жертвами
		case core.TraumaRoleInvisible:
			risk += 0.2 // невидимки могут быть незамеченными жертвами
		case core.TraumaRoleGoldenChild:
			risk += 0.1 // любимчики могут перенимать абьюзивные паттерны
		}
	}
	
	// Наличие абьюзера в близком окружении
	abuseAnalysis := AnalyzeAbuseDynamics(fs)
	if abuseAnalysis.Dynamics.PrimaryAbuserID != "" {
		if member.HasPartner(abuseAnalysis.Dynamics.PrimaryAbuserID) ||
		   containsString(member.Parents, abuseAnalysis.Dynamics.PrimaryAbuserID) {
			risk += 0.5
		}
	}
	
	// Защитные факторы снижают риск
	if len(member.SupportSystems) > 2 {
		risk -= 0.2
	}
	if len(member.TherapyHistory) > 0 {
		risk -= 0.3
	}
	
	if risk < 0 {
		risk = 0
	}
	if risk > 1 {
		risk = 1
	}
	
	return risk
}

// containsString проверяет наличие строки в слайсе
func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
