package biz

import (
	"errors"
	"fmt"

	"github.com/idagoras/TsRoute/model"
	"golang.org/x/sync/errgroup"
)

var(
	ErrGridIsNull = errors.New("grid is null")
)

const(
	ConcernTypeCost = iota
	ConcernTypeDuration
	ConcernTypeTransformNum
	ConcernTypeDistance
)

const(
	GaoDeKey = "39884b28c6fbf518b2c486b60121fb38"
)

type TsRouter struct{
	lbs model.LbsServer
	gridManager *GridPartitionManager
}

func NewTsRouter(lbs model.LbsServer, gridManager *GridPartitionManager) *TsRouter{
	return &TsRouter{
		lbs: lbs,
		gridManager: gridManager,
	}
}

func(r *TsRouter) GetOptimalRoute(epsilon0, epsilon1 float32,origin model.TsPoint,destination model.TsPoint,id string,concernType int)(*model.Route,error){
	var result *model.Route
	var err error
	originGrid := r.gridManager.GetGrid(id,origin)
	if originGrid == nil{
		return nil, ErrGridIsNull
	}
	originPoints, _ := originGrid.GetAllPoints()
	originDistances := originGrid.GetDistances(origin)
	var originHashs []string
	var originFloatValues []float64
	for _, p := range originPoints{
		originHashs = append(originHashs, p.Hash())
		originFloatValues = append(originFloatValues, originDistances[p.Hash()])
	}

	destGrid := r.gridManager.GetGrid(id,destination)
	if destGrid == nil{
		return nil, ErrGridIsNull
	}
	destPoints, _ := destGrid.GetAllPoints()
	destDistances := destGrid.GetDistances(destination)
	
	var destHashs []string
	var destFloatValues []float64
	for _, p := range destPoints{
		destHashs = append(destHashs, p.Hash())
		destFloatValues = append(destFloatValues, destDistances[p.Hash()])
	}
	
	originPb, err := NewPlanerExponentialMechanism(originHashs,originFloatValues,epsilon0)
	if err != nil{
		return nil, err
	}
	destPb, err := NewPlanerExponentialMechanism(destHashs,destFloatValues,epsilon0)
	if err != nil{
		return nil, err
	}

	k := int(epsilon1/epsilon0)
	// k is the number of queries
	var routes []*model.Route
	for i := 0 ; i < k ; i ++ {
		pbOriginHash := originPb.RandomString()
		pbDestHash := destPb.RandomString()
		pbOrigin, _ := originGrid.GetPoint(pbOriginHash)
		pbDestination, _ := destGrid.GetPoint(pbDestHash)
		if pbOriginHash == origin.Hash() && pbDestHash == destination.Hash(){
			rep, err := r.lbs.Serve(int(GaoDeLbsTransportationRouteQuery),model.GaoDeLbsTransportationRouteQueryLbsRequest{
				Key: GaoDeKey,
				Origin: fmt.Sprintf("%f,%f",origin.Lon(),origin.Lat()),
				Destination: fmt.Sprintf("%f,%f",destination.Lon(),destination.Lat()),
				City1:"020",
				City2:"020",
			})
			if err != nil{
				return  nil, err
			}
			rp, _ := rep.(model.GaoDeLbsTransportationRouteQueryResponse)
			result = &rp.Route
			return result, nil
		}
		var originRoute *model.Route
		var destRoute *model.Route
		var route *model.Route

		group := errgroup.Group{}
		group.Go(func() error {
			rep, err := r.lbs.Serve(int(GaoDeLbsTransportationRouteQuery),model.GaoDeLbsTransportationRouteQueryLbsRequest{
				Key: GaoDeKey,
				Origin: fmt.Sprintf("%f,%f",origin.Lon(),origin.Lat()),
				Destination: fmt.Sprintf("%f,%f",pbOrigin.Lon(),pbOrigin.Lat()),
				City1:"020",
				City2:"020",
			})
			if err != nil{
				return  err
			}
			rp, _ := rep.(model.GaoDeLbsTransportationRouteQueryResponse)
			originRoute = &rp.Route
			return nil
		})

		group.Go(func() error{
			rep, err := r.lbs.Serve(int(GaoDeLbsTransportationRouteQuery),model.GaoDeLbsTransportationRouteQueryLbsRequest{
				Key: GaoDeKey,
				Origin: fmt.Sprintf("%f,%f",pbOrigin.Lon(),pbOrigin.Lat()),
				Destination: fmt.Sprintf("%f,%f",pbDestination.Lon(),pbDestination.Lat()),
				City1:"020",
				City2:"020",
			})
			if err != nil{
				return  err
			}
			rp, _ := rep.(model.GaoDeLbsTransportationRouteQueryResponse)
			route = &rp.Route
			return nil
		})

		group.Go(func () error {
			rep, err := r.lbs.Serve(int(GaoDeLbsTransportationRouteQuery),model.GaoDeLbsTransportationRouteQueryLbsRequest{
				Key: GaoDeKey,
				Origin: fmt.Sprintf("%f,%f",pbDestination.Lon(),pbDestination.Lat()),
				Destination: fmt.Sprintf("%f,%f",destination.Lon(),destination.Lat()),
				City1:"020",
				City2:"020",
			})
			if err != nil{
				return  err
			}
			rp, _ := rep.(model.GaoDeLbsTransportationRouteQueryResponse)
			destRoute = &rp.Route
			return nil
		})
		
		err := group.Wait()
		if err != nil{
			return nil, err
		}
		routes = append(routes,originRoute.Join(route).Join(destRoute))
	}

	result = routes[0]
	for _, rt := range routes{
		switch concernType{
		case ConcernTypeCost:
			if rt.GetCost() < result.GetCost() && rt.GetCost() > 0{
				result = rt
			}
		case ConcernTypeDistance:
			if rt.GetDistance() < result.GetDistance() && rt.GetDistance() > 0{
				result = rt
			}
		case ConcernTypeDuration:
			if rt.GetDuration() < result.GetDuration() && rt.GetDistance() > 0{
				result = rt
			}
		case ConcernTypeTransformNum:
			if rt.GetTransformNum() < result.GetTransformNum() && rt.GetTransformNum() > 0{
				result = rt
			}
		}		
	}
	
	return result, err
}
