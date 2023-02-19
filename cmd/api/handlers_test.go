package main

import (
	"encoding/json"
	"testing"
)

var unmarshalRes error

func anonymousStructBench(d []byte) error {
	var req struct {
		Test  string `json:"test"`
		Param string `json:"param"`
	}

	err := json.Unmarshal(d, &req)
	if err != nil {
		return err
	}

	return nil
}

func BenchmarkAnonymousStruct(b *testing.B) {
	reqData := []byte(`{ "test": "onomatopeja", "param": "123" }`)
	for n := 0; n < b.N; n++ {
		unmarshalRes = anonymousStructBench(reqData)
	}
}

type testJSONReq struct {
	Test  string `json:"test"`
	Param string `json:"param"`
}

func declaredStructBench(d []byte) error {
	var req testJSONReq

	err := json.Unmarshal(d, &req)
	if err != nil {
		return err
	}

	return nil
}

func BenchmarkDeclaredStruct(b *testing.B) {
	reqData := []byte(`{ "test": "onomatopeja", "param": "123" }`)
	for n := 0; n < b.N; n++ {
		unmarshalRes = declaredStructBench(reqData)
	}
}
