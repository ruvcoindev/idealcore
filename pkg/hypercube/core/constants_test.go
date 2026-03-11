package core

import (
	"testing"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant interface{}
		expected interface{}
	}{
		{"CubeSize", CubeSize, 26},
		{"GenderMale", GenderMale, Gender("male")},
		{"GenderFemale", GenderFemale, Gender("female")},
		{"RoleMother", RoleMother, FamilyRole("mother")},
		{"RoleFather", RoleFather, FamilyRole("father")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("got %v, want %v", tt.constant, tt.expected)
			}
		})
	}
}

func TestTraumaThresholds(t *testing.T) {
	if TraumaThresholds.ErraticMin != 3 {
		t.Errorf("ErraticMin = %d, want 3", TraumaThresholds.ErraticMin)
	}
	if TraumaThresholds.FrozenMax != 1 {
		t.Errorf("FrozenMax = %d, want 1", TraumaThresholds.FrozenMax)
	}
	if TraumaThresholds.RepetitiveCount != 2 {
		t.Errorf("RepetitiveCount = %d, want 2", TraumaThresholds.RepetitiveCount)
	}
}

func TestRelationshipTypes(t *testing.T) {
	types := []RelationshipType{
		RelationshipMotherDaughter,
		RelationshipMotherSon,
		RelationshipFatherDaughter,
		RelationshipFatherSon,
		RelationshipSisterSister,
		RelationshipBrotherBrother,
		RelationshipSisterBrother,
	}

	for i, rt := range types {
		if rt == "" {
			t.Errorf("Relationship type at index %d is empty", i)
		}
	}
}
