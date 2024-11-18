package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
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

// ポール
type Pole struct {
	id      string
	stop_id string
}

// 系統
type Route struct {
	id   string
	name string
}

// 便
type Trip struct {
	id        string
	route_id  string
	yobi      string
	stopTimes []StopTime
}

// バス停と着発時刻
type StopTime struct {
	stop_id        string
	arrival_time   string
	departure_time string
}

// Stop連想配列
// キー stop_id
// 値 Stop
var stopMap map[string]Stop = make(map[string]Stop)

// Pole配列
var poleList []Pole

// Route連想配列
// キー route_id
// 値 Route
var routeMap map[string]Route = make(map[string]Route)

// Trip配列
var tripList []Trip

func main() {
	fmt.Println("処理開始")

	// inpput/StopMaster.tsvを読み込んで、stopをstopMapに格納
	readStopMasterTsv()
	// inpput/StopPoleMaster.tsvを読み込んで、poleをpole配列に格納
	readStopPoleMasterTsv()
	// pole配列の要素をstops.txtに出力
	writeStopsTxt()
	// stopMap連想配列の要素をtranslations.txtに出力
	writeTranslationsTxt()

	// inpput/RouteMaster.tsvを読み込んで、routeをrouteMapに格納
	readRouteMasterTsv()
	// routeMap連想配列の要素をroutes.txtに出力
	writeRoutesTxt()

	// inpput/DiaMaster.tsvを読み込んで、tripListに格納
	readDiaMasterTsv()

	// tripListの要素をtrips.txtに出力
	writeTripsTxt()

	// tripListの要素をstop_times.txtに出力
	writeStopTimesTxt()

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

// inpput/StopPoleMaster.tsvを読み込んで、poleをpole配列に格納
func readStopPoleMasterTsv() {
	fmt.Println("StopPoleMaster.tsv読み込み")
	var file string = "input/StopPoleMaster.tsv"
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
		// pole構造体を作成
		var pole Pole = Pole{}
		// stop構造体に分割された要素を格納
		pole.id = elements[1]
		pole.stop_id = elements[2]
		poleList = append(poleList, pole)
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

// inpput/DiaMaster.tsvを読み込んで、tripListに格納
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
		var trip Trip = Trip{}
		// dia構造体に分割された要素を格納
		var yobi string = elements[1]
		if yobi == "1" {
			trip.yobi = "平"
		} else if yobi == "2" {
			trip.yobi = "日"
		} else if yobi == "4" {
			trip.yobi = "土"
		}
		trip.route_id = elements[4]

		// 5列目 stop_id 6列名 着時刻 7列名 発時刻 を stopTimeに格納  8列名 9列目は捨てる
		// 10列名以降はその繰り返し
		// 5列名以降の繰り返しの数を計算
		var elementSize int = len(elements)
		var blockCnt int = (len(elements) - 4) / 5

		for i := 0; i < blockCnt; i++ {
			var stopTime StopTime
			stopTime.stop_id = elements[5+i*5]
			if 6+i*3 < elementSize {
				stopTime.arrival_time = toTime(elements[6+i*5])
			}
			if 7+i*3 < elementSize {
				stopTime.departure_time = toTime(elements[7+i*5])
			}
			trip.stopTimes = append(trip.stopTimes, stopTime)

			if i == 0 {
				trip.id = trip.route_id + "_" + trip.yobi + "_" + stopTime.departure_time
			}
		}

		// dia配列にdiaを追加
		tripList = append(tripList, trip)
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
	for _, pole := range poleList {

		if stop, ok := stopMap[pole.stop_id]; ok {
			// stopをstops.txtに出力
			data := []string{
				pole.id,
				stop.name,
			}
			writer.Write(data)
		}

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

// tripListの要素をtrips.txtに出力
func writeTripsTxt() {
	fmt.Println("trips.txt出力")
	file, _ := os.Create("output/trips.txt")
	defer file.Close()
	var writer *csv.Writer = csv.NewWriter(transform.NewWriter(file, japanese.ShiftJIS.NewEncoder()))
	writer.UseCRLF = true //改行コードを\r\nにする
	// 見出し行を出力
	data := []string{
		"route_id",
		"service_id",
		"trip_id",
	}
	writer.Write(data)
	// tripListの要素を取り出しながらループ
	for _, trip := range tripList {
		// tripをtrips.txtに出力
		data := []string{
			trip.route_id,
			trip.yobi,
			trip.id,
		}
		writer.Write(data)
	}
	writer.Flush()
}

// tripListの要素をstop_times.txtに出力
func writeStopTimesTxt() {
	fmt.Println("stop_times.txt出力")
	file, _ := os.Create("output/stop_times.txt")
	defer file.Close()
	var writer *csv.Writer = csv.NewWriter(transform.NewWriter(file, japanese.ShiftJIS.NewEncoder()))
	writer.UseCRLF = true //改行コードを\r\nにする
	// 見出し行を出力
	data := []string{
		"trip_id",
		"arrival_time",
		"departure_time",
		"stop_id",
		"stop_sequence",
	}
	writer.Write(data)
	// tripListの要素を取り出しながらループ
	for _, trip := range tripList {
		var sequence int = 1
		for _, stopTime := range trip.stopTimes {
			data := []string{
				trip.id,
				stopTime.arrival_time,
				stopTime.departure_time,
				stopTime.stop_id,
				strconv.Itoa(sequence),
			}
			writer.Write(data)
			sequence++
		}
	}
	writer.Flush()
}

// 時刻文字列を返す
// 610→6:10
// 1725→17:25
func toTime(str string) string {
	var len int = len(str)
	var time string
	if len == 3 {
		time = str[0:1] + ":" + str[1:]
	} else if len == 4 {
		time = str[0:2] + ":" + str[1:]
	}
	return time
}

// SJISをUTF8に変換
func sjis_to_utf8(str string) string {
	var iostr = transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder())
	var iobyte, _ = io.ReadAll(iostr)
	return string(iobyte)
}
