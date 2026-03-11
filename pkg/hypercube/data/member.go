package data

import (
	"time"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/core"
)

// FamilyMember представляет члена семьи с базовыми данными
type FamilyMember struct {
	ID             string
	Name           string
	Gender         core.Gender
	Role           core.FamilyRole
	BirthDate      time.Time
	DeathDate      *time.Time
	Events         []LifeEvent
	Parents        []string
	Children       []string
	Partners       []string
	CurrentPartner string
	IsDeceased     bool
	Generation     int
	BirthOrder     int
	TotalSiblings  int
}

// ExtendedFamilyMember расширяет базовые данные дополнительной информацией
type ExtendedFamilyMember struct {
	FamilyMember
	
	AbortionCount    int
	MiscarriageCount int
	LostChildren     []LostChild
	
	IsReplacement    bool
	ReplacesWho      string
	ReplacementData  *ReplacementChildData
	
	Diseases         []Disease
	Infidelities     []Infidelity
	KarmicDebts      []KarmicDebt
	
	ShamePattern     ShamePattern
	EmotionalStability float64
	TraumaLoad       float64
	ConsciousnessLevel float64
	IdentityDiffusion float64
	SurvivorGuilt    float64
}

// CreateFamilyMember создает нового члена семьи
func CreateFamilyMember(id, name string, gender core.Gender, role core.FamilyRole, birthDateStr string) (*ExtendedFamilyMember, error) {
	birthDate, err := core.ParseDate(birthDateStr)
	if err != nil {
		return nil, err
	}
	
	return &ExtendedFamilyMember{
		FamilyMember: FamilyMember{
			ID:         id,
			Name:       name,
			Gender:     gender,
			Role:       role,
			BirthDate:  birthDate,
			Events:     make([]LifeEvent, 0),
			Parents:    make([]string, 0),
			Children:   make([]string, 0),
			Partners:   make([]string, 0),
		},
		Diseases:           make([]Disease, 0),
		Infidelities:       make([]Infidelity, 0),
		KarmicDebts:        make([]KarmicDebt, 0),
		LostChildren:       make([]LostChild, 0),
		EmotionalStability: 0.5,
		TraumaLoad:         0,
		ConsciousnessLevel: 0.3,
		IdentityDiffusion:  0,
		SurvivorGuilt:      0,
	}, nil
}

// GetAge возвращает возраст
func (m *ExtendedFamilyMember) GetAge() int {
	if m.IsDeceased && m.DeathDate != nil {
		return int(m.DeathDate.Sub(m.BirthDate).Hours() / 24 / 365)
	}
	return int(time.Now().Sub(m.BirthDate).Hours() / 24 / 365)
}

// HasParent проверяет родителя
func (m *ExtendedFamilyMember) HasParent(parentID string) bool {
	for _, p := range m.Parents {
		if p == parentID {
			return true
		}
	}
	return false
}

// HasChild проверяет ребенка
func (m *ExtendedFamilyMember) HasChild(childID string) bool {
	for _, c := range m.Children {
		if c == childID {
			return true
		}
	}
	return false
}

// HasPartner проверяет партнера
func (m *ExtendedFamilyMember) HasPartner(partnerID string) bool {
	for _, p := range m.Partners {
		if p == partnerID {
			return true
		}
	}
	return false
}
