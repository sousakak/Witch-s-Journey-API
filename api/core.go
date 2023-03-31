package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

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

func Main(params Params) string {
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
					single_datas = "[\"" + rowresp.Values[0][elemIndex].(string) + "\"]"
				}
			} else {
				v := rowresp.Values[0]
				for len(v) < 5 {
					v = append(v, "")
				}
				raw_datas = append(raw_datas, Data{Name: v[0].(string), Witch_name: v[1].(string), Called_name: v[2].(string), Desc: v[3].(string), Chap: v[4].(string)})
			}
		}
	case params.char == "" && params.elem != "":
		var single_datas string
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
			single_datas = "["
			for index, content := range resp.Values[0] {
				if index == 0 {
					continue
				}
				single_datas += "\"" + content.(string) + "\","
			}
			single_datas = single_datas[:len(single_datas)-1]
			single_datas += "]"
		} else {
			err := errors.New("index out of range: jumped out from a number of elements")
			fmt.Println(err.Error())
		}
		return single_datas
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
