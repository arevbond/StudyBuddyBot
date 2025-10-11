package telegram

const msgHelp = `*Доступные команды:*

/add - добавить домашнее задание 📖
/cancel - отменить добавление домашнего задания
/get [number] [subject] - может вызываться без параметров, тогда выведет последние 5 добавленных записей, либо с одним из параметров, number - число последних записей, subject - название предмета
/delete id - удалить запись по id

/schedule - получить расписание из Google Calendar (_рабоает только если привязан calendar-id группы_)
/add\_calendar *[calendar-id]* - привязать расписание из Google Calendar (_возможно только для админов группы_)
(*ВАЖНО*  - _не забудьте в настройках календаря открыть доступ пользователю: calendar-manager@flash-spark-404006.iam.gserviceaccount.com_)

/dick - узнай всё про свой 🍌
/duel _@username_ - вызвать на (или принять) бой ⚔️
/top\_dick - статистика всех 🍆

/auf - зарядиться уверенностью на весь день 😎

/gay - узнать у кого сегодня удачный день 🤡 (_работает, только на админов чата_)
/top\_gay - статистика по бедолагам в чате 🔞

/start\_auction - заустить аукцион!
/deposit *[amount]* - если аук. запущен, то данной командой можно в нём поучавствовать, amount - обязательный параметр. 
/auction - если аукцион запущен, покажет список его участников.

/flip - подбросить монетку 🪙
/xkcd - случайный xkcd комикс 😂
`

