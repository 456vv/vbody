package vbody
	
import(
	"encoding/json"
	"reflect"
	"errors"
	"io"
	"bytes"
	"sync"
)
	
type Reader struct{
	M	map[string]interface{}				// 记录对象
	A	[]interface{}						// 记录数组
	m	sync.RWMutex						// 安全锁
	noZero bool								// 不可以为零
	err	error								// 错误
}

//读取
//	i interface{}	数据，支持map，array，slice，io.Reader, *string，[]byte。你也可以对 Reader.M 或 Reader.A 进行赋值。	
//	*Reader			读取对象
func NewReader(i interface{}) *Reader {
	bodyr := &Reader{M:make(map[string]interface{}), A:make([]interface{},0)}
	if i != nil {
		bodyr.err = bodyr.Reset(i)
	}
	return bodyr
}

//值不可以为零值，如果为零值，则读取是def设置值。
// y bool	true不可以为零，false可以为零值
func (T *Reader) NoZero(y bool) *Reader {
	t := new(Reader)
	*t = *T
	t.noZero=y
	return t
}

//错误
//	error	如果 Reader.Reset 重置有误，可以从这里得到相关错误。
func (T *Reader) Err() error {
	return T.err
}

//变更
func (T *Reader) Change() *Change {
	return &Change{T}
}

