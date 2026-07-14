package main

import (
	"qqbot/bot"
	"qqbot/config"
)

func main() {
	bot.Run(config.WsURL)
}
