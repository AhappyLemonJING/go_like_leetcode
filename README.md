## 技术栈

后台：Gin、GORM

前台：Vue、ElementUI

## 内容（调用顺序从上往下）

#### main

* 调用router，并开启程序执行

#### router

* 设置访问路径以及service中对应的相关的调用方法

#### service

* 对models中对应的处理结果进行反馈
  * 前端显示或者后台输出

#### models

* 创建表单
* 对数据库的各种操作（增删改查）

#### test（测试部分）

* 单元测试等

### 配置数据库和redis

```go
models/init.go

package models

import (
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB = Init()

var RDB = InitRedisDB()

func Init() *gorm.DB {
	dsn := "root:981122wzj@tcp(127.0.0.1:3306)/gin_gorm_oj?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("gorm init error:", err)
	}
	return db
}

func InitRedisDB() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
}

```

### 配置swagger

```go
// shell
go get "github.com/swaggo/gin-swagger"
swag init

// router中导入包
_ "gin_gorm_oj/docs"
swaggerfiles "github.com/swaggo/files"
ginSwagger "github.com/swaggo/gin-swagger"


func Router() *gin.Engine {
	r := gin.Default()

	// swagger配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

```

### 设计表单

### UserBasic

```go
package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model
	Identity  string `gorm:"column:identity;type:varchar(36);" json:"identity"`       // 用户的唯一标识
	Name      string `gorm:"column:name;type:varchar(100);" json:"name"`              // 姓名
	Password  string `gorm:"column:password;type:varchar(32);" json:"password"`       // 密码
	Phone     string `gorm:"column:phone;type:varchar(20);" json:"phone"`             // 电话
	Mail      string `gorm:"column:mail;type:varchar(100);" json:"mail"`              // 邮箱
	PassNum   int64  `gorm:"column:finish_problem_num;type:int(11);" json:"pass_num"` // 通过个数
	SubmitNum int64  `gorm:"column:submit_num;type:int(11);" json:"submit_num"`       // 提交次数
	IsAdmin   int    `gorm:"column:is_admin;type:tinyint(1);" json:"is_admin"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

```

### ProblemBasic

```go
package models

import (
	"gorm.io/gorm"
)

type ProblemBasic struct {
	gorm.Model
	Identity          string             `gorm:"column:identity;type:varchar(36);" json:"identity"` // 问题的唯一标识
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:problem_id;references:id"`
	Title             string             `gorm:"column:title;type:varchar(255);" json:"title"` // 题目的标题
	Content           string             `gorm:"column:content;type:text;" json:"content"`     // 题目正文描述
	MaxMem            int                `gorm:"column:max_mem;type:int;" json:"max_mem"`
	MaxRuntime        int                `gorm:"column:max_runtime;type:int;" json:"max_runtime"`
	TestCase          []*TestCase        `gorm:"foreignKey:problem_identity;references:identity"`
	PassNum           int64              `gorm:"column:pass_num;type:int(11);" json:"pass_num"`     // 通过个数
	SubmitNum         int64              `gorm:"column:submit_num;type:int(11);" json:"submit_num"` // 提交次数
}

func (table *ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword string, categoryIdentity string) *gorm.DB {
	tx := DB.Model(new(ProblemBasic)).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").Where("title like ? OR content like ?", "%"+keyword+"%", "%"+keyword+"%")

	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc on pc.problem_id = problem_basic.id").Where("pc.category_id = (SELECT cb.id FROM category_basic cb WHERE cb.identity = ?)", categoryIdentity)
	}
	return tx
}

```

### CategoryBasic

```go
package models

import (
	"gorm.io/gorm"
)

type ProblemBasic struct {
	gorm.Model
	Identity          string             `gorm:"column:identity;type:varchar(36);" json:"identity"` // 问题的唯一标识
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:problem_id;references:id"`
	Title             string             `gorm:"column:title;type:varchar(255);" json:"title"` // 题目的标题
	Content           string             `gorm:"column:content;type:text;" json:"content"`     // 题目正文描述
	MaxMem            int                `gorm:"column:max_mem;type:int;" json:"max_mem"`
	MaxRuntime        int                `gorm:"column:max_runtime;type:int;" json:"max_runtime"`
	TestCase          []*TestCase        `gorm:"foreignKey:problem_identity;references:identity"`
	PassNum           int64              `gorm:"column:pass_num;type:int(11);" json:"pass_num"`     // 通过个数
	SubmitNum         int64              `gorm:"column:submit_num;type:int(11);" json:"submit_num"` // 提交次数
}

