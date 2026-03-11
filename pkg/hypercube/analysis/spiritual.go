package analysis

import (
	"fmt"
	"math"
	"strings"
	"time"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/core"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/data"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/model"
)

// SpiritualDimension определяет духовное измерение
type SpiritualDimension string

const (
	SpiritualDimensionPurpose     SpiritualDimension = "purpose"      // Предназначение
	SpiritualDimensionMeaning     SpiritualDimension = "meaning"      // Смысл жизни
	SpiritualDimensionConnection  SpiritualDimension = "connection"   // Связь с высшим/вселенной
	SpiritualDimensionAncestral   SpiritualDimension = "ancestral"    // Родовая связь
	SpiritualDimensionKarmic      SpiritualDimension = "karmic"       // Кармические уроки
	SpiritualDimensionSoul        SpiritualDimension = "soul"         // Душа/самость
	SpiritualDimensionShadow      SpiritualDimension = "shadow"       // Теневая сторона
	SpiritualDimensionArchetype   SpiritualDimension = "archetype"    // Архетипы
)

// SpiritualPath определяет духовный путь
type SpiritualPath string

const (
	SpiritualPathTraditional   SpiritualPath = "traditional"   // Традиционная религия
	SpiritualPathMystical      SpiritualPath = "mystical"      // Мистический/эзотерический
	SpiritualPathPhilosophical SpiritualPath = "philosophical" // Философский
	SpiritualPathNature        SpiritualPath = "nature"        // Природный/языческий
	SpiritualPathIntegrative   SpiritualPath = "integrative"   // Интегративный
	SpiritualPathSecular       SpiritualPath = "secular"       // Светский (гуманизм)
	SpiritualPathCrisis        SpiritualPath = "crisis"        // Духовный кризис
	SpiritualPathAwakening     SpiritualPath = "awakening"     // Пробуждение
)

// SoulArchetype определяет архетип души
type SoulArchetype string

const (
	SoulArchetypeInnocent      SoulArchetype = "innocent"      // Невинный
	SoulArchetypeOrphan        SoulArchetype = "orphan"        // Сирота
	SoulArchetypeWarrior       SoulArchetype = "warrior"       // Воин
	SoulArchetypeCaregiver     SoulArchetype = "caregiver"     // Заботливый
	SoulArchetypeSeeker        SoulArchetype = "seeker"        // Искатель
	SoulArchetypeLover         SoulArchetype = "lover"         // Любовник
	SoulArchetypeDestroyer     SoulArchetype = "destroyer"     // Разрушитель
	SoulArchetypeCreator       SoulArchetype = "creator"       // Созидатель
	SoulArchetypeRuler         SoulArchetype = "ruler"         // Правитель
	SoulArchetypeMagician      SoulArchetype = "magician"      // Маг
	SoulArchetypeSage          SoulArchetype = "sage"          // Мудрец
	SoulArchetypeFool          SoulArchetype = "fool"          // Дурак (в смысле свободы)
	SoulArchetypeHealer        SoulArchetype = "healer"        // Целитель
	SoulArchetypeVisionary     SoulArchetype = "visionary"     // Визионер
	SoulArchetypeAlchemist     SoulArchetype = "alchemist"     // Алхимик
)

// KarmicLesson представляет кармический урок
type KarmicLesson struct {
	ID              string
	Type            string                    // "forgiveness", "love", "trust", "boundaries", "surrender"
	Description     string
	OriginID        string                    // ID человека или события, создавшего урок
	AffectedIDs     []string                  // кого касается
	Difficulty      float64                   // сложность (0-1)
	Progress        float64                   // прогресс прохождения (0-1)
	IsLearned       bool                       // пройден ли урок
	Manifestations  []string                   // как проявляется в жизни
	Blockages       []string                   // что блокирует прохождение
	Helpers         []string                   // кто помогает пройти
	PastLives       []string                   // возможные прошлые жизни (метафора)
}

// SoulMission представляет миссию души
type SoulMission struct {
	ID              string
	PersonID        string
	PrimaryMission  string                    // основная миссия
	SecondaryMissions []string                 // дополнительные миссии
	LifeAreas       []string                   // сферы жизни для реализации
	Challenges      []string                   // вызовы на пути
	Gifts           []string                   // дары/таланты для выполнения миссии
	Timing          string                      // время реализации
	IsActive        bool                        // активна ли сейчас
	Manifestations  []string                    // проявления в жизни
}

// AncestralLegacy представляет родовое наследие
type AncestralLegacy struct {
	ID              string
	PersonID        string
	Lineage         string                      // "mother", "father", "both"
	Gifts           []string                     // дары рода
	Curses          []string                     // родовые проклятия/паттерны
	Unresolved      []string                     // непроработанные темы
	Ancestors       []string                     // значимые предки
	Generations     int                          // сколько поколений затронуто
	HealingNeeded   []string                     // что нужно исцелить
	Rituals         []string                     // рекомендуемые родовые практики
}

