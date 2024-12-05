//go:generate goversioninfo

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	fmt.Println("処理開始")

	// inpput/StopMaster.tsvを読み込んで、stopをstopMapに格納
	readStopMasterTsv()
	// inpput/StopPoleMaster.tsvを読み込んで、poleをpole配列に格納
	readStopPoleMasterTsv()

	// inpput/RouteMaster.tsvを読み込んで、routeをrouteMapに格納
	readRouteMasterTsv()
	// routeMap連想配列の要素をroutes.txtに出力
	writeRoutesTxt()

	// inpput/DiaMaster.tsvを読み込んで、tripListに格納
	readDiaMasterTsv()

	// pole配列の要素をstops.txtに出力
	writeStopsTxt()
	// stopMap連想配列の要素をtranslations.txtに出力
	writeTranslationsTxt()

	// tripListの要素をtrips.txtに出力
	writeTripsTxt()

	// tripListの要素をstop_times.txtに出力
	writeStopTimesTxt()

	// calendar.txtに出力
	writeCalendarTxt()

	// agency.txtに出力
	writeAgencyTxt()

	// inpput/GenerationMaster.tsvを読み込んで、feed_info.txtに出力
	readGenerationMasterTsvAndWriteFeedInfoTxt()

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
		stop.id = maeZero(elements[1], 4) + "000"
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
		pole.id = maeZero(elements[1], 7)
		pole.stop_id = elements[2]
		pole.name = elements[4]
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
		// trip構造体を作成
		var trip Trip = Trip{}
		var yobi string = elements[1]
		if yobi == "1" {
			trip.yobi = "1_平日"
		} else if yobi == "2" {
			trip.yobi = "2_日祝"
		} else if yobi == "4" {
			trip.yobi = "4_土曜"
		} else if yobi == "3" {
			trip.yobi = "3_特殊"
		}
		trip.col3 = elements[3]
		trip.route_id = elements[4]

		// 5列目 stop_id 6列名 着時刻 7列名 発時刻 を stopTimeに格納  8列名 9列目は捨てる
		// 10列名以降はその繰り返し
		// 5列名以降の繰り返しの数を計算
		var elementSize int = len(elements)
		var blockCnt int = (len(elements) - 5) / 5

		for i := 0; i < blockCnt; i++ {
			// stopTime構造体に分割された要素を格納
			var stopTime StopTime
			stopTime.stop_id = elements[5+i*5]

			if 6+i*3 < elementSize {
				stopTime.arrival_time = toTime(elements[6+i*5])
			}
			if 7+i*3 < elementSize {
				stopTime.departure_time = toTime(elements[7+i*5])
			}
			//発駅が空の時着駅時刻を埋める
			if stopTime.arrival_time == "00:00:00" {
				stopTime.arrival_time = stopTime.departure_time
			}
			//着駅が空の時発駅時刻を埋める
			if stopTime.departure_time == "00:00:00" {
				stopTime.departure_time = stopTime.arrival_time
			}
			trip.stopTimes = append(trip.stopTimes, stopTime)

			if i == 0 {
				trip.id = trip.route_id + "_" + trip.yobi + "_" + stopTime.departure_time
			}

		}

		// trip配列にtripを追加
		tripList = append(tripList, trip)
	}
}

// poleListの要素をstops.txtに出力
func writeStopsTxt() {
	fmt.Println("stops.txt出力")

	// stop_times.txtに出力する stop_id をマップに格納
	var poleIdMap map[string]bool = make(map[string]bool)

	// tripListの要素を取り出しながらループ
	for _, trip := range tripList {
		for _, stopTime := range trip.stopTimes {
			if _, ok := poleIdMap[stopTime.stop_id]; ok {
				continue
			} else {
				poleIdMap[stopTime.stop_id] = false
			}
		}
	}

	file, _ := os.Create("output/stops.txt")
	defer file.Close()
	var writer *csv.Writer = csv.NewWriter(file)
	writer.UseCRLF = true //改行コードを\r\nにする
	// 見出し行を出力
	data := []string{
		"stop_id",
		"stop_name",
		"stop_lat",
		"stop_lon",
	}
	writer.Write(data)
	for _, pole := range poleList {
		// poleをstops.txtに出力
		if _, ok := poleIdMap[pole.id]; ok {
			data := []string{
				pole.id,
				pole.name,
				"",
				"",
			}
			writer.Write(data)
			poleIdMap[pole.id] = true
		}
	}
	for stop_id, ok := range poleIdMap {
		if !ok {
			if stop, ok := stopMap[stop_id]; ok {
				data := []string{
					stop.id,
					stop.name,
					"",
					"",
				}
				writer.Write(data)
			}
		}
	}

	writer.Flush()
}

