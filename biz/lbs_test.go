package biz

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/idagoras/TsRoute/model"
	"github.com/stretchr/testify/require"
)

func TestLbs(t *testing.T){
	lbs := NewGaoDeLbsServer(5)
	res, err := lbs.Serve(GaoDeLbsTransportationRouteQuery,model.GaoDeLbsTransportationRouteQueryLbsRequest{
		Key: "3bf4992397329e61be6b62a1399a2145",
		Origin: "113.26201,23.12899",
		Destination: "113.26474,23.11407",
		City1: "020",
		City2: "020",
	})
	require.NoError(t,err)
	rep, ok := res.(model.GaoDeLbsTransportationRouteQueryResponse)
	require.Equal(t,ok,true)
	jsonBytes, err := json.Marshal(rep)
	require.NoError(t,err)
	fmt.Println(string(jsonBytes))
}