package biz

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/idagoras/TsRoute/model"
	log "github.com/sirupsen/logrus"
)

type GridPartitionManager struct{
	grids map[string]map[string]*model.TsGrid
	partitionSize map[string]string
	points []model.TsPoint
	areaTopLeftPoint model.TsPoint
	areaBottomRightPoint model.TsPoint
	store model.TsStore
}

func NewGridPartitionManager(areaTopLeftPoint, areaBottomRightPoint model.TsPoint, store model.TsStore)(*GridPartitionManager,error){
	grids := make(map[string]map[string]*model.TsGrid)
	points , err := store.ListPointsBetween(areaTopLeftPoint.Lon(),areaTopLeftPoint.Lat(),areaBottomRightPoint.Lon(),areaBottomRightPoint.Lat())
	if err != nil{
		log.WithFields(log.Fields{"areaTopLeftPoint":areaTopLeftPoint,"areaBottomRightPoint":areaBottomRightPoint}).Error(err)
		return nil,err
	}
	var p []model.TsPoint
	for _, point := range points{
		p = append(p, point)
	}

	log.WithFields(log.Fields{
		"num":len(points),
	}).Info("get points from store")
	fmt.Println(len(points))
	return &GridPartitionManager{
		points: p,
		grids: grids,
		areaTopLeftPoint: areaTopLeftPoint,
		areaBottomRightPoint: areaBottomRightPoint,
		store: store,
		partitionSize: make(map[string]string),
	},nil
}

func(m *GridPartitionManager) AddPartition(xGridsNum, yGridsNum int, key string){
	var grids map[string]*model.TsGrid
	var ok bool
	if _, ok =  m.grids[key]; ok{
		return
	}else{
		grids = make(map[string]*model.TsGrid)
		m.grids[key] = grids
	}
	originX := m.areaTopLeftPoint.Lon()
	originY := m.areaBottomRightPoint.Lat()
	width := (m.areaBottomRightPoint.Lon() - m.areaTopLeftPoint.Lon())/float32(xGridsNum)
	height := (m.areaTopLeftPoint.Lat()- m.areaBottomRightPoint.Lat())/float32(yGridsNum)
	for _, point := range m.points{
		i := int((point.Lon() - originX)/width)
		j := int((point.Lat() - originY)/height)
		gridKey := fmt.Sprintf("%d,%d",i,j)
		if grid, ok:= grids[gridKey]; !ok{
			tfPoint := &model.TsStop{
				Longitude:float64(originX) + float64(i) * float64(width),
				Latitude: float64(originY) + float64(j+1) * float64(height),
			}
			brPoint := &model.TsStop{
				Longitude:float64(originX) + float64(i+1) * float64(width),
				Latitude: float64(originY) + float64(j) * float64(height),
			}
			grid = model.NewTsGrid(tfPoint,brPoint)
			grid.AddPoint(point)
			grids[gridKey] = grid
		}else{
			grid.AddPoint(point)
		}
	}
	size:= fmt.Sprintf("%f,%f",width,height)
	m.partitionSize[key] = size 
	log.WithFields(log.Fields{"key":key,"size":size}).Info("add partition success")
	var gridKeys []string
	for key := range m.grids[key]{
		gridKeys = append(gridKeys, key)
	}
	log.WithFields(log.Fields{"girdKeys":gridKeys}).Info("partiton add gridKeys")

}

func(m *GridPartitionManager) GetGrid(key string,point model.TsPoint) *model.TsGrid{
	size, ok := m.partitionSize[key]
	if !ok{
		log.WithFields(log.Fields{"key":key}).Error("key not exist in partitionSize Map")
		return nil
	}
	sizeArr := strings.Split(size,",")
	width, _ := strconv.ParseFloat(sizeArr[0],32)
	height, _ := strconv.ParseFloat(sizeArr[1],32)
	originX := m.areaTopLeftPoint.Lon()
	originY := m.areaBottomRightPoint.Lat()
	i := int((point.Lon() - originX)/float32(width))
	j := int((point.Lat() - originY)/float32(height))
	gridKey := fmt.Sprintf("%d,%d",i,j)

	grid, ok :=m.grids[key][gridKey]
	if !ok{
		log.WithFields(log.Fields{"key":key,"partition":size,"gridkey":gridKey,"point":point}).Error("key not exist in grids Map")
		return nil
	}

	return grid
}

func(m *GridPartitionManager)GetAllGrids(key string) []*model.TsGrid{
	var result []*model.TsGrid
	for _, grid := range m.grids[key]{
		result = append(result, grid)
	}
	return result
}

func(m *GridPartitionManager)GetEdgesAndLinesInGrids(grid *model.TsGrid)([]string,[]*model.TsEdge,error){
	points, err := grid.GetAllPoints()
	if err != nil{
		return nil,nil,err
	}
	var lines []string
	var edges []*model.TsEdge
	edgesMap := make(map[string]*model.TsEdge)
	linesMap := make(map[string]struct{})

	for _, point := range points{
		stop, _ := point.(*model.TsStop)
		edgesModel , err := m.store.ListPointsCanDirectedAchieved(stop.Id)
		//fmt.Println(edgesModel)
		if err != nil{
			return nil,nil,err
		}
		for _, edgeModel := range edgesModel{
			idArr := []string{edgeModel.StopId,edgeModel.ToStopId}
			linesMap[edgeModel.LineId] = struct{}{}
			sort.Strings(idArr)
			edgesMap[idArr[0]+idArr[1]] = edgeModel
		}
	}

	for line := range linesMap{
		lines = append(lines, line)
	}

	for _, edge := range edgesMap{
		edges = append(edges, edge)
	}

	return lines, edges, nil
}