package biz

import (

)


type DiscreteDistribution struct{
	dist map[any]float64
	xs []interface{}
}

func(d *DiscreteDistribution)Pr(x any)float64{
	return d.dist[x]
}

func NewDiscretedDistribution(xs []any,weights []float64)(*DiscreteDistribution,error){
	var sum = 0.0
	for _, weight := range weights{
		sum += weight
	}
	dist := make(map[any]float64)

	for i := range weights{
		dist[xs[i]] = weights[i]
	}
	return &DiscreteDistribution{
		dist: dist,
		xs: xs,
	},nil
}

func(d *DiscreteDistribution)GetStringXs() []string{
	var result []string = make([]string, 0,len(d.xs))
	for _,i := range d.xs{
		strVal, _ := i.(string)
		result = append(result, strVal)
	}
	return result
}


type LocationPrivacyCalculaterImpl struct{

}

func(c *LocationPrivacyCalculaterImpl)Calculate(dists []*DiscreteDistribution, priorityDist *DiscreteDistribution) float64{
	xs := priorityDist.GetStringXs()
	lp := 1.0
	for i := range xs{
		dt := dists[i]
		maxAE := 0.0
		for _, x := range xs{
			maxAE = max(maxAE,dt.Pr(x)*priorityDist.Pr(x))
		}
		lp -= maxAE
	}
	return lp
}