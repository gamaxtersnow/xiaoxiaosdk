package xapi

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk/utils"
	"context"
	"encoding/json"
	"io"
)

type ScheduleViewReq struct {
	ViewType   int64  `json:"viewtype"`   // 视图类型
	Search     string `json:"search"`     // 搜索关键字
	CourseIDs  string `json:"courseids"`  // 课程 ID 列表
	CampusIDs  string `json:"campusids"`  // 校区 ID 列表
	TeacherIDs string `json:"teacherids"` // 教师 ID 列表
	ZhuJiaoID  string `json:"zhujiaoid"`  // 助教 ID
	ClassIDs   string `json:"classids"`   // 班级 ID 列表
	Status     string `json:"status"`     // 状态
	PageID     int    `json:"pageid"`     // 当前页码
	PageSize   int    `json:"pagesize"`   // 每页大小
	ClzRooms   string `json:"clzrooms"`   // 教室列表
	ShowRange  int64  `json:"showrange"`  // 显示范围
	ExceptNull int64  `json:"exceptnull"` // 排除空值
	ShowCancel int64  `json:"showcancel"` // 显示取消
	Asc        int64  `json:"asc"`        // 升序或降序
	Dfrom      string `json:"dfrom"`      //日期起始时间
	Dto        string `json:"dto"`        //日期结束时间
	Dlall      int    `json:"dlall"`      // 是否全部
}

type ScheduleViewResp struct {
	Duration      int64      `json:"duration"`      // 持续时间
	Total         int        `json:"total"`         // 总数
	TheIngAccROLE int64      `json:"theIngAccROLE"` // 账户角色
	IpsActive     int64      `json:"ipsActive"`     // 活跃 IP 数量
	Data          []Schedule `json:"data"`          // 课程安排列表
	Success       bool       `json:"success"`       // 操作是否成功
	PageSize      int        `json:"pagesize"`      // 每页大小
	Messages      string     `json:"messages"`      // 消息，可能为 null
	PageID        int        `json:"pageid"`        // 页码
	ConflictCount int        `json:"conflictCount"` // 冲突计数
	Timestamp     int64      `json:"timestamp"`     // 时间戳
}

