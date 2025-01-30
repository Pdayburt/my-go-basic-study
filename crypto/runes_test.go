package crypto

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"math/big"
	"strconv"
	"strings"
	"testing"
)

// AddressType represents different Bitcoin address types
type AddressType int

// NetworkType represents the Bitcoin network type
type NetworkType string

const (
	MainNet NetworkType = "mainnet"
	TestNet NetworkType = "testnet"
)

// UnspentOutput represents a UTXO with additional metadata
type UnspentOutput struct {
	TxID         string
	Vout         uint32
	Satoshis     int64
	ScriptPubKey string
	PubKey       string
	AddressType  AddressType
	Inscriptions []Inscription
	Atomicals    []Atomical
	Runes        []Rune
	RawTx        string
}

// Inscription represents an inscription on a UTXO
type Inscription struct {
	InscriptionID     string
	InscriptionNumber int64
	Offset            int64
}

// Atomical represents an Atomical on a UTXO
type Atomical struct {
	AtomicalID     string
	AtomicalNumber int64
	Type           string // "FT" or "NFT"
	Ticker         string
	AtomicalValue  int64
}

// Rune represents a Rune on a UTXO
type Rune struct {
	RuneID string
	Amount string
}

// RuneTransferParams contains all parameters needed for a Rune transfer
type RuneTransferParams struct {
	AssetUTXOs   []UnspentOutput
	BTCUTXOs     []UnspentOutput
	AssetAddress string //asset找零地址
	BTCAddress   string //btc找零地址
	ToAddress    string
	NetworkType  NetworkType
	RuneID       string
	RuneAmount   string
	OutputValue  int64
	FeeRate      int64
	EnableRBF    bool
	PrivateKey   string
}

// ToSignInput represents an input that needs to be signed
type ToSignInput struct {
	TxID            string
	PublicKeyScript string
	Vout            uint32
	Value           int64
	PrivateKey      string
}

// EncodeVarInt encodes a number as a varint into a byte slice
func EncodeVarInt(n *big.Int, buf *bytes.Buffer) {
	for n.BitLen() > 7 {
		bit := new(big.Int).And(n, big.NewInt(0x7f))
		buf.WriteByte(byte(bit.Int64()) | 0x80)
		n.Rsh(n, 7)
	}
	buf.WriteByte(byte(n.Int64()))
}
func TestRunesTransfer(t *testing.T) {

	asseUtxo := []UnspentOutput{
		{
			TxID:         "c7c49fce8facb992b381a22c980e486ebec9d089c83f1a36387368783e01848e",
			Vout:         1,
			Satoshis:     546,
			ScriptPubKey: "5120a3e8a12e63e9df54000c6dc5113fc8bfea2d0c1a3206e799c25279b1f880153e",
			Runes: []Rune{
				{
					RuneID: "2584592:58",
					Amount: "100000000",
				},
			},
		},
	}
	btcUTXO := []UnspentOutput{
		{
			TxID:         "3d32c14fdd7536fd92464dd96ef1f2b6a37f3914b02530ecccaeee6f58336414",
			Vout:         3,
			Satoshis:     2596283,
			ScriptPubKey: "5120a3e8a12e63e9df54000c6dc5113fc8bfea2d0c1a3206e799c25279b1f880153e",
		},
	}

	rnTrfansferParam := RuneTransferParams{
		AssetUTXOs:   asseUtxo,
		BTCUTXOs:     btcUTXO,
		AssetAddress: "tb1p5052ztnra804gqqvdhz3z07ghl4z6rq6xgrw0xwz2fumr7yqz5lqjkwgzs",
		BTCAddress:   "tb1p5052ztnra804gqqvdhz3z07ghl4z6rq6xgrw0xwz2fumr7yqz5lqjkwgzs",
		ToAddress:    "tb1pva2w66x5r3jcq6t5a86kldyjgysuggtcek3mdzr7gauv4hhzfctszdc0mt",
		NetworkType:  "testnet3",
		RuneID:       "2584592:58",
		RuneAmount:   "50000000",
		OutputValue:  546,
		FeeRate:      3011,
		EnableRBF:    false,
		PrivateKey:   "5bbb59c4f0004715feab40ae9468fc52e2ac5c58462ff55a97bccb5c18ec89fb",
	}
	SendRunes(rnTrfansferParam)

}

