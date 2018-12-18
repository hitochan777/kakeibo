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
	re := regexp.MustCompile(`(\d+)月(\d+)日に(\d+)円で([\p{Hiragana}|\p{Katakana}|\p{Han}]+)の([\p{Hiragana}|\p{Katakana}|\p{Han}]+)に追加`)
	match := re.FindStringSubmatch(text)
	if match == nil {
		return nil, fmt.Errorf("Failed to parse raw text")
	}
	month, err := strconv.Atoi(match[1])
	if err != nil {
		return nil, err
	}
	date, err := strconv.Atoi(match[2])
	if err != nil {
		return nil, err
	}
	price, err := strconv.Atoi(match[3])
	if err != nil {
		return nil, err
	}
	payedAt := &service.PayedAt{
		Month: int32(month),
		Date:  int32(date),
	}
	category := &service.Category{
		Big:   match[4],
		Small: match[5],
	}
	return &service.Item{PayedAt: payedAt, Category: category, Price: int32(price)}, nil
}

func main() {
	// gRPC client setup
	conn, err := grpc.Dial("localhost:11111", grpc.WithInsecure()) // TODO: change to proper port when ready
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
		item, err := convert(string(body))
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