func (table *ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword string, categoryIdentity string) *gorm.DB {
	tx := DB.Model(new(ProblemBasic)).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").Where("title like ? OR content like ?", "%"+keyword+"%", "%"+keyword+"%")

	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc on pc.problem_id = problem_basic.id").Where("pc.category_id = (SELECT cb.id FROM category_basic cb WHERE cb.identity = ?)", categoryIdentity)
	}
	return tx
}

```

### ProblemCategory

```go
package models

import "gorm.io/gorm"

type ProblemCategory struct {
	gorm.Model
	ProblemId     uint           `gorm:"column:problem_id;type:varchar(36);" json:"problem_id"` // 问题的id
	CategoryId    uint           `gorm:"column:category_id;type:varchar(36);" json:"category_id"`
	CategoryBasic *CategoryBasic `gorm:"foreignKey:id;references:category_id"`
}

func (table *ProblemCategory) TableName() string {
	return "problem_category"
}

```

### SubmitBasic

```go
package models

import "gorm.io/gorm"

type SubmitBasic struct {
	gorm.Model
	Identity        string        `gorm:"column:identity;type:varchar(36);" json:"identity"`
	ProblemIdentity string        `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"`
	ProblemBasic    *ProblemBasic `gorm:"foreignKey:identity;references:problem_identity"`
	UserIdentity    string        `gorm:"column:user_identity;type:varchar(36);" json:"user_identity"`
	UserBasic       *UserBasic    `gorm:"foreignKey:identity;references:user_identity"`
	Path            string        `gorm:"column:path;type:varchar(255);" json:"path"`
	Status          int           `gorm:"column:status;type:tinyint(1);" json:"tinyint"`
}

func (table *SubmitBasic) TableName() string {
	return "submit_basic"
}

func GetSubmitList(problemIdentity string, userIdentity string, status int) *gorm.DB {
	tx := DB.Model(new(SubmitBasic)).Preload("ProblemBasic", func(db *gorm.DB) *gorm.DB {
		return db.Omit("content")
	}).Preload("UserBasic")

	if problemIdentity != "" {
		tx.Where("problem_identity = ?", problemIdentity)
	}
	if userIdentity != "" {
		tx.Where("user_identity = ?", userIdentity)
	}
	if status != 0 {
		tx.Where("status = ?", status)
	}
	return tx
}

```

### TestCase

```go
package models

import "gorm.io/gorm"

type TestCase struct {
	gorm.Model
	Identity        string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	ProblemIdentity string `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"`
	Input           string `gorm:"column:input;type:text;" json:"input"`
	Output          string `gorm:"column:output;type:text;" json:"output"`
}

func (table *TestCase) TableName() string {
	return "test_case"
}

```

### 获取题目列表

#### 1. 配置路由

```go
r.GET("/problem-list", service.GetProblemList) // 配置获取题目列表的路径及方法
```

#### 2. service包与models包中创建对应的方法

* 配置swagger：设置一些需要从前端获取的数据
* 方法实现：
  * 该方法主要实现根据keyword 和categoryIdentity从ProblemBasic数据表中查找相应的数据。
  * 获取前端输入的数据（包括分页设置避免数据太多全部显示、keyword、category_identity）
  * 调用`models.GetProblemList(keyword, categoryIdentity)`到数据库进行处理
  * 找到对应的数据进行返回到service中并显示

```go
service/problem.go

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "请输入当前页面，默认第一页"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-list [get]
func GetProblemList(ctx *gin.Context) {
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(ctx.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("get problem list page parse error:", err)
		return
	}
	page = (page - 1) * size
	var count int64
	keyword := ctx.Query("keyword")
	categoryIdentity := ctx.Query("category_identity")

	list := make([]*models.ProblemBasic, 0)
	tx := models.GetProblemList(keyword, categoryIdentity)
	err = tx.Count(&count).Omit("content").Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("get problem list error:", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}
```

### 获取题目详情

#### 1. 配置路由

```go
r.GET("/problem-detail", service.GetProblemDetail)
```

#### 2. service包中创建对应的方法

由于对应的数据库处理比较简单，因此没有单独放在models/problem_basic.go中

* 问题详情只需要前端给到问题的identity，并进行显示即可
* 拿到前端的identity，并从数据库中的Problem Basic表中比较有没有一样的identity，找到则进行返回

```go
service/problem.go

// GetProblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @Param identity query string false "problem identity"
// @Success 200 {string} json "{"code":"200","msg":"","data":""}"
// @Router /problem-detail [get]
func GetProblemDetail(ctx *gin.Context) {
	identity := ctx.Query("identity")
	if identity == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题唯一标识不能为空",
		})
		return
	}
	data := new(models.ProblemBasic)
	err := models.DB.Where("identity = ?", identity).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "当前问题不存在",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get problemDetail Error:" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}
