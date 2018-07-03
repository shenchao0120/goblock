package address

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"chaoshen.com/goblock/utils"
	"golang.org/x/crypto/ripemd160"
	"bytes"
)

const VERSION = byte(0x00)
const ADDRESS_CHECKSUM_LEN = 4

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	privKey, pubKey := newKeyPair()
	return &Wallet{privKey, pubKey}
}

func HashPublicKey(publicKey []byte) []byte {
	publicSHA256 := sha256.Sum256(publicKey)
	ripeMd := ripemd160.New()
	_, err := ripeMd.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	ripeRes := ripeMd.Sum(nil)
	return ripeRes
}

func newKeyPair() (*ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	return privKey, pubKey
}

func (w *Wallet) GetAddress() string {
	pubKeyHash := HashPublicKey(w.PublicKey)
	versionPubKeyHash := append([]byte{VERSION}, pubKeyHash...)
	checksum := checkSum(versionPubKeyHash)
	fullPayload := append(versionPubKeyHash, checksum...)
	return string(utils.Base58Encode(fullPayload))
}

func ValidateAddress (address string) bool {
	fullPayload:=utils.Base58Decode([]byte(address))
	checksum:=fullPayload[len(fullPayload)-ADDRESS_CHECKSUM_LEN:]
	pubKeyHash:=fullPayload[:len(fullPayload)-ADDRESS_CHECKSUM_LEN]
	targetChecksum:=checkSum(pubKeyHash)

	return bytes.Equal(checksum,targetChecksum)
}

func checkSum(input []byte) []byte {
	firstHash := sha256.Sum256(input)
	secondHash := sha256.Sum256(firstHash[:])
	if ADDRESS_CHECKSUM_LEN >= 32 {
		return secondHash[:]
	}
	return secondHash[:ADDRESS_CHECKSUM_LEN]
}
