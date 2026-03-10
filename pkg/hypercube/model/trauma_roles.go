package model

import (
	"github.com/ruvcoindev/idealcore/hypercube/core"
)

// isScapegoat проверяет, является ли человек "козлом отпущения"
// Признаки: частые конфликты, критика в его адрес, обвинения во всех проблемах семьи
func (fs *FamilySystem) isScapegoat(memberID string, birthOrder int, siblings []string) bool {
	member := fs.Members[memberID]
	if member == nil {
		return false
	}
	
	// Считаем негативные события
	negativeEvents := 0
	for _, event := range member.Events {
		if event.EventType == "conflict" || event.EventType == "criticism" {
			negativeEvents++
		}
	}
	
	// Считаем, сколько раз другие упоминали этого человека в негативном контексте
	criticizedByOthers := 0
	for _, sibID := range siblings {
		if sib, ok := fs.Members[sibID]; ok {
			for _, event := range sib.Events {
				if event.EventType == "criticism" && event.WithPerson == memberID {
					criticizedByOthers++
				}
			}
		}
	}
	
	// В больших семьях козлом отпущения часто становится средний ребенок
	if len(siblings) >= 5 && birthOrder > 2 && birthOrder < len(siblings)-1 {
		return negativeEvents > 1 || criticizedByOthers > 0
	}
	
	return negativeEvents > 2 || criticizedByOthers > 1
}

// isGoldenChild проверяет, является ли человек "любимчиком"
// Признаки: идеализация, высокие ожидания, давление быть совершенным
func (fs *FamilySystem) isGoldenChild(memberID string, birthOrder int, siblings []string) bool {
	member := fs.Members[memberID]
	if member == nil {
		return false
	}
	
	// Считаем позитивные события
	positiveEvents := 0
	for _, event := range member.Events {
		if event.EventType == "praise" || event.EventType == "achievement" {
			positiveEvents++
		}
	}
	
	// Считаем идеализацию со стороны родителей
	idealizedByParents := 0
	for _, parentID := range member.Parents {
		if parent, ok := fs.Members[parentID]; ok {
			for _, event := range parent.Events {
				if event.EventType == "praise" && event.WithPerson == memberID {
					idealizedByParents++
				}
			}
		}
	}
	
	// В больших семьях любимчиками часто становятся первый или последний
	if len(siblings) >= 5 && (birthOrder == 1 || birthOrder == len(siblings)) {
		return true
	}
	
	return positiveEvents > 3 || idealizedByParents > 2
}

// isLostChild проверяет, является ли человек "потерянным ребенком"
// Признаки: эмоциональная отстраненность, мало событий, мало отношений
func (fs *FamilySystem) isLostChild(memberID string, vectors core.PersonVectors) bool {
	member := fs.Members[memberID]
	if member == nil {
		return false
	}
	
	// Вычисляем амплитуду векторов (эмоциональную активность)
	amp := core.VectorAmplitude(vectors)
	
	// Считаем события и партнеров
	eventCount := len(member.Events)
	partnerCount := len(member.Partners)
	
	// В очень большой семье младшие дети часто "теряются"
	siblings := fs.GetSiblings(memberID)
	if len(siblings) >= 8 {
		birthOrder := fs.GetBirthOrder(memberID)
		if birthOrder > len(siblings)-2 {
			return true
		}
	}
	
	return amp < 3.0 && eventCount < 3 && partnerCount < 2
}

// isInvisible проверяет, является ли человек "невидимкой"
// Признаки: о нем не говорят, его события не запоминают, он на периферии
func (fs *FamilySystem) isInvisible(memberID string, siblings []string) bool {
	member := fs.Members[memberID]
	if member == nil {
		return false
	}
	
	// Считаем, сколько раз другие упоминали этого человека
	mentionedByOthers := 0
	for _, sibID := range siblings {
		if sib, ok := fs.Members[sibID]; ok {
			for _, event := range sib.Events {
				if event.WithPerson == memberID {
					mentionedByOthers++
				}
			}
		}
	}
	
	// В большой семье "средние" дети часто становятся невидимками
	if len(siblings) >= 5 {
		birthOrder := fs.GetBirthOrder(memberID)
		if birthOrder > 2 && birthOrder < len(siblings)-1 {
			return mentionedByOthers == 0
		}
	}
	
	return mentionedByOthers == 0 && len(member.Events) < 2
}

