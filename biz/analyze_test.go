package biz

import (
	"fmt"
	"testing"

	"github.com/idagoras/TsRoute/model"
	"github.com/stretchr/testify/require"
)

func TestAnalyze(t *testing.T){
	lbs := NewGaoDeLbsServer(5)
	res, err := lbs.Serve(GaoDeLbsTransportationRouteQuery,model.GaoDeLbsTransportationRouteQueryLbsRequest{
		Key: "3bf4992397329e61be6b62a1399a2145",
		Origin: "113.282936,23.137556",
		Destination: "113.35169,23.13555",
		City1: "020",
		City2: "020",
	})
	require.NoError(t,err)
	rep, ok := res.(model.GaoDeLbsTransportationRouteQueryResponse)
	require.Equal(t,ok,true)
	route := &rep.Route

	s := routeLCSS(route,route)
	require.Equal(t,s,1.0)

	analyst := NewAnalyst(lbs)
	origin1 := model.BusStop{Name:"广东电视台",Id: "BV11009136",Location: "113.282936,23.137556"}
	dest1 := model.BusStop{Name:"华师大南门",Id: "BV10019585",Location: "113.35169,23.13555"}
	stopPair := &model.StopPair{
			Origin: origin1,
			Destination: dest1,
		}
	result, err := analyst.Evaluate(
		map[string]string{
			"epsilon0":"0.01",
			"epsilon1":"0.03",
			"k":"3",
		},
		[]*model.StopPair{stopPair},
		[]*model.Route{route},
	)

	require.NoError(t,err)
	require.Equal(t,result.Num,1)
	result.Similarity.Vals.Range(func(key, value any) bool {
		fmt.Println(key,value)
		return true
	})
}