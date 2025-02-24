package jsonParse

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParsePerson(t *testing.T) {
	// 测试用例1：Token 为字符串
	jsonStr1 := `{
        "name": "Alice",
        "action": {
            "type": "transfer",
            "token": "abc123",
            "amount": "100"
        }
    }`

	// 测试用例2：Token 为数字
	jsonStr2 := `{
        "name": "Bob",
        "action": {
            "type": "transfer",
            "token": 12345,
            "amount": "200"
        }
    }`

	var p1, p2 Person
	err := json.Unmarshal([]byte(jsonStr1), &p1)
	if err != nil {
		t.Error(err)

	}
	fmt.Printf("p1: %v\n", p1)
	err = json.Unmarshal([]byte(jsonStr2), &p2)
	if err != nil {
		t.Error(err)

	}
	fmt.Printf("p2: %v\n", p2)

}

type Action struct {
	Type   string `json:"type"`
	Token  string `json:"token"`
	Amount string `json:"amount"`
}

type Person struct {
	Name    string `json:"name"`
	Actions Action `json:"action"`
}

func (a *Action) UnmarshalJSON(bytes []byte) error {
	type ActionTemp struct {
		Type   string      `json:"type"`
		Token  interface{} `json:"token"`
		Amount string      `json:"amount"`
	}
	var temp ActionTemp
	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return err
	}
	a.Type = temp.Type
	a.Amount = temp.Amount

	switch t := temp.Token.(type) {
	case string:
		a.Token = t
	case float64:
		a.Token = fmt.Sprintf("%f", t)
	default:
		a.Token = ""
	}
	return nil
}
