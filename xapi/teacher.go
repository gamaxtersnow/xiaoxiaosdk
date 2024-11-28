package xapi

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk/utils"
	"context"
	"encoding/json"
	"io"
)

type TeacherListReq struct {
	Role     int `json:"role"`     // 角色标识符
	PageSize int `json:"pagesize"` // 每页返回的记录数
	PageID   int `json:"pageid"`   // 当前页码
}

type TeacherListResp struct {
	Duration      int       `json:"duration"`      // 持续时间
	Total         int       `json:"total"`         // 总数
	TheIngAccROLE int       `json:"theIngAccROLE"` // 当前账户角色
	IpsActive     int       `json:"ipsActive"`     // 活跃 IP 数
	Data          []Teacher `json:"data"`          // 用户数据列表
	Success       bool      `json:"success"`       // 请求是否成功
	PageSize      int       `json:"pagesize"`      // 每页大小
	Messages      string    `json:"messages"`      // 消息（可能为空）
	PageID        int       `json:"pageid"`        // 当前页码
	Timestamp     int64     `json:"timestamp"`     // 时间戳
}

type Teacher struct {
	ID             int    `json:"id"`             // 用户 ID
	OrganizationID int    `json:"organizationid"` // 组织 ID
	Password       string `json:"password"`       // 密码
	Name           string `json:"name"`           // 姓名
	Role           int    `json:"role"`           // 角色
	Mobile         string `json:"mobile"`         // 手机号
	Teaching       string `json:"teaching"`       // 教学信息（可能为空）
	Username       string `json:"username"`       // 用户名
	SysUse         bool   `json:"sysuse"`         // 系统使用标志
	PhotoURL       string `json:"photourl"`       // 照片 URL（可能为空）
	Wechat         string `json:"wechat"`         // 微信（可能为空）
	Gender         string `json:"gender"`         // 性别（可能为空）
	Introduce      string `json:"introduce"`      // 介绍（可能为空）
	TeachLength    string `json:"teachlength"`    // 教学时长（可能为空）
	TeachStyle     string `json:"teachstyle"`     // 教学风格（可能为空）
	City           string `json:"city"`           // 城市（可能为空）
	PhotoPro       string `json:"photopro"`       // 照片专业（可能为空）
	VideoPro       string `json:"videopro"`       // 视频专业（可能为空）
	ShowYangyu     bool   `json:"showYangyu"`     // 是否显示杨语
	ShowExpYangyu  bool   `json:"showExpYangyu"`  // 是否显示杨语经验
	ClassInUID     int    `json:"classinuid"`     // ClassIn 用户 ID
	BDel           bool   `json:"bdel"`           // 删除标志
	RAccountID     int    `json:"raccountid"`     // 账户 ID（可能为空）
	TeaDuration    string `json:"teaduration"`    // 教学持续时间
	Token          string `json:"token"`          // 令牌（可能为空）
	Gzhvc          string `json:"gzhvc"`          // 公众号验证（可能为空）
	AccPriosList   string `json:"accPriosList"`   // 账户优先级列表（可能为空）
	Disabled       bool   `json:"disabled"`       // 是否禁用
	JustRecovery   bool   `json:"justRecovery"`   // 是否刚恢复
}

var _ TeacherModel = (*customTeacherModel)(nil)

type (
	TeacherModel interface {
		GetTeacherList(ctx context.Context, reg TeacherListReq) (*TeacherListResp, error)
	}
	customTeacherModel struct {
		client *xiaoxiaosdk.HttpClient
	}
)

func NewTeacherModel(client *xiaoxiaosdk.HttpClient) TeacherModel {
	return &customTeacherModel{
		client: client,
	}
}

// GetTeacherList 获取老师列表
func (t *customTeacherModel) GetTeacherList(ctx context.Context, req TeacherListReq) (*TeacherListResp, error) {
	resp, err := t.client.Get("/teacher/list", utils.StructToURLValues(req))
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
	teacherListResp := &TeacherListResp{}
	err = json.Unmarshal(body, teacherListResp)
	if err != nil {
		return nil, err
	}
	return teacherListResp, nil
}
