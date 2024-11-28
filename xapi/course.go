package xapi

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk"
	"context"
	"encoding/json"
	"io"
	"net/url"
)

type CourseListResp struct {
	Duration      int      `json:"duration"`      // 响应的持续时间
	Total         int      `json:"total"`         // 项目总数
	TheIngAccROLE int      `json:"theIngAccROLE"` // 账户角色
	IpsActive     int      `json:"ipsActive"`     // 活跃 IP 数量
	Data          []Course `json:"data"`          // 课程列表
	Success       bool     `json:"success"`       // 请求是否成功
	PageSize      int      `json:"pagesize"`      // 页面大小
	Messages      string   `json:"messages"`      // 消息（如果有）
	PageID        int      `json:"pageid"`        // 页面 ID
	Timestamp     int64    `json:"timestamp"`     // 响应的时间戳
}

type Course struct {
	ID                  int         `json:"id"`                  // 课程 ID
	OrganizationID      int         `json:"organizationid"`      // 组织 ID
	Name                string      `json:"name"`                // 课程名称
	Classification      string      `json:"classification"`      // 课程分类
	Category            int         `json:"category"`            // 课程类别
	SwPrice             bool        `json:"swprice"`             // 价格开关标志
	FeeForm             string      `json:"feeform"`             // 费用形式
	Period              string      `json:"period"`              // 期间
	PriceRaw            float64     `json:"priceraw"`            // 原始价格
	Price               float64     `json:"price"`               // 价格
	Description         string      `json:"description"`         // 描述
	Pics                string      `json:"pics"`                // 图片
	ContactWeixin       string      `json:"contactweixin"`       // 联系微信
	ContactWeixinV2Code string      `json:"contactweixinv2code"` // 联系微信 V2 代码
	CtContents          []CtContent `json:"ctcontents"`          // 课程内容
	IPkg                bool        `json:"ipkg"`                // IPkg 标志
	PkgItems            []string    `json:"pkgitems"`            // 包含项目
	Kind                int         `json:"kind"`                // 类型
	PricePerHour        float64     `json:"priceperhour"`        // 每小时价格
	JoinedClz           string      `json:"joinedClz"`           // 加入的班级
	JoinedClzList       string      `json:"joinedClzList"`       // 加入的班级列表
	PkgIncJson          string      `json:"pkgincjson"`          // 包含 JSON
	CotPkgCourse        int         `json:"cotPkgCourse"`        // 包课程数量
	CotPkgSundry        int         `json:"cotPkgSundry"`        // 包杂项数量
	IngForAdd           bool        `json:"ingForAdd"`           // 添加标志
}

type CtContent struct {
	ID             int    `json:"id"`             // 内容 ID
	OrganizationID int    `json:"organizationid"` // 组织 ID
	CourseID       int    `json:"courseid"`       // 课程 ID
	Content        string `json:"content"`        // 内容描述
}

var _ CourseModel = (*customCourseModel)(nil)

type (
	CourseModel interface {
		GetAllCourses(ctx context.Context) (*CourseListResp, error)
	}
	customCourseModel struct {
		client *xiaoxiaosdk.HttpClient
	}
)

func NewCourseModel(client *xiaoxiaosdk.HttpClient) CourseModel {
	return &customCourseModel{
		client: client,
	}
}

// GetAllCourses 获取所有课程
func (c *customCourseModel) GetAllCourses(ctx context.Context) (*CourseListResp, error) {
	resp, err := c.client.Get("/course/list", url.Values{"dlall": []string{"1"}})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	courseListResp := &CourseListResp{}
	err = json.Unmarshal(body, courseListResp)
	if err != nil {
		return nil, err
	}
	return courseListResp, nil
}
