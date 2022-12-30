package vbody

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"
)

type Reader struct {
	M   map[string]any // 记录对象
	A   []any          // 记录数组
	m   sync.RWMutex   // 安全锁
	err error          // 错误
}

// 读取
//
//	i any	数据，支持map，array，slice，io.Reader, *string，[]byte。你也可以对 Reader.M 或 Reader.A 进行赋值。
//	*Reader			读取对象
func NewReader(i any) *Reader {
	bodyr := &Reader{M: make(map[string]any), A: make([]any, 0)}
	if i != nil {
		bodyr.err = bodyr.Reset(i)
	}
	return bodyr
}

// 错误
//
//	error	如果 Reader.Reset 重置有误，可以从这里得到相关错误。
func (T *Reader) Err() error {
	return T.err
}

// 变更
func (T *Reader) Change() *Change {
	return &Change{T}
}

func (T *Reader) isNil(key any) bool {
	T.m.RLock()
	defer T.m.RUnlock()
	switch v := key.(type) {
	case string:
		if a, ok := T.M[v]; ok {
			return reflect.ValueOf(&a).Elem().IsNil()
		}
	case int:
		if len(T.A) > v {
			return T.A[v] == nil
		}
	}
	return true
}

// 是nil值
//
//	keys ...any		键名，如果需要判断切片的长度，可以传入int类型。
//	bool			是nil值，返回true
func (T *Reader) IsNil(keys ...any) bool {
	for _, key := range keys {
		if !T.isNil(key) {
			return false
		}
	}
	return true
}

func (T *Reader) has(key any) bool {
	T.m.RLock()
	defer T.m.RUnlock()
	switch v := key.(type) {
	case string:
		if _, ok := T.M[v]; ok {
			return ok
		}
	case int:
		return len(T.A) > v
	}
	return false
}

// 检查键是否存在
//
//	keys ...any		键名，如果需要判断切片的长度，可以传入int类型。当然也可以这样 len(Reader.A)
//	bool					存在，返回true
//
//	例如：{"a":"b"} Has("a") == true 或 Has("b") == false
func (T *Reader) Has(keys ...any) bool {
	for _, key := range keys {
		if !T.has(key) {
			return false
		}
	}
	return true
}

func (T *Reader) index(i int) any {
	if i >= len(T.A) {
		return nil
	}
	return T.A[i]
}

func (T *Reader) val(key any) any {
	T.m.RLock()
	defer T.m.RUnlock()

	var v any
	switch kv := key.(type) {
	case string:
		v = T.M[kv]
	case int:
		v = T.index(kv)
	}
	return v
}

// 判断值是否等于eq
//
//	eq any				判断keys是否等于
//	keys ... any		支持多个键名判断
//	bool				值等于eq，返回true
//
//	例如：{"a":"b"}  Equal("","a") == false 或 Equal("b","a") == true
func (T *Reader) Equal(eq any, keys ...any) bool {
	for _, key := range keys {
		if reflect.DeepEqual(T.val(key), eq) {
			return true
		}
	}
	return false
}

// 读取值是字符串类型的
//
//	key, def string	键名，默认值
//	string			读取的字符串
//
//	例如：{"a":"b"} String("a","") == "b"
func (T *Reader) String(key any, def ...string) string {
	v, _ := T.val(key).(string)
	if v == "" {
		for _, f := range def {
			if f != "" {
				return f
			}
		}
	}
	return v
}

// 读取值是布尔值类型的
//
//	key any, def bool	键名，默认值
//	bool			读取的布尔值
//
//	例如：{"a":true} Bool("a",false) == true
func (T *Reader) Bool(key any, def ...bool) bool {
	v, _ := T.val(key).(bool)
	if !v {
		for _, f := range def {
			if f {
				return f
			}
		}
	}
	return v
}

// 读取值是浮点数类型的
//
//	key, def float64	键名，默认值
//	float64				读取的浮点数
//
//	例如：{"a":123} Float64（"a",123） == 123
func (T *Reader) Float64(key any, def ...float64) float64 {
	var (
		rv = reflect.ValueOf(T.val(key))
		v  float64
	)

	switch rv.Kind() {
	case reflect.Float32, reflect.Float64:
		v = rv.Float()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = float64(rv.Int())
	}

	if v == 0 {
		for _, f := range def {
			if f != 0 {
				return f
			}
		}
	}
	return v
}

// 读取值是整数类型的
//
//	key, def int64	键名，默认值
//	int64				读取的整数
//
//	例如：{"a":123} Int64("a",0) == 123 或 Int64("b",456) == 456
func (T *Reader) Int64(key any, def ...int64) int64 {
	var (
		rv = reflect.ValueOf(T.val(key))
		v  int64
	)

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = rv.Int()
	case reflect.Float32, reflect.Float64:
		v = int64(rv.Float())
	}

	if v == 0 {
		for _, f := range def {
			if f != 0 {
				return f
			}
		}
	}
	return v
}

