package utils

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"mySparkler/pkg/constants"
	"mySparkler/pkg/file"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/shopspring/decimal"
)

// 文件存在就删除文件
func DirExistAndDel(autoPath string) error {
	if Exists(autoPath) { // 检查文件是否存在
		if file.IsDir(autoPath) { // 检查是否是文件夹
			if err := os.RemoveAll(autoPath); err != nil {
				return err
			}
		} else {
			if err := os.Remove(autoPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
func RandNum() int {
	n, _ := rand.Int(rand.Reader, big.NewInt(100))

	return int(n.Int64())
}

// var port = 0

// 随机端口
func RandPort() int {

	// if port <= 0 {
	// 	port = RandNum() + constants.Port
	// }
	// return port
	return constants.Port
}

var timeTemplates = []string{
	"2006-01-02 15:04:05", //常规类型
	"2006/01/02 15:04:05",
	"2006/01/02 15:04",
	"2006/1/2 15:04",
	"2006-01-02",
	"2006/01/02",
	"15:04:05",
	"01/02/06 15:04",
	"1/2/06 15:04",
}

func TimeStringToBeginTime(beginTime string) time.Time {
	Loc, _ := time.LoadLocation("Asia/Shanghai")
	endTime1, _ := time.ParseInLocation(constants.TimeFormat, beginTime, Loc)
	return endTime1
}

/* 时间格式字符串转换 */
func TimeStringToGoTime(tm string) time.Time {
	Loc, _ := time.LoadLocation("Asia/Shanghai")
	for i := range timeTemplates {
		t, err := time.ParseInLocation(timeTemplates[i], tm, time.Local)
		// Log(tm + "---->" + timeTemplates[i] + "  " + t.Local().Format(constants.TimeFormat))
		if nil == err && !t.IsZero() {
			return t
		}
		t, err = time.ParseInLocation(timeTemplates[i], tm, Loc)
		// Log(tm + "---->" + timeTemplates[i] + "  " + t.Local().Format(constants.TimeFormat) + "     Asia/Shanghai")
		if nil == err && !t.IsZero() {
			return t
		}
	}
	return time.Time{}
}

// log 增加日志记录
func Log(content string) {
	fmt.Println(content)
	path := file.PathExist(fmt.Sprintf("%s/logs", file.GetAppPath()))
	// 创建或打开日志文件
	logFile, err := os.OpenFile(fmt.Sprintf("%s/%s.log", path, time.Now().Format("2006-01-02")), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	//defer logFile.Close()
	//记录文件路径和行号
	_, file, line, _ := runtime.Caller(1)
	// 初始化日志
	logger := log.New(logFile, "", log.LstdFlags)
	logger.Printf("\n文件路径：%s:%d\n日志内容：%s\n", file, line, content)
}
func GetInterfaceToInt(t1 interface{}) int {
	var t2 int
	switch t1.(type) {
	case uint:
		t2 = int(t1.(uint))
		break
	case int8:
		t2 = int(t1.(int8))
		break
	case uint8:
		t2 = int(t1.(uint8))
		break
	case int16:
		t2 = int(t1.(int16))
		break
	case uint16:
		t2 = int(t1.(uint16))
		break
	case int32:
		t2 = int(t1.(int32))
		break
	case uint32:
		t2 = int(t1.(uint32))
		break
	case int64:
		t2 = int(t1.(int64))
		break
	case uint64:
		t2 = int(t1.(uint64))
		break
	case float32:
		t2 = int(t1.(float32))
		break
	case float64:
		t2 = int(t1.(float64))
		break
	case string:
		t2, _ = strconv.Atoi(t1.(string))
		break
	default:
		t2 = t1.(int)
	}
	return t2
}

// 格式化 float 并保留3位
func GetInterfaceToFloat64(t1 interface{}) float64 {
	var t2 float64
	switch t1.(type) {
	case float64:
		t2, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", t2), 64)
		break
	case string:
		t2, _ = strconv.ParseFloat(t1.(string), 64)
		t2, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", t2), 64)
		break
	default:
		t2 = t1.(float64)
	}
	return t2
}

// AddFloat decimal类型加法
// return d1 + d2
func AddFloat(d1, d2 float64) float64 {
	decimalD1 := decimal.NewFromFloat(d1)
	decimalD2 := decimal.NewFromFloat(d2)
	decimalResult := decimalD1.Add(decimalD2)
	float64Result, _ := decimalResult.Float64()
	return float64Result
}

// SubtractFloat decimal类型减法
// return d1 - d2
func SubtractFloat(d1, d2 float64) float64 {
	decimalD1 := decimal.NewFromFloat(d1)
	decimalD2 := decimal.NewFromFloat(d2)
	decimalResult := decimalD1.Sub(decimalD2)
	float64Result, _ := decimalResult.Float64()
	return float64Result
}

// MultiplyFloat decimal类型乘法
// return d1 * d2
func MultiplyFloat(d1, d2 float64) float64 {
	decimalD1 := decimal.NewFromFloat(d1)
	decimalD2 := decimal.NewFromFloat(d2)
	decimalResult := decimalD1.Mul(decimalD2)
	float64Result, _ := decimalResult.Float64()
	return float64Result
}

// DivideFloat decimal类型除法
// return d1 / d2
func DivideFloat(d1, d2 float64) float64 {
	decimalD1 := decimal.NewFromFloat(d1)
	decimalD2 := decimal.NewFromFloat(d2)
	decimalResult := decimalD1.Div(decimalD2)
	float64Result, _ := decimalResult.Float64()
	return float64Result
}

// Round 浮点类型保留小数点后n位精度
func Round(f interface{}, n int) (r float64, err error) {
	pow10N := math.Pow10(n)
	switch f.(type) {
	case float32:
		v := reflect.ValueOf(f).Interface().(float32)
		r = math.Trunc((float64(v)+0.5/pow10N)*pow10N) / pow10N
	case float64:
		v := reflect.ValueOf(f).Interface().(float64)
		r = math.Trunc((v+0.5/pow10N)*pow10N) / pow10N
	}
	return r, err
}

// GetRemoteClientIp 获取远程客户端IP
func GetRemoteClientIp(r *http.Request) string {
	remoteIp := r.RemoteAddr

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		remoteIp = ip
	} else if ip = r.Header.Get("X-Forwarded-For"); ip != "" {
		remoteIp = ip
	} else {
		remoteIp, _, _ = net.SplitHostPort(remoteIp)
	}

	//本地ip
	if remoteIp == "::1" {
		remoteIp = "127.0.0.1"
	}

	return remoteIp
}

type Tunit struct {
	Pro  string `json:"pro"`
	City string `json:"city"`
}

func GetRealAddressByIP(ip string) string {
	url := "http://whois.pconline.com.cn/ipJson.jsp?ip=" + ip + "&json=true"
	resp, err := http.Get(url)
	var result = "内网ip"
	if err != nil {
		result = "内网ip"
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			result = "内网ip"
		} else {
			dws := new(Tunit)
			json.Unmarshal(body, &dws)
			result = dws.Pro + " " + dws.City
		}
	}
	return result
}

// 数组转 string,
func Join(sep []int, sp string) string {
	sarr := make([]string, len(sep))
	for i, v := range sep {
		sarr[i] = fmt.Sprint(v)
	}
	return strings.Join(sarr, fmt.Sprint(sp))
}

// string转数组
/* 1,1,2,3,4,5*/
func Split(data string) []int {
	var sa = strings.Split(data, ",")
	var sarr []int
	for i := 0; i < len(sa); i++ {
		var v = sa[i]
		var v1, _ = strconv.Atoi(v)
		sarr = append(sarr, v1)
	}
	return sarr
}

// string转数组
/* 1,1,2,3,4,5*/
func SplitStr(data string) []string {
	var sa = strings.Split(data, ",")
	var sarr []string
	for i := 0; i < len(sa); i++ {
		var v = sa[i]
		sarr = append(sarr, v)
	}
	return sarr
}

// 字节的单位转换 保留两位小数
func FormatFileSize(fileSize uint64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fPB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func GetDbType(dt string) string {
	switch dt {
	case "varchar":
		return "string"
	case "datetime":
		return "time.Time"
	case "timestamp":
		return "time.Time"
	case "double":
		return "float64"
	case "float":
		return "float64"
	case "decimal":
		return "float64"
	case "int":
		return "int64"
	case "tinyint":
		return "int64"
	case "smallint":
		return "int64"
	case "mediumint":
		return "int64"
	case "integer":
		return "int64"
	case "bigint":
		return "int64"
	default:
		return "string"
	}
}

// 自己写了一个判空函数，你可以直接看判空部分的逻辑就好了。
func CheckType(args ...interface{}) {
	for _, arg := range args {
		fmt.Printf("数据类型为：%s\n", reflect.TypeOf(arg).Kind().String()) //先利用反射获取数据类型，再进入不同类型的判空逻辑
		switch reflect.TypeOf(arg).Kind().String() {
		case "int":
			if arg == 0 {
				fmt.Println("数据为int,是空值")
			}
		case "string":
			if arg == "" {
				fmt.Println("数据为string，为空值")
			} else {
				fmt.Println("数据为string，数值为:", arg)
			}
		case "int64":
			if arg == 0 {
				fmt.Println("数据为int64，为空值")
			}
		case "uint8":
			if arg == false {
				fmt.Println("数据为bool，为false")
			}
		case "float64":
			if arg == 0.0 {
				fmt.Println("数据为float，为空值")
			}
		case "byte":
			if arg == 0 {
				fmt.Println("数据为byte，为0")
			}
		case "ptr":
			if arg == nil { //接口状态下，它不认为自己是nil，所以要用反射判空
				fmt.Println("数据为指针，为nil")
			} else {
				fmt.Println("数据不为空，为", arg)
			}
			//反射判空逻辑
			if reflect.ValueOf(arg).IsNil() { //利用反射直接判空
				fmt.Println("反射判断：数据为指针，为nil")
				fmt.Println("nil:", reflect.ValueOf(nil).IsValid()) //利用反射判断是否是有效值
			}
		case "struct":
			if arg == nil {
				fmt.Println("数据为struct，为空值")
			} else {
				fmt.Println("数据为struct，默认有数，无法判空，只能判断对应指针有没有初始化，直接结构体无法判断")
			}
		case "slice":
			s := reflect.ValueOf(arg)
			if s.Len() == 0 {
				fmt.Println("数据为数组/切片，为空值")
			}
		case "array":
			s := reflect.ValueOf(arg)
			if s.Len() == 0 {
				fmt.Println("数据为数组/切片，为空值")
			} else {
				fmt.Println("数据为数组/切片，为", s.Len())
			}
		default:
			fmt.Println("奇怪的数据类型")
		}
	}
}

// 判断参数为空
func CheckTypeByReflectNil(arg interface{}) bool {
	if reflect.ValueOf(arg).IsNil() { //利用反射直接判空，指针用isNil
		// 函数解释：isNil() bool	判断值是否为 nil
		// 如果值类型不是通道（channel）、函数、接口、map、指针或 切片时发生 panic，类似于语言层的v== nil操作
		fmt.Printf("反射判断：数据类型为%s,数据值为：%v,nil：%v \n",
			reflect.TypeOf(arg).Kind(), reflect.ValueOf(arg), reflect.ValueOf(arg).IsValid())
	}
	return reflect.ValueOf(arg).IsNil()
}

// 判断参数为空
func CheckTypeByReflectZero(arg interface{}) bool {
	if reflect.ValueOf(arg).IsZero() { //利用反射直接判空，基础数据类型用isZero
		fmt.Printf("反射判断：数据类型为%s,数据值为：%v,nil：%v \n",
			reflect.TypeOf(arg).Kind(), reflect.ValueOf(arg), reflect.ValueOf(arg).IsValid())
	}
	return reflect.ValueOf(arg).IsZero()
}

func GetLocalIP() string {
	var ip string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				fmt.Println(ip)
			}
		}
	}
	return ip
}
func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
