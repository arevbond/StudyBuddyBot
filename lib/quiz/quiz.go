package quiz

type Level int

const (
	Easy Level = iota
	Medium
	Hard
)

type Quiz struct {
	Theme     string
	Questions []Question
	level     Level
}

func (q Quiz) Level() string {
	switch q.level {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	}
	return ""
}

var Quizzes = []Quiz{
	{
		"Разное (тестовый)",
		[]Question{
			{
				"Как зовут \"Отца пингвинов?\"", "", []string{"Линус Торвальдс", "Линус", "Торвальдс"}, 5,
			},
			{
				"Из какого фильма этот кадр?", "https://vlgfilm.ru/upload/resize_cache/iblock/c47/1600_800_1/011CAP4K.jpg",
				[]string{
					"Терминатор 2", "Terminator 2", "Терминатор 2: Судный день", "Терминатор 2 Судный день", "Терминатор два",
				}, 5,
			},
			{"В каком году вышел JavaScript?", "", []string{"1995"}, 5},
			{"Какой язык лучший?", "", []string{"Русский"}, 5},
			{"В каком году был создан телеграмм?", "", []string{"2013"}, 5},
		},
		0,
	},
}
