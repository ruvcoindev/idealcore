package model

import (
	"fmt"
	"sort"
	"github.com/ruvcoindev/idealcore/hypercube/core"
)

// FindFamilyPatterns выявляет все семейные паттерны
func (fs *FamilySystem) FindFamilyPatterns() []FamilyPattern {
	patterns := make([]FamilyPattern, 0)

	// Паттерны родитель-ребенок
	if p := fs.findMotherDaughterPattern(); p != nil {
		patterns = append(patterns, *p)
	}
	if p := fs.findMotherSonPattern(); p != nil {
		patterns = append(patterns, *p)
	}
	if p := fs.findFatherDaughterPattern(); p != nil {
		patterns = append(patterns, *p)
	}
	if p := fs.findFatherSonPattern(); p != nil {
		patterns = append(patterns, *p)
	}

	// Паттерны сиблингов
	if p := fs.findSisterSisterPattern(); p != nil {
		patterns = append(patterns, *p)
	}
	if p := fs.findBrotherBrotherPattern(); p != nil {
		patterns = append(patterns, *p)
	}
	if p := fs.findSisterBrotherPattern(); p != nil {
		patterns = append(patterns, *p)
	}

	// Межпоколенческие паттерны
	if p := fs.findGrandmotherGranddaughterPattern(); p != nil {
		patterns = append(patterns, *p)
	}
	if p := fs.findGrandfatherGrandsonPattern(); p != nil {
		patterns = append(patterns, *p)
	}
	if p := fs.findGrandmotherGrandsonPattern(); p != nil {
		patterns = append(patterns, *p)
	}
	if p := fs.findGrandfatherGranddaughterPattern(); p != nil {
		patterns = append(patterns, *p)
	}

	// Паттерны потерь и замещений
	if p := fs.findReplacementChildPattern(); p != nil {
		patterns = append(patterns, *p)
	}

	// Паттерны повторения (разводы, измены)
	if p := fs.findDivorcePattern(); p != nil {
		patterns = append(patterns, *p)
	}
	if p := fs.findInfidelityPattern(); p != nil {
		patterns = append(patterns, *p)
	}

	// Паттерны больших семей
	if p := fs.findLargeFamilyPattern(); p != nil {
		patterns = append(patterns, *p)
	}

	fs.Patterns = patterns
	return patterns
}

// findMotherDaughterPattern - передача травмы по женской линии
func (fs *FamilySystem) findMotherDaughterPattern() *FamilyPattern {
	for id, member := range fs.Members {
		if member.Gender != core.GenderFemale || member.Role != core.RoleMother {
			continue
		}
		
		for _, childID := range member.Children {
			child := fs.Members[childID]
			if child == nil || child.Gender != core.GenderFemale {
				continue
			}
			
			motherRole := fs.TraumaRoles[id]
			childRole := fs.TraumaRoles[childID]
			
			if motherRole != nil && childRole != nil {
				// Передача роли козла отпущения
				if motherRole.PrimaryRole == core.TraumaRoleScapegoat &&
					childRole.PrimaryRole == core.TraumaRoleScapegoat {
					return &FamilyPattern{
						PatternType: "mother-daughter-scapegoat",
						Description: "Передача роли 'козла отпущения' по женской линии",
						Members:     []string{id, childID},
						Severity:    0.8,
						Generations: []int{member.Generation, child.Generation},
						Recommendations: []string{
							"Терапия женской линии",
							"Работа с родовыми сценариями",
							"Сепарация от материнских проекций",
							"Восстановление собственной идентичности",
						},
					}
				}
				
				// Передача роли любимчика
				if motherRole.PrimaryRole == core.TraumaRoleGoldenChild &&
					childRole.PrimaryRole == core.TraumaRoleGoldenChild {
					return &FamilyPattern{
						PatternType: "mother-daughter-golden",
						Description: "Передача роли 'любимчика' по женской линии - перфекционизм и высокие ожидания",
						Members:     []string{id, childID},
						Severity:    0.7,
						Generations: []int{member.Generation, child.Generation},
						Recommendations: []string{
							"Работа с перфекционизмом",
							"Принятие несовершенства",
							"Отделение от родительских ожиданий",
						},
					}
				}
			}
		}
	}
	return nil
}