// SpiritualCrisis представляет духовный кризис
type SpiritualCrisis struct {
	ID              string
	PersonID        string
	StartDate       time.Time
	EndDate         *time.Time
	Triggers        []string                     // триггеры кризиса
	Symptoms        []string                     // симптомы
	Stage           string                       // "dark_night", "searching", "awakening", "integration"
	SupportNeeded   []string                     // необходимая поддержка
	Outcome         string                       // исход
	Transformation  string                       // произошедшая трансформация
}

// SpiritualProfile содержит духовный профиль человека
type SpiritualProfile struct {
	PersonID           string
	Archetype          SoulArchetype              // основной архетип
	SecondaryArchetypes []SoulArchetype           // дополнительные архетипы
	Path               SpiritualPath              // духовный путь
	Dimensions         map[SpiritualDimension]float64 // развитие по измерениям (0-1)
	KarmicLessons      []KarmicLesson              // кармические уроки
	Mission            *SoulMission                // миссия души
	Ancestral          *AncestralLegacy            // родовое наследие
	Crises             []SpiritualCrisis           // духовные кризисы
	Practices          []string                    // духовные практики
	Teachers           []string                    // учителя/наставники
	SacredSymbols      []string                    // значимые символы
	LifeThemes         []string                    // основные темы жизни
	ShadowAspects      []string                    // теневые аспекты
	IntegrationLevel   float64                      // уровень интеграции (0-1)
}

// SpiritualAnalysisResult содержит результаты духовного анализа
type SpiritualAnalysisResult struct {
	FamilyID            string
	Profiles            map[string]*SpiritualProfile
	CollectiveArchetype SoulArchetype               // архетип семьи в целом
	CollectiveMission   string                       // общая миссия семьи
	KarmicWeb           []KarmicLesson                // кармическая сеть
	AncestralPatterns   []string                      // родовые паттерны
	SpiritualAge        float64                       // духовный возраст семьи (0-1)
	Recommendations     []string
	Warnings            []string
	Rituals             []string                      // рекомендуемые ритуалы
}

// ExtendFamilyMemberWithSpiritual расширяет члена семьи духовными данными
type ExtendFamilyMemberWithSpiritual struct {
	*data.ExtendedFamilyMember
	SpiritualProfile    *SpiritualProfile
	AwakeningLevel      float64                       // уровень пробуждения (0-1)
	SoulAge             int                            // возраст души (1-10)
	LifePurpose         string                         // предназначение
	Dharma              string                         // дхарма (путь)
	Karma               float64                        // кармический баланс
}

// CalculateSoulArchetype определяет архетип души на основе даты рождения и травматической роли
func CalculateSoulArchetype(member *data.ExtendedFamilyMember, traumaRole core.TraumaRole) SoulArchetype {
	coords := core.ParseDateToCoords(member.BirthDate)
	vectors := core.CalculateVectors(coords)
	amp := core.VectorAmplitude(vectors)
	
	// Комбинация даты, векторов и травматической роли
	
	// По травматической роли
	switch traumaRole {
	case core.TraumaRoleScapegoat:
		return SoulArchetypeWarrior
	case core.TraumaRoleGoldenChild:
		return SoulArchetypeRuler
	case core.TraumaRoleLostChild:
		return SoulArchetypeOrphan
	case core.TraumaRoleInvisible:
		return SoulArchetypeSeeker
	case core.TraumaRoleParentified:
		return SoulArchetypeCaregiver
	case core.TraumaRoleMascot:
		return SoulArchetypeFool
	case core.TraumaRoleEmotionalSpouse:
		return SoulArchetypeLover
	case core.TraumaRoleHealer:
		return SoulArchetypeHealer
	case core.TraumaRoleReplacement:
		return SoulArchetypeAlchemist
	}
	
	// По амплитуде векторов
	if amp > 7.0 {
		return SoulArchetypeDestroyer
	} else if amp < 2.0 {
		return SoulArchetypeSage
	}
	
	// По дате (нумерологически)
	day := member.BirthDate.Day()
	switch day % 7 {
	case 0:
		return SoulArchetypeMagician
	case 1:
		return SoulArchetypeInnocent
	case 2:
		return SoulArchetypeSeeker
	case 3:
		return SoulArchetypeCreator
	case 4:
		return SoulArchetypeVisionary
	case 5:
		return SoulArchetypeWarrior
	case 6:
		return SoulArchetypeSage
	}
	
	return SoulArchetypeSeeker
}

