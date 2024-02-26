package handler

import (
	"math/rand"
	"time"
)

func Random(ran float64) bool {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	val := r.Float64() // 乱数を生成（rand.Rand のメソッド）
	return (val <= ran)
}
