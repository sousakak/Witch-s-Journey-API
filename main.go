// ███╗░░░███╗░█████╗░██╗███╗░░██╗░░░░██████╗░░█████╗░
// ████╗░████║██╔══██╗██║████╗░██║░░░██╔════╝░██╔══██╗
// ██╔████╔██║███████║██║██╔██╗██║░░░██║░░██╗░██║░░██║
// ██║╚██╔╝██║██╔══██║██║██║╚████║░░░██║░░╚██╗██║░░██║
// ██║░╚═╝░██║██║░░██║██║██║░╚███║██╗╚██████╔╝╚█████╔╝
// ╚═╝░░░░░╚═╝╚═╝░░╚═╝╚═╝╚═╝░░╚══╝╚═╝░╚═════╝░░╚════╝░

package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"

	"github.com/joho/godotenv"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

//go:embed assets/*
var assets embed.FS

type Data struct {
	Name        string `json:"name"`
	Witch_name  string `json:"witch_name"`
	Called_name string `json:"called_name"`
	Desc        string `json:"desc"`
	Chap        string `json:"chap"`
}

type IndexPage struct {
	Lang string
}

func Api(sheet string, char string, elem string) string {
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
	readRange := sheet + "!A:E"
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
		for len(row) < 5 {
			row = append(row, "")
		}
		single_data := Data{Name: row[0].(string), Witch_name: row[1].(string), Called_name: row[2].(string), Desc: row[3].(string), Chap: row[4].(string)}
		raw_datas = append(raw_datas, single_data)
	}

	bytes, err := json.Marshal(raw_datas)
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	t, err := template.ParseFiles("index.html")
	local_index := &IndexPage{"ja"}
	if err != nil {
		panic(err.Error())
	}
	if err := t.Execute(w, local_index); err != nil {
		panic(err.Error())
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	sheet := "イレイナ"
	char := ""
	elem := ""
	switch r.URL.Path {
	case "/api/":
		if r.URL.Query().Get("character") != "" {
			char = r.URL.Query().Get("character")
			fmt.Printf("char")
		}
		if r.URL.Query().Get("element") != "" {
			char = r.URL.Query().Get("element")
			fmt.Printf("elem")
		}
	case "/api/イレイナ", "/api/サヤ", "/api/白石定規", "/api/elaina", "/api/saya", "/api/jogi":
		sheet = r.URL.Path[5:]
		if m, _ := regexp.MatchString("[a-z]", sheet); m {
			switch sheet {
			case "elaina":
				sheet = "イレイナ"
			case "saya":
				sheet = "サヤ"
			case "jogi":
				sheet = "白石定規"
			}
		}
		if r.URL.Query().Get("character") != "" {
			char = r.URL.Query().Get("character")
			fmt.Printf("char")
		}
		if r.URL.Query().Get("element") != "" {
			char = r.URL.Query().Get("element")
			fmt.Printf("elem")
		}
	default:
		http.NotFound(w, r)
		return
	}
	json_file := Api(sheet, char, elem)
	_, err := fmt.Fprint(w, json_file)
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/api/", apiHandler)
	http.Handle("/assets/", http.FileServer(http.FS(assets)))
	http.ListenAndServe(":8080", nil)
}
