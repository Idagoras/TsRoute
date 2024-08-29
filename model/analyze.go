package model

import "sync"

type SetResult struct{
	Vals *sync.Map
	Mean float64
	Variance float64
	StdDev float64
}

type MapResult struct{
	KV map[string]any
}

func NewSetResult() *SetResult{
	return &SetResult{
		Vals: &sync.Map{},
	}
}

func NewMapResult() *MapResult{
	return &MapResult{
		KV:make(map[string]any),
	}
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

type GridAnalyzeResult struct{
	Num int
	Grids []*TsGrid
	Lines *SetResult
	Edges *SetResult
	Points *SetResult
	MetaData map[string]string
	Extra *MapResult
}

type KeyObservationResult struct{
	Similarity *MapResult
	MetaData map[string]string
}

