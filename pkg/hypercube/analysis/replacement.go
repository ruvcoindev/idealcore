package analysis

import (
	"time"
	"github.com/ruvcoindev/idealcore/hypercube/core"
	"github.com/ruvcoindev/idealcore/hypercube/data"
)

// ReplacementPattern представляет паттерн замещения
type ReplacementPattern struct {
	LostChild           *data.LostChild
	ReplacementChild    *data.ExtendedFamilyMember
	GapMonths           int
	SimilarityScore     float64
	Generations         []int
	Severity            float64
	IsTransgenerational bool
}

// ReplacementAnalysisResult содержит результаты анализа замещающих детей
type ReplacementAnalysisResult struct {
	Patterns      []ReplacementPattern
	TotalCount    int
	AverageBurden float64
	Warnings      []string
	Recommendations []string
}

// FindReplacementPatterns выявляет паттерны замещающих детей в семье
func FindReplacementPatterns(fs *model.FamilySystem) []ReplacementPattern {
	patterns := make([]ReplacementPattern, 0)
	
	// Ищем прямые замещения (ребенок помечен как замещающий)
	for _, member := range fs.Members {
		if member.IsReplacement && member.ReplacementData != nil {
			// Ищем потерянного ребенка, которого замещает
			var lostChild *data.LostChild
			for _, parentID := range member.Parents {
				if parent, ok := fs.Members[parentID]; ok {
					for _, lc := range parent.LostChildren {
						if lc.ID == member.ReplacesWho {
							lostChild = &lc
							break
						}
					}
				}
			}
			
			if lostChild != nil {
				// Вычисляем сходство
				similarity := 0.0
				
				// По дате (родился примерно в то же время года)
				if member.BirthDate.YearDay() == lostChild.LossDate.YearDay() {
					similarity += 0.3
				} else if abs(member.BirthDate.YearDay() - lostChild.LossDate.YearDay()) < 7 {
					similarity += 0.2
				}
				
				// По полу (если известен пол потерянного)
				// В реальном коде нужно хранить пол потерянного
				
				// Явное замещение
				if member.ReplacementData.IsExplicit {
					similarity += 0.5
				}
				
				pattern := ReplacementPattern{
					LostChild:           lostChild,
					ReplacementChild:    member,
					GapMonths:           member.ReplacementData.TimeGapMonths,
					SimilarityScore:     similarity,
					Generations:         []int{member.Generation},
					Severity:            member.ReplacementData.BurdenLevel,
					IsTransgenerational: false,
				}
				patterns = append(patterns, pattern)
			}
		}
	}
	
	// Ищем межпоколенческие замещения (мать сама была замещающей, и ее ребенок тоже)
	for _, member := range fs.Members {
		if member.Gender == core.GenderFemale && member.IsReplacement && len(member.Children) > 0 {
			for _, childID := range member.Children {
				child := fs.Members[childID]
				if child != nil && child.IsReplacement {
					// Два поколения замещающих детей
					pattern := ReplacementPattern{
						LostChild:           nil,
						ReplacementChild:    child,
						GapMonths:           child.ReplacementData.TimeGapMonths,
						SimilarityScore:     0.7,
						Generations:         []int{member.Generation, child.Generation},
						Severity:            0.9,
						IsTransgenerational: true,
					}
					patterns = append(patterns, pattern)
				}
			}
		}
	}
	
	return patterns
}

// AnalyzeReplacementDynamics анализирует динамику замещения
func AnalyzeReplacementDynamics(fs *model.FamilySystem) ReplacementAnalysisResult {
	patterns := FindReplacementPatterns(fs)
	
	result := ReplacementAnalysisResult{
		Patterns:      patterns,
		TotalCount:    len(patterns),
		AverageBurden: 0,
		Warnings:      make([]string, 0),
		Recommendations: make([]string, 0),
	}
	
	if len(patterns) == 0 {
		result.Recommendations = append(result.Recommendations,
			"Паттерны замещающих детей не обнаружены")
		return result
	}
	
	// Вычисляем среднюю нагрузку
	totalBurden := 0.0
	for _, p := range patterns {
		totalBurden += p.Severity
		
		// Предупреждения для критических случаев
		if p.Severity > 0.8 {
			result.Warnings = append(result.Warnings,
				"⚠️ Критическая нагрузка на замещающего ребенка")
		}
		if p.GapMonths < 6 {
			result.Warnings = append(result.Warnings,
				"⚠️ Слишком короткий разрыв между потерей и зачатием (<6 месяцев)")
		}
		if p.IsTransgenerational {
			result.Warnings = append(result.Warnings,
				"⚠️ Межпоколенческое замещение - паттерн повторяется")
		}
	}
	
	result.AverageBurden = totalBurden / float64(len(patterns))
	
	// Общие рекомендации
	result.Recommendations = append(result.Recommendations,
		"📋 Рекомендации для работы с замещающими детьми:",
		"   1. Терапия идентичности - отделение от образа потерянного",
		"   2. Работа с виной выжившего",
		"   3. Ритуал прощания с потерянным сиблингом",
		"   4. Признание права на собственную жизнь",
		"   5. Терапия для родителей - работа с горем")
	
	if result.AverageBurden > 0.7 {
		result.Recommendations = append(result.Recommendations,
			"   🔴 Рекомендуется интенсивная терапия для всех участников")
	}
	
	return result
}

// CalculateReplacementBurden вычисляет нагрузку на замещающего ребенка
func CalculateReplacementBurden(member *data.ExtendedFamilyMember) float64 {
	if !member.IsReplacement || member.ReplacementData == nil {
		return 0
	}
	
	burden := 0.5 // Базовый уровень
	
	// Чем меньше разрыв, тем выше нагрузка
	if member.ReplacementData.TimeGapMonths < 6 {
		burden += 0.3
	} else if member.ReplacementData.TimeGapMonths < 12 {
		burden += 0.2
	} else if member.ReplacementData.TimeGapMonths < 24 {
		burden += 0.1
	}
	
	// Явное замещение повышает нагрузку
	if member.ReplacementData.IsExplicit {
		burden += 0.2
	}
	
	// Неосознанное родителями замещение тоже повышает нагрузку
	if !member.ReplacementData.ParentsAware {
		burden += 0.2
	}
	
	// Вина выжившего
	burden += member.SurvivorGuilt * 0.3
	
	// Ограничиваем до 1.0
	if burden > 1.0 {
		burden = 1.0
	}
	
	return burden
}

// abs - вспомогательная функция для модуля числа
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
