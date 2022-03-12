package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/krupyansky/echo-bot/internal"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const BotApi = "https://api.telegram.org/bot"

const BotGetUpdates = "/getUpdates"
const BotSendMessage = "/sendMessage"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		panic("BOT_TOKEN not set or is empty")
	}

	botUrl := BotApi + botToken

	offset := 0
	for {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("getting updates was failed: ", err.Error())
		}

		for _, update := range updates {
			err = respond(botUrl, update)
			offset = update.UpdateId + 1
		}

		fmt.Println(updates)
	}
}

func getUpdates(botUrl string, offset int) ([]internal.Update, error) {
	resp, err := http.Get(botUrl + BotGetUpdates + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var restResponse internal.RestResponse

	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

func respond(botUrl string, update internal.Update) error {
	var botMessage internal.BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = update.Message.Text

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}

	_, err = http.Post(botUrl+BotSendMessage, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	return nil
}
