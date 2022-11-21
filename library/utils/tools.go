package utils

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/util/gconv"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/gogf/gf/crypto/gaes"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/encoding/gcharset"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/encoding/gurl"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
)

//字符串加密
func EncryptCBC(plainText, publicKey string) string {
	key := []byte(publicKey)
	b, e := gaes.EncryptCBC([]byte(plainText), key, key)
	if e != nil {
		g.Log().Error(e.Error())
		return ""
	}
	return gbase64.EncodeToString(b)
}

//字符串解密
func DecryptCBC(plainText, publicKey string) string {
	key := []byte(publicKey)
	plainTextByte, e := gbase64.DecodeString(plainText)
	if e != nil {
		g.Log().Error(e.Error())
		return ""
	}
	b, e := gaes.DecryptCBC(plainTextByte, key, key)
	if e != nil {
		g.Log().Error(e.Error())
		return ""
	}
	return gbase64.EncodeToString(b)
}

//服务端ip
func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String(), nil
	}
	return
}

//获取客户端IP
func GetClientIp(r *ghttp.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.GetClientIp()
	}
	return ip
}

//获取相差时间
func GetHourDiffer(startTime, endTime string) int64 {
	var hour int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		hour = diff / 3600
		return hour
	} else {
		return hour
	}
}

//日期字符串转时间戳（秒）
func StrToTimestamp(dateStr string) int64 {
	tm, err := gtime.StrToTime(dateStr)
	if err != nil {
		g.Log().Error(err)
		return 0
	}
	return tm.Timestamp()
}

//时间戳转 yyyy-MM-dd HH:mm:ss
func TimeStampToDateTime(timeStamp int64) string {
	tm := gtime.NewFromTimeStamp(timeStamp)
	return tm.Format("Y-m-d H:i:s")
}

//时间戳转 yyyy-MM-dd
func TimeStampToDate(timeStamp int64) string {
	tm := gtime.NewFromTimeStamp(timeStamp)
	return tm.Format("Y-m-d")
}

//获取ip所属城市
func GetCityByIp(ip string) string {
	if ip == "" {
		return ""
	}
	if ip == "[::1]" || ip == "127.0.0.1" {
		return "内网IP"
	}
	url := "http://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip
	bytes := ghttp.GetBytes(url)
	src := string(bytes)
	srcCharset := "GBK"
	tmp, _ := gcharset.ToUTF8(srcCharset, src)
	json, err := gjson.DecodeToJson(tmp)
	if err != nil {
		return ""
	}
	if json.GetInt("code") == 0 {
		city := json.GetString("city")
		return city
	} else {
		return ""
	}
}

//获取附件真实路径
func GetRealFilesUrl(r *ghttp.Request, path string) (realPath string, err error) {
	if gstr.ContainsI(path, "http") {
		realPath = path
		return
	}
	realPath, err = GetDomain(r)
	if err != nil {
		return
	}
	realPath = realPath + path
	return
}

//获取当前请求接口域名
func GetDomain(r *ghttp.Request) (string, error) {
	pathInfo, err := gurl.ParseURL(r.GetUrl(), -1)
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("解析附件路径失败")
		return "", err
	}
	return fmt.Sprintf("%s://%s:%s/", pathInfo["scheme"], pathInfo["host"], pathInfo["port"]), nil
}

//获取附件相对路径
func GetFilesPath(fileUrl string) (path string, err error) {
	if !gstr.ContainsI(fileUrl, "http") {
		path = fileUrl
		return
	}
	pathInfo, err := gurl.ParseURL(fileUrl, 32)
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("解析附件路径失败")
		return
	}
	path = gstr.TrimLeft(pathInfo["path"], "/")
	return
}

