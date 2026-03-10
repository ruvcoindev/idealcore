// Package diary предоставляет опросник "Кто я" для саморефлексии
package diary

import "fmt"

// Section — раздел опросника
type Section struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

// Question — вопрос опросника
type Question struct {
	ID          string   `json:"id"`
	Section     string   `json:"section"`
	Text        string   `json:"text"`
	Placeholder string   `json:"placeholder"`
	MinLength   int      `json:"min_length"`
	MaxLength   int      `json:"max_length"`
	Tags        []string `json:"tags"`
	Order       int      `json:"order"`
}

// GetSections возвращает все разделы опросника
func GetSections() []Section {
	return []Section{
		{
			ID:          "motivation",
			Title:       "Мотивация",
			Description: "Зачем я это делаю? Ради кого или чего?",
			Order:       1,
		},
		{
			ID:          "boundaries",
			Title:       "Границы",
			Description: "Где заканчиваюсь «я» и начинается «другой»?",
			Order:       2,
		},
		{
			ID:          "resources",
			Title:       "Ресурс",
			Description: "Работа, деньги, энергия — что я вкладываю и что получаю?",
			Order:       3,
		},
		{
			ID:          "patterns",
			Title:       "Паттерны",
			Description: "Узнаю ли я это? Я уже проходил этот сценарий?",
			Order:       4,
		},
		{
			ID:          "choice",
			Title:       "Выбор",
			Description: "Что я выбираю сегодня — ради себя?",
			Order:       5,
		},
	}
}

