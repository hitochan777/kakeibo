package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/hitochan777/kakeibo/backend/converter/service"
	"google.golang.org/grpc"
)

func convert(text string) (*service.Item, error) {
	// Do some conversion
	re := regexp.MustCompile(`([\p{Hiragana}|\p{Katakana}|\p{Han}]+)の([\p{Hiragana}|\p{Katakana}|\p{Han}]+)が(\d+)円`)
	match := re.FindStringSubmatch(text)
	if match == nil {
		return nil, fmt.Errorf("Failed to parse raw text")
	}
	category := &service.Category{
		Big:   match[1],
		Small: match[2],
	}
	price, err := strconv.Atoi(match[3])
	if err != nil {
		return nil, err
	}
	return &service.Item{Category: category, Price: int32(price)}, nil
}

func main() {
	// gRPC client setup
	conn, err := grpc.Dial("localhost:11111", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := service.NewKakeiboClient(conn)

	// Router setup
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rawText := string(body)
		log.Println("Received: ", rawText)
		item, err := convert(rawText)
		log.Println(item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err = client.AddItem(ctx, item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
