package store

import (
	"encoding/csv"
	"os"

	"github.com/idagoras/TsRoute/model"
)

type GuangzhouTsDataLoader struct{

}

func NewGuangzhouTsDataLoader() model.DataLoader{
	return &GuangzhouTsDataLoader{}
}

func(d *GuangzhouTsDataLoader) Load(path string,begin,num int)([]interface{},error){
	file ,err := os.Open(path)
	if err != nil{
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil{
		return nil, err
	}

	l := len(records)
	var result []interface{}
	if begin * num >= l{
		return result,nil
	}

	var end int
	if (begin+1) * num > l{
		end = l
	}else{
		end = (begin+1) * num
	}
	for i := begin * num ; i < end ; i ++{
		record := records[i]
		result = append(result, model.StopPair{
			Origin: model.BusStop{
				Id: record[0],
				Location: record[2] + "," + record[3],
				Name: record[6],
			},
			Destination: model.BusStop{
				Id: record[1],
				Location: record[4] + "," + record[5],
				Name: record[7],
			},
		})
	}
	return result, nil
}