package crypto

/*package crypto

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/stellar/strkey"
	"github.com/stellar/transaction"

	"github.com/blocktree/go-owcdrivers/ed25519WalletKey"
)

func init() {
	coins[XLM] = newXLM
}

type xlm struct {
	*btc
}

func newXLM(key *Key) Wallet {
	key.opt.Params = &BTCParams
	token := newBTC(key).(*btc)
	token.name = "Stellar"
	token.symbol = "XLM"

	return &xlm{btc: token}
}

func (c *xlm) SignRawTransaction(signIn *SignInput) (*SignResult, error) {
	wif := signIn.PrivateKey
	if len(wif) > 0 {
		privateKey := wif

		Account := signIn.SrcAddr

		Amount := uint64(signIn.Amount)
		Fee := uint64(100)
		if signIn.Fee > 0 {
			Fee = uint64(signIn.Fee)
		}
		Sequence := uint64(signIn.Sequence)
		Destination := signIn.DestAddr
		Memos := signIn.Memo
		op := transaction.PAYMENT
		if strings.ToLower(signIn.Type) != "payment" {
			op = transaction.CREATE_ACCOUNT

		}
		tx := transaction.CreateEmptyStellarTransaction(Account, Destination, Memos, Amount, Fee, Sequence, op) //PAYMENT CREATE_ACCOUNT
		fmt.Println(tx)

		signature, err := transaction.SignRawTransction(tx, privateKey)
		if err != nil {
			return nil, err
		}
		//fmt.Println(signature)
		return &SignResult{
			Res:   1,
			RawTX: signature,
		}, nil
	}
	return &SignResult{
		Res: 0,
	}, nil
}

func (c *xlm) GetWalletAccountFromWif() (*WalletAccount, error) {
	wif := c.GetKey().Wif
	if len(wif) > 0 {
		sk, err := hex.DecodeString(wif)
		if err != nil {
			return nil, err

		}
		pkByte := ed25519WalletKey.WalletPubKeyFromKeyBytes(sk)
		address, err := strkey.Encode(strkey.VersionByteAccountID, pkByte)

		if err != nil {
			return nil, err
		}
		return &WalletAccount{
			Res:        1,
			PrivateKey: wif,
			PublicKey:  hex.EncodeToString(pkByte),
			Address:    address,
		}, nil
	}
	return &WalletAccount{
		Res: 0,
	}, nil
}

func (c *xlm) GetWalletAccount() *WalletAccount {
	seedString := c.GetKey().Seed
	if len(seedString) == 0 {
		return &WalletAccount{Res: 0}

	}
	seed, err := hex.DecodeString(seedString)
	if err != nil {
		return &WalletAccount{Res: 0, ErrMsg: err.Error()}
	}

	if len(seed) > 0 {
		path := "m/44'/148'/0'"
		childKey, _ := ed25519WalletKey.NewWalletKeyFromMasterKey(seed, path)
		address, err := strkey.Encode(strkey.VersionByteAccountID, childKey.Key.PublicKey)
		//address, err := algoDecoder.AddressEncode(childKey.Key.PublicKey)
		if err != nil {
			return &WalletAccount{
				Res:    0,
				ErrMsg: err.Error(),
			}
		}

		return &WalletAccount{
			Res:        1,
			PrivateKey: hex.EncodeToString(childKey.Key.PrivateKey),
			Address:    address,
			PublicKey:  hex.EncodeToString(childKey.Key.PublicKey),
			Seed:       c.GetKey().Seed,
		}

	}

	return &WalletAccount{
		Res: 0,
	}
}
*/