// findMotherSonPattern - эмоциональный инцест
func (fs *FamilySystem) findMotherSonPattern() *FamilyPattern {
	for id, member := range fs.Members {
		if member.Gender != core.GenderFemale || member.Role != core.RoleMother {
			continue
		}
		
		for _, childID := range member.Children {
			child := fs.Members[childID]
			if child == nil || child.Gender != core.GenderMale {
				continue
			}
			
			childRole := fs.TraumaRoles[childID]
			if childRole != nil && childRole.PrimaryRole == core.TraumaRoleEmotionalSpouse {
				return &FamilyPattern{
					PatternType: "mother-son-emotional-incest",
					Description: "Эмоциональный инцест: сын замещает роль мужа для матери",
					Members:     []string{id, childID},
					Severity:    0.95,
					Generations: []int{member.Generation, child.Generation},
					Recommendations: []string{
						"Срочная терапия для разрыва симбиотической связи",
						"Работа с мужской идентичностью сына",
						"Отделение от материнских ожиданий",
						"Терапия для матери - восполнение супружеской роли",
						"Формирование здоровых границ",
					},
				}
			}
			
			// Менее выраженная форма - сын-любимчик
			if childRole != nil && childRole.PrimaryRole == core.TraumaRoleGoldenChild {
				// Проверяем, нет ли у матери мужа или есть проблемы в браке
				motherHasPartner := member.CurrentPartner != ""
				if !motherHasPartner {
					return &FamilyPattern{
						PatternType: "mother-son-golden",
						Description: "Сын-любимчик - компенсация отсутствия удовлетворительных отношений с мужем",
						Members:     []string{id, childID},
						Severity:    0.7,
						Generations: []int{member.Generation, child.Generation},
						Recommendations: []string{
							"Осознание роли сына как эмоциональной поддержки",
							"Поиск здоровых источников эмоциональной близости",
							"Терапия для формирования здоровой сепарации",
						},
					}
				}
			}
		}
	}
	return nil
}

// findFatherDaughterPattern - комплекс Электры
func (fs *FamilySystem) findFatherDaughterPattern() *FamilyPattern {
	for id, member := range fs.Members {
		if member.Gender != core.GenderMale || member.Role != core.RoleFather {
			continue
		}
		
		for _, childID := range member.Children {
			child := fs.Members[childID]
			if child == nil || child.Gender != core.GenderFemale {
				continue
			}
			
			childRole := fs.TraumaRoles[childID]
			if childRole != nil && childRole.PrimaryRole == core.TraumaRoleGoldenChild {
				return &FamilyPattern{
					PatternType: "father-daughter-golden",
					Description: "Дочь-любимица отца - комплекс Электры, влияние на выбор партнеров",
					Members:     []string{id, childID},
					Severity:    0.7,
					Generations: []int{member.Generation, child.Generation},
					Recommendations: []string{
						"Работа с образом идеального мужчины",
						"Осознание проекций отца на партнеров",
						"Формирование здоровых отношений с мужчинами",
						"Отделение от отцовских ожиданий",
					},
				}
			}
			
			// Отец-дочь с эмоциональной близостью (риск инцеста)
			if childRole != nil && childRole.PrimaryRole == core.TraumaRoleEmotionalSpouse {
				return &FamilyPattern{
					PatternType: "father-daughter-emotional",
					Description: "Дочь как эмоциональная замена жене - риск эмоционального инцеста",
					Members:     []string{id, childID},
					Severity:    0.9,
					Generations: []int{member.Generation, child.Generation},
					Recommendations: []string{
						"Срочная терапия границ",
						"Работа с супружеской подсистемой",
						"Терапия для дочери - отделение от роли 'маленькой жены'",
					},
				}
			}
		}
	}
	return nil
}

