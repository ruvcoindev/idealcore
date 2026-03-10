
# HyperCube Family Analysis v3.1

[![Go Version](https://img.shields.io/badge/Go-1.16+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-3.1.0-purple.svg)]()

Пакет для анализа семейных систем, травматических ролей и межпоколенческих паттернов, основанный на метафоре фильма "Куб" (Hypercube). Использует даты рождения для генерации 4-мерных координат и моделирования жизненных траекторий.

## 📖 О проекте

HyperCube — это математическая модель семейной системы, вдохновленная фильмом "Куб". Каждый человек представлен как точка в 4-мерном пространстве (X, Y, Z, W), где:

- **X** — день рождения (базовые реакции, инстинкты)
- **Y** — месяц рождения (эмоциональные паттерны)
- **Z** — год рождения (кармические/родовые программы)
- **W** — контрольная сумма (уровень осознанности, "мост" между измерениями)

Движение по гиперкубу моделирует жизненный путь, а комнаты-ловушки (111, 222, 333) символизируют повторяющиеся травматические паттерны. Комнаты-мосты (999) — точки выхода и исцеления.

## ✨ Ключевые возможности

### 🧬 Выявление травматических ролей
Автоматическое определение ролей в дисфункциональной семье:
- **Golden Child** (любимчик) — идеализируемый ребенок, страдающий от перфекционизма
- **Scapegoat** (козел отпущения) — ребенок, на которого проецируются все проблемы
- **Lost Child** (потерянный ребенок) — эмоционально отстраненный, незаметный
- **Invisible Child** (невидимка) — чьи потребности игнорируются
- **Glass Child** (стеклянный ребенок) — невидимый на фоне проблемного сиблинга
- **Parentified Child** (родифицированный) — вынужденный выполнять родительские функции
- **Mascot** (шут) — снимающий напряжение юмором
- **Emotional Spouse** (эмоциональный супруг) — заменяющий партнера родителю
- **Replacement Child** (замещающий ребенок) — рожденный после потери предыдущего
- **Ghost** (призрак) — потерянный ребенок, влияющий на семью
- **Ancestor** (предок) — носитель родовой травмы

### 👪 Анализ всех типов отношений
- **Родитель-ребенок**: мать-дочь, мать-сын, отец-дочь, отец-сын
- **Сиблинги**: сестра-сестра, брат-брат, сестра-брат
- **Межпоколенческие**: бабушка-внучка, дедушка-внук, бабушка-внук, дедушка-внучка
- **Партнерские**: супруги, гражданский брак, бывшие, романтические
- **Социальные**: дружеские, рабочие

### 🔄 Анализ замещающих детей и потерь
- Учет абортов, выкидышей, мертворождений, детских смертей
- Выявление паттернов замещения
- Расчет нагрузки на замещающего ребенка
- Межпоколенческое замещение

### 📊 Математические расчеты
- 4-мерные координаты из даты рождения
- Векторы движения и их амплитуда
- Совместимость двух людей (общие оси, синхронность, кармический фактор)
- Групповая совместимость (3+ человек)
- Расчет телесного стыда с учетом роли и событий

### 💔 Протоколы разделения
- Индивидуальные планы для разных типов отношений
- Этапы разделения с симптомами и вероятностью выживания
- Правила контакта
- Рекомендации по терапии

### 🏠 Анализ больших семей
- Специфическая динамика при 10+ детях
- Выявление родифицированных, потерянных и невидимых детей
- Распределение ролей по порядку рождения

## 🚀 Установка

```bash
go get github.com/idealcore/hypercube

📋 Быстрый старт
Минимальный пример

``` go
package main

import (
    "fmt"
    "github.com/ruvcoindev/idealcore/hypercube"
)

func main() {
    // Создаем семейную систему (базовый режим)
    fs := hypercube.NewFamilySystem(false)
    
    // Создаем членов семьи
    mother, _ := hypercube.CreateFamilyMember("m1", "Анна", hypercube.GenderFemale, "15.05.1970")
    daughter, _ := hypercube.CreateFamilyMember("d1", "Мария", hypercube.GenderFemale, "15.06.1996")
    
    // Устанавливаем связи
    daughter.Parents = []string{"m1"}
    mother.Children = []string{"d1"}
    
    // Добавляем в систему
    fs.AddMember(mother)
    fs.AddMember(daughter)
    
    // Выявляем травматические роли
    motherRole := fs.IdentifyTraumaRoles("m1")
    daughterRole := fs.IdentifyTraumaRoles("d1")
    
    fmt.Printf("Мать: %s\n", motherRole.PrimaryRole)
    fmt.Printf("Дочь: %s\n", daughterRole.PrimaryRole)
    
    // Анализ совместимости
    coords1 := hypercube.ParseDateToCoords(mother.BirthDate)
    coords2 := hypercube.ParseDateToCoords(daughter.BirthDate)
    vectors1 := hypercube.CalculateVectors(coords1)
    vectors2 := hypercube.CalculateVectors(coords2)
    
    compat := hypercube.CalculateCompatibility(coords1, coords2, vectors1, vectors2)
    fmt.Printf("Совместимость: %.1f%% (%s)\n", compat.Score*100, compat.Level)
    
    // Полный отчет
    report := fs.GenerateFamilyReport()
    fmt.Println(report)
}

Расширенный режим (с интимными данными)

```go

// Включаем расширенный режим для учета абортов, измен, потерь
fs := hypercube.NewFamilySystem(true)

// Добавляем информацию об абортах (только количество, без дат)
mother.AbortionCount = 2
mother.MiscarriageCount = 1

// Добавляем потерянного ребенка
lostID := "lost_1"
mother.AddLostChild(hypercube.LossTypeMiscarriage, "10.08.1995", 12, "")

// Помечаем дочь как замещающего ребенка
daughter.MarkAsReplacement(lostID, "01.01.1996", false, true)

// Добавляем измену
daughter.AddInfidelity("15.03.2020", "p2", "p1", true, 0.8, "Случайная связь")

// Добавляем заболевание
daughter.AddDisease("Депрессия", 0.6, "mental", false, "")

// Анализ замещающих детей
replacementPatterns := hypercube.FindReplacementPatterns(fs)
for _, p := range replacementPatterns {
    fmt.Printf("Замещение: разрыв %d мес, нагрузка %.0f%%\n", 
        p.GapMonths, p.Severity*100)
}

📚 Документация по модулям
Пакет core — базовые математические функции

```go

// Координаты в гиперкубе
type HypercubeCoords struct {
    X, Y, Z, W int32
}

// Преобразование даты в координаты
coords := hypercube.ParseDateToCoords(birthDate)

// Вычисление векторов
vectors := hypercube.CalculateVectors(coords)

// Перемещение по гиперкубу
newCoords := hypercube.MoveRoom(coords, vectors, step)

// Проверка комнат
isTrap := hypercube.IsTrapRoom(coords)  // ловушка (111, 222...)
isBridge := hypercube.IsBridgeRoom(coords) // мост (999)

// Поиск безопасных комнат
safeRooms := hypercube.FindSafeRooms(coords, 12)

Пакет data — структуры данных

```go

// Создание члена семьи
member, err := hypercube.CreateFamilyMember("id1", "Имя", hypercube.GenderFemale, "15.05.1970")

// Добавление жизненных событий
member.AddLifeEvent("marriage", "Свадьба", "15.06.1995", "partner_id", 0.5)

// Добавление потери (аборт, выкидыш)
member.AddLostChild(hypercube.LossTypeMiscarriage, "10.08.1995", 12, "father_id")

// Добавление заболевания
member.AddDisease("Диабет", 0.4, "chronic", true, "parent_id")

// Добавление измены
member.AddInfidelity("15.03.2020", "partner_id", "other_id", true, 0.8, "Описание")

Пакет model — семейная система и выявление ролей

```go

// Создание системы
fs := hypercube.NewFamilySystem(true)

// Добавление членов
fs.AddMember(member1)
fs.AddMember(member2)

// Выявление травматических ролей
role := fs.IdentifyTraumaRoles("id1")
fmt.Printf("Роль: %s (уверенность %.0f%%)\n", role.PrimaryRole, role.Confidence*100)
for _, evidence := range role.Evidence {
    fmt.Printf("  • %s\n", evidence)
}

// Поиск семейных паттернов
patterns := fs.FindFamilyPatterns()
for _, p := range patterns {
    fmt.Printf("Паттерн: %s\n", p.Description)
    fmt.Printf("Тяжесть: %.0f%%\n", p.Severity*100)
}

Пакет analysis — аналитические функции

```go 

// Анализ совместимости двух людей
compat := hypercube.CalculateCompatibility(coords1, coords2, vectors1, vectors2)
fmt.Printf("Совместимость: %.1f%% (%s)\n", compat.Score*100, compat.Level)

// Анализ совместимости группы
groupCompat := hypercube.CalculateGroupCompatibility(members)
fmt.Printf("Средняя совместимость: %.1f%%\n", groupCompat.AverageScore*100)

// Анализ стыда
shame := hypercube.CalculateBodyShame(coords, role.PrimaryRole, 30, member)
fmt.Printf("Уровень стыда: %.1f\n", shame)

// Анализ замещающих детей
replacementAnalysis := hypercube.AnalyzeReplacementDynamics(fs)
fmt.Printf("Найдено паттернов: %d\n", replacementAnalysis.TotalCount)

// Анализ большой семьи
largeFamilyAnalysis := hypercube.AnalyzeLargeFamily(fs)
if largeFamilyAnalysis.TotalChildren >= 10 {
    fmt.Printf("Большая семья: %d детей\n", largeFamilyAnalysis.TotalChildren)
}

Пакет protocol — протоколы разделения и терапии

```go

// Протокол разделения
protocol := hypercube.GenerateSeparationProtocol(person1, person2, 
    hypercube.RelationshipMotherDaughter, compat)

fmt.Printf("Требуется лет: %d\n", protocol.RequiredYears)
fmt.Printf("Вероятность успеха: %.0f%%\n", protocol.SuccessProbability*100)
for _, rule := range protocol.ContactRules {
    fmt.Println(rule)
}

// План терапии
therapyPlan := hypercube.GenerateTherapyPlan(member, role)
fmt.Println(hypercube.FormatTherapyPlan(therapyPlan))

// Краткий отчет
shortReport := hypercube.GenerateShortReport(fs)
fmt.Println(shortReport)

// Полный отчет
fullReport := fs.GenerateFamilyReport()
fmt.Println(fullReport)

📊 Интерпретация результатов

Уровни совместимости

Уровень	Балл	Рекомендации
Высокая	>0.7	Хороший потенциал для здоровых отношений при проработке травм
Средняя	0.4-0.7	Требуется осознанная работа, четкие границы
Низкая	<0.4	Высокий риск деструктивных паттернов, осторожность

Уровни стыда

Уровень	        Балл	               Рекомендации
Критический	   >150	           Срочная терапия, телесные практики
Высокий	       80-150	       Терапия, работа с внутренним критиком
Средний	       40-80	       Самостоятельная работа, дневник чувств
Низкий	       <40	           Профилактика, поддерживающие практики

Паттерны стыда

Burning (сжигание) — стыд проявляется как жар, гнев. Требуется работа с гневом, охлаждающие практики.

Freezing (заморозка) — стыд вызывает оцепенение. Требуется разогревающие практики, возвращение чувствительности.

Fleeing (бегство) — стыд вызывает избегание. Требуется заземление, постепенное приближение к триггерам.


🧪 Примеры использования
Пример 1: Анализ матери и дочери

package main

import (
    "fmt"
    "github.com/idealcore/hypercube"
)

func main() {
    fs := hypercube.NewFamilySystem(false)
    
    mother, _ := hypercube.CreateFamilyMember("m1", "Елена", hypercube.GenderFemale, "10.03.1970")
    daughter, _ := hypercube.CreateFamilyMember("d1", "Ольга", hypercube.GenderFemale, "15.08.1995")
    
    daughter.Parents = []string{"m1"}
    mother.Children = []string{"d1"}
    
    fs.AddMember(mother)
    fs.AddMember(daughter)
    
    motherRole := fs.IdentifyTraumaRoles("m1")
    daughterRole := fs.IdentifyTraumaRoles("d1")
    
    fmt.Println("=== АНАЛИЗ МАТЕРИ И ДОЧЕРИ ===")
    fmt.Printf("Мать: %s (уверенность %.0f%%)\n", motherRole.PrimaryRole, motherRole.Confidence*100)
    fmt.Printf("Дочь: %s (уверенность %.0f%%)\n", daughterRole.PrimaryRole, daughterRole.Confidence*100)
    
    coords1 := hypercube.ParseDateToCoords(mother.BirthDate)
    coords2 := hypercube.ParseDateToCoords(daughter.BirthDate)
    vectors1 := hypercube.CalculateVectors(coords1)
    vectors2 := hypercube.CalculateVectors(coords2)
    
    compat := hypercube.CalculateCompatibility(coords1, coords2, vectors1, vectors2)
    fmt.Printf("\nСовместимость: %.1f%% (%s)\n", compat.Score*100, compat.Level)
    
    if compat.Score < 0.5 {
        fmt.Println("\n⚠️ Рекомендации:")
        for _, rec := range compat.Recommendations {
            fmt.Printf("  • %s\n", rec)
        }
    }
    
    patterns := fs.FindFamilyPatterns()
    if len(patterns) > 0 {
        fmt.Println("\n🔄 Выявленные паттерны:")
        for _, p := range patterns {
            fmt.Printf("  • %s\n", p.Description)
        }
    }
}

Пример 2: Анализ семьи с замещающим ребенком

```go
package main

import (
    "fmt"
    "github.com/idealcore/hypercube"
)

func main() {
    fs := hypercube.NewFamilySystem(true)
    
    mother, _ := hypercube.CreateFamilyMember("m1", "Анна", hypercube.GenderFemale, "15.05.1970")
    daughter, _ := hypercube.CreateFamilyMember("d1", "Мария", hypercube.GenderFemale, "15.06.1996")
    
    daughter.Parents = []string{"m1"}
    mother.Children = []string{"d1"}
    
    fs.AddMember(mother)
    fs.AddMember(daughter)
    
    // Добавляем информацию о потере
    mother.AddLostChild(hypercube.LossTypeMiscarriage, "10.08.1995", 10, "")
    
    // Помечаем дочь как замещающего ребенка
    daughter.MarkAsReplacement("lost_m1_10.08.1995", "01.01.1996", false, true)
    
    // Выявляем роли
    motherRole := fs.IdentifyTraumaRoles("m1")
    daughterRole := fs.IdentifyTraumaRoles("d1")
    
    fmt.Println("=== АНАЛИЗ СЕМЬИ С ЗАМЕЩАЮЩИМ РЕБЕНКОМ ===")
    fmt.Printf("Мать: %s\n", motherRole.PrimaryRole)
    fmt.Printf("Дочь: %s\n", daughterRole.PrimaryRole)
    
    // Анализ замещающих детей
    replacementAnalysis := hypercube.AnalyzeReplacementDynamics(fs)
    fmt.Printf("\n👶 Замещающих детей: %d\n", replacementAnalysis.TotalCount)
    fmt.Printf("Средняя нагрузка: %.0f%%\n", replacementAnalysis.AverageBurden*100)
    
    for _, warning := range replacementAnalysis.Warnings {
        fmt.Printf("⚠️ %s\n", warning)
    }
    
    // План терапии для дочери
    therapyPlan := hypercube.GenerateTherapyPlan(daughter, daughterRole)
    fmt.Println("\n🧑‍⚕️ ПЛАН ТЕРАПИИ ДЛЯ ДОЧЕРИ")
    fmt.Println(hypercube.FormatTherapyPlan(therapyPlan))
}

Пример 3: Анализ большой семьи

```go

package main

import (
    "fmt"
    "github.com/idealcore/hypercube"
)

func main() {
    fs := hypercube.NewFamilySystem(false)
    
    // Родители
    mother, _ := hypercube.CreateFamilyMember("m1", "Екатерина", hypercube.GenderFemale, "01.01.1960")
    father, _ := hypercube.CreateFamilyMember("f1", "Иван", hypercube.GenderMale, "15.03.1958")
    fs.AddMember(mother)
    fs.AddMember(father)
    
    // 12 детей
    for i := 1; i <= 12; i++ {
        gender := hypercube.GenderFemale
        if i%2 == 0 {
            gender = hypercube.GenderMale
        }
        child, _ := hypercube.CreateFamilyMember(
            fmt.Sprintf("c%d", i), 
            fmt.Sprintf("Ребенок %d", i), 
            gender, 
            fmt.Sprintf("%02d.%02d.%d", i, i%12+1, 1980+i),
        )
        child.Parents = []string{"m1", "f1"}
        fs.AddMember(child)
        mother.Children = append(mother.Children, child.ID)
        father.Children = append(father.Children, child.ID)
    }
    
    // Выявляем роли для всех детей
    for i := 1; i <= 12; i++ {
        fs.IdentifyTraumaRoles(fmt.Sprintf("c%d", i))
    }
    
    // Анализ большой семьи
    largeFamilyAnalysis := hypercube.AnalyzeLargeFamily(fs)
    
    fmt.Println("=== АНАЛИЗ БОЛЬШОЙ СЕМЬИ (12 ДЕТЕЙ) ===")
    fmt.Printf("Родифицированных: %d\n", largeFamilyAnalysis.ParentifiedCount)
    fmt.Printf("Потерянных: %d\n", largeFamilyAnalysis.LostCount)
    fmt.Printf("Невидимок: %d\n", largeFamilyAnalysis.InvisibleCount)
    
    fmt.Println("\n📊 Распределение ролей по порядку рождения:")
    for order, role := range largeFamilyAnalysis.BirthOrderStats {
        fmt.Printf("  %d: %s\n", order, role)
    }
    
    fmt.Println("\n⚠️ Предупреждения:")
    for _, warning := range largeFamilyAnalysis.Warnings {
        fmt.Printf("  • %s\n", warning)
    }
    
    fmt.Println("\n📋 Рекомендации:")
    for _, rec := range largeFamilyAnalysis.Recommendations {
        fmt.Printf("  • %s\n", rec)
    }
}

🧠 Философия и метафоры
Гиперкуб как модель жизни
В фильме "Куб" комнаты двигались по циклическому алгоритму, и только комната 999 была выходом. В нашей модели:

Комнаты-ловушки (111, 222, 333) — повторяющиеся травматические паттерны, циклы, из которых сложно выйти

Комнаты-мосты (999) — моменты прозрения, исцеления, точки выхода из травмы

Векторы движения — жизненные изменения, направление развития

Амплитуда векторов — интенсивность жизни, эмоциональная турбулентность

Травматические роли как комнаты
Каждый человек в дисфункциональной семье занимает определенную "комнату" — роль, которая определяет его жизненную траекторию. Эти роли не выбираются сознательно, а присваиваются семейной системой.

Замещающие дети как "призраки"
Дети, рожденные после потери предыдущего, часто живут с "призраком" — незримым присутствием того, кого они должны заменить. Это создает диффузию идентичности и глубокую вину выжившего.

🤝 Вклад в развитие
Мы приветствуем вклад в развитие проекта! Вы можете помочь:

1. Сообщая об ошибках через Issues

2. Предлагая новые паттерны для анализа

3. Улучшая документацию

4. Добавляя тесты

5. Переводя на другие языки

📄 Лицензия
MIT License. См. файл LICENSE для деталей.

🙏 Благодарности
Фильму "Куб" (Cube, 1997) и его продолжениям за вдохновение

Всем исследователям семейных систем и травм

Пользователям, предоставляющим обратную связь

