package store

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/idagoras/TsRoute/model"
)

type GuangzhouTsDataSaver struct{

}

func NewGuangzhouTsDataSaver() model.DataSaver{
	return &GuangzhouTsDataSaver{}
}

func(s *GuangzhouTsDataSaver)Save(path string,data any) error{
	file, err := os.Create(path)
	if err != nil{
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	result, ok := data.(*model.AnalyzeResult)
	if !ok{
		return errors.New("type is not AnalyzeResult")
	}

	err = writer.Write([]string{
		"epsilon_small","epsilon_large","number of queries","topLeftPoint","bottomRightPoint","partition",
	})
	if err != nil{
		return err
	}

	err = writer.Write([]string{
		result.MetaData["epsilon0"],
		result.MetaData["epsilon1"],
		result.MetaData["k"],
		result.MetaData["tf"],
		result.MetaData["br"],
		result.MetaData["pt"],
	})
	if err != nil{
		return err
	}

	err = writer.Write([]string{
		"origin",
		"deatination",
		"originId",
		"destinationId",
		"originName",
		"destinationName",
		"similarity",
		"RDR",
		"RCR",
		"RTR",
	})
	if err != nil{
		return err
	}


	for _, pair := range result.Pairs{
		key := fmt.Sprintf("%f,%f,%f,%f",pair.Origin.Lon(),pair.Origin.Lat(),pair.Destination.Lon(),pair.Destination.Lat())
		fmt.Println(key)
		lcss, _ := result.Similarity.Vals.Load(key)
		similarity, _ := lcss.(float64)
		rdrp , _ := result.RDR.Vals.Load(key)
		rdr, _ := rdrp.(float64)
		rtrp , _ := result.RTR.Vals.Load(key)
		rtr, _ := rtrp.(float64)
		rcrp , _ := result.RCR.Vals.Load(key)
		rcr, _ := rcrp.(float64)
		record := []string{
			fmt.Sprintf("%f:%f",pair.Origin.Lon(),pair.Origin.Lat()),
			fmt.Sprintf("%f:%f",pair.Destination.Lon(),pair.Destination.Lat()),
			pair.Origin.Id,
			pair.Destination.Id,
			pair.Origin.Name,
			pair.Destination.Name,
			fmt.Sprintf("%f",similarity),
			fmt.Sprintf("%f",rdr),
			fmt.Sprintf("%f",rcr),
			fmt.Sprintf("%f",rtr),
		}
		err := writer.Write(record)
		if err != nil{
			return err
		}
	}

	return nil
}


type GridDataSaver struct{

}

func NewGridDataSaver() *GridDataSaver{
	return &GridDataSaver{}
}

func(s *GridDataSaver)Save(path string,data any) error{
	file, err := os.Create(path)
	if err != nil{
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	result, ok := data.(*model.GridAnalyzeResult)
	if !ok{
		return errors.New("type is not GridAnalyzeResult")
	}

	err = writer.Write([]string{
		"areaTopLeft","areaBottomRight","xGridNums","yGridNums",
	})

	if err != nil{
		return err
	}

	err = writer.Write([]string{
		result.MetaData["areaTopLeft"],
		result.MetaData["areaBottomRight"],
		result.MetaData["xGridNums"],
		result.MetaData["yGridNums"],
	})

	if err != nil{
		return err
	}

	result.Points.Vals.Range(func(key, value any) bool {
		strKey, _ := key.(string)
		pointVal, _:= value.(int)
		lineInterfaceVal, ok := result.Lines.Vals.Load(strKey)
		if !ok{
			return false
		} 
		lineVal, _ := lineInterfaceVal.(int)
		edgeInterfaceVal, ok := result.Edges.Vals.Load(strKey)
		if !ok{
			return false
		} 
		edgeVal, _ := edgeInterfaceVal.(int)
		err := writer.Write([]string{
			strKey,
			strconv.Itoa(pointVal),
			strconv.Itoa(lineVal),
			strconv.Itoa(edgeVal),
		})
		return err == nil
	})
	return nil
}

type KeyObservationSaver struct{

}

func NewKeyObservationSaver() model.DataSaver{
	return &KeyObservationSaver{}
}

func(s *KeyObservationSaver)Save(path string,data any) error{
	file, err := os.Create(path)
	if err != nil{
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	result, ok := data.(*model.KeyObservationResult)
	if !ok{
		return errors.New("type is not KeyObservationResult")
	}

	err = writer.Write([]string{"number of stops","similarity"})
	if err != nil{
		return err
	}
	for index, sims := range result.Similarity.KV{
		floatArr, _ := sims.([]float64)
		for _, fv := range floatArr{
			err := writer.Write([]string{
				index,
				fmt.Sprintf("%f",fv),
			})		
			if err != nil{
				return err
			}
		}
	}
	return nil
}


type LocationPrivacySaver struct{

}

func NewLocationPrivacySaver() *LocationPrivacySaver{
	return &LocationPrivacySaver{}
}

func(s *LocationPrivacySaver)Save(path string,data any) error{
	floatArr , _ := data.([]float64)
	file, err := os.Create(path)
	if err != nil{
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, value := range floatArr{
		err := writer.Write([]string{fmt.Sprintf("%f",value)})
		if err != nil{
			return err
		}
	}
	return nil
}