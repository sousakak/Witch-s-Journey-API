package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {
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
	var jsonfile string
	jsonfile += `{`
	for _, row := range resp.Values {
		data := make([]interface{}, 5)
		for i := 0; i < 5; i++ {
			if (i + 1) > len(row) {
				data[i] = ""
			} else {
				data[i] = row[i]
			}
		}
		data = append(data, row)
		jsonfile += `` + data[0].(string) + `: {"witch-name":` + data[1].(string) + `,"called-name":` + data[2].(string) + `,"description":` + data[3].(string) + `,"chapter":` + data[4].(string) + `},`
	}
	jsonfile += `}`
	fmt.Print(jsonfile)
}
