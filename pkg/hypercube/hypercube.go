// Package hypercube v3.1 - Полная версия с анализом замещающих детей и межпоколенческих паттернов
// Автоматически выявляет травматические роли в семье на основе дат рождения и жизненных событий
// Основано на метафоре фильма "Куб" - 4-мерное пространство жизненных траекторий
package hypercube

// Реэкспорт всех основных типов и функций для удобства использования
// Теперь можно импортировать один пакет и иметь доступ ко всему функционалу

import (
	"time"
	
	// Внутренние пакеты
	"github.com/ruvcoindev/idealcore/pkg/hypercube/core"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/data"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/model"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/analysis"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/protocol"
)

// ==================== КОНСТАНТЫ ====================

const CubeSize = core.CubeSize

type RelationshipType = core.RelationshipType

const (
	RelationshipMotherDaughter     = core.RelationshipMotherDaughter
	RelationshipMotherSon          = core.RelationshipMotherSon
	RelationshipFatherDaughter     = core.RelationshipFatherDaughter
	RelationshipFatherSon          = core.RelationshipFatherSon
	RelationshipSisterSister       = core.RelationshipSisterSister
	RelationshipBrotherBrother     = core.RelationshipBrotherBrother
	RelationshipSisterBrother      = core.RelationshipSisterBrother
	RelationshipGrandmotherGranddaughter = core.RelationshipGrandmotherGranddaughter
	RelationshipGrandmotherGrandson     = core.RelationshipGrandmotherGrandson
	RelationshipGrandfatherGranddaughter = core.RelationshipGrandfatherGranddaughter
	RelationshipGrandfatherGrandson     = core.RelationshipGrandfatherGrandson
	RelationshipHusbandWife         = core.RelationshipHusbandWife
	RelationshipCivilPartnership    = core.RelationshipCivilPartnership
	RelationshipExSpouses           = core.RelationshipExSpouses
	RelationshipFriendly            = core.RelationshipFriendly
	RelationshipWork                = core.RelationshipWork
	RelationshipRomantic            = core.RelationshipRomantic
)

type Gender = core.Gender

const (
	GenderMale   = core.GenderMale
	GenderFemale = core.GenderFemale
	GenderOther  = core.GenderOther
)

type FamilyRole = core.FamilyRole

const (
	RoleMother       = core.RoleMother
	RoleFather       = core.RoleFather
	RoleDaughter     = core.RoleDaughter
	RoleSon          = core.RoleSon
	RoleGrandmother  = core.RoleGrandmother
	RoleGrandfather  = core.RoleGrandfather
	RoleGranddaughter = core.RoleGranddaughter
	RoleGrandson     = core.RoleGrandson
	RoleHusband      = core.RoleHusband
	RoleWife         = core.RoleWife
	RolePartner      = core.RolePartner
)

type TraumaRole = core.TraumaRole

const (
	TraumaRoleGoldenChild     = core.TraumaRoleGoldenChild
	TraumaRoleScapegoat       = core.TraumaRoleScapegoat
	TraumaRoleLostChild       = core.TraumaRoleLostChild
	TraumaRoleInvisible       = core.TraumaRoleInvisible
	TraumaRoleGlassChild      = core.TraumaRoleGlassChild
	TraumaRoleParentified     = core.TraumaRoleParentified
	TraumaRoleMascot          = core.TraumaRoleMascot
	TraumaRoleEmotionalSpouse = core.TraumaRoleEmotionalSpouse
	TraumaRoleTruthTeller     = core.TraumaRoleTruthTeller
	TraumaRoleShadow          = core.TraumaRoleShadow
	TraumaRoleHealer          = core.TraumaRoleHealer
	TraumaRoleReplacement     = core.TraumaRoleReplacement
	TraumaRoleGhost           = core.TraumaRoleGhost
	TraumaRoleAncestor        = core.TraumaRoleAncestor
)

type LossType = core.LossType

const (
	LossTypeAbortion    = core.LossTypeAbortion
	LossTypeMiscarriage = core.LossTypeMiscarriage
	LossTypeStillborn   = core.LossTypeStillborn
	LossTypeInfantDeath = core.LossTypeInfantDeath
	LossTypeChildDeath  = core.LossTypeChildDeath
	LossTypeAdultDeath  = core.LossTypeAdultDeath
)

var TraumaThresholds = core.TraumaThresholds

// ==================== СТРУКТУРЫ ДАННЫХ ====================

type HypercubeCoords = core.HypercubeCoords
type PersonVectors = core.PersonVectors
type LifeEvent = data.LifeEvent
type Disease = data.Disease
type Infidelity = data.Infidelity
type KarmicDebt = data.KarmicDebt
type ShamePattern = data.ShamePattern
type LostChild = data.LostChild
type ReplacementChildData = data.ReplacementChildData
type ExtendedFamilyMember = data.ExtendedFamilyMember
type IdentifiedTraumaRoles = model.IdentifiedTraumaRoles
type FamilyPattern = model.FamilyPattern
type FamilySystem = model.FamilySystem
type CompatibilityResult = analysis.CompatibilityResult
type GroupCompatibilityResult = analysis.GroupCompatibilityResult
type ShameAnalysisResult = analysis.ShameAnalysisResult
type ReplacementPattern = analysis.ReplacementPattern
type ReplacementAnalysisResult = analysis.ReplacementAnalysisResult
type LargeFamilyAnalysisResult = analysis.LargeFamilyAnalysisResult
type SeparationStage = protocol.SeparationStage
type SeparationProtocol = protocol.SeparationProtocol
type TherapyRecommendation = protocol.TherapyRecommendation
type TherapyPlan = protocol.TherapyPlan
type FamilyReport = protocol.FamilyReport

