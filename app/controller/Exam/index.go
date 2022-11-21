package Exam

import (
	"aiyun_cloud_srv/app/controller/exam/model/exam_answer"
	"aiyun_cloud_srv/app/controller/exam/model/exam_user"
	"aiyun_cloud_srv/library/response"
	util_jwt "aiyun_cloud_srv/library/utils/jwt"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
	"github.com/gogf/gf/util/gvalid"
	"github.com/gogf/guuid"
	"io/ioutil"
	"math/rand"
	"strings"
)

var Index = index{}

type index struct{}

type TokenInfo struct {
	AppType   int
	Uid       int
	Uuid      string
	Username  string
	PaperType string
	Subject   string
	Num       string
}

func Auth(r *ghttp.Request) {
	if r.Router.Uri == "/api/exam/login" {
		r.Middleware.Next()
		return
	}

	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		response.Json(r, 9, "Token错误")
	}

	var info *exam_user.Entity
	if err := exam_user.M.Where("token", tokenStr).Scan(&info); err != nil {
		response.ErrorDb(r, err)
	}
	if info == nil {
		response.Json(r, 9, "Token无效或过期")
	}

	data, err := util_jwt.ParseToken(tokenStr, []byte(g.Cfg().GetString("jwt.sign", "qc_sign")))
	if err != nil {
		g.Log().Error(err.Error())
		response.Json(r, 9, "Token错误")
	}

	var u *TokenInfo
	if err = gconv.Struct(data, &u); err != nil {
		g.Log().Error(err.Error())
		response.Json(r, 9, "Token错误")
	}

	r.SetCtxVar("app_type", u.AppType)
	r.SetCtxVar("uid", u.Uid)
	r.SetCtxVar("uuid", u.Uuid)
	r.SetCtxVar("username", u.Username)
	r.SetCtxVar("paper_type", u.PaperType)
	r.SetCtxVar("subject", u.Subject)
	r.SetCtxVar("num", u.Num)
	//判断权限  目前只有登录权限
	// 执行下一步请求逻辑
	r.Middleware.Next()
}

func (*index) Login(r *ghttp.Request) {
	var req *exam_user.LoginReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	password := g.Cfg().GetString("exam.pass")
	if req.Password != password {
		response.Error(r, "密码错误")
	}

	//查询用户是否存在，不存在则保存
	var user *exam_user.Entity
	if err := exam_user.M.Where("username", req.Username).Scan(&user); err != nil {
		response.ErrorDb(r, err)
	}

	if user == nil {
		if _, err := exam_user.M.Data(g.Map{"username": req.Username}).Insert(); err != nil {
			response.ErrorDb(r, err)
		}
		if err := exam_user.M.Where("username", req.Username).Scan(&user); err != nil {
			response.ErrorDb(r, err)
		}
	}

	restart := true
	uuid := guuid.New().String()
	use_time := 0
	if !req.Restart && req.AppType == 1 {
		var answer_info *exam_answer.Entity
		if err := exam_answer.M.Where(g.Map{"uid": user.Id, "paper": req.PaperType}).OrderDesc("id").Scan(&answer_info); err != nil {
			response.ErrorDb(r, err)
		}

		if answer_info != nil && 20 > answer_info.Index {
			uuid = answer_info.GroupHash
			restart = false
			use_time = UseTimeMap[user.Id]
		} else {
			delete(QuestionMap, user.Id)
		}
	}

	//生成token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"app_type":   req.AppType,
		"uid":        user.Id,
		"uuid":       uuid,
		"username":   user.Username,
		"paper_type": req.PaperType,
		"subject":    req.Subject,
		"num":        req.Num,
		"rand":       grand.Letters(20),
	}).SignedString([]byte(g.Cfg().GetString("jwt.sign", "jwt_sign")))
	if err != nil {
		response.ErrorSys(r, err)
	}

	if _, err := exam_user.M.Where("id", user.Id).Data(g.Map{"token": token}).Update(); err != nil {
		response.ErrorDb(r, err)
	}

	if req.AppType == 0 {
		delete(QuestionMap, user.Id)
	}

	response.Success(r, g.Map{"id": user.Id, "username": user.Username, "token": token, "restart": restart, "use_time": use_time})
}

var QuestionMap = make(map[int][]*Question)

func (*index) Answer(r *ghttp.Request) {
	var req *exam_answer.AnswerReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	app_type := r.GetCtxVar("app_type").Int()
	uid := r.GetCtxVar("uid").Int()
	uuid := r.GetCtxVar("uuid").String()
	username := r.GetCtxVar("username").String()
	num := r.GetCtxVar("num").Int()
	paper_type := r.GetCtxVar("paper_type").String()
	subject := r.GetCtxVar("subject").Int()

	var info *exam_answer.Entity
	if err := exam_answer.M.Where("group_hash", uuid).Order("id DESC").Scan(&info); err != nil {
		response.ErrorDb(r, err)
	}

	index := 1
	if info != nil {
		index = info.Index + 1
	}
	question_list := QuestionMap[uid]

	is_correct := 0
	if req.UserAnswer == question_list[index-1].Answer {
		is_correct = 1
	}

	if app_type == 1 {
		num = len(QuestionMap[uid])
	}
	if _, err := exam_answer.M.Data(g.Map{
		"group_hash":  uuid,
		"uid":         uid,
		"username":    username,
		"paper":       paper_type,
		"subject":     subject,
		"use_time":    req.UseTime,
		"num":         num,
		"index":       index,
		"question":    question_list[index-1].Name,
		"answer":      question_list[index-1].Answer,
		"user_answer": req.UserAnswer,
		"is_correct":  is_correct,
	}).Insert(); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, "ok")
}

