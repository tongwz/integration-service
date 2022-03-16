package jz

import (
	"errors"
	"fmt"
	"integration.service/models/mongodb"
	"integration.service/models/mysql"
	"integration.service/pkg/db"
	"integration.service/pkg/logging"
	"integration.service/pkg/setting"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ExceptCompanyId 交流记录中有我们现有公司的聊天记录文件需要过滤不进行删除
// FileBasePath 文件位置的根目录
// TableNameMap 实际年月对应的 表格名称 生产与测试环境不同
type DelImageFiles struct {
	WechatSessionHistoryDao mysql.WechatSessionHistoryDao
	AttachmentDao           mongodb.AttachmentDao
	ExceptCompanyId         int64
	FileBasePath            string
	TableNameMap            map[string]string
}

func init() {

}

/**
 * @note: 操作删除多余文件
 * @auth: tongwz
 * @date  2022年2月18日10:04:41
**/
func (dif *DelImageFiles) DoIt(start string, end string) {
	months, err := dif.MonthSlice(start, end)
	if err != nil {
		fmt.Printf("月份日期发生错误：%s", err.Error())
		return
	}
	// fmt.Printf("日期的切片为：%v+, len:%d \n", months, len(months))
	pageSize := setting.PageSize
	delMsgType := map[string]string{
		"emotion":            "emotion",
		"image":              "image",
		"video":              "video",
		"file":               "file",
		"voice":              "voice",
		"meeting_voice_call": "meeting_voice_call",
	}
	// 先进行数据库查询
	for _, month := range months {
		// TODO:由于他们的表都在2021后缀的表里面因此单独所以像20222对应的表也是20212
		realMonth, ok := dif.TableNameMap[month]
		if !ok {
			fmt.Println("月份没找到：", month)
			continue
		}
		days := dif.MonthForEveryday(month)
		fmt.Printf("所有日期为：%v, len:%d \n", days, len(days))
		// TODO: 由于创建日期是表格key所以我们需要将每天的数据进行查询 不然后面将会特别慢
		for _, msgCreateAt := range days {
			/**
			TODO:由于文件服务器存储的时候对于相同的文件并没有多次存储，因此需要过滤掉jz 和 盈亚 都有的文件，不对它进行删除
			*/
			// redis 添加
			var wait = &sync.WaitGroup{}
			dif.addRedis(wait, realMonth, msgCreateAt, pageSize, delMsgType)

			fmt.Println(msgCreateAt + "已经redis准备好了数据-------------------------")
			// 删除文件
			var waitDel = &sync.WaitGroup{}
			dif.delFiles(waitDel, realMonth, msgCreateAt, pageSize, delMsgType)
			fmt.Println(msgCreateAt + "已经处理好了数据------------------------------")

			// 删除redis key
			dif.delRedis(msgCreateAt)
			fmt.Println(msgCreateAt + "已经删除了redis数据---------------------------")
		}
	}

}

/**
 * @note: 删除文件逻辑
 * @auth: tongwz
 * @date 2022年2月24日13:23:09
**/
func (dif *DelImageFiles) delFiles(wait *sync.WaitGroup, realMonth string, msgCreateAt string, pageSize int, delMsgType map[string]string) {
	page := 1
	ch := make(chan int, 1000)
	for {
		// TODO:由于没有用到索引，因此我们需要跑每天的数据
		historyFiles, err := dif.WechatSessionHistoryDao.SessionList(db.MysqlMb4, realMonth, msgCreateAt, pageSize, page)
		if err != nil {
			logging.Error("查询mysql数据报错", err.Error())
			return
		}
		if len(historyFiles) == 0 {
			break
		}
		wait.Add(len(historyFiles))
		for ii, WechatSessionHistory := range historyFiles {
			ch <- ii
			// TODO:使用多线程进行删除文件
			go func(WechatSessionHistory mysql.WechatSessionHistory) {
				// 如果不是这几个数据类型我们跳过，因为只有这几个是有文件的
				if _, ok := delMsgType[WechatSessionHistory.MsgType]; !ok {
					defer wait.Done()
					<-ch
					return
				}
				// 我们企业的需要过滤
				if WechatSessionHistory.CompanyId == dif.ExceptCompanyId {
					defer wait.Done()
					<-ch
					return
				}
				//fmt.Printf("获取到的文件数据为：%v \n", WechatSessionHistory)
				// 通过链接字符串 获取到_id
				fileUrl := WechatSessionHistory.Content
				idIndex := strings.Index(fileUrl, "id=")
				_id := fileUrl[idIndex+3:]
				// 如果redis中存在那么我们跳过
				if db.Redis.Exists(msgCreateAt + "_" + _id).Val() == int64(1) {
					defer wait.Done()
					<-ch
					logging.Info("redis-key = " + msgCreateAt + "_" + _id + "保护企业有同样文件------------------------")
					return
				}
				// fmt.Printf("获取到的_id=%s \n", _id)
				// 获取到文件内容
				attachment, _ := dif.AttachmentDao.GetFileById(db.MongodbClient, _id)
				if attachment.ID == "" {
					fmt.Println("查询不到这条文件内容！！！！！_id=", _id)
					defer wait.Done()
					<-ch
					return
				}
				// 如果名称需要替换 我们需要把地址改一下
				removeFile := dif.FileBasePath + attachment.RelativePath + "/" + attachment.SaveFileName

				// TODO：删除文件操作
				err := os.Remove(removeFile)
				if err != nil {
					// logging.Info("文件删除失败，失败原因：", err.Error())
					defer wait.Done()
					<-ch
					return
				}
				<-ch
				wait.Done()
				// fmt.Printf("文件删除成功：%+v \n", removeFile)
			}(WechatSessionHistory)
		}
		if len(historyFiles) < pageSize {
			break
		}
		page++
	}
	wait.Wait()
}

/**
 * @note: 给每日的我们公司企业加上缓存
 * @auth: tongwz
 * @date  2022年2月24日13:20:41
**/
func (dif *DelImageFiles) addRedis(wait *sync.WaitGroup, realMonth string, msgCreateAt string, pageSize int, delMsgType map[string]string) {
	exceptCompanyPage := 1
	ch := make(chan int, 1000)
	for {
		firstHistoryFiles, err := dif.WechatSessionHistoryDao.SessionList(db.MysqlMb4, realMonth, msgCreateAt, pageSize, exceptCompanyPage)
		if err != nil {
			logging.Error("查询mysql数据报错", err.Error())
			return
		}
		if len(firstHistoryFiles) == 0 {
			break
		}
		wait.Add(len(firstHistoryFiles))
		// 将过滤企业留下，之后把文件留下放入redis
		for ii, filesInfo := range firstHistoryFiles {
			ch <- ii
			// 并发执行但是 必须这边全部执行完才能执行下面的for循环
			go func(filesInfo mysql.WechatSessionHistory) {
				if _, ok := delMsgType[filesInfo.MsgType]; !ok {
					defer wait.Done()
					<-ch
					return
				}
				// 不是我们的企业需要过滤
				if filesInfo.CompanyId != dif.ExceptCompanyId {
					defer wait.Done()
					<-ch
					return
				}
				// 将留下的数据 id存下来 放入redis
				fileUrlDiy := filesInfo.Content
				idIndex := strings.Index(fileUrlDiy, "id=")
				_id := fileUrlDiy[idIndex+3:]
				db.Redis.Set(msgCreateAt+"_"+_id, 1, 120*time.Minute)
				defer wait.Done()
				<-ch
				return
			}(filesInfo)
		}
		// 最后一页也需要跳过
		if len(firstHistoryFiles) < pageSize {
			break
		}
		exceptCompanyPage++
	}
	wait.Wait()
}

// 某个月对应的所有日期
func (dif *DelImageFiles) MonthForEveryday(month string) []string {
	monthTime, _ := time.Parse("20061", month)
	nextMonth := monthTime.AddDate(0, 1, 0)
	days := make([]string, 0)
	days = append(days, monthTime.Format("2006-01-02"))
	for {
		if tempTime := monthTime.AddDate(0, 0, 1); tempTime.Unix() < nextMonth.Unix() {
			days = append(days, tempTime.Format("2006-01-02"))
			monthTime = tempTime
		} else {
			break
		}
	}
	return days
}

/**
 * @note: 删除redis数据
 * @auth: tongwz
 * @date  2022年2月24日14:12:43
**/
func (dif *DelImageFiles) delRedis(msgCreateAt string) {
	redisKeys, err := db.Redis.Do("keys", msgCreateAt+"*").Result()
	if err != nil {
		logging.Error("查询redis key失败", err.Error())
	}
	// TODO：通过反射查看数据的具体类型
	if reflect.TypeOf(redisKeys).Kind() == reflect.Slice {
		redisKeysVal := reflect.ValueOf(redisKeys)
		if redisKeysVal.Len() == 0 {
			return
		}
		for i := 0; i < redisKeysVal.Len(); i++ {
			db.Redis.Del(redisKeysVal.Index(i).Interface().(string))
			fmt.Printf("删除掉的redisKey = %s \n", redisKeysVal.Index(i).Interface().(string))
		}
	}
}

// 月份对应的表格生成
func (dif *DelImageFiles) MonthSlice(start string, end string) (slice []string, err error) {
	var months = make([]string, 1)
	// 获取月份范围内的表后缀
	intStart, err := strconv.Atoi(start)
	intEnd, err2 := strconv.Atoi(end)
	if err != nil || err2 != nil {
		return nil, err
	}
	if intStart > intEnd {
		return nil, errors.New("开始时间不能大于结束时间")
	}

	timeStart, _ := time.Parse("200601", start)
	timeEnd, _ := time.Parse("200601", end)

	months[0] = timeStart.Format("20061")
	var i = 0
	for {
		timeTemp := timeStart.AddDate(0, 1, 0)
		// fmt.Println(timeTemp.Format("200601"), timeTemp.Unix(), timeEnd.Unix())
		if timeTemp.Unix() > timeEnd.Unix() {
			break
		}
		i++
		months = append(months, timeTemp.Format("20061"))
		timeStart = timeTemp
	}
	return months, nil
}
