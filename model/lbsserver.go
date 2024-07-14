package model

type LbsServer interface{
	Serve(serviceType int,request any)(any, error)
}