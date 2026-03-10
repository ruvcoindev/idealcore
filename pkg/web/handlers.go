// Package web — HTTP обработчики idealcore
package web

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ruvcoindev/idealcore/pkg/ai"
	"github.com/ruvcoindev/idealcore/pkg/config"
	"github.com/ruvcoindev/idealcore/pkg/diary"
	"github.com/ruvcoindev/idealcore/pkg/vector"
)

// Handlers — коллекция обработчиков
type Handlers struct {
	cfg         *config.Config
	aiClient    *ai.OllamaClient
	diaryStore  *diary.Store
	vectorStore vector.VectorStore
	templates   *template.Template
}

// PageData — данные для шаблонов
type PageData struct {
	Title       string
	CurrentPath string
	UserID      string
	Description string
}

// NewHandlers создаёт обработчики
func NewHandlers(cfg *config.Config, aiClient *ai.OllamaClient, ds *diary.Store, vs vector.VectorStore) (*Handlers, error) {
	// Загружаем шаблоны
	tmplPath := filepath.Join("pkg", "web", "templates", "*.html")
	tmpl, err := template.ParseGlob(tmplPath)
	if err != nil {
		return nil, err
	}

	return &Handlers{
		cfg:         cfg,
		aiClient:    aiClient,
		diaryStore:  ds,
		vectorStore: vs,
		templates:   tmpl,
	}, nil
}

// RegisterRoutes регистрирует маршруты
func (h *Handlers) RegisterRoutes(router *gin.Engine) {
	router.Static("/static", "pkg/web/static")
	router.GET("/manifest.json", h.manifestHandler)
	router.GET("/sw.js", h.swHandler)
	router.GET("/", h.indexHandler)
	router.GET("/diary", h.diaryHandler)
	router.GET("/intention", h.intentionHandler)
	router.GET("/insights", h.insightsHandler)
	router.GET("/health", h.healthHandler)

	api := router.Group("/api")
	{
		api.POST("/diary/save", h.saveEntryHandler)
		api.GET("/diary/entries", h.getEntriesHandler)
		api.POST("/intention/generate", h.generateIntentionHandler)
		api.GET("/intention/themes", h.getThemesHandler)
		api.GET("/models", h.getModelsHandler)
		api.GET("/diary/questions", h.getQuestionsHandler)
	}
}

// === Страницы с шаблонами ===

func (h *Handlers) indexHandler(c *gin.Context) {
	data := PageData{
		Title:       "idealcore",
		Description: "твой путь к целостности",
		CurrentPath: "/",
	}
	h.templates.ExecuteTemplate(c.Writer, "index.html", data)
}

func (h *Handlers) diaryHandler(c *gin.Context) {
	userID := c.Query("user")
	if userID == "" {
		userID = "anonymous"
	}
	data := PageData{
		Title:       "Дневник «Кто я»",
		Description: "Опросник для саморефлексии",
		CurrentPath: "/diary",
		UserID:      userID,
	}
	h.templates.ExecuteTemplate(c.Writer, "diary.html", data)
}

func (h *Handlers) intentionHandler(c *gin.Context) {
	userID := c.Query("user")
	if userID == "" {
		userID = "anonymous"
	}
	data := PageData{
		Title:       "Генератор намерений",
		Description: "Намерения для практики на гвоздях",
		CurrentPath: "/intention",
		UserID:      userID,
	}
	h.templates.ExecuteTemplate(c.Writer, "intention.html", data)
}

func (h *Handlers) insightsHandler(c *gin.Context) {
	userID := c.Query("user")
	if userID == "" {
		userID = "anonymous"
	}
	data := PageData{
		Title:       "Инсайты",
		Description: "Твои паттерны и ресурсы",
		CurrentPath: "/insights",
		UserID:      userID,
	}
	h.templates.ExecuteTemplate(c.Writer, "insights.html", data)
}

func (h *Handlers) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "idealcore",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (h *Handlers) manifestHandler(c *gin.Context) {
	c.Header("Content-Type", "application/manifest+json")
	c.String(http.StatusOK, `{"name":"idealcore","start_url":"/","display":"standalone"}`)
}

func (h *Handlers) swHandler(c *gin.Context) {
	c.Header("Content-Type", "application/javascript")
	c.String(http.StatusOK, "// service worker stub")
}

// === API ===

type SaveEntryRequest struct {
	UserID  string   `json:"user_id"`
	Section string   `json:"section"`
	Answer  string   `json:"answer"`
	Tags    []string `json:"tags"`
}

func (h *Handlers) saveEntryHandler(c *gin.Context) {
	var req SaveEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.diaryStore.SaveEntry(req.UserID, req.Section, req.Answer, req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "saved"})
}

func (h *Handlers) getEntriesHandler(c *gin.Context) {
	userID := c.Query("user")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}
	entries, err := h.diaryStore.GetEntries(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"entries": entries})
}

type GenerateIntentionRequest struct {
	UserID string `json:"user_id"`
	Theme  string `json:"theme"`
}

func (h *Handlers) generateIntentionHandler(c *gin.Context) {
	var req GenerateIntentionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	intention, err := h.aiClient.GenerateIntention(ctx, ai.IntentionContext{
		UserID: req.UserID,
		Theme:  req.Theme,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"intention": intention})
}

func (h *Handlers) getThemesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"themes": []string{
		"границы с родителями", "право на отдых", "отношения с деньгами",
		"принятие себя", "страх брошенности", "вина и стыд",
		"право на ошибку", "выражение гнева", "доверие к себе", "сепарация",
	}})
}

func (h *Handlers) getModelsHandler(c *gin.Context) {
	models, err := h.aiClient.ListModels(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"models": models})
}

// getQuestionsHandler возвращает вопросы опросника
func (h *Handlers) getQuestionsHandler(c *gin.Context) {
	sections := diary.GetSections()
	questions := diary.GetQuestions()
	c.JSON(http.StatusOK, gin.H{
		"sections":  sections,
		"questions": questions,
	})
}
