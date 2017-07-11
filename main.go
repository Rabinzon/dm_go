package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type configType struct {
	Url     string
	Command string
}

const (
	tgBotToken = "410466859:AAE1dS8JpKtKmbV59-qjKrbsUh61gA0AlxY"
)

var (
	port = ""
	path = ""
)

func sendMessage(bot *tgbotapi.BotAPI, text string) {
	msg := tgbotapi.NewMessage(49307595, time.Now().Format("2006-01-02 15:04:05")+" "+text)
	bot.Send(msg)
}

func handler(item configType, bot *tgbotapi.BotAPI) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sendMessage(bot, "🚀 "+item.Url+" is running...")
		err := exec.Command("sh", "-c", item.Command).Run()
		if err != nil {
			sendMessage(bot, "🙊 "+item.Url+" filed")
			log.Print(err)
		}
		sendMessage(bot, "🎉 "+item.Url+" fineshed!")
		w.WriteHeader(http.StatusOK)
	}
}

func createBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(tgBotToken)

	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot
}

func getConf(confPath string) []configType {
	file, err := ioutil.ReadFile(confPath)

	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	config := make([]configType, 0)
	json.Unmarshal(file, &config)

	return config
}

func parse_args() bool {
	flag.StringVar(&port, "p", "", "Порт. Например, 3000")
	flag.StringVar(&path, "c", "", "Путь к конфиг файлу. Например, ./config.json")
	flag.Parse()
	if port == "" {
		fmt.Println("Не задан параметр -p port", port)
		return false
	}
	if path == "" {
		fmt.Println("Не задан параметр -c", path)
		return false
	}
	return true
}

func main() {
	if !parse_args() {
		return
	}

	bot := createBot()
	config := getConf(path)

	for _, item := range config {
		http.HandleFunc("/"+item.Url, handler(item, bot))
	}

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	fmt.Println("PORT:", os.Getenv("PORT"))
}
