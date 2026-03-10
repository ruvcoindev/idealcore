package data

import (
	"time"
	"github.com/ruvcoindev/idealcore/hypercube/core"
)

// LostChild представляет потерянного ребенка
// Это может быть аборт, выкидыш, мертворожденный или умерший ребенок
// Такие дети становятся "призраками", влияющими на семью годами
type LostChild struct {
	ID             string         // Уникальный идентификатор потери
	LossType       core.LossType  // Тип потери (аборт, выкидыш, смерть)
	LossDate       time.Time      // Дата потери (или дата аборта)
	DueDate        *time.Time     // Предполагаемая дата родов (для абортов/выкидышей)
	GestationWeeks int            // Неделя беременности (для абортов/выкидышей)
	MotherID       string         // ID матери
	FatherID       string         // ID отца (если известен)
	Name           string         // Имя, которое дали (если давали)
	IsNamed        bool           // Было ли дано имя
	WasBuried      bool           // Были ли похороны
	GriefProcessed bool           // Было ли горе проработано
	AgeAtDeath     int            // Возраст на момент смерти (для детей)
	CauseOfDeath   string         // Причина смерти (если известна)
	Notes          string         // Дополнительные примечания
}

// ReplacementChildData содержит данные о замещающем ребенке
// Замещающий ребенок - это ребенок, рожденный после потери предыдущего
// Он часто живет с чувством, что должен "заменить" умершего
type ReplacementChildData struct {
	ChildID        string    // ID замещающего ребенка
	ReplacesID     string    // ID потерянного ребенка, которого замещает
	ConceptionDate time.Time // Дата зачатия (примерная)
	BirthDate      time.Time // Дата рождения
	TimeGapMonths  int       // Разрыв в месяцах между потерей и зачатием
	ParentsAware   bool      // Осознают ли родители, что ребенок замещающий
	IsExplicit     bool      // Явное ли замещение (дали то же имя, ожидают того же)
	BurdenLevel    float64   // Уровень нагрузки на ребенка (0-1) - вычисляется
	// Чем меньше разрыв, тем выше нагрузка
	// Явное замещение повышает нагрузку
	// Неосознанное родителями замещение тоже повышает нагрузку
	Notes          string    // Дополнительные примечания
}

// AddLostChild добавляет запись о потерянном ребенке
func (m *ExtendedFamilyMember) AddLostChild(lossType core.LossType, lossDateStr string, gestationWeeks int, fatherID string) error {
	lossDate, err := time.Parse("02.01.2006", lossDateStr)
	if err != nil {
		return err
	}
	
	// Генерируем ID для потерянного ребенка
	lostID := "lost_" + m.ID + "_" + lossDateStr
	
	lostChild := LostChild{
		ID:             lostID,
		LossType:       lossType,
		LossDate:       lossDate,
		GestationWeeks: gestationWeeks,
		MotherID:       m.ID,
		FatherID:       fatherID,
		IsNamed:        false,
		WasBuried:      false,
		GriefProcessed: false,
	}
	
	// Примерная дата родов (если известна неделя)
	if gestationWeeks > 0 {
		dueDate := lossDate.AddDate(0, 0, (40-gestationWeeks)*7)
		lostChild.DueDate = &dueDate
	}
	
	m.LostChildren = append(m.LostChildren, lostChild)
	
	// Обновляем счетчики
	switch lossType {
	case core.LossTypeAbortion:
		m.AbortionCount++
	case core.LossTypeMiscarriage:
		m.MiscarriageCount++
	}
	
	// Добавляем кармический долг
	debt := KarmicDebt{
		DebtorID:    m.ID,
		CreditorID:  lostID,
		Nature:      "replacement",
		Weight:      0.9,
		IsResolved:  false,
		Description: "Потеря ребенка",
	}
	
	m.KarmicDebts = append(m.KarmicDebts, debt)
	
	return nil
}

// MarkAsReplacement помечает ребенка как замещающего
// Вызывается, когда ребенок рождается после потери
func (m *ExtendedFamilyMember) MarkAsReplacement(replacesID string, conceptionDateStr string, parentsAware, isExplicit bool) error {
	conceptionDate, err := time.Parse("02.01.2006", conceptionDateStr)
	if err != nil {
		return err
	}
	
	// Ищем потерянного ребенка, которого замещаем
	// В реальном коде здесь должен быть поиск по всем членам семьи
	
	m.IsReplacement = true
	m.ReplacesWho = replacesID
	
	// Вычисляем разрыв в месяцах (упрощенно)
	gapMonths := 0 // В реальном коде нужно вычислить разницу между датой потери и зачатием
	
	// Вычисляем уровень нагрузки
	burdenLevel := 0.5 // Базовый уровень
	
	if gapMonths < 6 {
		burdenLevel += 0.3 // Слишком быстрая замена
	}
	if isExplicit {
		burdenLevel += 0.2 // Явное замещение (то же имя)
	}
	if !parentsAware {
		burdenLevel += 0.2 // Неосознанное замещение
	}
	
	if burdenLevel > 1.0 {
		burdenLevel = 1.0
	}
	
	m.ReplacementData = &ReplacementChildData{
		ChildID:        m.ID,
		ReplacesID:     replacesID,
		ConceptionDate: conceptionDate,
		BirthDate:      m.BirthDate,
		TimeGapMonths:  gapMonths,
		ParentsAware:   parentsAware,
		IsExplicit:     isExplicit,
		BurdenLevel:    burdenLevel,
	}
	
	return nil
}

// CalculateReplacementBurden вычисляет нагрузку на замещающего ребенка
// Используется, когда не все данные известны
func (m *ExtendedFamilyMember) CalculateReplacementBurden() float64 {
	if !m.IsReplacement || m.ReplacementData == nil {
		return 0
	}
	
	burden := 0.5 // Базовый уровень
	
	if m.ReplacementData.TimeGapMonths < 6 {
		burden += 0.3
	}
	if m.ReplacementData.IsExplicit {
		burden += 0.2
	}
	if !m.ReplacementData.ParentsAware {
		burden += 0.2
	}
	
	// Вина выжившего добавляет нагрузку
	burden += m.SurvivorGuilt * 0.3
	
	if burden > 1.0 {
		burden = 1.0
	}
	
	return burden
}
