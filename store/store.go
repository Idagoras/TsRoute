package store

import (
	"strconv"

	"github.com/idagoras/TsRoute/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TsRouteStoreImpl struct{
	db *gorm.DB
}

func NewTsRouteStore(config *model.TsStoreConfig) (model.TsStore, error){
	var dsn string
	dsn += " host=" + config.Host
	dsn += " user=" + config.User
	dsn += " password=" + config.Password
	dsn += " dbname=" + config.Dbname
	dsn += " port=" + strconv.Itoa(config.Port)
	dsn += " sslmode="+ config.SSLMode
	dsn += " TimeZone" + config.TimeZone
	db, err := gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil{
		return nil, err
	}
	return &TsRouteStoreImpl{
		db: db,
	},nil
}

func(s *TsRouteStoreImpl)ListPointsBetween(tfLon,tfLat,brLon,brLat float32) ([]*model.TsStop,error){
	var stops []*model.TsStop
	err := s.db.Model(&model.TsStop{}).Where("longitude >= ? and longitude <= ? and latitude <= ? and latitude >= ?",tfLon,brLon,tfLat,brLat).Find(&stops).Error
	if err != nil{
		return nil, err
	}
	return stops, nil
}

func(s *TsRouteStoreImpl)GetOptimalRoute(origin string,destination string) (*model.TsOptimalRoute, error){
	var route model.TsOptimalRoute
	err := s.db.Model(&model.TsOptimalRoute{}).Where("origin = ? and destination = ?",origin,destination).First(&route).Error
	if err != nil{
		return nil, err
	}
	return &route, nil
}

func(s *TsRouteStoreImpl)InsertOptimalRoute(route *model.TsOptimalRoute) error{
	err := s.db.Model(&model.TsOptimalRoute{}).Create(route).Error
	if err != nil{
		return err
	}
	return nil
}

func(s *TsRouteStoreImpl)UpdateOptimalRoute(route *model.TsOptimalRoute) error{
	err := s.db.Model(&model.TsOptimalRoute{}).Updates(route).Error
	if err != nil{
		return err
	}
	return nil
}