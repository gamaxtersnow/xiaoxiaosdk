package utils

import (
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// StructToURLValues 将任意结构体转换为 url.Values。
func StructToURLValues(s interface{}) url.Values {
	values := url.Values{}
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	// 确保传入的是一个结构体
	if t.Kind() != reflect.Struct {
		return values
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("json")
		tagName := tag
		if commaIdx := len(tagName) - 1; commaIdx > 0 && tagName[commaIdx] == ',' {
			tagName = tagName[:commaIdx]
		}
		tagName = strings.ToLower(tagName)
		// 处理不同类型的字段
		switch field.Kind() {
		case reflect.String:
			values.Set(tagName, field.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			values.Set(tagName, strconv.FormatInt(field.Int(), 10))
		case reflect.Float32, reflect.Float64:
			values.Set(tagName, strconv.FormatFloat(field.Float(), 'f', -1, 64))
		case reflect.Bool:
			values.Set(tagName, strconv.FormatBool(field.Bool()))
		default:
		}
	}
	return values
}

// ConvertStruct 使用反射将一个结构体转换为另一个结构体
func ConvertStruct(src, dst interface{}) error {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		dstField := dstVal.FieldByName(srcVal.Type().Field(i).Name)

		if dstField.IsValid() && dstField.CanSet() {
			if srcField.Kind() == reflect.Ptr {
				// Handle pointer fields
				if !srcField.IsNil() {
					srcElem := srcField.Elem()
					dstElem := reflect.New(srcElem.Type()).Elem()
					dstElem.Set(srcElem)
					dstField.Set(dstElem.Addr())
				} else {
					dstField.Set(reflect.Zero(dstField.Type()))
				}
			} else if srcField.Kind() == reflect.Struct {
				// 处理结构体字段
				if err := ConvertStruct(srcField.Interface(), dstField.Addr().Interface()); err != nil {
					return err
				}
			} else if srcField.Kind() == reflect.Slice && dstField.Kind() == reflect.Slice {
				// 处理切片
				if err := convertSlice(srcField, dstField); err != nil {
					return err
				}
			} else {
				dstField.Set(srcField)
			}
		}
	}
	return nil
}

func convertSlice(srcField, dstField reflect.Value) error {
	dstSlice := reflect.MakeSlice(dstField.Type(), srcField.Len(), srcField.Cap())
	for j := 0; j < srcField.Len(); j++ {
		srcElem := srcField.Index(j)
		dstElem := reflect.New(dstField.Type().Elem()).Elem()
		for k := 0; k < srcElem.NumField(); k++ {
			srcField := srcElem.Field(k)
			dstField := dstElem.FieldByName(srcElem.Type().Field(k).Name)

			if dstField.IsValid() && dstField.CanSet() {
				dstField.Set(srcField)
			}
		}

		dstSlice.Index(j).Set(dstElem)
	}
	dstField.Set(dstSlice)
	return nil
}

// 获取格式化的时间
func GetTimeText(timestamp int64) (dateStr, timeStr string) {
	// 将毫秒级时间戳转换为time.Time
	t := time.UnixMilli(timestamp)

	// 中文周几
	weeks := []string{"日", "一", "二", "三", "四", "五", "六"}
	weekday := weeks[t.Weekday()]

	// 格式化日期：2024年7月1日[周五]
	dateStr = fmt.Sprintf("%d年%d月%d日[周%s]",
		t.Year(),
		t.Month(),
		t.Day(),
		weekday,
	)

	// 格式化时间：10:00
	timeStr = t.Format("15:04")

	return dateStr, timeStr
}
func GetCourseTypeText(CourseType int64) string {
	switch CourseType {
	case 1:
		return "主教课程"
	case 2:
		return "助教课程"
	default:
		return "未知课程"
	}
}

// 获取上课方式文本
func GetWayText(way int64) string {
	switch way {
	case 1:
		return "线下"
	case 2:
		return "线上"
	default:
		return "未知"
	}
}

// 获取课程状态文本
func GetStatusText(status int64) string {
	switch status {
	case 1:
		return "已上课"
	case 2:
		return "已取消"
	default:
		return "未上课"
	}
}
func GetAttendanceText(countJoin int64, countClz int64) string {
	return strconv.Itoa(int(countJoin)) + "/" + strconv.Itoa(int(countClz))
}
func GetUUID() string {
	newUUID := uuid.New()
	return strings.ReplaceAll(newUUID.String(), "-", "")
}