func (T *Reader) isNil(key interface{}) bool {
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

//是nil值
//	keys ...interface{}		键名，如果需要判断切片的长度，可以传入int类型。
//	bool					是nil值，返回true
func (T *Reader) IsNil(keys ... interface{}) bool {
	for _, key := range keys {
		if !T.isNil(key) {
			return false
		}
	}
	return true
}


func (T *Reader) has(key interface{}) bool {
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

//检查键是否存在
//	keys ...interface{}		键名，如果需要判断切片的长度，可以传入int类型。当然也可以这样 len(Reader.A)
//	bool					存在，返回true
//	
//	例如：{"a":"b"} Has("a") == true 或 Has("b") == false
func (T *Reader) Has(keys ... interface{}) bool {
	for _, key := range keys {
		if !T.has(key) {
			return false
		}
	}
	return true
}

//读取值是字符串类型的
//	key, def string	键名，默认值
//	string			读取的字符串
//	
//	例如：{"a":"b"} String("a","") == "b"
func (T *Reader) String(key string, def ...string) string {
	T.m.RLock()
	defer T.m.RUnlock()
	v, ok := T.M[key].(string)
	if !ok || T.noZero && v=="" {
		for _, f := range def {
			if f != "" {
				return f
			}
		}
	}
	return v
}

//判断值是否等于eq这个字符串
//	eq string			判断keys是否等于这个字符串
//	keys ... string		支持多个键名判断
//	bool				值等于eq，返回true
//	
//	例如：{"a":"b"}  StringAnyEqual("","a") == false 或 StringAnyEqual("b","a") == true
func (T *Reader) StringAnyEqual(eq string, keys ... string) bool {
	for _, key := range keys {
		if T.String(key) == eq {
			return true
		}
	}
	return false
}

//读取值是布尔值类型的
//	key string, def bool	键名，默认值
//	bool			读取的布尔值
//	
//	例如：{"a":true} Bool("a",false) == true
func (T *Reader) Bool(key string, def ... bool) bool {
	T.m.RLock()
	defer T.m.RUnlock()
	v, ok := T.M[key].(bool)
	if !ok {
		for _, f := range def {
			if f != false {
				return f
			}
		}
		return false
	}
	return v
}

//判断值是否等于eq这个布尔值
//	eq bool				判断keys是否等于这个布尔值
//	keys ... string		支持多个键名判断
//	bool				值等于eq，返回true
//	
//	例如：{"a":true}  BoolAnyEqual(false,"a") == false 或 BoolAnyEqual(true,"a") == true
func (T *Reader) BoolAnyEqual(eq bool, keys ... string) bool {
	for _, key := range keys {
		if T.Bool(key) == eq {
			return true
		}
	}
	return false
}

//读取值是浮点数类型的
//	key, def float64	键名，默认值
//	float64				读取的浮点数
//	
//	例如：{"a":123} Float64（"a",123） == 123
func (T *Reader) Float64(key string, def ... float64) float64 {
	T.m.RLock()
	defer T.m.RUnlock()
	rv := reflect.ValueOf(T.M[key])
	var (
		f64 float64
		notExists bool
	)
	switch rv.Kind() {
	case reflect.Float32,reflect.Float64:
		f64 = rv.Float()
	case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
		f64 = float64(rv.Int())
	default:
		notExists = true
	}
	//1，强制不为0
	//2，不存在该字段
	if T.noZero && f64 == 0 || notExists {
		for _, f := range def {
			if f != 0 {
				return f
			}
		}
	}
	return f64
}

//判断值是否等于eq这个浮点数
//	eq string			判断keys是否等于这个浮点数
//	keys ... string		支持多个键名判断
//	bool				值等于eq，返回true
//	
//	例如：{"a":123}  Float64AnyEqual(456,"a") == false 或 Float64AnyEqual(123,"a") == true
func (T *Reader) Float64AnyEqual(eq float64, keys ... string) bool {
	for _, key := range keys {
		if T.Float64(key) == eq {
			return true
		}
	}
	return false
}

//读取值是整数类型的
//	key, def int64	键名，默认值
//	int64				读取的整数
//	
//	例如：{"a":123} Int64("a",0) == 123 或 Int64("b",456) == 456
func (T *Reader) Int64(key string, def ... int64) int64 {
	T.m.RLock()
	defer T.m.RUnlock()
	rv := reflect.ValueOf(T.M[key])
	var (
		i64 int64
		notExists bool
	)
	switch rv.Kind() {
	case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
		i64 = rv.Int()
	case reflect.Float32,reflect.Float64:
		i64 = int64(rv.Float())
	default:
		notExists = true
	}
	
	//1，强制不为0
	//2，不存在该字段
	if T.noZero && i64 == 0 || notExists {
		for _, f := range def {
			if f != 0 {
				return f
			}
		}
	}
	return i64
}

//判断值是否等于eq这个整数
//	eq int64			判断keys是否等于这个整数
//	keys ... int64		支持多个键名判断
//	bool				值等于eq，返回true
//	
//	例如：{"a":123}  Int64AnyEqual(456,"a") == false 或 Int64AnyEqual(123,"a") == true
func (T *Reader) Int64AnyEqual(eq int64, keys ... string) bool {
	for _, key := range keys {
		if T.Int64(key) == eq {
			return true
		}
	}
	return false
}

//读取值是接口类型的
//	key string			键名
//	def interface{}		默认值
//	interface{}			读取的接口类型，需要转换
//	
//	例如：{"a":"b"} Interface("a","c") == "b" 或 Interface("b","c") == "c"
func (T *Reader) Interface(key string, def ... interface{}) interface{} {
	T.m.RLock()
	defer T.m.RUnlock()
	v, ok := T.M[key]
	if !ok || T.noZero && v==nil {
		for _, f := range def {
			if f != nil {
				return f
			}
		}
	}
	return v
}

//读取值是接口类型的
//	key string			键名
//	def interface{}		默认值
//	interface{}			读取的接口类型，需要转换
//	
//	例如：{"a":{"b":123}} NewInterface("a",*{"b":456}) == *{"b":123} 或 NewInterface("b",*{"b":456}) == *{"b":456}
func (T *Reader) NewInterface(key string, def ... interface{}) *Reader {
	r := NewReader(T.Interface(key, def...))
	r.noZero = T.noZero
	return r
}

//读取值是数组类型的
//	key string				键名
//	def []interface{}		默认值
//	[]interface{}			读取的数组类型
//	
//	例如：{"a":[1,3,4,5,6]} Array("a",[7,8,9,0]) == [1,3,4,5,6] 或 Array("b",[7,8,9,0]) == [7,8,9,0]
func (T *Reader) Array(key string, def ... []interface{}) []interface{} {
	T.m.RLock()
	defer T.m.RUnlock()
	v ,ok := T.M[key].([]interface{})
	if !ok || T.noZero && len(v) == 0 {
		for _, f := range def {
			if f != nil {
				return f
			}
		}
	}
	return v
}

//读取值是数组类型的
//	key string				键名
//	def []interface{}		默认值
//	[]interface{}			读取的数组类型
//	
//	例如：{"a":[1,3,4,5,6]} Array("a",[7,8,9,0]) == *[1,3,4,5,6] 或 Array("b",[7,8,9,0]) == *[7,8,9,0]
func (T *Reader) NewArray(key string, def ... []interface{}) *Reader {
	r := NewReader(T.Array(key, def...))
	r.noZero = T.noZero
	return r
}

//读取值是切片类型的，设定开始和结束位置来读取。
//	s, e int		开始，结束。
//	[]interface{}	读取到的切片
//
//	例如：[1,2,3,4,5,6] Slice(1,2) == [2]  或 Slice(8,9) == []
func (T *Reader) Slice(s, e int) []interface{} {
	T.m.RLock()
	defer T.m.RUnlock()
	l := len(T.A)
	if s > l {
		return []interface{}{}
	}
	if e > l {
		e = l
	}
	return T.A[s:e]
}

//读取值是切片类型的，设定开始和结束位置来读取。
//	s, e int		开始，结束。
//	*Reader			读取到的切片对象
//	
//	例如：[1,2,3,4,5,6] NewSlice(1,2) == *[2]  或 NewSlice(8,9) == *[]
func (T *Reader) NewSlice(s, e int) *Reader {
	r := NewReader(T.Slice(s, e))
	r.noZero = T.noZero
	return r
}

//读取值是切片类型的
//	i int				索引位置
//	def nterface{}		默认值
//	interface{}			读取到的切片值
//	
//	例如：[1,2,3,4,5,6] Index(1,11) == 1  或 Index(8,22) == 22
func (T *Reader) Index(i int, def ... interface{}) interface{} {
	as := T.Slice(i,i+1)
	if len(as) == 0 || T.noZero && as[0] == nil {
		for _, f := range def {
			if f != nil {
				return f
			}
		}
		return nil
	}
	return as[0]
}

//读取值是切片类型的
//	i int				索引位置
//	def nterface{}		默认值
//	*Reader				读取到的切片值
//	
//	例如：[1,2,[7,8,9,0],4,5,6] NewIndex(2,[11,22,33]) == *[7,8,9,0]  或 NewIndex(3,[11,22,33]) == *[] 或 NewIndex(33,[11,22,33]) == *[11,22,33]
func (T *Reader) NewIndex(i int, def ... interface{}) *Reader {
	r := NewReader(T.Index(i, def...))
	r.noZero = T.noZero
	return r
}

//读取切片类型的值是字符串类型的
//	i int				索引位置
//	def string			默认值
//	string				读取到的切片值
//	
//	例如：["1","2",[7,8,9,0],"4","5","6"] IndexString(1,"11") == "2"  或 IndexString(2,"22") == "22"
func (T *Reader) IndexString(i int, def ... string) string {
	v, ok := T.Index(i).(string)
	if !ok || T.noZero && v == "" {
		for _, f := range def {
			if f != "" {
				return f
			}
		}
	}
	return v
}


//读取切片类型的值是浮点数类型的
//	i int				索引位置
//	def float64			默认值
//	float64				读取到的浮点数
//	
//	例如：[1,2,[7,8,9,0],4,5,6] IndexInt64(1,11) == 2  或 IndexInt64(2,22) == 22
func (T *Reader) IndexFloat64(i int, def ... float64) float64 {
	rv := reflect.ValueOf(T.Index(i))
	var (
		f64 float64
		notExists bool
	)
	switch rv.Kind() {
	case reflect.Float32,reflect.Float64:
		f64 = rv.Float()
	case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
		f64 = float64(rv.Int())
	default:
		notExists = true
	}
	//1，强制不为0
	//2，不存在该字段
	if T.noZero && f64 == 0 || notExists {
		for _, f := range def {
			if f != 0 {
				return f
			}
		}
	}
	return f64
}
//读取切片类型的值是整数类型的
//	i int				索引位置
//	def int64			默认值
//	int64				读取到的整数
//	
//	例如：[1,2,[7,8,9,0],4,5,6] IndexInt64(1,11) == 2  或 IndexInt64(2,22) == 22
func (T *Reader) IndexInt64(i int, def ... int64) int64 {
	rv := reflect.ValueOf(T.Index(i))
	var (
		i64 int64
		notExists bool
	)
	switch rv.Kind() {
	case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
		i64 = rv.Int()
	case reflect.Float32,reflect.Float64:
		i64 = int64(rv.Float())
	default:
		notExists = true
	}
	//1，强制不为0
	//2，不存在该字段
	if T.noZero && i64 == 0 || notExists {
		for _, f := range def {
			if f != 0 {
				return f
			}
		}
	}
	return i64
}

//读取切片类型的值是数组类型的
//	i int				索引位置
//	def []interface{}	默认值
//	[]interface{}		读取到的切片值
//	
//	例如：[[1],[2]] IndexArray(1,[]interface{1,2}) == [2] 或 IndexArray(3,[]interface{1,2}) == [1,2]
func (T *Reader) IndexArray(i int, def ... []interface{}) []interface{} {
	v, ok := T.Index(i).([]interface{})
	if !ok || T.noZero && len(v) == 0 {
		for _, f := range def {
			if f != nil {
				return f
			}
		}
	}
	return v
}

//读取切片类型的值是数组类型的
//	i int				索引位置
//	def []interface{}	默认值
//	*Reader				读取到的切片值
//	
//	例如：[[1],[2]] NewIndexArray(1,[]interface{1,2}) == *[2] 或 NewIndexArray(3,[]interface{1,2}) == *[1,2]
func (T *Reader) NewIndexArray(i int, def ... []interface{}) *Reader {
	r := NewReader(T.IndexArray(i, def...))
	r.noZero = T.noZero
	return r
}

//重置，如果需要重置为空，需要先调用一次.Reset(nil)
//	i interface{}	支持格式，包括：map,array,slice,io.Reader,*string, []byte
//	error			错误
func (T *Reader) Reset(i interface{}) error {
	T.m.Lock()
	defer T.m.Unlock()
	var(
		err error
		tm = T.M
	)
	T.M = make(map[string]interface{})
	
	//原类型判断
	switch iv := i.(type) {
	case io.Reader:
		err = json.NewDecoder(iv).Decode(&T.M)
	case *string:
		err = json.NewDecoder(bytes.NewBufferString(*iv)).Decode(&T.M)
	case string:
		err = json.NewDecoder(bytes.NewBufferString(iv)).Decode(&T.M)
	case []byte:
		err = json.NewDecoder(bytes.NewBuffer(iv)).Decode(&T.M)
  	}
  	//发生错误，恢复
	if err != nil {
		T.M = tm
		return err
	}
	//成功后
	if len(T.M) != 0 {
		T.A = T.A[0:0]
		return nil
	}
	
	//其它类型
	rv := reflect.ValueOf(i)
	for rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch typ := rv.Kind(); typ {
	case reflect.Map:
		if m, ok := rv.Interface().(map[string]interface{}); ok {
			T.M = m
			T.A = T.A[0:0]
			return nil
		}
		return errors.New("vbody.Reader.Reset: 无法转换数据类型为map[string]interface{}")
	case reflect.Array, reflect.Slice:
		if a, ok := rv.Interface().([]interface{}); ok {
	 		T.A = a
	 		return nil
		}
		return errors.New("vbody.Reader.Reset: 无法转换数据类型为[]interface{}")
	case reflect.String:
		err = json.NewDecoder(bytes.NewBufferString(rv.String())).Decode(&T.M)
		if err != nil {
			T.M = tm
			return err
		}
	default:
		//恢复
		T.M = tm
		return errors.New("vbody.Reader.Reset: 无法转换数据类型为"+typ.String())
	}
	return nil
}

//从r读取字节串并解析成Reader
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

//Reader转字节串
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

//字节串解析成Reader
//	data []byte	字节串，如：[]byte(`{"A":1}`)
//	error		错误
func (T *Reader) UnmarshalJSON(data []byte) error {
	T.m.Lock()
	defer T.m.Unlock()
	err :=json.Unmarshal(data, &T.M)
	if err == nil {
		T.A = T.A[0:0]
	}
	return err
}
