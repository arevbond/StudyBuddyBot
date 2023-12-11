package telegram

import (
	"math/rand"
	"time"
)

const (
	urlImgHink = "https://avatars.mds.yandex.net/get-altay/3518606/2a00000179e2472a99931c431d308fd69e09/XXL"
	urlImgRoom = "https://www.vsu.ru/gallery/photos/study/dept_phys.jpg"
)

func hinkOrRoomCmd() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(2)
	if n == 1 {
		return urlImgHink
	}
	return urlImgRoom
}
