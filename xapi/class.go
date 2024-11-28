package xapi

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk/utils"
	"context"
	"encoding/json"
	"io"
)

type ClassListByRoleReq struct {
	Search    string `json:"search"`    // 搜索关键字（可能为空）
	CourseID  string `json:"courseid"`  // 课程 ID（可能为空）
	CampusIDs string `json:"campusids"` // 校区 ID 列表
	PageID    int    `json:"pageid"`    // 当前页码
	PageSize  int    `json:"pagesize"`  // 每页大小
	Status    string `json:"status"`    // 状态
}

type ClassListByRoleResp struct {
	Duration      int                  `json:"duration"`      // 持续时间
	TheIngAccROLE int                  `json:"theIngAccROLE"` // 当前账户角色
	Pages         ClassListByRolePages `json:"pages"`         // 分页信息
	IpsActive     int                  `json:"ipsActive"`     // 活跃 IP 数
	Data          []ClassCourse        `json:"data"`          // 课程数据列表
	Success       bool                 `json:"success"`       // 请求是否成功
	Messages      string               `json:"messages"`      // 消息（可能为空）
	Timestamp     int64                `json:"timestamp"`     // 时间戳
}

type ClassListByRolePages struct {
	Total    int `json:"total"`    // 总记录数
	PageSize int `json:"pageSize"` // 每页大小
	PageID   int `json:"pageId"`   // 当前页码
	Pages    int `json:"pages"`    // 总页数
}

type ClassCourse struct {
	ID         int    `json:"id"`         // 课程 ID
	Name       string `json:"name"`       // 课程名称
	CampusID   int    `json:"campusId"`   // 校区 ID
	CampusName string `json:"campusName"` // 校区名称
	CourseID   int    `json:"courseId"`   // 课程 ID
	CourseName string `json:"courseName"` // 课程名称
	ClassType  int    `json:"classtype"`  // 课程类型（可能为空）
	Category   string `json:"category"`   // 课程类别
	Status     int    `json:"status"`     // 课程状态
}

var _ ClassModel = (*customClassModel)(nil)

type (
	ClassModel interface {
		GetClassListByRole(ctx context.Context, req ClassListByRoleReq) (*ClassListByRoleResp, error)
	}
	customClassModel struct {
		client *xiaoxiaosdk.HttpClient
	}
)

func NewClassModel(client *xiaoxiaosdk.HttpClient) ClassModel {
	return &customClassModel{
		client: client,
	}
}

// GetClassListByRole 通过角色获取班级列表
func (t *customClassModel) GetClassListByRole(ctx context.Context, req ClassListByRoleReq) (*ClassListByRoleResp, error) {
	params := utils.StructToURLValues(req)
	resp, err := t.client.Get("/clz/list/page/by/role", params)
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
	classListByRoleResp := &ClassListByRoleResp{}
	err = json.Unmarshal(body, classListByRoleResp)
	if err != nil {
		return nil, err
	}
	return classListByRoleResp, nil
}
