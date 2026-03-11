package testhelper

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/core"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/data"
	"github.com/ruvcoindev/idealcore/pkg/hypercube/model"
)

// TestFamily представляет тестовую семью
type TestFamily struct {
	Name    string         `json:"name"`
	Members []TestMember   `json:"members"`
}

// TestMember представляет тестового члена семьи
type TestMember struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Gender    string   `json:"gender"`
	BirthDate string   `json:"birth_date"`
	Parents   []string `json:"parents,omitempty"`
}

// LoadTestFamily загружает тестовую семью из JSON файла
func LoadTestFamily(t *testing.T, path string) *TestFamily {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	var family TestFamily
	err = json.Unmarshal(data, &family)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	return &family
}

// BuildFamilySystem создает FamilySystem из тестовых данных
func BuildFamilySystem(t *testing.T, testFamily *TestFamily, extended bool) *model.FamilySystem {
	fs := model.NewFamilySystem(extended)

	// Сначала создаем всех членов
	for _, tm := range testFamily.Members {
		gender := core.GenderFemale
		if tm.Gender == "male" {
			gender = core.GenderMale
		}

		member, err := data.CreateFamilyMember(tm.ID, tm.Name, gender, tm.BirthDate)
		if err != nil {
			t.Fatalf("Failed to create member %s: %v", tm.ID, err)
		}

		fs.AddMember(member)
	}

	// Потом устанавливаем связи
	for _, tm := range testFamily.Members {
		member := fs.GetMember(tm.ID)
		if member == nil {
			continue
		}

		for _, parentID := range tm.Parents {
			member.Parents = append(member.Parents, parentID)
			if parent := fs.GetMember(parentID); parent != nil {
				parent.Children = append(parent.Children, tm.ID)
			}
		}
	}

	return fs
}

// AssertEqualCoords проверяет равенство координат
func AssertEqualCoords(t *testing.T, got, want core.HypercubeCoords) {
	if got.X != want.X || got.Y != want.Y || got.Z != want.Z || got.W != want.W {
		t.Errorf("Coords = %v, want %v", got, want)
	}
}

// AssertEqualStrings проверяет равенство строк
func AssertEqualStrings(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// AssertEqualInts проверяет равенство целых чисел
func AssertEqualInts(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

// AssertEqualFloats проверяет равенство чисел с плавающей точкой
func AssertEqualFloats(t *testing.T, got, want, epsilon float64) {
	if got < want-epsilon || got > want+epsilon {
		t.Errorf("got %f, want %f (±%f)", got, want, epsilon)
	}
}

// AssertTrue проверяет истинность
func AssertTrue(t *testing.T, condition bool, msg string) {
	if !condition {
		t.Error(msg)
	}
}

// AssertFalse проверяет ложность
func AssertFalse(t *testing.T, condition bool, msg string) {
	if condition {
		t.Error(msg)
	}
}

// AssertNoError проверяет отсутствие ошибки
func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// AssertError проверяет наличие ошибки
func AssertError(t *testing.T, err error) {
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
