package web

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"golang.org/x/crypto/bcrypt"
	"math"
	"testing"
)

func TestEncrypt(t *testing.T) {

	p := "123d"

	password, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(password))

	err = bcrypt.CompareHashAndPassword(password, []byte(p))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCeil(t *testing.T) {
	ceil1 := math.Ceil(float64(28) / float64(50))
	fmt.Println(ceil1)
	ceil2 := math.Ceil(28 / 50)
	fmt.Println(ceil2)
}

func TestQuery(t *testing.T) {

	request := resty.New().SetBaseURL("https://api.xion-testnet-1.burnt.com").R()

	for i := 0; i < 10; i++ {
		response, err := request.Get("cosmos/tx/v1beta1/txs/654e331886aa24fd712ac32e4449a6c5634551a22b5c8da30ce5db61cd33f8b0")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(response.StatusCode())
	}

}