// findFatherSonPattern - передача по мужской линии
func (fs *FamilySystem) findFatherSonPattern() *FamilyPattern {
	for id, member := range fs.Members {
		if member.Gender != core.GenderMale || member.Role != core.RoleFather {
			continue
		}
		
		for _, childID := range member.Children {
			child := fs.Members[childID]
			if child == nil || child.Gender != core.GenderMale {
				continue
			}
			
			fatherRole := fs.TraumaRoles[id]
			childRole := fs.TraumaRoles[childID]
			
			if fatherRole != nil && childRole != nil {
				if fatherRole.PrimaryRole == core.TraumaRoleScapegoat &&
					childRole.PrimaryRole == core.TraumaRoleScapegoat {
					return &FamilyPattern{
						PatternType: "father-son-scapegoat",
						Description: "Передача роли 'козла отпущения' по мужской линии",
						Members:     []string{id, childID},
						Severity:    0.8,
						Generations: []int{member.Generation, child.Generation},
						Recommendations: []string{
							"Терапия мужской линии",
							"Работа с родовыми проклятиями",
							"Восстановление мужской идентичности",
							"Преодоление чувства неполноценности",
						},
					}
				}
				
				if fatherRole.PrimaryRole == core.TraumaRoleGoldenChild &&
					childRole.PrimaryRole == core.TraumaRoleGoldenChild {
					return &FamilyPattern{
						PatternType: "father-son-golden",
						Description: "Передача роли 'любимчика' по мужской линии - высокие ожидания, перфекционизм",
						Members:     []string{id, childID},
						Severity:    0.7,
						Generations: []int{member.Generation, child.Generation},
						Recommendations: []string{
							"Работа с перфекционизмом",
							"Принятие права на ошибку",
							"Отделение от отцовских ожиданий",
						},
					}
				}
				
				if fatherRole.PrimaryRole == core.TraumaRoleLostChild &&
					childRole.PrimaryRole == core.TraumaRoleLostChild {
					return &FamilyPattern{
						PatternType: "father-son-lost",
						Description: "Передача паттерна 'потерянного ребенка' по мужской линии",
						Members:     []string{id, childID},
						Severity:    0.7,
						Generations: []int{member.Generation, child.Generation},
						Recommendations: []string{
							"Активация эмоциональной жизни",
							"Поиск своего места в семье и мире",
							"Терапия привязанности",
						},
					}
				}
			}
		}
	}
	return nil
}

// findSisterSisterPattern - динамика сестер
func (fs *FamilySystem) findSisterSisterPattern() *FamilyPattern {
	// Группируем сестер по родителям
	sisterGroups := make(map[string][]string)
	
	for id, member := range fs.Members {
		if member.Gender != core.GenderFemale {
			continue
		}
		
		// Ключ - набор родителей
		key := fmt.Sprintf("%v", member.Parents)
		sisterGroups[key] = append(sisterGroups[key], id)
	}
	
	// Анализируем каждую группу сестер
	for _, group := range sisterGroups {
		if len(group) >= 2 {
			// Получаем роли
			roles := make([]core.TraumaRole, 0)
			for _, sid := range group {
				if role, ok := fs.TraumaRoles[sid]; ok {
					roles = append(roles, role.PrimaryRole)
				}
			}
			
			// Проверяем классическую динамику: любимица + козел отпущения
			hasGolden := false
			hasScapegoat := false
			for _, r := range roles {
				if r == core.TraumaRoleGoldenChild {
					hasGolden = true
				}
				if r == core.TraumaRoleScapegoat {
					hasScapegoat = true
				}
			}
			
			if hasGolden && hasScapegoat {
				// Найдем конкретных участников
				goldenID := ""
				scapegoatID := ""
				for _, sid := range group {
					if role, ok := fs.TraumaRoles[sid]; ok {
						if role.PrimaryRole == core.TraumaRoleGoldenChild {
							goldenID = sid
						}
						if role.PrimaryRole == core.TraumaRoleScapegoat {
							scapegoatID = sid
						}
					}
				}
				
				if goldenID != "" && scapegoatID != "" {
					m1 := fs.Members[goldenID]
					m2 := fs.Members[scapegoatID]
					
					return &FamilyPattern{
						PatternType: "sister-sister-golden-scapegoat",
						Description: fmt.Sprintf("Классическая динамика сестер: %s (любимица) и %s (козел отпущения)",
							m1.Name, m2.Name),
						Members:     []string{goldenID, scapegoatID},
						Severity:    0.85,
						Generations: []int{fs.Members[group[0]].Generation},
						Recommendations: []string{
							"Терапия сиблинговых отношений",
							"Работа с завистью и соперничеством",
							"Восстановление сестринской связи",
							"Осознание роли, навязанной родителями",
						},
					}
				}
			}
			
			// Проверка на "потерянную" и "невидимку"
			hasLost := false
			hasInvisible := false
			for _, r := range roles {
				if r == core.TraumaRoleLostChild {
					hasLost = true
				}
				if r == core.TraumaRoleInvisible {
					hasInvisible = true
				}
			}
			
			if hasLost && hasInvisible {
				return &FamilyPattern{
					PatternType: "sister-sister-lost-invisible",
					Description: "Одна сестра 'потерянная', другая 'невидимка' - обе не получают внимания",
					Members:     group,
					Severity:    0.7,
					Generations: []int{fs.Members[group[0]].Generation},
					Recommendations: []string{
						"Активация обеих сестер",
						"Поиск голоса и места в семье",
						"Групповая терапия для сиблингов",
					},
				}
			}
		}
	}
	
	return nil
}

