package vbody

import (
	"testing"
)

func Test_Change_Update_1(t *testing.T){
	
	tests := []struct{
		name string
		value interface{}
		result interface{}
	}{
		{name:"a", value:"11", result:"11"},
		{name:"c", value:"11", result:""},
	}
	
	r := NewReader(`{"a":"1","b":[1,2,3,4]}`)
	c := r.Change()
	
	for index,test := range tests {
		c.Update(test.name, test.value)
		s := r.String(test.name)
		if s != test.result {
			t.Fatalf("update fail, in %d result %v", index, s)
		}
	}
}
func Test_Change_Update_2(t *testing.T){
	
	tests := []struct{
		name int
		value interface{}
		result interface{}
	}{
		{name:1, value:"11", result:"11"},
		{name:11, value:"11", result:nil},
	}
	
	r := NewReader(`{"a":"1","b":[1,2,3,4]}`).NewInterface("b")
	c := r.Change()
	
	for index,test := range tests {
		c.Update(test.name, test.value)
		i := r.Index(test.name)
		if i != test.result {
			t.Fatalf("update fail, in %d result %v", index, i)
		}
	}
}

func Test_Change_Set_1(t *testing.T){
	tests := []struct{
		name string
		value interface{}
	}{
		{name:"a", value:"11"},
		{name:"c", value:"11"},
	}
	r := NewReader(`{"a":"1","b":[1,2,3,4]}`)
	c := r.Change()

	for index, test := range tests {
		c.Set(test.name, test.value)
		s := r.String(test.name)
		if s != test.value {
			t.Fatalf("set fail, in %d result %v", index, s)
		}
	}
}

func Test_Change_Set_2(t *testing.T){
	tests := []struct{
		name int
		value interface{}
	}{
		{name:1, value:"a"},
		{name:11, value:"11"},
		{name:4, value:"4"},
	}
	r := NewReader(`{"a":"1","b":[0,1,2,3]}`)
	br := r.NewInterface("b")
	c := br.Change()

	for index, test := range tests {
		c.Set(test.name, test.value)
		s := br.Index(test.name)
		if s != test.value {
			t.Fatalf("set fail, in %d result %v", index, s)
		}
	}
}
func Test_Change_Delete_1(t *testing.T){
	tests := []struct{
		name string
	}{
		{name:"a"},
		{name:"b"},
	}
	r := NewReader(`{"a":"1","b":[0,1,2,3]}`)
	c := r.Change()
	
	for index, test := range tests {
		c.Delete(test.name)
		s := r.String(test.name)
		if s != "" {
			t.Fatalf("set fail, in %d result %v", index, s)
		}
	}
}

func Test_Change_Delete_2(t *testing.T){
	tests := []struct{
		name int
	}{
		{name:1},
		{name:11},
	}
	r := NewReader(`{"a":"1","b":[0,1,2,3]}`)
	
	br := r.NewInterface("b")
	c := br.Change()
	
	for index, test := range tests {
		c.Delete(test.name)
		s := r.Index(test.name)
		if s != nil {
			t.Fatalf("set fail, in %d result %v", index, s)
		}
	}
}