// CalculateSpiritualPath определяет духовный путь
func CalculateSpiritualPath(member *data.ExtendedFamilyMember) SpiritualPath {
	// По событиям жизни
	for _, event := range member.Events {
		if strings.Contains(event.EventType, "religious") ||
		   strings.Contains(event.EventType, "church") {
			return SpiritualPathTraditional
		}
		if strings.Contains(event.EventType, "meditation") ||
		   strings.Contains(event.EventType, "yoga") {
			return SpiritualPathIntegrative
		}
		if strings.Contains(event.EventType, "crisis") &&
		   strings.Contains(event.EventType, "spiritual") {
			return SpiritualPathCrisis
		}
	}
	
	// По координатам
	coords := core.ParseDateToCoords(member.BirthDate)
	if coords.W == 27 { // 999 в сумме
		return SpiritualPathAwakening
	}
	
	return SpiritualPathSecular
}

// CalculateKarmicLessons вычисляет кармические уроки на основе семейных паттернов
func CalculateKarmicLessons(member *data.ExtendedFamilyMember, fs *model.FamilySystem) []KarmicLesson {
	lessons := make([]KarmicLesson, 0)
	
	// Урок на основе травматической роли
	if role, ok := fs.TraumaRoles[member.ID]; ok {
		var lessonType string
		switch role.PrimaryRole {
		case core.TraumaRoleScapegoat:
			lessonType = "self-worth"
		case core.TraumaRoleGoldenChild:
			lessonType = "humility"
		case core.TraumaRoleLostChild:
			lessonType = "connection"
		case core.TraumaRoleInvisible:
			lessonType = "visibility"
		case core.TraumaRoleParentified:
			lessonType = "surrender"
		case core.TraumaRoleEmotionalSpouse:
			lessonType = "boundaries"
		case core.TraumaRoleReplacement:
			lessonType = "identity"
		default:
			lessonType = "integration"
		}
		
		lesson := KarmicLesson{
			ID:           fmt.Sprintf("karma_%s_%s", member.ID, lessonType),
			Type:         lessonType,
			Description:  getLessonDescription(lessonType, role.PrimaryRole),
			OriginID:     findLessonOrigin(member, fs, lessonType),
			AffectedIDs:  []string{member.ID},
			Difficulty:   calculateLessonDifficulty(member, role),
			Progress:     calculateLessonProgress(member, lessonType),
			IsLearned:    false,
			Manifestations: getLessonManifestations(lessonType, member),
		}
		
		if lesson.Progress > 0.8 {
			lesson.IsLearned = true
		}
		
		lessons = append(lessons, lesson)
	}
	
	// Уроки на основе семейных паттернов
	for _, pattern := range fs.Patterns {
		if containsString(pattern.Members, member.ID) {
			lesson := KarmicLesson{
				ID:           fmt.Sprintf("karma_pattern_%s_%s", member.ID, pattern.PatternType),
				Type:         "family_pattern",
				Description:  fmt.Sprintf("Проработка семейного паттерна: %s", pattern.Description),
				OriginID:     findPatternOrigin(pattern, fs),
				AffectedIDs:  pattern.Members,
				Difficulty:   pattern.Severity,
				Progress:     0.3, // начальный прогресс
				Manifestations: []string{pattern.Description},
			}
			lessons = append(lessons, lesson)
		}
	}
	
	return lessons
}

// getLessonDescription возвращает описание кармического урока
func getLessonDescription(lessonType string, role core.TraumaRole) string {
	descriptions := map[string]map[core.TraumaRole]string{
		"self-worth": {
			core.TraumaRoleScapegoat: "Научиться ценить себя независимо от мнения других. Перестать быть 'козлом отпущения' и принять свою ценность.",
		},
		"humility": {
			core.TraumaRoleGoldenChild: "Научиться принимать несовершенство свое и других. Отпустить перфекционизм и право быть просто человеком.",
		},
		"connection": {
			core.TraumaRoleLostChild: "Научиться устанавливать глубокие связи с людьми. Выйти из изоляции и позволить себе быть видимым.",
		},
		"visibility": {
			core.TraumaRoleInvisible: "Научиться занимать место в этом мире. Заявить о своих потребностях и праве на существование.",
		},
		"surrender": {
			core.TraumaRoleParentified: "Научиться отпускать контроль и принимать помощь. Позволить себе быть уязвимым и получать заботу.",
		},
		"boundaries": {
			core.TraumaRoleEmotionalSpouse: "Научиться выстраивать здоровые границы. Отделить себя от родителя и найти свою идентичность.",
		},
		"identity": {
			core.TraumaRoleReplacement: "Научиться быть собой, а не заменой. Присвоить право на собственную жизнь и уникальность.",
		},
	}
	
	if descMap, ok := descriptions[lessonType]; ok {
		if desc, ok := descMap[role]; ok {
			return desc
		}
	}
	
	return "Интеграция жизненного опыта и обретение целостности."
}

// findLessonOrigin находит источник кармического урока
func findLessonOrigin(member *data.ExtendedFamilyMember, fs *model.FamilySystem, lessonType string) string {
	// Ищем в родителях
	for _, parentID := range member.Parents {
		if parent, ok := fs.Members[parentID]; ok {
			if role, ok := fs.TraumaRoles[parentID]; ok {
				// Проверяем соответствие типов уроков
				if (lessonType == "self-worth" && role.PrimaryRole == core.TraumaRoleScapegoat) ||
				   (lessonType == "boundaries" && role.PrimaryRole == core.TraumaRoleEmotionalSpouse) {
					return parentID
				}
			}
		}
	}
	return "ancestral"
}

