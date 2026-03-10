package protocol

import (
	"fmt"
	"strings"
	"github.com/ruvcoindev/idealcore/hypercube/core"
	"github.com/ruvcoindev/idealcore/hypercube/model"
	"github.com/ruvcoindev/idealcore/hypercube/analysis"
)

// FamilyReport содержит полный отчет по семье
type FamilyReport struct {
	System          *model.FamilySystem
	MemberCount     int
	Generations     int
	TraumaRoles     map[string]core.TraumaRole
	Patterns        []model.FamilyPattern
	AverageTrauma   float64
	KarmicLoad      float64
	Warnings        []string
	Recommendations []string
}

// GenerateFamilyReport создает полный отчет по семейной системе
func GenerateFamilyReport(fs *model.FamilySystem) string {
	var report strings.Builder
	
	report.WriteString("========================================\n")
	report.WriteString("СЕМЕЙНЫЙ ОТЧЕТ - HyperCube Analysis v3.1\n")
	report.WriteString("========================================\n\n")
	
	// Основная статистика
	report.WriteString(fmt.Sprintf("📊 ОСНОВНАЯ СТАТИСТИКА:\n"))
	report.WriteString(fmt.Sprintf("   Всего членов семьи: %d\n", len(fs.Members)))
	report.WriteString(fmt.Sprintf("   Режим расширенного анализа: %v\n", fs.ExtendedMode))
	
	// Поколения
	generations := make(map[int][]string)
	for id, m := range fs.Members {
		generations[m.Generation] = append(generations[m.Generation], m.Name)
	}
	
	report.WriteString(fmt.Sprintf("   Поколений: %d\n", len(generations)))
	for g := 1; g <= 4; g++ {
		if members, ok := generations[g]; ok {
			report.WriteString(fmt.Sprintf("      Поколение %d: %d человек\n", g, len(members)))
		}
	}
	report.WriteString("\n")
	
	// Травматические роли
	report.WriteString("🧬 ВЫЯВЛЕННЫЕ ТРАВМАТИЧЕСКИЕ РОЛИ:\n")
	for id, role := range fs.TraumaRoles {
		member := fs.Members[id]
		if member != nil {
			report.WriteString(fmt.Sprintf("   %s: %s (уверенность %.0f%%)\n",
				member.Name, role.PrimaryRole, role.Confidence*100))
			for _, evidence := range role.Evidence {
				report.WriteString(fmt.Sprintf("      • %s\n", evidence))
			}
			report.WriteString("\n")
		}
	}
	
	// Семейные паттерны
	if len(fs.Patterns) > 0 {
		report.WriteString("🔄 ВЫЯВЛЕННЫЕ СЕМЕЙНЫЕ ПАТТЕРНЫ:\n")
		for _, p := range fs.Patterns {
			report.WriteString(fmt.Sprintf("   • %s: %s (тяжесть %.0f%%)\n",
				p.PatternType, p.Description, p.Severity*100))
			
			// Имена участников
			names := make([]string, 0, len(p.Members))
			for _, mid := range p.Members {
				if m, ok := fs.Members[mid]; ok {
					names = append(names, m.Name)
				}
			}
			if len(names) > 0 {
				report.WriteString(fmt.Sprintf("     Участники: %s\n", strings.Join(names, ", ")))
			}
			
			for _, rec := range p.Recommendations {
				report.WriteString(fmt.Sprintf("     - %s\n", rec))
			}
			report.WriteString("\n")
		}
	}
	
	// Анализ замещающих детей
	replacementAnalysis := analysis.AnalyzeReplacementDynamics(fs)
	if replacementAnalysis.TotalCount > 0 {
		report.WriteString("👶 АНАЛИЗ ЗАМЕЩАЮЩИХ ДЕТЕЙ:\n")
		report.WriteString(fmt.Sprintf("   Найдено паттернов: %d\n", replacementAnalysis.TotalCount))
		report.WriteString(fmt.Sprintf("   Средняя нагрузка: %.0f%%\n", replacementAnalysis.AverageBurden*100))
		
		for _, warning := range replacementAnalysis.Warnings {
			report.WriteString(fmt.Sprintf("   ⚠️ %s\n", warning))
		}
		
		for _, rec := range replacementAnalysis.Recommendations {
			report.WriteString(fmt.Sprintf("   • %s\n", rec))
		}
		report.WriteString("\n")
	}
	
	// Анализ больших семей
	largeFamilyAnalysis := analysis.AnalyzeLargeFamily(fs)
	if largeFamilyAnalysis.TotalChildren >= 10 {
		report.WriteString("🏠 АНАЛИЗ БОЛЬШОЙ СЕМЬИ:\n")
		report.WriteString(fmt.Sprintf("   Количество детей: %d\n", largeFamilyAnalysis.TotalChildren))
		report.WriteString(fmt.Sprintf("   Родифицированных: %d\n", largeFamilyAnalysis.ParentifiedCount))
		report.WriteString(fmt.Sprintf("   Потерянных: %d\n", largeFamilyAnalysis.LostCount))
		report.WriteString(fmt.Sprintf("   Невидимок: %d\n", largeFamilyAnalysis.InvisibleCount))
		
		for _, warning := range largeFamilyAnalysis.Warnings {
			report.WriteString(fmt.Sprintf("   ⚠️ %s\n", warning))
		}
		
		for _, rec := range largeFamilyAnalysis.Recommendations {
			report.WriteString(fmt.Sprintf("   • %s\n", rec))
		}
		report.WriteString("\n")
	}
	
	// Кармическая нагрузка
	if fs.ExtendedMode {
		report.WriteString("🌀 КАРМИЧЕСКАЯ НАГРУЗКА:\n")
		report.WriteString(fmt.Sprintf("   Общий уровень: %.0f%%\n", fs.KarmicLoad*100))
		
		// Собираем все кармические долги
		debtCount := 0
		for _, member := range fs.Members {
			debtCount += len(member.KarmicDebts)
		}
		report.WriteString(fmt.Sprintf("   Незакрытых долгов: %d\n", debtCount))
		report.WriteString("\n")
	}
	
	// Общие рекомендации
	report.WriteString("💡 ОБЩИЕ РЕКОМЕНДАЦИИ:\n")
	report.WriteString("   • Индивидуальная терапия для каждого члена семьи\n")
	report.WriteString("   • Семейная терапия для проработки общих паттернов\n")
	report.WriteString("   • Работа с родовыми сценариями\n")
	report.WriteString("   • Восстановление здоровых границ\n")
	
	if fs.GenerationalTrauma > 70 {
		report.WriteString("   • 🔴 Критический уровень межпоколенческой травмы - требуется глубокая родовая терапия\n")
	}
	
	return report.String()
}

