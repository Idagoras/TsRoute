package biz

import (
	"fmt"

	"github.com/idagoras/TsRoute/model"
)

type KeyObservationExperiment struct{
	origin model.TsPoint
	destination model.TsPoint
	store model.TsStore
	analyst *KeyObservationAnalyst
	lbs model.LbsServer
	savePath string
	saver model.DataSaver
}

func NewKeyObservationExperiemnt(origin model.TsPoint,destination model.TsPoint,store model.TsStore,analyst *KeyObservationAnalyst,lbs model.LbsServer, savePath string,saver model.DataSaver)*KeyObservationExperiment{
	return &KeyObservationExperiment{
		origin: origin,
		destination: destination,
		store: store,
		analyst: analyst,
		lbs: lbs,
		savePath: savePath,
		saver: saver,
	}
}

func(e *KeyObservationExperiment) Once(num int) error{
	var err error
	index2Routes := make(map[int][]*model.Route)
	originStop , _ := e.origin.(*model.BusStop)

	var baseRoute *model.Route

	repBase, err := e.lbs.Serve(int(GaoDeLbsTransportationRouteQuery),model.GaoDeLbsTransportationRouteQueryLbsRequest{
		Key: GaoDeKey,
		Origin: fmt.Sprintf("%f,%f",e.origin.Lon(),e.origin.Lat()),
		Destination: fmt.Sprintf("%f,%f",e.destination.Lon(),e.destination.Lat()),
		City1:"020",
		City2:"020",
	})
	if err != nil{
		return err
	}
	rpBase, _ := repBase.(model.GaoDeLbsTransportationRouteQueryResponse)
	baseRoute = &rpBase.Route

	for i := 1; i <= num; i ++ {
		stops, err := e.store.ListPointsCanAchieved(originStop.Id,i)
		if err != nil{
			return err
		}
		for _, stop := range stops{
			rep, err := e.lbs.Serve(int(GaoDeLbsTransportationRouteQuery),model.GaoDeLbsTransportationRouteQueryLbsRequest{
				Key: GaoDeKey,
				Origin: fmt.Sprintf("%f,%f",stop.Lon(),stop.Lat()),
				Destination: fmt.Sprintf("%f,%f",e.destination.Lon(),e.destination.Lat()),
				City1:"020",
				City2:"020",
			})
			if err != nil{
				return err
			}
			rp, _ := rep.(model.GaoDeLbsTransportationRouteQueryResponse)
			index2Routes[i] = append(index2Routes[i], &rp.Route)
		}
	}

	result, err := e.analyst.Evaluate(num,baseRoute,index2Routes,map[string]string{
		"origin":fmt.Sprintf("%f,%f",e.origin.Lon(),e.origin.Lat()),
		"destination":fmt.Sprintf("%f,%f",e.destination.Lon(),e.destination.Lat()),
	})
	if err != nil{
		return err
	}
	savePath := e.savePath + "/" +  "keyobservation" + GetNowTimeStampString() + ".csv"
	err = e.saver.Save(savePath,result)
	if err != nil{
		return err
	}
	return err
}