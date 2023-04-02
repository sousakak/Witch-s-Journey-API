// ███╗░░░███╗░█████╗░██╗███╗░░██╗░░░░██████╗░░█████╗░
// ████╗░████║██╔══██╗██║████╗░██║░░░██╔════╝░██╔══██╗
// ██╔████╔██║███████║██║██╔██╗██║░░░██║░░██╗░██║░░██║
// ██║╚██╔╝██║██╔══██║██║██║╚████║░░░██║░░╚██╗██║░░██║
// ██║░╚═╝░██║██║░░██║██║██║░╚███║██╗╚██████╔╝╚█████╔╝
// ╚═╝░░░░░╚═╝╚═╝░░╚═╝╚═╝╚═╝░░╚══╝╚═╝░╚═════╝░░╚════╝░

package main

import (
	"Witchs-Journey-API/utilities"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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

type EachData struct {
	Name        string `json:"name"`
	Witch_name  string `json:"witch_name"`
	Called_name string `json:"called_name"`
	Desc        string `json:"desc"`
	Chap        string `json:"chap"`
}

type IndexPage struct {
	Lang string
}

func Handle404(w http.ResponseWriter, r *http.Request) {
	f, err := ioutil.ReadFile("404.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, string(f))
	}
}

func Api(params Params) []byte {
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

	raw_datas := map[string]EachData{}
	var single_datas []string
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
				rtCstStruct := EachData{Name: "", Witch_name: "", Called_name: "", Desc: "", Chap: ""}
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
					single_datas = append(single_datas, rowresp.Values[0][elemIndex].(string))
				}
			} else {
				v := rowresp.Values[0]
				for len(v) < 5 {
					v = append(v, "")
				}
				raw_datas[v[0].(string)] = EachData{Name: v[0].(string), Witch_name: v[1].(string), Called_name: v[2].(string), Desc: v[3].(string), Chap: v[4].(string)}
			}
		} else {
			err := errors.New("index out of range: jumped out from a number of elements")
			fmt.Println(err.Error())
		}
	case params.char == "" && params.elem != "":
		rtCstStruct := EachData{Name: "", Witch_name: "", Called_name: "", Desc: "", Chap: ""}
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
		target := params.sheet + "!"
		if elemIndex != -1 {
			target = target + "ABCDE"[(elemIndex):(elemIndex+1)] + ":" + "ABCDE"[(elemIndex):(elemIndex+1)]
			resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, target).MajorDimension("COLUMNS").Do()
			if err != nil {
				log.Fatalln(err)
			}
			if len(resp.Values) == 0 {
				log.Fatalln("data not found")
			}
			for index, content := range resp.Values[0] {
				if index == 0 {
					continue
				}
				single_datas = append(single_datas, content.(string))
			}
		} else {
			err := errors.New("index out of range: jumped out from a number of elements")
			fmt.Println(err.Error())
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
			single_data := EachData{Name: row[0].(string), Witch_name: row[1].(string), Called_name: row[2].(string), Desc: row[3].(string), Chap: row[4].(string)}
			raw_datas[row[0].(string)] = single_data
		}
	}

	var bytes []byte
	if len(single_datas) != 0 {
		bytes, err = json.Marshal(single_datas)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		bytes, err = json.Marshal(raw_datas)
		if err != nil {
			log.Fatal(err)
		}
	}
	return bytes
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Handle404(w, r)
		return
	}
	var local_index *IndexPage
	if r.URL.Query().Get("lang") != "" {
		local_index = &IndexPage{r.URL.Query().Get("lang")}
		utilities.SetLangCookies(w, r.URL.Query().Get("lang"))
	} else if utilities.GetCookie(r) != nil {
		local_index = &IndexPage{utilities.GetCookie(r).Value}
	} else {
		local_index = &IndexPage{"en"}
	}

	t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		panic(err.Error())
	}
	if err := t.Execute(w, local_index); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		panic(err.Error())
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	sheet := "イレイナ"
	char := ""
	elem := ""
	if r.URL.Path != "/api/" {
		http.NotFound(w, r)
		return
	}
	switch r.URL.Query().Get("introducer") {
	case "":
		if r.URL.Query().Get("character") != "" {
			char = strings.Title(r.URL.Query().Get("character"))
		}
		if r.URL.Query().Get("element") != "" {
			elem = strings.Title(r.URL.Query().Get("element"))
		}
	case "イレイナ", "サヤ", "白石定規", "Elaina", "Saya", "Jogi":
		sheet = r.URL.Query().Get("introducer")
		if m, _ := regexp.MatchString("[a-z]", sheet); m {
			switch sheet {
			case "Elaina":
				sheet = "イレイナ"
			case "Saya":
				sheet = "サヤ"
			case "Jogi":
				sheet = "白石定規"
			}
		}
		if r.URL.Query().Get("character") != "" {
			char = strings.Title(r.URL.Query().Get("character"))
		}
		if r.URL.Query().Get("element") != "" {
			elem = strings.Title(r.URL.Query().Get("element"))
		}
	default:
		var emptyData []string
		bytes, err := json.Marshal(emptyData)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(bytes)
		return
	}
	params := Params{sheet: sheet, char: char, elem: elem}
	json_file := Api(params)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write(json_file)
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
