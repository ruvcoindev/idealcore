// Package hypercube реализует математику из фильма "Куб/Гиперкуб"
// с добавлением измерений для моделирования семейных систем, травм
// и межпоколенческих паттернов
package hypercube

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// CubeSize — размер куба (26 комнат по каждой оси)
const CubeSize = 26

// RelationshipType определяет тип отношений
type RelationshipType string

const (
	RelationshipMotherDaughter     RelationshipType = "mother-daughter"
	RelationshipMotherSon           RelationshipType = "mother-son"
	RelationshipFatherDaughter      RelationshipType = "father-daughter"
	RelationshipFatherSon           RelationshipType = "father-son"
	RelationshipSisterSister        RelationshipType = "sister-sister"
	RelationshipBrotherBrother      RelationshipType = "brother-brother"
	RelationshipSisterBrother       RelationshipType = "sister-brother"
	RelationshipGrandmotherGranddaughter RelationshipType = "grandmother-granddaughter"
	RelationshipGrandmotherGrandson     RelationshipType = "grandmother-grandson"
	RelationshipGrandfatherGranddaughter RelationshipType = "grandfather-granddaughter"
	RelationshipGrandfatherGrandson     RelationshipType = "grandfather-grandson"
	RelationshipHusbandWife         RelationshipType = "husband-wife"
	RelationshipCivilPartnership    RelationshipType = "civil-partnership"
	RelationshipExSpouses           RelationshipType = "ex-spouses"
	RelationshipFriendly            RelationshipType = "friendly"
	RelationshipWork                RelationshipType = "work"
)

// Gender определяет пол
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// FamilyRole определяет базовую роль в семье
type FamilyRole string

const (
	RoleMother      FamilyRole = "mother"
	RoleFather      FamilyRole = "father"
	RoleDaughter    FamilyRole = "daughter"
	RoleSon         FamilyRole = "son"
	RoleGrandmother FamilyRole = "grandmother"
	RoleGrandfather FamilyRole = "grandfather"
	RoleGranddaughter FamilyRole = "granddaughter"
	RoleGrandson    FamilyRole = "grandson"
	RoleHusband     FamilyRole = "husband"
	RoleWife        FamilyRole = "wife"
	RolePartner     FamilyRole = "partner"
)

// TraumaRole определяет роль в травматической системе (выявляется автоматически)
type TraumaRole string

const (
	TraumaRoleGoldenChild   TraumaRole = "golden_child"   // любимчик
	TraumaRoleScapegoat     TraumaRole = "scapegoat"      // козел отпущения
	TraumaRoleLostChild     TraumaRole = "lost_child"     // потерянный ребенок
	TraumaRoleInvisible     TraumaRole = "invisible"      // невидимка
	TraumaRoleGlassChild    TraumaRole = "glass_child"    // стеклянный ребенок
	TraumaRoleParentified   TraumaRole = "parentified"    // родифицированный
	TraumaRoleMascot        TraumaRole = "mascot"         // шут
	TraumaRoleEmotionalSpouse TraumaRole = "emotional_spouse" // эмоциональный супруг
	TraumaRoleTruthTeller   TraumaRole = "truth_teller"   // говорящий правду
	TraumaRoleShadow        TraumaRole = "shadow"         // тень семьи
	TraumaRoleHealer        TraumaRole = "healer"         // целитель/миротворец
)

// LifeEvent представляет значимое событие в жизни
type LifeEvent struct {
	Date        time.Time
	EventType   string // "birth", "death", "marriage", "divorce", "civil_union", "separation"
	Description string
	WithPerson  string // ID другого человека (для свадьбы и т.д.)
}

// FamilyMember представляет члена семьи с полными данными
type FamilyMember struct {
	ID              string
	Name            string
	Gender          Gender
	BirthDate       time.Time
	DeathDate       *time.Time      // может быть nil если жив
	Events          []LifeEvent     // все жизненные события
	Parents         []string        // ID родителей
	Children        []string        // ID детей
	Partners        []string        // ID партнеров (включая бывших)
	CurrentPartner  string          // ID текущего партнера
	IsDeceased      bool
	Generation      int             // вычисляется автоматически
}

// IdentifiedTraumaRoles содержит выявленные травматические роли
type IdentifiedTraumaRoles struct {
	PrimaryRole      TraumaRole
	SecondaryRoles   []TraumaRole
	Confidence       float64        // уверенность в определении (0-1)
	Evidence         []string       // доказательства/признаки
}

// FamilySystem представляет полную семейную систему
type FamilySystem struct {
	Members           map[string]*FamilyMember
	NarcissisticMembers []string     // ID членов с нарциссическими чертами (выявленные)
	EnablerMembers    []string       // ID созависимых (выявленные)
	TraumaRoles       map[string]*IdentifiedTraumaRoles // выявленные роли для каждого
	GenerationalTrauma float64       // общий уровень межпоколенческой травмы
	Patterns          []FamilyPattern // выявленные паттерны
}

// FamilyPattern представляет выявленный семейный паттерн
type FamilyPattern struct {
	PatternType     string
	Description     string
	Members         []string
	Severity        float64
	Generations     []int
	Recommendations []string
}

// NewFamilySystem создает новую семейную систему
func NewFamilySystem() *FamilySystem {
	return &FamilySystem{
		Members:           make(map[string]*FamilyMember),
		NarcissisticMembers: make([]string, 0),
		EnablerMembers:    make([]string, 0),
		TraumaRoles:       make(map[string]*IdentifiedTraumaRoles),
		Patterns:          make([]FamilyPattern, 0),
	}
}

// AddMember добавляет члена семьи с автоматическим определением поколения
func (fs *FamilySystem) AddMember(member *FamilyMember) {
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
	fs.Members[member.ID] = member
}

// ParseDate парсит дату в формате ДД.ММ.ГГГГ
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("02.01.2006", dateStr)
}

// CreateFamilyMember создает нового члена семьи с обязательными полями
func CreateFamilyMember(id, name string, gender Gender, birthDateStr string) (*FamilyMember, error) {
	birthDate, err := ParseDate(birthDateStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат даты рождения: %v", err)
	}
	return &FamilyMember{
		ID:         id,
		Name:       name,
		Gender:     gender,
		BirthDate:  birthDate,
		Events:     make([]LifeEvent, 0),
		Parents:    make([]string, 0),
		Children:   make([]string, 0),
		Partners:   make([]string, 0),
	}, nil
}

// AddLifeEvent добавляет жизненное событие
func (fm *FamilyMember) AddLifeEvent(eventType, description string, dateStr string, withPerson string) error {
	date, err := ParseDate(dateStr)
	if err != nil {
		return err
	}
	event := LifeEvent{
		Date:        date,
		EventType:   eventType,
		Description: description,
		WithPerson:  withPerson,
	}
	fm.Events = append(fm.Events, event)

	// Специальная обработка для свадьбы/развода
	switch eventType {
	case "marriage", "civil_union":
		if withPerson != "" {
			fm.Partners = append(fm.Partners, withPerson)
			fm.CurrentPartner = withPerson
		}
	case "divorce", "separation":
		fm.CurrentPartner = ""
	}
	return nil
}

// HypercubeCoords представляет координаты в гиперкубе
type HypercubeCoords struct {
	X int32
	Y int32
	Z int32
	W int32
}

// ParseDateToCoords преобразует дату в координаты гиперкуба
func ParseDateToCoords(date time.Time) HypercubeCoords {
	day := date.Day()
	month := int(date.Month())
	year := date.Year()

	// Преобразование в трехзначные числа с сохранением позиций
	x := toThreeDigitPreserve(day)
	y := toThreeDigitPreserve(month)
	z := toThreeDigitPreserve(year)
	// W - контрольная сумма
	w := sumDigits(x) + sumDigits(y) + sumDigits(z)
	return HypercubeCoords{
		X: int32(x),
		Y: int32(y),
		Z: int32(z),
		W: int32(w),
	}
}

