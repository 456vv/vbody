package vbody

import(
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"net/http"
	"bytes"
)

type Writer struct{
	M map[string]interface{}						// 写记录
}

//写，这是http响应的一种特定格式。采用格式为：{"Status":200,"Message":"内容","Result":[{...},{...}]}
func NewWriter() *Writer {
	return &Writer{M:make(map[string]interface{})}
}

//状态
//	d int	状态码
func (T *Writer) Status(d int) {
	T.M["Status"]=d
}

//提示内容
//	s interface{}	内容
func (T *Writer) Message(s interface{}) {
	T.M["Message"]=s
}

//提示内容，支持fmt.Sprintf 格式
//	f string			格式
//	a ...interface{}	参数
func (T *Writer) Messagef(f string, a ...interface{}) {
	T.M["Message"]=fmt.Sprintf(f, a...)
}

//设置结果
//	i interface{}	结果
func (T *Writer) SetResult(i interface{}) {
	T.M["Result"]=i
}

//设置结果
//	key string		键名
//	i interface{}	值
func (T *Writer) Result(key string, i interface{}) {
	v := reflect.ValueOf(T.M["Result"])
	for ;v.Kind() == reflect.Ptr; v = v.Elem() {}
	switch v.Kind() {
	case reflect.Map:
		v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(i))
	default:
		m := make(map[string]interface{})
		m[key]=i
		T.M["Result"]=m
	}
}

//写入到w
//	w io.Writer	写入接口
//	n int64		写入长度
//	err error	错误
func (T *Writer) WriteTo(w io.Writer) (n int64, err error) {
	_, ok := T.M["Status"]
	if !ok {
		T.M["Status"]=200
	}
	resultInf, ok := T.M["Result"]
	if !ok {
		T.M["Result"]=nil
	}else if result, ok := resultInf.(json.Marshaler); ok {
		rbyte, err := result.MarshalJSON()
		if err != nil {
			return 0, err
		}
		T.M["Result"]=json.RawMessage(rbyte)
	}else if rbyte, ok := resultInf.([]byte); ok {
		T.M["Result"]=json.RawMessage(rbyte)
	}
	if rw, ok := w.(http.ResponseWriter); ok {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	buf := bytes.NewBuffer(nil)
	buf.Grow(1024)
	err = json.NewEncoder(buf).Encode(T.M)
	if err != nil {
		return 0, err
	}
	
	return buf.WriteTo(w)
}
