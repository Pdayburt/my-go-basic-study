package crypto

import (
	"fmt"
	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
	"testing"
)

func TestSUIGetFaucet(t *testing.T) {
	faucetHost, err := sui.GetFaucetHost(constant.SuiTestnet)
	if err != nil {
		fmt.Println("GetFaucetHost err:", err)
		return
	}

	fmt.Println("faucetHost:", faucetHost)

	recipient := "0xdcf1fe3947cc9d189e5535bafa9de2d2e44a41a41400a71405902eb31061f216"

	header := map[string]string{}
	err = sui.RequestSuiFromFaucet(faucetHost, recipient, header)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// the successful transaction block url: https://suiexplorer.com/txblock/91moaxbXsQnJYScLP2LpbMXV43ZfngS2xnRgj1CT7jLQ?network=devnet
	fmt.Println("Request DevNet Sui From Faucet success")
}
