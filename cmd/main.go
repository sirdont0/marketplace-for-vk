package main

import (
	"marketplace/internal/auth"
	"marketplace/internal/handlers"
	"marketplace/internal/storage"

	"strings"
	"os"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
    //"github.com/gin-contrib/zap"
)

func main() {
	// Инициализация базы данных
	if err := storage.InitDB(); err != nil {
		panic(err)
	}

	// Создание роутера Gin
	r := gin.Default()

	// Регистрация middleware для валидации изображений
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("image_url", func(fl validator.FieldLevel) bool {
			url := fl.Field().String()
			validExts := []string{".jpg", ".png", ".webp", ".jpeg"}
			for _, ext := range validExts {
				if strings.HasSuffix(strings.ToLower(url), ext) {
					return true
				}
			}
			return false
		})
	}

	// Группа публичных эндпоинтов
	public := r.Group("/api")
	{
		public.POST("/auth/register", handlers.Register)
		public.POST("/auth/login", handlers.Login)
		public.GET("/ads", handlers.GetAds)
	}

	// Группа приватных эндпоинтов (требует авторизации)
	private := r.Group("/api")
	private.Use(auth.AuthMiddleware())
	{
		private.POST("/ads", handlers.CreateAd)
	}

	// Запуск сервера
	//r.Run(":8080")
	certFile := os.Getenv("SSL_CERT_PATH")  // Путь к cert.pem
    keyFile := os.Getenv("SSL_KEY_PATH")    // Путь к key.pem

    if certFile != "" && keyFile != "" {
        // Автоматический редирект с HTTP на HTTPS
        go func() {
            if err := http.ListenAndServe(":80", http.HandlerFunc(redirectToHTTPS)); err != nil {
                zap.L().Fatal("HTTP server error", zap.Error(err))
            }
        }()

        // Запуск HTTPS
        zap.L().Info("Starting HTTPS server")
        if err := r.RunTLS(":443", certFile, keyFile); err != nil {
            zap.L().Fatal("HTTPS server error", zap.Error(err))
        }
    } else {
        zap.L().Warn("Running in HTTP mode (no SSL certificates detected)")
        if err := r.Run(":8080"); err != nil {
            zap.L().Fatal("Server error", zap.Error(err))
        }
    }
}

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
    target := "https://" + r.Host + r.URL.Path
    http.Redirect(w, r, target, http.StatusMovedPermanently)
}
