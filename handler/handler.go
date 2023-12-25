package handler

import (
	"log"
	"strconv"

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
	SimplePost(h.bot, p.Message.ChannelID, ":@"+p.Message.User.Name+":さん\n"+"## ようこそtraP飯テロ部へ\n"+"飯テロ候補先リストに無事登録されました。今の所解除する方法はないです:gomen:")
}

// 通常攻撃:db上に存在するユーザーから1人を選んで爆撃します
func (h *Handler) Attack(p *payload.MessageCreated, meshiurl string, attackNum int) {
	var attackTo, attackName string

	//初の攻撃なら自分に飛ぶ
	if attackNum == 0 {
		attackTo = "97d954a2-695b-466d-9d94-cf4ad88dd262"
		log.Println("InitAttack実行")
	} else {
		//ランダム選択1名
		err := h.db.Get(&attackTo, "SELECT `id` FROM `users` ORDER BY RAND() LIMIT 1")
		if err != nil {
			SimplePost(h.bot, p.Message.ChannelID, "Internal error: "+err.Error())
			log.Println("Internal error: " + err.Error())
			return
		}
		log.Println("Attack実行")

	}

	attackNum++
	attackNumstr := strconv.Itoa(attackNum)
	attackId, attackName := GetUserHome(h.bot, attackTo)
	attackmesid := SimplePost(h.bot, attackId, ":@"+p.Message.User.Name+":"+":oisu-1::oisu-2::oisu-3::oisu-4yoko:"+meshiurl)
	SimplePost(h.bot, p.Message.ChannelID, ":@"+attackName+":"+"に爆撃しました。\n累積攻撃回数:"+attackNumstr+"回\n"+"https://q.trap.jp/messages/"+attackmesid)
	_, err := h.db.Exec("UPDATE `users` SET `attack`=? WHERE `id`=?", attackNum, p.Message.User.ID)
	if err != nil {
		SimplePost(h.bot, p.Message.ChannelID, "Internal error: "+err.Error())
		log.Println("Internal error: " + err.Error())
		return
	}
	log.Println("Attack完了")

}

// テスト:自爆
func (h *Handler) SelfAttack(p *payload.MessageCreated, meshiurl string) {
	log.Println("SelfAttack実行")
	attackId, _ := GetUserHome(h.bot, p.Message.User.ID)
	SimplePost(h.bot, attackId, ":@"+p.Message.User.Name+":"+":oisu-1::oisu-2::oisu-3::oisu-4yoko:"+meshiurl)
	log.Println("SelfAttack完了")
}
