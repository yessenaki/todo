package main

import (
	"log"

	"github.com/yesseneon/todo-app"
	"github.com/yesseneon/todo-app/pkg/handler"
	"github.com/yesseneon/todo-app/pkg/repository"
	"github.com/yesseneon/todo-app/pkg/service"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	var srv todo.Server
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server %s", err.Error())
	}
}
