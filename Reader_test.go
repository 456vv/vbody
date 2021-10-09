package vbody

import(
	"testing"
	"bytes"
)


func Test_Reader_NewReader(t *testing.T) {
	br := &bytes.Reader{}
	br.Reset([]byte(`{"username":"yourname","email":"yourname@yourdomain.com","password":"yourpassword","a":[1,2.1,"3"]}`))

	bodyr := NewReader(nil)
	bodyr.ReadFrom(br)
	bras := bodyr.NewArray("a", nil)
	bras.A=append(bras.A,int(4))

	if f := bras.IndexFloat64(0, -1); f != 1 {
		t.Fatalf("预测 1，结果 %f", f)
	}
	if f := bras.IndexFloat64(1, -1); f != 2.1 {
		t.Fatalf("预测 2.1，结果 %f", f)
	}
	if s := bras.IndexString(2, "-3"); s != "3" {
		t.Fatalf("预测 3，结果 %s", s)
	}
	if d := bras.IndexInt64(3, -3); d != 4 {
		t.Fatalf("预测 4，结果 %d",d)
	}
}

func Test_Reader_IsNil(t *testing.T){
	bodyr := NewReader(nil)
	err := bodyr.Reset([]byte(`{"a":null,"c":0, "d":[null,0]}`))
	if err != nil  {
		t.Fatal(err)
	}
	if !bodyr.IsNil("a") {
		t.Fatal("错误")
	}
	if !bodyr.IsNil("b") {
		t.Fatal("错误")
	}
	if bodyr.IsNil("c") {
		t.Fatal("错误")
	}
	d := bodyr.NewArray("d", nil)
	if !d.Has(0,1) {
		t.Fatal("错误")
	}
	
	if d.Index(0, 1) != nil {
		t.Fatal("错误")
	}
	if d.Index(1, 1) == nil {
		t.Fatal("错误")
	}
}

func Test_Reader_NoZero(t *testing.T){
	bodyr := NewReader(nil)
	err := bodyr.Reset([]byte(`{"a":null,"c":0, "d":[null,0]}`))
	if err != nil  {
		t.Fatal(err)
	}
	nr := bodyr.NoZero(true)
	if !nr.noZero {
		t.Fatal("错误")
	}
	i := nr.Int64("c", 1)
	if i != 1 {
		t.Fatal("错误")
	}
	if bodyr.noZero {
		t.Fatal("错误")
	}
	err = bodyr.Reset([]byte(`{"b":"1"}`))
	if err != nil  {
		t.Fatal(err)
	}
	if !nr.StringAnyEqual("", "b") {
		t.Fatal("错误")
	}
	//t.Log(bodyr.M, nr.M)
}


func Test_Reader_Reset(t *testing.T){
	bodyr := NewReader(nil)
	err := bodyr.Reset([]byte(`{"a":[1,2,3,4]}`))
	if err != nil  {
		t.Fatal(err)
	}
	bodyr = bodyr.NewArray("a",nil)
	nr := bodyr.NoZero(true)
	
	bodyr.Reset(map[string]interface{}{"1":"1"})
	if err != nil  {
		t.Fatal(err)
	}
	if len(bodyr.A) != 0 || len(nr.A) == 0 {
		t.Fatal(bodyr.A, nr.A)
	}
}

func Test_Reader_1(t *testing.T){
	bodyr := NewReader(nil)
	err := bodyr.Reset([]byte(`{"a":[1,2,3,4]}`))
	if err != nil  {
		t.Fatal(err)
	}
	if bodyr.String("a","123") != "123" {
		t.Fatal("错误")
	}
}
