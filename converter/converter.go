package service

import (
	"fmt"
	"net/http"
)

type Converter struct {
	client *KakeiboClient
}

func (c *Converter) Convert(text string) error {
	// Do some conversion
	itemInput := &ItemInput{}
	return nil
}

func Start() {
	converter := &Converter{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}
		data := converter.Convert(body)
		response, err := c.client.AddItem(context.Background(), itemInput)
		if err != nil {
			return err
		}
		if !response.Ok {
			return fmt.Errorf("Got failed response")
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
