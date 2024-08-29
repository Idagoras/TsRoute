package biz

import (
	"testing"

	"github.com/idagoras/TsRoute/model"
	"github.com/idagoras/TsRoute/store"
	"github.com/stretchr/testify/require"
)

func TestTsRouteExperiment(t *testing.T){
	st , err := store.NewTsRouteStore(&model.TsStoreConfig{
		Host: "127.0.0.1",
		User: "idagoras",
		Password: "314159",
		Dbname: "shiftroute",
		Port: 5432,
		SSLMode: "disable",
		TimeZone: "Asia/Shanghai",
	})
	require.NoError(t,err)

	gridManager, err := NewGridPartitionManager(
		&model.GeoPoint{
			Longitude: 113.2445,
			Latitude: 23.2077,
		},
		&model.GeoPoint{
			Longitude: 113.4445,
			Latitude: 23.0077,
		},
		st,
	)
	gridManager.AddPartition(10,10,"exp1")
	require.NoError(t,err)

	lbs := NewGaoDeLbsServer(10)
	tsRoute := NewTsRouter(lbs, gridManager)
	dataLoader := store.NewGuangzhouTsDataLoader()
	dataSaver := store.NewGuangzhouTsDataSaver()
	analyst := NewAnalyst(lbs)
	filePath := "../dataset/gz200.csv"
	savePath := "../result/same"
	experiment := NewTsRouteExperiment(analyst,tsRoute,dataLoader,dataSaver,filePath,savePath)
	err = experiment.Once("exp1",1,5,ConcernTypeTransformNum)
	require.NoError(t,err)
}

func TestGridExperiment(t *testing.T){
	analyst := NewGridAnalyst()
	saver := store.NewGridDataSaver()
	st , err := store.NewTsRouteStore(&model.TsStoreConfig{
		Host: "127.0.0.1",
		User: "idagoras",
		Password: "314159",
		Dbname: "shiftroute",
		Port: 5432,
		SSLMode: "disable",
		TimeZone: "Asia/Shanghai",
	})
	require.NoError(t,err)
	gridManager, err := NewGridPartitionManager(
		&model.GeoPoint{
			Longitude: 113.2445,
			Latitude: 23.2077,
		},
		&model.GeoPoint{
			Longitude: 113.4445,
			Latitude: 23.0077,
		},
		st,
	)
	require.NoError(t,err)
	savePath := "../result"
	exp := NewGridExperiment(analyst,gridManager,saver,savePath)
	keys := []string{"px1","px2","px3","px4","px5","px6","px7","px8","px9"}
	xs := []int{10,12,13,14,15,16,17,18,19}
	ys := []int{10,12,13,14,15,16,17,18,19}
	err = exp.MultipleTimes(1,func(i int) (string, int, int) {
		return keys[i],xs[i],ys[i]
	})
	require.NoError(t,err)
}

func TestKeyObservationExperiment(t *testing.T){
	st , err := store.NewTsRouteStore(&model.TsStoreConfig{
		Host: "127.0.0.1",
		User: "idagoras",
		Password: "314159",
		Dbname: "shiftroute",
		Port: 5432,
		SSLMode: "disable",
		TimeZone: "Asia/Shanghai",
	})
	require.NoError(t,err)
	saver := store.NewKeyObservationSaver()
	lbs := NewGaoDeLbsServer(10)
	savePath := "../result"
	origin := model.BusStop{Name:"广东电视台",Id: "BV11009136",Location: "113.282936,23.137556"}
	destination := model.BusStop{Name:"华师大南门",Id: "BV10019585",Location: "113.35169,23.13555"}
	analyst := NewKeyObservationAnalyst()
	exp := NewKeyObservationExperiemnt(&origin,&destination,st,analyst,lbs,savePath,saver)
	err = exp.Once(3)
	require.NoError(t,err)
}

func TestLocationPrivacyExperiment(t *testing.T) {
	st , err := store.NewTsRouteStore(&model.TsStoreConfig{
		Host: "127.0.0.1",
		User: "idagoras",
		Password: "314159",
		Dbname: "shiftroute",
		Port: 5432,
		SSLMode: "disable",
		TimeZone: "Asia/Shanghai",
	})
	require.NoError(t,err)
	gridManager, err := NewGridPartitionManager(
		&model.GeoPoint{
			Longitude: 113.2445,
			Latitude: 23.2077,
		},
		&model.GeoPoint{
			Longitude: 113.4445,
			Latitude: 23.0077,
		},
		st,
	)
	gridManager.AddPartition(20,20,"exp1")
	require.NoError(t,err)
	lpc := &LocationPrivacyCalculaterImpl{}
	e := NewLocationPrivacyExperiment(lpc,gridManager)
	epsilons := []float64{0,0.1,0.5,1,2,5,10,15,20}
	for i := 0; i < len(epsilons); i++{
		err = e.Once("exp1",float32(epsilons[i]))
	}
	require.NoError(t,err)
}