// findBrotherBrotherPattern - динамика братьев
func (fs *FamilySystem) findBrotherBrotherPattern() *FamilyPattern {
	// Группируем братьев по родителям
	brotherGroups := make(map[string][]string)
	
	for id, member := range fs.Members {
		if member.Gender != core.GenderMale {
			continue
		}
		
		key := fmt.Sprintf("%v", member.Parents)
		brotherGroups[key] = append(brotherGroups[key], id)
	}
	
	for _, group := range brotherGroups {
		if len(group) >= 2 {
			roles := make([]core.TraumaRole, 0)
			for _, bid := range group {
				if role, ok := fs.TraumaRoles[bid]; ok {
					roles = append(roles, role.PrimaryRole)
				}
			}
			
			// Треугольник братьев: любимчик, козел отпущения, потерянный
			hasGolden := false
			hasScapegoat := false
			hasLost := false
			
			for _, r := range roles {
				if r == core.TraumaRoleGoldenChild {
					hasGolden = true
				}
				if r == core.TraumaRoleScapegoat {
					hasScapegoat = true
				}
				if r == core.TraumaRoleLostChild {
					hasLost = true
				}
			}
			
			if hasGolden && hasScapegoat && hasLost {
				return &FamilyPattern{
					PatternType: "brother-brother-triangle",
					Description: "Треугольник братьев: любимчик, козел отпущения и потерянный - классическая мужская динамика",
					Members:     group,
					Severity:    0.9,
					Generations: []int{fs.Members[group[0]].Generation},
					Recommendations: []string{
						"Мужская групповая терапия",
						"Работа с конкуренцией и иерархией",
						"Восстановление братской связи",
						"Осознание родительских проекций",
					},
				}
			}
			
			// Конкуренция между старшим и младшим
			if len(group) == 2 {
				b1 := fs.Members[group[0]]
				b2 := fs.Members[group[1]]
				order1 := fs.GetBirthOrder(group[0])
				order2 := fs.GetBirthOrder(group[1])
				
				if order1 < order2 {
					// Старший и младший
					return &FamilyPattern{
						PatternType: "brother-brother-competition",
						Description: fmt.Sprintf("Конкуренция между старшим (%s) и младшим (%s) братом",
							b1.Name, b2.Name),
						Members:     group,
						Severity:    0.6,
						Generations: []int{fs.Members[group[0]].Generation},
						Recommendations: []string{
							"Осознание роли порядка рождения",
							"Работа с завистью и соперничеством",
							"Поиск индивидуального пути каждого",
						},
					}
				}
			}
		}
	}
	
	return nil
}

