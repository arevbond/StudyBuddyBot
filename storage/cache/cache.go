package cache

import "tg_ics_useful_bot/storage"

type UserCache map[[2]int]*storage.DBUser

func NewUserCache() UserCache {
	return UserCache{}
}

func (c UserCache) GetUser(tgID int, chatID int) (*storage.DBUser, error) {
	user, ok := c[[2]int{tgID, chatID}]
	if !ok {
		return nil, storage.ErrUserNotExist
	}
	return user, nil
}

func (c UserCache) AddUser(user *storage.DBUser) {
	c[[2]int{user.TgID, user.ChatID}] = user
}
