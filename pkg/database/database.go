package database

import (
	"context"
	"crypto/rand"
	"math/big"
	"fmt"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/virus242/telegram-bot-for-storing-articles/pkg/config"
	"github.com/virus242/telegram-bot-for-storing-articles/internal/structs"
)



var ctx = context.TODO()
var DB = structs.Database{}


func showResDeletedArticle(nameArticle string){
	fmt.Printf("delete: %s\n", nameArticle)
}


func deleteOneArticle(ctx context.Context, nameArticle string){
	_, err := DB.CollectionDB.DeleteOne(ctx, bson.M{"article": nameArticle})
	if err != nil{
		panic(err)
	}

	showResDeletedArticle(nameArticle)
}


func DeleteArticle(nameArticle string){
	startDB("WhatMeReadDB", "users")
	defer DB.Client.Disconnect(ctx)
	deleteOneArticle(ctx, nameArticle)
}


func createSearchFilterForChatID(ChatID int64)primitive.D{
	return bson.D{{Key: "chatID", Value: ChatID}}
}


func setCursorByArticle(ctx context.Context, chatID int64) *mongo.Cursor{
	cursor, err := DB.CollectionDB.Find(ctx, createSearchFilterForChatID(chatID))
	if err != nil{
		panic(err)
	}
	return cursor
}


func getChooseArticles(cursor *mongo.Cursor) []bson.M{
	var articlesbson []bson.M
	if err := cursor.All(ctx, &articlesbson); err != nil{
		panic(err)
	}
	return articlesbson
}


func findArticle(ctx context.Context, chatID int64)[]bson.M{
	return getChooseArticles(setCursorByArticle(ctx,chatID))
}


func getRundomNumBelowMax(Max int64)int64{

	if Max > 0{
		randNum, err := rand.Int(rand.Reader, big.NewInt(Max))
		if err != nil {
			panic(err)
		}
		return randNum.Int64()
	}
	return 0
}


func findArticleByIndexInArr(index int64, articles []primitive.M)(string, bool){
	for i, e := range articles{
		if i == int(index){
			return fmt.Sprintf("%s", e["article"]), true
		}
	}
	return "Your article collection is empty", false
}


func GetRandomArticle(chatID int64)(string, bool){
	startDB("WhatMeReadDB", "users")
	defer DB.Client.Disconnect(ctx)
	
	articlesbson := findArticle(ctx, chatID)
	var randNum int64 = getRundomNumBelowMax(int64(len(articlesbson)))

	return findArticleByIndexInArr(randNum, articlesbson)
}


func showResAddNewArticle(id string){
	fmt.Printf("insert ID:%s\n", id)
}


func AddNewArticleToDB(newArticle primitive.D){
	r, err := DB.CollectionDB.InsertOne(ctx, newArticle)
	if err != nil{
		panic(err)
	}
	showResAddNewArticle(fmt.Sprintf("%s", r.InsertedID))
}


func CreateNewArticleToDB(articleURL string, chatID int64){
	startDB("WhatMeReadDB", "users")
	defer DB.Client.Disconnect(ctx)

	newArticle := createArtcile(articleURL, chatID)
	AddNewArticleToDB(newArticle)
}


func createArtcile(articleURL string, chatID int64) primitive.D{
	return bson.D{
		{Key: "article", Value: articleURL},
		{Key: "chatID", Value: chatID},
	}
}


func startDB(nameDB, nameCollection string){
	var err error

	DB.ClientOptions = options.Client().ApplyURI(config.LinkToDB)
	DB.Client, err = mongo.Connect(ctx, DB.ClientOptions)
	if err != nil{
		panic(err)
	}

	DB.DBMongo = DB.Client.Database("WhatMeReadDB")
	DB.CollectionDB = DB.DBMongo.Collection("users")
}