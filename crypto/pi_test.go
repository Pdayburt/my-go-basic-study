package crypto

import (
	"fmt"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
	"testing"
)

func TestGenerateBaseAccount(t *testing.T) {
	address, seed, err := GenerateBaseAccount()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(seed, address)

	seed2 := "SDE3VECXGDGHCYGRGZTFPQYQYWWEBZDIQT2X5ZPXA4ZQTFYWY6PHDM66"
	fromSeed2, err := GenerateBaseAccountFromSeed(seed2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(seed2, fromSeed2)
}

func TestMutedAddressFromBaseAddress(t *testing.T) {

	baseAddress := "GC57TPEPATEU4I2AT5P6DUPQ25YAUXUZ7HRLWZ343FR46Z3RZG7UCJUU"
	address, err := GenerateMuxedAddress(baseAddress, 10086)
	if err != nil {
		t.Errorf("Error generating Muxed Address: %v", err)
	}
	fmt.Printf("Address: %s\n", address)
	muxedAddress := "MD62PWT2SSS5F7CXNC4MCGHYIG2XIVWDSA4ZH6XD23D7VCGKJQRQOAAAAAAAAABHM24GK"
	a, id, err := parseMuxedAddress(muxedAddress)
	if err != nil {
		t.Errorf("Error parsing Muxed Address: %v", err)
	}
	fmt.Printf("Address: %s %d\n", a, id)

}

func TestSignPiTransferOffline(t *testing.T) {
	fromSecret := "SASBYXSYWVRZCZHW4CLNHNP5NRLEYCT5FVWERJ7ENSWNVRJUVRDW5IEG"
	toAddress := "MBXC3RUMZLFIEDNQ5ZQETCWDB7OLVWUBXOBXTK7FKA2HC5G5WVKTAAAAAAAAAABFG6MQQ"
	amount := "8.7654321" //一个pi
	var (
		sequence int64 = 78833003731615747
		baseFee  int64 = 100000
	)
	memo := "noah 20250214"
	passPhrase := "Pi Testnet"
	transferOffline, err := signPiTransferOffline(fromSecret, toAddress, amount, sequence, baseFee, memo, passPhrase)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Pi transferOffline :", transferOffline)

}

func TestCreateNewAccountTransaction(t *testing.T) {
	fromSecret := "SASBYXSYWVRZCZHW4CLNHNP5NRLEYCT5FVWERJ7ENSWNVRJUVRDW5IEG"
	newAccountAddress := "GBXC3RUMZLFIEDNQ5ZQETCWDB7OLVWUBXOBXTK7FKA2HC5G5WVKTAERJ"
	startingBalance := "22"
	var (
		sequence int64 = 78833003731615750
		baseFee  int64 = 100000
	)
	passPhrase := "Pi Testnet"
	//passPhrase := "Pi Network"
	transaction, err := ActiveAccountTransaction(fromSecret, newAccountAddress, startingBalance, sequence, baseFee, passPhrase)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Transaction :", transaction)
}

func TestActiveAndTransferPI(t *testing.T) {
	fromSecret := "SASBYXSYWVRZCZHW4CLNHNP5NRLEYCT5FVWERJ7ENSWNVRJUVRDW5IEG"
	toAddress := "MA47SETHDSCWY4JTKYYIZ4PNGBW4IJWISBCZNT6Y2MSZQFW77F7RAAAAAAAAAABHM3FCK"
	amount := "3.21"
	var (
		sequence int64 = 78833003731615750
		baseFee  int64 = 100000
	)
	passPhrase := "Pi Testnet"
	memo := "noah 2025021502"
	res, err := activeAndTransferPI(fromSecret, toAddress, amount, sequence, baseFee, memo, passPhrase, false)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Transaction :", res)

}

// toAddress位muxed类型的地址
func activeAndTransferPI(fromSecret string, toAddress string, amount string,
	sequence int64, baseFee int64, memo string, passPhrase string, needActive bool) (string, error) {

	baseAddress, _, err := parseMuxedAddress(toAddress)
	if err != nil {
		return "-1", fmt.Errorf("muxed地址无法解析出Base地址: %w", err)
	}
	// 解析密钥对
	kp, err := keypair.Parse(fromSecret)
	if err != nil {
		return "-1", fmt.Errorf("pi私钥错误，无法解析: %w", err)
	}
	var operations []txnbuild.Operation

	//如果账号需要激活激活，则先激活再转账
	if needActive {
		createAccountOp := txnbuild.CreateAccount{
			Destination: baseAddress,
			Amount:      "20", // 测试网最小激活金额20 ，主网为1
		}
		operations = append(operations, &createAccountOp)
	}
	//不需要激活则直接转账
	// 创建支付操作
	paymentOp := txnbuild.Payment{
		Destination: toAddress,
		Amount:      amount,
		Asset:       txnbuild.NativeAsset{},
	}
	operations = append(operations, &paymentOp)

	// 创建交易
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: kp.Address(), Sequence: sequence},
			IncrementSequenceNum: true,
			Operations:           operations,
			BaseFee:              baseFee,
			Memo:                 txnbuild.MemoText(memo),
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewTimeout(300)},
		},
	)
	if err != nil {
		return "-1", fmt.Errorf("构建Pi交易失败：: %w", err)
	}

	// 签名交易
	tx, err = tx.Sign(passPhrase, kp.(*keypair.Full))
	if err != nil {
		return "-1", fmt.Errorf("创建Pi交易失败： %w", err)
	}
	// 获取签名后的XDR
	return tx.Base64()
}

