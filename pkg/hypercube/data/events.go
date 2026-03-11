package data

import (
	"time"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/core"
)

// LifeEvent представляет значимое событие в жизни
type LifeEvent struct {
	Date        time.Time
	EventType   string
	Description string
	WithPerson  string
	Impact      float64
}

// Disease представляет заболевание
type Disease struct {
	Name        string
	Severity    float64
	Type        string
	DiagnosisDate *time.Time
	FromParent  string
	IsGenetic   bool
}

// Infidelity представляет измену
type Infidelity struct {
	Date          time.Time
	PartnerAtTime string
	WithPerson    string
	WasDiscovered bool
	Impact        float64
	Description   string
}

// KarmicDebt представляет кармический долг
type KarmicDebt struct {
	CreditorID  string
	DebtorID    string
	Nature      string
	Weight      float64
	IsResolved  bool
	Description string
}

// ShamePattern описывает паттерн стыда
type ShamePattern struct {
	Pattern     string
	Score       float64
	BodyZone    string
	Triggers    []string
}

// AddLifeEvent добавляет событие в историю
func (m *ExtendedFamilyMember) AddLifeEvent(eventType, description string, dateStr string, withPerson string, impact float64) error {
	date, err := core.ParseDate(dateStr)
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
	date, err := core.ParseDate(dateStr)
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
