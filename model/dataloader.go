package model

type StopPair struct{
	Origin BusStop
	Destination BusStop
}

type DataLoader interface{
	Load(path string,begin,num int)([]interface{},error)
}