// Package psychology содержит базу психологических знаний: травмы, защиты, привязанности
package psychology

import "github.com/ruvcoindev/idealcore/internal/domain"

// TraumaDB — база знаний о травмах
type TraumaDB struct {
	Entries map[domain.TraumaType]TraumaInfo
}

// TraumaInfo — подробное описание травмы
type TraumaInfo struct {
	Type        domain.TraumaType
	NameRU      string
	Description string
	Symptoms    []string
	Origins     []string
	Defenses    []domain.DefenseMechanism
	HealingPath []string
	RelatedChakras []domain.ChakraType
}

// NewTraumaDB создает и заполняет базу травм
func NewTraumaDB() *TraumaDB {
	db := &TraumaDB{
		Entries: make(map[domain.TraumaType]TraumaInfo),
	}
	db.load()
	return db
}

func (db *TraumaDB) load() {
	// ПТСР
	db.Entries[domain.TraumaPTSD] = TraumaInfo{
		Type:        domain.TraumaPTSD,
		NameRU:      "Посттравматическое стрессовое расстройство",
		Description: "Реакция на единичное или повторяющееся травматическое событие, угрожающее жизни или целостности.",
		Symptoms: []string{
			"Флешбеки, навязчивые воспоминания",
			"Избегание триггеров",
			"Гипервозбуждение, тревожность",
			"Эмоциональное онемение",
		},
		Origins: []string{
			"Насилие (физическое, сексуальное, эмоциональное)",
			"Военные действия, катастрофы",
			"Внезапная потеря близкого",
		},
		Defenses: []domain.DefenseMechanism{
			domain.DefenseDenial,
			domain.DefenseDisplacement,
			domain.DefenseIntellectualization,
		},
		HealingPath: []string{
			"EMDR-терапия",
			"Соматическая проработка",
			"Безопасное повторное переживание в терапии",
		},
		RelatedChakras: []domain.ChakraType{
			domain.ChakraRoot,   // безопасность
			domain.ChakraHeart,  // доверие
			domain.ChakraThirdEye, // интеграция опыта
		},
	}

	// Комплексное ПТСР
	db.Entries[domain.TraumaCPTSD] = TraumaInfo{
		Type:        domain.TraumaCPTSD,
		NameRU:      "Комплексное ПТСР",
		Description: "Длительная травма развития, часто в детстве, в условиях отсутствия безопасной привязанности.",
		Symptoms: []string{
			"Нарушение регуляции эмоций",
			"Негативная самоконцепция ('я дефектный')",
			"Трудности в отношениях",
			"Диссоциация, деперсонализация",
		},
		Origins: []string{
			"Нарциссическое воспитание",
			"Эмоциональное пренебрежение",
			"Хроническое насилие в семье",
		},
		Defenses: []domain.DefenseMechanism{
			domain.DefenseProjection,
			domain.DefenseRationalization,
			domain.DefenseTriangulation,
		},
		HealingPath: []string{
			"Терапия, фокусированная на привязанности",
			"Работа с внутренним ребенком",
			"Развитие самосострадания",
		},
		RelatedChakras: []domain.ChakraType{
			domain.ChakraRoot,   // базовая безопасность
			domain.ChakraSacral, // эмоции
			domain.ChakraHeart,  // принятие себя
		},
	}

	// Травма привязанности
	db.Entries[domain.TraumaAttachment] = TraumaInfo{
		Type:        domain.TraumaAttachment,
		NameRU:      "Травма привязанности",
		Description: "Нарушение формирования безопасной привязанности в раннем детстве.",
		Symptoms: []string{
			"Страх брошенности или поглощения",
			"Трудности с доверием",
			"Циклы сближения-отдаления в отношениях",
			"Проекция родительских фигур на партнеров",
		},
		Origins: []string{
			"Непредсказуемость опекуна",
			"Эмоциональная недоступность родителей",
			"Ранняя сепарация, потеря",
		},
		Defenses: []domain.DefenseMechanism{
			domain.DefenseProjection,
			domain.DefenseReactionFormation,
			domain.DefenseRegression,
		},
		HealingPath: []string{
			"Корректирующий эмоциональный опыт",
			"Осознание паттернов привязанности",
			"Практика безопасной уязвимости",
		},
		RelatedChakras: []domain.ChakraType{
			domain.ChakraRoot,   // базовое доверие
			domain.ChakraHeart,  // способность любить
			domain.ChakraThroat, // выражение потребностей
		},
	}

	// Травма от нарциссического родителя
	db.Entries[domain.TraumaNarcissisticParent] = TraumaInfo{
		Type:        domain.TraumaNarcissisticParent,
		NameRU:      "Травма от нарциссического родителя",
		Description: "Воспитание родителем с нарциссическими чертами, где ребенок используется для регуляции самооценки родителя.",
		Symptoms: []string{
			"Чувство, что любовь нужно заслужить",
			"Размытые границы, трудности с 'нет'",
			"Гиперответственность за чувства других",
			"Страх быть 'эгоистичным' при заботе о себе",
			"Путаница: любовь = жертва / контроль",
		},
		Origins: []string{
			"Родитель, который любит условно",
			"Триангуляция: 'только я тебя понимаю'",
			"Обесценивание потребностей ребенка",
			"Проекция стыда родителя на ребенка",
		},
		Defenses: []domain.DefenseMechanism{
			domain.DefenseNarcissisticDevaluation, // усвоенное от родителя
			domain.DefenseTriangulation,
			domain.DefenseIntellectualization,
			domain.DefenseSublimation, // "я спасу других"
		},
		HealingPath: []string{
			"Сепарация: эмоциональная и/или физическая",
			"Работа со стыдом и виной",
			"Разрешение на эгоизм (здоровый)",
			"Построение отношений с четкими границами",
		},
		RelatedChakras: []domain.ChakraType{
			domain.ChakraSolar,  // воля, право на свои желания
			domain.ChakraHeart,  // любовь без условий
			domain.ChakraThroat, // право говорить 'нет'
		},
	}

	// ... можно добавить BPD, NPD, Abandonment и др. по аналогии
}

// Get возвращает информацию о травме
func (db *TraumaDB) Get(t domain.TraumaType) (TraumaInfo, bool) {
	info, ok := db.Entries[t]
	return info, ok
}

// GetBySymptoms находит травмы по симптомам (упрощенный поиск)
func (db *TraumaDB) GetBySymptoms(symptoms []string) []TraumaInfo {
	var results []TraumaInfo
	for _, info := range db.Entries {
		for _, sym := range symptoms {
			for _, entrySym := range info.Symptoms {
				if contains(entrySym, sym) {
					results = append(results, info)
					break
				}
			}
		}
	}
	return results
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		containsHelper(s, substr))))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
