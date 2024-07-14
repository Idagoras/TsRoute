package biz

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/idagoras/TsRoute/model"
)

type TsRouteExperiment struct{
	analyst *Analyst
	tsRoute *TsRouter
	dataLoader model.DataLoader
	dataSaver model.DataSaver 
	readPath string
	savePath string
}

func NewTsRouteExperiment(
	analyst *Analyst,
	tsRoute *TsRouter,
	dataLoader model.DataLoader,
	dataSaver model.DataSaver,
	readPath string,
	savePath string,
) *TsRouteExperiment{
	return &TsRouteExperiment{
		analyst: analyst,
		tsRoute: tsRoute,
		dataLoader: dataLoader,
		dataSaver: dataSaver,
		readPath: readPath,
		savePath: savePath,
	}
}

func(e *TsRouteExperiment) Once(id string, epsilon_small,epsilon_large float32,concernType int) error{
	data, err := e.dataLoader.Load(e.readPath,0,1000)
	if err != nil{
		return err
	}
	var pairs []*model.StopPair
	var routes []*model.Route
	metaData := map[string]string{
		"epsilon0":fmt.Sprintf("%f",epsilon_small),
		"epsilon1":fmt.Sprintf("%f",epsilon_large),
		"k":fmt.Sprintf("%d",int(epsilon_large/epsilon_small)),
		"tf":fmt.Sprintf("%f:%f",e.tsRoute.gridManager.areaTopLeftPoint.Lon(),e.tsRoute.gridManager.areaTopLeftPoint.Lat()),
		"br":fmt.Sprintf("%f:%f",e.tsRoute.gridManager.areaBottomRightPoint.Lon(),e.tsRoute.gridManager.areaBottomRightPoint.Lat())	,
		"pt":id,
	}
	for _, dt := range data{
		pair, ok := dt.(model.StopPair)
		pairs = append(pairs, &pair)
		if !ok{
			return errors.New("read data failed,type error")
		}
		origin := pair.Origin
		destination := pair.Destination
		guessedRoute, err := e.tsRoute.GetOptimalRoute(
			epsilon_small,
			epsilon_large,
			&origin,
			&destination,
			id,
			concernType,
		)
		if err != nil{
			return err
		}
		routes = append(routes, guessedRoute)
	}

	analyseResult, err := e.analyst.Evaluate(metaData,pairs,routes)
	if err != nil{
		return err
	}
	savePath := e.savePath + "/" + GetNowTimeStampString() 
	err = e.dataSaver.Save(savePath,analyseResult)
	if err != nil{
		return nil
	}
	return nil
}

func(e *TsRouteExperiment)MulipleTimes(k int,f func(int)(float32,float32,string,int)) error{
	for i := 0 ; i < k ; k ++{
		epsilon_0,epsilon_1,id,concernType := f(i)
		err := e.Once(id,epsilon_0,epsilon_1,concernType)
		if err != nil{
			return err
		}
	}
	return nil
}

func GetNowTimeStampString() string{
	t := time.Now()
	nanosecond := t.UnixMilli()
	timeStr := strconv.Itoa(int(nanosecond))
	return timeStr
}