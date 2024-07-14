package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type BusStop struct{
	Name string `json:"name"`
	Id string `json:"id"`
	Location string `json:"location"`
}

func(s *BusStop) Lon() float32{
	vals := strings.Split(s.Location,",")
	lon, _ := strconv.ParseFloat(vals[0],32)
	return float32(lon)
}

func(s *BusStop) Lat() float32{
	vals := strings.Split(s.Location,",")
	lat, _ := strconv.ParseFloat(vals[1],32)
	return float32(lat)
}

func(s *BusStop) Hash() string{
	return s.Name
}

type Busline struct{
	DepartureStop BusStop `json:"departure_stop"`
	ArrivalStop BusStop `json:"arrival_stop"`
	Name string `json:"name"`
	Id string `json:"id"`
	Type string `json:"type"`
	Distance string `json:"distance"`
	BusTimeTips string `json:"bus_time_tips"`
	BusTimeTag string `json:"bustimetag"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
	ViaNum string `json:"via_num"`
	ViaStops []BusStop `json:"via_stops"`
}

type WalkingStep struct{
	Instruction string `json:"instruction"`
	Road string `json:"road"`
	Distance string `json:"distance"`
}

type WalingSegment struct{
	Destination string `json:"destination"`
	Distance string `json:"distance"`
	Origin string `json:"origin"`
	Steps []WalkingStep  `json:"steps"`
}

type BusSegment struct{
	Buslines []Busline `json:"buslines"`
}

type Segment struct{
	Walking WalingSegment `json:"walking"`
	Bus BusSegment `json:"bus"`
}

type Transit struct{
	Cost string `json:"cost"`
	Duration string `json:"duration"`
	Nightflag string `json:"nightflag"`
	WalkingDistance string `json:"walking_distance"`
	Distance string `json:"distance"`
	Missed string `json:"missed"`
	Segments []Segment `json:"segments"`
}

type Route struct{
	Origin string `json:"origin"`
	Destination string `json:"destination"`
	Distance string `json:"distance"`
	TaxiCost string `json:"taxi_cost"`
	Transits []Transit `json:"transits"`
}

func(r *Route)ToDSModel() *TsOptimalRoute{
	jsonByte, _ := json.Marshal(r)
	return &TsOptimalRoute{
		Origin: r.Origin,
		Destination: r.Destination,
		Route: string(jsonByte),
	}
}

func(r *Route)Join(ar *Route) *Route{
	var result *Route = &Route{}
	if r.Origin == ar.Origin && r.Destination == ar.Destination{
		return r
	}
	if r.Destination == ar.Origin{
		rd, _ := strconv.ParseFloat(r.Distance,64)
		ard, _ := strconv.ParseFloat(ar.Distance,64)
		result.Distance = fmt.Sprintf("%f",rd+ard)
		rTrans := r.Transits[0]
		arTrans := ar.Transits[1]
		var rLastBusLine *Busline
		var rIndex int
		for i := len(rTrans.Segments) - 1; i >= 0 ; i -= 1 {
			l := len(rTrans.Segments[i].Bus.Buslines)
			if l > 0{
				rLastBusLine = &rTrans.Segments[i].Bus.Buslines[l - 1]
				rIndex = i
				break
			}
		}
		var arFirstBusline *Busline
		var arIndex int
		for i := 0; i < len(arTrans.Segments); i++{
			l := len(arTrans.Segments[i].Bus.Buslines)
			if l > 0{
				arFirstBusline = &arTrans.Segments[i].Bus.Buslines[0]
				arIndex = i
				break
			}
		}

		result.Transits = append(result.Transits, r.Transits[0])
		result.Transits[0].Segments = make([]Segment, 0)
		var busline Busline
		if rLastBusLine.Id == arFirstBusline.Id && rLastBusLine.ArrivalStop.Id == arFirstBusline.Id{
			busline = *rLastBusLine
			busline.ArrivalStop = arFirstBusline.ArrivalStop
			rn , _ := strconv.Atoi(rLastBusLine.ViaNum)
			arn , _ := strconv.Atoi(arFirstBusline.ViaNum)
			busline.ViaNum = strconv.Itoa(rn + arn - 1)
			busline.ViaStops = append(busline.ViaStops,arFirstBusline.DepartureStop)
			busline.ViaStops = append(busline.ViaStops,arFirstBusline.ViaStops...)
			result.Transits[0].Segments = append(result.Transits[0].Segments, rTrans.Segments[:rIndex]...)
			segement := Segment{
					Walking: rTrans.Segments[rIndex].Walking,
					Bus: BusSegment{
						Buslines: []Busline{busline},
					},
			}
			result.Transits[0].Segments = append(result.Transits[0].Segments, segement)
			result.Transits[0].Segments = append(result.Transits[0].Segments, arTrans.Segments[arIndex:]...)
			rwalkingdistance, _ := strconv.ParseInt(rTrans.WalkingDistance,10,64)
			arwalkingdistance, _ := strconv.ParseInt(arTrans.WalkingDistance,10,64)
			result.Transits[0].WalkingDistance = strconv.Itoa(int(rwalkingdistance+arwalkingdistance)) 
			rduration, _ := strconv.ParseInt(rTrans.Duration,10,64)
			arduration, _ := strconv.ParseInt(rTrans.Duration,10,64)
			result.Transits[0].Duration = strconv.Itoa(int(rduration)+int(arduration))
			rcost, _ := strconv.ParseInt(rTrans.Cost,10,64)
			arcost,_ := strconv.ParseInt(arTrans.Cost,10,64)
			result.Transits[0].Cost = strconv.Itoa(int(rcost)+int(arcost)-2)
		}else{
			result.Transits[0].Segments = append(result.Transits[0].Segments, rTrans.Segments[:rIndex]...)
			result.Transits[0].Segments = append(result.Transits[0].Segments, rTrans.Segments[rIndex])
			result.Transits[0].Segments = append(result.Transits[0].Segments, arTrans.Segments[arIndex])
			result.Transits[0].Segments = append(result.Transits[0].Segments, arTrans.Segments[arIndex:]...)
			rcost, _ := strconv.ParseInt(rTrans.Cost,10,64)
			arcost,_ := strconv.ParseInt(arTrans.Cost,10,64)
			result.Transits[0].Cost = strconv.Itoa(int(rcost)+int(arcost))
			rduration, _ := strconv.ParseInt(rTrans.Duration,10,64)
			arduration, _ := strconv.ParseInt(rTrans.Duration,10,64)
			result.Transits[0].Duration = strconv.Itoa(int(rduration)+int(arduration))
			rwalkingdistance, _ := strconv.ParseInt(rTrans.WalkingDistance,10,64)
			arwalkingdistance, _ := strconv.ParseInt(arTrans.WalkingDistance,10,64)
			result.Transits[0].WalkingDistance = strconv.Itoa(int(rwalkingdistance+arwalkingdistance))
		}
	}else{
		return r
	}
	return result
}

func(r *Route)GetCost() int{
	cost, _ := strconv.Atoi(r.Transits[0].Cost)
	return cost
}

func(r *Route)GetDuration() int{
	duration, _ := strconv.Atoi(r.Transits[0].Duration)
	return duration
}

func(r *Route)GetTransformNum() int{
	var num int
	for _, seg := range r.Transits[0].Segments{
		num += len(seg.Bus.Buslines)
	}
	return num
}

func(r *Route)GetDistance() float64{
	distance, _ := strconv.ParseFloat(r.Transits[0].Distance,64)
	return distance
}

type GaoDeLbsTransportationRouteQueryResponse struct{
	Status string `json:"status"`
	Info string	`json:"info"`
	InfoCode string `json:"infocode"`
	Count string `json:"count"`
	Route Route `json:"route"`
}

