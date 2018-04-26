package main

type Message struct {
	speed    int
	text     string
	autoplay bool
}

var Intro = []Message{
	Message{100, "Привет, Софья.", false},
	Message{40, "Это Андрей.\nТвой друг.", false},
	Message{10, "Мне не сиделось на месте, и поэтому я придумал небольшое откровение для тебя.", true},
	Message{50, "В последнее время у меня затаились некоторые чувства и мысли. И я чувствую, что должен исповедовать оные в ясной и понятной форме.", false},
}

var Test = []Message{
	Message{50, "Тестовое сообщение.", false},
	Message{50, "Второе тестовое сообщение.", false},
}
