package tools

import (
	"mySparkler/pkg/utils/R"
	"os"

	"github.com/xuri/excelize/v2"
)

// 导出表格

// var (
// 	defaultSheetName = "Sheet1" //默认Sheet名称
// 	defaultHeight    = 25.0     //默认行高度
// )

type ExcelImport struct {
	// file      *excelize.File
	sheetName string //可定义默认sheet名称
}

func NewExcelImport() *ExcelImport {
	return &ExcelImport{sheetName: defaultSheetName}
}

func (l *ExcelImport) ImportToDb(filepath string) (resp R.Result) {

	file, err := os.Open(filepath)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	defer file.Close()

	data, err := l.ImportData(file)

	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Data = data
	resp.Code = 200
	return
}

// 读取数据

func (l *ExcelImport) ImportData(file *os.File) ([][]string, error) {
	// defaultFileName := createFileName()

	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		// return R.ReturnFailMsg("请选择文件")
		return nil, err
	}

	rows, err := xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))
	// fmt.Println(rows)
	if len(rows) <= 0 {
		rows, err = xlsx.GetRows(l.sheetName)
	}
	// fmt.Println(rows)
	return rows, err

}

// Letter 遍历a-z
// func Letter(length int) []string {
// 	var str []string
// 	for i := 0; i < length; i++ {
// 		str = append(str, string(rune('A'+i)))
// 	}
// 	return str
// }
