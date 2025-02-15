package util

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
)

func TestName(t *testing.T) {

	amount, err := ConvertRuneAmount("123", 2)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(amount)
}
func ConvertRuneAmount(amount string, divisibility int) (*big.Int, error) {
	// Parse the decimal string
	parts := strings.Split(amount, ".")
	if len(parts) > 2 {
		return nil, fmt.Errorf("invalid amount format")
	}

	result := new(big.Int)

	// Handle the integer part
	if len(parts) > 0 {
		result.SetString(parts[0], 10)
		result.Mul(result, big.NewInt(int64(math.Pow10(divisibility)))) // Multiply by 10^divisibility
	}

	// Handle the decimal part
	if len(parts) == 2 {
		decimal := parts[1]
		if len(decimal) > divisibility {
			decimal = decimal[:divisibility] // Truncate to divisibility decimal places
		}
		// Pad with zeros if necessary
		for len(decimal) < divisibility {
			decimal += "0"
		}

		decimalValue := new(big.Int)
		decimalValue.SetString(decimal, 10)
		result.Add(result, decimalValue)
	}

	return result, nil
}
