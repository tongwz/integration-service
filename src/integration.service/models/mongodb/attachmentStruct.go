package mongodb

import (
	"integration.service/pkg/setting"
)

type Attachment struct {
	ID           string `json:"_id" bson:"_id"`
	FileName     string `json:"file_name" bson:"file_name"`
	SaveFileName string `json:"save_file_name" bson:"save_file_name"`
	RelativePath string `json:"relative_path" bson:"relative_path"`
	FullPath     string `json:"full_path" bson:"full_path"`
}

func (Attachment) TableName() string {
	config, _ := setting.Cfg.GetSection("mongodb")
	return config.Key("prefix").MustString("sscf_") + "attachment"
}
