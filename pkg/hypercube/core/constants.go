package core

// CubeSize — размер куба (26 комнат по каждой оси, как в фильме "Куб")
// Каждая координата может принимать значения от 0 до 25, что символизирует
// ограниченность жизненного пространства и повторяющиеся паттерны
const CubeSize = 26

// RelationshipType определяет тип отношений между двумя людьми
// Используется для анализа специфической динамики в разных типах связей
type RelationshipType string

const (
	// Родительско-детские отношения
	RelationshipMotherDaughter RelationshipType = "mother-daughter" // Мать-дочь: особая связь с передачей женских сценариев
	RelationshipMotherSon      RelationshipType = "mother-son"      // Мать-сын: риск эмоционального инцеста и симбиоза
	RelationshipFatherDaughter RelationshipType = "father-daughter" // Отец-дочь: формирование образа мужчины, комплекс Электры
	RelationshipFatherSon      RelationshipType = "father-son"      // Отец-сын: передача мужских сценариев, конкуренция

	// Сиблинговые отношения (братья и сестры)
	RelationshipSisterSister   RelationshipType = "sister-sister"   // Сестра-сестра: конкуренция, зависть, поддержка
	RelationshipBrotherBrother  RelationshipType = "brother-brother"  // Брат-брат: иерархия, соперничество, союз
	RelationshipSisterBrother   RelationshipType = "sister-brother"   // Сестра-брат: комплементарные роли, защита

	// Межпоколенческие отношения (бабушки-дедушки и внуки)
	RelationshipGrandmotherGranddaughter RelationshipType = "grandmother-granddaughter" // Бабушка-внучка: передача через поколение
	RelationshipGrandmotherGrandson      RelationshipType = "grandmother-grandson"     // Бабушка-внук: компенсация отношений с сыном
	RelationshipGrandfatherGranddaughter RelationshipType = "grandfather-granddaughter" // Дедушка-внучка: особая нежность
	RelationshipGrandfatherGrandson      RelationshipType = "grandfather-grandson"     // Дедушка-внук: передача мужской линии

	// Партнерские отношения
	RelationshipHusbandWife      RelationshipType = "husband-wife"       // Муж-жена: супружеские отношения
	RelationshipCivilPartnership RelationshipType = "civil-partnership"   // Гражданский брак: официально не зарегистрированные отношения
	RelationshipExSpouses        RelationshipType = "ex-spouses"         // Бывшие супруги: непроработанные связи
	RelationshipRomantic         RelationshipType = "romantic"           // Романтические отношения без совместного проживания

	// Социальные отношения
	RelationshipFriendly RelationshipType = "friendly" // Дружеские отношения
	RelationshipWork     RelationshipType = "work"     // Рабочие/профессиональные отношения
)

// Gender определяет пол человека
// Используется для анализа гендерных паттернов в семье
type Gender string

const (
	GenderMale   Gender = "male"   // Мужской
	GenderFemale Gender = "female" // Женский
	GenderOther  Gender = "other"  // Другой (небинарный, не указан)
)

// FamilyRole определяет базовую роль человека в семейной системе
// Эти роли задаются пользователем, в отличие от травматических ролей,
// которые вычисляются автоматически
type FamilyRole string

const (
	RoleMother       FamilyRole = "mother"        // Мать
	RoleFather       FamilyRole = "father"        // Отец
	RoleDaughter     FamilyRole = "daughter"      // Дочь
	RoleSon          FamilyRole = "son"           // Сын
	RoleGrandmother  FamilyRole = "grandmother"   // Бабушка
	RoleGrandfather  FamilyRole = "grandfather"   // Дедушка
	RoleGranddaughter FamilyRole = "granddaughter" // Внучка
	RoleGrandson     FamilyRole = "grandson"      // Внук
	RoleHusband      FamilyRole = "husband"       // Муж
	RoleWife         FamilyRole = "wife"          // Жена
	RolePartner      FamilyRole = "partner"       // Партнер (не состоящий в браке)
)

// TraumaRole определяет роль в травматической системе семьи
// Эти роли ВЫЯВЛЯЮТСЯ АВТОМАТИЧЕСКИ на основе анализа семейной динамики,
// порядка рождения, дат рождения и жизненных событий
type TraumaRole string