const (
	msgCreateUserWithUsername     = "@%s, только что обнаружил(а) свой пенис 🤣\n"
	msgCreateUserWithFullName     = "%s, только что обнаружил(а) свой пенис 🤣\n"
	msgAlreadyPlaysWithUsername   = "@%s, сегодня все твои попытки закончились 🚨\n"
	msgAlreadyPlaysWithFullName   = "%s, сегодня все твои попытки закончились 🚨\n"
	msgDickSize                   = "Теперь размер его пениса: %d см 🍌"
	msgDickIncrease               = "@%s, твой пенис увеличился на %d см 😍\n"
	msgDickDecrease               = "@%s, твой пенис уменьшился на %d см 😭\n"
	msgChangeDickSizeWithUsername = "@%s, %d ➜ %d см 🍌"
	msgChangeDickSizeWithFullName = "%s, %d ➜ %d см 🍌"

	msgTargetNotFound = "@%s этот пользователь не имеет писюна 🍆"
	msgVictoryInDuel  = "@%s победил в дуели @%s\n"
	msgUserHasBanned  = "@%s получает бан на %d секунд 🚫\n"

	msgChanceDuel = "@%s имеет пенис %d см и шансы на победу %.2f%%\n@%s имеет пенис %d см и шансы на победу %.2f%%\n"

	msgNewGayOfDayUsername     = "Новый пидор дня - @%s"
	msgCurrentGayOfDayUsername = "Текущий пидор дня - @%s"

	msgNewGayOfDayFullName     = "Новый пидор дня - %s"
	msgCurrentGayOfDayFullName = "Текущий пидор дня - %s"

	msgCalendarNotExists       = "К вашей группе не привязан Google Calendar\nЧтобы приязвать Google Calendar воспользуйтесь командой /add_calendar"
	msgForbiddenCalendarUpdate = "Обновить id календаря может только администратор группы"
	msgErrorSendMessage        = "Не удалось получить расписание из календаря с id: %s"
	msgErrorUpdateCalendarID   = "Не удалось добавить данное id календаря: \"%s\""
	msgSuccessUpdateCalendarID = "Теперь к вашему чату привязано новое расписание\nНе забудь открыть доступ к своему календарю пользователю:\ncalendar-manager@flash-spark-404006.iam.gserviceaccount.com\n"

	msgSuccess                  = "Успешно"
	msgError                    = "Ошибка"
	msgErrorAdminChangeDickSize = "Не удалось поменять значение пениса данного пользователя"

	msgHomeworkCanceled       = "Галя, у нас отмена!"
	msgHomeworkWithoutSubject = "Пожалуйста после команды укажите название предмета в формате:\n/add_homework #ЗащитаИнформации Лабораторная 7..."
	msgHomeworkWithoutData    = "Пожалуйста после команды и названия предмета укажите само задание в формате:\n/add_homework #ФизическаяКультура Задали пробежать 100 км на выходных"
	msgHomeworkSuccessAdded   = "ДЗ: %s - %s\n Успешно добавлено"

	msgErrorZeroUsersInSeason = "Не удалось заввершить сезон, т.к. было 0 участников"

	msgEndSeason = "🏆🎉 *Сезон завершён!* Поздравляем всех участников! 🏁🥳\n\n"

	msgSeasonResult = "Результаты:\n%s"

	msgStartSeason = "_Победитель получает +1 крутку в день!_\n🚀🌟 Начинаем *новый сезон игр*! 🎮🔥 Готовы к вызовам и победам? 📅⏳ Присоединяйтесь к нам! 🥳💫\n"

	msgUserStats = "Количество сообщений в данном чате: %d\nКоличество прибавлений к пенису: %d\nКоличество уменьшений пениса: %d\nПизда моментов: %d\nПидора ответ моментов: %d\n" +
		"Всего дуелей: %d\nВыиграно дуелей: %d\nПроиграно дуелей: %d\nУбийств в дуелях: %d\nСмертей в дуелях: %d\n"

	// DUEL
	msgDuelWithYourself = "@%s засунул пенис себе в рот 🍆"

	msgChallengeToDuel = "@%s вызывает на дуель @%s"

	msgAcceptDuel = "@%s %s %d см 🍌 %.2f%%\n⚔️\n@%s %s %d см 🍌 %.2f%%\n"

	// msgFinishDuel = "\n\n🤼‍♂️🤼‍♂️🤼‍♂️🤼‍♂️🤼‍♂️🤼‍♂️🤼‍♂️🤼‍♂️🤼‍♂️\n\n🏆 @%s %s %d см ➕ %d \n🤕@%s %s %d см ➖ %d"
	msgFinishDuel = "\n\n🤼‍♂️🤼‍♂️🤼‍♂️🤼‍♂️🤼‍♂️🤼‍♂️🤼‍♂️🤼‍♂️🤼‍♂️\n\n🏆 @%s %s %d см ➕ %d \n🤕@%s %s %d см"

	// msgPlayerDie = "\n\nεつ▄█▀█ ●\n\n🏆 @%s %s %d см ➕ %d \n🤕@%s %s %d см ➖ %d"
	msgPlayerDie = "\n\nεつ▄█▀█ ●\n\n🏆 @%s %s %d см ➕ %d \n🤕@%s %s %d см"

	msgCantCreateDuel = "Невозможно создать дуель\n@%s - %s \n@%s - %s"

	// HP
	msgCantGetHP = "@%s %s - твоё хп."
	msgGetHp     = "@%s пополнил HP 🥰\nТекущее здоровье  %s"

	// HOMEWORK
	msgAddSubject = "Введите название предмета"
	msgAddTask    = "Введите задание"
	msgSomethingWrong
	msgSuccessDelete    = "Запись №%d успешно удалена"
	msgErrorDelete      = "Не удалось удалить запись №%d"
	msgIncorrectValue   = "%s - некоректное значение id"
	msgErrorAddHomework = "Не удалось добавить задание"

	// auction
	msgStartAuction = `Итак дорогие друзья!
Объявляю вашему внимаю, что запускается *аукцион*!

Чтобы учавстовать ставь на кон часть совего пениса и увеличивай шансы победы в аукционе!
Максимальная ставка: %d см.

_Команда для участия:_
/deposit _{amount}_ - amount является обязательным параметром! И не должен превышать размеры вашего члена!

*УДАЧИ!!!*`
	msgAuctionIsStarted  = "В данном чате уже запущен аукцион!"
	msgErrorDeposit      = "Столько вашего пениса в аукцион не влезет..."
	msgAuctionNotStarted = "Аукцион пока что не запущен."
	msgSuccessDeposit    = "Вы успешно внесли в аукцион %d см своего пениса!"
	msgErrorDepositCmd   = "Чтобы внести депозит в аукцион введите\n/deposit {amount} - amount обязательный аргумент, не превышающий размер вашего пениса"
	msgNotEnoughPlayers  = "Увы, в аукционе никто не участвовал..."
	msgZeroPlayers       = "На данный момент никто не учавствует в аукционе\nКоманда для участия:\n/deposit {amount} - amount является обязательным параметром! И не должен превышать размеры вашего члена!"

	msgWinner = "@%s побеждает в аукционе!\nИ прибавляет %d см к своему пенису!"

	// Holiday
	msgAlreadyUsedHoliday = "Вы уже испытали удачу!"
	msgWinHoliday         = "Поздравляю с Новым Годом, друг!\nЖелаю чтобы и в жизни твой пенис увеличился в два раза!"
	msgLoseHoliday        = "Поздравляю с Новым Годом, друг!\nТвой пенис увеличивается в 0.5 раз!!!!"

	// Quiz
	msgErrorQuiz        = "Не удаётся запустить quit из файла: %s"
	msgWriteQuizName    = "Введите название quit"
	msgCorrectAnswer    = "👍"
	msgStartQuiz        = "*Запускается Квиз!*\nТема квиза: _%s_\nУровень сложности: %s\nКоличество вопросов: %d\n_За каждый правильный ответ начисляются сантиметры!_"
	msgFinishQuiz       = "Квиз окончен!"
	msgStoppedQuiz      = "Квиз остановлен!"
	msgQuizNotAvailable = "Квиз не запущен!"
)
