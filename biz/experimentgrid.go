package biz

import (
	"fmt"
	"strconv"

	"github.com/idagoras/TsRoute/model"
	log "github.com/sirupsen/logrus"
)

type GridExperiment struct{
	analyst *GridAnalyst
	gridManager *GridPartitionManager 
	dataSaver model.DataSaver 
	savePath string
}

func NewGridExperiment(analyst *GridAnalyst, gridManager *GridPartitionManager,saver model.DataSaver,savePath string) *GridExperiment{
	return &GridExperiment{
		analyst: analyst,
		gridManager: gridManager,
		dataSaver: saver,
		savePath: savePath,
	}
}

func(e *GridExperiment)Once(partition string, xGridsNum int, yGridsNum int) error{
	e.gridManager.AddPartition(xGridsNum,yGridsNum,partition)
	grids := e.gridManager.GetAllGrids(partition)
	lineMap := make(map[string][]string)
	edgeMap := make(map[string][]*model.TsEdge)
	totalLine := make(map[string]struct{})
	for _, grid := range grids{
		tf := fmt.Sprintf("%f,%f",grid.TopLeftPoint.Lon(),grid.TopLeftPoint.Lat())
		br := fmt.Sprintf("%f,%f",grid.BottomRightPoint.Lon(),grid.BottomRightPoint.Lat())
		key := tf + ":" + br
		lines, edges, err := e.gridManager.GetEdgesAndLinesInGrids(grid)
		fmt.Println(len(lines),len(edges))
		if err != nil{
			log.Errorf("get edges and lines in grid failed :%s",err)
			return err
		}
		for _, line := range lines{
			totalLine[line] = struct{}{}
		}
		lineMap[key] = lines
		edgeMap[key] = edges
	}
	fmt.Println(len(totalLine))
	metadata := map[string]string{
		"areaTopLeft":fmt.Sprintf("%f,%f",e.gridManager.areaTopLeftPoint.Lon(),e.gridManager.areaTopLeftPoint.Lat()),
		"areaBottomRight":fmt.Sprintf("%f,%f",e.gridManager.areaBottomRightPoint.Lon(),e.gridManager.areaBottomRightPoint.Lat()),
		"xGridNums":strconv.Itoa(xGridsNum),
		"yGridNums":strconv.Itoa(yGridsNum),
	}
	result, err := e.analyst.Evaluate(metadata,grids,lineMap,edgeMap)
	if err != nil{
		log.Errorf("analyse resuly failed:%s",err)
		return err
	}
	savePath := e.savePath + "/" +  "grid" + GetNowTimeStampString() + ".csv"
	e.dataSaver.Save(savePath,result)	
	return nil
}

func(e *GridExperiment)MultipleTimes(k int, f func(int)(string,int,int))error{
	for i := 0; i < k; i ++{
		err := e.Once(f(i))
		if err != nil{
			return err
		}
	}
	return nil
}