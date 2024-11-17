package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// バス停
type Stop struct {
	id   string
	name string
	yomi string
}

// 系統
type Route struct {
	id   string
	name string
}

// ダイヤ
type Dia struct {
	trip_id string
	bins    []Bin
}

// 便
type Bin struct {
	arrival_time   string
	departure_time string
	stop_id        string
}

// Stop連想配列
// キー stop_id
// 値 Stop
var stopMap map[string]Stop = make(map[string]Stop)

// Route連想配列
// キー route_id
// 値 Route
var routeMap map[string]Route = make(map[string]Route)

// Dia配列
var diaList []Dia

func main() {
	fmt.Println("処理開始")

	// inpput/StopMaster.tsvを読み込んで、stopをstopMapに格納
	readStopMasterTsv()
	// stopMap連想配列の要素をstops.txtに出力
	writeStopsTxt()
	// stopMap連想配列の要素をtranslations.txtに出力
	writeTranslationsTxt()

	// inpput/RouteMaster.tsvを読み込んで、routeをrouteMapに格納
	readRouteMasterTsv()
	// routeMap連想配列の要素をroutes.txtに出力
	writeRoutesTxt()

	// inpput/DiaMaster.tsvを読み込んで、diaをdiaListに格納
	readDiaMasterTsv()

	fmt.Println("処理終了")
}

// inpput/StopMaster.tsvを読み込んで、stopをstopMap連想配列に格納
func readStopMasterTsv() {
	fmt.Println("StopMaster.tsv読み込み")
	var file string = "input/StopMaster.tsv"
	if _, err := os.Stat(file); err != nil {
		fmt.Println("ファイルは存在しません！" + file)
		os.Exit(1)
	}
	data, _ := os.Open(file)
	defer data.Close()

	var line string

	scanner := bufio.NewScanner(data)
	// 1行ずつ読み込み
	for scanner.Scan() {
		// 1行読み込み
		line = sjis_to_utf8(scanner.Text())
		// 1行をタブで分割
		elements := strings.Split(line, "\t")
		// stop構造体を作成
		var stop Stop = Stop{}
		// stop構造体に分割された要素を格納
		stop.id = elements[1]
		stop.name = elements[2]
		stop.yomi = elements[3]
		stopMap[stop.id] = stop
	}
}

// inpput/RouteMaster.tsvを読み込んで、routeをrouteMap連想配列に格納
func readRouteMasterTsv() {
	fmt.Println("RouteMaster.tsv読み込み")
	var file string = "input/RouteMaster.tsv"
	if _, err := os.Stat(file); err != nil {
		fmt.Println("ファイルは存在しません！" + file)
		os.Exit(1)
	}
	data, _ := os.Open(file)
	defer data.Close()

	var line string

	scanner := bufio.NewScanner(data)
	// 1行ずつ読み込み
	for scanner.Scan() {
		// 1行読み込み
		line = sjis_to_utf8(scanner.Text())
		// 1行をタブで分割
		elements := strings.Split(line, "\t")
		// route構造体を作成
		var route Route = Route{}
		// stop構造体に分割された要素を格納
		route.id = elements[1]
		route.name = elements[5]
		routeMap[route.id] = route
	}
}

// inpput/DiaMaster.tsvを読み込んで、diaをdiaListに格納
func readDiaMasterTsv() {
	fmt.Println("DiaMaster.tsv読み込み")
	var file string = "input/DiaMaster.tsv"
	if _, err := os.Stat(file); err != nil {
		fmt.Println("ファイルは存在しません！" + file)
		os.Exit(1)
	}
	data, _ := os.Open(file)
	defer data.Close()

	var line string

	scanner := bufio.NewScanner(data)
	// 1行ずつ読み込み
	for scanner.Scan() {
		// 1行読み込み
		line = sjis_to_utf8(scanner.Text())
		// 1行をタブで分割
		elements := strings.Split(line, "\t")
		// stop構造体を作成
		var dia Dia = Dia{}
		// dia構造体に分割された要素を格納
		dia.trip_id = elements[1]
		// dia配列にdiaを追加
		diaList = append(diaList, dia)
	}
}

// stopMap連想配列の要素をstops.txtに出力
func writeStopsTxt() {
	fmt.Println("stops.txt出力")
	file, _ := os.Create("output/stops.txt")
	defer file.Close()
	var writer *csv.Writer = csv.NewWriter(transform.NewWriter(file, japanese.ShiftJIS.NewEncoder()))
	writer.UseCRLF = true //改行コードを\r\nにする
	// 見出し行を出力
	data := []string{
		"stop_id",
		"stop_name",
	}
	writer.Write(data)
	// stopMap連想配列の要素を取り出しながらループ
	for _, stop := range stopMap {
		// stopをstops.txtに出力
		data := []string{
			stop.id,
			stop.name,
		}
		writer.Write(data)
	}
	writer.Flush()
}

// stopMap連想配列の要素をtranslations.txtに出力
func writeTranslationsTxt() {
	fmt.Println("translations.txt出力")
	file, _ := os.Create("output/translations.txt")
	defer file.Close()
	var writer *csv.Writer = csv.NewWriter(transform.NewWriter(file, japanese.ShiftJIS.NewEncoder()))
	writer.UseCRLF = true //改行コードを\r\nにする
	// 見出し行を出力
	data := []string{
		"table_name",
		"field_name",
		"language",
		"translation",
	}
	writer.Write(data)
	// stopMap連想配列の要素を取り出しながらループ
	for _, stop := range stopMap {
		// stopをstops.txtに出力
		data := []string{
			"stops",
			"stop_name",
			"ja-Hrkt",
			stop.yomi,
		}
		writer.Write(data)
	}
	writer.Flush()
}

// routeMap連想配列の要素をroutes.txtに出力
func writeRoutesTxt() {
	fmt.Println("routes.txt出力")
	file, _ := os.Create("output/routes.txt")
	defer file.Close()
	var writer *csv.Writer = csv.NewWriter(transform.NewWriter(file, japanese.ShiftJIS.NewEncoder()))
	writer.UseCRLF = true //改行コードを\r\nにする
	// 見出し行を出力
	data := []string{
		"route_id",
		"route_long_name",
	}
	writer.Write(data)
	// routeMap連想配列の要素を取り出しながらループ
	for _, route := range routeMap {
		// routeをsroutes.txtに出力
		data := []string{
			route.id,
			route.name,
		}
		writer.Write(data)
	}
	writer.Flush()
}

// SJISをUTF8に変換
func sjis_to_utf8(str string) string {
	var iostr = transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder())
	var iobyte, _ = io.ReadAll(iostr)
	return string(iobyte)
}
