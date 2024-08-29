package store

import (
	"fmt"
	"testing"

	"github.com/idagoras/TsRoute/model"
	"github.com/stretchr/testify/require"
)

func TestDataLoader(t *testing.T){
	dataLoader := NewGuangzhouTsDataLoader()
	data, err := dataLoader.Load("../dataset/guangzhou.csv",0,100)
	require.NoError(t,err)
	for _, record := range data{
		stopPair, ok := record.(model.StopPair)
		require.Equal(t,ok,true)
		fmt.Println(stopPair.Origin.Id)
		fmt.Println(stopPair.Origin.Name)
		fmt.Println(stopPair.Origin.Location)
		fmt.Println(stopPair.Destination.Id)
		fmt.Println(stopPair.Destination.Name)
		fmt.Println(stopPair.Destination.Location)
	}
}