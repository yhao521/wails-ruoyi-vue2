package system

import (
	"context"
	"fmt"
	"mySparkler/backend/api/baseAPI"
	"mySparkler/backend/model/system"
	"mySparkler/backend/model/tools"
	"mySparkler/pkg/utils"
	"mySparkler/pkg/utils/R"
	"mySparkler/pkg/utils/jwt"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
)

type NoticeAPI struct {
	ctx context.Context
	baseAPI.Base
}

var noticeAPI *NoticeAPI
var onceNoticeAPI sync.Once

// NewApp creates a new App application struct
func NewNoticeAPI() *NoticeAPI {
	if noticeAPI == nil {
		onceNoticeAPI.Do(func() {
			noticeAPI = &NoticeAPI{}
		})
	}
	return noticeAPI
}

func (a *NoticeAPI) ListNotice(params map[string]interface{}) tools.TableDataInfo {

	var param = tools.SearchTableDataParam{
		PageNum:  1,
		PageSize: 10,
		Other: system.SysNotice{
			NoticeTitle: "",
			CreateBy:    "",
			NoticeType:  "",
		},
		Params: tools.Params{
			BeginTime: "",
			EndTime:   "",
		},
	}
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}
	var result = system.SelectSysNoticeList(param, true)
	return result
}

func (a *NoticeAPI) GetNotice(noticeId int) R.Result {
	result := system.FindNoticeInfoById(noticeId)
	return R.ReturnSuccess(result)
}

func (a *NoticeAPI) SaveNotice(params map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	var noticeParam = system.SysNotice{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &noticeParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := system.FindUserById(userId)
	noticeParam.CreateBy = user.UserName
	noticeParam.CreateTime = time.Now()
	result := system.SaveNotice(noticeParam)
	return result
}

func (a *NoticeAPI) UploadNotice(params map[string]interface{}) R.Result {
	userId := jwt.CacheGetUserId()
	var noticeParam = system.SysNotice{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &noticeParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := system.FindUserById(userId)
	noticeParam.UpdateBy = user.UserName
	noticeParam.UpdateTime = time.Now()
	result := system.UploadNotice(noticeParam)
	return result
}

func (a *NoticeAPI) DeleteNotice(noticeIds string) R.Result {
	// var noticeIds = context.Param("noticeIds")
	result := system.DeleteNotice(noticeIds)
	return result
}