func ActiveAccountTransaction(fromSecret string, newAccountAddress string, startingBalance string,
	sequence int64, baseFee int64, passPhrase string) (string, error) {

	// 解析密钥对
	kp, err := keypair.Parse(fromSecret)
	if err != nil {
		return "", err
	}

	// 创建账户操作
	op := txnbuild.CreateAccount{
		Destination: newAccountAddress,
		Amount:      startingBalance,
	}

	// 创建交易
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: kp.Address(), Sequence: sequence},
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              baseFee,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewTimeout(300)},
		},
	)
	if err != nil {
		return "", err
	}

	// 签名交易
	tx, err = tx.Sign(passPhrase, kp.(*keypair.Full))
	if err != nil {
		return "", err
	}

	// 获取签名后的XDR
	return tx.Base64()
}

func signPiTransferOffline(fromSecret string, toAddress string, amount string,
	sequence int64, baseFee int64, memo string, passPhrase string) (string, error) {

	// 解析密钥对
	kp, err := keypair.Parse(fromSecret)
	if err != nil {
		return "Pi私钥错误，无法解析", fmt.Errorf("pi私钥错误，无法解析: %w", err)
	}

	// 创建支付操作
	op := txnbuild.Payment{
		Destination: toAddress,
		Amount:      amount,
		Asset:       txnbuild.NativeAsset{},
	}

	// 创建交易
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: kp.Address(), Sequence: sequence},
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              baseFee,
			Memo:                 txnbuild.MemoText(memo),
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewTimeout(300)},
		},
	)
	if err != nil {
		return "", err
	}

	// 签名交易
	tx, err = tx.Sign(passPhrase, kp.(*keypair.Full))
	if err != nil {
		return "", err
	}

	// 获取签名后的XDR
	return tx.Base64()
}

func GenerateBaseAccount() (address string, seed string, err error) {
	kp, err := keypair.Random()
	if err != nil {
		return "", "", err
	}
	return kp.Address(), kp.Seed(), nil
}

func GenerateBaseAccountFromSeed(seed string) (string, error) {
	kp, err := keypair.Parse(seed)
	if err != nil {
		return "", err
	}

	return kp.Address(), nil
}

func GenerateMuxedAddress(baseAddress string, id uint64) (string, error) {
	// 解析基础地址
	accountID := xdr.AccountId{}
	err := accountID.SetAddress(baseAddress)
	if err != nil {
		return "", fmt.Errorf("invalid base address: %v", err)
	}

	// 创建 Muxed 账户
	muxedAccount := xdr.MuxedAccount{
		Type: xdr.CryptoKeyTypeKeyTypeMuxedEd25519,
		Med25519: &xdr.MuxedAccountMed25519{
			Id:      xdr.Uint64(id),
			Ed25519: *accountID.Ed25519,
		},
	}

	// 编码为 M... 格式的地址
	muxedAddress := muxedAccount.Address()
	return muxedAddress, nil
}

func parseMuxedAddress(muxedAddress string) (baseAddress string, id uint64, err error) {
	// 解析 Muxed 地址
	muxedAccount := xdr.MuxedAccount{}
	err = muxedAccount.SetAddress(muxedAddress)
	if err != nil {
		return "", 0, fmt.Errorf("invalid muxed address: %v", err)
	}

	// 获取基础地址
	baseAddress = muxedAccount.ToAccountId().Address()

	// 获取 ID
	if muxedAccount.Type == xdr.CryptoKeyTypeKeyTypeMuxedEd25519 {
		id = uint64(muxedAccount.Med25519.Id)
	}
	return baseAddress, id, nil
}