// GetQuestions возвращает все вопросы опросника
func GetQuestions() []Question {
	return []Question{
		// === РАЗДЕЛ 1: МОТИВАЦИЯ ===
		{
			ID:          "motivation_1",
			Section:     "motivation",
			Text:        "Ради кого я это делаю? (конкретное имя или «для себя/работы»)",
			Placeholder: "Напиши имя или «для себя»",
			MinLength:   5,
			MaxLength:   500,
			Tags:        []string{"motivation", "purpose"},
			Order:       1,
		},
		{
			ID:          "motivation_2",
			Section:     "motivation",
			Text:        "Если этот человек исчезнет из моей жизни завтра — буду ли я делать это всё равно?",
			Placeholder: "Честно: да или нет?",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"motivation", "dependency"},
			Order:       2,
		},
		{
			ID:          "motivation_3",
			Section:     "motivation",
			Text:        "Я жду, что это изменит отношение Х ко мне?",
			Placeholder: "Да/нет + пояснение",
			MinLength:   5,
			MaxLength:   500,
			Tags:        []string{"motivation", "expectations"},
			Order:       3,
		},
		{
			ID:          "motivation_4",
			Section:     "motivation",
			Text:        "Что я чувствую, когда думаю об этом: спокойствие или напряжение?",
			Placeholder: "Опиши телесное ощущение",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"motivation", "body"},
			Order:       4,
		},
		{
			ID:          "motivation_5",
			Section:     "motivation",
			Text:        "Это выбор или надежда?",
			Placeholder: "Выбираю я или надеюсь, что они изменятся?",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"motivation", "choice"},
			Order:       5,
		},

		// === РАЗДЕЛ 2: ГРАНИЦЫ ===
		{
			ID:          "boundaries_1",
			Section:     "boundaries",
			Text:        "Где заканчиваюсь «я» и начинается «другой» в этой ситуации?",
			Placeholder: "Что моё, а что не моё?",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"boundaries", "self"},
			Order:       1,
		},
		{
			ID:          "boundaries_2",
			Section:     "boundaries",
			Text:        "Я беру на себя ответственность за чужие чувства/выбор/последствия? (что именно)",
			Placeholder: "Перечисли, за что ты отвечаешь вместо другого",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"boundaries", "responsibility"},
			Order:       2,
		},
		{
			ID:          "boundaries_3",
			Section:     "boundaries",
			Text:        "Если я скажу «нет» — что я потеряю? А если скажу «да» — что я потеряю?",
			Placeholder: "Сравни цену обоих решений",
			MinLength:   20,
			MaxLength:   1500,
			Tags:        []string{"boundaries", "choice"},
			Order:       3,
		},
		{
			ID:          "boundaries_4",
			Section:     "boundaries",
			Text:        "Это моё дело или я лезу не в свою зону контроля?",
			Placeholder: "Честно: можешь ли ты это контролировать?",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"boundaries", "control"},
			Order:       4,
		},
		{
			ID:          "boundaries_5",
			Section:     "boundaries",
			Text:        "Что я разрешаю делать со мной, что не разрешил бы другу?",
			Placeholder: "Сравни отношение к себе и к близкому",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"boundaries", "self-compassion"},
			Order:       5,
		},

		// === РАЗДЕЛ 3: РЕСУРС ===
		{
			ID:          "resources_1",
			Section:     "resources",
			Text:        "Я работаю ради потока или ради доказательства?",
			Placeholder: "Деньги или подтверждение ценности?",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"resources", "work"},
			Order:       1,
		},
		{
			ID:          "resources_2",
			Section:     "resources",
			Text:        "Если сегодня не будет отдачи — буду ли я делать это всё равно?",
			Placeholder: "Внутренняя мотивация или внешняя?",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"resources", "motivation"},
			Order:       2,
		},
		{
			ID:          "resources_3",
			Section:     "resources",
			Text:        "Я вкладываю в проект/человека больше, чем получаю обратно? (время, эмоции, деньги)",
			Placeholder: "Подсчитай баланс",
			MinLength:   20,
			MaxLength:   1500,
			Tags:        []string{"resources", "balance"},
			Order:       3,
		},
		{
			ID:          "resources_4",
			Section:     "resources",
			Text:        "Что я отложил «на потом» ради этого? (здоровье, отдых, другие клиенты)",
			Placeholder: "Какая цена в твоей жизни?",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"resources", "sacrifice"},
			Order:       4,
		},
		{
			ID:          "resources_5",
			Section:     "resources",
			Text:        "Это инвестиция или ставка на один билет?",
			Placeholder: "Диверсификация или все яйца в одну корзину?",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"resources", "risk"},
			Order:       5,
		},

		// === РАЗДЕЛ 4: ПАТТЕРНЫ ===
		{
			ID:          "patterns_1",
			Section:     "patterns",
			Text:        "Я уже проходил этот сценарий? Когда? Чем закончилось?",
			Placeholder: "Вспомни похожие ситуации из прошлого",
			MinLength:   20,
			MaxLength:   2000,
			Tags:        []string{"patterns", "history"},
			Order:       1,
		},
		{
			ID:          "patterns_2",
			Section:     "patterns",
			Text:        "Я пытаюсь «переиграть» прошлое через этого человека/ситуацию?",
			Placeholder: "Что ты хочешь исправить?",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"patterns", "repetition"},
			Order:       2,
		},
		{
			ID:          "patterns_3",
			Section:     "patterns",
			Text:        "Я жду, что в этот раз будет иначе, хотя условия те же?",
			Placeholder: "Честно: что изменилось?",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"patterns", "expectations"},
			Order:       3,
		},
		{
			ID:          "patterns_4",
			Section:     "patterns",
			Text:        "Что я игнорирую, потому что «очень хочу»?",
			Placeholder: "Красные флаги, которые ты не замечаешь",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"patterns", "denial"},
			Order:       4,
		},
		{
			ID:          "patterns_5",
			Section:     "patterns",
			Text:        "Если бы мой лучший друг рассказал мне это — что бы я ему сказал?",
			Placeholder: "Посмотри со стороны",
			MinLength:   20,
			MaxLength:   1500,
			Tags:        []string{"patterns", "perspective"},
			Order:       5,
		},

		// === РАЗДЕЛ 5: ВЫБОР ===
		{
			ID:          "choice_1",
			Section:     "choice",
			Text:        "Я действую из страха или из интереса?",
			Placeholder: "Что движет тобой сейчас?",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"choice", "motivation"},
			Order:       1,
		},
		{
			ID:          "choice_2",
			Section:     "choice",
			Text:        "Что я выберу, если отброшу надежду, что «они поймут»?",
			Placeholder: "Без надежды на их изменение",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"choice", "acceptance"},
			Order:       2,
		},
		{
			ID:          "choice_3",
			Section:     "choice",
			Text:        "Это шаг вперёд или шаг в сторону от себя?",
			Placeholder: "Приближает или отдаляет от целостности?",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"choice", "alignment"},
			Order:       3,
		},
		{
			ID:          "choice_4",
			Section:     "choice",
			Text:        "Что я могу сделать сегодня, что зависит только от меня?",
			Placeholder: "Конкретное действие в твоей зоне контроля",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"choice", "action"},
			Order:       4,
		},
		{
			ID:          "choice_5",
			Section:     "choice",
			Text:        "Если я не получу подтверждения, что «всё правильно» — смогу ли я идти дальше?",
			Placeholder: "Твоя внутренняя опора",
			MinLength:   10,
			MaxLength:   1000,
			Tags:        []string{"choice", "self-trust"},
			Order:       5,
		},
	}
}

