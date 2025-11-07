package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Stenoliv/didlydoodash_api/internal/config"
	"github.com/Stenoliv/didlydoodash_api/internal/db"
	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
	"github.com/Stenoliv/didlydoodash_api/internal/handlers"
	"github.com/Stenoliv/didlydoodash_api/internal/middleware"
	"github.com/Stenoliv/didlydoodash_api/internal/repositories"
	"github.com/Stenoliv/didlydoodash_api/internal/services"
	"github.com/Stenoliv/didlydoodash_api/pkg/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment
	_ = godotenv.Load()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Initialize structured logger
	logger := logging.New(cfg.Mode)
	logger.Infof("Starting DidlyDooDash API in %s mode", cfg.Mode)

	// Create gin instance
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logging.Middleware(logger))
	r.Use(middleware.ErrorHandler())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CorsOrigins,
		AllowMethods:     []string{"POST", "PUT", "PATCH", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "_retry"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Init DB
	pgx, err := db.Load()
	if err != nil {
		logger.Fatalf("failed ot connect to database: %v", err)
	}

	// Repositories
	repo := repository.New(pgx)
	txManager := repositories.NewTxManager(pgx)
	orgRepo := repositories.NewOrganisationRepo(repo, logger)
	userRepo := repositories.NewUserRepository(repo, logger)
	roleRepo := repositories.NewRoleRepo(repo, logger)
	memberRepo := repositories.NewMemberRepo(repo, logger)

	// Services
	checkerService := services.NewChecker(memberRepo, roleRepo, logger)
	authService := services.NewAuthService(userRepo, txManager, cfg, logger)
	orgService := services.NewOrganisationService(services.OrganisationServiceRepos{
		Org:    orgRepo,
		Member: memberRepo,
		Role:   roleRepo,
	}, txManager, logger)
	membershipService := services.NewMembershipService(services.MembershipRepos{
		Role:   roleRepo,
		Member: memberRepo,
		User:   userRepo,
	}, txManager, logger)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService, cfg)
	orgHandler := handlers.NewOrganisationHandler(handlers.OrganisationHandlerServices{
		Org:     orgService,
		Checker: checkerService,
	}, cfg)
	membershipHandler := handlers.NewMembershipHandler(handlers.MembershipHandlerServices{
		Member:       membershipService,
		Organisation: orgService,
		Checker:      checkerService,
	}, cfg)

	// API routes
	api := r.Group("/api/v1")
	authHandler.Routes(api)
	orgHandler.Routes(api)
	membershipHandler.Routes(api)

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server with graceful shutdown
	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server error: %v", err)
		}
	}()
	logger.Infof("Server listening on port %s", cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited cleanly")
}
