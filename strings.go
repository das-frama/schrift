package main

type Message struct {
	speed    int
	text     string
	autoplay bool
}

var Intro = []Message{
	Message{100, "Привет, Софья.", true},
	Message{100, "Это Андрей.\nТвой друг.", true},
	Message{100, "Мне не сиделось на месте, и поэтому я придумал небольшое откровение для тебя.", false},
	Message{100, "В последнее время у меня затаились некоторые чувства и мысли. И я чувствую, что должен исповедовать оные в ясной и понятной форме.", false},
}

var Test = []Message{
	Message{10, "Тестовое сообщение.", true},
	Message{10, "Второе тестовое сообщение.", true},
}
