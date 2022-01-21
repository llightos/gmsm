package sm4_test

import (
	"bytes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/emmansun/gmsm/sm4"
)

func paddingPKCS7(buf []byte, blockSize int) []byte {
	bufLen := len(buf)
	padLen := blockSize - bufLen%blockSize
	padded := make([]byte, bufLen+padLen)
	copy(padded, buf)
	for i := 0; i < padLen; i++ {
		padded[bufLen+i] = byte(padLen)
	}
	return padded
}

func unpaddingPKCS7(padded []byte, size int) ([]byte, error) {
	if len(padded)%size != 0 {
		return nil, errors.New("pkcs7: Padded value wasn't in correct size")
	}
	paddedByte := int(padded[len(padded)-1])
	if (paddedByte > size) || (paddedByte < 1) {
		return nil, fmt.Errorf("Invalid decrypted text, no padding")
	}
	bufLen := len(padded) - paddedByte
	buf := make([]byte, bufLen)
	copy(buf, padded[:bufLen])
	return buf, nil
}

var cbcSM4Tests = []struct {
	name string
	key  []byte
	iv   []byte
	in   []byte
	out  []byte
}{
	{
		"from internet",
		[]byte("0123456789ABCDEF"),
		[]byte("0123456789ABCDEF"),
		[]byte("Hello World"),
		[]byte{0x0a, 0x67, 0x06, 0x2f, 0x0c, 0xd2, 0xdc, 0xe2, 0x6a, 0x7b, 0x97, 0x8e, 0xbf, 0x21, 0x34, 0xf9},
	},
	{
		"Three blocks",
		[]byte("0123456789ABCDEF"),
		[]byte("0123456789ABCDEF"),
		[]byte("Hello World Hello World Hello World Hello Worldd"),
		[]byte{
			0xd3, 0x1e, 0x36, 0x83, 0xe4, 0xfc, 0x9b, 0x51, 0x6a, 0x2c, 0x0f, 0x98, 0x36, 0x76, 0xa9, 0xeb,
			0x1f, 0xdc, 0xc3, 0x2a, 0xf3, 0x84, 0x08, 0x97, 0x81, 0x57, 0xa2, 0x06, 0x5d, 0xe3, 0x4c, 0x6a,
			0x06, 0x8d, 0x0f, 0xef, 0x4e, 0x2b, 0xfa, 0xb4, 0xbc, 0xab, 0xa6, 0x64, 0x41, 0xfd, 0xe0, 0xfe,
			0x92, 0xc1, 0x64, 0xec, 0xa1, 0x70, 0x24, 0x75, 0x72, 0xde, 0x12, 0x02, 0x95, 0x2e, 0xc7, 0x27,
		},
	},
	{
		"Four blocks",
		[]byte("0123456789ABCDEF"),
		[]byte("0123456789ABCDEF"),
		[]byte("Hello World Hello World Hello World Hello World Hello World Hell"),
		[]byte{
			0xd3, 0x1e, 0x36, 0x83, 0xe4, 0xfc, 0x9b, 0x51, 0x6a, 0x2c, 0x0f, 0x98, 0x36, 0x76, 0xa9, 0xeb,
			0x1f, 0xdc, 0xc3, 0x2a, 0xf3, 0x84, 0x08, 0x97, 0x81, 0x57, 0xa2, 0x06, 0x5d, 0xe3, 0x4c, 0x6a,
			0xe0, 0x02, 0xd6, 0xe4, 0xf5, 0x66, 0x87, 0xc4, 0xcc, 0x54, 0x1d, 0x1f, 0x1c, 0xc4, 0x2f, 0xe6,
			0xe5, 0x1d, 0xea, 0x52, 0xb8, 0x0c, 0xc8, 0xbe, 0xae, 0xcc, 0x44, 0xa8, 0x51, 0x81, 0x08, 0x60,
			0x34, 0x6e, 0x9d, 0xad, 0xe1, 0x8a, 0xf4, 0xa1, 0x83, 0x69, 0x57, 0xb9, 0x37, 0x26, 0x7e, 0x03,
		},
	},
	{
		"Five blocks",
		[]byte("0123456789ABCDEF"),
		[]byte("0123456789ABCDEF"),
		[]byte("Hello World Hello World Hello World Hello World Hello World Hello World Hello Wo"),
		[]byte{
			0xd3, 0x1e, 0x36, 0x83, 0xe4, 0xfc, 0x9b, 0x51, 0x6a, 0x2c, 0x0f, 0x98, 0x36, 0x76, 0xa9, 0xeb,
			0x1f, 0xdc, 0xc3, 0x2a, 0xf3, 0x84, 0x08, 0x97, 0x81, 0x57, 0xa2, 0x06, 0x5d, 0xe3, 0x4c, 0x6a,
			0xe0, 0x02, 0xd6, 0xe4, 0xf5, 0x66, 0x87, 0xc4, 0xcc, 0x54, 0x1d, 0x1f, 0x1c, 0xc4, 0x2f, 0xe6,
			0xe5, 0x1d, 0xea, 0x52, 0xb8, 0x0c, 0xc8, 0xbe, 0xae, 0xcc, 0x44, 0xa8, 0x51, 0x81, 0x08, 0x60,
			0xb6, 0x09, 0x7b, 0xb8, 0x7e, 0xdb, 0x53, 0x4b, 0xea, 0x2a, 0xc6, 0xa1, 0xe5, 0xa0, 0x2a, 0xe9,
			0x62, 0xb5, 0xe7, 0x50, 0x44, 0xea, 0x24, 0xcc, 0x9b, 0x5e, 0x07, 0x48, 0x04, 0x89, 0xa2, 0x74,
		},
	},
	{
		"9 blocks",
		[]byte("0123456789ABCDEF"),
		[]byte("0123456789ABCDEF"),
		[]byte("Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World"),
		[]byte{
			0xd3, 0x1e, 0x36, 0x83, 0xe4, 0xfc, 0x9b, 0x51, 0x6a, 0x2c, 0x0f, 0x98, 0x36, 0x76, 0xa9, 0xeb,
			0x1f, 0xdc, 0xc3, 0x2a, 0xf3, 0x84, 0x08, 0x97, 0x81, 0x57, 0xa2, 0x06, 0x5d, 0xe3, 0x4c, 0x6a,
			0xe0, 0x02, 0xd6, 0xe4, 0xf5, 0x66, 0x87, 0xc4, 0xcc, 0x54, 0x1d, 0x1f, 0x1c, 0xc4, 0x2f, 0xe6,
			0xe5, 0x1d, 0xea, 0x52, 0xb8, 0x0c, 0xc8, 0xbe, 0xae, 0xcc, 0x44, 0xa8, 0x51, 0x81, 0x08, 0x60,
			0xb6, 0x09, 0x7b, 0xb8, 0x7e, 0xdb, 0x53, 0x4b, 0xea, 0x2a, 0xc6, 0xa1, 0xe5, 0xa0, 0x2a, 0xe9,
			0x22, 0x65, 0x5b, 0xa3, 0xb9, 0xcc, 0x63, 0x92, 0x16, 0x0e, 0x2f, 0xf4, 0x3b, 0x93, 0x06, 0x82,
			0xb3, 0x8c, 0x26, 0x2e, 0x06, 0x51, 0x34, 0x2c, 0xe4, 0x3d, 0xd0, 0xc7, 0x2b, 0x8f, 0x31, 0x15,
			0x30, 0xa8, 0x96, 0x1c, 0xbc, 0x8e, 0xf7, 0x4f, 0x6b, 0x69, 0x9d, 0xc9, 0x40, 0x89, 0xd7, 0xe8,
			0xf7, 0x90, 0x47, 0x74, 0xaf, 0x40, 0xfd, 0x72, 0xc6, 0x17, 0xeb, 0xc0, 0x8b, 0x01, 0x71, 0x5c,
		},
	},
	{
		"17 blocks",
		[]byte("0123456789ABCDEF"),
		[]byte("0123456789ABCDEF"),
		[]byte("Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World"),
		[]byte{
			0xd3, 0x1e, 0x36, 0x83, 0xe4, 0xfc, 0x9b, 0x51, 0x6a, 0x2c, 0x0f, 0x98, 0x36, 0x76, 0xa9, 0xeb,
			0x1f, 0xdc, 0xc3, 0x2a, 0xf3, 0x84, 0x08, 0x97, 0x81, 0x57, 0xa2, 0x06, 0x5d, 0xe3, 0x4c, 0x6a,
			0xe0, 0x02, 0xd6, 0xe4, 0xf5, 0x66, 0x87, 0xc4, 0xcc, 0x54, 0x1d, 0x1f, 0x1c, 0xc4, 0x2f, 0xe6,
			0xe5, 0x1d, 0xea, 0x52, 0xb8, 0x0c, 0xc8, 0xbe, 0xae, 0xcc, 0x44, 0xa8, 0x51, 0x81, 0x08, 0x60,
			0xb6, 0x09, 0x7b, 0xb8, 0x7e, 0xdb, 0x53, 0x4b, 0xea, 0x2a, 0xc6, 0xa1, 0xe5, 0xa0, 0x2a, 0xe9,
			0x22, 0x65, 0x5b, 0xa3, 0xb9, 0xcc, 0x63, 0x92, 0x16, 0x0e, 0x2f, 0xf4, 0x3b, 0x93, 0x06, 0x82,
			0xb3, 0x8c, 0x26, 0x2e, 0x06, 0x51, 0x34, 0x2c, 0xe4, 0x3d, 0xd0, 0xc7, 0x2b, 0x8f, 0x31, 0x15,
			0x30, 0xa8, 0x96, 0x1c, 0xbc, 0x8e, 0xf7, 0x4f, 0x6b, 0x69, 0x9d, 0xc9, 0x40, 0x89, 0xd7, 0xe8,
			0x2a, 0xe8, 0xc3, 0x3d, 0xcb, 0x8a, 0x1c, 0xb3, 0x70, 0x7d, 0xe9, 0xe6, 0x88, 0x36, 0x65, 0x21,
			0x7b, 0x34, 0xac, 0x73, 0x8d, 0x4f, 0x11, 0xde, 0xd4, 0x21, 0x45, 0x9f, 0x1f, 0x3e, 0xe8, 0xcf,
			0x50, 0x92, 0x8c, 0xa4, 0x79, 0x58, 0x3a, 0x26, 0x01, 0x7b, 0x99, 0x5c, 0xff, 0x8d, 0x66, 0x5b,
			0x07, 0x86, 0x0e, 0x22, 0xb4, 0xb4, 0x83, 0x74, 0x33, 0x79, 0xd0, 0x54, 0x9f, 0x03, 0x6b, 0x60,
			0xa1, 0x52, 0x3c, 0x61, 0x1d, 0x91, 0xbf, 0x50, 0x00, 0xfb, 0x62, 0x58, 0xfa, 0xd3, 0xbd, 0x17,
			0x7d, 0x6f, 0xda, 0x76, 0x9a, 0xdb, 0x01, 0x96, 0x97, 0xc9, 0x5f, 0x64, 0x20, 0x3c, 0x70, 0x7a,
			0x40, 0x1f, 0x35, 0xc8, 0x22, 0xf2, 0x76, 0x6d, 0x8e, 0x4a, 0x78, 0xd7, 0x8d, 0x52, 0x51, 0x60,
			0x39, 0x14, 0xd8, 0xcd, 0xc7, 0x4b, 0x3f, 0xb3, 0x16, 0xdf, 0x52, 0xba, 0xcb, 0x98, 0x56, 0xaa,
			0x97, 0x8b, 0xab, 0xa7, 0xbf, 0xe8, 0x0f, 0x16, 0x27, 0xbb, 0x56, 0xce, 0x10, 0xe5, 0x90, 0x05,
		},
	},
	{
		"A.1",
		[]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10},
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10},
		[]byte{
			0x68, 0x1e, 0xdf, 0x34, 0xd2, 0x06, 0x96, 0x5e, 0x86, 0xb3, 0xe9, 0x4f, 0x53, 0x6e, 0x42, 0x46,
			0x67, 0x7d, 0x30, 0x7e, 0x84, 0x4d, 0x7a, 0xa2, 0x45, 0x79, 0xd5, 0x56, 0x49, 0x0d, 0xc7, 0xaa},
	},
}