// calendar.txtを出力
func writeCalendarTxt() {
	fmt.Println("calendar.txtを出力")
	file, _ := os.Create("output/calendar.txt")
	defer file.Close()

	var writer *csv.Writer = csv.NewWriter(file)
	writer.UseCRLF = true //改行コードを\r\nにする
	// 見出し行を出力
	data := []string{
		"service_id",
		"monday",
		"tuesday",
		"wednesday",
		"thursday",
		"friday",
		"saturday",
		"sunday",
		"start_date",
		"end_date",
	}
	writer.Write(data)

	today := time.Now()
	// 月初
	gessho := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.UTC)
	// 1年後の月末
	after1Year := gessho.AddDate(1, 0, -1)

	data = []string{
		"1_平日",
		"1",
		"1",
		"1",
		"1",
		"1",
		"",
		"",
		strconv.Itoa(gessho.Year()) + maeZero(strconv.Itoa((int)(gessho.Month())), 2) + "01",
		strconv.Itoa(after1Year.Year()) + maeZero(strconv.Itoa((int)(after1Year.Month())), 2) + maeZero(strconv.Itoa(after1Year.Day()), 2),
	}
	writer.Write(data)

	data = []string{
		"2_日祝",
		"",
		"",
		"",
		"",
		"",
		"",
		"1",
		strconv.Itoa(gessho.Year()) + maeZero(strconv.Itoa((int)(gessho.Month())), 2) + "01",
		strconv.Itoa(after1Year.Year()) + maeZero(strconv.Itoa((int)(after1Year.Month())), 2) + maeZero(strconv.Itoa(after1Year.Day()), 2),
	}
	writer.Write(data)

	data = []string{
		"4_土曜",
		"",
		"",
		"",
		"",
		"",
		"1",
		"",
		strconv.Itoa(gessho.Year()) + maeZero(strconv.Itoa((int)(gessho.Month())), 2) + "01",
		strconv.Itoa(after1Year.Year()) + maeZero(strconv.Itoa((int)(after1Year.Month())), 2) + maeZero(strconv.Itoa(after1Year.Day()), 2),
	}
	writer.Write(data)

	data = []string{
		"3_特殊",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		strconv.Itoa(gessho.Year()) + maeZero(strconv.Itoa((int)(gessho.Month())), 2) + "01",
		strconv.Itoa(after1Year.Year()) + maeZero(strconv.Itoa((int)(after1Year.Month())), 2) + maeZero(strconv.Itoa(after1Year.Day()), 2),
	}
	writer.Write(data)

	writer.Flush()
}

// agency.txtを出力
func writeAgencyTxt() {
	fmt.Println("agency.txtを出力")
	file, _ := os.Create("output/agency.txt")
	defer file.Close()
	var writer *csv.Writer = csv.NewWriter(file)
	writer.UseCRLF = true //改行コードを\r\nにする
	// 見出し行を出力
	data := []string{
		"agency_id",
		"agency_name",
		"agency_url",
		"agency_timezone",
	}
	writer.Write(data)

	data = []string{
		"八戸市交通部",
		"八戸市交通部",
		"",
		"Asia/Tokyo",
	}
	writer.Write(data)

	writer.Flush()
}

// poleListの要素をtranslations.txtに出力
func writeTranslationsTxt() {
	fmt.Println("translations.txt出力")
	file, _ := os.Create("output/translations.txt")
	defer file.Close()
	var writer *csv.Writer = csv.NewWriter(file)
	writer.UseCRLF = true //改行コードを\r\nにする
	// 見出し行を出力
	data := []string{
		"table_name",
		"field_name",
		"language",
		"field_value",
		"translation",
	}
	writer.Write(data)

	var poleNameMap map[string]string = make(map[string]string)

	// stopMap連想配列の要素を取り出しながらループ
	for _, pole := range poleList {

		if _, ok := poleNameMap[pole.name]; ok {
			continue
		} else {
			poleNameMap[pole.name] = pole.name
		}

		if stop, ok := stopMap[pole.stop_id]; ok {
			// stopをstops.txtに出力
			data := []string{
				"stops",
				"stop_name",
				"ja",
				pole.name,
				pole.name,
			}
			writer.Write(data)
			data = []string{
				"stops",
				"stop_name",
				"ja-Hrkt",
				pole.name,
				stop.yomi,
			}
			writer.Write(data)
		}
	}
	writer.Flush()
}

