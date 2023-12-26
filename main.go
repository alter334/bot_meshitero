package main

import (
	"database/sql"
	"errors"
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

	"bot_meshitero/handler"
)

var (
	db *sqlx.DB
)

func main() {
	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: os.Getenv("TRAQ_BOT_TOKEN"), // Required
	})
	if err != nil {
		panic(err)
	}
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("NS_MARIADB_USER") == "" {
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(os.Getenv("NS_MARIADB_USER"))
	fmt.Println("aa")
	conf := mysql.Config{
		User:                 os.Getenv("NS_MARIADB_USER"),
		Passwd:               os.Getenv("NS_MARIADB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("NS_MARIADB_HOSTNAME") + ":" + os.Getenv("NS_MARIADB_PORT"),
		DBName:               os.Getenv("NS_MARIADB_DATABASE"),
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

	h := handler.NewHandler(db, bot)

	bot.OnError(func(message string) {
		log.Println("Received ERROR message: " + message)
	})

	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		log.Println("Received MESSAGE_CREATED event: " + p.Message.Text)
		var user User
		log.Println("A")
		err := db.Get(&user, "SELECT * FROM `users` WHERE `id`=?", p.Message.User.ID)
		log.Println("B")
		//----------------------------------------------------------------
		//ユーザーが見つからなかったらエントリー(usersdb登録)実行
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("C")
			h.Entry(p)
			log.Println("D")
			user.Attack = 0
		} else if err != nil {
			handler.SimplePost(bot, p.Message.ChannelID, "Internal error: "+err.Error())
			return
		}

		//----------------------------------------------------------------
		// 投稿されたチャンネルが過去に投稿されたことのないチャンネルの場合placesdb登録処理を実行
		var exists int
		row := db.QueryRowx("SELECT EXISTS (SELECT * FROM `places` WHERE `channelid`=?)", p.Message.ChannelID)
		log.Println(row)
		row.Scan(&exists)
		log.Println(exists)
		//----------------------------------------------------------------
		//場所が登録されていないとき
		if exists == 0 {
			log.Println("F")
			h.MonitorInsert(p)
		} else if err != nil {
			handler.SimplePost(bot, p.Message.ChannelID, "Internal error: "+err.Error())
			return
		}

		//----------------------------------------------------------------
		//ユーザーが存在した場合コマンド処理に応じて実行
		//コマンドは/区切り
		//画像url取得
		log.Println("E")
		cmd := strings.Fields(p.Message.Text)
		meshiurl := cmd[len(cmd)-1]

		//コマンドなし->通常モード(attackコマンドでも同様)
		switch len(cmd) {
		case 1:
			handler.SimplePost(bot, p.Message.ChannelID, "Input commands or photo")
		default: //現在はコマンド機能は導入していないので
			h.Attack(p, meshiurl, user.Attack)
			if err != nil {
				log.Println(err)
			}
		}

	})

	if err := bot.Start(); err != nil {
		panic(err)
	}

}
