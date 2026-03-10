// Package main — точка входа idealcore
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ruvcoindev/idealcore/pkg/ai"
	"github.com/ruvcoindev/idealcore/pkg/config"
	"github.com/ruvcoindev/idealcore/pkg/diary"
	"github.com/ruvcoindev/idealcore/pkg/vector"
	"github.com/ruvcoindev/idealcore/pkg/web"
	"github.com/ruvcoindev/idealcore/pkg/yggdrasil"
)

func main() {
	log.Println("🚀 idealcore сервер запускается...")

	// 1. Конфигурация
	cfg := config.Load()
	log.Printf("📋 Порт: %s, Ollama: %s:%s", cfg.ServerPort, cfg.OllamaHost, cfg.OllamaPort)

	// 2. Хранилища
	diaryStore, err := diary.NewStore(cfg.DataDir)
	if err != nil {
		log.Fatalf("❌ diary store: %v", err)
	}
	log.Println("✅ diary store ready")

	vectorStore := vector.NewStore(cfg.VectorDim)
	log.Println("✅ vector store ready")

	// 3. AI клиент
	ollamaURL := fmt.Sprintf("http://%s:%s", cfg.OllamaHost, cfg.OllamaPort)
	aiClient := ai.NewOllamaClient(ollamaURL, vectorStore)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := aiClient.HealthCheck(ctx); err != nil {
		log.Printf("⚠️ Ollama: %v", err)
	} else {
		log.Println("✅ Ollama connected")
	}

	// 4. Yggdrasil (опционально)
	if cfg.YggdrasilEnable {
		ygg := yggdrasil.NewTransport(nil)
		_ = ygg.Start()
		log.Printf("✅ Yggdrasil: %s", ygg.GetNodeID())
	}

	// 5. Веб-обработчики
	handlers, err := web.NewHandlers(cfg, aiClient, diaryStore, vectorStore)
	if err != nil {
		log.Fatalf("❌ handlers: %v", err)
	}
	log.Println("✅ handlers ready")

	// 6. Router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	handlers.RegisterRoutes(router)
	log.Println("✅ routes registered")

	// 7. Server
	addr := ":" + cfg.ServerPort
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 8. Start
	go func() {
		log.Printf("🌐 Listening on http://%s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ server: %v", err)
		}
	}()

	// 9. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	_ = server.Shutdown(shutdownCtx)

	log.Println("✅ stopped")
}