const (
	// Любимчик (Golden Child) - ребенок, которого идеализируют
	// На него возложены ожидания, он должен быть совершенным
	// Часто страдает от перфекционизма и страха неудачи
	TraumaRoleGoldenChild TraumaRole = "golden_child"

	// Козел отпущения (Scapegoat) - ребенок, на которого проецируются все проблемы семьи
	// Его критикуют, обвиняют, он "плохой" в семье
	// Часто становится независимым, но с глубоким чувством стыда
	TraumaRoleScapegoat TraumaRole = "scapegoat"

	// Потерянный ребенок (Lost Child) - ребенок, который "исчезает" эмоционально
	// Проводит много времени один, не доставляет хлопот
	// Страдает от одиночества и неспособности устанавливать близкие связи
	TraumaRoleLostChild TraumaRole = "lost_child"

	// Невидимка (Invisible Child) - ребенок, чьи потребности игнорируются
	// Его как бы не замечают, он существует на периферии семейного внимания
	TraumaRoleInvisible TraumaRole = "invisible"

	// Стеклянный ребенок (Glass Child) - ребенок, который становится "невидимым"
	// на фоне проблемного сиблинга (с инвалидностью, болезнью, зависимостью)
	// Выглядит "прозрачным", потому что все внимание уходит другому
	TraumaRoleGlassChild TraumaRole = "glass_child"

	// Родифицированный ребенок (Parentified Child) - ребенок, который вынужден
	// выполнять родительские функции, заботиться о родителях или сиблингах
	// Лишен детства, становится гиперответственным
	TraumaRoleParentified TraumaRole = "parentified"

	// Шут (Mascot) - ребенок, который разряжает напряжение юмором
	// Использует смех как защиту, минимизирует серьезные проблемы
	TraumaRoleMascot TraumaRole = "mascot"

	// Эмоциональный супруг (Emotional Spouse) - ребенок, который заменяет
	// партнера одному из родителей (обычно противоположного пола)
	// Находится в симбиотической связи, не может создать здоровые отношения
	TraumaRoleEmotionalSpouse TraumaRole = "emotional_spouse"

	// Говорящий правду (Truth Teller) - ребенок, который рискует говорить
	// о семейных проблемах, которые все игнорируют
	// Часто становится изгоем за свою честность
	TraumaRoleTruthTeller TraumaRole = "truth_teller"

	// Тень семьи (Shadow) - ребенок, который несет непрожитые эмоции семьи
	// Часто депрессивен, тревожен, выражает то, что подавлено в других
	TraumaRoleShadow TraumaRole = "shadow"

	// Целитель (Healer) - ребенок, который пытается "исцелить" семью
	// Становится психологом, врачом, спасателем
	TraumaRoleHealer TraumaRole = "healer"

	// Замещающий ребенок (Replacement Child) - ребенок, рожденный после потери
	// предыдущего ребенка (смерть, выкидыш, аборт)
	// Живет с "призраком" потерянного сиблинга, чувствует, что должен его заменить
	TraumaRoleReplacement TraumaRole = "replacement_child"

	// Призрак (Ghost) - потерянный ребенок, который продолжает влиять на семью
	// Его присутствие ощущается, о нем говорят, его оплакивают годами
	TraumaRoleGhost TraumaRole = "ghost"

	// Предок-носитель травмы (Ancestor) - предок, чья травма передается
	// через поколения, даже если сам предок уже умер
	TraumaRoleAncestor TraumaRole = "ancestor"
)

// LossType определяет тип потери ребенка
// Используется для анализа травм, связанных с детской смертностью и абортами
type LossType string

const (
	LossTypeAbortion      LossType = "abortion"       // Аборт - искусственное прерывание беременности
	LossTypeMiscarriage   LossType = "miscarriage"    // Выкидыш - самопроизвольное прерывание
	LossTypeStillborn     LossType = "stillborn"      // Мертворожденный - ребенок родился мертвым
	LossTypeInfantDeath   LossType = "infant_death"   // Смерть младенца (до 1 года)
	LossTypeChildDeath    LossType = "child_death"    // Смерть ребенка (до 18 лет)
	LossTypeAdultDeath    LossType = "adult_death"    // Смерть взрослого (после 18 лет)
)

// TraumaThresholds — пороговые значения для определения травматических паттернов
// Используются в алгоритмах автоматического определения ролей
var TraumaThresholds = struct {
	ErraticMin       int // Минимальная амплитуда колебаний для определения "нестабильности"
	FrozenMax        int // Максимальное изменение для определения "заморозки"
	RepetitiveCount  int // Количество повторений для определения "повторяющегося паттерна"
}{
	ErraticMin:      3,  // Если амплитуда больше 3 - паттерн считается нестабильным
	FrozenMax:       1,  // Если изменения меньше 2 - паттерн считается "замороженным"
	RepetitiveCount: 2,  // Если паттерн повторяется 2+ раза - считается цикличным
}
