package data

import (
	"time"
)

// LifeEvent представляет значимое событие в жизни человека
// События влияют на траекторию движения и могут быть триггерами травм
type LifeEvent struct {
	Date        time.Time // Дата события
	EventType   string    // Тип: "birth", "death", "marriage", "divorce", "civil_union", "separation", "trauma", "achievement", "loss"
	Description string    // Описание (для пользователя)
	WithPerson  string    // ID другого человека (если событие связано с кем-то)
	Impact      float64   // Влияние на жизнь (0-1) - вычисляется или указывается пользователем
}

// Disease представляет наследственное или хроническое заболевание
// Учитывается при расчете кармической нагрузки и травм
type Disease struct {
	Name        string    // Название заболевания
	Severity    float64   // Тяжесть (0-1): 0.1 - легкое, 0.5 - среднее, 1.0 - критическое
	Type        string    // Тип: "mental" (психическое), "physical" (физическое), "chronic" (хроническое), "genetic" (генетическое)
	DiagnosisDate *time.Time // Дата диагностики (если известна)
	FromParent  string    // ID родителя, от которого унаследовано (если известно)
	IsGenetic   bool      // Является ли генетическим/наследственным
}

// Infidelity представляет измену в отношениях
// Учитывается для кармического анализа и паттернов повторения
type Infidelity struct {
	Date          time.Time // Дата измены
	PartnerAtTime string    // ID партнера, которому изменили
	WithPerson    string    // ID человека, с которым изменили
	WasDiscovered bool      // Было ли раскрыто
	Impact        float64   // Влияние на отношения (0-1)
	Description   string    // Описание (для контекста)
}

// KarmicDebt представляет кармический долг между членами семьи
// Вычисляется автоматически на основе паттернов повторения
type KarmicDebt struct {
	CreditorID  string    // ID того, кому должны (кто пострадал)
	DebtorID    string    // ID того, кто должен (кто нанес вред)
	Nature      string    // Природа долга: "replacement" (замещение), "guilt" (вина), "shame" (стыд), "obligation" (обязательство), "betrayal" (предательство)
	Weight      float64   // Вес долга (0-1)
	IsResolved  bool      // Отработан ли долг
	Description string    // Описание паттерна
	OriginEvent *LifeEvent // Событие, породившее долг (если есть)
}

// ShamePattern описывает паттерн стыда, характерный для человека
// Стыд может проявляться по-разному в теле и поведении
type ShamePattern struct {
	Pattern     string   // Тип: "burning" (сжигание), "freezing" (заморозка), "fleeing" (бегство)
	Score       float64  // Уровень стыда (0-100)
	BodyZone    string   // Зона тела: "heart" (сердце), "throat" (горло), "gut" (живот), "head" (голова), "whole" (все тело)
	Triggers    []string // События или ситуации, вызывающие стыд
}

// AddLifeEvent добавляет событие в историю
func (m *ExtendedFamilyMember) AddLifeEvent(eventType, description string, dateStr string, withPerson string, impact float64) error {
	date, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return err
	}
	
	event := LifeEvent{
		Date:        date,
		EventType:   eventType,
		Description: description,
		WithPerson:  withPerson,
		Impact:      impact,
	}
	
	m.Events = append(m.Events, event)
	
	// Автоматические обновления в зависимости от типа события
	switch eventType {
	case "marriage", "civil_union":
		if withPerson != "" && !m.HasPartner(withPerson) {
			m.Partners = append(m.Partners, withPerson)
			m.CurrentPartner = withPerson
		}
	case "divorce", "separation":
		m.CurrentPartner = ""
	}
	
	return nil
}

// AddDisease добавляет заболевание
func (m *ExtendedFamilyMember) AddDisease(name string, severity float64, diseaseType string, isGenetic bool, fromParent string) {
	disease := Disease{
		Name:      name,
		Severity:  severity,
		Type:      diseaseType,
		IsGenetic: isGenetic,
		FromParent: fromParent,
	}
	
	now := time.Now()
	disease.DiagnosisDate = &now
	
	m.Diseases = append(m.Diseases, disease)
}

// AddInfidelity добавляет измену
func (m *ExtendedFamilyMember) AddInfidelity(dateStr string, partnerAtTime string, withPerson string, wasDiscovered bool, impact float64, description string) error {
	date, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return err
	}
	
	infidelity := Infidelity{
		Date:          date,
		PartnerAtTime: partnerAtTime,
		WithPerson:    withPerson,
		WasDiscovered: wasDiscovered,
		Impact:        impact,
		Description:   description,
	}
	
	m.Infidelities = append(m.Infidelities, infidelity)
	
	// Добавляем кармический долг
	debt := KarmicDebt{
		DebtorID:    m.ID,
		CreditorID:  partnerAtTime,
		Nature:      "betrayal",
		Weight:      impact * 0.8,
		IsResolved:  false,
		Description: "Измена в отношениях",
	}
	
	m.KarmicDebts = append(m.KarmicDebts, debt)
	
	return nil
}
