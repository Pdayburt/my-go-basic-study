package crypto

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/stellar/go/exp/crypto/derivation"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"regexp"
	"testing"
)

func TestReg(t *testing.T) {
	// 编译正则表达式
	pattern := `^M[A-Za-z0-9]{68}$`
	re, err := regexp.Compile(pattern)

	if err != nil {
		fmt.Printf("正则表达式编译错误: %v\n", err)
		return
	}

	// 测试一些字符串
	testStrings := []string{
		"MBG63IFS3EDNQSD44UTJG3NPUEAHZ5MZ43WJTVU4FV7UW2KFY2TLCAAAAAAAAABHM2GXQ", // 68位字符
		"MBG63IFS3EDNQSD44UTJG3NPUEAHZ5MZ43WJTVU4FV7UW2KFY2TLCAAAAAAAAABHDKLMO", // 太短
		"MARIRPA55VFGI7MYUW6MCZKK4F7BTXS6CFHXXS4DOVB6CQTHDTLLGAAAAAAAAABHMZA6O", // 包含特殊字符
		"MAF4Z2TCIM6KRZSGGVR7V7RBT2OXKWGK5GH6QGK47U6GRMSV33ZX6AAAAAAAAABHM2QXY", // 包含特殊字符
	}

	for _, s := range testStrings {
		fmt.Printf("字符串 '%s' 是否匹配: %v\n", s, re.MatchString(s))
	}
}

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

	baseAddress := "GDMDXVLQAZ6AIAK6YOG7UIFOITBIHZYNBAHTUUQ5XJ5XWWEGBSH5P7L4"
	address, err := GenerateMuxedAddress(baseAddress, 10086)
	if err != nil {
		t.Errorf("Error generating Muxed Address: %v", err)
	}
	fmt.Printf("Address: %s\n", address)
	muxedAddress := "MBG63IFS3EDNQSD44UTJG3NPUEAHZ5MZ43WJTVU4FV7UW2KFY2TLCAAAAAAAAABHM2GXQ"
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

func TestPubk(t *testing.T) {
	//mnemonic := strings.Join(words, " ")
	//println("Mnemonic:", mnemonic)
	mnemonic := "useless begin slam ribbon suggest kid acquire joy split middle heavy reject knife media result permit robust above enact broken history imitate deer idle"
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		panic(err)
	}

	println("BIP39 Seed:", hex.EncodeToString(seed))

	//314159
	//masterKey, err := derivation.DeriveForPath(derivation.StellarAccountPrefix, seed)
	masterKey, err := derivation.DeriveForPath("m/44'/314159'", seed)
	if err != nil {
		panic(err)
	}

	println("m/44'/148' key:", hex.EncodeToString(masterKey.Key))
	key, err := masterKey.Derive(derivation.FirstHardenedIndex + 0)
	kp, err := keypair.FromRawSeed(key.RawSeed())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Address:%v ,seed :%v", kp.Address(), kp.Seed())
	//GDJ7K7HKBSDHDRCBC7ACHBD6BZXM5TID24RHJBUY37O2IN5DDMIKKNS2

	println("")
}

func GeneratePiKeysFromMnemonic111(mnemonic string) (string, string, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", "", fmt.Errorf("无效的助记词")
	}

	// 生成种子
	seed := bip39.NewSeed(mnemonic, "")

	// 生成主密钥
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", "", fmt.Errorf("生成主密钥失败: %v", err)
	}

	// 按照 BIP44 路径派生：m/44'/314159'/0'/0/0
	// 注意：bip32.FirstHardenedChild 是 0x80000000
	path := []uint32{
		44 + bip32.FirstHardenedChild,     // purpose: 44'
		314159 + bip32.FirstHardenedChild, // coinType: 314159' (Pi Network)
		0 + bip32.FirstHardenedChild,      // account: 0'
		0,                                 // change: 0
		0,                                 // addressIndex: 0
	}

	// 按照路径派生子密钥
	key := masterKey
	for i, childNum := range path {
		key, err = key.NewChildKey(childNum)
		if err != nil {
			return "", "", fmt.Errorf("派生第 %d 层子密钥失败: %v", i, err)
		}
	}

	// 使用派生的密钥作为 ED25519 私钥种子
	privateKey := ed25519.NewKeyFromSeed(key.Key[:32])
	publicKey := privateKey.Public().(ed25519.PublicKey)

	// 编码为 Stellar 格式（因为 Pi Network 使用相同的地址格式）
	stellarPublicKey, err := strkey.Encode(strkey.VersionByteAccountID, publicKey)
	if err != nil {
		return "", "", fmt.Errorf("编码公钥失败: %v", err)
	}

	stellarPrivateKey, err := strkey.Encode(strkey.VersionByteSeed, privateKey.Seed())
	if err != nil {
		return "", "", fmt.Errorf("编码私钥失败: %v", err)
	}

	return stellarPublicKey, stellarPrivateKey, nil
}

func GenerateFromMnemonic(mnemonic string) (*keypair.Full, error) {
	// 验证助记词
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("无效的助记词")
	}

	// 从助记词生成种子
	seed := bip39.NewSeed(mnemonic, "")
	hash := sha256.Sum256(seed) // 生成32字节的哈希

	// 从种子生成 Stellar 密钥对
	// 注意：这里我们只使用种子的前 32 字节，因为 Stellar 的私钥长度是 32 字节
	kp, err := keypair.FromRawSeed(hash)
	if err != nil {
		return nil, fmt.Errorf("生成密钥对失败: %v", err)
	}

	return kp, nil
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
