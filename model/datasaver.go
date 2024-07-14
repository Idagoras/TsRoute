package model

type DataSaver interface{
	Save(path string, data any) error
}