// calculateLessonDifficulty вычисляет сложность урока
func calculateLessonDifficulty(member *data.ExtendedFamilyMember, role *core.IdentifiedTraumaRoles) float64 {
	difficulty := 0.5
	
	// Чем выше уверенность в роли, тем сложнее урок
	difficulty += role.Confidence * 0.3
	
	// Чем старше человек, тем сложнее менять паттерны
	age := member.GetAge()
	if age > 50 {
		difficulty += 0.2
	} else if age > 30 {
		difficulty += 0.1
	}
	
	// Наличие поддержки снижает сложность
	if len(member.Partners) > 0 {
		difficulty -= 0.1
	}
	
	if difficulty > 1.0 {
		difficulty = 1.0
	}
	if difficulty < 0.1 {
		difficulty = 0.1
	}
	
	return difficulty
}

// calculateLessonProgress вычисляет прогресс в прохождении урока
func calculateLessonProgress(member *data.ExtendedFamilyMember, lessonType string) float64 {
	progress := 0.3 // базовый прогресс
	
	// По событиям жизни
	for _, event := range member.Events {
		if strings.Contains(event.EventType, "therapy") {
			progress += 0.1
		}
		if strings.Contains(event.EventType, "awakening") {
			progress += 0.2
		}
		if strings.Contains(event.EventType, "breakthrough") {
			progress += 0.15
		}
	}
	
	// По количеству партнеров (опыт в отношениях)
	if len(member.Partners) > 3 {
		progress += 0.1
	}
	
	if progress > 1.0 {
		progress = 1.0
	}
	
	return progress
}

// getLessonManifestations возвращает проявления урока в жизни
func getLessonManifestations(lessonType string, member *data.ExtendedFamilyMember) []string {
	manifestations := make([]string, 0)
	
	switch lessonType {
	case "self-worth":
		manifestations = append(manifestations,
			"Трудности с отстаиванием границ",
			"Склонность угождать другим",
			"Чувство, что недостаточно хорош",
			"Привлечение критикующих партнеров")
	case "boundaries":
		manifestations = append(manifestations,
			"Симбиотические отношения",
			"Трудности с разделением своих и чужих чувств",
			"Чувство ответственности за других",
			"Страх близости")
	case "identity":
		manifestations = append(manifestations,
			"Чувство 'не на своем месте'",
			"Поиск себя",
			"Сравнение с другими",
			"Неустойчивая самооценка")
	}
	
	return manifestations
}

// findPatternOrigin находит происхождение паттерна
func findPatternOrigin(pattern model.FamilyPattern, fs *model.FamilySystem) string {
	// Ищем самого старшего участника паттерна
	oldestGen := 10
	oldestID := ""
	
	for _, memberID := range pattern.Members {
		if member, ok := fs.Members[memberID]; ok {
			if member.Generation < oldestGen {
				oldestGen = member.Generation
				oldestID = memberID
			}
		}
	}
	
	if oldestID != "" {
		return oldestID
	}
	
	return "ancestral"
}

// CalculateSoulMission определяет миссию души
func CalculateSoulMission(member *data.ExtendedFamilyMember, archetype SoulArchetype, lessons []KarmicLesson) *SoulMission {
	mission := &SoulMission{
		PersonID:        member.ID,
		PrimaryMission:  getPrimaryMission(archetype),
		SecondaryMissions: make([]string, 0),
		LifeAreas:       make([]string, 0),
		Challenges:      make([]string, 0),
		Gifts:           make([]string, 0),
		IsActive:        true,
		Manifestations:  make([]string, 0),
	}
	
	// Дополнительные миссии на основе уроков
	for _, lesson := range lessons {
		if !lesson.IsLearned {
			mission.SecondaryMissions = append(mission.SecondaryMissions,
				fmt.Sprintf("Проработка %s", lesson.Type))
		}
	}
	
	// Сферы жизни
	switch archetype {
	case SoulArchetypeWarrior:
		mission.LifeAreas = []string{"карьера", "спорт", "защита других"}
		mission.Gifts = []string{"сила", "мужество", "решительность"}
		mission.Challenges = []string{"гнев", "конфликтность", "нетерпимость"}
	case SoulArchetypeHealer:
		mission.LifeAreas = []string{"медицина", "психология", "помощь другим"}
		mission.Gifts = []string{"эмпатия", "интуиция", "терпение"}
		mission.Challenges = []string{"эмоциональное выгорание", "спасательство"}
	case SoulArchetypeCreator:
		mission.LifeAreas = []string{"искусство", "творчество", "инновации"}
		mission.Gifts = []string{"воображение", "оригинальность", "чувствительность"}
		mission.Challenges = []string{"нестабильность", "чувствительность к критике"}
	case SoulArchetypeSage:
		mission.LifeAreas = []string{"образование", "философия", "наставничество"}
		mission.Gifts = []string{"мудрость", "аналитическое мышление", "обучение"}
		mission.Challenges = []string{"отстраненность", "склонность к назиданию"}
	}
	
	return mission
}

