package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"integration.service/pkg/logging"
	"integration.service/pkg/setting"
)

type AttachmentDao struct {

}

/**
 * @note: 通过id查询文件位置
 * @auth: tongwz
 * @date  2022年2月18日19:21:04
**/
func (ad *AttachmentDao) GetFileById(tx *mongo.Client, id string) (Attachment, error) {
	var attachment = Attachment{}
	config, _ := setting.Cfg.GetSection("mongodb")
	objId, err := primitive.ObjectIDFromHex(id)
	//fmt.Printf("获取到的_id=%+v \n", objId)
	if err != nil {
		logging.Debug("通过id生成objId失败")
		return attachment, err
	}

	// 通过库名 和 collection 获取连接句柄
	collection := tx.Database(config.Key("name").String()).Collection(Attachment{}.TableName())
	filter := bson.D{
		{"_id", objId},
	}

	err = collection.FindOne(context.TODO(), filter).Decode(&attachment)
	if err != nil {
		logging.Debug("通过_id查询文件数据失败~_id = ", id)
		return attachment, err
	}
	return attachment, nil
}