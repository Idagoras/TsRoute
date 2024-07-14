package biz

import (
	"math"

	"github.com/peterli110/discreteprobability"
)

type PlanerExponentialMechanism struct{
	rng *discreteprobability.Generator
}

func NewPlanerExponentialMechanism(keys []string,values []float64, epsilon float32) (*PlanerExponentialMechanism, error){
	var sum float64
	for i := range values{
		values[i] = math.Exp(-values[i]*float64(epsilon))
		sum += float64(values[i])
	}
	for i:= range values{
		values[i] = values[i]/sum
	}
	stringRNG, err := discreteprobability.New(keys,values)
	if err != nil{
		return nil, err
	}
	return &PlanerExponentialMechanism{
		rng: stringRNG,
	},nil
}

func(m *PlanerExponentialMechanism)RandomString() string{
	return m.rng.RandomString()
}