// SendRunes creates a transaction to transfer Runes
func SendRunes(params RuneTransferParams) {

	// Create new transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	// Track inputs that need signing
	var toSignInputs []ToSignInput

	// Add asset inputs
	fromRuneAmount := big.NewInt(0)
	runesMap := make(map[string]bool)

	for _, utxo := range params.AssetUTXOs {
		// Add input
		txHash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			fmt.Printf("TxID解析错误: %v, TxID: %s\n", err, utxo.TxID)
			continue
		}
		prevOut := wire.NewOutPoint(txHash, utxo.Vout)
		tx.AddTxIn(wire.NewTxIn(prevOut, nil, nil))

		// Track signing
		toSignInputs = append(toSignInputs, ToSignInput{
			TxID:            utxo.TxID,
			Vout:            utxo.Vout,
			Value:           utxo.Satoshis,
			PublicKeyScript: utxo.ScriptPubKey,
			PrivateKey:      params.PrivateKey,
		})

		// Track Rune amounts
		for _, rune := range utxo.Runes {
			runesMap[rune.RuneID] = true
			if rune.RuneID == params.RuneID {
				runeAmt := new(big.Int)
				runeAmt.SetString(rune.Amount, 10)
				fromRuneAmount.Add(fromRuneAmount, runeAmt)
			}
		}
	}

	// Calculate change amount
	transferAmount := new(big.Int)
	transferAmount.SetString(params.RuneAmount, 10)
	changeAmount := new(big.Int).Sub(fromRuneAmount, transferAmount)

	if changeAmount.Sign() < 0 {

	}

	// Determine if change output is needed
	needChange := len(runesMap) > 1 || changeAmount.Sign() > 0

	// Create OP_RETURN output
	var payload bytes.Buffer
	EncodeVarInt(big.NewInt(0), &payload) // Version

	// Parse RuneID
	//	runeBlock, runeTx := ParseRuneID(params.RuneID)
	runeData := strings.Split(params.RuneID, ":")
	runeBlock, _ := strconv.ParseInt(runeData[0], 10, 64)
	runeTx, _ := strconv.ParseInt(runeData[1], 10, 64)

	EncodeVarInt(big.NewInt(runeBlock), &payload)
	EncodeVarInt(big.NewInt(runeTx), &payload)
	EncodeVarInt(transferAmount, &payload)

	if needChange {
		EncodeVarInt(big.NewInt(2), &payload) // 2 = send with change
	} else {
		EncodeVarInt(big.NewInt(1), &payload) // 1 = send without change
	}

	// Add OP_RETURN output
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(txscript.OP_13)
	builder.AddData(payload.Bytes())
	opReturnScript, _ := builder.Script()
	tx.AddTxOut(wire.NewTxOut(0, opReturnScript))

	// Add change output if needed
	if needChange {
		addr, _ := btcutil.DecodeAddress(params.AssetAddress, GetNetwork(params.NetworkType))
		script, _ := txscript.PayToAddrScript(addr)
		tx.AddTxOut(wire.NewTxOut(params.OutputValue, script))
	}

	// Add recipient output
	toAddr, _ := btcutil.DecodeAddress(params.ToAddress, GetNetwork(params.NetworkType))
	toScript, _ := txscript.PayToAddrScript(toAddr)
	tx.AddTxOut(wire.NewTxOut(params.OutputValue, toScript))

	// Add BTC UTXOs for fees
	btcToSignInputs, err := AddSufficientUTXOsForFee(tx, params.BTCUTXOs, params.FeeRate, params.BTCAddress, params.PrivateKey)
	if err != nil {

	}
	toSignInputs = append(toSignInputs, btcToSignInputs...)

	var signedTxBuf1 bytes.Buffer
	tx.Serialize(&signedTxBuf1)
	unsignedTxHex := hex.EncodeToString(signedTxBuf1.Bytes())
	fmt.Println("unsignedTxHex:", unsignedTxHex)
	//err = SigVins(GetNetwork(params.NetworkType), tx, toSignInputs, params.PrivateKey)
	err = SignVinUtxo(GetNetwork(params.NetworkType), tx, toSignInputs, params.PrivateKey)
	if err != nil {
		fmt.Printf("SignVinUtxo error: %v\n", err)
	}
	// 序列化签名后的交易
	var signedTxBuf2 bytes.Buffer
	tx.Serialize(&signedTxBuf2)
	signedTxHex := hex.EncodeToString(signedTxBuf2.Bytes())
	fmt.Println("signedTxHex:", signedTxHex)
}

