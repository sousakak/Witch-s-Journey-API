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
	"reflect"
	"regexp"
	"strconv"
	"text/template"

	"github.com/joho/godotenv"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

//go:embed assets/*
var assets embed.FS

type Params struct {
	sheet string
	char  string
	elem  string
}

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

func Api(params Params) string {
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

	raw_datas := []Data{}
	var single_datas string
	switch {
	case params.char != "":
		readRange := params.sheet + "!A:E"
		resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).MajorDimension("COLUMNS").Do()
		if err != nil {
			log.Fatalln(err)
		}
		if len(resp.Values) == 0 {
			log.Fatalln("data not found")
		}

		var rownum int = -1
		for index, name := range resp.Values[0] {
			if name.(string) == params.char {
				rownum = index + 1
				break
			}
		}
		println(params.elem)
		if rownum != -1 {
			targetRow := params.sheet + "!A" + strconv.Itoa(rownum) + ":E" + strconv.Itoa(rownum)
			rowresp, err := srv.Spreadsheets.Values.Get(spreadsheetID, targetRow).Do()
			if err != nil {
				log.Fatalln(err)
			}
			if len(rowresp.Values) == 0 {
				log.Fatalln("data not found")
			}
			if params.elem != "" {
				rtCstStruct := Data{Name: "", Witch_name: "", Called_name: "", Desc: "", Chap: ""}
				rtCst := reflect.TypeOf(rtCstStruct)
				var elemList = []string{}
				for i := 0; i < rtCst.NumField(); i++ {
					elemList = append(elemList, rtCst.Field(i).Name)
				}
				var elemIndex int = -1
				for i, elemName := range elemList {
					if elemName == params.elem {
						elemIndex = i
						break
					} else {
						continue
					}
				}
				if elemIndex != -1 {
					single_datas = rowresp.Values[0][elemIndex].(string)
				}
			} else {
				v := rowresp.Values[0]
				for len(v) < 5 {
					v = append(v, "")
				}
				raw_datas = append(raw_datas, Data{Name: v[0].(string), Witch_name: v[1].(string), Called_name: v[2].(string), Desc: v[3].(string), Chap: v[4].(string)})
			}
		}
	default:
		readRange := params.sheet + "!A:E"
		resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
		if err != nil {
			log.Fatalln(err)
		}
		if len(resp.Values) == 0 {
			log.Fatalln("data not found")
		}

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
	}

	if single_datas != "" {
		return single_datas
	} else {
		bytes, err := json.Marshal(raw_datas)
		if err != nil {
			log.Fatal(err)
		}
		return string(bytes)
	}
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
			fmt.Printf("char\n")
		}
		if r.URL.Query().Get("element") != "" {
			elem = r.URL.Query().Get("element")
			fmt.Printf("elem\n")
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
	params := Params{sheet: sheet, char: char, elem: elem}
	json_file := Api(params)
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
