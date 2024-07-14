package model

import "sync"

type SetResult struct{
	Vals sync.Map
	Mean float64
	Variance float64
	StdDev float64
}

type AnalyzeResult struct{
	Num int
	Pairs []*StopPair
	Similarity *SetResult
	RDR *SetResult
	RTR *SetResult
	RCR *SetResult
	MetaData map[string]string
}
