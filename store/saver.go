package store

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"

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