// isGlassChild проверяет, является ли человек "стеклянным ребенком"
// Признаки: есть проблемный сиблинг, на фоне которого он становится "прозрачным"
func (fs *FamilySystem) isGlassChild(memberID string, siblings []string) bool {
	// Проверяем, есть ли сиблинг с серьезными проблемами
	for _, sibID := range siblings {
		if sib, ok := fs.Members[sibID]; ok {
			problematicEvents := 0
			for _, event := range sib.Events {
				if event.EventType == "hospitalization" ||
					event.EventType == "addiction" ||
					event.EventType == "crisis" ||
					event.EventType == "disability" {
					problematicEvents++
				}
			}
			if problematicEvents > 2 {
				return true
			}
		}
	}
	
	return false
}

// isParentified проверяет, является ли человек "родифицированным ребенком"
// Признаки: заботится о родителях/сиблингах, рано повзрослел
func (fs *FamilySystem) isParentified(memberID string, birthOrder int) bool {
	member := fs.Members[memberID]
	if member == nil {
		return false
	}
	
	// Считаем события, связанные с заботой о других
	careEvents := 0
	for _, event := range member.Events {
		if event.EventType == "care" || event.EventType == "responsibility" || event.EventType == "provider" {
			careEvents++
		}
	}
	
	// В большой семье старшие дети всегда родифицированы
	siblings := fs.GetSiblings(memberID)
	if len(siblings) >= 5 && birthOrder == 1 {
		return true
	}
	
	return birthOrder == 1 && careEvents > 2
}

// isMascot проверяет, является ли человек "шутом"
// Признаки: использует юмор для снятия напряжения, минимизирует проблемы
func (fs *FamilySystem) isMascot(memberID string, vectors core.PersonVectors) bool {
	member := fs.Members[memberID]
	if member == nil {
		return false
	}
	
	amp := core.VectorAmplitude(vectors)
	
	// Считаем события, связанные с юмором и развлечениями
	humorEvents := 0
	for _, event := range member.Events {
		if event.EventType == "humor" || event.EventType == "entertainment" || event.EventType == "joke" {
			humorEvents++
		}
	}
	
	return amp > 6.0 && humorEvents > 2
}

// isEmotionalSpouse проверяет, является ли человек "эмоциональным супругом"
// Признаки: особо близкие отношения с родителем противоположного пола,
// замещает роль партнера для этого родителя
func (fs *FamilySystem) isEmotionalSpouse(memberID string) bool {
	member := fs.Members[memberID]
	if member == nil {
		return false
	}
	
	// Проверяем близость с родителем противоположного пола
	closeWithParent := false
	for _, parentID := range member.Parents {
		if parent, ok := fs.Members[parentID]; ok {
			// Проверяем, противоположного ли пола родитель
			if parent.Gender != member.Gender {
				// Считаем интимные события (доверительные разговоры, совместное время)
				intimateCount := 0
				for _, event := range parent.Events {
					if event.WithPerson == memberID &&
						(event.EventType == "intimate" || event.EventType == "confidant" || event.EventType == "date") {
						intimateCount++
					}
				}
				if intimateCount > 3 {
					closeWithParent = true
				}
			}
		}
	}
	
	return closeWithParent && len(member.Partners) < 2
}

// isReplacementChild проверяет, является ли человек замещающим ребенком
// Просто проверяет флаг, установленный при вводе данных
func (fs *FamilySystem) isReplacementChild(memberID string) bool {
	member := fs.Members[memberID]
	return member != nil && member.IsReplacement
}