// toThreeDigitPreserve дополняет число до трех знаков
func toThreeDigitPreserve(n int) int {
	if n < 10 {
		return n * 100
	}
	if n < 100 {
		return n * 10
	}
	return n % 1000
}

// sumDigits считает сумму цифр числа
func sumDigits(n int) int {
	sum := 0
	n = int(math.Abs(float64(n)))
	for n > 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

// extractDigits извлекает цифры из числа
func extractDigits(n int32) []int32 {
	n = int32(math.Abs(float64(n)))
	return []int32{
		(n / 100) % 10,
		(n / 10) % 10,
		n % 10,
	}
}

// PersonVectors содержит векторы для человека
type PersonVectors struct {
	X     []int32
	Y     []int32
	Z     []int32
	Full  [3][]int32
}

// CalculateVectors вычисляет векторы для координат
func CalculateVectors(coords HypercubeCoords) PersonVectors {
	digitsX := extractDigits(coords.X)
	digitsY := extractDigits(coords.Y)
	digitsZ := extractDigits(coords.Z)
	return PersonVectors{
		X: []int32{
			digitsX[0] - digitsX[1],
			digitsX[1] - digitsX[2],
			digitsX[2] - digitsX[0],
		},
		Y: []int32{
			digitsY[0] - digitsY[1],
			digitsY[1] - digitsY[2],
			digitsY[2] - digitsY[0],
		},
		Z: []int32{
			digitsZ[0] - digitsZ[1],
			digitsZ[1] - digitsZ[2],
			digitsZ[2] - digitsZ[0],
		},
		Full: [3][]int32{
			{
				digitsX[0] - digitsX[1],
				digitsX[1] - digitsX[2],
				digitsX[2] - digitsX[0],
			},
			{
				digitsY[0] - digitsY[1],
				digitsY[1] - digitsY[2],
				digitsY[2] - digitsY[0],
			},
			{
				digitsZ[0] - digitsZ[1],
				digitsZ[1] - digitsZ[2],
				digitsZ[2] - digitsZ[0],
			},
		},
	}
}

// MoveRoom перемещает комнату по векторам
func MoveRoom(coords HypercubeCoords, vectors PersonVectors, step int) HypercubeCoords {
	newCoords := coords
	vectorIdx := step % 3
	newCoords.X += vectors.Full[0][vectorIdx]
	newCoords.Y += vectors.Full[1][vectorIdx]
	newCoords.Z += vectors.Full[2][vectorIdx]
	newCoords.W = int32(sumDigits(int(newCoords.X)) +
		sumDigits(int(newCoords.Y)) +
		sumDigits(int(newCoords.Z)))

	// Модуль 26
	newCoords.X = ((newCoords.X % 26) + 26) % 26
	newCoords.Y = ((newCoords.Y % 26) + 26) % 26
	newCoords.Z = ((newCoords.Z % 26) + 26) % 26
	return newCoords
}

// IsTrapRoom проверяет, является ли комната ловушкой
func IsTrapRoom(coords HypercubeCoords) bool {
	digitsX := extractDigits(coords.X)
	digitsY := extractDigits(coords.Y)
	digitsZ := extractDigits(coords.Z)
	allEqual := func(d []int32) bool {
		return len(d) == 3 && d[0] == d[1] && d[1] == d[2]
	}
	return allEqual(digitsX) || allEqual(digitsY) || allEqual(digitsZ)
}

// CalculateBodyShame вычисляет уровень телесного стыда
func CalculateBodyShame(coords HypercubeCoords, traumaRole TraumaRole, ageAtEvent int) float64 {
	digits := extractDigits(coords.X + coords.Y + coords.Z)
	baseShame := float64(digits[0] * 10)

	// Коэффициенты для разных травматических ролей
	roleMultiplier := map[TraumaRole]float64{
		TraumaRoleScapegoat:       1.8,
		TraumaRoleGoldenChild:     1.3,
		TraumaRoleLostChild:       1.6,
		TraumaRoleInvisible:       1.5,
		TraumaRoleGlassChild:      1.4,
		TraumaRoleParentified:     1.7,
		TraumaRoleMascot:          1.2,
		TraumaRoleEmotionalSpouse: 1.9,
		TraumaRoleTruthTeller:     1.5,
		TraumaRoleShadow:          1.8,
		TraumaRoleHealer:          1.3,
	}[traumaRole]

	if roleMultiplier == 0 {
		roleMultiplier = 1.0
	}

	accumulated := baseShame * roleMultiplier * math.Exp(float64(ageAtEvent)/20.0)
	return math.Min(accumulated, 100.0)
}

// CompatibilityResult содержит результат анализа совместимости
type CompatibilityResult struct {
	Score           float64
	CommonAxes      int
	AmpDiff         float64
	SyncSteps       int
	Level           string
	Recommendations []string
}

// CalculateCompatibility вычисляет совместимость двух людей
func CalculateCompatibility(coords1, coords2 HypercubeCoords, vectors1, vectors2 PersonVectors) CompatibilityResult {
	// Общие оси
	commonAxes := 0
	if coords1.X == coords2.X {
		commonAxes++
	}
	if coords1.Y == coords2.Y {
		commonAxes++
	}
	if coords1.Z == coords2.Z {
		commonAxes++
	}
	if coords1.W == coords2.W {
		commonAxes++
	}

	// Амплитуда векторов
	amp1 := calculateVectorAmplitude(vectors1)
	amp2 := calculateVectorAmplitude(vectors2)
	ampDiff := math.Abs(amp1 - amp2)

	// Синхронные шаги
	syncSteps := 0
	pos1, pos2 := coords1, coords2
	for step := 0; step < 12; step++ {
		pos1 = MoveRoom(pos1, vectors1, step)
		pos2 = MoveRoom(pos2, vectors2, step)
		distance := euclideanDistance(pos1, pos2)
		if distance < 5.0 {
			syncSteps++
		}
	}

	// Общий счет совместимости
	score := float64(commonAxes)*0.25 +
		(1.0-math.Min(ampDiff/10.0, 1.0))*0.35 +
		float64(syncSteps)/12.0*0.4

	// Определение уровня
	level := "Низкая"
	recommendations := make([]string, 0)
	if score >= 0.7 {
		level = "Высокая"
		recommendations = append(recommendations,
			"Хороший потенциал для здоровых отношений при условии проработки индивидуальных травм")
	} else if score >= 0.4 {
		level = "Средняя"
		recommendations = append(recommendations,
			"Требуется осознанная работа над отношениями и четкие границы")
	} else {
		level = "Низкая"
		recommendations = append(recommendations,
			"Высокий риск деструктивных паттернов, рекомендуется осторожность")
	}

	return CompatibilityResult{
		Score:           score,
		CommonAxes:      commonAxes,
		AmpDiff:         ampDiff,
		SyncSteps:       syncSteps,
		Level:           level,
		Recommendations: recommendations,
	}
}

func calculateVectorAmplitude(vectors PersonVectors) float64 {
	sum := 0.0
	count := 0
	for _, vec := range vectors.Full {
		for _, v := range vec {
			sum += math.Abs(float64(v))
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

func euclideanDistance(a, b HypercubeCoords) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	dz := float64(a.Z - b.Z)
	dw := float64(a.W - b.W)
	return math.Sqrt(dx*dx + dy*dy + dz*dz + dw*dw)
}

// IdentifyTraumaRoles выявляет травматические роли на основе данных
func (fs *FamilySystem) IdentifyTraumaRoles(memberID string) *IdentifiedTraumaRoles {
	member := fs.Members[memberID]
	if member == nil {
		return nil
	}

	roles := &IdentifiedTraumaRoles{
		SecondaryRoles: make([]TraumaRole, 0),
		Evidence:       make([]string, 0),
	}

	coords := ParseDateToCoords(member.BirthDate)
	vectors := CalculateVectors(coords)

	// Анализ позиции в семье (порядок рождения)
	birthOrder := fs.getBirthOrder(memberID)
	siblings := fs.getSiblings(memberID)

	// 1. Проверка на козла отпущения (Scapegoat)
	if fs.isScapegoat(memberID, birthOrder, siblings) {
		roles.PrimaryRole = TraumaRoleScapegoat
		roles.Confidence = 0.8
		roles.Evidence = append(roles.Evidence,
			"Негативные проекции семьи",
			"Частая критика и обвинения",
			"Роль 'проблемного ребенка'")
	}

	// 2. Проверка на любимчика (Golden Child)
	if fs.isGoldenChild(memberID, birthOrder, siblings) {
		if roles.PrimaryRole == "" {
			roles.PrimaryRole = TraumaRoleGoldenChild
			roles.Confidence = 0.8
		} else {
			roles.SecondaryRoles = append(roles.SecondaryRoles, TraumaRoleGoldenChild)
		}
		roles.Evidence = append(roles.Evidence,
			"Идеализация родителями",
			"Высокие ожидания и давление",
			"Роль 'гордости семьи'")
	}

	// 3. Проверка на потерянного ребенка (Lost Child)
	if fs.isLostChild(memberID, vectors) {
		if roles.PrimaryRole == "" {
			roles.PrimaryRole = TraumaRoleLostChild
			roles.Confidence = 0.7
		} else {
			roles.SecondaryRoles = append(roles.SecondaryRoles, TraumaRoleLostChild)
		}
		roles.Evidence = append(roles.Evidence,
			"Эмоциональная отстраненность",
			"Склонность к уединению",
			"Минимальные потребности и запросы")
	}

	// 4. Проверка на невидимку (Invisible)
	if fs.isInvisible(memberID, siblings) {
		if roles.PrimaryRole == "" {
			roles.PrimaryRole = TraumaRoleInvisible
			roles.Confidence = 0.6
		} else {
			roles.SecondaryRoles = append(roles.SecondaryRoles, TraumaRoleInvisible)
		}
		roles.Evidence = append(roles.Evidence,
			"Заметен только в негативе",
			"Игнорирование потребностей",
			"Отсутствие личного пространства")
	}

	// 5. Проверка на стеклянного ребенка (Glass Child)
	if fs.isGlassChild(memberID, siblings) {
		if roles.PrimaryRole == "" {
			roles.PrimaryRole = TraumaRoleGlassChild
			roles.Confidence = 0.7
		} else {
			roles.SecondaryRoles = append(roles.SecondaryRoles, TraumaRoleGlassChild)
		}
		roles.Evidence = append(roles.Evidence,
			"На фоне проблемного сиблинга",
			"Сверх-самостоятельность",
			"Скрытые проблемы и потребности")
	}

	// 6. Проверка на родифицированного (Parentified)
	if fs.isParentified(memberID, birthOrder) {
		if roles.PrimaryRole == "" {
			roles.PrimaryRole = TraumaRoleParentified
			roles.Confidence = 0.8
		} else {
			roles.SecondaryRoles = append(roles.SecondaryRoles, TraumaRoleParentified)
		}
		roles.Evidence = append(roles.Evidence,
			"Забота о родителях/сиблингах",
			"Отсутствие детства",
			"Чувство ответственности за всех")
	}

	// 7. Проверка на шута (Mascot)
	if fs.isMascot(memberID, vectors) {
		if roles.PrimaryRole == "" {
			roles.PrimaryRole = TraumaRoleMascot
			roles.Confidence = 0.6
		} else {
			roles.SecondaryRoles = append(roles.SecondaryRoles, TraumaRoleMascot)
		}
		roles.Evidence = append(roles.Evidence,
			"Снятие напряжения юмором",
			"Минимизация проблем",
			"Роль 'клоуна' в семье")
	}

	// 8. Проверка на эмоционального супруга (Emotional Spouse)
	if fs.isEmotionalSpouse(memberID) {
		if roles.PrimaryRole == "" {
			roles.PrimaryRole = TraumaRoleEmotionalSpouse
			roles.Confidence = 0.9
		} else {
			roles.SecondaryRoles = append(roles.SecondaryRoles, TraumaRoleEmotionalSpouse)
		}
		roles.Evidence = append(roles.Evidence,
			"Замещение супружеской роли",
			"Эмоциональная близость с родителем",
			"Ревность со стороны другого родителя")
	}

	// 9. Феномен "девочка с яйцами" (Girl with Balls)
	if member.Gender == GenderFemale &&
		(roles.PrimaryRole == TraumaRoleLostChild || roles.PrimaryRole == TraumaRoleScapegoat) &&
		calculateVectorAmplitude(vectors) > 5.0 {
		roles.Evidence = append(roles.Evidence,
			"Феномен 'девочка с яйцами' - маскулинные стратегии выживания",
			"Компенсаторная гипер-независимость",
			"Отказ от феминности как защита")
	}

	// Если роль не определена, ставим Shadow (тень)
	if roles.PrimaryRole == "" {
		roles.PrimaryRole = TraumaRoleShadow
		roles.Confidence = 0.5
		roles.Evidence = append(roles.Evidence, "Невыраженная травма, требуется дополнительный анализ")
	}

	fs.TraumaRoles[memberID] = roles
	return roles
}

// Вспомогательные методы для выявления ролей
func (fs *FamilySystem) getBirthOrder(memberID string) int {
	member := fs.Members[memberID]
	if member == nil {
		return 0
	}
	// Находим всех сиблингов
	siblings := fs.getSiblings(memberID)
	// Сортируем по дате рождения
	type siblingInfo struct {
		id   string
		date time.Time
	}
	sibList := make([]siblingInfo, 0)
	for _, sibID := range siblings {
		if sib, ok := fs.Members[sibID]; ok {
			sibList = append(sibList, siblingInfo{id: sibID, date: sib.BirthDate})
		}
	}
	sort.Slice(sibList, func(i, j int) bool {
		return sibList[i].date.Before(sibList[j].date)
	})
	for i, sib := range sibList {
		if sib.id == memberID {
			return i + 1
		}
	}
	return 0
}

func (fs *FamilySystem) getSiblings(memberID string) []string {
	member := fs.Members[memberID]
	if member == nil {
		return nil
	}
	siblings := make([]string, 0)
	// Ищем через общих родителей
	for _, parentID := range member.Parents {
		if parent, ok := fs.Members[parentID]; ok {
			for _, childID := range parent.Children {
				if childID != memberID && !contains(siblings, childID) {
					siblings = append(siblings, childID)
				}
			}
		}
	}
	return siblings
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (fs *FamilySystem) isScapegoat(memberID string, birthOrder int, siblings []string) bool {
	member := fs.Members[memberID]
	if member == nil {
		return false
	}
	// Козел отпущения часто:
	// - Не первый и не последний ребенок
	// - Часто тот же пол, что и нарциссический родитель
	// - Имеет больше негативных событий в истории
	negativeEvents := 0
	for _, event := range member.Events {
		if event.EventType == "conflict" || event.EventType == "criticism" {
			negativeEvents++
		}
	}
	// Проверка на критику в свидетельствах других членов семьи
	criticizedByOthers := 0
	for _, otherID := range siblings {
		if other, ok := fs.Members[otherID]; ok {
			for _, event := range other.Events {
				if event.EventType == "criticism" && event.WithPerson == memberID {
					criticizedByOthers++
				}
			}
		}
	}
	return negativeEvents > 2 || criticizedByOthers > 1
}

func (fs *FamilySystem) isGoldenChild(memberID string, birthOrder int, siblings []string) bool {
	member := fs.Members[memberID]
	if member == nil {
		return false
	}
	// Любимчик часто:
	// - Первый или последний ребенок
	// - Часто противоположного пола с нарциссическим родителем
	// - Имеет больше позитивных проекций
	positiveEvents := 0
	for _, event := range member.Events {
		if event.EventType == "praise" || event.EventType == "achievement" {
			positiveEvents++
		}
	}
	// Проверка на идеализацию
	idealizedByParents := 0
	for _, parentID := range member.Parents {
		if parent, ok := fs.Members[parentID]; ok {
			for _, event := range parent.Events {
				if event.EventType == "praise" && event.WithPerson == memberID {
					idealizedByParents++
				}
			}
		}
	}
	return positiveEvents > 3 || idealizedByParents > 2
}

func (fs *FamilySystem) isLostChild(memberID string, vectors PersonVectors) bool {
	member := fs.Members[memberID]
	if member == nil {
		return false
	}
	// Потерянный ребенок:
	// - Низкая амплитуда векторов (эмоциональная заморозка)
	// - Мало событий в истории
	// - Отсутствие значимых отношений
	amp := calculateVectorAmplitude(vectors)
	eventCount := len(member.Events)
	return amp < 3.0 && eventCount < 3 && len(member.Partners) < 2
}

func (fs *FamilySystem) isInvisible(memberID string, siblings []string) bool {
	member := fs.Members[memberID]
	if member == nil {
		return false
	}
	// Невидимка:
	// - Нет значимых событий
	// - Не упоминается в событиях других
	// - Средний ребенок в многодетной семье
	mentionedByOthers := 0
	for _, sibID := range siblings {
		if sib, ok := fs.Members[sibID]; ok {
			for _, event := range sib.Events {
				if event.WithPerson == memberID {
					mentionedByOthers++
				}
			}
		}
	}
	return mentionedByOthers == 0 && len(member.Events) < 2
}

func (fs *FamilySystem) isGlassChild(memberID string, siblings []string) bool {
	// Стеклянный ребенок растет рядом с проблемным сиблингом
	// и становится "невидимым" на его фоне
	for _, sibID := range siblings {
		if sib, ok := fs.Members[sibID]; ok {
			// Проверяем, есть ли сиблинг с серьезными проблемами
			problematicEvents := 0
			for _, event := range sib.Events {
				if event.EventType == "hospitalization" ||
					event.EventType == "addiction" ||
					event.EventType == "crisis" {
					problematicEvents++
				}
			}
			if problematicEvents > 2 {
				return true
			}
		}
	}
	return false
}

func (fs *FamilySystem) isParentified(memberID string, birthOrder int) bool {
	member := fs.Members[memberID]
	if member == nil {
		return false
	}
	// Родифицированный ребенок:
	// - Часто старший
	// - Заботится о родителях/сиблингах
	// - Рано начал работать/брать ответственность
	careEvents := 0
	for _, event := range member.Events {
		if event.EventType == "care" || event.EventType == "responsibility" {
			careEvents++
		}
	}
	return birthOrder == 1 && careEvents > 2
}

func (fs *FamilySystem) isMascot(memberID string, vectors PersonVectors) bool {
	member := fs.Members[memberID]
	if member == nil {
		return false
	}
	// Шут:
	// - Высокая амплитуда векторов (гиперактивность)
	// - Много событий, но поверхностных
	// - Использует юмор в stressful ситуациях
	amp := calculateVectorAmplitude(vectors)
	humorEvents := 0
	for _, event := range member.Events {
		if event.EventType == "humor" || event.EventType == "entertainment" {
			humorEvents++
		}
	}
	return amp > 6.0 && humorEvents > 2
}

func (fs *FamilySystem) isEmotionalSpouse(memberID string) bool {
	member := fs.Members[memberID]
	if member == nil {
		return false
	}
	// Эмоциональный супруг:
	// - Особо близкие отношения с родителем противоположного пола
	// - Отсутствие здоровых романтических отношений
	// - Ревность со стороны другого родителя
	closeWithParent := false
	for _, parentID := range member.Parents {
		if parent, ok := fs.Members[parentID]; ok {
			// Проверяем интенсивность взаимодействия
			interactionCount := 0
			for _, event := range parent.Events {
				if event.WithPerson == memberID &&
					(event.EventType == "intimate" || event.EventType == "confidant") {
					interactionCount++
				}
			}
			if interactionCount > 3 {
				closeWithParent = true
			}
		}
	}
	return closeWithParent && len(member.Partners) < 2
}

// IdentifyNarcissisticMembers выявляет членов с нарциссическими чертами
func (fs *FamilySystem) IdentifyNarcissisticMembers() {
	for id, member := range fs.Members {
		score := 0.0
		evidence := make([]string, 0)
		// Анализ отношений с детьми
		if len(member.Children) > 0 {
			// Проверка на наличие любимчика и козла отпущения среди детей
			hasGoldenChild := false
			hasScapegoat := false
			for _, childID := range member.Children {
				if role, ok := fs.TraumaRoles[childID]; ok {
					if role.PrimaryRole == TraumaRoleGoldenChild {
						hasGoldenChild = true
					}
					if role.PrimaryRole == TraumaRoleScapegoat {
						hasScapegoat = true
					}
				}
			}
			if hasGoldenChild && hasScapegoat {
				score += 0.6
				evidence = append(evidence, "Создание раскола между детьми (любимчик vs козел отпущения)")
			}
		}
		// Анализ партнерских отношений
		if len(member.Partners) > 1 {
			// Частая смена партнеров может указывать на нарциссические черты
			score += 0.2
			evidence = append(evidence, "Нестабильные партнерские отношения")
		}
		// Анализ векторов (высокая амплитуда может указывать на нестабильность)
		coords := ParseDateToCoords(member.BirthDate)
		vectors := CalculateVectors(coords)
		amp := calculateVectorAmplitude(vectors)
		if amp > 7.0 {
			score += 0.3
			evidence = append(evidence, "Высокая эмоциональная нестабильность (векторный анализ)")
		}
		// Если есть события грандиозности или обесценивания
		for _, event := range member.Events {
			if event.EventType == "grandiosity" {
				score += 0.4
				evidence = append(evidence, "Проявления грандиозности")
			}
			if event.EventType == "devaluation" {
				score += 0.3
				evidence = append(evidence, "Склонность к обесцениванию")
			}
		}
		if score > 0.7 {
			fs.NarcissisticMembers = append(fs.NarcissisticMembers, id)
		}
	}
}

// IdentifyEnablerMembers выявляет созависимых членов
func (fs *FamilySystem) IdentifyEnablerMembers() {
	for id, member := range fs.Members {
		score := 0.0
		evidence := make([]string, 0)
		// Защита/оправдание нарциссического партнера
		for _, partnerID := range member.Partners {
			if contains(fs.NarcissisticMembers, partnerID) {
				// Проверяем события, где member защищает нарцисса
				for _, event := range member.Events {
					if event.EventType == "defense" && event.WithPerson == partnerID {
						score += 0.4
						evidence = append(evidence, "Защита нарциссического партнера")
					}
				}
			}
		}
		// Принятие на себя вины за проблемы семьи
		for _, event := range member.Events {
			if event.EventType == "self_blame" {
				score += 0.3
				evidence = append(evidence, "Склонность к самообвинению")
			}
		}
		// Низкая амплитуда векторов (подавление себя)
		coords := ParseDateToCoords(member.BirthDate)
		vectors := CalculateVectors(coords)
		amp := calculateVectorAmplitude(vectors)
		if amp < 3.0 {
			score += 0.3
			evidence = append(evidence, "Эмоциональное подавление")
		}
		if score > 0.5 {
			fs.EnablerMembers = append(fs.EnablerMembers, id)
		}
	}
}

// FindFamilyPatterns выявляет семейные паттерны
func (fs *FamilySystem) FindFamilyPatterns() []FamilyPattern {
	patterns := make([]FamilyPattern, 0)
	// 1. Паттерн "мать-дочь" (передача травмы по женской линии)
	motherDaughterPattern := fs.findMotherDaughterPattern()
	if motherDaughterPattern != nil {
		patterns = append(patterns, *motherDaughterPattern)
	}
	// 2. Паттерн "мать-сын" (эмоциональный инцест)
	motherSonPattern := fs.findMotherSonPattern()
	if motherSonPattern != nil {
		patterns = append(patterns, *motherSonPattern)
	}
	// 3. Паттерн "отец-дочь"
	fatherDaughterPattern := fs.findFatherDaughterPattern()
	if fatherDaughterPattern != nil {
		patterns = append(patterns, *fatherDaughterPattern)
	}
	// 4. Паттерн "отец-сын"
	fatherSonPattern := fs.findFatherSonPattern()
	if fatherSonPattern != nil {
		patterns = append(patterns, *fatherSonPattern)
	}
	// 5. Паттерн "бабушка-мать-внучка" (три поколения)
	grandmotherPattern := fs.findGrandmotherPattern()
	if grandmotherPattern != nil {
		patterns = append(patterns, *grandmotherPattern)
	}
	// 6. Паттерн "дедушка-отец-внук"
	grandfatherPattern := fs.findGrandfatherPattern()
	if grandfatherPattern != nil {
		patterns = append(patterns, *grandfatherPattern)
	}
	// 7. Паттерн "бабушка-внук" (через поколение)
	grandmotherGrandsonPattern := fs.findGrandmotherGrandsonPattern()
	if grandmotherGrandsonPattern != nil {
		patterns = append(patterns, *grandmotherGrandsonPattern)
	}
	// 8. Паттерн "дедушка-внучка" (через поколение)
	grandfatherGranddaughterPattern := fs.findGrandfatherGranddaughterPattern()
	if grandfatherGranddaughterPattern != nil {
		patterns = append(patterns, *grandfatherGranddaughterPattern)
	}
	// 9. Паттерн "сестра-сестра" (сиблинговая динамика)
	sisterSisterPattern := fs.findSisterSisterPattern()
	if sisterSisterPattern != nil {
		patterns = append(patterns, *sisterSisterPattern)
	}
	// 10. Паттерн "брат-брат"
	brotherBrotherPattern := fs.findBrotherBrotherPattern()
	if brotherBrotherPattern != nil {
		patterns = append(patterns, *brotherBrotherPattern)
	}
	// 11. Паттерн "сестра-брат"
	sisterBrotherPattern := fs.findSisterBrotherPattern()
	if sisterBrotherPattern != nil {
		patterns = append(patterns, *sisterBrotherPattern)
	}
	// 12. Паттерн "повторение браков" (серия разводов)
	divorcePattern := fs.findDivorcePattern()
	if divorcePattern != nil {
		patterns = append(patterns, *divorcePattern)
	}
	fs.Patterns = patterns
	return patterns
}

func (fs *FamilySystem) findMotherDaughterPattern() *FamilyPattern {
	// Ищем матерей и дочерей с похожими травматическими ролями
	for id, member := range fs.Members {
		if member.Gender != GenderFemale || member.Role != RoleMother {
			continue
		}
		for _, childID := range member.Children {
			child := fs.Members[childID]
			if child == nil || child.Gender != GenderFemale {
				continue
			}
			motherRole := fs.TraumaRoles[id]
			childRole := fs.TraumaRoles[childID]
			if motherRole != nil && childRole != nil {
				// Проверяем передачу ролей
				if motherRole.PrimaryRole == TraumaRoleScapegoat &&
					childRole.PrimaryRole == TraumaRoleScapegoat {
					return &FamilyPattern{
						PatternType: "mother-daughter-scapegoat",
						Description: "Передача роли козла отпущения по женской линии",
						Members:     []string{id, childID},
						Severity:    0.8,
						Generations: []int{member.Generation, child.Generation},
						Recommendations: []string{
							"Терапия женской линии",
							"Работа с родовыми сценариями",
							"Сепарация от материнских проекций",
						},
					}
				}
			}
		}
	}
	return nil
}

func (fs *FamilySystem) findMotherSonPattern() *FamilyPattern {
	for id, member := range fs.Members {
		if member.Gender != GenderFemale || member.Role != RoleMother {
			continue
		}
		// Проверяем, есть ли сын, который является эмоциональным супругом
		for _, childID := range member.Children {
			child := fs.Members[childID]
			if child == nil || child.Gender != GenderMale {
				continue
			}
			childRole := fs.TraumaRoles[childID]
			if childRole != nil && childRole.PrimaryRole == TraumaRoleEmotionalSpouse {
				return &FamilyPattern{
					PatternType: "mother-son-emotional-incest",
					Description: "Эмоциональный инцест: сын замещает роль мужа",
					Members:     []string{id, childID},
					Severity:    0.9,
					Generations: []int{member.Generation, child.Generation},
					Recommendations: []string{
						"Терапия для разрыва симбиотической связи",
						"Работа с мужской идентичностью",
						"Отделение от материнских ожиданий",
					},
				}
			}
		}
	}
	return nil
}

func (fs *FamilySystem) findFatherDaughterPattern() *FamilyPattern {
	for id, member := range fs.Members {
		if member.Gender != GenderMale || member.Role != RoleFather {
			continue
		}
		for _, childID := range member.Children {
			child := fs.Members[childID]
			if child == nil || child.Gender != GenderFemale {
				continue
			}
			childRole := fs.TraumaRoles[childID]
			if childRole != nil && childRole.PrimaryRole == TraumaRoleGoldenChild {
				return &FamilyPattern{
					PatternType: "father-daughter-golden",
					Description: "Дочь-любимица отца",
					Members:     []string{id, childID},
					Severity:    0.7,
					Generations: []int{member.Generation, child.Generation},
					Recommendations: []string{
						"Работа с комплексом Электры",
						"Формирование здоровых отношений с мужчинами",
						"Отделение от отцовских проекций",
					},
				}
			}
		}
	}
	return nil
}

func (fs *FamilySystem) findFatherSonPattern() *FamilyPattern {
	for id, member := range fs.Members {
		if member.Gender != GenderMale || member.Role != RoleFather {
			continue
		}
		for _, childID := range member.Children {
			child := fs.Members[childID]
			if child == nil || child.Gender != GenderMale {
				continue
			}
			fatherRole := fs.TraumaRoles[id]
			childRole := fs.TraumaRoles[childID]
			if fatherRole != nil && childRole != nil {
				if fatherRole.PrimaryRole == TraumaRoleScapegoat &&
					childRole.PrimaryRole == TraumaRoleScapegoat {
					return &FamilyPattern{
						PatternType: "father-son-scapegoat",
						Description: "Передача роли козла отпущения по мужской линии",
						Members:     []string{id, childID},
						Severity:    0.8,
						Generations: []int{member.Generation, child.Generation},
						Recommendations: []string{
							"Терапия мужской линии",
							"Работа с родовыми проклятиями",
							"Восстановление мужской идентичности",
						},
					}
				}
			}
		}
	}
	return nil
}

func (fs *FamilySystem) findGrandmotherPattern() *FamilyPattern {
	// Паттерн бабушка-мать-внучка (три поколения)
	for gmID, gm := range fs.Members {
		if gm.Gender != GenderFemale || gm.Role != RoleGrandmother {
			continue
		}
		for _, childID := range gm.Children {
			mother := fs.Members[childID]
			if mother == nil || mother.Gender != GenderFemale || mother.Role != RoleMother {
				continue
			}
			for _, grandchildID := range mother.Children {
				grandchild := fs.Members[grandchildID]
				if grandchild == nil || grandchild.Gender != GenderFemale {
					continue
				}
				// Проверяем передачу травмы через поколения
				gmRole := fs.TraumaRoles[gmID]
				motherRole := fs.TraumaRoles[childID]
				gcRole := fs.TraumaRoles[grandchildID]
				if gmRole != nil && motherRole != nil && gcRole != nil {
					// Если роли совпадают или усиливаются
					if gmRole.PrimaryRole == motherRole.PrimaryRole ||
						motherRole.PrimaryRole == gcRole.PrimaryRole {
						return &FamilyPattern{
							PatternType: "grandmother-mother-granddaughter",
							Description: "Передача травмы через три поколения женщин",
							Members:     []string{gmID, childID, grandchildID},
							Severity:    0.9,
							Generations: []int{1, 2, 3},
							Recommendations: []string{
								"Глубокая родовая терапия",
								"Работа с женской линией",
								"Разрыв родового проклятия",
								"Индивидуальная терапия для каждой",
							},
						}
					}
				}
			}
		}
	}
	return nil
}

func (fs *FamilySystem) findGrandfatherPattern() *FamilyPattern {
	// Паттерн дедушка-отец-внук
	for gfID, gf := range fs.Members {
		if gf.Gender != GenderMale || gf.Role != RoleGrandfather {
			continue
		}
		for _, childID := range gf.Children {
			father := fs.Members[childID]
			if father == nil || father.Gender != GenderMale || father.Role != RoleFather {
				continue
			}
			for _, grandchildID := range father.Children {
				grandchild := fs.Members[grandchildID]
				if grandchild == nil || grandchild.Gender != GenderMale {
					continue
				}
				gfRole := fs.TraumaRoles[gfID]
				fatherRole := fs.TraumaRoles[childID]
				gcRole := fs.TraumaRoles[grandchildID]
				if gfRole != nil && fatherRole != nil && gcRole != nil {
					if gfRole.PrimaryRole == fatherRole.PrimaryRole ||
						fatherRole.PrimaryRole == gcRole.PrimaryRole {
						return &FamilyPattern{
							PatternType: "grandfather-father-grandson",
							Description: "Передача травмы через три поколения мужчин",
							Members:     []string{gfID, childID, grandchildID},
							Severity:    0.9,
							Generations: []int{1, 2, 3},
							Recommendations: []string{
								"Мужская родовая терапия",
								"Работа с родовыми сценариями",
								"Восстановление мужской линии",
							},
						}
					}
				}
			}
		}
	}
	return nil
}

func (fs *FamilySystem) findGrandmotherGrandsonPattern() *FamilyPattern {
	// Паттерн бабушка-внук (через поколение)
	for gmID, gm := range fs.Members {
		if gm.Gender != GenderFemale || gm.Role != RoleGrandmother {
			continue
		}
		// Ищем внуков напрямую (минуя родителей)
		for _, childID := range gm.Children {
			parent := fs.Members[childID]
			if parent == nil {
				continue
			}
			for _, grandchildID := range parent.Children {
				grandchild := fs.Members[grandchildID]
				if grandchild == nil || grandchild.Gender != GenderMale {
					continue
				}
				// Проверяем особую связь бабушка-внук
				// Часто внук становится "заместителем" нелюбимого сына
				parentRole := fs.TraumaRoles[childID]
				if parentRole != nil && parentRole.PrimaryRole == TraumaRoleScapegoat {
					return &FamilyPattern{
						PatternType: "grandmother-grandson-compensation",
						Description: "Бабушка компенсирует через внука отвергнутого сына",
						Members:     []string{gmID, grandchildID},
						Severity:    0.7,
						Generations: []int{1, 3},
						Recommendations: []string{
							"Осознание переноса",
							"Работа с сыном (пропущенным поколением)",
							"Терапия для внука от слияния",
						},
					}
				}
			}
		}
	}
	return nil
}

func (fs *FamilySystem) findGrandfatherGranddaughterPattern() *FamilyPattern {
	// Паттерн дедушка-внучка
	for gfID, gf := range fs.Members {
		if gf.Gender != GenderMale || gf.Role != RoleGrandfather {
			continue
		}
		for _, childID := range gf.Children {
			parent := fs.Members[childID]
			if parent == nil {
				continue
			}
			for _, grandchildID := range parent.Children {
				grandchild := fs.Members[grandchildID]
				if grandchild == nil || grandchild.Gender != GenderFemale {
					continue
				}
				// Особая связь дедушка-внучка
				return &FamilyPattern{
					PatternType: "grandfather-granddaughter-bond",
					Description: "Особая связь дедушки и внучки (часто компенсаторная)",
					Members:     []string{gfID, grandchildID},
					Severity:    0.6,
					Generations: []int{1, 3},
					Recommendations: []string{
						"Проверка на эмоциональный инцест",
						"Работа с женской идентичностью",
						"Отделение от дедушкиных проекций",
					},
				}
			}
		}
	}
	return nil
}

func (fs *FamilySystem) findSisterSisterPattern() *FamilyPattern {
	// Ищем сестер с комплементарными ролями
	sisters := make([]string, 0)
	for id, member := range fs.Members {
		if member.Gender == GenderFemale {
			sisters = append(sisters, id)
		}
	}
	// Группируем по родителям
	sisterGroups := make(map[string][]string)
	for _, s1 := range sisters {
		m1 := fs.Members[s1]
		for _, s2 := range sisters {
			if s1 >= s2 {
				continue
			}
			m2 := fs.Members[s2]
			// Проверяем, есть ли общие родители
			commonParents := false
			for _, p1 := range m1.Parents {
				for _, p2 := range m2.Parents {
					if p1 == p2 {
						commonParents = true
						break
					}
				}
			}
			if commonParents {
				key := fmt.Sprintf("%v", m1.Parents)
				if _, ok := sisterGroups[key]; !ok {
					sisterGroups[key] = make([]string, 0)
				}
				if !contains(sisterGroups[key], s1) {
					sisterGroups[key] = append(sisterGroups[key], s1)
				}
				if !contains(sisterGroups[key], s2) {
					sisterGroups[key] = append(sisterGroups[key], s2)
				}
			}
		}
	}
	// Анализируем каждую группу сестер
	for _, group := range sisterGroups {
		if len(group) >= 2 {
			roles := make([]TraumaRole, 0)
			for _, sid := range group {
				if role, ok := fs.TraumaRoles[sid]; ok {
					roles = append(roles, role.PrimaryRole)
				}
			}
			// Проверяем классическую динамику: любимчик + козел отпущения
			hasGolden := false
			hasScapegoat := false
			for _, r := range roles {
				if r == TraumaRoleGoldenChild {
					hasGolden = true
				}
				if r == TraumaRoleScapegoat {
					hasScapegoat = true
				}
			}
			if hasGolden && hasScapegoat {
				return &FamilyPattern{
					PatternType: "sister-sister-golden-scapegoat",
					Description: "Классическая динамика сестер: любимица и козел отпущения",
					Members:     group,
					Severity:    0.8,
					Generations: []int{fs.Members[group[0]].Generation},
					Recommendations: []string{
						"Терапия сиблинговых отношений",
						"Работа с завистью и соперничеством",
						"Восстановление сестринской связи",
					},
				}
			}
		}
	}
	return nil
}

func (fs *FamilySystem) findBrotherBrotherPattern() *FamilyPattern {
	// Аналогично для братьев
	brothers := make([]string, 0)
	for id, member := range fs.Members {
		if member.Gender == GenderMale {
			brothers = append(brothers, id)
		}
	}
	brotherGroups := make(map[string][]string)
	for _, b1 := range brothers {
		m1 := fs.Members[b1]
		for _, b2 := range brothers {
			if b1 >= b2 {
				continue
			}
			m2 := fs.Members[b2]
			commonParents := false
			for _, p1 := range m1.Parents {
				for _, p2 := range m2.Parents {
					if p1 == p2 {
						commonParents = true
						break
					}
				}
			}
			if commonParents {
				key := fmt.Sprintf("%v", m1.Parents)
				if _, ok := brotherGroups[key]; !ok {
					brotherGroups[key] = make([]string, 0)
				}
				if !contains(brotherGroups[key], b1) {
					brotherGroups[key] = append(brotherGroups[key], b1)
				}
				if !contains(brotherGroups[key], b2) {
					brotherGroups[key] = append(brotherGroups[key], b2)
				}
			}
		}
	}
	for _, group := range brotherGroups {
		if len(group) >= 2 {
			roles := make([]TraumaRole, 0)
			for _, bid := range group {
				if role, ok := fs.TraumaRoles[bid]; ok {
					roles = append(roles, role.PrimaryRole)
				}
			}
			hasGolden := false
			hasScapegoat := false
			hasLost := false
			for _, r := range roles {
				if r == TraumaRoleGoldenChild {
					hasGolden = true
				}
				if r == TraumaRoleScapegoat {
					hasScapegoat = true
				}
				if r == TraumaRoleLostChild {
					hasLost = true
				}
			}
			if hasGolden && hasScapegoat && hasLost {
				return &FamilyPattern{
					PatternType: "brother-brother-triangle",
					Description: "Треугольник братьев: любимчик, козел отпущения и потерянный",
					Members:     group,
					Severity:    0.9,
					Generations: []int{fs.Members[group[0]].Generation},
					Recommendations: []string{
						"Мужская групповая терапия",
						"Работа с конкуренцией",
						"Восстановление братской связи",
					},
				}
			}
		}
	}
	return nil
}

func (fs *FamilySystem) findSisterBrotherPattern() *FamilyPattern {
	// Динамика между сестрами и братьями
	pairs := make([][2]string, 0)
	for sid, sister := range fs.Members {
		if sister.Gender != GenderFemale {
			continue
		}
		for bid, brother := range fs.Members {
			if brother.Gender != GenderMale || sid == bid {
				continue
			}
			// Проверяем, есть ли общие родители
			commonParents := false
			for _, p1 := range sister.Parents {
				for _, p2 := range brother.Parents {
					if p1 == p2 {
						commonParents = true
						break
					}
				}
			}
			if commonParents {
				pairs = append(pairs, [2]string{sid, bid})
			}
		}
	}
	for _, pair := range pairs {
		sisterRole := fs.TraumaRoles[pair[0]]
		brotherRole := fs.TraumaRoles[pair[1]]
		if sisterRole != nil && brotherRole != nil {
			// Специфические паттерны
			if sisterRole.PrimaryRole == TraumaRoleParentified &&
				brotherRole.PrimaryRole == TraumaRoleLostChild {
				return &FamilyPattern{
					PatternType: "sister-brother-parentified-lost",
					Description: "Сестра-родитель и брат-потерянный ребенок",
					Members:     []string{pair[0], pair[1]},
					Severity:    0.7,
					Generations: []int{fs.Members[pair[0]].Generation},
					Recommendations: []string{
						"Освобождение сестры от родительской роли",
						"Активация брата",
						"Баланс ответственности",
					},
				}
			}
			// Феномен "девочка с яйцами" и ее брат
			if sisterRole.PrimaryRole == TraumaRoleScapegoat &&
				calculateVectorAmplitude(CalculateVectors(ParseDateToCoords(fs.Members[pair[0]].BirthDate))) > 5.0 &&
				brotherRole.PrimaryRole == TraumaRoleLostChild {
				return &FamilyPattern{
					PatternType: "sister-brother-girl-with-balls",
					Description: "Сестра с маскулинными стратегиями и потерянный брат",
					Members:     []string{pair[0], pair[1]},
					Severity:    0.8,
					Generations: []int{fs.Members[pair[0]].Generation},
					Recommendations: []string{
						"Интеграция маскулинной и феминной энергии у сестры",
						"Активация мужского начала у брата",
						"Баланс ролей",
					},
				}
			}
		}
	}
	return nil
}

func (fs *FamilySystem) findDivorcePattern() *FamilyPattern {
	// Паттерн повторения разводов в поколениях
	divorceGenerations := make(map[int][]string)
	for id, member := range fs.Members {
		divorceCount := 0
		for _, event := range member.Events {
			if event.EventType == "divorce" {
				divorceCount++
			}
		}
		if divorceCount >= 1 {
			divorceGenerations[member.Generation] = append(divorceGenerations[member.Generation], id)
		}
	}
	// Проверяем наличие разводов в нескольких поколениях
	if len(divorceGenerations) >= 2 {
		members := make([]string, 0)
		gens := make([]int, 0)
		for gen, ids := range divorceGenerations {
			gens = append(gens, gen)
			members = append(members, ids...)
		}
		sort.Ints(gens)
		return &FamilyPattern{
			PatternType: "generational-divorce",
			Description: fmt.Sprintf("Повторение разводов в поколениях %v", gens),
			Members:     members,
			Severity:    0.7,
			Generations: gens,
			Recommendations: []string{
				"Анализ родовых сценариев отношений",
				"Терапия привязанности",
				"Работа с выбором партнера",
				"Осознание повторяющихся паттернов",
			},
		}
	}
	return nil
}

// CalculateGenerationalTrauma вычисляет общий уровень межпоколенческой травмы
func (fs *FamilySystem) CalculateGenerationalTrauma() float64 {
	if len(fs.Members) == 0 {
		return 0
	}
	totalTrauma := 0.0
	count := 0
	for id, member := range fs.Members {
		coords := ParseDateToCoords(member.BirthDate)
		role := fs.TraumaRoles[id]
		if role != nil {
			ageAtEvent := time.Now().Year() - member.BirthDate.Year()
			trauma := CalculateBodyShame(coords, role.PrimaryRole, ageAtEvent)
			totalTrauma += trauma
			count++
		}
	}
	if count > 0 {
		fs.GenerationalTrauma = totalTrauma / float64(count)
	}
	return fs.GenerationalTrauma
}

// AnalyzeRelationship анализирует отношения между двумя людьми
func (fs *FamilySystem) AnalyzeRelationship(id1, id2 string, relType RelationshipType) *RelationshipAnalysis {
	m1 := fs.Members[id1]
	m2 := fs.Members[id2]
	if m1 == nil || m2 == nil {
		return nil
	}
	coords1 := ParseDateToCoords(m1.BirthDate)
	coords2 := ParseDateToCoords(m2.BirthDate)
	vectors1 := CalculateVectors(coords1)
	vectors2 := CalculateVectors(coords2)
	compatibility := CalculateCompatibility(coords1, coords2, vectors1, vectors2)
	role1 := fs.TraumaRoles[id1]
	role2 := fs.TraumaRoles[id2]
	analysis := &RelationshipAnalysis{
		Person1:        m1,
		Person2:        m2,
		Type:           relType,
		Compatibility:  compatibility,
		TraumaDynamics: make([]string, 0),
		Warnings:       make([]string, 0),
		Recommendations: make([]string, 0),
	}
	// Анализ динамики в зависимости от типа отношений
	switch relType {
	case RelationshipMotherDaughter:
		fs.analyzeMotherDaughter(analysis, role1, role2)
	case RelationshipMotherSon:
		fs.analyzeMotherSon(analysis, role1, role2)
	case RelationshipFatherDaughter:
		fs.analyzeFatherDaughter(analysis, role1, role2)
	case RelationshipFatherSon:
		fs.analyzeFatherSon(analysis, role1, role2)
	case RelationshipSisterSister:
		fs.analyzeSisterSister(analysis, role1, role2)
	case RelationshipBrotherBrother:
		fs.analyzeBrotherBrother(analysis, role1, role2)
	case RelationshipSisterBrother:
		fs.analyzeSisterBrother(analysis, role1, role2)
	case RelationshipGrandmotherGranddaughter:
		fs.analyzeGrandmotherGranddaughter(analysis, role1, role2)
	case RelationshipGrandmotherGrandson:
		fs.analyzeGrandmotherGrandson(analysis, role1, role2)
	case RelationshipGrandfatherGranddaughter:
		fs.analyzeGrandfatherGranddaughter(analysis, role1, role2)
	case RelationshipGrandfatherGrandson:
		fs.analyzeGrandfatherGrandson(analysis, role1, role2)
	case RelationshipHusbandWife, RelationshipCivilPartnership:
		fs.analyzePartners(analysis, role1, role2)
	case RelationshipExSpouses:
		fs.analyzeExPartners(analysis, role1, role2)
	}
	return analysis
}

type RelationshipAnalysis struct {
	Person1        *FamilyMember
	Person2        *FamilyMember
	Type           RelationshipType
	Compatibility  CompatibilityResult
	TraumaDynamics []string
	Warnings       []string
	Recommendations []string
}

func (fs *FamilySystem) analyzeMotherDaughter(analysis *RelationshipAnalysis, role1, role2 *IdentifiedTraumaRoles) {
	if role1 != nil && role2 != nil {
		if role1.PrimaryRole == TraumaRoleScapegoat && role2.PrimaryRole == TraumaRoleScapegoat {
			analysis.TraumaDynamics = append(analysis.TraumaDynamics,
				"Передача роли козла отпущения по женской линии")
			analysis.Warnings = append(analysis.Warnings,
				"Высокий риск повторения материнского сценария")
		}
		if role1.PrimaryRole == TraumaRoleGoldenChild && role2.PrimaryRole == TraumaRoleScapegoat {
			analysis.TraumaDynamics = append(analysis.TraumaDynamics,
				"Классическая нарциссическая диада: любимица vs козел отпущения")
			analysis.Warnings = append(analysis.Warnings,
				"Дочь в роли козла отпущения, мать проецирует свою тень")
		}
	}
	analysis.Recommendations = append(analysis.Recommendations,
		"Терапия женской линии",
		"Работа с сепарацией",
		"Отделение от материнских проекций")
}

func (fs *FamilySystem) analyzeMotherSon(analysis *RelationshipAnalysis, role1, role2 *IdentifiedTraumaRoles) {
	if role2 != nil && role2.PrimaryRole == TraumaRoleEmotionalSpouse {
		analysis.TraumaDynamics = append(analysis.TraumaDynamics,
			"Эмоциональный инцест: сын замещает роль мужа")
		analysis.Warnings = append(analysis.Warnings,
			"КРИТИЧНО: Высокий риск симбиотической связи",
			"Сын не способен на здоровые романтические отношения")
	}
	if role2 != nil && role2.PrimaryRole == TraumaRoleGoldenChild {
		analysis.TraumaDynamics = append(analysis.TraumaDynamics,
			"Сын-любимчик, компенсация неудовлетворенности в браке")
	}
	analysis.Recommendations = append(analysis.Recommendations,
		"Срочная терапия для разрыва симбиоза",
		"Работа с мужской идентичностью",
		"Отделение от материнских ожиданий")
}

func (fs *FamilySystem) analyzeSisterSister(analysis *RelationshipAnalysis, role1, role2 *IdentifiedTraumaRoles) {
	if role1 != nil && role2 != nil {
		if (role1.PrimaryRole == TraumaRoleGoldenChild && role2.PrimaryRole == TraumaRoleScapegoat) ||
			(role2.PrimaryRole == TraumaRoleGoldenChild && role1.PrimaryRole == TraumaRoleScapegoat) {
			analysis.TraumaDynamics = append(analysis.TraumaDynamics,
				"Классическая динамика сестер: любимица и козел отпущения")
		}
	}
	analysis.Recommendations = append(analysis.Recommendations,
		"Терапия сиблинговых отношений",
		"Работа с завистью и соперничеством",
		"Восстановление сестринской связи")
}

func (fs *FamilySystem) analyzeBrotherBrother(analysis *RelationshipAnalysis, role1, role2 *IdentifiedTraumaRoles) {
	if role1 != nil && role2 != nil {
		if (role1.PrimaryRole == TraumaRoleGoldenChild && role2.PrimaryRole == TraumaRoleScapegoat) ||
			(role2.PrimaryRole == TraumaRoleGoldenChild && role1.PrimaryRole == TraumaRoleScapegoat) {
			analysis.TraumaDynamics = append(analysis.TraumaDynamics,
				"Классическая динамика братьев: любимица и козел отпущения")
		}
	}
	analysis.Recommendations = append(analysis.Recommendations,
		"Терапия братьев",
		"Работа с конкуренцией",
		"Восстановление братской связи")
}

func (fs *FamilySystem) analyzeSisterBrother(analysis *RelationshipAnalysis, role1, role2 *IdentifiedTraumaRoles) {
	if role1 != nil && role2 != nil {
		if role1.PrimaryRole == TraumaRoleParentified && role2.PrimaryRole == TraumaRoleLostChild {
			analysis.TraumaDynamics = append(analysis.TraumaDynamics,
				"Сестра-родитель и брат-потерянный ребенок")
		}
	}
	analysis.Recommendations = append(analysis.Recommendations,
		"Освобождение сестры от родительской роли",
		"Активация брата",
		"Баланс ролей")
}

func (fs *FamilySystem) analyzeGrandmotherGranddaughter(analysis *RelationshipAnalysis, role1, role2 *IdentifiedTraumaRoles) {
	analysis.Recommendations = append(analysis.Recommendations,
		"Работа с передачей травмы через поколения",
		"Сепарация от бабушкиных установок")
}

func (fs *FamilySystem) analyzeGrandmotherGrandson(analysis *RelationshipAnalysis, role1, role2 *IdentifiedTraumaRoles) {
	analysis.Recommendations = append(analysis.Recommendations,
		"Осознание компенсации через внука",
		"Работа с мужским началом внука")
}

func (fs *FamilySystem) analyzeGrandfatherGranddaughter(analysis *RelationshipAnalysis, role1, role2 *IdentifiedTraumaRoles) {
	analysis.Recommendations = append(analysis.Recommendations,
		"Проверка на эмоциональный инцест",
		"Работа с женской идентичностью внучки")
}

func (fs *FamilySystem) analyzeGrandfatherGrandson(analysis *RelationshipAnalysis, role1, role2 *IdentifiedTraumaRoles) {
	analysis.Recommendations = append(analysis.Recommendations,
		"Работа с мужской линией",
		"Восстановление мужской идентичности внука")
}

func (fs *FamilySystem) analyzePartners(analysis *RelationshipAnalysis, role1, role2 *IdentifiedTraumaRoles) {
	analysis.Recommendations = append(analysis.Recommendations,
		"Работа с индивидуальными травмами обоих",
		"Установление здоровых границ",
		"Формирование равноправного союза")
}

func (fs *FamilySystem) analyzeExPartners(analysis *RelationshipAnalysis, role1, role2 *IdentifiedTraumaRoles) {
	analysis.Recommendations = append(analysis.Recommendations,
		"Работа с болью расставания",
		"Завершение отношений",
		"Обеспечение безопасности после разрыва")
}
