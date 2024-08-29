package biz

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/idagoras/TsRoute/model"
	"github.com/idagoras/TsRoute/store"
	"github.com/stretchr/testify/require"
)

func TestTsRoute(t *testing.T){
	lbs := NewGaoDeLbsServer(5)
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

	gridPartitionManager, err := NewGridPartitionManager(
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
	require.NotNil(t,gridPartitionManager)
	key := "10*10"
	gridPartitionManager.AddPartition(10,10,key)

	tsroute := NewTsRouter(lbs,gridPartitionManager)
	origin1 := model.BusStop{Name:"广东电视台",Id: "BV11009136",Location: "113.282936,23.137556"}
	dest1 := model.BusStop{Name:"华师大南门",Id: "BV10019585",Location: "113.35169,23.13555"}
	route, err := tsroute.GetOptimalRoute(0.01,0.05,&origin1,&dest1,key,ConcernTypeDuration)
	require.NoError(t,err)
	require.NotNil(t,route)

	jsonByte, err := json.Marshal(route)
	require.NoError(t,err)
	fmt.Println(string(jsonByte))
}