// ==================== ФУНКЦИИ РАБОТЫ С ДАТАМИ ====================

// ParseDate парсит дату в формате ДД.ММ.ГГГГ или ДД.М.ГГГГ
func ParseDate(dateStr string) (time.Time, error) {
	return core.ParseDate(dateStr)
}

// ParseDateToCoords преобразует дату в координаты гиперкуба
func ParseDateToCoords(date time.Time) HypercubeCoords {
	return core.ParseDateToCoords(date)
}

// ==================== ФУНКЦИИ СОЗДАНИЯ ЧЛЕНОВ СЕМЬИ ====================

// CreateFamilyMember создает нового члена семьи
func CreateFamilyMember(id, name string, gender Gender, birthDateStr string) (*ExtendedFamilyMember, error) {
	return data.CreateFamilyMember(id, name, gender, birthDateStr)
}

// NewFamilySystem создает новую семейную систему
func NewFamilySystem(extendedMode bool) *FamilySystem {
	return model.NewFamilySystem(extendedMode)
}

// ==================== МЕТОДЫ РАБОТЫ С СЕМЕЙНОЙ СИСТЕМОЙ ====================

// AddMember добавляет члена семьи
func (fs *FamilySystem) AddMember(member *ExtendedFamilyMember) {
	fs.AddMember(member)
}

// GetMember возвращает члена семьи по ID
func (fs *FamilySystem) GetMember(id string) *ExtendedFamilyMember {
	return fs.GetMember(id)
}

// GetSiblings возвращает список сиблингов
func (fs *FamilySystem) GetSiblings(memberID string) []string {
	return fs.GetSiblings(memberID)
}

// GetBirthOrder возвращает порядок рождения
func (fs *FamilySystem) GetBirthOrder(memberID string) int {
	return fs.GetBirthOrder(memberID)
}

// IdentifyTraumaRoles выявляет травматические роли
func (fs *FamilySystem) IdentifyTraumaRoles(memberID string) *IdentifiedTraumaRoles {
	return fs.IdentifyTraumaRoles(memberID)
}

// FindFamilyPatterns выявляет семейные паттерны
func (fs *FamilySystem) FindFamilyPatterns() []FamilyPattern {
	return fs.FindFamilyPatterns()
}

// GenerateFamilyReport генерирует полный отчет
func (fs *FamilySystem) GenerateFamilyReport() string {
	return protocol.GenerateFamilyReport(fs)
}

// ==================== АНАЛИТИЧЕСКИЕ ФУНКЦИИ ====================

// CalculateVectors вычисляет векторы для координат
func CalculateVectors(coords HypercubeCoords) PersonVectors {
	return core.CalculateVectors(coords)
}

// CalculateCompatibility вычисляет совместимость двух людей
func CalculateCompatibility(coords1, coords2 HypercubeCoords, vectors1, vectors2 PersonVectors) CompatibilityResult {
	return analysis.CalculateCompatibility(coords1, coords2, vectors1, vectors2)
}

// CalculateGroupCompatibility анализирует совместимость группы
func CalculateGroupCompatibility(members []*ExtendedFamilyMember) GroupCompatibilityResult {
	return analysis.CalculateGroupCompatibility(members)
}

// CalculateBodyShame вычисляет уровень телесного стыда
func CalculateBodyShame(coords HypercubeCoords, traumaRole TraumaRole, ageAtEvent int, member *ExtendedFamilyMember) float64 {
	return analysis.CalculateBodyShame(coords, traumaRole, ageAtEvent, member)
}

// AnalyzeShamePattern определяет паттерн стыда
func AnalyzeShamePattern(vectors PersonVectors, traumaRole TraumaRole) *ShamePattern {
	return analysis.AnalyzeShamePattern(vectors, traumaRole)
}

// FindReplacementPatterns выявляет паттерны замещающих детей
func FindReplacementPatterns(fs *FamilySystem) []ReplacementPattern {
	return analysis.FindReplacementPatterns(fs)
}

// AnalyzeReplacementDynamics анализирует динамику замещения
func AnalyzeReplacementDynamics(fs *FamilySystem) ReplacementAnalysisResult {
	return analysis.AnalyzeReplacementDynamics(fs)
}

// AnalyzeLargeFamily анализирует большие семьи
func AnalyzeLargeFamily(fs *FamilySystem) *LargeFamilyAnalysisResult {
	return analysis.AnalyzeLargeFamily(fs)
}

// ==================== ФУНКЦИИ РАБОТЫ С ТЕРАПИЕЙ ====================

// GenerateSeparationProtocol создает протокол разделения
func GenerateSeparationProtocol(p1, p2 *ExtendedFamilyMember, relType RelationshipType, compat CompatibilityResult) SeparationProtocol {
	return protocol.GenerateSeparationProtocol(p1, p2, relType, compat)
}

// GenerateTherapyPlan создает план терапии
func GenerateTherapyPlan(member *ExtendedFamilyMember, role *IdentifiedTraumaRoles) TherapyPlan {
	return protocol.GenerateTherapyPlan(member, role)
}

// FormatTherapyPlan форматирует план терапии
func FormatTherapyPlan(plan TherapyPlan) string {
	return protocol.FormatTherapyPlan(plan)
}

// GenerateShortReport генерирует краткий отчет
func GenerateShortReport(fs *FamilySystem) string {
	return protocol.GenerateShortReport(fs)
}
