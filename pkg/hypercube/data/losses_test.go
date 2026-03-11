package data

import (
	"testing"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/core" // ИСПРАВЛЕНО: добавил pkg/
)

func TestAddLostChild(t *testing.T) {
	member, _ := CreateFamilyMember("mother", "Мать", core.GenderFemale, core.RoleMother, "15.05.1970") // ИСПРАВЛЕНО: добавил роль
	
	err := member.AddLostChild(core.LossTypeMiscarriage, "10.08.1995", 12, "father1")
	if err != nil {
		t.Fatalf("AddLostChild failed: %v", err)
	}

	if len(member.LostChildren) != 1 {
		t.Errorf("LostChildren length = %d, want 1", len(member.LostChildren))
	}

	lost := member.LostChildren[0]
	if lost.LossType != core.LossTypeMiscarriage {
		t.Errorf("LossType = %s, want miscarriage", lost.LossType)
	}
	if lost.GestationWeeks != 12 {
		t.Errorf("GestationWeeks = %d, want 12", lost.GestationWeeks)
	}
	if lost.MotherID != "mother" {
		t.Errorf("MotherID = %s, want mother", lost.MotherID)
	}
	if lost.FatherID != "father1" {
		t.Errorf("FatherID = %s, want father1", lost.FatherID)
	}

	// Проверяем счетчики
	if member.MiscarriageCount != 1 {
		t.Errorf("MiscarriageCount = %d, want 1", member.MiscarriageCount)
	}

	// Проверяем кармический долг
	if len(member.KarmicDebts) != 1 {
		t.Errorf("KarmicDebts length = %d, want 1", len(member.KarmicDebts))
	}
}

func TestAddAbortion(t *testing.T) {
	member, _ := CreateFamilyMember("woman", "Женщина", core.GenderFemale, core.RoleDaughter, "15.05.1970") // ИСПРАВЛЕНО: добавил роль
	
	err := member.AddLostChild(core.LossTypeAbortion, "10.08.1995", 8, "partner1")
	if err != nil {
		t.Fatalf("AddLostChild failed: %v", err)
	}

	if member.AbortionCount != 1 {
		t.Errorf("AbortionCount = %d, want 1", member.AbortionCount)
	}
}

func TestMarkAsReplacement(t *testing.T) {
	// Создаем мать с потерянным ребенком
	mother, _ := CreateFamilyMember("mother", "Мать", core.GenderFemale, core.RoleMother, "15.05.1970") // ИСПРАВЛЕНО: добавил роль
	mother.AddLostChild(core.LossTypeMiscarriage, "10.08.1995", 12, "father1")
	
	// Создаем ребенка
	child, _ := CreateFamilyMember("child", "Ребенок", core.GenderFemale, core.RoleDaughter, "15.06.1996") // ИСПРАВЛЕНО: добавил роль
	child.Parents = []string{"mother"}
	
	// В реальном коде нужно добавить мать в систему, но для теста упростим
	// Помечаем как замещающего
	err := child.MarkAsReplacement("lost_mother_10.08.1995", "01.01.1996", false, true)
	if err != nil {
		t.Fatalf("MarkAsReplacement failed: %v", err)
	}

	if !child.IsReplacement {
		t.Error("IsReplacement should be true")
	}
	if child.ReplacementData == nil {
		t.Fatal("ReplacementData should not be nil")
	}

	if child.ReplacementData.IsExplicit != true {
		t.Error("IsExplicit should be true")
	}
	if child.ReplacementData.ParentsAware != false {
		t.Error("ParentsAware should be false")
	}
}

func TestCalculateReplacementBurden(t *testing.T) {
	child, _ := CreateFamilyMember("child", "Ребенок", core.GenderFemale, core.RoleDaughter, "15.06.1996")
	
	// Сначала без данных
	burden := child.CalculateReplacementBurden()
	if burden != 0 {
		t.Errorf("Burden without data = %f, want 0", burden)
	}

	// Помечаем как замещающего
	child.IsReplacement = true
	child.ReplacementData = &ReplacementChildData{
		TimeGapMonths: 5,
		IsExplicit:    true,
		ParentsAware:  false,
	}
	
	burden = child.CalculateReplacementBurden()
	// Из-за ограничения burden <= 1.0 в коде, ожидаем 1.0
	if burden < 0.9 || burden > 1.0 {
		t.Errorf("Burden = %f, want around 1.0", burden)
	}

	// С виной выжившего
	child.SurvivorGuilt = 0.5
	burden = child.CalculateReplacementBurden()
	if burden > 1.0 {
		t.Errorf("Burden > 1.0: %f", burden)
	}
}