// routeMap連想配列の要素をroutes.txtに出力
func writeRoutesTxt() {
	fmt.Println("routes.txt出力")
	file, _ := os.Create("output/routes.txt")
	defer file.Close()
	var writer *csv.Writer = csv.NewWriter(file)
	writer.UseCRLF = true //改行コードを\r\nにする
	// 見出し行を出力
	data := []string{
		"route_id",
		"agency_id",
		"route_long_name",
	}
	writer.Write(data)
	// routeMap連想配列の要素を取り出しながらループ
	for _, route := range routeMap {
		// routeをsroutes.txtに出力
		data := []string{
			route.id,
			"八戸市交通部",
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
	var writer *csv.Writer = csv.NewWriter(file)
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

		if trip.id == "" {
			continue
		}

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
	var writer *csv.Writer = csv.NewWriter(file)
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

		if trip.col3 == "2" {
			continue
		}

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

// inpput/GenerationMaster.tsvを読み込んで、feed_info.txtに出力
func readGenerationMasterTsvAndWriteFeedInfoTxt() {
	fmt.Println("GenerationMaster.tsv読み込み")
	var file string = "input/GenerationMaster.tsv"
	if _, err := os.Stat(file); err != nil {
		fmt.Println("ファイルは存在しません！" + file)
		os.Exit(1)
	}
	data, _ := os.Open(file)
	defer data.Close()

	var line string

	scanner := bufio.NewScanner(data)
	// 1行ずつ読み込み
	scanner.Scan()
	// 1行読み込み
	line = sjis_to_utf8(scanner.Text())
	// 1行をタブで分割
	elements := strings.Split(line, "\t")

	fmt.Println("feed_info.txt出力")
	wfile, _ := os.Create("output/feed_info.txt")
	defer wfile.Close()
	var writer *csv.Writer = csv.NewWriter(wfile)
	writer.UseCRLF = true //改行コードを\r\nにする

	// 見出し行を出力
	wdata := []string{
		"feed_start_date",
		"feed_version",
		"feed_start_date",
		"feed_end_date",
	}
	writer.Write(wdata)

	today := time.Now()
	// 月初
	gessho := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.UTC)
	// 1年後の月末
	after1Year := gessho.AddDate(1, 0, -1)

	wdata = []string{
		elements[1],
		elements[2],
		strconv.Itoa(gessho.Year()) + maeZero(strconv.Itoa((int)(gessho.Month())), 2) + "01",
		strconv.Itoa(after1Year.Year()) + maeZero(strconv.Itoa((int)(after1Year.Month())), 2) + maeZero(strconv.Itoa(after1Year.Day()), 2),
	}
	writer.Write(wdata)
	writer.Flush()
}

// 時刻文字列を返す
// 610→10:10:00
// 301→05:01:00
// hが10より小さかったら0を先頭に足す
// mが10より小さかったら0を先頭に足す
func toTime(str string) string {
	var time int
	time, _ = strconv.Atoi(str)
	h := time / 60
	m := time % 60
	hh := fmt.Sprintf("%02d", h)
	mm := fmt.Sprintf("%02d", m)
	hhmm := hh + ":" + mm + ":" + "00"
	//hhmm := strconv.Itoa(hh) + ":" + strconv.Itoa(m)

	return hhmm
}

// 前ゼロ埋め
// str 文字列
// size 桁数
func maeZero(str string, size int) string {
	var len = size - len(str)
	for i := 0; i < len; i++ {
		str = "0" + str
	}
	return str
}

// SJISをUTF8に変換
func sjis_to_utf8(str string) string {
	var iostr = transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder())
	var iobyte, _ = io.ReadAll(iostr)
	return string(iobyte)
}