// findSisterBrotherPattern - динамика сестра-брат
func (fs *FamilySystem) findSisterBrotherPattern() *FamilyPattern {
	// Ищем пары сестра-брат с общими родителями
	pairs := make([][2]string, 0)
	
	for sid, sister := range fs.Members {
		if sister.Gender != core.GenderFemale {
			continue
		}
		
		for bid, brother := range fs.Members {
			if brother.Gender != core.GenderMale || sid == bid {
				continue
			}
			
			// Проверяем общих родителей
			commonParents := false
			for _, p1 := range sister.Parents {
				for _, p2 := range brother.Parents {
					if p1 == p2 {
						commonParents = true
						break
					}
				}
			}
			
			if commonParents {
				pairs = append(pairs, [2]string{sid, bid})
			}
		}
	}
	
	for _, pair := range pairs {
		sisterRole := fs.TraumaRoles[pair[0]]
		brotherRole := fs.TraumaRoles[pair[1]]
		
		if sisterRole != nil && brotherRole != nil {
			// Сестра-родитель и брат-потерянный
			if sisterRole.PrimaryRole == core.TraumaRoleParentified &&
				brotherRole.PrimaryRole == core.TraumaRoleLostChild {
				sister := fs.Members[pair[0]]
				brother := fs.Members[pair[1]]
				
				return &FamilyPattern{
					PatternType: "sister-brother-parentified-lost",
					Description: fmt.Sprintf("Сестра (%s) в роли родителя для потерянного брата (%s)",
						sister.Name, brother.Name),
					Members:     []string{pair[0], pair[1]},
					Severity:    0.75,
					Generations: []int{sister.Generation},
					Recommendations: []string{
						"Освобождение сестры от родительской роли",
						"Активация брата, поиск его голоса",
						"Баланс ответственности в паре",
						"Терапия для сестры - возвращение себе детства",
					},
				}
			}
			
			// Феномен "девочка с яйцами" и ее брат
			if sisterRole.PrimaryRole == core.TraumaRoleScapegoat &&
				brotherRole.PrimaryRole == core.TraumaRoleLostChild {
				sister := fs.Members[pair[0]]
				brother := fs.Members[pair[1]]
				
				// Проверяем амплитуду векторов сестры
				coords := core.ParseDateToCoords(sister.BirthDate)
				vectors := core.CalculateVectors(coords)
				amp := core.VectorAmplitude(vectors)
				
				if amp > 5.0 {
					return &FamilyPattern{
						PatternType: "sister-brother-girl-with-balls",
						Description: fmt.Sprintf("Феномен 'девочка с яйцами' - сестра (%s) с маскулинными стратегиями и потерянный брат (%s)",
							sister.Name, brother.Name),
						Members:     []string{pair[0], pair[1]},
						Severity:    0.8,
						Generations: []int{sister.Generation},
						Recommendations: []string{
							"Интеграция маскулинной и феминной энергии у сестры",
							"Активация мужского начала у брата",
							"Баланс ролей в паре",
							"Терапия для обоих - восстановление естественных ролей",
						},
					}
				}
			}
			
			// Защитник и защищаемая
			if brotherRole.PrimaryRole == core.TraumaRoleGoldenChild &&
				sisterRole.PrimaryRole == core.TraumaRoleInvisible {
				brother := fs.Members[pair[1]]
				sister := fs.Members[pair[0]]
				
				return &FamilyPattern{
					PatternType: "sister-brother-protector",
					Description: fmt.Sprintf("Брат (%s) защищает невидимую сестру (%s)",
						brother.Name, sister.Name),
					Members:     []string{pair[0], pair[1]},
					Severity:    0.6,
					Generations: []int{sister.Generation},
					Recommendations: []string{
						"Осознание роли защитника братом",
						"Активация сестры, поиск собственного голоса",
						"Баланс в отношениях без спасательства",
					},
				}
			}
		}
	}
	
	return nil
}

