package analysis

import (
	"fmt"
	"github.com/ruvcoindev/idealcore/hypercube/core"
	"github.com/ruvcoindev/idealcore/hypercube/model"
)

// LargeFamilyAnalysisResult содержит результаты анализа большой семьи
type LargeFamilyAnalysisResult struct {
	TotalChildren     int
	BirthOrderStats   map[int]string      // Порядок рождения и роли
	RoleDistribution  map[core.TraumaRole]int
	ParentifiedCount  int
	LostCount         int
	InvisibleCount    int
	Warnings          []string
	Recommendations   []string
}

// AnalyzeLargeFamily анализирует семьи с большим количеством детей (10+)
func AnalyzeLargeFamily(fs *model.FamilySystem) *LargeFamilyAnalysisResult {
	// Находим семьи с 10+ детьми
	largeFamilies := make([]string, 0)
	for _, member := range fs.Members {
		if len(member.Children) >= 10 {
			largeFamilies = append(largeFamilies, member.ID)
		}
	}
	
	if len(largeFamilies) == 0 {
		return &LargeFamilyAnalysisResult{
			Warnings: []string{"Большие семьи (10+ детей) не обнаружены"},
		}
	}
	
	// Анализируем первую найденную большую семью (для простоты)
	parentID := largeFamilies[0]
	parent := fs.Members[parentID]
	children := parent.Children
	
	result := &LargeFamilyAnalysisResult{
		TotalChildren:    len(children),
		BirthOrderStats:  make(map[int]string),
		RoleDistribution: make(map[core.TraumaRole]int),
		Warnings:         make([]string, 0),
		Recommendations:  make([]string, 0),
	}
	
	// Анализируем каждого ребенка
	for i, childID := range children {
		child := fs.Members[childID]
		if child == nil {
			continue
		}
		
		birthOrder := i + 1
		role := fs.TraumaRoles[childID]
		
		if role != nil {
			result.RoleDistribution[role.PrimaryRole]++
			
			// Отмечаем особые роли
			switch role.PrimaryRole {
			case core.TraumaRoleParentified:
				result.ParentifiedCount++
				result.BirthOrderStats[birthOrder] = "родифицированный"
			case core.TraumaRoleLostChild:
				result.LostCount++
				result.BirthOrderStats[birthOrder] = "потерянный"
			case core.TraumaRoleInvisible:
				result.InvisibleCount++
				result.BirthOrderStats[birthOrder] = "невидимка"
			case core.TraumaRoleGoldenChild:
				result.BirthOrderStats[birthOrder] = "любимчик"
			case core.TraumaRoleScapegoat:
				result.BirthOrderStats[birthOrder] = "козел отпущения"
			default:
				result.BirthOrderStats[birthOrder] = string(role.PrimaryRole)
			}
		}
	}
	
	// Формируем предупреждения
	if result.ParentifiedCount > 2 {
		result.Warnings = append(result.Warnings,
			"⚠️ Множество родифицированных детей - старшие лишены детства")
	}
	
	if result.LostCount > 3 {
		result.Warnings = append(result.Warnings,
			"⚠️ Много 'потерянных' детей - младшие не получают внимания")
	}
	
	if result.InvisibleCount > 2 {
		result.Warnings = append(result.Warnings,
			"⚠️ Есть 'невидимки' - средние дети игнорируются")
	}
	
	// Рекомендации
	result.Recommendations = append(result.Recommendations,
		fmt.Sprintf("📋 Анализ большой семьи (%d детей):", result.TotalChildren),
		"   1. Распределение внимания между всеми детьми",
		"   2. Предотвращение родификации старших",
		"   3. Индивидуальный подход к каждому ребенку",
		"   4. Групповая терапия для сиблингов",
		"   5. Работа с 'потерянными' и 'невидимыми' детьми",
		"   6. Осознание родителями особенностей большой семьи")
	
	if result.ParentifiedCount > 0 {
		result.Recommendations = append(result.Recommendations,
			"   🟠 Терапия для родифицированных детей - возвращение детства")
	}
	
	if result.LostCount > 0 {
		result.Recommendations = append(result.Recommendations,
			"   🔵 Активация 'потерянных' детей - поиск голоса и места")
	}
	
	return result
}
