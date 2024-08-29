package store

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/idagoras/TsRoute/model"
	//"github.com/peterli110/discreteprobability"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T){
	store, err := NewTsRouteStore(&model.TsStoreConfig{
		Host: "127.0.0.1",
		User: "idagoras",
		Password: "314159",
		Dbname: "shiftroute",
		Port: 5432,
		SSLMode: "disable",
		TimeZone: "Asia/Shanghai",
	})
	require.NoError(t,err)
	stops, err := store.ListPointsBetween(113.2445,23.2077,113.4445,23.0077)
	require.NoError(t,err)
	file, err := os.Create("../result/points.csv")
	require.NoError(t,err)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, st := range stops{
		writer.Write([]string{
			fmt.Sprintf("%f",st.Longitude),
			fmt.Sprintf("%f",st.Latitude),
			st.Id,
			st.StopName,
		})
	}
	fmt.Println(len(stops))
	require.NoError(t,err)

	/*
	st_map := make(map[string]*model.TsStop)
	keys := make([]string,0)
	vals := make([]float64,0)
	for _, stop := range stops{
		st_map[stop.Id] = stop
		keys = append(keys, stop.Id)
		vals = append(vals, 1/float64(len(stops)))
	}
	rng, err := discreteprobability.New(keys,vals)
	require.NoError(t,err)

	file, err = os.Create("../dataset/gz100.200.csv")
	require.NoError(t,err)
	defer file.Close()
	writer = csv.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < 200; i ++{
		str1 := rng.RandomString()
		str2 := rng.RandomString()
		if str1 == str2{
			i -= 1
			continue
		}
		st1 := st_map[str1]
		st2 := st_map[str2]
		err = writer.Write([]string{
			st1.Id,
			st2.Id,
			fmt.Sprintf("%f",st1.Longitude),
			fmt.Sprintf("%f",st1.Latitude),
			fmt.Sprintf("%f",st2.Longitude),
			fmt.Sprintf("%f",st2.Latitude),
			st1.StopName,
			st2.StopName,
		})
		require.NoError(t,err)
	}*/

}

func TestListStopsNearBy(t *testing.T){
	store, err := NewTsRouteStore(&model.TsStoreConfig{
		Host: "127.0.0.1",
		User: "idagoras",
		Password: "314159",
		Dbname: "shiftroute",
		Port: 5432,
		SSLMode: "disable",
		TimeZone: "Asia/Shanghai",
	})
	require.NoError(t,err)
	edges, err := store.ListPointsCanDirectedAchieved("BV11009136")
	require.NoError(t,err)
	jsonByte, err:= json.Marshal(edges)
	require.NoError(t,err)
	fmt.Println(string(jsonByte))

	points, err := store.ListPointsCanAchieved("BV11009136",1)
	require.NoError(t,err)
	jsonByte, err = json.Marshal(points)
	require.NoError(t,err)
	fmt.Println(string(jsonByte))
}