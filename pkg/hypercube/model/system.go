package model

import (
	"sort"
	"github.com/ruvcoindev/idealcore/hypercube/core"
	"github.com/ruvcoindev/idealcore/hypercube/data"
)

// FamilySystem представляет полную семейную систему
// Хранит всех членов семьи и результаты анализа
type FamilySystem struct {
	Members            map[string]*data.ExtendedFamilyMember // Все члены семьи по ID
	NarcissisticMembers []string                               // ID членов с нарциссическими чертами (выявленные)
	EnablerMembers     []string                                // ID созависимых членов (выявленные)
	TraumaRoles        map[string]*IdentifiedTraumaRoles       // Выявленные травматические роли для каждого
	GenerationalTrauma float64                                  // Общий уровень межпоколенческой травмы (0-100)
	Patterns           []FamilyPattern                          // Выявленные семейные паттерны
	ExtendedMode       bool                                      // Режим расширенного анализа (с интимными данными)
}

// IdentifiedTraumaRoles содержит выявленные травматические роли для одного человека
type IdentifiedTraumaRoles struct {
	PrimaryRole   core.TraumaRole   // Основная роль (наиболее выраженная)
	SecondaryRoles []core.TraumaRole // Дополнительные роли (могут сочетаться)
	Confidence    float64            // Уверенность в определении (0-1)
	Evidence      []string            // Доказательства/признаки, на основе которых сделаны выводы
}

// FamilyPattern представляет выявленный семейный паттерн
type FamilyPattern struct {
	PatternType     string   // Тип паттерна (например, "mother-daughter-scapegoat")
	Description     string   // Описание на понятном языке
	Members         []string // ID участников паттерна
	Severity        float64  // Тяжесть/выраженность (0-1)
	Generations     []int    // Задействованные поколения
	Recommendations []string // Рекомендации по работе с паттерном
	KarmicOrigin    *data.KarmicDebt // Кармическое происхождение (если есть)
}

// NewFamilySystem создает новую семейную систему
func NewFamilySystem(extendedMode bool) *FamilySystem {
	return &FamilySystem{
		Members:            make(map[string]*data.ExtendedFamilyMember),
		NarcissisticMembers: make([]string, 0),
		EnablerMembers:     make([]string, 0),
		TraumaRoles:        make(map[string]*IdentifiedTraumaRoles),
		Patterns:           make([]FamilyPattern, 0),
		ExtendedMode:       extendedMode,
	}
}

// AddMember добавляет члена семьи с автоматическим определением поколения
func (fs *FamilySystem) AddMember(member *data.ExtendedFamilyMember) {
	// Вычисляем поколение на основе родителей
	if len(member.Parents) > 0 {
		maxParentGen := 0
		for _, parentID := range member.Parents {
			if parent, ok := fs.Members[parentID]; ok {
				if parent.Generation > maxParentGen {
					maxParentGen = parent.Generation
				}
			}
		}
		member.Generation = maxParentGen + 1
	} else {
		// Если нет родителей, считаем первым поколением
		member.Generation = 1
	}
	
	// Добавляем в карту
	fs.Members[member.ID] = member
	
	// Обновляем связи родителей и детей
	for _, parentID := range member.Parents {
		if parent, ok := fs.Members[parentID]; ok {
			// Добавляем ребенка к родителю, если еще не добавлен
			alreadyHas := false
			for _, childID := range parent.Children {
				if childID == member.ID {
					alreadyHas = true
					break
				}
			}
			if !alreadyHas {
				parent.Children = append(parent.Children, member.ID)
			}
		}
	}
}

// GetMember возвращает члена семьи по ID
func (fs *FamilySystem) GetMember(id string) *data.ExtendedFamilyMember {
	return fs.Members[id]
}

// GetSiblings возвращает список ID всех сиблингов (братьев и сестер)
func (fs *FamilySystem) GetSiblings(memberID string) []string {
	member := fs.Members[memberID]
	if member == nil {
		return nil
	}
	
	siblings := make([]string, 0)
	siblingMap := make(map[string]bool)
	
	// Ищем через общих родителей
	for _, parentID := range member.Parents {
		if parent, ok := fs.Members[parentID]; ok {
			for _, childID := range parent.Children {
				if childID != memberID && !siblingMap[childID] {
					siblingMap[childID] = true
					siblings = append(siblings, childID)
				}
			}
		}
	}
	
	return siblings
}

// GetBirthOrder возвращает порядок рождения среди сиблингов
func (fs *FamilySystem) GetBirthOrder(memberID string) int {
	member := fs.Members[memberID]
	if member == nil {
		return 0
	}
	
	siblings := fs.GetSiblings(memberID)
	
	// Добавляем самого себя для сортировки
	allChildren := append(siblings, memberID)
	
	// Сортируем по дате рождения
	type childInfo struct {
		id   string
		date string
	}
	
	childrenList := make([]childInfo, 0, len(allChildren))
	for _, childID := range allChildren {
		if child, ok := fs.Members[childID]; ok {
			childrenList = append(childrenList, childInfo{
				id:   childID,
				date: child.BirthDate.Format("2006-01-02"),
			})
		}
	}
	
	sort.Slice(childrenList, func(i, j int) bool {
		return childrenList[i].date < childrenList[j].date
	})
	
	// Находим свою позицию
	for i, child := range childrenList {
		if child.id == memberID {
			return i + 1 // +1 потому что порядок начинается с 1
		}
	}
	
	return 0
}

// GetGenerationSpread возвращает количество поколений в семье
func (fs *FamilySystem) GetGenerationSpread() int {
	genMap := make(map[int]bool)
	for _, member := range fs.Members {
		genMap[member.Generation] = true
	}
	return len(genMap)
}

// GetMembersByGeneration возвращает членов семьи по поколению
func (fs *FamilySystem) GetMembersByGeneration(generation int) []*data.ExtendedFamilyMember {
	result := make([]*data.ExtendedFamilyMember, 0)
	for _, member := range fs.Members {
		if member.Generation == generation {
			result = append(result, member)
		}
	}
	return result
}