```

### 获取用户详情

#### 1. 配置路由

```go
r.GET("/user-detail", service.GetUserDetail)
```

#### 2. service包中创建对应的方法

这部分和获取题目详情类似

* 前端只需传identity参数，后台拿到这个参数再从数据库中查找并返回

```go
// GetUserDetail
// @Tags 公共方法
// @Summary 用户详情
// @Param identity query string false "user identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user-detail [get]
func GetUserDetail(ctx *gin.Context) {
	identity := ctx.Query("identity")
	if identity == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户唯一标识不能为空",
		})
		return
	}
	data := new(models.UserBasic)
	err := models.DB.Omit("password").Where("identity = ?", identity).Find(&data).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get usersDetail Error:" + err.Error() + "identity:" + identity,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}
```

### 登陆

#### 1. 配置路由

```go
r.POST("/login", service.Login)
```

#### 2. 登陆和注册需要对密码进行加密

```go
helper/helper.go

// 生成Md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
```

#### 3. 登陆token的生成，该token表示登陆状态，也用作鉴权

* 传入identity，name，isAdmin
* 设置myKey作为token的种子，让每个identity任何时候生成的token都一样
* 使用jwt生成token

```go
helper/helper.go
 
go get "github.com/dgrijalva/jwt-go"  // 引入jwt生成token

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}
var myKey = []byte("gin-gorm-oj-key")
func GenerateToken(identity string, name string, isAdim int) (string, error) {
	userClaim := &UserClaims{
		Identity:       identity,
		Name:           name,
		IsAdmin:        isAdim,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
```

#### 4. service包中创建对应的方法

* 前端输入用户名和密码
* 后台拿到用户名和密码，对密码进行md5的加密再和数据库中进行匹配（因为数据库中存储的是加密之后的密码）
* 然后调用`helper.GenerateToken(data.Identity, data.Name, data.IsAdmin)`生成该用户的token表示该用户登陆成功，该token可用于后续的鉴权

```go
// Login
// @Tags 公共方法
// @Summary 用户登陆
// @Param username formData string false "username"
// @Param password formData string false "password"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /login [post]
func Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	if username == "" || password == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码不能为空",
		})
	}
	// md5
	password = helper.GetMd5(password)
	data := new(models.UserBasic)
	err := models.DB.Where("name = ? and password = ?", username, password).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户名或密码错误",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get user error:" + err.Error(),
		})
		return
	}
	token, err := helper.GenerateToken(data.Identity, data.Name, data.IsAdmin)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GenerateToken err :" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg": map[string]interface{}{
			"token": token,
		},
	})

}
```

### 用户验证码

#### 1. 配置路由

```go
r.POST("/send-code", service.SendCode)
```

#### 2. 随机生成验证码

* 使用当前时间作为随机种子，随机6位0-9的数字

```go
helper/helper.go

// 生成验证码
func GetRand() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s = s + strconv.Itoa(rand.Intn(10))
	}
	return s
}
```

#### 3. 给指定的邮箱发验证码

* 设置发送验证码的管理员邮箱，从该邮箱发出验证码

```go
helper/helper.go

// 发送验证码
func SendCode(toUserEmal, code string) error {
	e := email.NewEmail()
	e.From = "wzj <wzj2010624@163.com>"
	e.To = []string{toUserEmal}
	e.Subject = "验证码已发送，请查收"
	e.HTML = []byte("你的验证码是：<b>" + code + "</b>")
	err := e.SendWithTLS("smtp.163.com:587", smtp.PlainAuth("", "wzj2010624@163.com", "RUSCZFDRNLMUYJZA", "smtp.163.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	return err
}
```

#### 4. service包中创建对应的方法

* 获取前端的数据，这里只需要指定邮箱即可
* 调用`helper.GetRand()`生成随机验证码
* 将指定邮箱和验证码发布到redis进行缓存`models.RDB.Set(ctx, email, code, time.Second*300)`设置300秒过期，用于后期注册阶段校验验证码是否正确
* 调用`helper.SendCode(email, code)`发送验证码

```go
// SendCode
// @Tags 公共方法
// @Summary 发送验证码
// @Param email formData string false "email"
// @Success 200 {string} json "{"code":"200","msg":""}"
// @Router /send-code [post]
func SendCode(ctx *gin.Context) {
	email := ctx.PostForm("email")
	if email == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱为空，无法发送！",
		})
		return
	}
	code := helper.GetRand()
	models.RDB.Set(ctx, email, code, time.Second*300)
	err := helper.SendCode(email, code)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "send code error:" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "发送成功",
	})
}
```

### 注册

#### 1. 配置路由

```go
r.POST("/register", service.Register)
```

#### 2. 自动生成uuid

```go
helper/helper.go

