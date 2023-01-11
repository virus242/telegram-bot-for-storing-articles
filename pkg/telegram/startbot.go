package telegram

import (
	"flag"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/virus242/telegram-bot-for-storing-articles/internal"
	"github.com/virus242/telegram-bot-for-storing-articles/internal/structs"
	"github.com/virus242/telegram-bot-for-storing-articles/pkg/database"
)

var bot structs.Bot
var currentArticleText = map[int64]string{}
var msg tgbotapi.MessageConfig

var numbericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Delete", "Delete"),
		tgbotapi.NewInlineKeyboardButtonData("Leave", "Leave"),
	),
)


func getTOKENFromArg()string{
	tokenPrt := flag.String("TOKEN", "", "TOKEN for your telegram bot")
	flag.Parse()
	return *tokenPrt
}


func createTgBot(token string) *tgbotapi.BotAPI{
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil{
		panic(err)
	}
	return bot
}


func sendMessage(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig){
	if _, err := bot.Send(msg); err != nil{
		panic(err)
	}
}


func InitTgBot(){
	bot.Token = getTOKENFromArg()
	bot.BotAPI = createTgBot(bot.Token)
	bot.BotAPI.Debug = false
	bot.UpdateConfig = tgbotapi.NewUpdate(1)
	bot.UpdateConfig.Timeout = 30 
	bot.Updates = bot.BotAPI.GetUpdatesChan(bot.UpdateConfig)
}


func isCallbackQuety(update tgbotapi.Update)bool{
	if update.CallbackQuery != nil{
		return true
	}
	return false
}


func getButtonFromQuety(update tgbotapi.Update)string{
	if isCallbackQuety(update){
		return update.CallbackQuery.Data
	}
	return "None"
}


func putButtonDelete(update tgbotapi.Update){
	database.DeleteArticle(currentArticleText[update.CallbackQuery.From.ID])
	msg = tgbotapi.NewMessage(update.CallbackQuery.From.ID, "The link has been removed from your collection")
	sendMessage(bot.BotAPI, msg)
}

func processingButtonQuery(update tgbotapi.Update, msg tgbotapi.MessageConfig)bool{
	switch getButtonFromQuety(update){
	case "Delete":
		putButtonDelete(update)
		return true
	default:
		return false
	}
}


func isMessageEmpty(update tgbotapi.Update)bool{
	if update.Message == nil{
		return true
	}
	return false
}


func reqCommandGetRandArticle(update tgbotapi.Update){
	var ok bool
	currentArticleText[update.Message.Chat.ID], ok = database.GetRandomArticle(update.Message.Chat.ID)
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, currentArticleText[update.Message.Chat.ID])
	if ok{
		msg.ReplyMarkup = numbericKeyboard
	}
}


func reqCommandHelpOrStart(update tgbotapi.Update){
	text := `Bot is designed to save articles that have interested you. 
	To save the article just send the link to the bot. 
	To get the link type /getRandomArticle`
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
}


func processingCommands(update tgbotapi.Update){
	
	switch update.Message.Command(){

	case "getRandomArticle":
		reqCommandGetRandArticle(update)

	case "help", "start":
		reqCommandHelpOrStart(update)
	}
}


func getUpdates(){

	for update := range bot.Updates{
		if processingButtonQuery(update, msg){
			continue
		}

		if isMessageEmpty(update){
			continue
		}

		if update.Message.IsCommand(){
			processingCommands(update)
		} else if internal.CheckURL(update.Message.Text){
			
			go database.CreateNewArticleToDB(update.Message.Text, update.Message.Chat.ID)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Link is save!")
		}

		sendMessage(bot.BotAPI, msg)
	}
}


func StartBot(){
	InitTgBot()
	getUpdates()
}