// IdentifyTraumaRoles выявляет травматические роли для человека
func (fs *FamilySystem) IdentifyTraumaRoles(memberID string) *IdentifiedTraumaRoles {
	member := fs.Members[memberID]
	if member == nil {
		return nil
	}

	roles := &IdentifiedTraumaRoles{
		SecondaryRoles: make([]core.TraumaRole, 0),
		Evidence:       make([]string, 0),
	}

	// Получаем данные для анализа
	coords := core.ParseDateToCoords(member.BirthDate)
	vectors := core.CalculateVectors(coords)
	siblings := fs.GetSiblings(memberID)
	birthOrder := fs.GetBirthOrder(memberID)

	// Специальная проверка: замещающий ребенок (приоритет)
	if member.IsReplacement {
		roles.PrimaryRole = core.TraumaRoleReplacement
		roles.Confidence = 0.95
		roles.Evidence = append(roles.Evidence,
			"Ребенок рожден после потери предыдущего",
			"Живет с чувством, что должен заменить умершего/потерянного",
			"Может испытывать вину выжившего и диффузию идентичности")
		
		if member.ReplacementData != nil {
			if member.ReplacementData.TimeGapMonths < 6 {
				roles.Evidence = append(roles.Evidence, "Очень короткий разрыв между потерей и зачатием - экстренное замещение")
			}
			if member.ReplacementData.IsExplicit {
				roles.Evidence = append(roles.Evidence, "Явное замещение (то же имя, те же ожидания)")
			}
		}
		
		fs.TraumaRoles[memberID] = roles
		return roles
	}

	// Проверяем остальные роли по очереди
	if fs.isScapegoat(memberID, birthOrder, siblings) {
		roles.PrimaryRole = core.TraumaRoleScapegoat
		roles.Confidence = 0.8
		roles.Evidence = append(roles.Evidence,
			"Частая критика и обвинения со стороны семьи",
			"На него проецируются все семейные проблемы")
	} else if fs.isGoldenChild(memberID, birthOrder, siblings) {
		roles.PrimaryRole = core.TraumaRoleGoldenChild
		roles.Confidence = 0.8
		roles.Evidence = append(roles.Evidence,
			"Идеализация родителями",
			"Высокие ожидания и давление быть совершенным")
	} else if fs.isLostChild(memberID, vectors) {
		roles.PrimaryRole = core.TraumaRoleLostChild
		roles.Confidence = 0.7
		roles.Evidence = append(roles.Evidence,
			"Эмоциональная отстраненность",
			"Склонность к уединению и незаметности")
	} else if fs.isInvisible(memberID, siblings) {
		roles.PrimaryRole = core.TraumaRoleInvisible
		roles.Confidence = 0.6
		roles.Evidence = append(roles.Evidence,
			"Игнорирование потребностей",
			"Существует на периферии семейного внимания")
	} else if fs.isGlassChild(memberID, siblings) {
		roles.PrimaryRole = core.TraumaRoleGlassChild
		roles.Confidence = 0.7
		roles.Evidence = append(roles.Evidence,
			"Есть проблемный сиблинг, на фоне которого стал невидимым",
			"Сверх-самостоятельность, скрытые проблемы")
	} else if fs.isParentified(memberID, birthOrder) {
		roles.PrimaryRole = core.TraumaRoleParentified
		roles.Confidence = 0.8
		roles.Evidence = append(roles.Evidence,
			"Забота о родителях или младших сиблингах",
			"Отсутствие детства, гиперответственность")
	} else if fs.isMascot(memberID, vectors) {
		roles.PrimaryRole = core.TraumaRoleMascot
		roles.Confidence = 0.6
		roles.Evidence = append(roles.Evidence,
			"Снятие напряжения юмором",
			"Минимизация серьезных проблем")
	} else if fs.isEmotionalSpouse(memberID) {
		roles.PrimaryRole = core.TraumaRoleEmotionalSpouse
		roles.Confidence = 0.9
		roles.Evidence = append(roles.Evidence,
			"Замещение супружеской роли для родителя",
			"Эмоциональная близость, заменяющая партнерские отношения")
	} else {
		// Если роль не определена, ставим "тень"
		roles.PrimaryRole = core.TraumaRoleShadow
		roles.Confidence = 0.5
		roles.Evidence = append(roles.Evidence,
			"Невыраженная травма, требуется дополнительный анализ",
			"Возможно, несет подавленные эмоции семьи")
	}

	// Проверка на феномен "девочка с яйцами"
	// Женщина в маскулинной роли, с высокой амплитудой векторов
	if member.Gender == core.GenderFemale &&
		(roles.PrimaryRole == core.TraumaRoleLostChild || roles.PrimaryRole == core.TraumaRoleScapegoat) &&
		core.VectorAmplitude(vectors) > 5.0 {
		roles.Evidence = append(roles.Evidence,
			"Феномен 'девочка с яйцами' - маскулинные стратегии выживания",
			"Компенсаторная гипер-независимость, отказ от феминности как защита")
	}

	fs.TraumaRoles[memberID] = roles
	return roles
}