// findGrandmotherGranddaughterPattern - бабушка-внучка (3 поколения)
func (fs *FamilySystem) findGrandmotherGranddaughterPattern() *FamilyPattern {
	for gmID, gm := range fs.Members {
		if gm.Gender != core.GenderFemale || gm.Role != core.RoleGrandmother {
			continue
		}
		
		for _, childID := range gm.Children {
			mother := fs.Members[childID]
			if mother == nil || mother.Gender != core.GenderFemale || mother.Role != core.RoleMother {
				continue
			}
			
			for _, grandchildID := range mother.Children {
				grandchild := fs.Members[grandchildID]
				if grandchild == nil || grandchild.Gender != core.GenderFemale {
					continue
				}
				
				gmRole := fs.TraumaRoles[gmID]
				motherRole := fs.TraumaRoles[childID]
				gcRole := fs.TraumaRoles[grandchildID]
				
				if gmRole != nil && motherRole != nil && gcRole != nil {
					// Проверяем передачу роли
					if gmRole.PrimaryRole == motherRole.PrimaryRole ||
						motherRole.PrimaryRole == gcRole.PrimaryRole {
						
						roleName := ""
						switch gmRole.PrimaryRole {
						case core.TraumaRoleScapegoat:
							roleName = "козла отпущения"
						case core.TraumaRoleGoldenChild:
							roleName = "любимчика"
						case core.TraumaRoleLostChild:
							roleName = "потерянного ребенка"
						default:
							roleName = string(gmRole.PrimaryRole)
						}
						
						return &FamilyPattern{
							PatternType: "grandmother-mother-granddaughter",
							Description: fmt.Sprintf("Передача роли '%s' через три поколения женщин", roleName),
							Members:     []string{gmID, childID, grandchildID},
							Severity:    0.95,
							Generations: []int{1, 2, 3},
							Recommendations: []string{
								"Глубокая родовая терапия",
								"Работа с женской линией",
								"Разрыв родового проклятия",
								"Индивидуальная терапия для каждой",
								"Восстановление женской идентичности",
							},
						}
					}
				}
			}
		}
	}
	return nil
}

// findGrandfatherGrandsonPattern - дедушка-отец-внук
func (fs *FamilySystem) findGrandfatherGrandsonPattern() *FamilyPattern {
	for gfID, gf := range fs.Members {
		if gf.Gender != core.GenderMale || gf.Role != core.RoleGrandfather {
			continue
		}
		
		for _, childID := range gf.Children {
			father := fs.Members[childID]
			if father == nil || father.Gender != core.GenderMale || father.Role != core.RoleFather {
				continue
			}
			
			for _, grandchildID := range father.Children {
				grandchild := fs.Members[grandchildID]
				if grandchild == nil || grandchild.Gender != core.GenderMale {
					continue
				}
				
				gfRole := fs.TraumaRoles[gfID]
				fatherRole := fs.TraumaRoles[childID]
				gcRole := fs.TraumaRoles[grandchildID]
				
				if gfRole != nil && fatherRole != nil && gcRole != nil {
					if gfRole.PrimaryRole == fatherRole.PrimaryRole ||
						fatherRole.PrimaryRole == gcRole.PrimaryRole {
						
						return &FamilyPattern{
							PatternType: "grandfather-father-grandson",
							Description: "Передача травмы через три поколения мужчин",
							Members:     []string{gfID, childID, grandchildID},
							Severity:    0.95,
							Generations: []int{1, 2, 3},
							Recommendations: []string{
								"Мужская родовая терапия",
								"Работа с родовыми сценариями",
								"Восстановление мужской линии",
								"Ритуалы передачи силы",
							},
						}
					}
				}
			}
		}
	}
	return nil
}