// 读取值是接口类型的
//
//	key any			键名
//	def any		默认值
//	any			读取的接口类型，需要转换
//
//	例如：{"a":"b"} Interface("a","c") == "b" 或 Interface("b","c") == "c"
func (T *Reader) Any(key any, def ...any) any {
	v := T.val(key)
	if v == nil {
		for _, f := range def {
			if f != nil {
				return f
			}
		}
	}
	return v
}

// 读取值是接口类型的
//
//	key any			键名
//	def any		默认值
//	any			读取的接口类型，需要转换
//
//	例如：{"a":{"b":123}} NewInterface("a",*{"b":456}) == *{"b":123} 或 NewInterface("b",*{"b":456}) == *{"b":456}
func (T *Reader) NewAny(key any, def ...any) *Reader {
	return NewReader(T.Any(key, def...))
}

// 读取值是数组类型的
//
//	key any				键名
//	def []any		默认值
//	[]any			读取的数组类型
//
//	例如：{"a":[1,3,4,5,6]} Array("a",[7,8,9,0]) == [1,3,4,5,6] 或 Array("b",[7,8,9,0]) == [7,8,9,0]
func (T *Reader) Slice(key any, def ...[]any) []any {
	v, _ := T.val(key).([]any)
	if len(v) == 0 {
		for _, f := range def {
			if f != nil {
				return f
			}
		}
	}
	return v
}

// 读取值是数组类型的
//
//	key any				键名
//	def []any		默认值
//	[]any			读取的数组类型
//
//	例如：{"a":[1,3,4,5,6]} Array("a",[7,8,9,0]) == *[1,3,4,5,6] 或 Array("b",[7,8,9,0]) == *[7,8,9,0]
func (T *Reader) NewSlice(key any, def ...[]any) *Reader {
	return NewReader(T.Slice(key, def...))
}

func (T *Reader) otherLoad(i any) error {
	rv := reflect.ValueOf(i)
	tm := make(map[string]any)
	for rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch typ := rv.Kind(); typ {
	case reflect.Map:
		mr := rv.MapRange()
		for mr.Next() {
			tm[fmt.Sprint(mr.Key().Interface())] = mr.Value().Interface()
		}
		T.A = T.A[0:0]
		T.M = tm
	case reflect.Array, reflect.Slice:
		T.M = tm
		T.A = T.A[0:0]
		for i := 0; i < rv.Len(); i++ {
			T.A = append(T.A, rv.Index(i).Interface())
		}
	default:
		return errors.New("cannot convert data type " + typ.String())
	}
	return nil
}

// 重置，如果需要重置为空，需要先调用一次.Reset(nil)
//
//	i any	支持格式，包括：map,array,slice,io.Reader,*string, []byte
//	error			错误
func (T *Reader) Reset(i any) error {
	T.m.Lock()
	defer T.m.Unlock()

	var (
		err error
		tm  = make(map[string]any)
	)

	// 原类型判断
	switch iv := i.(type) {
	case io.Reader:
		err = json.NewDecoder(iv).Decode(&tm)
	case *string:
		err = json.NewDecoder(bytes.NewBufferString(*iv)).Decode(&tm)
	case string:
		err = json.NewDecoder(bytes.NewBufferString(iv)).Decode(&tm)
	case []byte:
		err = json.NewDecoder(bytes.NewBuffer(iv)).Decode(&tm)
	default:
		return T.otherLoad(i)
	}

	if err != nil {
		return err
	}

	T.M = tm
	T.A = T.A[0:0]

	return nil
}

// 从r读取字节串并解析成Reader
//
//	r io.Reader	字节串读接口
//	int64		读取长度
//	error		错误
func (T *Reader) ReadFrom(r io.Reader) (int64, error) {
	T.m.Lock()
	defer T.m.Unlock()
	err := json.NewDecoder(r).Decode(&T.M)
	if err == nil {
		T.A = T.A[0:0]
	}
	return 0, err
}

// Reader转字节串
//
//	[]byte	字节串，如：[]byte(`{"A":1}`)
//	error	错误
func (T *Reader) MarshalJSON() ([]byte, error) {
	T.m.RLock()
	defer T.m.RUnlock()
	b, err := json.Marshal(&T.M)
	if err == nil {
		T.A = T.A[0:0]
	}
	return b, err
}

// 字节串解析成Reader
//
//	data []byte	字节串，如：[]byte(`{"A":1}`)
//	error		错误
func (T *Reader) UnmarshalJSON(data []byte) error {
	T.m.Lock()
	defer T.m.Unlock()
	err := json.Unmarshal(data, &T.M)
	if err == nil {
		T.A = T.A[0:0]
	}
	return err
}
