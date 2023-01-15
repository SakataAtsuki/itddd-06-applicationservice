package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/SakataAtsuki/itddd-06-applicationservice/domain/model/user"
)

var command = flag.String("usecase", "", "usercase of application")

func main() {
	// DBに接続
	uri := fmt.Sprintf("postgres://%s/%s?sslmode=disable&user=%s&password=%s&port=%s&timezone=Asia/Tokyo",
		os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"))
	db, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	log.Println("successfully connected to database")

	// アプリケーションサービスにリポジトリとドメインサービスを挿入
	userRepository, err := user.NewUserRepository(db)
	if err != nil {
		panic(err)
	}
	userService, err := user.NewUserService(userRepository)
	if err != nil {
		panic(err)
	}
	userApplicationService, err := user.NewUserApplicationService(userRepository, *userService)
	if err != nil {
		panic(err)
	}

	// ユースケースの実行
	flag.Parse()
	log.Println(*command)
	switch *command {
	case "register":
		if err := userApplicationService.Register("test-user"); err != nil {
			log.Println(err)
		}
	case "get":
		userData, err := userApplicationService.Get("test-id")
		if err != nil {
			log.Println(err)
		}
		log.Println(userData)
	case "update":
		userUpdateCommand := &user.UserUpdateCommand{Id: "test-id", Name: "test-updated-user"}
		if err := userApplicationService.Update(*userUpdateCommand); err != nil {
			log.Println(err)
		}
	case "delete":
		userDeleteCommand := &user.UserDeleteCommand{Id: "test-id"}
		if err := userApplicationService.Delete(*userDeleteCommand); err != nil {
			log.Println(err)
		}
	default:
		log.Printf("%s is not command. choose in ('register', 'get', 'update', 'delete')", *command)
	}
}