func (*index) Score(r *ghttp.Request) {
	uid := r.GetCtxVar("uid").Int()
	uuid := r.GetCtxVar("uuid").String()

	var list []*exam_answer.Entity
	if err := exam_answer.M.Where(g.Map{"group_hash": uuid, "uid": uid}).Order("id").Scan(&list); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, list)
}

//生成题目
func (*index) Question(r *ghttp.Request) {
	uid := r.GetCtxVar("uid").Int()
	_, ok := QuestionMap[uid]

	paper_type := ""
	if !ok {
		num := r.GetCtxVar("num").Int()
		if num == 0 {
			num = 20
		}
		QuestionMap[uid] = genQuestion(num)
	}

	uuid := r.GetCtxVar("uuid").String()

	var info *exam_answer.Entity
	if err := exam_answer.M.Where("group_hash", uuid).Order("id DESC").Scan(&info); err != nil {
		response.ErrorDb(r, err)
	}
	index := 1
	if info != nil {
		index = info.Index + 1
	}

	question_list := QuestionMap[uid]
	if index > len(question_list) {
		response.Success(r, g.Map{
			"total":       len(question_list),
			"index":       index,
			"file_type":   question_list[len(question_list)-1].FileType,
			"question":    question_list[len(question_list)-1].Name,
			"question_ai": question_list[len(question_list)-1].AiName,
			"paper_type":  paper_type,
			"finish":      true,
			"type":        question_list[len(question_list)-1].Type,
		})
	} else {
		response.Success(r, g.Map{
			"total":       len(question_list),
			"index":       index,
			"file_type":   question_list[index-1].FileType,
			"question":    question_list[index-1].Name,
			"question_ai": question_list[index-1].AiName,
			"paper_type":  paper_type,
			"finish":      false,
			"type":        question_list[index-1].Type,
		})
	}
}

type Question struct {
	Name     string
	AiName   string
	Type     int // 1：无: 2：同 3：异
	Answer   string
	FileType int // 0 ：图片   1： 视频
}

// 获取题目
//limit 0 获取全部数据
func genQuestion(limit int) (question_list []*Question) {
	// 读取图片文件
	for i := 1; i <= 3; i++ {
		url := fmt.Sprintf("exam/image/image/image_%d", i)
		ai_url := fmt.Sprintf("exam/image/imageAI/image_%d", i)
		path := "./public/" + url
		dirs, _ := ioutil.ReadDir(path)
		for _, f := range dirs {
			chil_dir := fmt.Sprintf("%s/%s/", path, f.Name())
			chilFiles, _ := ioutil.ReadDir(chil_dir)
			for _, c := range chilFiles {
				question_list = append(question_list, &Question{
					Name:     url + "/" + f.Name() + "/" + c.Name(),
					AiName:   ai_url + "/" + f.Name() + "/" + strings.Replace(c.Name(), ".jpg", "_AI.jpg", 1),
					Type:     i,
					Answer:   f.Name(),
					FileType: 0,
				})
			}

		}
	}

	// 读取视频文件
	for i := 1; i <= 3; i++ {
		url := fmt.Sprintf("exam/video/video/video_%d", i)
		ai_url := fmt.Sprintf("exam/video/videoAI/video_%d", i)
		path := "./public/" + url
		dirs, _ := ioutil.ReadDir(path)
		for _, f := range dirs {
			chil_dir := fmt.Sprintf("%s/%s/", path, f.Name())
			chilFiles, _ := ioutil.ReadDir(chil_dir)
			for _, c := range chilFiles {
				question_list = append(question_list, &Question{
					Name:     url + "/" + f.Name() + "/" + c.Name(),
					AiName:   ai_url + "/" + f.Name() + "/" + c.Name(),
					Type:     i,
					Answer:   f.Name(),
					FileType: 1,
				})
			}

		}
	}

	// 打乱数组
	rand.Shuffle(len(question_list), func(i, j int) {
		question_list[i], question_list[j] = question_list[j], question_list[i]
	})

	if limit > 0 {
		question_list = question_list[:limit]
	}
	return
}

var UseTimeMap = make(map[int]int)

//退出 记录use_time
func (*index) Logout(r *ghttp.Request) {
	uid := r.GetCtxVar("uid").Int()
	use_time := r.GetQueryInt("use_time")
	UseTimeMap[uid] = use_time
	response.Success(r)
}
