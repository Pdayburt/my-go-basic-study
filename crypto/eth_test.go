package crypto

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
	"testing"
)

func TestETHPK(t *testing.T) {
	privateKey := "e908f86dbb4d55ac876378565aafeabc187f6690f046459397b17d9b9a19688e"
	testPrivateKey(privateKey)

}

func testPrivateKey(privateKeyHex string) {
	// 移除可能存在的 "0x" 前缀
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	// 将十六进制字符串转换为字节数组
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		fmt.Println("私钥格式错误:", err)
		return
	}

	// 转换为 ECDSA 私钥
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		fmt.Println("无效的私钥:", err)
		return
	}

	// 获取公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("无法获取公钥")
		return
	}

	// 获取地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	//0xcD49bbAc6E85fdEB167EB7cA41A945d2b8758F6F
	//0x0D1d9635D0640821d15e323ac8AdADfA9c111414

	fmt.Println("私钥验证成功!")
	fmt.Println("地址:", address.Hex())
	fmt.Println("公钥:", hexutil.Encode(crypto.FromECDSAPub(publicKeyECDSA)))
}