// 生成uuid
func GetUUID() string {
	return uuid.NewV4().String()
}
```

#### 3. service包中创建对应的方法

* 获取前端输入的注册邮箱、验证码、姓名、密码、手机号
* 调用`models.RDB.Get(ctx, mail).Result()`获取发布到redis上的验证码，将其和前端获取的验证码进行比对，查验是否一致
* 判断邮箱有没有已经被注册（在UserBasic表中的mail字段查询）
* 调用`helper.GetUUID()`以及`helper.GetMd5(password)`以及前端获取到的信息生成用户数据，并插入UserBasic表中
* 调用`helper.GenerateToken`直接生成token，变成登陆状态

```go
// Register
// @Tags 公共方法
// @Summary 用户注册
// @Param mail formData string true "mail"
// @Param code formData string true "code"
// @Param name formData string true "name"
// @Param password formData string true "password"
// @Param phone formData string false "phone"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /register [post]
func Register(ctx *gin.Context) {
	mail := ctx.PostForm("mail")
	userCode := ctx.PostForm("code")
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")
	phone := ctx.PostForm("phone")
	if mail == "" || userCode == "" || name == "" || password == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	// 验证码是否正确
	sysCode, err := models.RDB.Get(ctx, mail).Result()
	if err != nil {
		log.Println("get code err :", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "code get error:" + err.Error(),
		})
		return
	}
	if sysCode != userCode {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码不正确",
		})
		return
	}

	// 判断邮箱是否已经注册
	var cnt int64
	err = models.DB.Where("mail = ?", mail).Model(new(models.UserBasic)).Count(&cnt).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get user error:" + err.Error(),
		})
		return
	}
	if cnt > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该邮箱已经被注册",
		})
		return
	}

	// 数据插入 password生成md5
	userIdentity := helper.GetUUID()
	data := &models.UserBasic{
		Identity: userIdentity,
		Name:     name,
		Password: helper.GetMd5(password),
		Mail:     mail,
		Phone:    phone,
	}
	err = models.DB.Create(data).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "create user error:" + err.Error(),
		})
		return
	}
	// 生成token
	token, err := helper.GenerateToken(userIdentity, name, data.IsAdmin)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "generate token error:" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": token,
		},
	})

}

```

### 用户排行榜

#### 1. 配置路由

```go
r.GET("/rank-list", service.GetRankList)
```

#### 2. service包中创建对应的方法

* 和问题列表、用户列表类似，根据用户的通过数量和提交数量进行排行显示

```go
// GetRankList
// @Tags 公共方法
// @Summary 用户排行榜
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /rank-list [get]
func GetRankList(ctx *gin.Context) {
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(ctx.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("get rank list page parse error:", err)
		return
	}
	page = (page - 1) * size
	var count int64
	list := make([]*models.UserBasic, 0)
	err = models.DB.Model(new(models.UserBasic)).Count(&count).Order("pass_num DESC, submit_num ASC").Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "getRankList error" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

```

### 提交列表

#### 1. 配置路由

```go
r.GET("/submit-list", service.GetSubmitList)
```

#### 2. service包中创建对应的方法

* 获取前端输入的数据，包括problem_identity、user_identity、status
  * 查看某道题目的提交列表或者是某个用户的提交列表或者是某种状态的提交列表
* 调用`models.GetSubmitList(problemIdentity, userIdentity, status)`查看数据库并返回数据

```go
// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
// @Param page query int false "请输入当前页面，默认第一页"
// @Param size query int false "size"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query string false "user_identity"
// @Param status query int false "status"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /submit-list [get]
func GetSubmitList(ctx *gin.Context) {
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(ctx.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("get problem list page parse error:", err)
		return
	}
	page = (page - 1) * size
	var count int64
	problemIdentity := ctx.Query("problem_identity")
	userIdentity := ctx.Query("user_identity")
	status, _ := strconv.Atoi(ctx.Query("status"))

	list := make([]*models.SubmitBasic, 0)
	tx := models.GetSubmitList(problemIdentity, userIdentity, status)

	err = tx.Count(&count).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("Get submitlist error:", err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get submitlist error:" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})

}
```

### 问题创建（这部分开始是管理员私有方法）

#### 1. 设计中间件用于鉴权

* 获取header中存储的token`ctx.GetHeader("Authorization")`
* 对该token进行解析`helper.AnalyseToken(auth)`
* 解析出有权限则可进行下一步操作

```go
middlewares/auth_admin.go

package middlewares

import (
	"gin_gorm_oj/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthAdminCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: check if user is admin
		auth := ctx.GetHeader("Authorization")
		userClaim, err := helper.AnalyseToken(auth)
		if err != nil {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			return
		}
		if userClaim == nil || userClaim.IsAdmin != 1 {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			return
		}
		ctx.Next()
	}
}

// 解析token
helper/helper.go

func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(t *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse token error:%v", err)
	}
	return userClaim, nil
}