func TestCBCEncrypterSM4(t *testing.T) {
	for _, test := range cbcSM4Tests {
		c, err := sm4.NewCipher(test.key)
		if err != nil {
			t.Errorf("%s: NewCipher(%d bytes) = %s", test.name, len(test.key), err)
			continue
		}

		encrypter := cipher.NewCBCEncrypter(c, test.iv)

		plainText := paddingPKCS7(test.in, sm4.BlockSize)
		data := make([]byte, len(plainText))
		copy(data, plainText)

		encrypter.CryptBlocks(data, data)
		if !bytes.Equal(test.out, data) {
			t.Errorf("%s: CBCEncrypter\nhave %s\nwant %x", test.name, hex.EncodeToString(data), test.out)
			for i := 0; i < len(data); i++ {
				fmt.Printf("0x%02x, ", data[i])
				if (i+1)%16 == 0 {
					fmt.Println()
				}
			}
		}
	}
}

func TestCBCDecrypterSM4(t *testing.T) {
	for _, test := range cbcSM4Tests {
		c, err := sm4.NewCipher(test.key)
		if err != nil {
			t.Errorf("%s: NewCipher(%d bytes) = %s", test.name, len(test.key), err)
			continue
		}

		decrypter := cipher.NewCBCDecrypter(c, test.iv)

		data := make([]byte, len(test.out))
		copy(data, test.out)

		decrypter.CryptBlocks(data, data)
		data, err = unpaddingPKCS7(data, sm4.BlockSize)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(test.in, data) {
			t.Errorf("%s: CBCDecrypter\nhave %x\nwant %x", test.name, data, test.in)
		}
	}
}
