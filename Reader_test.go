package vbody

import(
	"testing"
	"bytes"
//	"io"
)


func Test_Reader_NewReader(t *testing.T) {
	br := &bytes.Reader{}
	br.Reset([]byte(`{"username":"yourname","email":"yourname@yourdomain.com","password":"yourpassword","a":[1,2,"3"]}`))

	bodyr := NewReader(nil)
	bodyr.ReadFrom(br)
	
	bras := bodyr.NewArray("a", nil).NewSlice(1,20)
	
	if f := bras.IndexFloat64(0, -1); f != float64(2) {
		t.Fatalf("预测 2，结果 %f", f)
	}
	if s := bras.IndexString(0, "-1"); s != "-1" {
		t.Fatalf("预测 2，结果 %s", s)
	}
}