// getPrimaryMission возвращает основную миссию для архетипа
func getPrimaryMission(archetype SoulArchetype) string {
	switch archetype {
	case SoulArchetypeInnocent:
		return "Вернуть веру в добро и сохранить чистоту намерений"
	case SoulArchetypeOrphan:
		return "Найти свое место и создать настоящую семью"
	case SoulArchetypeWarrior:
		return "Защищать слабых и отстаивать справедливость"
	case SoulArchetypeCaregiver:
		return "Заботиться о других, не теряя себя"
	case SoulArchetypeSeeker:
		return "Искать истину и делиться находками с другими"
	case SoulArchetypeLover:
		return "Учиться любить безусловно и принимать любовь"
	case SoulArchetypeDestroyer:
		return "Разрушать отжившее, чтобы освободить место новому"
	case SoulArchetypeCreator:
		return "Созидать красоту и вдохновлять других"
	case SoulArchetypeRuler:
		return "Создавать порядок и гармонию в своем мире"
	case SoulArchetypeMagician:
		return "Трансформировать реальность через намерение"
	case SoulArchetypeSage:
		return "Постигать мудрость и передавать знания"
	case SoulArchetypeFool:
		return "Жить в моменте и доверять потоку жизни"
	case SoulArchetypeHealer:
		return "Исцелять себя и помогать исцеляться другим"
	case SoulArchetypeVisionary:
		return "Видеть будущее и вести за собой"
	case SoulArchetypeAlchemist:
		return "Трансформировать страдания в золото мудрости"
	default:
		return "Прожить жизнь в гармонии с собой и миром"
	}
}

// CalculateAncestralLegacy вычисляет родовое наследие
func CalculateAncestralLegacy(member *data.ExtendedFamilyMember, fs *model.FamilySystem) *AncestralLegacy {
	legacy := &AncestralLegacy{
		PersonID:      member.ID,
		Gifts:         make([]string, 0),
		Curses:        make([]string, 0),
		Unresolved:    make([]string, 0),
		Ancestors:     make([]string, 0),
		HealingNeeded: make([]string, 0),
		Rituals:       make([]string, 0),
	}
	
	// Анализируем родителей
	for _, parentID := range member.Parents {
		if parent, ok := fs.Members[parentID]; ok {
			legacy.Ancestors = append(legacy.Ancestors, parent.Name)
			legacy.Generations++
			
			// Дары от родителей
			if role, ok := fs.TraumaRoles[parentID]; ok {
				switch role.PrimaryRole {
				case core.TraumaRoleHealer:
					legacy.Gifts = append(legacy.Gifts, "способность исцелять")
				case core.TraumaRoleWarrior:
					legacy.Gifts = append(legacy.Gifts, "сила и стойкость")
				case core.TraumaRoleCreator:
					legacy.Gifts = append(legacy.Gifts, "творческие способности")
				}
				
				// Непроработанные паттерны
				legacy.Unresolved = append(legacy.Unresolved, 
					fmt.Sprintf("непроработанная роль %s в роду", role.PrimaryRole))
			}
		}
	}
	
	// Анализируем семейные паттерны
	for _, pattern := range fs.Patterns {
		if containsString(pattern.Members, member.ID) {
			legacy.Curses = append(legacy.Curses, pattern.Description)
			legacy.HealingNeeded = append(legacy.HealingNeeded, 
				fmt.Sprintf("исцеление паттерна %s", pattern.PatternType))
		}
	}
	
	// Рекомендуемые родовые ритуалы
	if len(legacy.Unresolved) > 0 {
		legacy.Rituals = append(legacy.Rituals,
			"Составление genogramma (родового древа)",
			"Медитация на связь с родом",
			"Письма благодарности предкам",
			"Родовая расстановка")
	}
	
	return legacy
}

