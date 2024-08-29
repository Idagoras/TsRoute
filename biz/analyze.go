package biz

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hbollon/go-edlib"
	"github.com/idagoras/TsRoute/model"
	"golang.org/x/sync/errgroup"
	"gonum.org/v1/gonum/stat"
	log "github.com/sirupsen/logrus"
)



type Analyst struct{
	lbs model.LbsServer
}

func NewAnalyst(lbs model.LbsServer)*Analyst{
	return &Analyst{
		lbs: lbs,
	}
}



func(a *Analyst)Evaluate(metaData map[string]string,pairs []*model.StopPair, guessedOptimalRoutes []*model.Route)(*model.AnalyzeResult,error){
	log.Info("begin evaluate")
	result := &model.AnalyzeResult{}
	lcssSet := model.NewSetResult()
	RDR := model.NewSetResult()
	RTR := model.NewSetResult()
	RCR := model.NewSetResult()
	g := errgroup.Group{}
	for i := range pairs{
		pair := pairs[i]
		index := i
		g.Go(func() error {
			rep, err := a.lbs.Serve(GaoDeLbsTransportationRouteQuery,model.GaoDeLbsTransportationRouteQueryLbsRequest{
				Key: GaoDeKey,
				Origin: fmt.Sprintf("%f,%f",pair.Origin.Lon(),pair.Origin.Lat()),
				Destination: fmt.Sprintf("%f,%f",pair.Destination.Lon(),pair.Destination.Lat()),
				City1:"020",
				City2:"020",
			})
			if err != nil{
				return err
			}
			rp, _ := rep.(model.GaoDeLbsTransportationRouteQueryResponse)
			route := rp.Route
			key := fmt.Sprintf("%f,%f,%f,%f",pair.Origin.Lon(),pair.Origin.Lat(),pair.Destination.Lon(),pair.Destination.Lat())
			jsonByte , _ := json.Marshal(route)
			fmt.Println(string(jsonByte))
			lcssSet.Vals.Store(key,routeLCSS(&route,guessedOptimalRoutes[index]))
			RDR.Vals.Store(key,routeRDR(&route,guessedOptimalRoutes[index]))
			RTR.Vals.Store(key,routeRTR(&route,guessedOptimalRoutes[index]))
			RCR.Vals.Store(key,routeRCR(&route,guessedOptimalRoutes[index]))
			return nil
		})
	}
	err := g.Wait()
	if err != nil{
		return nil, err
	}
	var lcssVals []float64
	var rdrVals []float64
	var rtrVals []float64
	var rcrVals []float64

	lcssSet.Vals.Range(func(key, value any) bool {
		floatval, _ := value.(float64)
		lcssVals = append(lcssVals, floatval)
		return true
	})
	RDR.Vals.Range(func(key, value any) bool {
		floatval, _ := value.(float64)
		rdrVals = append(rdrVals, floatval)
		return true
	})
	RTR.Vals.Range(func(key, value any) bool {
		floatval, _ := value.(float64)
		rtrVals = append(rtrVals, floatval)
		return true
	})
	RCR.Vals.Range(func(key, value any) bool {
		floatval, _ := value.(float64)
		rcrVals = append(rcrVals, floatval)
		return true
	})

	lcssSet.Mean = stat.Mean(lcssVals,nil)
	lcssSet.Variance = stat.Variance(lcssVals,nil)
	lcssSet.StdDev = stat.StdDev(lcssVals,nil)

	RDR.Mean = stat.Mean(rdrVals,nil)
	RDR.Variance = stat.Variance(rdrVals,nil)
	RDR.StdDev = stat.StdDev(rdrVals,nil)

	RTR.Mean = stat.Mean(rtrVals,nil)
	RTR.Variance = stat.Variance(rtrVals,nil)
	RTR.StdDev = stat.StdDev(rtrVals,nil)

	RCR.Mean = stat.Mean(rcrVals,nil)
	RCR.Variance = stat.Variance(rcrVals,nil)
	RCR.StdDev = stat.StdDev(rcrVals,nil)

	result.Similarity = lcssSet
	result.RDR = RDR
	result.RTR = RTR
	result.RCR = RCR
	result.Num = len(lcssVals)
	result.Pairs = pairs
	result.MetaData = metaData
	log.Info("end evaluate")
	return result, nil
}


