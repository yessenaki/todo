package main

import (
	"log"

	"github.com/yesseneon/todo-app"
	"github.com/yesseneon/todo-app/pkg/handler"
)

func main() {
	var handlers handler.Handler
	var srv todo.Server
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server %s", err.Error())
	}
}
