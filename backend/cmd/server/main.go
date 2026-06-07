package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/water-factory/api/config"
	"github.com/water-factory/api/internal/database"
	"github.com/water-factory/api/internal/handler"
	"github.com/water-factory/api/internal/repository"
	"github.com/water-factory/api/internal/service"
	jwtpkg "github.com/water-factory/api/pkg/jwt"
	"github.com/water-factory/api/router"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("load config")
	}

	pool, err := database.NewPool(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("connect database")
	}
	defer pool.Close()

	if err := database.RunMigrations(context.Background(), pool); err != nil {
		log.Fatal().Err(err).Msg("run migrations")
	}

	jwtMgr, err := jwtpkg.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	if err != nil {
		log.Fatal().Err(err).Msg("init jwt manager")
	}

	// Repositories
	userRepo := repository.NewUserRepo(pool)
	mdRepo := repository.NewMasterDataRepo(pool)
	rmInvRepo := repository.NewRMInventoryRepo(pool)
	prodRepo := repository.NewProductionRepo(pool)
	fgInvRepo := repository.NewFGInventoryRepo(pool)
	salesRepo := repository.NewSalesRepo(pool)

	// Services
	authSvc := service.NewAuthService(userRepo, jwtMgr)
	userSvc := service.NewUserService(userRepo)
	mdSvc := service.NewMasterDataService(mdRepo)
	rmInvSvc := service.NewRMInventoryService(rmInvRepo)
	prodSvc := service.NewProductionService(prodRepo)
	fgInvSvc := service.NewFGInventoryService(fgInvRepo)
	salesSvc := service.NewSalesService(salesRepo)
	dashboardSvc := service.NewDashboardService(repository.NewDashboardRepo(pool))

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	userH := handler.NewUserHandler(userSvc)
	mdH := handler.NewMasterDataHandler(mdSvc)
	rmInvH := handler.NewRMInventoryHandler(rmInvSvc)
	prodH := handler.NewProductionHandler(prodSvc)
	fgInvH := handler.NewFGInventoryHandler(fgInvSvc)
	salesH := handler.NewSalesHandler(salesSvc)
	dashboardH := handler.NewDashboardHandler(dashboardSvc)

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	router.Setup(engine, jwtMgr, authH, userH, mdH, rmInvH, prodH, fgInvH, salesH, dashboardH)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Info().Str("addr", addr).Msg("starting server")
	if err := engine.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}
