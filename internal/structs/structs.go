package structs

import(
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type Bot struct{
	BotAPI *tgbotapi.BotAPI
	Token string
	UpdateConfig tgbotapi.UpdateConfig
	Updates tgbotapi.UpdatesChannel
}


type Database struct{
	ClientOptions *options.ClientOptions
	Client *mongo.Client
	DBMongo *mongo.Database
	CollectionDB *mongo.Collection
}