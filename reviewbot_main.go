package main

import (
	"log"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	token, ok := os.LookupEnv("TOKEN")
	if !ok {
		log.Fatalln("TOKEN is empty. Please set it in the environment.")
	}

	updates := getUpdatesChan(token)

	revbot := NewReviewBot()

	var wg sync.WaitGroup
	wg.Add(1)

	go receiveUpdates(updates, revbot)

	log.Printf("%s is ready...\n", os.Args[0])

	wg.Wait()
}

func getUpdatesChan(token string) tgbotapi.UpdatesChannel {
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	// Set this to true to log all interactions with telegram servers
	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	return updates
}

func receiveUpdates(updates tgbotapi.UpdatesChannel, revbot *ReviewBot) {
	for {
		// receive update from channel and then handle it
		update := <-updates
		handleUpdate(update, revbot)
	}
}

func handleUpdate(update tgbotapi.Update, revbot *ReviewBot) {
	if update.Message != nil {
		recvMessage(update.Message, revbot)
	}
}

func recvMessage(message *tgbotapi.Message, revbot *ReviewBot) {
	chatId := message.Chat.ID
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	userId := UserId(chatId)
	log.Printf("userId = %d: %s wrote %s", userId, user.FirstName, text)

	response := revbot.HandleText(text, userId, user.FirstName)

	sendMessage(userId, response)
}

func sendMessage(userId UserId, text string) {
	if text == "" {
		return
	}

	chatId := int64(userId)
	mesg := tgbotapi.NewMessage(chatId, text)
	_, err := bot.Send(mesg)
	if err != nil {
		log.Panic(err)
	}
}
