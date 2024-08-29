package store

import (
	"fmt"
	"sync"
	"testing"

	"github.com/idagoras/TsRoute/model"
	"github.com/stretchr/testify/require"
)

func TestDataSaver(t *testing.T){
	origin1 := model.BusStop{Name:"广东电视台",Id: "BV11009136",Location: "113.282936,23.137556"}
	dest1 := model.BusStop{Name:"华师大南门",Id: "BV10019585",Location: "113.35169,23.13555"}
	origin2 := model.BusStop{Name:"迎宾馆",Id: "BV11091393",Location: "113.26201,23.12899"}
	dest2 := model.BusStop{Name:"海珠广场总站(侨光西)",Id: "BV11206247",Location: "113.26474,23.11407"}
	saver := NewGuangzhouTsDataSaver()
	key1 := fmt.Sprintf("%f,%f,%f,%f",origin1.Lon(),origin1.Lat(),dest1.Lon(),dest1.Lat())
	key2 := fmt.Sprintf("%f,%f,%f,%f",origin2.Lon(),origin2.Lat(),dest2.Lon(),dest2.Lat())
	fmt.Println(key1)
	fmt.Println(key2)
	err := saver.Save("../result/test.csv",&model.AnalyzeResult{
		Num: 2,
		Pairs: []*model.StopPair{
			{
				Origin: origin1,
				Destination: dest1,
			},
			{
				Origin: origin2,
				Destination: dest2,
			},
		},
		Similarity: &model.SetResult{
			Vals: func()*sync.Map{
				mp := sync.Map{}
				mp.Store(key1,0.95)
				mp.Store(key2,0.91)
				return &mp
			}(),
			Mean: 0.93,
			Variance: 0.93,
			StdDev: 0.93,
		},
		RDR: &model.SetResult{
			Vals: func()*sync.Map{
				mp := sync.Map{}
				mp.Store(key1,0.95)
				mp.Store(key2,0.91)
				return &mp
			}(),
			Mean: 0.93,
			Variance: 0.93,
			StdDev: 0.93,
		},
		RTR: &model.SetResult{
			Vals: func()*sync.Map{
				mp := sync.Map{}
				mp.Store(key1,0.95)
				mp.Store(key2,0.91)
				return &mp
			}(),
			Mean: 0.93,
			Variance: 0.93,
			StdDev: 0.93,
		},
		RCR: &model.SetResult{
			Vals: func()*sync.Map{
				mp := sync.Map{}
				mp.Store(key1,0.95)
				mp.Store(key2,0.91)
				return &mp
			}(),
			Mean: 0.93,
			Variance: 0.93,
			StdDev: 0.93,
		},
		MetaData: map[string]string{
			"epsilon0":"0.01",
			"epsilon1":"0.05",
			"k":"5",
			"tf":"113.2445,23.2077",
			"br":"113.4445,23.0077",
			"pt":"10*10",
		},

	})
	require.NoError(t,err)
}