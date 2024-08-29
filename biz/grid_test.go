package biz

import (
	"testing"
	"fmt"
	"github.com/idagoras/TsRoute/model"
	"github.com/idagoras/TsRoute/store"
	"github.com/stretchr/testify/require"
)

func TestGridManager(t *testing.T){
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
	grid := gridPartitionManager.GetGrid(key,&model.BusStop{
		Name:"华师大南门",Id: "BV10019585",Location: "113.35169,23.13555",
	})
	fmt.Println(grid.GetPointsNum())
	fmt.Printf("tf:%f,%f br:%f,%f\n",grid.TopLeftPoint.Lon(),grid.TopLeftPoint.Lat(),grid.BottomRightPoint.Lon(),grid.BottomRightPoint.Lat())
	fmt.Println(grid.GetDistances(&model.BusStop{
		Name:"华师大南门",Id: "BV10019585",Location: "113.35169,23.13555",
	}))

}