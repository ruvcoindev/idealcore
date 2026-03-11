package core

import (
	"time"
	"fmt"
)

// CubeSize — размер куба (26 комнат по каждой оси, как в фильме "Куб")
// Каждая координата может принимать значения от 0 до 25, что символизирует
// ограниченность жизненного пространства и повторяющиеся паттерны
const CubeSize = 26

// RelationshipType определяет тип отношений между двумя людьми
type RelationshipType string

const (
	RelationshipMotherDaughter RelationshipType = "mother-daughter"
	RelationshipMotherSon      RelationshipType = "mother-son"
	RelationshipFatherDaughter RelationshipType = "father-daughter"
	RelationshipFatherSon      RelationshipType = "father-son"
	RelationshipSisterSister   RelationshipType = "sister-sister"
	RelationshipBrotherBrother  RelationshipType = "brother-brother"
	RelationshipSisterBrother   RelationshipType = "sister-brother"
	RelationshipGrandmotherGranddaughter RelationshipType = "grandmother-granddaughter"
	RelationshipGrandmotherGrandson      RelationshipType = "grandmother-grandson"
	RelationshipGrandfatherGranddaughter RelationshipType = "grandfather-granddaughter"
	RelationshipGrandfatherGrandson      RelationshipType = "grandfather-grandson"
	RelationshipHusbandWife      RelationshipType = "husband-wife"
	RelationshipCivilPartnership RelationshipType = "civil-partnership"
	RelationshipExSpouses        RelationshipType = "ex-spouses"
	RelationshipRomantic         RelationshipType = "romantic"
	RelationshipFriendly         RelationshipType = "friendly"
	RelationshipWork             RelationshipType = "work"
)

// Gender определяет пол человека
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// FamilyRole определяет базовую роль в семье
type FamilyRole string

const (
	RoleMother       FamilyRole = "mother"
	RoleFather       FamilyRole = "father"
	RoleDaughter     FamilyRole = "daughter"
	RoleSon          FamilyRole = "son"
	RoleGrandmother  FamilyRole = "grandmother"
	RoleGrandfather  FamilyRole = "grandfather"
	RoleGranddaughter FamilyRole = "granddaughter"
	RoleGrandson     FamilyRole = "grandson"
	RoleHusband      FamilyRole = "husband"
	RoleWife         FamilyRole = "wife"
	RolePartner      FamilyRole = "partner"
)

// TraumaRole определяет роль в травматической системе
type TraumaRole string

const (
	TraumaRoleGoldenChild     TraumaRole = "golden_child"
	TraumaRoleScapegoat       TraumaRole = "scapegoat"
	TraumaRoleLostChild       TraumaRole = "lost_child"
	TraumaRoleInvisible       TraumaRole = "invisible"
	TraumaRoleGlassChild      TraumaRole = "glass_child"
	TraumaRoleParentified     TraumaRole = "parentified"
	TraumaRoleMascot          TraumaRole = "mascot"
	TraumaRoleEmotionalSpouse TraumaRole = "emotional_spouse"
	TraumaRoleTruthTeller     TraumaRole = "truth_teller"
	TraumaRoleShadow          TraumaRole = "shadow"
	TraumaRoleHealer          TraumaRole = "healer"
	TraumaRoleReplacement     TraumaRole = "replacement_child"
	TraumaRoleGhost           TraumaRole = "ghost"
	TraumaRoleAncestor        TraumaRole = "ancestor"
)

// LossType определяет тип потери
type LossType string

const (
	LossTypeAbortion    LossType = "abortion"
	LossTypeMiscarriage LossType = "miscarriage"
	LossTypeStillborn   LossType = "stillborn"
	LossTypeInfantDeath LossType = "infant_death"
	LossTypeChildDeath  LossType = "child_death"
	LossTypeAdultDeath  LossType = "adult_death"
)

// TraumaThresholds — пороговые значения для определения травматических паттернов
var TraumaThresholds = struct {
	ErraticMin       int
	FrozenMax        int
	RepetitiveCount  int
}{
	ErraticMin:      3,
	FrozenMax:       1,
	RepetitiveCount: 2,
}

// IdentifiedTraumaRoles содержит выявленные травматические роли
type IdentifiedTraumaRoles struct {
	PrimaryRole   TraumaRole
	SecondaryRoles []TraumaRole
	Confidence    float64
	Evidence      []string
}

// ParseDate парсит дату в формате ДД.ММ.ГГГГ или ДД.М.ГГГГ
// ЭКСПОРТИРУЕМАЯ функция для использования в других пакетах
func ParseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"02.01.2006",
		"2.1.2006",
		"02.1.2006",
		"2.01.2006",
		"2006-01-02",
		"2006-1-2",
	}
	
	var err error
	var t time.Time
	for _, f := range formats {
		t, err = time.Parse(f, dateStr)
		if err == nil {
			return t, nil
		}
	}
	return t, fmt.Errorf("неверный формат даты: %s", dateStr)
}
