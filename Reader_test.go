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