// AnalyzeSpiritualProfile создает полный духовный профиль человека
func AnalyzeSpiritualProfile(member *data.ExtendedFamilyMember, fs *model.FamilySystem) *SpiritualProfile {
	if member == nil {
		return nil
	}
	
	role := fs.TraumaRoles[member.ID]
	if role == nil {
		role = &core.IdentifiedTraumaRoles{PrimaryRole: core.TraumaRoleShadow, Confidence: 0.5}
	}
	
	archetype := CalculateSoulArchetype(member, role.PrimaryRole)
	lessons := CalculateKarmicLessons(member, fs)
	mission := CalculateSoulMission(member, archetype, lessons)
	ancestral := CalculateAncestralLegacy(member, fs)
	
	// Измерения духовного развития
	dimensions := make(map[SpiritualDimension]float64)
	
	coords := core.ParseDateToCoords(member.BirthDate)
	vectors := core.CalculateVectors(coords)
	amp := core.VectorAmplitude(vectors)
	
	// Целостность (интеграция) на основе амплитуды
	dimensions[SpiritualDimensionSoul] = 1.0 - math.Abs(amp-5.0)/10.0
	
	// Кармическое развитие на основе прогресса уроков
	karmicProgress := 0.0
	for _, lesson := range lessons {
		karmicProgress += lesson.Progress
	}
	if len(lessons) > 0 {
		dimensions[SpiritualDimensionKarmic] = karmicProgress / float64(len(lessons))
	} else {
		dimensions[SpiritualDimensionKarmic] = 0.5
	}
	
	// Связь с родом на основе поколений
	if member.Generation > 1 {
		dimensions[SpiritualDimensionAncestral] = 0.5 + float64(member.Generation)*0.1
	} else {
		dimensions[SpiritualDimensionAncestral] = 0.3
	}
	
	// Предназначение на основе наличия миссии
	if mission != nil {
		dimensions[SpiritualDimensionPurpose] = 0.6
	} else {
		dimensions[SpiritualDimensionPurpose] = 0.3
	}
	
	// Теневая сторона на основе травматической роли
	dimensions[SpiritualDimensionShadow] = role.Confidence
	
	// Ограничиваем значения
	for dim, val := range dimensions {
		if val < 0 {
			dimensions[dim] = 0
		}
		if val > 1 {
			dimensions[dim] = 1
		}
	}
	
	profile := &SpiritualProfile{
		PersonID:           member.ID,
		Archetype:          archetype,
		SecondaryArchetypes: make([]SoulArchetype, 0),
		Path:               CalculateSpiritualPath(member),
		Dimensions:         dimensions,
		KarmicLessons:      lessons,
		Mission:            mission,
		Ancestral:          ancestral,
		Crises:             make([]SpiritualCrisis, 0),
		Practices:          make([]string, 0),
		Teachers:           make([]string, 0),
		SacredSymbols:      make([]string, 0),
		LifeThemes:         make([]string, 0),
		ShadowAspects:      make([]string, 0),
		IntegrationLevel:   dimensions[SpiritualDimensionSoul],
	}
	
	// Добавляем теневые аспекты на основе роли
	switch role.PrimaryRole {
	case core.TraumaRoleGoldenChild:
		profile.ShadowAspects = append(profile.ShadowAspects, 
			"перфекционизм", "страх неудачи", "потребность в одобрении")
	case core.TraumaRoleScapegoat:
		profile.ShadowAspects = append(profile.ShadowAspects,
			"чувство неполноценности", "гнев", "самообвинение")
	case core.TraumaRoleLostChild:
		profile.ShadowAspects = append(profile.ShadowAspects,
			"изоляция", "страх близости", "эмоциональная отстраненность")
	case core.TraumaRoleReplacement:
		profile.ShadowAspects = append(profile.ShadowAspects,
			"диффузия идентичности", "вина выжившего", "чувство 'не своего места'")
	}
	
	// Жизненные темы
	profile.LifeThemes = append(profile.LifeThemes,
		fmt.Sprintf("Основной архетип: %s", archetype),
		fmt.Sprintf("Духовный путь: %s", profile.Path),
		fmt.Sprintf("Ключевой урок: %s", lessons[0].Type))
	
	return profile
}

