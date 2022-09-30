package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type FindItem struct {
	Key string `json:"key"`
}


type JsonFind struct {
	FindItems []FindItem
	MatchStr string
	MaxNum int
}


func (jf *JsonFind) addInfo(key string) {
	jf.FindItems = append(jf.FindItems, FindItem{
		Key:  key,
	})
}

func (jf *JsonFind) Find(sourceDst interface{}, path string)  {
	if sourceDst == nil {
		return
	}
	dstType := reflect.TypeOf(sourceDst)
	if dstType.Kind() == reflect.Slice {
		jf.findSlice(sourceDst.([]interface{}), path)
	} else if dstType.Kind() == reflect.Map {
		jf.findMap(sourceDst.(map[string]interface{}), path)
	} else {
		if jf.DoMatch(sourceDst, path) {
			return
		}
	}
}

func (jf *JsonFind) findMap(dst  map[string]interface{}, path string) {
	for key, val := range dst {
		nextKey := fmt.Sprintf("%s.%s", path, key)
		if jf.DoMatch(key, nextKey) {
			return
		}
		jf.Find(val, nextKey)
	}
}


func (jf *JsonFind) findSlice(dst  []interface{}, path string) {
	dstLen := len(dst)
	for i := 0; i < dstLen; i += 1 {
		nextKey := fmt.Sprintf("%s.%d", path, i)
		jf.Find(dst[i], nextKey)
	}
}

func (jf *JsonFind) DoMatch(dst interface{}, key string) bool {
	isMore := len(jf.FindItems) >= jf.MaxNum + 1

	if len(jf.FindItems) == jf.MaxNum {
		jf.addInfo("匹配的内容太多了，请修改关键字")
	}
	if isMore {
		return false
	}
	if jf.IsMatch(dst) {
		jf.addInfo(strings.TrimPrefix(key, "."))
	}
	return isMore
}


func (jf *JsonFind) IsMatch(str interface{}) bool {
	resultStr := jf.stringify(str)
	if len(resultStr) > 0 {
		reg := fmt.Sprintf("^.*%s.*$", jf.MatchStr)
		matched, _ := regexp.Match(reg, []byte(resultStr))
		return matched
	}
	return  false
}

func (jf *JsonFind) stringify(value interface{}) string {
	var str string
	if value == nil {
		return str
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		str = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		str = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		str = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		str = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		str = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		str = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		str = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		str = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		str = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		str = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		str = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		str = strconv.FormatUint(it, 10)
	case string:
		str = value.(string)
	case []byte:
		str = string(value.([]byte))
	default:
	}

	return str
}


func FindPath(sourceStr string, key string, maxNum int) {
	finder := JsonFind{
		MatchStr: key,
		FindItems: make([]FindItem, 0),
		MaxNum: maxNum,
	}
	var dstArrayValue []interface{}
	var dstMapValue map[string]interface{}
	if strings.HasPrefix(sourceStr, "{") {
		err := json.Unmarshal([]byte(sourceStr), &dstMapValue)
		if err != nil {
			fmt.Println("内部不是JSON格式")
			return
		}
		finder.Find(dstMapValue, "")
	}
	if strings.HasPrefix(sourceStr, "[") {
		err := json.Unmarshal([]byte(sourceStr), &dstArrayValue)
		if err != nil {
			fmt.Println("内部不是JSON格式")
			return
		}
		finder.Find(dstArrayValue, "")
	}
	for _, item := range finder.FindItems{
		var Green  = "\033[32m"
		println(Green + item.Key + "")
	}
}