// findGrandmotherGrandsonPattern - бабушка-внук (через поколение)
func (fs *FamilySystem) findGrandmotherGrandsonPattern() *FamilyPattern {
	for gmID, gm := range fs.Members {
		if gm.Gender != core.GenderFemale || gm.Role != core.RoleGrandmother {
			continue
		}
		
		for _, childID := range gm.Children {
			parent := fs.Members[childID]
			if parent == nil {
				continue
			}
			
			for _, grandchildID := range parent.Children {
				grandchild := fs.Members[grandchildID]
				if grandchild == nil || grandchild.Gender != core.GenderMale {
					continue
				}
				
				// Проверяем, не был ли родитель "потерянным" для бабушки
				parentRole := fs.TraumaRoles[childID]
				if parentRole != nil && parentRole.PrimaryRole == core.TraumaRoleLostChild {
					return &FamilyPattern{
						PatternType: "grandmother-grandson-compensation",
						Description: "Бабушка компенсирует через внука отвергнутого сына",
						Members:     []string{gmID, grandchildID},
						Severity:    0.8,
						Generations: []int{1, 3},
						Recommendations: []string{
							"Осознание переноса с сына на внука",
							"Работа с сыном (пропущенным поколением)",
							"Терапия для внука от слияния",
							"Восстановление связи с отцом",
						},
					}
				}
			}
		}
	}
	return nil
}

// findGrandfatherGranddaughterPattern - дедушка-внучка
func (fs *FamilySystem) findGrandfatherGranddaughterPattern() *FamilyPattern {
	for gfID, gf := range fs.Members {
		if gf.Gender != core.GenderMale || gf.Role != core.RoleGrandfather {
			continue
		}
		
		for _, childID := range gf.Children {
			parent := fs.Members[childID]
			if parent == nil {
				continue
			}
			
			for _, grandchildID := range parent.Children {
				grandchild := fs.Members[grandchildID]
				if grandchild == nil || grandchild.Gender != core.GenderFemale {
					continue
				}
				
				// Особая связь дедушка-внучка
				gfRole := fs.TraumaRoles[gfID]
				gcRole := fs.TraumaRoles[grandchildID]
				
				if gfRole != nil && gcRole != nil {
					return &FamilyPattern{
						PatternType: "grandfather-granddaughter-bond",
						Description: "Особая связь дедушки и внучки (часто компенсаторная или идеализированная)",
						Members:     []string{gfID, grandchildID},
						Severity:    0.6,
						Generations: []int{1, 3},
						Recommendations: []string{
							"Проверка на эмоциональный инцест",
							"Работа с женской идентичностью",
							"Отделение от дедушкиных проекций",
							"Формирование здоровых отношений с мужчинами",
						},
					}
				}
			}
		}
	}
	return nil
}

// findReplacementChildPattern - паттерн замещающих детей
func (fs *FamilySystem) findReplacementChildPattern() *FamilyPattern {
	for id, member := range fs.Members {
		if member.IsReplacement && member.ReplacementData != nil {
			// Находим, кого замещает
			replacesID := member.ReplacesWho
			var lostParent string
			
			for _, parentID := range member.Parents {
				if parent, ok := fs.Members[parentID]; ok {
					for _, lost := range parent.LostChildren {
						if lost.ID == replacesID {
							lostParent = parentID
							break
						}
					}
				}
			}
			
			members := []string{id}
			if lostParent != "" {
				members = append(members, lostParent)
			}
			
			severity := 0.7
			if member.ReplacementData.TimeGapMonths < 6 {
				severity = 0.9
			}
			if member.ReplacementData.IsExplicit {
				severity += 0.1
			}
			if severity > 1.0 {
				severity = 1.0
			}
			
			return &FamilyPattern{
				PatternType: "replacement-child",
				Description: "Замещающий ребенок - рожден после потери предыдущего",
				Members:     members,
				Severity:    severity,
				Generations: []int{member.Generation},
				Recommendations: []string{
					"Терапия идентичности - отделение от образа потерянного",
					"Работа с виной выжившего",
					"Ритуал прощания с потерянным сиблингом",
					"Признание права на собственную жизнь",
					"Терапия для родителей - работа с горем",
				},
			}
		}
	}
	return nil
}

