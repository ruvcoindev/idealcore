// Package domain содержит основные доменные модели
package domain

import "time"

// User — профиль пользователя
type User struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	BirthDate       string                 `json:"birth_date"` // YYYY-MM-DD
	BirthPlace      string                 `json:"birth_place"`
	AttachmentType  AttachmentType         `json:"attachment_type"`
	TraumaFlags     []TraumaType           `json:"trauma_flags"`
	DefenseMechanisms []DefenseMechanism   `json:"defense_mechanisms"`
	ChakraBlocks    []ChakraType           `json:"chakra_blocks"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// DiaryEntry — запись дневника "Кто я"
type DiaryEntry struct {
	ID          string            `json:"id"`
	UserID      string            `json:"user_id"`
	Section     string            `json:"section"` // motivation, boundaries, resources, patterns, choice
	Question    string            `json:"question"`
	Answer      string            `json:"answer"`
	Tags        []string          `json:"tags"`
	Sentiment   float32           `json:"sentiment"` // -1..1
	Vector      []float32         `json:"vector,omitempty"` // эмбеддинг
	CreatedAt   time.Time         `json:"created_at"`
}

// Intention — намерение для практики на гвоздях
type Intention struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Text            string    `json:"text"`
	Theme           string    `json:"theme"` // money, boundaries, self_worth...
	TraumaContext   []TraumaType `json:"trauma_context"`
	AttachmentStyle AttachmentType `json:"attachment_style"`
	ChakraFocus     []ChakraType `json:"chakra_focus"`
	Complexity      int       `json:"complexity"` // 1..5
	UsedAt          time.Time `json:"used_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// HypercubeCoords — координаты в 4D пространстве (модель из фильма)
type HypercubeCoords struct {
	X, Y, Z, W int32 `json:"x,y,z,w"` // W — четвертое измерение
}

// RoomVector — вектор движения "комнаты" (человека)
type RoomVector struct {
	Coords  HypercubeCoords `json:"coords"`
	Vectors [3][3]int32     `json:"vectors"` // 3 вектора × 3 компоненты
	Cycle   int             `json:"cycle"`   // номер шага в цикле
}

// Типы привязанности
type AttachmentType string

const (
	AttachmentSecure      AttachmentType = "secure"
	AttachmentAnxious     AttachmentType = "anxious"
	AttachmentAvoidant    AttachmentType = "avoidant"
	AttachmentDisorganized AttachmentType = "disorganized"
)

// Типы травм
type TraumaType string

const (
	TraumaPTSD      TraumaType = "ptsd"       // Посттравматическое стрессовое расстройство
	TraumaCPTSD     TraumaType = "cptsd"      // Комплексное ПТСР
	TraumaBPD       TraumaType = "bpd"        // Пограничное расстройство личности
	TraumaNPD       TraumaType = "npd"        // Нарциссическое расстройство личности
	TraumaAttachment TraumaType = "attachment" // Травма привязанности
	TraumaNarcissisticParent TraumaType = "narcissistic_parent"
	TraumaAbandonment TraumaType = "abandonment"
	TraumaBetrayal  TraumaType = "betrayal"
	TraumaShame     TraumaType = "shame"
)

// Защитные механизмы
type DefenseMechanism string

const (
	DefenseDenial         DefenseMechanism = "denial"          // Отрицание
	DefenseProjection     DefenseMechanism = "projection"      // Проекция
	DefenseRationalization DefenseMechanism = "rationalization" // Рационализация
	DefenseDisplacement   DefenseMechanism = "displacement"    // Смещение
	DefenseSublimation    DefenseMechanism = "sublimation"     // Сублимация
	DefenseRegression     DefenseMechanism = "regression"      // Регрессия
	DefenseReactionFormation DefenseMechanism = "reaction_formation" // Реактивное образование
	DefenseIntellectualization DefenseMechanism = "intellectualization" // Интеллектуализация
	DefenseNarcissisticGrandiosity DefenseMechanism = "narcissistic_grandiosity"
	DefenseNarcissisticDevaluation DefenseMechanism = "narcissistic_devaluation"
	DefenseTriangulation  DefenseMechanism = "triangulation"   // Триангуляция
)

// Типы чакр
type ChakraType string

const (
	ChakraRoot      ChakraType = "root"      // Муладхара — выживание, безопасность
	ChakraSacral    ChakraType = "sacral"    // Свадхистана — эмоции, творчество
	ChakraSolar     ChakraType = "solar"     // Манипура — воля, контроль
	ChakraHeart     ChakraType = "heart"     // Анахата — любовь, принятие
	ChakraThroat    ChakraType = "throat"    // Вишудха — выражение, правда
	ChakraThirdEye  ChakraType = "third_eye" // Аджна — интуиция, видение
	ChakraCrown     ChakraType = "crown"     // Сахасрара — связь, трансценденция
)