type Schedule struct {
	ID                    int64         `json:"id"`                    // 课程安排 ID
	LockVersion           int64         `json:"lockversion"`           // 锁版本
	ClzID                 int64         `json:"clzid"`                 // 班级 ID
	OrganizationID        int64         `json:"organizationid"`        // 组织 ID
	CourseID              int64         `json:"courseid"`              // 课程 ID
	Teacher               Person        `json:"teacher"`               // 教师信息
	StartTime             int64         `json:"starttime"`             // 开始时间（毫秒）
	EndTime               int64         `json:"endtime"`               // 结束时间（毫秒）
	Way                   int64         `json:"way"`                   // 方式
	Content               string        `json:"content"`               // 内容
	Place                 string        `json:"place"`                 // 地点
	URLTeacher            interface{}   `json:"urlteacher"`            // 教师 URL
	URLStudent            interface{}   `json:"urlstudent"`            // 学生 URL
	Note                  interface{}   `json:"note"`                  // 备注
	Duration              string        `json:"duration"`              // 持续时间（格式化）
	DurationMinutes       int64         `json:"durationminutes"`       // 持续时间（分钟）
	Knowledge             interface{}   `json:"knowledge"`             // 知识点
	Training              interface{}   `json:"training"`              // 培训
	Scene                 interface{}   `json:"scene"`                 // 场景
	Genseeid              interface{}   `json:"genseeid"`              // 生成 ID
	FromSmart             bool          `json:"fromsmart"`             // 是否来自智能
	TDay                  int64         `json:"tday"`                  // 天
	Category              int64         `json:"category"`              // 类别
	CourseType            int64         `json:"coursetype"`            // 课程类型
	TagContent            string        `json:"tagcontent"`            // 标签内容
	TagCIndex             int64         `json:"tagcindex"`             // 标签索引
	Adjusts               []interface{} `json:"adjusts"`               // 调整 未获取具体类型暂时设置为interface
	Ended                 bool          `json:"ended"`                 // 是否结束
	LiveRCDs              interface{}   `json:"livercds"`              // 直播记录
	IsConflict            bool          `json:"isconflict"`            // 是否冲突
	PreoccupyConflict     bool          `json:"preoccupyconflict"`     // 预占冲突
	StudentConflict       bool          `json:"studentconflict"`       // 学生冲突
	IgnorePerfms          interface{}   `json:"ignoreperfms"`          // 忽略性能
	CorpWxScheduleID      interface{}   `json:"corpwxScheduleId"`      // 企业微信日程 ID
	IFromSmart            bool          `json:"ifromsmart"`            // 是否来自智能
	Status                int64         `json:"status"`                // 状态
	PlaceID               interface{}   `json:"placeId"`               // 地点 ID
	EdStatus              bool          `json:"edstatus"`              // 编辑状态
	Attendanced           bool          `json:"attendanced"`           // 是否考勤
	Attendance            int64         `json:"attendance"`            // 考勤
	SCountClz             int64         `json:"scountclz"`             // 班级计数
	SCountJoin            int64         `json:"scountjoin"`            // 加入计数
	SCountLeave           int64         `json:"scountleave"`           // 离开计数
	StartTimeStr          string        `json:"starttimeStr"`          // 开始时间（字符串）
	EndTimeStr            string        `json:"endtimeStr"`            // 结束时间（字符串）
	CourseName            string        `json:"courseName"`            // 课程名称
	CourseCategory        string        `json:"courseCategory"`        // 课程类别
	ClzName               string        `json:"clzName"`               // 班级名称
	ScheduleCount         int64         `json:"scheduleCount"`         // 日程计数
	PerformanceInSchedule interface{}   `json:"performanceInSchedule"` // 日程中的表现
	URLAssistant          interface{}   `json:"urlassistant"`          // 助教 URL
	IngStatus             int64         `json:"ingstatus"`             // 进行状态
	Adjust                interface{}   `json:"adjust"`                // 调整
	ClassType             int64         `json:"classtype"`             // 班级类型
	Force                 interface{}   `json:"force"`                 // 强制
	CampusID              int64         `json:"campusid"`              // 校区 ID
	CampusName            interface{}   `json:"campusname"`            // 校区名称
	Conficts              interface{}   `json:"conficts"`              // 冲突
	TheSmart              bool          `json:"thesmart"`              // 智能
	ClzZhujiao            []Person      `json:"clzZhujiao"`            // 班级助教
	ViewSrc               interface{}   `json:"viewsrc"`               // 视图来源
	Diseditable           bool          `json:"diseditable"`           // 是否可编辑
	PeriodOfClz           float64       `json:"periodOfClz"`           // 班级周期
	PeriodOfCur           float64       `json:"periodOfCur"`           // 当前周期
	ClassPlatform         int64         `json:"classplatform"`         // 班级平台
	DayList               interface{}   `json:"dayList"`               // 天列表
	ScheduleType          interface{}   `json:"scheduleType"`          // 日程类型
	TeacherID             interface{}   `json:"teacherid"`             // 教师 ID
}
type Person struct {
	ID             int64       `json:"id"`             // 人员 ID
	OrganizationID int64       `json:"organizationid"` // 组织 ID
	Name           string      `json:"name"`           // 姓名
	Mobile         string      `json:"mobile"`         // 手机号
	Role           int64       `json:"role"`           // 角色
	Teaching       interface{} `json:"teaching"`       // 教学
	PhotoURL       interface{} `json:"photourl"`       // 照片 URL
	Wechat         interface{} `json:"wechat"`         // 微信
	Gender         interface{} `json:"gender"`         // 性别
	Introduce      interface{} `json:"introduce"`      // 介绍
	ClassInUID     int64       `json:"classinuid"`     // ClassIn 用户 ID
	Token          interface{} `json:"token"`          // 令牌
	Gzhvc          interface{} `json:"gzhvc"`          // 公众号验证
	FeedBack       interface{} `json:"feedBack"`       // 反馈
	AccPriosList   interface{} `json:"accPriosList"`   // 账户优先级列表
}

var _ ScheduleModel = (*customScheduleModel)(nil)

type (
	ScheduleModel interface {
		GetScheduleList(ctx context.Context, req *ScheduleViewReq) (*ScheduleViewResp, error)
	}
	customScheduleModel struct {
		client *xiaoxiaosdk.HttpClient
	}
)

func NewScheduleModel(client *xiaoxiaosdk.HttpClient) ScheduleModel {
	return &customScheduleModel{
		client: client,
	}
}

// GetScheduleList 获取排课列表
func (s *customScheduleModel) GetScheduleList(ctx context.Context, req *ScheduleViewReq) (*ScheduleViewResp, error) {
	param := ScheduleViewReq{
		ViewType:   req.ViewType,
		ExceptNull: req.ExceptNull,
		ShowCancel: req.ShowCancel,
		PageID:     req.PageID,
		PageSize:   req.PageSize,
		Asc:        req.Asc,
		Search:     req.Search,
		ShowRange:  req.ShowRange,
		Status:     req.Status,
		ClassIDs:   req.ClassIDs,
		CampusIDs:  req.CampusIDs,
		ClzRooms:   req.ClzRooms,
		CourseIDs:  req.CourseIDs,
		ZhuJiaoID:  req.ZhuJiaoID,
		TeacherIDs: req.TeacherIDs,
		Dfrom:      req.Dfrom,
		Dto:        req.Dto,
		Dlall:      req.Dlall,
	}
	resParams := utils.StructToURLValues(param)
	resp, err := s.client.Get("/schedule/views", resParams)
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
	scheduleViewResp := &ScheduleViewResp{}
	err = json.Unmarshal(body, scheduleViewResp)
	if err != nil {
		return nil, err
	}
	return scheduleViewResp, nil
}
