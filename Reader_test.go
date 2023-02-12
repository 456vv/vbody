package vbody

import (
	"bytes"
	"testing"
)

func Test_Reader_NewReader(t *testing.T) {
	br := new(bytes.Reader)
	br.Reset([]byte(`{"username":"yourname","email":"yourname@yourdomain.com","password":"yourpassword","a":[1,2.1,"3"]}`))

	bodyr := NewReader(nil)
	bodyr.Decode(br)
	bras := bodyr.NewSlice("a")
	bras.A = append(bras.A, int(4))

	if f := bras.Float64(0); f != 1 {
		t.Fatalf("预测 1，结果 %f", f)
	}
	if f := bras.Float64(1); f != 2.1 {
		t.Fatalf("预测 2.1，结果 %f", f)
	}
	if s := bras.String(2); s != "3" {
		t.Fatalf("预测 3，结果 %s", s)
	}
	if d := bras.Int64(3); d != 4 {
		t.Fatalf("预测 4，结果 %d", d)
	}
}

func Test_Reader_IsNil(t *testing.T) {
	bodyr := NewReader(nil)
	err := bodyr.Reset([]byte(`{"a":null,"c":0, "d":[null,0]}`))
	if err != nil {
		t.Fatal(err)
	}

	if !bodyr.IsNil("a") {
		t.Fatal("error")
	}
	if !bodyr.IsNil("b") {
		t.Fatal("error")
	}
	if bodyr.IsNil("c") {
		t.Fatal("error")
	}
	d := bodyr.NewSlice("d")
	if !d.Has(0, 1) {
		t.Fatal("error")
	}

	if d.Any(0) != nil {
		t.Fatal("error")
	}
	if d.Any(1) == nil {
		t.Fatal("error")
	}
}

func Test_Reader_Reset(t *testing.T) {
	bodyr := NewReader(nil)

	if err := bodyr.Reset([]byte(`{"a":[1,2,3,4]}`)); err != nil {
		t.Fatal(err)
	}

	br := bodyr.NewSlice("a")
	if err := br.Reset(map[string]any{"1": "1"}); err != nil {
		t.Fatal(err)
	}

	t.Log(br)

	if len(br.M) == 0 {
		t.Fatal("error")
	}
}

func Test_Reader_String(t *testing.T) {
	bodyr := NewReader(nil)
	err := bodyr.Reset([]byte(`{"a":[1,2,3,4]}`))
	if err != nil {
		t.Fatal(err)
	}
	if bodyr.String("a", "123") != "123" {
		t.Fatal("error")
	}
}
