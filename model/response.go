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
	Cost Cost `json:"cost"`
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
	Cost Cost `json:"cost"`
	Steps []WalkingStep  `json:"steps"`
}

type BusSegment struct{
	Buslines []Busline `json:"buslines"`
}

type Segment struct{
	Walking WalingSegment `json:"walking"`
	Bus BusSegment `json:"bus"`
}

type Cost struct{
	Duration string `json:"duration"`
	TransitFee string `json:"transit_fee"`
	TaxiCost string `json:"taxi_cost"`
}

type Transit struct{
	Cost Cost `json:"cost"`
	Nightflag string `json:"nightflag"`
	WalkingDistance string `json:"walking_distance"`
	Distance string `json:"distance"`
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
	if len(r.Transits) == 0 {
		rd, _ := strconv.ParseFloat(r.Distance,64)
		ard, _ := strconv.ParseFloat(ar.Distance,64)
		rt := *ar
		rt.Distance = fmt.Sprintf("%f",rd+ard)
		rt.Origin = r.Origin
		return &rt
	}

	if len(ar.Transits) == 0{
		rd, _ := strconv.ParseFloat(r.Distance,64)
		ard, _ := strconv.ParseFloat(ar.Distance,64)
		rt := *r
		rt.Distance = fmt.Sprintf("%f",rd+ard)
		rt.Destination = ar.Destination
		return &rt
	}

	if r.Destination == ar.Origin{
		result.Origin = r.Origin
		result.Destination = ar.Destination
		rd, _ := strconv.ParseFloat(r.Distance,64)
		ard, _ := strconv.ParseFloat(ar.Distance,64)
		result.Distance = fmt.Sprintf("%f",rd+ard)
		rTrans := r.Transits[0]
		arTrans := ar.Transits[0]
		var rIndex int
		for i := len(rTrans.Segments) - 1; i >= 0 ; i -= 1 {
			l := len(rTrans.Segments[i].Bus.Buslines)
			if l > 0{
				rIndex = i
				break
			}
		}
		var arIndex int
		for i := 0; i < len(arTrans.Segments); i++{
			l := len(arTrans.Segments[i].Bus.Buslines)
			if l > 0{
				arIndex = i
				break
			}
		}

		result.Transits = append(result.Transits, r.Transits[0])
		result.Transits[0].Segments = make([]Segment, 0)
		var publicBusline *Busline

		rTransBuslines := make(map[string]int)
		for i, busline := range rTrans.Segments[rIndex].Bus.Buslines{
			rTransBuslines[busline.Id] = i
		}

		for _, busline := range arTrans.Segments[arIndex].Bus.Buslines{
			if index, ok := rTransBuslines[busline.Id]; ok{
				rBusline := &rTrans.Segments[rIndex].Bus.Buslines[index]
				if rBusline.ArrivalStop.Id == busline.DepartureStop.Id{
					publicBusline = &Busline{
						ArrivalStop: busline.ArrivalStop,
						DepartureStop: rBusline.DepartureStop,
					}
					copy(publicBusline.ViaStops,rBusline.ViaStops)
					rn , _ := strconv.Atoi(rBusline.ViaNum)
					arn , _ := strconv.Atoi(busline.ViaNum)
					publicBusline.ViaNum = strconv.Itoa(rn + arn - 1)
					publicBusline.ViaStops = append(publicBusline.ViaStops,busline.DepartureStop)
					publicBusline.ViaStops = append(publicBusline.ViaStops,busline.ViaStops...)
					result.Transits[0].Segments = append(result.Transits[0].Segments, rTrans.Segments[:rIndex]...)
					segement := Segment{
						Walking: rTrans.Segments[rIndex].Walking,
						Bus: BusSegment{
							Buslines: []Busline{*publicBusline},
						},
					}
					result.Transits[0].Segments = append(result.Transits[0].Segments, segement)
					result.Transits[0].Segments = append(result.Transits[0].Segments, arTrans.Segments[arIndex:]...)
					rwalkingdistance, _ := strconv.ParseInt(rTrans.WalkingDistance,10,64)
					arwalkingdistance, _ := strconv.ParseInt(arTrans.WalkingDistance,10,64)
					result.Transits[0].WalkingDistance = strconv.Itoa(int(rwalkingdistance+arwalkingdistance)) 
					rduration, _ := strconv.ParseInt(rTrans.Cost.Duration,10,64)
					arduration, _ := strconv.ParseInt(arTrans.Cost.Duration,10,64)
					result.Transits[0].Cost.Duration = strconv.Itoa(int(rduration)+int(arduration))
					rcost, _ := strconv.ParseFloat(rTrans.Cost.TransitFee,64)
					arcost,_ := strconv.ParseFloat(arTrans.Cost.TransitFee,64)
					result.Transits[0].Cost.TransitFee = strconv.Itoa(int(rcost)+int(arcost)-2)
					return result
				}
			}
		}

		result.Transits[0].Segments = append(result.Transits[0].Segments, rTrans.Segments[:rIndex]...)
		result.Transits[0].Segments = append(result.Transits[0].Segments, rTrans.Segments[rIndex])
		result.Transits[0].Segments = append(result.Transits[0].Segments, arTrans.Segments[arIndex])
		result.Transits[0].Segments = append(result.Transits[0].Segments, arTrans.Segments[arIndex:]...)
		rcost, _ := strconv.ParseFloat(rTrans.Cost.TransitFee,64)
		arcost,_ := strconv.ParseFloat(arTrans.Cost.TransitFee,64)
		result.Transits[0].Cost.TransitFee = strconv.Itoa(int(rcost)+int(arcost))
		rduration, _ := strconv.ParseInt(rTrans.Cost.Duration,10,64)
		arduration, _ := strconv.ParseInt(arTrans.Cost.Duration,10,64)
		result.Transits[0].Cost.Duration = strconv.Itoa(int(rduration)+int(arduration))
		rwalkingdistance, _ := strconv.ParseInt(rTrans.WalkingDistance,10,64)
		arwalkingdistance, _ := strconv.ParseInt(arTrans.WalkingDistance,10,64)
		result.Transits[0].WalkingDistance = strconv.Itoa(int(rwalkingdistance+arwalkingdistance))
		
	}else{
		return r
	}
	return result
}

func(r *Route)GetCost() float64{
	if len(r.Transits) == 0{
		return 0.0
	}
	cost, _ := strconv.ParseFloat(r.Transits[0].Cost.TransitFee,64)
	return cost
}

func(r *Route)GetDuration() int{
	if len(r.Transits) == 0{
		return 0
	}
	duration, _ := strconv.Atoi(r.Transits[0].Cost.Duration)
	return duration
}

func(r *Route)GetTransformNum() int{
	var num int
	if len(r.Transits) == 0{
		return 0
	}
	for _, seg := range r.Transits[0].Segments{
		if len(seg.Bus.Buslines) >= 1{
			num += 1
		}
	}
	return num
}

func(r *Route)GetDistance() float64{
	if len(r.Transits) == 0{
		return 0.0
	}
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

