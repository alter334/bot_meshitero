package handler

import (
	"log"

	"github.com/jmoiron/sqlx"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

type Handler struct {
	db  *sqlx.DB
	bot *traqwsbot.Bot
}

func NewHandler(db *sqlx.DB, bot *traqwsbot.Bot) *Handler {
	return &Handler{db: db, bot: bot}
}

// エントリー:テロ会員でない場合(db上に存在しなかった場合)はここにきてエントリーメッセージを投稿する
func (h *Handler) Entry(p *payload.MessageCreated) {
	_, err := h.db.Exec("INSERT INTO `users`(`name`,`id`,`attack`,`rate`) VALUES(?,?,0,0)", p.Message.User.Name, p.Message.User.ID)
	if err != nil {
		SimplePost(h.bot, p.Message.ChannelID, "Internal error: "+err.Error())
		log.Println("Internal error: " + err.Error())
		return
	}
	SimplePost(h.bot, p.Message.ChannelID, ":@"+p.Message.User.Name+":さん\n"+"## ようこそtraP飯テロ部へ\n"+"飯テロ候補先リストに無事登録されました。今の所解除する方法はないです:gomen:(不都合ある場合は個人的に連絡ください)")
}

// 通常攻撃:db上に存在するユーザーから1人を選んで爆撃します
func (h *Handler) Attack(p *payload.MessageCreated, meshiurl string, attackNum int) {
	log.Println("Attack実行")
	var attackTo string

	//初の攻撃なら
	if attackNum == 0 {
		log.Println("InitAttack実行")
		SimplePost(h.bot, "402a1c2c-878e-40ef-ae14-011354394e36", ":@"+p.Message.User.Name+":"+"oisu-"+meshiurl)
		log.Println("InitAttack完了")
		return
	}

	//ランダム選択1名
	err := h.db.Get(&attackTo, "SELECT `id` FROM `users` ORDER BY RAND() LIMIT 1")
	if err != nil {
		SimplePost(h.bot, p.Message.ChannelID, "Internal error: "+err.Error())
		log.Println("Internal error: " + err.Error())
		return
	}

	attackId := GetUserHome(h.bot, attackTo)
	SimplePost(h.bot, attackId, ":@"+p.Message.User.Name+":"+"oisu-"+meshiurl)
	log.Println("Attack完了")
	
}

// テスト:自爆
func (h *Handler) SelfAttack(p *payload.MessageCreated, meshiurl string) {
	log.Println("SelfAttack実行")
	attackId := GetUserHome(h.bot, p.Message.User.ID)
	SimplePost(h.bot, attackId, ":@"+p.Message.User.Name+":"+"oisu-"+meshiurl)
	log.Println("SelfAttack完了")
}
