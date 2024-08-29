package biz

import (
	"encoding/json"
	"errors"
	"fmt"

	"strconv"
	"strings"
	"time"

	"github.com/idagoras/TsRoute/model"
	log "github.com/sirupsen/logrus"
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
	data, err := e.dataLoader.Load(e.readPath,0,20)
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
		if !(model.Between(e.tsRoute.gridManager.areaTopLeftPoint,e.tsRoute.gridManager.areaBottomRightPoint,&pair.Origin) && model.Between(e.tsRoute.gridManager.areaTopLeftPoint,e.tsRoute.gridManager.areaBottomRightPoint,&pair.Destination)){
			continue
		}
		pairs = append(pairs, &pair)
		if !ok{
			return errors.New("read data failed,type error")
		}
		log.WithFields(log.Fields{"origin":pair.Origin,"destination":pair.Destination}).Info("Begin compute")
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
		jsonByte, _ :=json.Marshal(guessedRoute)
		fmt.Println(string(jsonByte))
		if err != nil{
			if strings.Contains(err.Error(), "Client.Timeout exceeded") {
				log.Error("request time out")
				time.Sleep(2 * time.Second)
				continue
			}
			log.WithFields(log.Fields{"origin":pair.Origin,"destination":pair.Destination}).Info("Compute Get Error")
			//return err
			continue
		}
		log.WithFields(log.Fields{"origin":pair.Origin,"destination":pair.Destination}).Info("Compute Success")
		routes = append(routes, guessedRoute)
		//time.Sleep(1*time.Second)
	}

	analyseResult, err := e.analyst.Evaluate(metaData,pairs,routes)
	if err != nil{
		log.WithFields(
			log.Fields{
				"msg":"analyse failed",
			},
		).Error(err)
		return err
	}
	savePath := e.savePath + "/" + GetNowTimeStampString() + ".csv"
	//fmt.Println(savePath)
	err = e.dataSaver.Save(savePath,analyseResult)
	if err != nil{
		log.WithFields(
			log.Fields{
				"msg":"saved failed",
			},
		).Error(err)
		return nil
	}
	return nil
}

func(e *TsRouteExperiment)MulipleTimes(k int,f func(int)(float32,float32,string,int)) error{
	for i := 0 ; i < k ; i ++{
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



