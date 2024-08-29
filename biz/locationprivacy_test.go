package biz

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocationPrivacy(t *testing.T) {
	lpCalculator := &LocationPrivacyCalculaterImpl{}
	var dists []*DiscreteDistribution
	xs := []any{"A","B","C","D","E"}
	prs := []float64{0.2,0.2,0.2,0.2,0.2}
	for i := 0 ; i < len(xs); i ++{
		dist, err := NewDiscretedDistribution(xs,prs)
		require.NoError(t,err)
		dists = append(dists, dist)
	}
	priorityDist , err := NewDiscretedDistribution(xs,prs)
	require.NoError(t,err)
	lp := lpCalculator.Calculate(dists,priorityDist)
	fmt.Println(lp)
}