//获取reflect的方法的键
func GetReflectkeys(m map[string]reflect.Value) []string {
	l := len(m)
	keys := make([]string, 0, l)
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// ========文件名相关start========

//下划线 to 小驼峰转换
func Case2SmallCamel(name string) (str string) {
	resKey := Lcfirst(Case2Camel(name))
	return resKey
}

//下划线 to 大驼峰转换
func Case2BigCamel(name string) (str string) {
	resKey := Case2Camel(name)
	return resKey
}

//大小驼峰 to 下划线
// 驼峰式写法转为下划线写法；大小驼峰都是用这个
func Camel2Case(name string) string {
	// buffer := bytes.NewBufferString("")
	buffer := bytes.NewBuffer([]byte{})
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteRune('_')
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

// ---工具写法--
// 下划线写法转为驼峰写法
func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// 首字母大写
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// 首字母小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

//把a=1&b=2转成 "{"a":1,"b":2}"
func Query2JsonString(origin string) (s string, err error) {
	if origin == "" {
		return "", err
	}
	arr := gstr.Split(origin, "&")
	var objMap map[string]string
	objMap = make(map[string]string, len(arr))
	for _, a := range arr {
		a1 := gstr.Split(a, "=")
		objMap[a1[0]] = a1[1]
	}
	s1, err := gjson.Encode(objMap)

	return gconv.String(s1), nil
}

// ========文件名相关end========

//代码生成-注册路由
/**
 * @Description: 加入路由
 * @param objRouterStr 加入路由字符串
 * @param codeFlag 目标标记
 * @return error
 */
func RegisterGenCodeRouter(originRouterStr string, objRouterStr string, codeFlag string) error {
	//读取文件
	routerStr := gfile.GetContents(originRouterStr)
	//判断是否已有要加入的路由字符串
	pd := gstr.Contains(routerStr, objRouterStr)
	//不存在，在加入路由
	if pd == false {
		router1 := gstr.StrTillEx(routerStr, codeFlag)
		router2 := gstr.StrEx(routerStr, codeFlag)
		//插入对应字符串
		routerStr = router1 + "\n\t\t" + objRouterStr + "\n\t\t" + codeFlag + "\n\t\t" + router2
		gfile.PutContents(originRouterStr, routerStr)
	}
	return nil
}

/**
 * @Description: 生成路由。admin api的路由
 * @param tableName 输入表名，大驼峰
 * @return error
 */
func RegisterTableNameAllRouter(tableName string, serviceGenPath string) error {
	tableNameLower := Lcfirst(tableName)
	err := RegisterGenCodeRouter(serviceGenPath+"router/adminRouter.go", `group.ALL("/`+tableNameLower+`", Admin.`+tableName+`)`, "/**insert_router_end**/")
	if err != nil {
		return gerror.New("加入admin路由出错：" + `group.ALL("/` + tableNameLower + `", Admin.` + tableName + `)`)
	}
	err = RegisterGenCodeRouter(serviceGenPath+"router/apiRouter.go", `group.ALL("/`+tableNameLower+`", Api.`+tableName+`)`, "/**api_insert_router_end**/")
	if err != nil {
		return gerror.New("加入api路由出错：" + `group.ALL("/` + tableNameLower + `", Api.` + tableName + `)`)
	}

	return err
}

//把数组变成字符串，必须数组字符类型
//  ArrToString
//  @Description: 把数组变成字符串，必须数组字符类型
// arr [a,b,c]
// sp 分割符号
//  @return string a,b,v
//
func ArrToString(arr []string, sp string) string {
	var newArr = garray.NewArrayFrom(gconv.Interfaces(arr))

	newStr := newArr.Join(sp)
	return newStr
}

//
//  StringSplitToArr
//  @Description: 字符串分割成数据
// origin "a,b,c"
// sp	分割符号
//  @return []string [a,b,c]
//
func StringSplitToArr(origin string, sp string) []string {
	arr := gstr.Split(origin, sp)
	return arr
}

//
//  checkPerms
//  @Description: 判断是否有权限
// perms 已有的权限
// uri	访问的权限
//  @return bool
//
func CheckPerms(cachePerms string, uri string) bool {
	perms, err := gjson.Decode(cachePerms)
	if err != nil {
		return false
	}
	if gstr.Contains(uri, "/admin/") {
		uri = gstr.Replace(uri, "/admin/", "")
	}

	arr := gstr.Split(uri, "/")
	uri = gstr.Join(arr, ":")
	var newArr = garray.NewArrayFrom(gconv.Interfaces(perms))

	if newArr.Contains(uri) {
		return true
	}
	return false
}

//
//  checkPerms
//  @Description: 获取文件扩展名
// 	path 文件名或文件全路径
//  @return string
//
func Ext(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return strings.TrimLeft(path[i:], ".")
		}
	}
	return ""
}

//密码生成规则
func GenPwd(password, salt string) string {
	return gmd5.MustEncryptString(password + salt + password)
}

