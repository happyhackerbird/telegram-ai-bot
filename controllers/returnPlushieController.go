package controllers

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	teddy = "\U0001F9F8"
	r, _  = utf8.DecodeRuneInString(teddy)
)

func ReturnPlushies(str string, msg *tgbotapi.MessageConfig) {
	var s string
	if strings.Contains(str, "plushies") {
		s = "plushies"
	} else if strings.Contains(str, "plushie") {
		s = "plushie"
	} else if strings.Contains(str, "plushes") {
		s = "plushes"
	} else if strings.Contains(str, "plushe") {
		s = "plushe"
	} else if strings.Contains(str, "plushy") {
		s = "plushy"
	} else if strings.Contains(str, "plush") {
		s = "plush"
	} else {
		msg.Text = "I don't know that command. Here is a plushie: " + string(r)
		return
	}
	count := strings.Count(str, s)
	if count == 1 {
		if s == "plushies" {
			s = "plushie"
		} else if s == "plushes" {
			s = "plushe"
		}
		msg.Text = fmt.Sprintf("Here is one %v: %v", s, strings.Repeat(string(r), count))
	} else {
		msg.Text = fmt.Sprintf("Here are %d %v: %v", count, s, strings.Repeat(string(r), count))
	}
}

func ReturnRandomPlushies(msg *tgbotapi.MessageConfig) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	s := strings.Repeat("plushies", r1.Intn(10)+1)
	log.Println(s)
	ReturnPlushies(s, msg)
}
