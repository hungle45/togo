package repository_test

import (
	"os"
	"testing"
	"togo/config"
	"togo/database"
	"togo/domain"
	"togo/app/repository"
)

var userRepository domain.UserRepository
var taskRepository domain.TaskReponsitory

func TestMain(m *testing.M) {
	cfg := config.LoadConfig("../../config.yml")

	db := database.NewMySQLDatabase(&cfg)
	conn := db.GetConn()

	db.GetConn().AutoMigrate(&domain.User{})

	userRepository = repository.NewUserReposity(conn)
	taskRepository = repository.NewTaskRepository(conn)

	os.Exit(m.Run())
}
