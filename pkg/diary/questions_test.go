package diary

import (
	"testing"
)

func TestGetSections(t *testing.T) {
	sections := GetSections()

	if len(sections) != 5 {
		t.Errorf("expected 5 sections, got %d", len(sections))
	}

	expectedSections := []string{"motivation", "boundaries", "resources", "patterns", "choice"}
	for i, expected := range expectedSections {
		if sections[i].ID != expected {
			t.Errorf("section %d: expected %s, got %s", i, expected, sections[i].ID)
		}
	}
}

func TestGetQuestions(t *testing.T) {
	questions := GetQuestions()

	// 5 разделов × 5 вопросов = 25 вопросов
	if len(questions) != 25 {
		t.Errorf("expected 25 questions, got %d", len(questions))
	}

	// Проверяем, что у каждого вопроса есть обязательные поля
	for _, q := range questions {
		if q.ID == "" {
			t.Error("question ID is empty")
		}
		if q.Section == "" {
			t.Error("question section is empty")
		}
		if q.Text == "" {
			t.Errorf("question %s text is empty", q.ID)
		}
		if q.MinLength < 0 {
			t.Errorf("question %s has negative minLength", q.ID)
		}
		if q.MaxLength <= q.MinLength {
			t.Errorf("question %s: maxLength must be greater than minLength", q.ID)
		}
	}
}

func TestGetQuestionsBySection(t *testing.T) {
	questions := GetQuestionsBySection("motivation")

	if len(questions) != 5 {
		t.Errorf("expected 5 questions in motivation section, got %d", len(questions))
	}

	for _, q := range questions {
		if q.Section != "motivation" {
			t.Errorf("expected section motivation, got %s", q.Section)
		}
	}
}

func TestGetQuestionByID(t *testing.T) {
	question, ok := GetQuestionByID("motivation_1")
	if !ok {
		t.Fatal("expected to find question motivation_1")
	}

	if question.Text == "" {
		t.Error("expected non-empty question text")
	}

	// Проверяем несуществующий вопрос
	_, ok = GetQuestionByID("nonexistent")
	if ok {
		t.Error("expected false for nonexistent question")
	}
}

func TestGetSectionByID(t *testing.T) {
	section, ok := GetSectionByID("boundaries")
	if !ok {
		t.Fatal("expected to find section boundaries")
	}

	if section.Title == "" {
		t.Error("expected non-empty section title")
	}

	// Проверяем несуществующий раздел
	_, ok = GetSectionByID("nonexistent")
	if ok {
		t.Error("expected false for nonexistent section")
	}
}

func TestGetFinalQuestion(t *testing.T) {
	question := GetFinalQuestion()

	if question.ID != "final_choice" {
		t.Errorf("expected ID final_choice, got %s", question.ID)
	}
	if question.Section != "choice" {
		t.Errorf("expected section choice, got %s", question.Section)
	}
}

func TestGetStopSignals(t *testing.T) {
	signals := GetStopSignals()

	if len(signals) == 0 {
		t.Error("expected at least one stop signal")
	}

	for _, signal := range signals {
		if signal == "" {
			t.Error("empty stop signal found")
		}
	}
}

func TestValidateAnswer(t *testing.T) {
	tests := []struct {
		name        string
		questionID  string
		answer      string
		shouldPass  bool
	}{
		{"valid short answer", "motivation_1", "Для себя", true},
		{"too short", "motivation_1", "Да", false}, // меньше 5 символов
		{"valid long answer", "motivation_2", "Да, буду делать всё равно, потому что это моё", true},
		{"empty answer", "motivation_1", "", false},
		{"nonexistent question", "nonexistent", "answer", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := ValidateAnswer(tt.questionID, tt.answer)
			if valid != tt.shouldPass {
				t.Errorf("expected valid=%v, got %v (error: %v)", tt.shouldPass, valid, err)
			}
		})
	}
}

func TestQuestionOrder(t *testing.T) {
	sections := GetSections()

	// Проверяем, что порядок секций правильный
	for i := 1; i < len(sections); i++ {
		if sections[i].Order <= sections[i-1].Order {
			t.Errorf("sections not in correct order: %d <= %d", sections[i].Order, sections[i-1].Order)
		}
	}

	// Проверяем порядок вопросов в каждом разделе
	sectionIDs := []string{"motivation", "boundaries", "resources", "patterns", "choice"}
	for _, sectionID := range sectionIDs {
		questions := GetQuestionsBySection(sectionID)
		for i := 1; i < len(questions); i++ {
			if questions[i].Order <= questions[i-1].Order {
				t.Errorf("questions in %s not in correct order: %d <= %d", sectionID, questions[i].Order, questions[i-1].Order)
			}
		}
	}
}

func TestQuestionTags(t *testing.T) {
	questions := GetQuestions()

	for _, q := range questions {
		if len(q.Tags) == 0 {
			t.Errorf("question %s has no tags", q.ID)
		}
	}
}
