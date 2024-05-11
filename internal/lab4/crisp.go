package lab4

import (
	"bytes"
	"crypto/hmac"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gost-ling/internal/lab1"
	"gost-ling/internal/lab2"
	"gost-ling/internal/lab3"
	"hash"
	"runtime"
)

var (
	ExternalKeyIdFlagWithVersion = []byte{ // 1 bit + 15 bits
		0x00, 0x00,
	}
	CS = []byte{ // 8 bits
		0xf1,
	}
	KeyId = []byte{ // 8 bits
		0x80,
	}
)

const (
	BlockSize  = 16 // bytes
	KeySize    = 32 // bytes
	PacketSize = 56 // bytes
)

type Crisp struct {
	Decoder    *Decoder
	Encoder    *Encoder
	randomSeed [32]byte
}

func (c *Crisp) Close() {
	for i := 0; i < len(c.randomSeed); i++ {
		c.randomSeed[i] = 0x00
	}
	c.Decoder.kdf.Close()
	c.Decoder.seqNum = 0
	c.Decoder.cipher.Close()
	c.Encoder.kdf.Close()
	c.Encoder.seqNum = 0
	c.Encoder.cipher.Close()
	runtime.GC()
	fmt.Printf("Clear mem [Crisp]: %p\n", &c)
}

type Decoder struct {
	random *lab3.XoroShiro256StarStar
	kdf    *lab2.KDF
	cipher *lab1.Ofb
	seqNum uint32
}

type Encoder struct {
	random *lab3.XoroShiro256StarStar
	kdf    *lab2.KDF
	cipher *lab1.Ofb
	seqNum uint32
}

type Message struct {
	ExternalKeyIdFlagWithVersion []byte
	CS                           []byte
	KeyId                        []byte
	SeqNum                       []byte
	Payload                      []byte
	ICV                          []byte
	Digits                       []byte
}

func New(key []byte, randomSeed [32]byte) *Crisp {
	if len(key) != KeySize {
		panic("Key size should be 32 bytes")
	}

	kdf := lab2.NewKDF(key[:])
	cipher := lab1.NewOfb(key)
	return &Crisp{
		Decoder: &Decoder{
			random: lab3.New(randomSeed[:]),
			kdf:    kdf,
			cipher: cipher,
			seqNum: 0,
		},
		Encoder: &Encoder{
			random: lab3.New(randomSeed[:]),
			kdf:    kdf,
			cipher: cipher,
			seqNum: 0,
		},
		randomSeed: randomSeed,
	}
}

func (c *Crisp) Reset() {
	c.Encoder.seqNum = 0
	c.Encoder.random = lab3.New(c.randomSeed[:])
	c.Decoder.seqNum = 0
	c.Decoder.random = lab3.New(c.randomSeed[:])
}

func (c *Crisp) Encode(plainText []byte) []Message {
	var res []Message

	c.Reset()
	for i := 0; i < len(plainText); i += BlockSize {
		message := c.EncodeNextBlock(plainText[i : i+BlockSize])
		res = append(res, message)
	}

	return res
}

func (c *Crisp) EncodeNextBlock(plainText []byte) Message {
	if len(plainText) != BlockSize {
		panic("Block size should be 16 bytes")
	}
	e := c.Encoder

	var seqNum [4]byte
	var seed [8]byte
	var iv []byte
	binary.BigEndian.PutUint32(seqNum[:], e.seqNum)
	binary.BigEndian.PutUint64(seed[:], e.random.Next())
	for i := 0; i < 2; i++ {
		var buf [8]byte
		binary.BigEndian.PutUint64(buf[:], e.random.Next())
		iv = append(iv, buf[:]...)
	}

	// Key(N)
	key := e.kdf.Derive(seqNum[:], seed[:], 1)

	// text[N]
	block := plainText[:]

	// Payload(N)
	var ciphertext []byte
	for _, blck := range e.cipher.Encrypt(block, iv) {
		ciphertext = append(ciphertext, blck[:]...)
	}

	// Mac(N)
	h := hmac.New(newHash, key)
	h.Write(ciphertext)
	mac := h.Sum(nil)

	var message []byte
	message = append(message, ExternalKeyIdFlagWithVersion...)
	message = append(message, CS...)
	message = append(message, KeyId...)
	message = append(message, seqNum[:]...)
	message = append(message, ciphertext[:]...)
	message = append(message, mac...)

	e.seqNum += 1 // complete current iteration and prepare next
	return Message{
		ExternalKeyIdFlagWithVersion: ExternalKeyIdFlagWithVersion,
		CS:                           CS,
		KeyId:                        KeyId,
		SeqNum:                       seqNum[:],
		Payload:                      ciphertext[:],
		ICV:                          mac[:],
		Digits:                       message,
	}
}

func (c *Crisp) Decode(cipherText [][]byte) [][]byte {
	for i, b := range cipherText {
		if len(b) != PacketSize {
			panic(fmt.Sprintf("Block size of block [%d] should be 56 bytes", i))
		}
	}

	var res [][]byte
	for _, b := range cipherText {
		decoded := c.DecodeNextBlock(b)
		res = append(res, decoded)
	}

	return res
}

func (c *Crisp) DecodeNextBlock(cipherText []byte) []byte {
	if len(cipherText) != PacketSize {
		panic("Block size should be equal 56 bytes")
	}
	d := c.Decoder

	var seqNum [4]byte
	var seed [8]byte
	var iv []byte
	binary.BigEndian.PutUint64(seed[:], d.random.Next())
	for i := 0; i < 2; i++ {
		var buf [8]byte
		binary.BigEndian.PutUint64(buf[:], d.random.Next())
		iv = append(iv, buf[:]...)
	}

	// parse seqnum and payload
	seqNum = [4]byte(cipherText[4:8])
	payload := cipherText[8:24]

	// get iter key
	key := d.kdf.Derive(seqNum[:], seed[:], 1)

	// check mac
	mac := cipherText[24:56]
	h := hmac.New(newHash, key)
	h.Write(payload)
	checkedMac := h.Sum(nil)
	if !bytes.Equal(mac, checkedMac) {
		panic("mac is not equal")
	}

	// decrypt
	decryptBlock := d.cipher.Decrypt(payload, iv)
	var decryptText []byte
	for _, block := range decryptBlock {
		decryptText = append(decryptText, block[:]...)
	}

	return decryptText
}

func (m *Message) String() string {
	format :=
		`Message:
    ExternalKeyIdFlagWithVersion: %s
    CS:                           %s
    KeyId:                        %s
    SeqNum:                       %s
    Payload:                      %s
    ICV:                          %s
    As block:                     %s`

	return fmt.Sprintf(format,
		hex.EncodeToString(m.ExternalKeyIdFlagWithVersion),
		hex.EncodeToString(m.CS),
		hex.EncodeToString(m.KeyId),
		hex.EncodeToString(m.SeqNum),
		hex.EncodeToString(m.Payload),
		hex.EncodeToString(m.ICV),
		hex.EncodeToString(m.Digits))
}

func newHash() hash.Hash {
	return lab2.NewHash()
}
