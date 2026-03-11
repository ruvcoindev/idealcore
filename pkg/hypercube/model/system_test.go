package model

import (
	"testing"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/core"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/data"
)

func TestNewFamilySystem(t *testing.T) {
	fs := NewFamilySystem(false)
	
	if fs == nil {
		t.Fatal("NewFamilySystem returned nil")
	}
	if len(fs.Members) != 0 {
		t.Errorf("Members length = %d, want 0", len(fs.Members))
	}
	if fs.ExtendedMode != false {
		t.Error("ExtendedMode should be false")
	}

	fsExtended := NewFamilySystem(true)
	if !fsExtended.ExtendedMode {
		t.Error("ExtendedMode should be true")
	}
}

func TestAddMember(t *testing.T) {
	fs := NewFamilySystem(false)
	
	member1, _ := data.CreateFamilyMember("m1", "Мать", core.GenderFemale, "15.05.1970")
	member2, _ := data.CreateFamilyMember("c1", "Ребенок", core.GenderFemale, "15.06.1995")
	member2.Parents = []string{"m1"}
	member1.Children = []string{"c1"}
	
	fs.AddMember(member1)
	fs.AddMember(member2)
	
	if len(fs.Members) != 2 {
		t.Errorf("Members length = %d, want 2", len(fs.Members))
	}
	
	// Проверяем поколения
	if member1.Generation != 1 {
		t.Errorf("Mother generation = %d, want 1", member1.Generation)
	}
	if member2.Generation != 2 {
		t.Errorf("Child generation = %d, want 2", member2.Generation)
	}
}

func TestGetMember(t *testing.T) {
	fs := NewFamilySystem(false)
	member, _ := data.CreateFamilyMember("test1", "Тест", core.GenderMale, "15.05.1990")
	fs.AddMember(member)
	
	result := fs.GetMember("test1")
	if result == nil {
		t.Fatal("GetMember returned nil")
	}
	if result.ID != "test1" {
		t.Errorf("Member ID = %s, want test1", result.ID)
	}
	
	notFound := fs.GetMember("nonexistent")
	if notFound != nil {
		t.Error("GetMember for nonexistent ID should return nil")
	}
}

func TestGetSiblings(t *testing.T) {
	fs := NewFamilySystem(false)
	
	parent, _ := data.CreateFamilyMember("p1", "Родитель", core.GenderFemale, "15.05.1970")
	child1, _ := data.CreateFamilyMember("c1", "Ребенок1", core.GenderFemale, "15.06.1995")
	child2, _ := data.CreateFamilyMember("c2", "Ребенок2", core.GenderFemale, "15.06.1998")
	child3, _ := data.CreateFamilyMember("c3", "Ребенок3", core.GenderFemale, "15.06.2001")
	
	child1.Parents = []string{"p1"}
	child2.Parents = []string{"p1"}
	child3.Parents = []string{"p1"}
	parent.Children = []string{"c1", "c2", "c3"}
	
	fs.AddMember(parent)
	fs.AddMember(child1)
	fs.AddMember(child2)
	fs.AddMember(child3)
	
	siblings := fs.GetSiblings("c1")
	if len(siblings) != 2 {
		t.Errorf("Siblings length = %d, want 2", len(siblings))
	}
	
	// Проверяем, что среди сиблингов нет самого себя
	for _, sib := range siblings {
		if sib == "c1" {
			t.Error("GetSiblings returned self")
		}
	}
}

func TestGetBirthOrder(t *testing.T) {
	fs := NewFamilySystem(false)
	
	parent, _ := data.CreateFamilyMember("p1", "Родитель", core.GenderFemale, "15.05.1970")
	child1, _ := data.CreateFamilyMember("c1", "Ребенок1", core.GenderFemale, "15.06.1995")
	child2, _ := data.CreateFamilyMember("c2", "Ребенок2", core.GenderFemale, "15.06.1998")
	child3, _ := data.CreateFamilyMember("c3", "Ребенок3", core.GenderFemale, "15.06.2001")
	
	child1.Parents = []string{"p1"}
	child2.Parents = []string{"p1"}
	child3.Parents = []string{"p1"}
	parent.Children = []string{"c1", "c2", "c3"}
	
	fs.AddMember(parent)
	fs.AddMember(child1)
	fs.AddMember(child2)
	fs.AddMember(child3)
	
	order1 := fs.GetBirthOrder("c1")
	if order1 != 1 {
		t.Errorf("Birth order of c1 = %d, want 1", order1)
	}
	
	order2 := fs.GetBirthOrder("c2")
	if order2 != 2 {
		t.Errorf("Birth order of c2 = %d, want 2", order2)
	}
	
	order3 := fs.GetBirthOrder("c3")
	if order3 != 3 {
		t.Errorf("Birth order of c3 = %d, want 3", order3)
	}
}

func TestGetGenerationSpread(t *testing.T) {
	fs := NewFamilySystem(false)
	
	g1, _ := data.CreateFamilyMember("g1", "Поколение1", core.GenderFemale, "15.05.1950")
	g2, _ := data.CreateFamilyMember("g2", "Поколение2", core.GenderFemale, "15.05.1975")
	g3, _ := data.CreateFamilyMember("g3", "Поколение3", core.GenderFemale, "15.05.2000")
	
	g2.Parents = []string{"g1"}
	g3.Parents = []string{"g2"}
	
	fs.AddMember(g1)
	fs.AddMember(g2)
	fs.AddMember(g3)
	
	spread := fs.GetGenerationSpread()
	if spread != 3 {
		t.Errorf("Generation spread = %d, want 3", spread)
	}
}

func TestGetMembersByGeneration(t *testing.T) {
	fs := NewFamilySystem(false)
	
	g1_1, _ := data.CreateFamilyMember("g1_1", "Поколение1_1", core.GenderFemale, "15.05.1950")
	g1_2, _ := data.CreateFamilyMember("g1_2", "Поколение1_2", core.GenderMale, "15.05.1952")
	g2, _ := data.CreateFamilyMember("g2", "Поколение2", core.GenderFemale, "15.05.1975")
	
	g2.Parents = []string{"g1_1", "g1_2"}
	
	fs.AddMember(g1_1)
	fs.AddMember(g1_2)
	fs.AddMember(g2)
	
	gen1 := fs.GetMembersByGeneration(1)
	if len(gen1) != 2 {
		t.Errorf("Generation 1 members = %d, want 2", len(gen1))
	}
	
	gen2 := fs.GetMembersByGeneration(2)
	if len(gen2) != 1 {
		t.Errorf("Generation 2 members = %d, want 1", len(gen2))
	}
	
	gen3 := fs.GetMembersByGeneration(3)
	if len(gen3) != 0 {
		t.Errorf("Generation 3 members = %d, want 0", len(gen3))
	}
}
