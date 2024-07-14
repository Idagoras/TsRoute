package biz

import (
	"fmt"

	"github.com/hbollon/go-edlib"
	"github.com/idagoras/TsRoute/model"
	"golang.org/x/sync/errgroup"
	"gonum.org/v1/gonum/stat"
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
	result := &model.AnalyzeResult{}
	lcssSet := &model.SetResult{}
	RDR := &model.SetResult{}
	RTR := &model.SetResult{}
	RCR := &model.SetResult{}
	g := errgroup.Group{}
	for i := range pairs{
		pair := pairs[i]
		index := i
		g.Go(func() error {
			rep, err := a.lbs.Serve(GaoDeLbsTransportationRouteQuery,&model.GaoDeLbsTransportationRouteQueryLbsRequest{
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

	return result, nil
}


func routeLCSS(r1 *model.Route,r2 *model.Route) float64{
	var lcss int
	var r1Str string
	var r2Str string
	for _, seg := range r1.Transits[0].Segments{
		for _, busline := range seg.Bus.Buslines{
			r1Str += busline.Id
		}
	}
	for _, seg := range r2.Transits[0].Segments{
		for _, busline := range seg.Bus.Buslines{
			r2Str += busline.Id
		}
	}
	lcss = edlib.LCS(r1Str,r2Str)
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