```

#### 2. 配置路由

```go
authAdmin := r.Group("/admin", middlewares.AuthAdminCheck())
authAdmin.POST("/problem-create", service.ProblemCreate)
```

#### 3. service包中创建对应的方法

* 前端输入token（用于鉴别该登陆的用户有没有管理员权限）、问题的基本信息
* 后台获取前端输入的信息
* 调用`helper.GetUUID()`生成该问题的uuid
* 将问题的基本信息封装成ProblemBasic的格式存入data中
* 获取前端输入的该问题的分类类别，将类别和问题关联起来的条目添加到data.ProblemCategory中
* 获取前端输入的测试用例`{"input":"1 2\n","output":"3\n"}`将该字符串使用`json.Unmarshal([]byte(testCase), &caseMap)`的方式转化为map存入caseMap中，调用`helper.GetUUID(),`生成该测试用例的uuid，将caseMap的input和output一起封装成一个条目存入到data.TestCase表中
* 将data存入数据库，顺便ProblemCategory和TestCase也一并更新

```go
// ProblemCreate
// @Tags 管理员私有方法
// @Summary 问题创建
// @Param authorization header string true "authorization"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_mem formData int true "max_mem"
// @Param max_runtime formData int true "max_runtime"
// @Param category_ids formData []string false "category_ids" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200","msg":"","data":""}"
// @Router /admin/problem-create [post]
func ProblemCreate(ctx *gin.Context) {
	title := ctx.PostForm("title")
	content := ctx.PostForm("content")
	maxMem, _ := strconv.Atoi(ctx.PostForm("max_mem"))
	maxRuntime, _ := strconv.Atoi(ctx.PostForm("max_runtime"))
	categoryIds := ctx.PostFormArray("category_ids")
	testCases := ctx.PostFormArray("test_cases")

	if title == "" || content == "" || len(testCases) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	identity := helper.GetUUID()
	data := models.ProblemBasic{
		Title:      title,
		Content:    content,
		MaxMem:     maxMem,
		MaxRuntime: maxRuntime,
		Identity:   identity,
	}
	// 处理分类
	categoryBasic := make([]*models.ProblemCategory, 0)
	for _, id := range categoryIds {
		intId, _ := strconv.Atoi(id)
		categoryBasic = append(categoryBasic, &models.ProblemCategory{
			ProblemId:  data.ID,
			CategoryId: uint(intId),
		})
	}
	data.ProblemCategories = categoryBasic

	// 处理测试用例
	testCaseBasics := make([]*models.TestCase, 0)
	for _, testCase := range testCases {
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(testCase), &caseMap)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式错误",
			})
			return
		}
		if _, ok := caseMap["input"]; !ok {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式错误 input",
			})
			return
		}
		if _, ok := caseMap["output"]; !ok {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式错误 output",
			})
			return
		}
		testCaseBasic := &models.TestCase{
			Identity:        helper.GetUUID(),
			ProblemIdentity: identity,
			Input:           caseMap["input"],
			Output:          caseMap["output"],
		}
		testCaseBasics = append(testCaseBasics, testCaseBasic)

	}
	data.TestCase = testCaseBasics

	// 创建问题
	err := models.DB.Create(&data).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "problem create err:" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"identity": data.Identity,
		},
	})

}
```

### 问题修改

#### 1. 配置路由

```go
authAdmin.PUT("/problem-modify", service.ProblemMotify)
```

#### 2. service包中创建对应的方法

* 前端同样输入token用于鉴权
* 从前端获取问题identity进行该问题的修改，修改内容包括title、content、max_mem、max_runtime、category_ids、test_cases
* 先将前端获取的基础信息保存到problemBasic中，根据identity查找数据库中的条目进行更新，并将该条目更新后保存到problemBasic让后台获取
* 对问题的分类进行修改：
  * 删除原来的问题分类：在ProblemCategory表中找到对应的problemBasic.ID进行删除
  * 添加新的问题分类：遍历前端输入的categoryIds，将problemId与categoryIds关联成Problem Category条目，进行存储
* 对测试用例进行修改：
  * 删除原来的测试用例：在TestCase表中根据问题的problem_identity进行删除条目
  * 添加新的测试用例：遍历前端输入的testCases，分别使用`json.Unmarshal([]byte(testCase), &caseMap)`的方式转化为map存入caseMap中，调用`helper.GetUUID(),`生成该测试用例的uuid，将caseMap的input和output一起封装成一个条目存入到TestCase表中

```go
// ProblemModify
// @Tags 管理员私有方法
// @Summary 问题修改
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_mem formData int true "max_mem"
// @Param max_runtime formData int true "max_runtime"
// @Param category_ids formData []string false "category_ids" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200","msg":"","data":""}"
// @Router /admin/problem-modify [put]
func ProblemMotify(ctx *gin.Context) {
	identity := ctx.PostForm("identity")
	title := ctx.PostForm("title")
	content := ctx.PostForm("content")
	maxMem, _ := strconv.Atoi(ctx.PostForm("max_mem"))
	maxRuntime, _ := strconv.Atoi(ctx.PostForm("max_runtime"))
	categoryIds := ctx.PostFormArray("category_ids")
	testCases := ctx.PostFormArray("test_cases")

	if identity == "" || title == "" || content == "" || len(testCases) == 0 || maxMem == 0 || maxRuntime == 0 || len(categoryIds) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 问题基础信息保存
		problemBasic := &models.ProblemBasic{
			Identity:   identity,
			Title:      title,
			Content:    content,
			MaxMem:     maxMem,
			MaxRuntime: maxRuntime,
		}
		err := tx.Where("identity = ?", identity).Updates(problemBasic).Error
		if err != nil {
			return err
		}
		// 查询问题详情
		err = tx.Where("identity = ?", identity).Find(problemBasic).Error
		if err != nil {
			return err
		}

		// 关联问题分类的保存
		// 1. 删除已存在的关联关系
		err = tx.Where("problem_id=?", problemBasic.ID).Delete(new(models.ProblemCategory)).Error
		if err != nil {
			return err
		}
		// 2. 新增新的关联关系
		pcs := make([]*models.ProblemCategory, 0)
		for _, id := range categoryIds {
			intid, _ := strconv.Atoi(id)
			procat := &models.ProblemCategory{
				ProblemId:  problemBasic.ID,
				CategoryId: uint(intid),
			}
			pcs = append(pcs, procat)
		}
		err = tx.Model(new(models.ProblemCategory)).Create(&pcs).Error
		if err != nil {
			return err
		}

		// 关联测试用例的保存
		// 1. 删除已存在的关联关系
		err = tx.Where("problem_identity = ?", identity).Delete(new(models.TestCase)).Error
		if err != nil {
			return err
		}
		// 2. 增加新的关联关系
		tcs := make([]*models.TestCase, 0)
		for _, testCase := range testCases {
			caseMap := make(map[string]string)
			err = json.Unmarshal([]byte(testCase), &caseMap)
			if err != nil {
				return err
			}
			if _, ok := caseMap["input"]; !ok {
				return errors.New("测试案例格式错误")
			}
			if _, ok := caseMap["output"]; !ok {
				return errors.New("测试案例格式错误")
			}
			tcs = append(tcs, &models.TestCase{
				Identity:        helper.GetUUID(),
				ProblemIdentity: identity,
				Input:           caseMap["input"],
				Output:          caseMap["output"],
			})

		}
		err = tx.Create(&tcs).Model(new(models.TestCase)).Error
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题修改失败,err :" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "问题修改成功",
	})

}
```

### 分类列表

#### 1. 配置路由

```go
authAdmin.GET("/category-list", service.GetCategoryList)
```

#### 2. service包中创建对应的方法

* 通过前端输入的keyword查找类别

```go
// GetCategoryList
// @Tags 管理员私有方法
// @Summary 分类列表
// @Param authorization header string true "authorization"
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-list [get]
func GetCategoryList(ctx *gin.Context) {
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(ctx.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("get category list page parse error:", err)
		return
	}
	page = (page - 1) * size
	var count int64
	keyword := ctx.Query("keyword")

	categorylist := make([]*models.CategoryBasic, 0)
	err = models.DB.Model(new(models.CategoryBasic)).Where("name like ?", "%"+keyword+"%").Count(&count).Offset(page).Limit(size).Find(&categorylist).Error

	if err != nil {
		log.Println("get category list error:", err)
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "获取分类列表失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"categorylist": categorylist,
			"count":        count,
		},
	})

}
```

### 分类创建

#### 1. 配置路由

```go
authAdmin.POST("/category-create", service.CategoryCreate)
```

#### 2. service包中创建对应的方法

* 前端传入token鉴权，看看有没有创建分类的权限
* 后台获取前端输入的需要创建的类别名称和父级id
* 调用`helper.GetUUID(),`生成新创建的类别的uuid，与前端输入的名称和父级id一起封装成CategoryBasic条目进行存储

```go
// CategoryCreate
// @Tags 管理员私有方法
// @Summary 分类创建
// @Param authorization header string true "authorization"
// @Param name formData string true "name"
// @Param parentId formData int true "parentId"
// @Success 200 {string} json "{"code":"200","msg":"","data":""}"
// @Router /admin/category-create [post]
func CategoryCreate(ctx *gin.Context) {
	name := ctx.PostForm("name")
	parentId, _ := strconv.Atoi(ctx.PostForm("parentId"))
	category := &models.CategoryBasic{
		Identity: helper.GetUUID(),
		Name:     name,
		ParentId: parentId,
	}
	err := models.DB.Create(category).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "创建分类失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})

}
```

### 分类修改

#### 1. 配置路由

```go
authAdmin.PUT("/category-modify", service.CategoryModify)
```

#### 2. service包中创建对应的方法

* 前端输入token进行鉴权
* 后台获取前端输入的需要修改的分类的唯一标识identity，以及修改内容
* 先将前端获取的数据封装成CategoryBasic条目，并且根据identity在数据库中找到对应的条目进行修改
* 

```go
// CategoryModify
// @Tags 管理员私有方法
// @Summary 分类修改
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param name formData string true "name"
// @Param parentId formData int true "parentId"
// @Success 200 {string} json "{"code":"200","msg":"","data":""}"
// @Router /admin/category-modify [put]
func CategoryModify(ctx *gin.Context) {
	name := ctx.PostForm("name")
	identity := ctx.PostForm("identity")
	parentId, _ := strconv.Atoi(ctx.PostForm("parentId"))
	if name == "" || identity == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	category := &models.CategoryBasic{
		Identity: identity,
		Name:     name,
		ParentId: parentId,
	}
	err := models.DB.Model(new(models.CategoryBasic)).Where("identity = ?", identity).Updates(category).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "分类修改失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分类修改成功",
	})

}
```

### 分类删除

#### 1. 配置路由

```go
authAdmin.DELETE("/category-delete", service.CategoryDelete)
```

#### 2. service包中创建对应的方法

* 前端输入token进行鉴权
* 后台获取前端输入的需要删除的分类的identity
* 通过查找ProblemCategory表中有没有题目是该分类下的，如果有则无法删除
* 如果可以删除，则根据identity在CategoryBasic表中找到对应的条目进行删除

```go
// CategoryDelete
// @Tags 管理员私有方法
// @Summary 分类删除
// @Param authorization header string true "authorization"
// @Param identity query string true "identity"
// @Success 200 {string} json "{"code":"200","msg":"","data":""}"
// @Router /admin/category-delete [delete]
func CategoryDelete(ctx *gin.Context) {
	identity := ctx.Query("identity")
	if identity == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确,identity",
		})
		return
	}
	var cnt int64
	err := models.DB.Model(new(models.ProblemCategory)).Where("category_id = (SELECT id from category_basic WHERE identity = ? LIMIT 1)", identity).Count(&cnt).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取分类关联的问题失败",
		})
		return
	}
	if cnt > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该分类下有题目，不能删除",
		})
		return
	}
	err = models.DB.Model(new(models.CategoryBasic)).Where("identity = ?", identity).Delete(&models.CategoryBasic{}).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除分类失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分类删除成功",
	})
}

