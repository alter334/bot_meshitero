package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"

	//"bot_meshitero/handler"
)

var (
	db *sqlx.DB
)

func main() {
	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: os.Getenv("ACCESS_TOKEN"), // Required
	})
	if err != nil {
		panic(err)
	}
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("MARIADB_USER") == "" {
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(os.Getenv("MARIADB_USER"))
	fmt.Println("aa")
	conf := mysql.Config{
		User:                 os.Getenv("MARIADB_USER"),
		Passwd:               os.Getenv("MARIADB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MARIADB_HOSTNAME") + ":" + os.Getenv("MARIADB_PORT"),
		DBName:               os.Getenv("MARIADB_DATABASE"),
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  jst,
		AllowNativePasswords: true,
	}

	_db, err := sqlx.Open("mysql", conf.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("conntected")

	db = _db

	bot.OnError(func(message string) {
		log.Println("Received ERROR message: " + message)
	})

	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		log.Println("Received MESSAGE_CREATED event: " + p.Message.Text)
		cmd := strings.Split(p.Message.Text, " ")
		log.Println(cmd[1])
		//handler.SimplePost(bot, p.Message.ChannelID, "oisu- https://q.trap.jp/files/645e3498-e982-4c59-ba7c-576f6f1baad9")
		if err != nil {
			log.Println(err)
		}
	})

	if err := bot.Start(); err != nil {
		panic(err)
	}

}