// GetQuestionsBySection возвращает вопросы конкретного раздела
func GetQuestionsBySection(sectionID string) []Question {
	all := GetQuestions()
	var result []Question
	for _, q := range all {
		if q.Section == sectionID {
			result = append(result, q)
		}
	}
	return result
}

// GetQuestionByID возвращает вопрос по ID
func GetQuestionByID(id string) (Question, bool) {
	all := GetQuestions()
	for _, q := range all {
		if q.ID == id {
			return q, true
		}
	}
	return Question{}, false
}

// GetSectionByID возвращает раздел по ID
func GetSectionByID(id string) (Section, bool) {
	sections := GetSections()
	for _, s := range sections {
		if s.ID == id {
			return s, true
		}
	}
	return Section{}, false
}

// GetFinalQuestion возвращает финальный вопрос опросника
func GetFinalQuestion() Question {
	return Question{
		ID:          "final_choice",
		Section:     "choice",
		Text:        "Что я выбираю сегодня — ради себя?",
		Placeholder: "Не «надо». Не «правильно». Не «ради будущего». А сейчас.",
		MinLength:   10,
		MaxLength:   1000,
		Tags:        []string{"choice", "final"},
		Order:       6,
	}
}

// GetStopSignals возвращает стоп-сигналы для самопроверки
func GetStopSignals() []string {
	return []string{
		"Я проверяю статус/онлайн человека, которому обещал «не писать»",
		"Я отменяю свои планы, потому что «вдруг она напишет»",
		"Я объясняю одно и то же человеку, который уже сказал «нет»",
		"Я чувствую, что «должен» доказать, что я хороший/правильный/достойный",
		"Я жертвую работой/здоровьем/деньгами ради «шанса» с одним человеком",
		"Я жду, что человек «очнётся», «поймёт», «оценит» — без моих слов",
		"Я чувствую обиду, что «я всё дал, а она не взяла»",
	}
}

// ValidateAnswer проверяет ответ на соответствие требованиям вопроса
func ValidateAnswer(questionID, answer string) (bool, error) {
	question, ok := GetQuestionByID(questionID)
	if !ok {
		return false, fmt.Errorf("question not found: %s", questionID)
	}

	if len(answer) < question.MinLength {
		return false, fmt.Errorf("answer too short: minimum %d characters", question.MinLength)
	}

	if len(answer) > question.MaxLength {
		return false, fmt.Errorf("answer too long: maximum %d characters", question.MaxLength)
	}

	return true, nil
}