func ReverseBytes(source []byte) []byte {
	byteLen := len(source)
	array := make([]byte, byteLen)
	for i, v := range source {
		array[i] = v
	}
	for i := 0; i < byteLen/2; i++ {
		array[i], array[byteLen-i-1] = array[byteLen-i-1], array[i]
	}
	return array
}

// EncodeMD5 md5 encryption
func EncodeMD5(value, salt string) string {
	m := md5.New()
	m.Write([]byte(salt))
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

// sha256 encryption
func EncodeSHA256(value, salt string) string {
	h := sha256.New()
	h.Write([]byte(salt))
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr { // 如果是指针，则获取其所指向的元素
		t = t.Elem()
		v = v.Elem()
	}
	var data = make(map[string]interface{})
	if t.Kind() == reflect.Struct { // 只有结构体可以获取其字段信息
		for i := 0; i < t.NumField(); i++ {
			if len(t.Field(i).Tag.Get("json")) > 0 { // 获取public field
				tempField := strings.Split(t.Field(i).Tag.Get("json"), ",")
				data[tempField[0]] = v.Field(i).Interface()
			}
		}
	}
	return data
}

func Map2List(obj map[string]interface{}) []string {
	keys := make([]string, len(obj))
	var out []string
	i := 0
	for k := range obj {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		if k != "signature" {
			out = append(out, fmt.Sprintf("%v=%v", k, obj[k]))
		}
	}
	return out
}

func GenSignMsg(data interface{}) string {
	sourceStr := Map2List(Struct2Map(data))
	hashStr := strings.Join(sourceStr, "&")
	return hashStr
}

func InArray(target int, array []int) bool {
	for _, v := range array {
		if target == v {
			return true
		}
	}
	return false
}

func StrInArray(target string, array []string) bool {
	for _, v := range array {
		if target == v {
			return true
		}
	}
	return false
}

//func main() {
//	fileList := []string{
//		"/Users/Jack/Downloads/QQreceiver/72bian_log/anal.txt",
//		"/Users/Jack/Downloads/QQreceiver/72bian_log/config.json",
//		"/Users/Jack/Downloads/QQreceiver/72bian_log/xzwx_client/md5.txt",
//	}
//	//保留原来文件的结构
//	err := ZipFiles("/Users/Jack/Downloads/QQreceiver/test.zip", fileList, "","")
//	if err != nil {
//		fmt.Println(err)
//	}
//}
// ZipFiles compresses one or many files into a single zip archive file.
//压缩多个文件到一个文件里面
// Param 1: 输出的zip文件的名字
// Param 2: 需要添加到zip文件里面的文件
// Param 3: 由于file是绝对路径，打包后可能不是想要的目录，oldform就是filename中需要被替换的掉的路径
// Param 4: 要替换成的路径
func ZipFiles(filename string, files []string, oldform, newform string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 把files添加到zip中
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// 获取file的基础信息
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		//使用上面的FileInforHeader() 就可以把文件保存的路径替换成我们自己想要的了，如下面
		header.Name = strings.Replace(file, oldform, newform, -1)

		// 优化压缩
		// 更多参考see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}

func ListDir(folder string, file_list []string) ([]string, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			ListDir(folder+"/"+file.Name(), file_list)
		} else {
			// 输出绝对路径
			strAbsPath := filepath.Join(folder + "/" + file.Name())
			if err != nil {
				return nil, err
			}

			file_list = append(file_list, strAbsPath)
		}
	}

	return file_list, err
}

// 递归获取指定目录下的所有文件名
func GetAllFile(pathname string) ([]string, error) {
	result := []string{}

	fis, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Printf("读取文件目录失败，pathname=%v, err=%v \n", pathname, err)
		return result, err
	}

	// 所有文件/文件夹
	for _, fi := range fis {
		fullname := pathname + "/" + fi.Name()
		// 是文件夹则递归进入获取;是文件，则压入数组
		if fi.IsDir() {
			temp, err := GetAllFile(fullname)
			if err != nil {
				fmt.Printf("读取文件目录失败,fullname=%v, err=%v", fullname, err)
				return result, err
			}
			result = append(result, temp...)
		} else {
			result = append(result, fullname)
		}
	}

	return result, nil
}

// 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}
