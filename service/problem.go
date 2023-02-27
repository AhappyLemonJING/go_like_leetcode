package service

import (
	"encoding/json"
	"errors"
	"gin_gorm_oj/define"
	"gin_gorm_oj/helper"
	"gin_gorm_oj/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
