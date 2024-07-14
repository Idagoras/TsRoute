package main

import (
	"github.com/idagoras/TsRoute/biz"
	"github.com/idagoras/TsRoute/model"
	"github.com/idagoras/TsRoute/store"
)

func main(){
	st , err := store.NewTsRouteStore(&model.TsStoreConfig{
		Host: "127.0.0.1",
		User: "idagoras",
		Password: "314159",
		Dbname: "shiftroute",
		Port: 5432,
		SSLMode: "disabled",
		TimeZone: "Asia/Shanghai",
	})
	if err != nil{
		panic("initial store failed")
	}

	gridManager, err := biz.NewGridPartitionManager(
		&model.GeoPoint{
			Longitude: 113.2445,
			Latitude: 23.2077,
		},
		&model.GeoPoint{
			Longitude: 113.4445,
			Latitude: 23.2077,
		},
		st,
	)
	gridManager.AddPartition(10,10,"exp1")
	if err != nil{
		panic("inital gridManager failed")
	}

	lbs := biz.NewGaoDeLbsServer(5)
	tsRoute := biz.NewTsRouter(lbs, gridManager)
	dataLoader := store.NewGuangzhouTsDataLoader()
	dataSaver := store.NewGuangzhouTsDataSaver()
	analyst := biz.NewAnalyst(lbs)
	filePath := "../dataset/guangzhou.csv"
	savePath := "../result"
	experiment := biz.NewTsRouteExperiment(analyst,tsRoute,dataLoader,dataSaver,filePath,savePath)

	err = experiment.Once("exp1",0.01,0.06,biz.ConcernTypeDuration)
	if err != nil{
		panic("experiment run failed")
	}
}