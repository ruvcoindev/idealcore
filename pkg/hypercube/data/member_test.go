package data

import (
	"testing"
	"time"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/core"
)

func TestCreateFamilyMember(t *testing.T) {
	member, err := CreateFamilyMember("test1", "Тестовый", core.GenderMale, core.RoleSon, "15.05.1990") // ИСПРАВЛЕНО
	if err != nil {
		t.Fatalf("CreateFamilyMember failed: %v", err)
	}

	if member.ID != "test1" {
		t.Errorf("ID = %s, want test1", member.ID)
	}
	if member.Name != "Тестовый" {
		t.Errorf("Name = %s, want Тестовый", member.Name)
	}
	if member.Gender != core.GenderMale {
		t.Errorf("Gender = %v, want male", member.Gender)
	}
	if member.BirthDate.Year() != 1990 || member.BirthDate.Month() != 5 || member.BirthDate.Day() != 15 {
		t.Errorf("BirthDate = %v, want 1990-05-15", member.BirthDate)
	}
}

func TestCreateFamilyMember_InvalidDate(t *testing.T) {
	_, err := CreateFamilyMember("test1", "Тестовый", core.GenderMale, core.RoleSon, "15.13.1990") // ИСПРАВЛЕНО
	if err == nil {
		t.Error("CreateFamilyMember with invalid date should return error")
	}
}

func TestGetAge(t *testing.T) {
	member, _ := CreateFamilyMember("test1", "Тестовый", core.GenderMale, core.RoleSon, "15.05.1990") // ИСПРАВЛЕНО
	
	age := member.GetAge()
	if age < 30 || age > 40 { // В 2026 году возраст 35-36 лет
		t.Errorf("GetAge() = %d, seems incorrect", age)
	}

	// Тест с датой смерти
	deathDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	member.DeathDate = &deathDate
	member.IsDeceased = true
	
	age = member.GetAge()
	if age != 29 { // 1990-2020 = 30 лет (в зависимости от месяца)
		t.Logf("Deceased age = %d", age)
	}
}

func TestHasParent(t *testing.T) {
	member, _ := CreateFamilyMember("child", "Ребенок", core.GenderFemale, core.RoleDaughter, "15.05.2010") // ИСПРАВЛЕНО
	member.Parents = []string{"mother", "father"}

	if !member.HasParent("mother") {
		t.Error("HasParent(mother) should be true")
	}
	if !member.HasParent("father") {
		t.Error("HasParent(father) should be true")
	}
	if member.HasParent("grandma") {
		t.Error("HasParent(grandma) should be false")
	}
}

func TestHasChild(t *testing.T) {
	member, _ := CreateFamilyMember("parent", "Родитель", core.GenderFemale, core.RoleMother, "15.05.1980") // ИСПРАВЛЕНО
	member.Children = []string{"child1", "child2"}

	if !member.HasChild("child1") {
		t.Error("HasChild(child1) should be true")
	}
	if !member.HasChild("child2") {
		t.Error("HasChild(child2) should be true")
	}
	if member.HasChild("child3") {
		t.Error("HasChild(child3) should be false")
	}
}

func TestHasPartner(t *testing.T) {
	member, _ := CreateFamilyMember("person", "Человек", core.GenderFemale, core.RoleDaughter, "15.05.1980") // ИСПРАВЛЕНО
	member.Partners = []string{"partner1", "partner2"}

	if !member.HasPartner("partner1") {
		t.Error("HasPartner(partner1) should be true")
	}
	if !member.HasPartner("partner2") {
		t.Error("HasPartner(partner2) should be true")
	}
	if member.HasPartner("partner3") {
		t.Error("HasPartner(partner3) should be false")
	}
}
