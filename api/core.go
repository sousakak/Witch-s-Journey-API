package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Data struct {
	name        string `json:"name"`
	witch_name  string `json:"witch-name"`
	called_name string `json:"called-name"`
	desc        string `json:"description"`
	chap        string `json:"chapter"`
}

func Main() string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("環境変数の読み込みに失敗しました: %v", err)
	}
	var spreadsheetID = os.Getenv("DATABASE")
	credential := option.WithCredentialsFile("client_secret.json")
	srv, err := sheets.NewService(context.TODO(), credential)
	if err != nil {
		log.Fatal(err)
	}
	readRange := "イレイナ!A:E"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalln(err)
	}
	if len(resp.Values) == 0 {
		log.Fatalln("data not found")
	}

	raw_datas := []Data{}
	for index, row := range resp.Values {
		if index == 0 {
			continue
		}
		for len(row) == 5 {
			row = append(row, "")
		}
		single_data := Data{row[0].(string), row[1].(string), row[2].(string), row[3].(string), row[4].(string)}
		raw_datas = append(raw_datas, single_data)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(&raw_datas); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}
