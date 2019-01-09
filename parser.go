package schrift

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

func parseFile(filename string) ([]*Message, error) {
	// Open the file.
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var messages []*Message

	scanner := bufio.NewScanner(file)
	scanner.Split(scanMessages)
	for scanner.Scan() {
		m, err := newMessage(scanner.Text())
		if err != nil {
			log.Fatalln(err)
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func newMessage(s string) (*Message, error) {
	start := strings.IndexByte(s, '[')
	end := strings.IndexByte(s, ']')
	if start < 0 && end >= 0 || start >= 0 && end < 0 {
		return nil, errors.New("ошибка формирования сообщения")
	}

	var message Message

	// Numbers.
	if start+1 != end && start >= 0 && end >= 0 {
		nn := strings.SplitN(s[start+1:end], ";", 2)
		if len(nn) == 1 {
			delay, err := strconv.Atoi(nn[0])
			if err != nil {
				return nil, err
			}
			message.Delay = time.Duration(delay) * time.Millisecond
		} else if len(nn) == 2 {
			speed, err := strconv.Atoi(nn[1])
			if err != nil {
				return nil, err
			}
			message.Speed = time.Duration(speed) * time.Millisecond
		}
	}

	// Questions.
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		if !strings.HasPrefix(l, "- ") {
			continue
		}
		question, err := newQuestion(l)
		if err != nil {
			return nil, err
		}
		message.Questions = append(message.Questions, question)
	}
	start = end + 1
	if end = strings.Index(s, "- "); end == -1 {
		end = len(s)
	}
	message.NewLines = strings.Count(s[start:end], "\n")
	message.Text = strings.TrimSpace(s[start:end])

	return &message, nil
}

func scanMessages(data []byte, atEOF bool) (int, []byte, error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(r) {
			break
		}
	}

	// Check if data starts with '[', then advance to ']'.
	search := 0
	if start == bytes.IndexByte(data, '[') {
		search = bytes.IndexByte(data, ']') + 1
	}
	if i := bytes.IndexByte(data[search:], '['); i >= 0 {
		return search + i, data[start : search+i], nil
	}

	// If we're at EOF, we have a final, non-empty, non-terminated message. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	return 0, nil, nil
}

func newQuestion(s string) (*Question, error) {
	var question Question

	s = strings.TrimPrefix(s, "- ")
	lines := strings.SplitN(s, " : ", 3)
	if len(lines) != 3 {
		return nil, errors.New("ошибка парсинга ответа")
	}
	question.UserAnswers = strings.Split(lines[0], ";")
	switch lines[1] {
	case "continue":
		question.Command = ContinueCommand
	case "quit":
		question.Command = QuitCommand
	}
	question.Answer = strings.TrimSpace(lines[2])

	return &question, nil
}
