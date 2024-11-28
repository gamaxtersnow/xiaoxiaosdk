package xapi

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk"
	"context"
	"encoding/json"
	"io"
	"net/url"
)

type CampusListResponse struct {
	Duration   int64        `json:"duration"`
	IpsActive  int64        `json:"ips_active"`
	Messages   string       `json:"messages"`
	Success    bool         `json:"success"`
	TimeStamp  int64        `json:"timestamp"`
	CampusList []CampusList `json:"campusList"`
}

type CampusList struct {
	Id         int64  `json:"id"`
	OrgId      int64  `json:"orgid"`
	Type       int64  `json:"type"`
	Name       string `json:"name"`
	CreateTime int64  `json:"createtime"`
}

type CampusRoomListResponse struct {
	Duration      int      `json:"duration"`      // 事件或过程的持续时间
	TheIngAccROLE int      `json:"theIngAccROLE"` // 账户的角色
	IpsActive     int      `json:"ipsActive"`     // 活跃 IP 的数量
	Data          []Campus `json:"data"`          // 校区列表
	Success       bool     `json:"success"`       // 操作是否成功
	Messages      string   `json:"messages"`      // 可选的消息，可能为 null
	Timestamp     int64    `json:"timestamp"`     // 数据的时间戳
}

type Campus struct {
	ID             int         `json:"id"`             // 校区的唯一标识符
	Name           string      `json:"name"`           // 校区名称
	OutClzRoomList []Classroom `json:"outClzRoomList"` // 校区中的教室列表
}

type Classroom struct {
	ID   int    `json:"id"`   // 教室的唯一标识符
	Name string `json:"name"` // 教室名称
}

var _ CampusModel = (*ConCampusModel)(nil)

type (
	CampusModel interface {
		GetAllCampuses(ctx context.Context) (*CampusListResponse, error)
		GetCampusRoomList(ctx context.Context) (*CampusRoomListResponse, error)
	}
	ConCampusModel struct {
		client *xiaoxiaosdk.HttpClient
	}
)

func NewCampusModel(client *xiaoxiaosdk.HttpClient) *ConCampusModel {
	return &ConCampusModel{
		client: client,
	}
}

// GetAllCampuses 获取校区列表
func (c *ConCampusModel) GetAllCampuses(ctx context.Context) (*CampusListResponse, error) {
	resp, err := c.client.Get("/getCampusList", url.Values{"dlall": []string{"1"}})
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
	campusListResp := &CampusListResponse{}
	err = json.Unmarshal(body, campusListResp)
	if err != nil {
		return nil, err
	}
	return campusListResp, nil
}

// GetCampusRoomList 获取校区教室列表
func (m *ConCampusModel) GetCampusRoomList(ctx context.Context) (*CampusRoomListResponse, error) {
	resp, err := m.client.Get("/campus/clz/room/list", nil)
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
	CampusRoomListResponse := &CampusRoomListResponse{}
	err = json.Unmarshal(body, CampusRoomListResponse)
	if err != nil {
		return nil, err
	}
	return CampusRoomListResponse, nil
}