func SignVinUtxo(chainParams *chaincfg.Params, tx *wire.MsgTx, vins []ToSignInput, privateKey string) error {

	prevOuts := make(map[wire.OutPoint]*wire.TxOut)
	for _, vin := range vins {
		hash, err := chainhash.NewHashFromStr(vin.TxID)
		if err != nil {
			return fmt.Errorf("chainhash.NewHashFromStr err : %w", err)
		}
		pk, err := hex.DecodeString(vin.PublicKeyScript)
		if err != nil {
			return fmt.Errorf("hex.DecodeString err : %w", err)
		}
		outPoint := wire.NewOutPoint(hash, vin.Vout)
		out := wire.NewTxOut(vin.Value, pk)
		prevOuts[*outPoint] = out
	}
	privKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return fmt.Errorf("hex.DecodeString err : %w", err)
	}
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

	txSigHash := txscript.NewTxSigHashes(tx, txscript.NewMultiPrevOutFetcher(prevOuts))

	for i, vin := range vins {
		tx.TxIn[i].SignatureScript = nil
		tx.TxIn[i].Witness = nil

		txInPkScript, err := hex.DecodeString(vin.PublicKeyScript)
		if err != nil {
			return fmt.Errorf("decode script error: %v", err)
		}

		scriptClass, _, _, err := txscript.ExtractPkScriptAddrs(txInPkScript, chainParams)
		if err != nil {
			return fmt.Errorf("extract script type error: %v", err)
		}

		switch scriptClass {
		case txscript.PubKeyHashTy:
			script, err := txscript.SignatureScript(
				tx,
				i,
				txInPkScript,
				txscript.SigHashAll,
				privKey,
				true,
			)
			if err != nil {
				return fmt.Errorf("sign legacy input error: %v", err)
			}
			tx.TxIn[i].SignatureScript = script

		case txscript.ScriptHashTy:
			witnessProg := btcutil.Hash160(privKey.PubKey().SerializeCompressed())
			// 创建赎回脚本
			redeemScript := []byte{
				txscript.OP_0,
				byte(len(witnessProg)),
			}
			redeemScript = append(redeemScript, witnessProg...)
			tx.TxIn[i].SignatureScript = redeemScript

			witness, err := txscript.WitnessSignature(
				tx,
				txSigHash,
				i,
				vin.Value,
				txInPkScript,
				txscript.SigHashAll,
				privKey,
				true,
			)
			if err != nil {
				return fmt.Errorf("sign p2sh input error: %v", err)
			}
			tx.TxIn[i].Witness = witness

		case txscript.WitnessV0PubKeyHashTy:
			witness, err := txscript.WitnessSignature(
				tx,
				txSigHash,
				i,
				vin.Value,
				txInPkScript,
				txscript.SigHashAll,
				privKey,
				true,
			)
			if err != nil {
				return fmt.Errorf("sign segwit input error: %v", err)
			}
			tx.TxIn[i].Witness = witness

		case txscript.WitnessV1TaprootTy:
			witness, err := txscript.TaprootWitnessSignature(
				tx,
				txSigHash,
				i,
				vin.Value,
				txInPkScript,
				txscript.SigHashDefault,
				privKey,
			)
			if err != nil {
				return fmt.Errorf("sign taproot input error: %v", err)
			}
			tx.TxIn[i].Witness = witness

		default:
			return fmt.Errorf("unsupported script type: %s", scriptClass.String())
		}

		if vin.Value <= 0 {
			return fmt.Errorf("invalid input amount: %d", vin.Value)
		}

		vm, err := txscript.NewEngine(
			txInPkScript,
			tx,
			i,
			txscript.StandardVerifyFlags,
			nil,
			txSigHash,
			vin.Value,
			nil,
		)
		if err != nil {
			return fmt.Errorf("create script engine error: %w", err)
		}

		if err := vm.Execute(); err != nil {
			return fmt.Errorf("verify signature error: %w", err)
		}
	}

	return nil
}

// Helper functions that would need to be implemented:
func GetNetwork(networkType NetworkType) *chaincfg.Params {
	if networkType == MainNet {
		return &chaincfg.MainNetParams
	}
	return &chaincfg.TestNet3Params
}

func AddSufficientUTXOsForFee(tx *wire.MsgTx, btcUTXOs []UnspentOutput, feeRate int64, changeAddress, privateKey string) ([]ToSignInput, error) {
	// 创建找零地址
	addr, err := btcutil.DecodeAddress(changeAddress, &chaincfg.TestNet3Params)
	if err != nil {
		return nil, err
	}
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	// 计算当前费用
	currentFee := int64(tx.SerializeSize()) * feeRate
	toSignInputs := []ToSignInput{}

	// 选择UTXO直到满足费用
	for _, utxo := range btcUTXOs {
		if currentFee <= 0 {
			break
		}
		txHash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}

		// 添加输入
		prevOut := wire.NewOutPoint(txHash, utxo.Vout)
		tx.AddTxIn(wire.NewTxIn(prevOut, nil, nil))
		toSignInputs = append(toSignInputs, ToSignInput{
			TxID:            utxo.TxID,
			PublicKeyScript: utxo.ScriptPubKey,
			Vout:            utxo.Vout,
			Value:           utxo.Satoshis,
			PrivateKey:      privateKey,
		})

		currentFee -= utxo.Satoshis
	}

	// 添加找零输出
	if currentFee < 0 {
		tx.AddTxOut(&wire.TxOut{
			Value:    -currentFee,
			PkScript: pkScript,
		})
	}
	return toSignInputs, nil
}
