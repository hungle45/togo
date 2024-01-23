package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	httpHandler "togo/app/delivery/http"
	"togo/app/delivery/http/middleware"
	repo "togo/app/repository"
	service "togo/app/service"
	"togo/config"
	"togo/database"
	"togo/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg config.Config) *gin.Engine {
	db := database.NewMySQLDatabase(&cfg)
	conn := db.GetConn()

	if err := db.GetConn().AutoMigrate(&domain.User{}); err != nil {
		log.Fatal(err)
	}
	if err := db.GetConn().AutoMigrate(&domain.Task{}); err != nil {
		log.Fatal(err)
	}
	if err := db.GetConn().AutoMigrate(&domain.TaskManager{}); err != nil {
		log.Fatal(err)
	}

	userRepo := repo.NewUserReposity(conn)
	userService := service.NewUserService(userRepo)
	userService.CreateAdmin(&cfg)

	taskRepo := repo.NewTaskRepository(conn)
	taskService := service.NewTaskService(userRepo, taskRepo)

	r := gin.Default()
	jwtMiddleware := middleware.JWTMiddleware(userService)

	v1 := r.Group("/v1")
	httpHandler.NewUserHTTPHandler(v1, userService)
	httpHandler.NewTaskHTTPHandler(v1, taskService, jwtMiddleware)

	return r
}

func StartServer(cfg config.Config) {
	r := SetupRouter(cfg)

	s := &http.Server{
		Addr:           fmt.Sprintf("%v:%v", cfg.App.Server.Host, cfg.App.Server.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Server running at %v:%v", cfg.App.Server.Host, cfg.App.Server.Port)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
