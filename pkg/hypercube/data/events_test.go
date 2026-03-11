package data

import (
	"testing"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/core"
)

func TestAddLifeEvent(t *testing.T) {
	member, _ := CreateFamilyMember("test1", "Тестовый", core.GenderMale, core.RoleSon, "15.05.1990")
	
	err := member.AddLifeEvent("marriage", "Свадьба", "15.06.2015", "partner1", 0.5)
	if err != nil {
		t.Fatalf("AddLifeEvent failed: %v", err)
	}

	if len(member.Events) != 1 {
		t.Errorf("Events length = %d, want 1", len(member.Events))
	}

	event := member.Events[0]
	if event.EventType != "marriage" {
		t.Errorf("EventType = %s, want marriage", event.EventType)
	}
	if event.WithPerson != "partner1" {
		t.Errorf("WithPerson = %s, want partner1", event.WithPerson)
	}
	if event.Impact != 0.5 {
		t.Errorf("Impact = %f, want 0.5", event.Impact)
	}

	// Проверяем, что партнер добавился
	if len(member.Partners) != 1 {
		t.Errorf("Partners length = %d, want 1", len(member.Partners))
	}
	if member.CurrentPartner != "partner1" {
		t.Errorf("CurrentPartner = %s, want partner1", member.CurrentPartner)
	}
}

func TestAddDivorce(t *testing.T) {
	member, _ := CreateFamilyMember("test1", "Тестовый", core.GenderMale, core.RoleSon, "15.05.1990")
	member.Partners = []string{"partner1"}
	member.CurrentPartner = "partner1"
	
	err := member.AddLifeEvent("divorce", "Развод", "15.06.2018", "partner1", 0.7)
	if err != nil {
		t.Fatalf("AddLifeEvent failed: %v", err)
	}

	if member.CurrentPartner != "" {
		t.Errorf("CurrentPartner = %s, want empty", member.CurrentPartner)
	}
}

func TestAddDisease(t *testing.T) {
	member, _ := CreateFamilyMember("test1", "Тестовый", core.GenderMale, core.RoleSon, "15.05.1990")
	
	member.AddDisease("Диабет", 0.6, "chronic", true, "parent1")

	if len(member.Diseases) != 1 {
		t.Errorf("Diseases length = %d, want 1", len(member.Diseases))
	}

	disease := member.Diseases[0]
	if disease.Name != "Диабет" {
		t.Errorf("Disease name = %s, want Диабет", disease.Name)
	}
	if disease.Severity != 0.6 {
		t.Errorf("Severity = %f, want 0.6", disease.Severity)
	}
	if !disease.IsGenetic {
		t.Error("IsGenetic should be true")
	}
	if disease.FromParent != "parent1" {
		t.Errorf("FromParent = %s, want parent1", disease.FromParent)
	}
}

func TestAddInfidelity(t *testing.T) {
	member, _ := CreateFamilyMember("test1", "Тестовый", core.GenderMale, core.RoleSon, "15.05.1990")
	
	err := member.AddInfidelity("15.06.2015", "partner1", "other1", true, 0.8, "Случайная связь")
	if err != nil {
		t.Fatalf("AddInfidelity failed: %v", err)
	}

	if len(member.Infidelities) != 1 {
		t.Errorf("Infidelities length = %d, want 1", len(member.Infidelities))
	}

	infidelity := member.Infidelities[0]
	if infidelity.PartnerAtTime != "partner1" {
		t.Errorf("PartnerAtTime = %s, want partner1", infidelity.PartnerAtTime)
	}
	if !infidelity.WasDiscovered {
		t.Error("WasDiscovered should be true")
	}
	if infidelity.Impact != 0.8 {
		t.Errorf("Impact = %f, want 0.8", infidelity.Impact)
	}

	// Проверяем, что добавился кармический долг
	if len(member.KarmicDebts) != 1 {
		t.Errorf("KarmicDebts length = %d, want 1", len(member.KarmicDebts))
	}
}
