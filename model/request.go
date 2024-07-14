package model

type GaoDeLbsTransportationRouteQueryLbsRequest struct{
	Key string 
	Origin string 
	Destination string 
	OriginPoi string 
	DestinationPoi string
	Ad1 string
	Ad2 string
	City1 string
	City2 string
	Strategy string
	AlternativeRoute string
	Multiexport int
	Nightflag int
	Date string
	Time string
	ShowFields []string
	Sig string
	Output string
	Callback string
}