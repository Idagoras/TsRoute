package biz

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPerturb(t *testing.T){
	pb, err := NewPlanerExponentialMechanism(
		[]string{"BRT师大暨大S1","东方一路","东方一路口","东方三路"},
		[]float64{0,0.6,0.7,1.2},
		0.01,
	)
	require.NoError(t,err)
	for i := 0; i < 100; i++{
		fmt.Println(pb.RandomString())
	}
}