// AnalyzeFamilySpirituality анализирует духовное состояние семьи
func AnalyzeFamilySpirituality(fs *model.FamilySystem) *SpiritualAnalysisResult {
	result := &SpiritualAnalysisResult{
		FamilyID:          "family_main",
		Profiles:          make(map[string]*SpiritualProfile),
		KarmicWeb:         make([]KarmicLesson, 0),
		AncestralPatterns: make([]string, 0),
		Recommendations:   make([]string, 0),
		Warnings:          make([]string, 0),
		Rituals:           make([]string, 0),
	}
	
	// Создаем профили для всех членов семьи
	archetypeCount := make(map[SoulArchetype]int)
	
	for id, member := range fs.Members {
		profile := AnalyzeSpiritualProfile(member, fs)
		result.Profiles[id] = profile
		archetypeCount[profile.Archetype]++
		
		// Собираем все кармические уроки
		result.KarmicWeb = append(result.KarmicWeb, profile.KarmicLessons...)
	}
	
	// Определяем коллективный архетип семьи
	maxCount := 0
	for archetype, count := range archetypeCount {
		if count > maxCount {
			maxCount = count
			result.CollectiveArchetype = archetype
		}
	}
	
	// Коллективная миссия
	switch result.CollectiveArchetype {
	case SoulArchetypeHealer:
		result.CollectiveMission = "Исцеление родовых травм и помощь другим"
	case SoulArchetypeWarrior:
		result.CollectiveMission = "Преодоление трудностей и защита ценностей"
	case SoulArchetypeCreator:
		result.CollectiveMission = "Созидание нового и передача творческого наследия"
	default:
		result.CollectiveMission = "Интеграция опыта и гармонизация отношений"
	}
	
	// Родовые паттерны
	for _, pattern := range fs.Patterns {
		result.AncestralPatterns = append(result.AncestralPatterns, pattern.Description)
	}
	
	// Духовный возраст семьи (средний уровень интеграции)
	totalIntegration := 0.0
	for _, profile := range result.Profiles {
		totalIntegration += profile.IntegrationLevel
	}
	if len(result.Profiles) > 0 {
		result.SpiritualAge = totalIntegration / float64(len(result.Profiles))
	}
	
	// Предупреждения
	if result.SpiritualAge < 0.3 {
		result.Warnings = append(result.Warnings,
			"⚠️ Низкий уровень духовной интеграции в семье - высокий риск конфликтов")
	}
	if len(result.KarmicWeb) > len(result.Profiles)*2 {
		result.Warnings = append(result.Warnings,
			"⚠️ Высокая кармическая нагрузка - множество непроработанных уроков")
	}
	
	// Рекомендации
	result.Recommendations = []string{
		"🧘 Практики осознанности для всех членов семьи",
		"📿 Родовые расстановки для проработки родовых паттернов",
		"📖 Изучение духовных традиций, созвучных семье",
		"🌿 Время на природе вместе для восстановления связи",
		"🎭 Творческое самовыражение как путь к архетипам",
	}
	
	// Ритуалы
	result.Rituals = []string{
		"🔥 Ритуал благодарности предкам (ежемесячно)",
		"🕯️ Медитация на связь с родом (еженедельно)",
		"🌅 Встреча рассвета вместе в важные даты",
		"📝 Создание семейной книги мудрости",
		"🎨 Совместное творчество как медитация",
	}
	
	return result
}

// InterpretLifePath дает духовную интерпретацию жизненного пути
func InterpretLifePath(member *data.ExtendedFamilyMember, vectors core.PersonVectors) string {
	amp := core.VectorAmplitude(vectors)
	coords := core.ParseDateToCoords(member.BirthDate)
	
	interpretation := fmt.Sprintf("✨ ДУХОВНАЯ ИНТЕРПРЕТАЦИЯ ЖИЗНЕННОГО ПУТИ %s ✨\n\n", member.Name)
	
	// Интерпретация амплитуды
	if amp > 7.0 {
		interpretation += "Ваша жизнь отмечена высокой интенсивностью и страстью. "
		interpretation += "Вы подобны вулкану - мощная энергия ищет выход. "
		interpretation += "Ваша задача - научиться направлять эту силу в созидательное русло, "
		interpretation += "не разрушая себя и других.\n\n"
	} else if amp < 2.0 {
		interpretation += "Ваша жизнь течет спокойно и размеренно. "
		interpretation += "Вы подобны глубокому озеру - внешне тихому, но хранящему глубину. "
		interpretation += "Ваша задача - не бояться глубины, исследовать свои внутренние миры, "
		interpretation += "и позволять тишине говорить.\n\n"
	} else {
		interpretation += "Ваша жизнь имеет сбалансированный ритм. "
		interpretation += "Вы подобны реке - иногда быстрой, иногда медленной, но всегда движущейся. "
		interpretation += "Ваша задача - сохранять этот баланс, не давая крайностям увести вас с пути.\n\n"
	}
	
	// Интерпретация W (ключа)
	if coords.W == 27 {
		interpretation += "🌟 ОСОБАЯ ОТМЕТКА: Ваша контрольная сумма равна 27 - числу моста. "
		interpretation += "В фильме 'Куб' комната 999 была выходом. Вы обладаете особым даром - "
		interpretation += "способностью находить выход из самых сложных ситуаций и вести за собой других. "
		interpretation += "Вы - проводник между мирами, между прошлым и будущим, между травмой и исцелением.\n\n"
	}
	
	// Интерпретация ловушек
	if core.IsTrapRoom(coords) {
		interpretation += "⚠️ Вы находитесь в ловушке повторяющихся паттернов. "
		interpretation += "Числа-близнецы (111, 222, 333) указывают на циклы, которые вам нужно разорвать. "
		interpretation += "Повторяющиеся ситуации - это не наказание, а приглашение к пробуждению. "
		interpretation += "Обратите внимание: что именно повторяется? Какие уроки вы не усвоили?\n\n"
	}
	
	// Интерпретация мостов
	if core.IsBridgeRoom(coords) {
		interpretation += "🌈 Вы находитесь в точке перехода. "
		interpretation += "То, что казалось тупиком, оказывается мостом в новую реальность. "
		interpretation += "Сейчас важны доверие и готовность отпустить старое.\n\n"
	}
	
	return interpretation
}

