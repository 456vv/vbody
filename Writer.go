package vbody

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"text/template"
	"sync"
)

type Writer struct {
	M map[string]any // 写记录
	m sync.Mutex
}

// 写，这是http响应的一种特定格式。采用格式为：{"Code":0,"Message":"内容","Result":[{...},{...}]}
func NewWriter() *Writer {
	return &Writer{M: make(map[string]any)}
}

// 状态
//
//	d int	状态码
func (T *Writer) Code(d int) *Writer {
	T.m.Lock()
	defer T.m.Unlock()
	
	T.M["Code"] = d
	return T
}

// 提示内容
//
//	s any	内容
func (T *Writer) Message(s any) *Writer {
	T.m.Lock()
	defer T.m.Unlock()
	
	T.M["Message"] = template.JSEscaper(s)
	return T
}

// 提示内容，支持fmt.Sprintf 格式
//
//	f string			格式
//	a ...any	参数
func (T *Writer) Messagef(f string, a ...any) *Writer {
	return T.Message(fmt.Sprintf(f, a...))
}

// 设置结果
//
//	val any	结果
func (T *Writer) SetResult(val any) *Writer {
	T.m.Lock()
	defer T.m.Unlock()
	
	T.M["Result"] = val
	return T
}

// 设置结果
//
//	key string		键名
//	val any	值
func (T *Writer) Result(key string, val any) *Writer {
	T.m.Lock()
	defer T.m.Unlock()
	
	if t, ok := val.([]byte); ok {
		val = json.RawMessage(t)
	}
	v := reflect.ValueOf(T.M["Result"])
	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
	}
	switch v.Kind() {
	case reflect.Map:
		v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	default:
		T.M["Result"] = map[string]any{
			key: val,
		}
	}
	return T
}

// 写入到w
//
//	w io.Writer	写入接口
//	n int64		写入长度
//	err error	错误
func (T *Writer) WriteTo(w io.Writer) (n int64, err error) {
	if err := T.ready(); err != nil {
		return 0, err
	}

	buf := bytes.NewBuffer(nil)
	buf.Grow(1024)
	err = json.NewEncoder(buf).Encode(T.M)
	if err != nil {
		return 0, err
	}

	if rw, ok := w.(http.ResponseWriter); ok {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	return buf.WriteTo(w)
}

func (T *Writer) ready() error {
	T.m.Lock()
	defer T.m.Unlock()
	
	if _, ok := T.M["Code"]; !ok {
		T.M["Code"] = 0
	}
	resultInf, ok := T.M["Result"]
	if !ok {
		T.M["Result"] = nil
	} else if result, ok := resultInf.(json.Marshaler); ok {
		rbyte, err := result.MarshalJSON()
		if err != nil {
			return err
		}
		T.M["Result"] = json.RawMessage(rbyte)
	} else if rbyte, ok := resultInf.([]byte); ok {
		T.M["Result"] = json.RawMessage(rbyte)
	}
	return nil
}

// 字符串json格式
//
//	string		json格式
func (T *Writer) String() string {
	if err := T.ready(); err != nil {
		return err.Error()
	}
	b, err := json.Marshal(T.M)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
