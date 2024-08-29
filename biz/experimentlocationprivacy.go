package biz

import (
	"fmt"
	"math"
	"time"

	"github.com/idagoras/TsRoute/store"
)

type LocationPrivacyExperiment struct{
	lpc *LocationPrivacyCalculaterImpl
	gridManager *GridPartitionManager
}

func NewLocationPrivacyExperiment(lpc *LocationPrivacyCalculaterImpl,gridManager *GridPartitionManager)*LocationPrivacyExperiment{
	return &LocationPrivacyExperiment{
		lpc:lpc,
		gridManager: gridManager,
	}
}

func(e *LocationPrivacyExperiment) Once(key string,epsilon float32) error{
	grids := e.gridManager.GetAllGrids(key)
	var lps []float64
	for _, grid := range grids{
		points, err := grid.GetAllPoints()
		if err != nil{
			return err
		}
		dists := make([]*DiscreteDistribution,0,len(points))
		var priorityxs []any
		var priorityprs []float64
		for _, point := range points{
			priorityxs = append(priorityxs, point.Hash())
			priorityprs = append(priorityprs, 1/float64(len(points)))
			var xs []any
			var prs []float64
			distances := grid.GetDistances(point)
			sum := 0.0
			for key, value := range distances{
				val := math.Exp(-value*float64(epsilon))
				xs = append(xs, key)
				prs = append(prs, val)
				sum += val
			}

			for i := range prs{
				prs[i] = prs[i]/sum
			}
			dist,err := NewDiscretedDistribution(xs,prs)
			if err != nil{
				return err
			}
			dists = append(dists, dist)

		}
		priorityDist , err := NewDiscretedDistribution(priorityxs,priorityprs)
		if err != nil{
			return err
		}
		lp := e.lpc.Calculate(dists,priorityDist)
		lps = append(lps, lp)
	}
	saver := store.NewLocationPrivacySaver()
	saver.Save("../result/2020/lp"+time.Now().Format("20060504150201")+fmt.Sprintf("%f",epsilon)+".csv",lps)
	return nil
}