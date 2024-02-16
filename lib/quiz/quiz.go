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
	{
		"Программирование",
		[]Question{
			{"Как называется символическое имя, которое используется для хранения данных в программе?", "", []string{"Переменная"}, 30},
			{"Что является основным строительным блоком программ и содержит инструкции для выполнения определенной задачи?", "", []string{"Код"}, 30},
			{"Что представляет собой блок кода, который выполняет определенную задачу и может быть вызван из других частей программы?", "", []string{"Функция"}, 30},
			{"Что представляет собой структура данных, которая хранит набор элементов одного типа?", "", []string{"Массив"}, 30},
			{"Как называется символ или ключевое слово, используемое для выполнения определенной операции над данными?", "", []string{"Оператор"}, 30},
			{"Какая конструкция программирования используется для повторения выполнения определенного блока кода до тех пор, пока выполняется заданное условие?", "", []string{"Цикл"}, 30},
			{"Как называется значение, передаваемое в функцию при ее вызове для использования внутри этой функции?", "", []string{"Параметр"}, 30},
			{"Что представляет собой набор функций или классов, предназначенных для решения определенной задачи в программировании?", "", []string{"Библиотека"}, 30},
			{"Как называется процесс изменения исходного кода программы в целях устранения ошибок?", "", []string{"Дебаггинг", "Дебагинг", "Отладка"}, 30},
			{"Как называется техника программирования, при которой большие задачи разбиваются на более мелкие, более управляемые части?", "", []string{"Декомпозиция"}, 30},
			{"Как называется специальная переменная, которая содержит адрес памяти другой переменной?", "", []string{"Указатель"}, 30},
			{"Какой алгоритм сортировки работает, выбирая на каждом шаге наименьший элемент и перемещая его в начало списка?", "", []string{"Выбором", "Сортировка выбором"}, 30},
			{"Как называется процесс преобразования исходного кода программы в машинный код, который может быть выполнен компьютером?", "", []string{"Компиляция"}, 30},
			{"Какая вторая по популярности программа разработанная Линусом Торвальдсом?", "", []string{"Git", "Гит"}, 30},
			{"Какой термин обозначает структуру данных, представляющую собой набор узлов, каждый из которых содержит ссылку на следующий узел?", "", []string{"Связной список", "Linked list"}, 30},
			{"Это объектно-ориентированный язык программирования, разработанный компанией Sun Microsystems, который используется для создания мобильных и веб-приложений.", "", []string{"Java"}, 30},
			{"Это способ организации кода, при котором действия и данные, которые ими оперируют, объединены в единый объект.", "", []string{"ООП", "Объектно-ориентированное программирование"}, 30},
			{"Это инструмент, используемый для управления версиями исходного кода в проектах разработки программного обеспечения.", "", []string{"Система контроля версий", "Система управления версиями"}, 30},
			{"Как называется метод, который позволяет объектам разных классов иметь одинаковое имя для метода, но различное поведение в зависимости от класса?", "", []string{"Полиморфиз"}, 30},
		},
		Easy,
	},
	{
		"Программирование",
		[]Question{
			{"Как называется концепция программирования, при которой операции выполняются независимо друг от друга и не ожидают завершения других операций?", "", []string{"Асинхронность"}, 60},
			{"Как называется метод разрешения конфликтов при хэшировании, который использует связанные списки для хранения элементов с одинаковым хэшем?", "", []string{"Метод цепочек", "Цепочек"}, 60},
			{"Какая структура данных используется в функциональном программировании для обработки последовательных вычислений с возможностью обработки ошибок и управлением побочными эффектами?", "", []string{"Монада"}, 60},
			{"Как называется паттерн проектирования, который обеспечивает простой интерфейс для взаимодействия с различными частями программы?", "", []string{"Фасад"}, 60},
			{"Как называется процесс внедрения ошибки в программу с целью выявления уязвимостей в системе без реального воздействия на пользователей?", "", []string{"Фаззинг", "Фазинг"}, 60},
			{"Как называется паттерн проектирования, который предоставляет способ создания семейств взаимосвязанных объектов без указания их конкретных классов?", "", []string{"Абстрактная фабрика"}, 60},
			{"Какой алгоритм используется для определения минимального остовного дерева в связном взвешенном графе?", "", []string{"Алгоритм Прима", "Прима"}, 60},
			{"Как называется алгоритм сортировки, который разбивает массив на части, называемые \"кучами\", и последовательно превращает массив в кучу? ", "", []string{"Пирамидальная сортировка", "Сортировка кучей"}, 60},
			{"Как называется техника управления памятью в операционных системах, при которой фрагментированная память объединяется в один блок для повышения эффективности ее использования?", "", []string{"Компактация памяти", "Компактация"}, 60},
			{"Какой шаблон проектирования предлагает разделение системы на множество независимых объектов, что позволяет изменять их поведение без изменения кода клиента?", "", []string{"Стратегия"}, 60},
			{"Как называется алгоритм сжатия данных, который использует кодирование последовательностей символов с использованием дерева кодирования?", "", []string{"Алгоритм Хаффмана", "Хаффмана"}, 60},
			{"Как называется алгоритм сортировки, который разбивает массив на меньшие подмассивы, сортирует их, а затем сливает в один отсортированный массив?", "", []string{"Слиянием", "Сортировка слиянием"}, 60},
			{"Какой алгоритм используется для поиска кратчайшего пути в графе, если веса ребер могут быть отрицательными?", "", []string{"Алгоритм Беллмана-Форда", "Беллмана-Форда", "Беллман-Форд"}, 60},
		},
		Hard,
	},
}
