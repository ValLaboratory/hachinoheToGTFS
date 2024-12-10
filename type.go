package main

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
	name    string
}

// 系統
type Route struct {
	id       string
	name     string
	stop_ids []string
	ikisaki  string
}

// 便
type Trip struct {
	id        string
	route_id  string
	yobi      string
	stopTimes []StopTime
	col3      string // 3:通常  2:列車備考(MarkMaster)
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
var routeMap map[string]*Route = make(map[string]*Route)

// Trip配列
var tripList []Trip
