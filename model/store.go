package model

import "time"

type TsStoreConfig struct{
	Host string
	User string
	Password string
	Dbname string
	Port int
	SSLMode string
	TimeZone string
}

type TsStore interface{
	ListPointsBetween(tfLon,tfLat,brLon,brLat float32)([]*TsStop, error)
	GetOptimalRoute(origin string,destination string)(*TsOptimalRoute,error)
	InsertOptimalRoute(route *TsOptimalRoute) error
	UpdateOptimalRoute(route *TsOptimalRoute) error
	ListPointsCanDirectedAchieved(id string)([]*TsEdge, error)
	ListPointsCanAchieved(id string, num int)([]*TsStop,error)
}

type TsStop struct{
	Id string `json:"id" gorm:"column:id;primaryKey"`
	StopName string `json:"stop_name" gorm:"column:stop_name"`
	Longitude float64 `json:"longitude" gorm:"column:longitude"`
	Latitude float64 `json:"latitude" gorm:"column:latitude"`
	LineIdToSequence string `json:"line_id_to_sequence" gorm:"column:line_id_to_sequence"`
} 

func(s *TsStop) Lon() float32{
	return float32(s.Longitude)
}

func(s *TsStop) Lat() float32{
	return float32(s.Latitude)
}

func(s *TsStop) Hash() string{
	return s.StopName
}

func(s *TsStop)TableName() string{
	return "guangzhou_stop"
}

type TsLine struct{
	Id string `json:"id" gorm:"column:id;primaryKey"`
	LineName string `json:"line_name" gorm:"column:line_name"`
	StopNum int `json:"stop_num" gorm:"column:stop_num"`
	StartTime time.Time `json:"start_time" gorm:"column:start_time"`
	EndTime time.Time `json:"end_time" gorm:"column:end_time"`
	VehicleType int `json:"vehicle_type" gorm:"column:vehicle_type"`
	Polyline string `json:"polyline" gorm:"column:polyline"`
	BasicPrice int `json:"basic_price" gorm:"column:basic_price"`
	TotalPrice int `json:"total_price" gorm:"column:total_price"`
	Direc string `json:"direc" gorm:"column:direc"`
}

type TsEdge struct{
	StopId string `json:"stopId" gorm:"column:stop_id"`
	ToStopId string `json:"toStopId" gorm:"column:to"`
	LineId string `json:"lineId" gorm:"column:line_id"`
}


func(l *TsLine)TableName() string{
	return "guangzhou_vehicle"
}


type TsOptimalRoute struct{
	Origin string `json:"origin"`
	Destination string `json:"destination"`
	Route string `json:"route"`
}