// FormatSpiritualProfile форматирует духовный профиль для вывода
func FormatSpiritualProfile(profile *SpiritualProfile) string {
	var result strings.Builder
	
	result.WriteString("🧘 ДУХОВНЫЙ ПРОФИЛЬ\n")
	result.WriteString("==================\n\n")
	
	result.WriteString(fmt.Sprintf("Архетип души: %s\n", profile.Archetype))
	result.WriteString(fmt.Sprintf("Духовный путь: %s\n", profile.Path))
	result.WriteString(fmt.Sprintf("Уровень интеграции: %.0f%%\n\n", profile.IntegrationLevel*100))
	
	result.WriteString("📊 Духовные измерения:\n")
	for dim, val := range profile.Dimensions {
		bar := strings.Repeat("█", int(val*10))
		space := strings.Repeat("░", 10-int(val*10))
		result.WriteString(fmt.Sprintf("  %s: %s%s %.0f%%\n", 
			dim, bar, space, val*100))
	}
	
	if profile.Mission != nil {
		result.WriteString(fmt.Sprintf("\n🎯 Миссия души: %s\n", profile.Mission.PrimaryMission))
		if len(profile.Mission.SecondaryMissions) > 0 {
			result.WriteString("  Дополнительные задачи:\n")
			for _, m := range profile.Mission.SecondaryMissions {
				result.WriteString(fmt.Sprintf("    • %s\n", m))
			}
		}
	}
	
	if len(profile.KarmicLessons) > 0 {
		result.WriteString("\n🔄 Кармические уроки:\n")
		for _, lesson := range profile.KarmicLessons {
			status := "🟡 в процессе"
			if lesson.IsLearned {
				status = "🟢 пройден"
			}
			result.WriteString(fmt.Sprintf("  • %s [%s] (прогресс %.0f%%)\n", 
				lesson.Description, status, lesson.Progress*100))
		}
	}
	
	if profile.Ancestral != nil && len(profile.Ancestral.Gifts) > 0 {
		result.WriteString("\n🎁 Дары рода:\n")
		for _, gift := range profile.Ancestral.Gifts {
			result.WriteString(fmt.Sprintf("  • %s\n", gift))
		}
	}
	
	if len(profile.ShadowAspects) > 0 {
		result.WriteString("\n🌑 Теневые аспекты для интеграции:\n")
		for _, shadow := range profile.ShadowAspects {
			result.WriteString(fmt.Sprintf("  • %s\n", shadow))
		}
	}
	
	if len(profile.Rituals) > 0 {
		result.WriteString("\n🕯️ Рекомендуемые практики:\n")
		for _, ritual := range profile.Rituals {
			result.WriteString(fmt.Sprintf("  • %s\n", ritual))
		}
	}
	
	return result.String()
}

// FormatSpiritualFamilyReport форматирует духовный отчет семьи
func FormatSpiritualFamilyReport(result *SpiritualAnalysisResult) string {
	var report strings.Builder
	
	report.WriteString("🌟 ДУХОВНЫЙ ОТЧЕТ СЕМЬИ 🌟\n")
	report.WriteString("==========================\n\n")
	
	report.WriteString(fmt.Sprintf("Коллективный архетип семьи: %s\n", result.CollectiveArchetype))
	report.WriteString(fmt.Sprintf("Коллективная миссия: %s\n", result.CollectiveMission))
	report.WriteString(fmt.Sprintf("Духовный возраст семьи: %.0f%%\n\n", result.SpiritualAge*100))
	
	report.WriteString("👥 Духовные профили членов семьи:\n")
	for id, profile := range result.Profiles {
		memberName := id // в реальном коде нужно получать имя
		report.WriteString(fmt.Sprintf("  %s: %s (интеграция %.0f%%)\n", 
			memberName, profile.Archetype, profile.IntegrationLevel*100))
	}
	
	if len(result.AncestralPatterns) > 0 {
		report.WriteString("\n🔄 Родовые паттерны:\n")
		for _, pattern := range result.AncestralPatterns {
			report.WriteString(fmt.Sprintf("  • %s\n", pattern))
		}
	}
	
	if len(result.KarmicWeb) > 0 {
		report.WriteString(fmt.Sprintf("\n⚖️ Кармическая сеть: %d активных уроков\n", len(result.KarmicWeb)))
	}
	
	if len(result.Warnings) > 0 {
		report.WriteString("\n⚠️ Предупреждения:\n")
		for _, w := range result.Warnings {
			report.WriteString(fmt.Sprintf("  %s\n", w))
		}
	}
	
	if len(result.Recommendations) > 0 {
		report.WriteString("\n💡 Рекомендации:\n")
		for _, r := range result.Recommendations {
			report.WriteString(fmt.Sprintf("  • %s\n", r))
		}
	}
	
	if len(result.Rituals) > 0 {
		report.WriteString("\n🕯️ Семейные ритуалы:\n")
		for _, r := range result.Rituals {
			report.WriteString(fmt.Sprintf("  • %s\n", r))
		}
	}
	
	return report.String()
}
