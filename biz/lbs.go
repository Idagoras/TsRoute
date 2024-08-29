package biz

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/idagoras/TsRoute/model"
)

var(
	ErrInvalidServiceType = errors.New("invalid service type")
	ErrInvalidRequest = errors.New("invalid request type")
)

const(
	GaoDeLbsTransportationRouteQuery = iota
)

type GaoDeLbsServer struct{
	client *http.Client
	urls map[int]string
}

func NewGaoDeLbsServer(timeout int) model.LbsServer{
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	return &GaoDeLbsServer{
		client: client,
		urls: map[int]string{
			GaoDeLbsTransportationRouteQuery:"https://restapi.amap.com/v5/direction/transit/integrated",
		},
	}
}

func joinUrl(baseUrl string,form map[string]string) string{
	baseUrl += "?"
	for key, value := range form{
		baseUrl += key+"="+value+"&"
	}
	return baseUrl[:len(baseUrl)-1]
}

func(s *GaoDeLbsServer)Serve(serviceType int,request any)(response any, err error){
	var url string
	url, ok := s.urls[serviceType]
	if !ok{
		return nil, ErrInvalidServiceType
	}

	req, ok := request.(model.GaoDeLbsTransportationRouteQueryLbsRequest)
	if !ok{
		return nil,ErrInvalidRequest
	}

	url =  joinUrl(url,map[string]string{
		"key" : req.Key,
		"origin":req.Origin,
		"destination":req.Destination,
		"city1":req.City1,
		"city2":req.City2,
		"AlternativeRoute":"1",
		"show_fields":"cost",
	})

	rep, err := s.client.Get(url)
	if err != nil{
		return nil, err
	}
	defer rep.Body.Close()
	body, err := io.ReadAll(rep.Body)

	var responseObj model.GaoDeLbsTransportationRouteQueryResponse
	err = json.Unmarshal(body,&responseObj)
	if err != nil{
		return nil, err
	}
	time.Sleep(1 * time.Second)
	return responseObj, nil
}
