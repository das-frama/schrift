package main

type Message struct {
	delay    int
	autoplay bool
	text     string
}

var Intro = []Message{
	Message{0000, true, "Привет, Софья."},
	Message{2000, true, "Это Андрей.\nТвой друг."},
	Message{2000, true, "Мне не сиделось на месте, и поэтому я придумал небольшое откровение для тебя."},
	Message{2000, true, "В последнее время у меня затаились некоторые мысли, и я чувствую, что должен исповедовать оные в ясной и понятной форме."},
	Message{2000, true, "Этим «актом выпендрёжничества» я хочу тебе кое-что показать."},
	Message{2000, true, "Что-то, что беспокоит меня очень давно на протяжении многих лет."},
	Message{2000, true, "Это «что-то» делает меня собой."},
	Message{2000, true, "Но это «что-то» мне так же пришлось очень хорошо спрятать от твоего внимания."},
	Message{2000, true, "Чтобы найти то, тебе придётся отключить свой разум и рассеять сознание. Так можно стать свидетелем дао."},
	Message{4000, true, "Извини за пафос."},
	Message{2000, true, "Я не расстроюсь, если ты не захочешь читать это дальше или продолжать пытаться искать."},
	Message{2000, true, "Я так же не расстроюсь, если тебе всё надоест, и ты захочешь уйти."},
	Message{2000, true, "Я смогу это понять."},
	Message{4000, true, "Ты готова продолжить? (Д/н):"},
}

var Test = []Message{
	Message{10, true, "Тестовое сообщение."},
	Message{10, true, "Второе тестовое сообщение."},
}
