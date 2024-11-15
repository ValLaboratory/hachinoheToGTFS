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

// Stop辞書
// キー stop_id
// 値 Stop
var stopMap map[string]Stop = make(map[string]Stop)

func main() {
	fmt.Println("処理開始")

	// inpput/StopMaster.tsvを読み込んで、stopをstopMap辞書に格納
	readStopMasterTsv()

	// stopMap辞書の要素をstops.txtに出力
	writeStopsTxt()

	fmt.Println("処理終了")
}

// inpput/StopMaster.tsvを読み込んで、stopをstopMap辞書に格納
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

// stops.txtを出力
func writeStopsTxt() {
	fmt.Println("stops.txt出力")
	file, _ := os.Create("output/stops.txt")
	defer file.Close()
	var writer *csv.Writer = csv.NewWriter(transform.NewWriter(file, japanese.ShiftJIS.NewEncoder()))
	writer.UseCRLF = true //改行コードを\r\nにする
	data := []string{
		"aaa",
		"bbb",
		"ccc",
	}
	writer.Write(data)
	writer.Flush()
}

// SJISをUTF8に変換
func sjis_to_utf8(str string) string {
	var iostr = transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder())
	var iobyte, _ = io.ReadAll(iostr)
	return string(iobyte)
}
