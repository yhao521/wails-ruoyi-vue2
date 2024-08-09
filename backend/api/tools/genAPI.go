package tools

import (
	"context"
	"embed"
	"fmt"
	"log"
	"mySparkler/backend/api/baseAPI"
	"mySparkler/backend/model/tools"
	tools2 "mySparkler/backend/service/tools"
	"mySparkler/pkg/utils"
	"mySparkler/pkg/utils/R"
	"mySparkler/pkg/utils/jwt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type GenAPI struct {
	ctx context.Context
	baseAPI.Base
	templates embed.FS
}

var genApi *GenAPI
var onceGen sync.Once

// NewApp creates a new App application struct
func NewGenAPI() *GenAPI {
	if genApi == nil {
		onceGen.Do(func() {
			genApi = &GenAPI{}
		})
	}
	return genApi
}

// setCtx 设置上下文对象
func (b *GenAPI) SetCtx(ctx context.Context, templates embed.FS) {
	b.ctx = ctx
	b.templates = templates
}

func (g *GenAPI) GenList(params map[string]interface{}) R.Result {

	var param = tools.SearchTableDataParam{
		// PageNum:  pageNum,
		// PageSize: pageSize,
		Other: tools.GenTable{
			// Table_Name:   tableName,
			// TableComment: tableComment,
		},
		// Params: tools.Params{
		// 	BeginTime: beginTime,
		// 	EndTime:   endTime,
		// },
	}
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}

	return R.ReturnSuccess(tools2.SelectGenList(param, true))
}

func (g *GenAPI) GenInfo(tableId string) R.Result {
	result := tools2.GenInfoService(tableId)
	return R.ReturnSuccess(result)
}

func (g *GenAPI) GenDbList(params map[string]interface{}) R.Result {
	var param = tools.SearchTableDataParam{
		PageNum:  1,
		PageSize: 10,
		Other:    tools.GenTable{},
		Params:   tools.Params{
			// 	BeginTime: beginTime,
			// 	EndTime:   endTime,
		},
	}

	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}

	return R.ReturnSuccess(tools2.GenDbList(param))
}

func (g *GenAPI) GenColumnInfo(context *gin.Context) {
	tableId := context.Param("tableId")
	log.Println(tableId)
}

func (g *GenAPI) ImportTable(tables string) R.Result {

	userId := jwt.CacheGetUserId()
	// var tables = context.DefaultQuery("tables", "")
	table := utils.SplitStr(tables)
	result := tools2.SelectDbTableListByNames(table)
	tools2.ImportGenTable(result, utils.GetInterfaceToInt(userId))
	return R.ReturnSuccess("成功")
}

func (g *GenAPI) GenEdit(params map[string]interface{}) R.Result {
	var req = tools.EditGenTableVO{}
	if utils.CheckTypeByReflectNil(req) {
		return R.ReturnFailMsg("获取参数失败")
	} else {

		err := mapstructure.Decode(params, &req)
		if err != nil {
			fmt.Println(err.Error())
		}
		res := tools2.UpdateGenTableService(req)
		return res
	}
}

func (g *GenAPI) GenDelete(tableIds []int) {
	// tableIds := context.Param("tableIds")
	tools2.DeleteGenTableByIds(tableIds)
	tools2.DeleteGenTableColumnByIds(tableIds)
}

func (g *GenAPI) PreviewGenTable(tableId int) R.Result {
	// tableId := context.Param("tableId")
	appPath := g.GetAppPath()
	result := tools2.PreviewGenTableCode(tableId, g.templates, appPath)
	return R.ReturnSuccess(result)
}

func (g *GenAPI) GenDownload(tableName string) {
	// tableName := context.Param("tableName")
	log.Println(tableName)
}

func (g *GenAPI) GenBatchCode(tables string, dirPath string) {
	table := utils.SplitStr(tables)
	const fileName = "Go.zip"
	appPath := g.GetAppPath()
	// dirPath := g.SaveFilePath(fileName, "")
	utils.DirExistAndDel(appPath + "/view/template/" + fileName)
	tools2.GenBatchCode(table, fileName, g.templates, appPath)
	// utils.DonwloadFile(context, fileName, filePath)

	g.OpenMacDir(appPath + "/view/template")
}

func (g *GenAPI) Gen(tableName string) R.Result {
	// tableName := context.Param("tableName")
	log.Println(tableName)
	return R.ReturnFailMsg("暂不支持")
}

// CreateTable

func (g *GenAPI) CreateTable(params map[string]string) R.Result {
	// tableName := context.Param("tableName")
	sql := params["sql"]
	log.Println(sql)
	return R.ReturnFailMsg("暂不支持")
}

// SynchDb
func (g *GenAPI) SynchDb(tableName string) R.Result {
	// tableName := context.Param("tableName")
	log.Println(tableName)
	return R.ReturnFailMsg("暂不支持")
}
