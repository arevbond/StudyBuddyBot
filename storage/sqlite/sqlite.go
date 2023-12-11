package sqlite

//
//import (
//	"context"
//	"database/sql"
//	"fmt"
//	_ "github.com/mattn/go-sqlite3"
//	"log"
//	"tg_ics_useful_bot/lib/e"
//	"tg_ics_useful_bot/storage"
//	"time"
//)
//
//type Storage struct {
//	db *sql.DB
//}
//
//// New creates new SQLite storage.
//func New(path string) (storage *Storage, err error) {
//	db, err := sql.Open("sqlite3", path)
//	if err != nil {
//		return nil, e.Wrap("[ERROR] can't open db (probably wrong path): ", err)
//	}
//
//	if err := db.Ping(); err != nil {
//		return nil, e.Wrap("[ERROR] can't ping db: ", err)
//	}
//	return &Storage{db: db}, nil
//}
//
//// Init creates tables to storage.
//func (s *Storage) Init(ctx context.Context) error {
//	q1 := `CREATE TABLE IF NOT EXISTS gays (chat_id int, tg_id int, username TEXT, date_last_used DATE)`
//	q2 := `CREATE TABLE IF NOT EXISTS users (tg_id int, chat_id int, is_bot BIT, first_name TEXT, last_name TEXT,
//			username TEXT, is_premium BIT, dick_size INT DEFAULT 0, count_gay_of_day int DEFAULT 0 , last_try_change_dick DATE)`
//	q3 := `CREATE TABLE IF NOT EXISTS calendars (chat_id int UNIQUE, calendar_id TEXT)`
//	q4 := `CREATE TABLE IF NOT EXISTS homeworks (chat_id int, subject TEXT, task TEXT, created_at DATE)`
//	q5 := `CREATE TABLE IF NOT EXISTS user_stats (tg_user_id int, chat_id int, username TEXT, first_name TEXT, last_name TEXT,
//message_count int, dick_plus_count int, dick_minus_count int, yes_count int, no_count int)`
//
//	_, err := s.db.ExecContext(ctx, q1)
//	if err != nil {
//		return e.Wrap("[ERROR] can't create table gays", err)
//	}
//
//	_, err = s.db.ExecContext(ctx, q2)
//	if err != nil {
//		return e.Wrap("[ERROR] can't create table users", err)
//	}
//
//	_, err = s.db.ExecContext(ctx, q3)
//	if err != nil {
//		return e.Wrap("[ERROR] can't create table calendars", err)
//	}
//
//	_, err = s.db.ExecContext(ctx, q4)
//	if err != nil {
//		return e.Wrap("[ERROR] can't create table homeworks", err)
//	}
//
//	_, err = s.db.ExecContext(ctx, q5)
//	if err != nil {
//		return e.Wrap("[ERROR] can't create table user_stats", err)
//	}
//
//	return nil
//}
//
//// CreateUser new user by chatID and telegramID.
//func (s *Storage) CreateUser(ctx context.Context, u *storage.DBUser) error {
//	q := `INSERT INTO users (tg_id, chat_id, is_bot, first_name, last_name, username, is_premium, dick_size, last_try_change_dick)
//							VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
//
//	log.Printf("[INFO] create user #%d '%s' '%s' '%s', chat_id = %d, dick size = %d", u.TgID, u.Username, u.FirstName,
//		u.LastName, u.ChatID, u.DickSize)
//
//	if _, err := s.db.ExecContext(ctx, q, u.TgID, u.ChatID, u.IsBot, u.FirstName,
//		u.LastName, u.Username, u.IsPremium, u.DickSize, u.DateChangeDick); err != nil {
//		return e.Wrap(fmt.Sprintf("can't create user %d %s: ", u.TgID, u.Username), err)
//	}
//	return nil
//}
//
//// User get user by chatID and telegram ID.
//func (s *Storage) GetUser(ctx context.Context, tgID, chatID int) (*storage.DBUser, error) {
//	q := `SELECT * FROM users WHERE tg_id = ? AND chat_id = ?`
//
//	user := &storage.DBUser{}
//
//	err := s.db.QueryRowContext(ctx, q, tgID, chatID).Scan(&user.TgID, &user.ChatID, &user.IsBot, &user.FirstName, &user.LastName,
//		&user.Username, &user.IsPremium, &user.DickSize, &user.CountGayOfDay, &user.DateChangeDick)
//
//	if err == sql.ErrNoRows {
//		return nil, storage.ErrUserNotExist
//	}
//
//	if err != nil {
//		return nil, e.Wrap(fmt.Sprintf("[ERROR] can't get user from storage tg id: %d, chat id: %d", tgID, chatID), err)
//	}
//
//	// log.Printf("from storage get user: tg id = %d, chat id = %d, dick size = %d", user.TgID, user.ChatID, user.DickSize)
//
//	return user, nil
//}
//
//func (s *Storage) UserByUsername(ctx context.Context, username string, chatID int) (*storage.DBUser, error) {
//	q := `SELECT * FROM users WHERE username = ? AND chat_id = ?`
//	user := &storage.DBUser{}
//
//	err := s.db.QueryRowContext(ctx, q, username, chatID).Scan(&user.TgID, &user.ChatID, &user.IsBot, &user.FirstName, &user.LastName,
//		&user.Username, &user.IsPremium, &user.DickSize, &user.CountGayOfDay, &user.DateChangeDick)
//
//	if err == sql.ErrNoRows {
//		return nil, storage.ErrUserNotExist
//	}
//
//	if err != nil {
//		return nil, e.Wrap(fmt.Sprintf("[ERROR] can't get user from storage username: %s, chat id: %d", username, chatID), err)
//	}
//
//	// log.Printf("from storage get user: tg id = %d, chat id = %d, dick size = %d", user.TgID, user.ChatID, user.DickSize)
//
//	return user, nil
//}
//
//func (s *Storage) UpdateUserDickSize(ctx context.Context, u *storage.DBUser, dickSize int) error {
//	q := `UPDATE users SET dick_size = ? WHERE tg_id = ? AND chat_id = ?`
//	oldDickSize := u.DickSize
//	if _, err := s.db.ExecContext(ctx, q, dickSize, u.TgID, u.ChatID); err != nil {
//		return e.Wrap(fmt.Sprintf("[ERROR] can't update dick size user %d chat id %d from %d to %d",
//			u.TgID, u.ChatID, u.DickSize, dickSize), err)
//	}
//	u.DickSize = dickSize
//	log.Printf("[INFO] user %d %s change his dick from %d to %d", u.TgID, u.Username, oldDickSize, u.DickSize)
//	return nil
//}
//
//func (s *Storage) UpdateDateLastTryChangeDickToNow(ctx context.Context, u *storage.DBUser) error {
//	q := `UPDATE users SET last_try_change_dick = ? WHERE tg_id = ? AND chat_id = ?`
//	currentTime := time.Now()
//	if _, err := s.db.ExecContext(ctx, q, currentTime, u.TgID, u.ChatID); err != nil {
//		return e.Wrap(fmt.Sprintf("[ERROR] can't update date last try change dick to now user %d chat id",
//			u.TgID, u.ChatID), err)
//	}
//	log.Printf("[INFO] user #%d %s change his date last try change dick to %s", u.TgID, u.Username, currentTime.Format("02-Jan-2006"))
//	return nil
//}
//
//func (s *Storage) UsersByChat(ctx context.Context, chatID int) ([]*storage.DBUser, error) {
//	q := `SELECT * FROM users WHERE chat_id = ? ORDER BY -dick_size`
//	rows, err := s.db.QueryContext(ctx, q, chatID)
//	if err != nil {
//		return nil, e.Wrap(fmt.Sprintf("can't get users by chat id: %s", chatID), err)
//	}
//	defer rows.Close()
//
//	var users []*storage.DBUser
//
//	for rows.Next() {
//		user := &storage.DBUser{}
//		if err := rows.Scan(&user.TgID, &user.ChatID, &user.IsBot, &user.FirstName, &user.LastName,
//			&user.Username, &user.IsPremium, &user.DickSize, &user.CountGayOfDay, &user.DateChangeDick); err != nil {
//			return users, e.Wrap(fmt.Sprintf("can't get users by chat id: %s", chatID), err)
//		}
//		users = append(users, user)
//	}
//	if err = rows.Err(); err != nil {
//		return users, err
//	}
//	return users, nil
//}
//
//func (s *Storage) GetGayOfDay(ctx context.Context, chatID int) (*storage.DBGay, error) {
//	q := `SELECT * FROM gays WHERE chat_id = ?`
//
//	gay := &storage.DBGay{}
//
//	err := s.db.QueryRowContext(ctx, q, chatID).Scan(&gay.ChatID, &gay.TgID, &gay.Username, &gay.CreatedAt)
//
//	if err == sql.ErrNoRows {
//		return nil, storage.ErrUserNotExist
//	}
//
//	if err != nil {
//		return nil, e.Wrap(fmt.Sprintf("[ERROR] can't get gay from table gays chat id: %d", chatID), err)
//	}
//
//	// log.Printf("from storage get user: tg id = %d, chat id = %d, dick size = %d", user.TgID, user.ChatID, user.DickSize)
//
//	return gay, nil
//}
//
//func (s *Storage) CreateGayOfDay(ctx context.Context, gay *storage.DBGay) error {
//	q := `INSERT INTO gays (chat_id, tg_id, username, date_last_used)
//							VALUES (?, ?, ?, ?)`
//
//	log.Printf("[INFO] create gay of day #%d '%s', chat_id = %d", gay.TgID, gay.Username, gay.ChatID)
//
//	if _, err := s.db.ExecContext(ctx, q, gay.ChatID, gay.TgID, gay.Username, gay.CreatedAt); err != nil {
//		return e.Wrap(fmt.Sprintf("can't create gay %d %s: ", gay.TgID, gay.Username), err)
//	}
//	return nil
//}
//
//func (s *Storage) RemoveGayOfDay(ctx context.Context, chatID int) error {
//	q := `DELETE FROM gays WHERE chat_id = ?`
//
//	if _, err := s.db.ExecContext(ctx, q, chatID); err != nil {
//		return e.Wrap(fmt.Sprintf("[ERROR] can't remove gay %d %s: ", chatID), err)
//	}
//	return nil
//}
//
//func (s *Storage) IncreaseCountOfGay(ctx context.Context, u *storage.DBUser) error {
//	q := `UPDATE users SET count_gay_of_day = ? WHERE tg_id = ? AND chat_id = ?`
//	oldCount := u.CountGayOfDay
//	if _, err := s.db.ExecContext(ctx, q, oldCount+1, u.TgID, u.ChatID); err != nil {
//		return e.Wrap(fmt.Sprintf("[ERROR] can't update count gay of day user %d chat id %d",
//			u.TgID, u.ChatID), err)
//	}
//	u.CountGayOfDay += 1
//	return nil
//}
//
//func (s *Storage) GetCalendarID(ctx context.Context, chatID int) (string, error) {
//	q := `SELECT * from calendars WHERE chat_id = ?`
//	var id int
//	var calendarID string
//	err := s.db.QueryRowContext(ctx, q, chatID).Scan(&id, &calendarID)
//	if err == sql.ErrNoRows {
//		return "", storage.ErrUserNotExist
//	}
//
//	if err != nil {
//		return "", e.Wrap(fmt.Sprintf("[ERROR] can't get calendar_id from table calendars chat id: %d", chatID), err)
//	}
//	return calendarID, nil
//}
//
//func (s *Storage) AddCalendarID(ctx context.Context, chatID int, calendarID string) error {
//	q := `INSERT INTO calendars (chat_id, calendar_id) VALUES (?, ?) ON CONFLICT (chat_id) DO UPDATE SET calendar_id = ?`
//	if _, err := s.db.ExecContext(ctx, q, chatID, calendarID, calendarID); err != nil {
//		return e.Wrap(fmt.Sprintf("can't update or create calendar_id in chat #%d: ", chatID), err)
//	}
//	return nil
//}
//
//func (s *Storage) addHomeworkCmd(ctx context.Context, chatID int, subject string, task string) error {
//	q := `INSERT INTO homeworks (chat_id, subject, task, created_at) VALUES (?, ?, ?, ?)`
//	if _, err := s.db.ExecContext(ctx, q, chatID, subject, task, time.Now()); err != nil {
//		return e.Wrap("can't add homework:", err)
//	}
//	return nil
//}
//
//func (s *Storage) GetHomeworkByChatID(ctx context.Context, chatID int, limit int) ([]*storage.DBHomework, error) {
//	q := `SELECT rowid, *  from homeworks WHERE chat_id = ? ORDER BY -created_at LIMIT ?`
//
//	rows, err := s.db.QueryContext(ctx, q, chatID, limit)
//	if err != nil {
//		return nil, e.Wrap(fmt.Sprintf("can't get homeworks by chat id: %s", chatID), err)
//	}
//	defer rows.Close()
//
//	var homeworks []*storage.DBHomework
//
//	for rows.Next() {
//		homework := &storage.DBHomework{}
//		if err := rows.Scan(&homework.ID, &homework.ChatID, &homework.Subject, &homework.Task, &homework.CreatedAT); err != nil {
//			return nil, e.Wrap(fmt.Sprintf("can't get homeworks by chat id: %s", chatID), err)
//		}
//		homeworks = append(homeworks, homework)
//	}
//	if err = rows.Err(); err != nil {
//		return homeworks, err
//	}
//	return homeworks, nil
//}
//
//func (s *Storage) GetHomeworkBySubject(ctx context.Context, chatID int, subject string) ([]*storage.DBHomework, error) {
//	q := `SELECT rowid, * from homeworks WHERE subject = ? ORDER BY -created_at`
//
//rows, err := s.db.QueryContext(ctx, q, subject)
//if err != nil {
//	return nil, e.Wrap(fmt.Sprintf("can't get homeworks by chat id: %s", chatID), err)
//}
//defer rows.Close()
//
//var homeworks []*storage.DBHomework
//
//for rows.Next() {
//	homework := &storage.DBHomework{}
//	if err := rows.Scan(&homework.ID, &homework.ChatID, &homework.Subject, &homework.Task, &homework.CreatedAT); err != nil {
//		return nil, e.Wrap(fmt.Sprintf("can't get homeworks by subject: %d", chatID), err)
//	}
//	homeworks = append(homeworks, homework)
//}
//if err = rows.Err(); err != nil {
//	return homeworks, err
//}
//return homeworks, nil
//}
//
//func (s *Storage) DeleteHomework(ctx context.Context, rowID int) error {
//	q := `DELETE FROM homeworks WHERE rowid = ?`
//	_, err := s.db.ExecContext(ctx, q, rowID)
//	if err != nil {
//		return e.Wrap("can't delete row:", err)
//	}
//	return nil
//}
//
//func (s *Storage) CreateUserStats(ctx context.Context, u *storage.DBUserStat) error {
//	q := `INSERT INTO user_stats (tg_user_id, chat_id, username, first_name, last_name, message_count, dick_plus_count,
//                        dick_minus_count, yes_count, no_count)
//							VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
//
//	log.Printf("[INFO] create user_stats #%d chat_id = %d", u.TelegramID, u.ChatID)
//
//	if _, err := s.db.ExecContext(ctx, q, u.TelegramID, u.ChatID, u.UserName, u.FirstName, u.LastName, u.MessageCount, u.DickPlusCount,
//		u.DickMinusCount, u.YesCount, u.NoCount); err != nil {
//		return e.Wrap(fmt.Sprintf("can't create user stats %d %s: ", u.TelegramID, u.UserName), err)
//	}
//	return nil
//}
//
//func (s *Storage) UserStatsByTelegramIDAndChatID(ctx context.Context, tgID, chatID int) (*storage.DBUserStat, error) {
//	q := `SELECT * FROM user_stats WHERE tg_user_id = ? AND chat_id = ?`
//
//	user := &storage.DBUserStat{}
//
//	err := s.db.QueryRowContext(ctx, q, tgID, chatID).Scan(&user.TelegramID, &user.ChatID, &user.UserName, &user.FirstName, &user.LastName,
//		&user.MessageCount, &user.DickPlusCount, &user.DickMinusCount, &user.YesCount, &user.NoCount)
//
//	if err == sql.ErrNoRows {
//		return nil, storage.ErrUserNotExist
//	}
//
//	if err != nil {
//		return nil, e.Wrap(fmt.Sprintf("[ERROR] can't get user stats from storage tg id: %d, chat id: %d", tgID, chatID), err)
//	}
//	return user, nil
//}
//
//func (s *Storage) IncreaseMessageCount(ctx context.Context, u *storage.DBUserStat) error {
//	q := `UPDATE user_stats SET message_count = ? WHERE tg_user_id = ? AND chat_id = ?`
//	oldCount := u.MessageCount
//	if _, err := s.db.ExecContext(ctx, q, oldCount+1, u.TelegramID, u.ChatID); err != nil {
//		return e.Wrap(fmt.Sprintf("[ERROR] can't update message count of user %d chat id %d",
//			u.TelegramID, u.ChatID), err)
//	}
//	u.MessageCount += 1
//	return nil
//}
//
//func (s *Storage) IncreaseDickPlusCount(ctx context.Context, u *storage.DBUserStat) error {
//	q := `UPDATE user_stats SET dick_plus_count = ? WHERE tg_user_id = ? AND chat_id = ?`
//	oldCount := u.DickPlusCount
//	if _, err := s.db.ExecContext(ctx, q, oldCount+1, u.TelegramID, u.ChatID); err != nil {
//		return e.Wrap(fmt.Sprintf("[ERROR] can't update dick_plus_count of user %d chat id %d",
//			u.TelegramID, u.ChatID), err)
//	}
//	u.DickPlusCount += 1
//	return nil
//}
//
//func (s *Storage) IncreaseDickMinusCount(ctx context.Context, u *storage.DBUserStat) error {
//	q := `UPDATE user_stats SET dick_minus_count = ? WHERE tg_user_id = ? AND chat_id = ?`
//	oldCount := u.DickMinusCount
//	if _, err := s.db.ExecContext(ctx, q, oldCount+1, u.TelegramID, u.ChatID); err != nil {
//		return e.Wrap(fmt.Sprintf("[ERROR] can't update dick_minus_count of user %d chat id %d",
//			u.TelegramID, u.ChatID), err)
//	}
//	u.DickMinusCount += 1
//	return nil
//}
//
//func (s *Storage) IncreaseYesCount(ctx context.Context, u *storage.DBUserStat) error {
//	q := `UPDATE user_stats SET yes_count = ? WHERE tg_user_id = ? AND chat_id = ?`
//	oldCount := u.YesCount
//	if _, err := s.db.ExecContext(ctx, q, oldCount+1, u.TelegramID, u.ChatID); err != nil {
//		return e.Wrap(fmt.Sprintf("[ERROR] can't update yes_count of user %d chat id %d",
//			u.TelegramID, u.ChatID), err)
//	}
//	u.YesCount += 1
//	return nil
//}
//
//func (s *Storage) IncreaseNoCount(ctx context.Context, u *storage.DBUserStat) error {
//	q := `UPDATE user_stats SET no_count = ? WHERE tg_user_id = ? AND chat_id = ?`
//	oldCount := u.NoCount
//	if _, err := s.db.ExecContext(ctx, q, oldCount+1, u.TelegramID, u.ChatID); err != nil {
//		return e.Wrap(fmt.Sprintf("[ERROR] can't update no_count of user %d chat id %d",
//			u.TelegramID, u.ChatID), err)
//	}
//	u.NoCount += 1
//	return nil
//}
//
//func (s *Storage) AllUsersStatsInChat(ctx context.Context, chatID int) ([]*storage.DBUserStat, error) {
//	q := `SELECT * FROM user_stats WHERE chat_id = ?`
//
//	rows, err := s.db.QueryContext(ctx, q, chatID)
//
//	if err != nil {
//		return nil, e.Wrap(fmt.Sprintf("can't get users stats by chat id: %s", chatID), err)
//	}
//	defer rows.Close()
//
//	var users []*storage.DBUserStat
//
//	for rows.Next() {
//		user := &storage.DBUserStat{}
//		if err := rows.Scan(&user.TelegramID, &user.ChatID, &user.UserName, &user.FirstName,
//			&user.LastName, &user.MessageCount, &user.DickPlusCount, &user.DickMinusCount, &user.YesCount,
//			&user.NoCount); err != nil {
//			return users, e.Wrap(fmt.Sprintf("can't get users stats by chat id: %s", chatID), err)
//		}
//		users = append(users, user)
//	}
//	if err = rows.Err(); err != nil {
//		return users, err
//	}
//	return users, nil
//}
