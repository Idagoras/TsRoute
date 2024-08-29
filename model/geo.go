package model

import (
	"errors"
	"math"

	log "github.com/sirupsen/logrus"
	"github.com/dominikbraun/graph"
)

// 23.0077-23.2077 113.2445 113.4445

var(
	ErrPointNotInGrid = errors.New("point should not be in this grid")
)

type GeoPoint struct{
	Longitude float32
	Latitude float32
}

func(p *GeoPoint)Lon() float32{
	return p.Longitude
}

func(p *GeoPoint)Lat() float32{
	return p.Latitude
}

func(p *GeoPoint)Hash() string{
	return ""
}

type TsPoint interface{
	Lon() float32
	Lat() float32
	Hash() string
}

type TsRoute interface{
	Lines() []string
	Distance() float32
	GetDuration(lineId string) int
	GetDistance(lineId string) float32
	GetCost(lineId string) int
	Origin() TsPoint
	Destination() TsPoint
}

type TsGrid struct{
	TopLeftPoint TsPoint
	BottomRightPoint TsPoint
	g graph.Graph[string,TsPoint]
	points map[string]TsPoint
}

func NewTsGrid(topleft , bottomRight TsPoint) *TsGrid{
	g := graph.New(func(p TsPoint)string{
		return p.Hash()
	})
	return &TsGrid{
		TopLeftPoint: topleft,
		BottomRightPoint: bottomRight,
		g: g,
		points: make(map[string]TsPoint),
	}
}

func (gd *TsGrid)AddPoint(p TsPoint) error{
	if _, ok := gd.points[p.Hash()]; ok{
		return nil
	}
	if Between(gd.TopLeftPoint,gd.BottomRightPoint,p){
		err := gd.g.AddVertex(p)
		if err != nil{
			return err
		}
		gd.points[p.Hash()] = p
		return nil
	}else{
		return ErrPointNotInGrid
	}

}

func (gd *TsGrid)AddRoute(r TsRoute) error{
	return gd.g.AddEdge(r.Origin().Hash(),r.Destination().Hash(),graph.EdgeData(r),graph.EdgeWeight(int(1)))
}

func (gd *TsGrid)GetPoint(id string)(TsPoint, error){
	return gd.g.Vertex(id)
}

func (gd *TsGrid)GetAllPoints()([]TsPoint,error){
	var result []TsPoint
	for _, value := range gd.points{
		result = append(result, value)
	}
	return result,nil
}

func (gd *TsGrid)GetPointsNum()int{
	return len(gd.points)
}

func (gd *TsGrid)GetDistances(p TsPoint)map[string]float64{
	distanceMap := make(map[string]float64)
	for key, value := range gd.points{
		distanceMap[key] = GetDistance(float64(p.Lat()),float64(value.Lat()),float64(p.Lon()),float64(value.Lon()))
	}
	return distanceMap
}

func GetDistance(lat1, lat2, lng1, lng2 float64) float64 {
	if lat1 == lat2 && lng1 == lng2{
		return 0.0
	}
    radius := 6371000.0 //6378137.0
	LON1 := lng1
	LON2 := lng2
	LAT1 := lat1
	LAT2 := lat2
    rad := math.Pi / 180.0
    lat1 = lat1 * rad
    lng1 = lng1 * rad
    lat2 = lat2 * rad
    lng2 = lng2 * rad
    theta := lng2 - lng1
    dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	if math.IsNaN(dist){
		log.WithFields(log.Fields{
			"lat1":LAT1,
			"lon1":LON1,
			"lat2":LAT2,
			"lon2":LON2,
		}).Error("get NaN")
	}
    return dist * radius / 1000
}

func Between(topLeft TsPoint,bottomRight TsPoint,point TsPoint) bool{
	tf_lat := topLeft.Lat()
	tf_lon := topLeft.Lon()
	br_lat := bottomRight.Lat()
	br_lon := bottomRight.Lon()
	lon := point.Lon()
	lat := point.Lat()
	return lon <= br_lon && lon >= tf_lon && lat <= tf_lat && lat >= br_lat
}