package utils

import (
	"math/rand"
	"time"
)

var usedNames = make(map[string]bool)

func RandomPlayerName() string {
	// 隨機名稱生成器的字符集
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// 隨機種子
	rand.Seed(time.Now().UnixNano())

	// 不斷生成名稱，直到找到一個不重複的名稱
	for {
		// 生成隨機名稱，假設名稱長度為 8
		name := make([]byte, 8)
		for i := range name {
			name[i] = charSet[rand.Intn(len(charSet))]
		}

		// 檢查名稱是否已經使用過
		if !usedNames[string(name)] {
			usedNames[string(name)] = true
			return string(name)
		}
	}
}