// findDivorcePattern - паттерн повторяющихся разводов
func (fs *FamilySystem) findDivorcePattern() *FamilyPattern {
	divorceGenerations := make(map[int][]string)
	
	for id, member := range fs.Members {
		divorceCount := 0
		for _, event := range member.Events {
			if event.EventType == "divorce" {
				divorceCount++
			}
		}
		
		if divorceCount >= 1 {
			divorceGenerations[member.Generation] = append(divorceGenerations[member.Generation], id)
		}
	}
	
	if len(divorceGenerations) >= 2 {
		// Собираем всех участников
		members := make([]string, 0)
		gens := make([]int, 0)
		
		for gen, ids := range divorceGenerations {
			gens = append(gens, gen)
			members = append(members, ids...)
		}
		
		sort.Ints(gens)
		
		genStr := ""
		for i, gen := range gens {
			if i > 0 {
				genStr += ", "
			}
			genStr += fmt.Sprintf("%d", gen)
		}
		
		return &FamilyPattern{
			PatternType: "generational-divorce",
			Description: fmt.Sprintf("Повторение разводов в поколениях %s", genStr),
			Members:     members,
			Severity:    0.8,
			Generations: gens,
			Recommendations: []string{
				"Анализ родовых сценариев отношений",
				"Терапия привязанности",
				"Работа с выбором партнера",
				"Осознание повторяющихся паттернов",
				"Проработка страха близости",
			},
		}
	}
	
	return nil
}

// findInfidelityPattern - паттерн измен
func (fs *FamilySystem) findInfidelityPattern() *FamilyPattern {
	if !fs.ExtendedMode {
		return nil
	}
	
	infidelityMembers := make([]string, 0)
	
	for id, member := range fs.Members {
		if len(member.Infidelities) > 0 {
			infidelityMembers = append(infidelityMembers, id)
		}
	}
	
	if len(infidelityMembers) >= 2 {
		// Проверяем, есть ли повторение в поколениях
		genMap := make(map[int]bool)
		for _, id := range infidelityMembers {
			if member, ok := fs.Members[id]; ok {
				genMap[member.Generation] = true
			}
		}
		
		gens := make([]int, 0, len(genMap))
		for g := range genMap {
			gens = append(gens, g)
		}
		sort.Ints(gens)
		
		if len(gens) >= 2 {
			return &FamilyPattern{
				PatternType: "generational-infidelity",
				Description: "Повторяющиеся измены в нескольких поколениях",
				Members:     infidelityMembers,
				Severity:    0.75,
				Generations: gens,
				Recommendations: []string{
					"Терапия доверия",
					"Работа с семейными тайнами",
					"Восстановление целостности",
					"Проработка травмы предательства",
					"Осознание причин измен в каждом поколении",
				},
			}
		}
	}
	
	return nil
}

// findLargeFamilyPattern - паттерны больших семей (10+ детей)
func (fs *FamilySystem) findLargeFamilyPattern() *FamilyPattern {
	largeFamilies := make(map[string][]string) // parentID -> children
	
	for _, member := range fs.Members {
		if len(member.Children) >= 10 {
			largeFamilies[member.ID] = member.Children
		}
	}
	
	if len(largeFamilies) == 0 {
		return nil
	}
	
	for parentID, children := range largeFamilies {
		parent := fs.Members[parentID]
		
		// Анализ ролей в большой семье
		roleCount := make(map[core.TraumaRole]int)
		for _, childID := range children {
			if role, ok := fs.TraumaRoles[childID]; ok {
				roleCount[role.PrimaryRole]++
			}
		}
		
		// Формируем описание
		description := fmt.Sprintf("Большая семья (%d детей) - ", len(children))
		if roleCount[core.TraumaRoleParentified] > 0 {
			description += "старшие дети родифицированы, "
		}
		if roleCount[core.TraumaRoleLostChild] > 2 {
			description += "многие дети 'потеряны', "
		}
		if roleCount[core.TraumaRoleInvisible] > 2 {
			description += "есть 'невидимки', "
		}
		
		return &FamilyPattern{
			PatternType: "large-family-dynamics",
			Description: description,
			Members:     children,
			Severity:    0.7,
			Generations: []int{parent.Generation, parent.Generation + 1},
			Recommendations: []string{
				"Распределение внимания между всеми детьми",
				"Предотвращение родификации старших",
				"Индивидуальный подход к каждому ребенку",
				"Групповая терапия для сиблингов",
				"Работа с 'потерянными' и 'невидимыми' детьми",
				"Осознание родителями особенностей большой семьи",
			},
		}
	}
	
	return nil
}
