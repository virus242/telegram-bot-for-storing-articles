package internal

import (
	"log"
	"net/http"
)

func makeRequestGetStatusCode(URL string)int{
	resp, err := http.Get(URL)
	if err != nil{
		log.Printf(err.Error())
		return -1
	}
	return resp.StatusCode
}

func CheckURL(URL string)bool{
	if makeRequestGetStatusCode(URL) != 200{
		return false
	}
	return true
}