// GenerateShortReport генерирует краткий отчет для быстрого ознакомления
func GenerateShortReport(fs *model.FamilySystem) string {
	var report strings.Builder
	
	report.WriteString("📋 КРАТКИЙ СЕМЕЙНЫЙ ОТЧЕТ\n")
	report.WriteString("========================\n\n")
	
	report.WriteString(fmt.Sprintf("Членов семьи: %d\n", len(fs.Members)))
	report.WriteString(fmt.Sprintf("Поколений: %d\n", fs.GetGenerationSpread()))
	
	// Ключевые роли
	roleSummary := make(map[core.TraumaRole]int)
	for _, role := range fs.TraumaRoles {
		roleSummary[role.PrimaryRole]++
	}
	
	report.WriteString("\nКлючевые роли:\n")
	for role, count := range roleSummary {
		if count > 0 {
			report.WriteString(fmt.Sprintf("   %s: %d\n", role, count))
		}
	}
	
	// Основные паттерны
	if len(fs.Patterns) > 0 {
		report.WriteString("\nОсновные паттерны:\n")
		for i, p := range fs.Patterns {
			if i >= 3 {
				break
			}
			report.WriteString(fmt.Sprintf("   • %s\n", p.Description))
		}
	}
	
	report.WriteString(fmt.Sprintf("\nУровень травмы: %.0f%%\n", fs.GenerationalTrauma))
	
	return report.String()
}
