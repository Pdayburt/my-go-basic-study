package crypto

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

func TestDecimal(t *testing.T) {

	amtDec, err := decimal.NewFromString("1000000000000000000")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(amtDec)

}