```

### 代码提交（这个是用户私有方法）

#### 1. 设计中间件用于鉴权

* 解析前端拿到的token`helper.AnalyseToken(auth)`，如果有token则`ctx.Set("user", userClaim)`，用于记录该用户，方便后台获取该用户的信息

```go
middlewares/auth_user.go

package middlewares

import (
	"gin_gorm_oj/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthUserCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: check if user is not admin
		auth := ctx.GetHeader("Authorization")
		userClaim, err := helper.AnalyseToken(auth)
		if err != nil {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			return
		}
		if userClaim == nil {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			return
		}

		ctx.Set("user", userClaim)
		ctx.Next()
	}
}

```

#### 2. 配置路由

```go
authUser := r.Group("/user", middlewares.AuthUserCheck())
authUser.POST("/submit", service.Submit)
```

#### 3. service包中创建对应的方法

* 前端输入token进行鉴权，输入问题的唯一标识identity、输入该问题代码
* 代码的输入方式是body，后台通过`ioutil.ReadAll(ctx.Request.Body)`来获取代码
* 保存代码`helper.CodeSave(code)`到path路径下
* 后台获取当前用户的信息`ctx.Get("user")`
* 创建SubmitBasic条目，使用`helper.GetUUID(),`自动生成该条目的uuid
* 通过ProblemBasic关联的TestCase找到问题与测试用例
* 设置三个通道分别表示三种代码状态，超内存、错误、编译不通过，定义互斥锁，遍历该问题的TestCase，通过协程执行测试，使用`exec.Command("go", "run", path)`运行path路径下的代码，根据测试的输入案例进行运行拿到输出结果和标准输出结果是否匹配，并且在运行前后都设置内存`var bm runtime.MemStats`用于判断超内存，最后通过select和三种通道状态判断提交的代码的状态。
* 通过`gorm.Expr("submit_num + ?", 1)`对数据库中的数据进行累加，更新用户列表（该用户的submit_num和pass_num的改变）、问题列表（该提交的问题的submit_num和pass_num的改变）

```go
// Submit
// @Tags 用户私有方法
// @Summary 代码提交
// @Param authorization header string true "authorization"
// @Param problem_identity query string true "problem_identity"
// @Param code body string true "code"
// @Success 200 {string} json "{"code":"200","msg":"","data":""}"
// @Router /user/submit [post]
func Submit(ctx *gin.Context) {
	problemIdentity := ctx.Query("problem_identity")
	code, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "read code err:" + err.Error(),
		})
		return
	}
	// 代码保存
	path, err := helper.CodeSave(code)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "read code err:" + err.Error(),
		})
		return
	}
	// 提交
	u, _ := ctx.Get("user")
	userClaim := u.(*helper.UserClaims)
	sb := &models.SubmitBasic{
		Identity:        helper.GetUUID(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    userClaim.Identity,
		Path:            path,
	}

	// 代码判断
	pb := new(models.ProblemBasic)
	err = models.DB.Where("identity = ?", problemIdentity).Preload("TestCase").First(pb).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "read code err:" + err.Error(),
		})
		return
	}
	WA := make(chan int)  // 错误答案的情况
	OOM := make(chan int) // 超内存
	CE := make(chan int)  // 编译错误
	passCount := 0
	var lock sync.Mutex // 定义互斥锁
	// 提示信息
	var msg string

	for _, testcode := range pb.TestCase {
		go func() {
			// 通过协程执行测试
			cmd := exec.Command("go", "run", path)
			var out, stderr bytes.Buffer
			cmd.Stderr = &stderr
			cmd.Stdout = &out

			stdinPipe, err := cmd.StdinPipe()
			if err != nil {
				log.Fatalln(err)
			}
			io.WriteString(stdinPipe, testcode.Input)
			// 根据测试的输入案例进行运行拿到输出结果和标准输出结果是否匹配
			var bm runtime.MemStats
			runtime.ReadMemStats(&bm)
			if err := cmd.Run(); err != nil {
				log.Println(err, stderr.String())
				if err.Error() == "exit status 2" {
					msg = stderr.String()
					CE <- 1
					return
				}
			}
			var em runtime.MemStats
			runtime.ReadMemStats(&em)
			// 答案错误情况
			if testcode.Output != out.String() {
				msg = "答案错误"
				WA <- 1
				return
			}
			// 运行超内存情况
			if em.Alloc/1024-bm.Alloc/1024 > uint64(pb.MaxMem) {
				msg = "运行超内存"
				OOM <- 1
				return
			}
			lock.Lock()
			passCount++

			lock.Unlock()

		}()
	}

	select {
	// -1-待判断，1-正确，2-错误，3-超时，4-超内存， 5-编译错误
	case <-WA:
		sb.Status = 2
	case <-OOM:
		sb.Status = 4
	case <-CE:
		sb.Status = 5
	case <-time.After(time.Millisecond * time.Duration(pb.MaxRuntime)):
		if passCount == len(pb.TestCase) {
			sb.Status = 1
			msg = "答案正确"
		} else {
			sb.Status = 3
			msg = "运行超时"
		}
	}

	if err = models.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(sb).Error
		if err != nil {
			return errors.New("userbasic create err:" + err.Error())
		}
		m := make(map[string]interface{})
		m["submit_num"] = gorm.Expr("submit_num + ?", 1)
		if sb.Status == 1 {
			m["pass_num"] = gorm.Expr("pass_num + ?", 1)
		}
		// 更新userbasic
		err = tx.Model(new(models.UserBasic)).Where("identity = ?", userClaim.Identity).Updates(m).Error
		if err != nil {
			return errors.New("userbasic modify err:" + err.Error())
		}
		// 更新problembasic
		err = tx.Model(new(models.ProblemBasic)).Where("identity = ?", problemIdentity).Updates(m).Error
		if err != nil {
			return errors.New("problembasic modify err:" + err.Error())
		}

		return nil
	}); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "read code err:" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg": map[string]interface{}{
			"status": sb.Status,
			"msg":    msg,
		},
	})
}

```

