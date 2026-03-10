package data

import (
	"time"
	"github.com/ruvcoindev/idealcore/hypercube/core"
)

// FamilyMember представляет члена семьи с базовыми данными
// Это минимальный набор информации, необходимый для анализа
type FamilyMember struct {
	ID             string            // Уникальный идентификатор
	Name           string            // Имя (для удобства восприятия)
	Gender         core.Gender       // Пол
	BirthDate      time.Time         // Дата рождения (основная для расчетов)
	DeathDate      *time.Time        // Дата смерти (если есть)
	Events         []LifeEvent       // Жизненные события
	Parents        []string          // ID родителей
	Children       []string          // ID детей
	Partners       []string          // ID партнеров (включая бывших)
	CurrentPartner string            // ID текущего партнера
	IsDeceased     bool              // Флаг смерти
	Generation     int               // Поколение (вычисляется автоматически)
	BirthOrder     int               // Порядок рождения среди сиблингов
	TotalSiblings  int               // Общее количество сиблингов
}

// ExtendedFamilyMember расширяет базовые данные дополнительной информацией
// Используется в расширенном режиме анализа, когда пользователь согласен
// предоставить более глубокие, интимные данные для точности прогноза
type ExtendedFamilyMember struct {
	FamilyMember                       // Встраиваем базовые поля
	
	// Данные о потерянных детях (аборты, выкидыши, смерти)
	// Количество абортов указывается без дат - только статистика для кармического анализа
	AbortionCount    int               // Количество абортов (суммарно, без дат)
	MiscarriageCount int               // Количество выкидышей
	LostChildren     []LostChild       // Детальная информация о потерянных детях (если доступна)
	
	// Данные о замещении
	IsReplacement    bool              // Является ли замещающим ребенком
	ReplacesWho      string            // ID потерянного ребенка, которого замещает
	ReplacementData  *ReplacementChildData // Детальные данные о замещении
	
	// Медицинские и психологические данные
	Diseases         []Disease         // Наследственные заболевания
	Infidelities     []Infidelity      // Измены (для анализа кармических паттернов)
	KarmicDebts      []KarmicDebt      // Кармические долги (вычисляются)
	
	// Психологический профиль (вычисляется)
	ShamePattern     ShamePattern      // Паттерн стыда (сжигание, заморозка, бегство)
	EmotionalStability float64         // Эмоциональная стабильность (0-1)
	TraumaLoad       float64            // Общая нагрузка травмы (0-1)
	ConsciousnessLevel float64          // Уровень осознанности (0-1)
	IdentityDiffusion float64           // Диффузия идентичности (для замещающих)
	SurvivorGuilt    float64            // Вина выжившего (0-1)
}

// CreateFamilyMember создает нового члена семьи с минимальными данными
func CreateFamilyMember(id, name string, gender core.Gender, birthDateStr string) (*ExtendedFamilyMember, error) {
	birthDate, err := core.ParseDate(birthDateStr)
	if err != nil {
		return nil, err
	}
	
	return &ExtendedFamilyMember{
		FamilyMember: FamilyMember{
			ID:         id,
			Name:       name,
			Gender:     gender,
			BirthDate:  birthDate,
			Events:     make([]LifeEvent, 0),
			Parents:    make([]string, 0),
			Children:   make([]string, 0),
			Partners:   make([]string, 0),
		},
		Diseases:     make([]Disease, 0),
		Infidelities: make([]Infidelity, 0),
		KarmicDebts:  make([]KarmicDebt, 0),
		LostChildren: make([]LostChild, 0),
	}, nil
}

// GetAge возвращает возраст на текущий момент или на момент смерти
func (m *ExtendedFamilyMember) GetAge() int {
	if m.IsDeceased && m.DeathDate != nil {
		return int(m.DeathDate.Sub(m.BirthDate).Hours() / 24 / 365)
	}
	return int(time.Now().Sub(m.BirthDate).Hours() / 24 / 365)
}

// HasParent проверяет, является ли указанный человек родителем
func (m *ExtendedFamilyMember) HasParent(parentID string) bool {
	for _, p := range m.Parents {
		if p == parentID {
			return true
		}
	}
	return false
}

// HasChild проверяет, является ли указанный человек ребенком
func (m *ExtendedFamilyMember) HasChild(childID string) bool {
	for _, c := range m.Children {
		if c == childID {
			return true
		}
	}
	return false
}

// HasPartner проверяет, является ли указанный человек партнером (текущим или бывшим)
func (m *ExtendedFamilyMember) HasPartner(partnerID string) bool {
	for _, p := range m.Partners {
		if p == partnerID {
			return true
		}
	}
	return false
}
