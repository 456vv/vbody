# vbody [![Build Status](https://travis-ci.org/456vv/vbody.svg?branch=master)](https://travis-ci.org/456vv/vbody)
golang vbody，在web的请求与响应中，body使用的JSON数据格式。

# **列表：**
```go
type Reader struct{													// 读取
	M	map[string]interface{}											// 记录对象
	A	[]interface{}													// 记录数组
	m	sync.RWMutex													// 安全锁
	err	error															// 错误
}
	func NewReader(i interface{}) *Reader								// 读取
	func (T *Reader) Err() error										// 错误
    func (T *Reader) IsNil(keys ...interface{}) bool                    // 是nil
	func (T *Reader) Has(keys ...interface{}) bool 						// 判断
	func (T *Reader) String(key, def string) string 					// 读取字符串
	func (T *Reader) StringAnyEqual(eq string, keys ... string) bool	// 字符串判断
	func (T *Reader) Bool(key string, def bool) bool 					// 读取布尔值
	func (T *Reader) BoolAnyEqual(eq bool, keys ... string) bool		// 布尔值判断
	func (T *Reader) Float64(key string, def float64) float64			// 读取浮点数
	func (T *Reader) Float64AnyEqual(eq float64, keys ... string) bool	// 浮点数判断
	func (T *Reader) Int64(key string, def int64) int64					// 读取整数
	func (T *Reader) Int64AnyEqual(eq int64, keys ... string) bool		// 整数判断
	func (T *Reader) Interface(key string, def interface{}) interface{}	// 读取接口
	func (T *Reader) NewInterface(key string, def interface{}) *Reader	// 读取接口*
	func (T *Reader) Array(key string, def []interface{}) []interface{}	// 读取数组
	func (T *Reader) NewArray(key string, def []interface{}) *Reader	// 读取数组*
	func (T *Reader) Slice(s, e int) []interface{}						// 读取切片
	func (T *Reader) NewSlice(s, e int) *Reader							// 读取切片*
	func (T *Reader) Index(i int, def interface{}) interface{}			// 读取数组的单个值
	func (T *Reader) NewIndex(i int, def interface{}) *Reader			// 读取数组的单个值*
	func (T *Reader) IndexString(i int, def string) string				// 读取数组的单个字符串
	func (T *Reader) IndexFloat64(i int, def float64) float64			// 读取数组的单个浮点数
	func (T *Reader) IndexInt64(i int, def int64) int64					// 读取数组的单个整数
	func (T *Reader) IndexArray(i int, def []interface{}) []interface{}	// 读取数组的单个数组
	func (T *Reader) NewIndexArray(i int, def []interface{}) *Reade		// 读取数组的单个数组*
	func (T *Reader) Reset(i interface{}) error							// 重置
	func (T *Reader) ReadFrom(r io.Reader) error						// 从r读取导入
	func (T *Reader) MarshalJSON() ([]byte, error)						// 编码json
	func (T *Reader) UnmarshalJSON(data []byte) error					// 解码json

type Writer struct{													// 写入
	M map[string]interface{}											// 写记录
}
	func NewWriter() *Writer											// 写入器
	func (T *Writer) Status(d int)										// 状态
	func (T *Writer) Message(s interface{})								// 信息
	func (T *Writer) Messagef(f string, a ...interface{})				// 信息（支持格式）
	func (T *Writer) SetResult(i interface{})							// 设置结果
	func (T *Writer) Result(key string, i interface{})					// 结果
	func (T *Writer) WriteTo(w io.Writer) (n int64, err error) 			// 写入到w
```