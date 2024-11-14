package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	fmt.Println("処理開始")

	// inpput/StopMaster.tsvを読み込んで、stops.txtを出力
	readStopMasterTsv()

	fmt.Println("処理終了")
}

// inpput/StopMaster.tsvを読み込んで、stops.txtを出力
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
		line = sjis_to_utf8(scanner.Text())
		elements := strings.Split(line, "\t")
		var stop_id string = elements[1]
		var stop_name string = elements[2]
		var stop_yomi string = elements[3]

		fmt.Printf("%s,%s,%s\n", stop_id, stop_name, stop_yomi)
	}
}

// SJISをUTF8に変換
func sjis_to_utf8(str string) string {
	var iostr = transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder())
	var iobyte, _ = io.ReadAll(iostr)
	return string(iobyte)
}