func routeLCSS(r1 *model.Route,r2 *model.Route) float64{
	var lcss int
	var r1Str string
	var r2Str string

	if len(r1.Transits) == 0 && len(r2.Transits) == 0{

		return 1.0
	}
	if len(r1.Transits) == 0 || len(r2.Transits) == 0{
		return 0.0
	}
	for _, seg := range r1.Transits[0].Segments{
		if len(seg.Bus.Buslines) > 0{
			r1Str += seg.Bus.Buslines[0].DepartureStop.Id
			for _, stop := range seg.Bus.Buslines[0].ViaStops{
				r1Str+=stop.Id
			}
			r1Str += seg.Bus.Buslines[0].ArrivalStop.Id
		}
	}
	for _, seg := range r2.Transits[0].Segments{
		if len(seg.Bus.Buslines) > 0{
			r2Str += seg.Bus.Buslines[0].DepartureStop.Id
			for _, stop := range seg.Bus.Buslines[0].ViaStops{
				r2Str+=stop.Id
			}
			r2Str += seg.Bus.Buslines[0].ArrivalStop.Id
		}
	}
	lcss = edlib.LCS(r1Str,r2Str)
	//fmt.Println(r1Str,r2Str)
	return float64(lcss)/float64(len(r2Str))
}

func routeRDR(r1 *model.Route,r2 *model.Route) float64{
	return float64(r1.GetDuration())/float64(r2.GetDuration())
}

func routeRTR(r1 *model.Route,r2 *model.Route) float64{
	return float64(r1.GetTransformNum())/float64(r2.GetTransformNum())
}

func routeRCR(r1 *model.Route,r2 *model.Route) float64{
	return float64(r1.GetCost())/float64(r2.GetCost())
}

type GridAnalyst struct{

}

func NewGridAnalyst() *GridAnalyst{
	return &GridAnalyst{

	}
}

func(g *GridAnalyst)Evaluate(metaData map[string]string,girds []*model.TsGrid,lineMap map[string][]string,edgeMap map[string][]*model.TsEdge)(*model.GridAnalyzeResult, error){
	var result *model.GridAnalyzeResult = &model.GridAnalyzeResult{
		Num: len(girds),
		Grids: girds,
		MetaData: metaData,
	}
	var err error
	lineSet := model.NewSetResult()
	edgeSet := model.NewSetResult()
	pointSet := model.NewSetResult()

	for _, grid := range girds{
		tf := fmt.Sprintf("%f,%f",grid.TopLeftPoint.Lon(),grid.TopLeftPoint.Lat())
		br := fmt.Sprintf("%f,%f",grid.BottomRightPoint.Lon(),grid.BottomRightPoint.Lat())
		key := tf + ":" + br
		lines := lineMap[key]
		edges := edgeMap[key]
		lineSet.Vals.Store(key,len(lines))
		edgeSet.Vals.Store(key,len(edges))
		num := grid.GetPointsNum()
		pointSet.Vals.Store(key,num)
	}
	
	result.Points = pointSet
	result.Lines = lineSet
	result.Edges = edgeSet

	var lineNums []float64
	var edgeNums []float64
	var pointNums []float64

	result.Lines.Vals.Range(func(key, value any) bool {
		intVal, _ := value.(int)
		lineNums = append(lineNums, float64(intVal))
		return true
	})

	result.Edges.Vals.Range(func(key, value any) bool {
		intVal, _ := value.(int)
		edgeNums = append(edgeNums, float64(intVal))
		return true
	})

	result.Points.Vals.Range(func(key, value any) bool {
		intVal, _ := value.(int)
		pointNums = append(pointNums, float64(intVal))
		return true
	})

	result.Lines.Mean = stat.Mean(lineNums,nil)
	result.Lines.Variance = stat.Variance(lineNums,nil)
	result.Lines.StdDev = stat.StdDev(lineNums,nil)

	result.Edges.Mean = stat.Mean(edgeNums,nil)
	result.Edges.Variance = stat.Variance(edgeNums,nil)
	result.Edges.StdDev = stat.StdDev(edgeNums,nil)

	result.Points.Mean =stat.Mean(pointNums,nil)
	result.Points.Variance = stat.Variance(pointNums,nil)
	result.Points.StdDev = stat.StdDev(pointNums,nil)

	return result,err
}

type KeyObservationAnalyst struct{

}

func NewKeyObservationAnalyst() *KeyObservationAnalyst{
	return &KeyObservationAnalyst{
	}
}

func(a *KeyObservationAnalyst) Evaluate(num int,baseRoute *model.Route,numToRoutes map[int][]*model.Route,metaData map[string]string)(*model.KeyObservationResult,error){
	var result *model.KeyObservationResult = &model.KeyObservationResult{
		MetaData: metaData,
	}
	i2Sim := model.NewMapResult()
	for i := 1; i <= num; i++{
		var vals []float64
		for _, route := range numToRoutes[i]{
			lcss := routeLCSS(baseRoute,route)
			vals = append(vals, lcss)
		}
		i2Sim.KV[strconv.Itoa(i)] = vals
	}
	result.Similarity = i2Sim
	return 	result, nil
}  