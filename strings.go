package main

type Message struct {
	duration int64
	text     string
}

var Intro = []Message{
	Message{100, "Привет, Софья."},
	Message{100, "Это Андрей.\nТвой друг."},
	Message{100, "Мне не сиделось на месте, и поэтому я придумал небольшое откровение для тебя."},
	Message{100, "В последнее время у меня затаились некоторые чувства и мысли. И я чувствую, что должен исповедовать оные в ясной и понятной форме."},
}

var Test = []Message{
	Message{100, "Тестовое сообщение."},
	Message{100, "Второе тестовое сообщение."},
}

// func GetMessage(category string, position int) Message {
// 	switch category {
// 	case "intro":
// 		return Intro[position]
// 	default:
// 		return Message{}
// 	}
// }