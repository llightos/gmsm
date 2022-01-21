package cipher

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"testing"
)

// These test vectors have been taken from IEEE P1619/D16, Annex B.
var xtsTestVectors = []struct {
	key        string
	sector     uint64
	plaintext  string
	ciphertext string
}{
	{ // XTS-AES-128 applied for a data unit of 32 bytes
		"0000000000000000000000000000000000000000000000000000000000000000",
		0,
		"0000000000000000000000000000000000000000000000000000000000000000",
		"917cf69ebd68b2ec9b9fe9a3eadda692cd43d2f59598ed858c02c2652fbf922e",
	}, {
		"1111111111111111111111111111111122222222222222222222222222222222",
		0x3333333333,
		"4444444444444444444444444444444444444444444444444444444444444444",
		"c454185e6a16936e39334038acef838bfb186fff7480adc4289382ecd6d394f0",
	}, {
		"fffefdfcfbfaf9f8f7f6f5f4f3f2f1f022222222222222222222222222222222",
		0x3333333333,
		"4444444444444444444444444444444444444444444444444444444444444444",
		"af85336b597afc1a900b2eb21ec949d292df4c047e0b21532186a5971a227a89",
	}, { // XTS-AES-128 applied for a data unit of 512 bytes
		"2718281828459045235360287471352631415926535897932384626433832795",
		0,
		"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedfe0e1e2e3e4e5e6e7e8e9eaebecedeeeff0f1f2f3f4f5f6f7f8f9fafbfcfdfeff000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedfe0e1e2e3e4e5e6e7e8e9eaebecedeeeff0f1f2f3f4f5f6f7f8f9fafbfcfdfeff",
		"27a7479befa1d476489f308cd4cfa6e2a96e4bbe3208ff25287dd3819616e89cc78cf7f5e543445f8333d8fa7f56000005279fa5d8b5e4ad40e736ddb4d35412328063fd2aab53e5ea1e0a9f332500a5df9487d07a5c92cc512c8866c7e860ce93fdf166a24912b422976146ae20ce846bb7dc9ba94a767aaef20c0d61ad02655ea92dc4c4e41a8952c651d33174be51a10c421110e6d81588ede82103a252d8a750e8768defffed9122810aaeb99f9172af82b604dc4b8e51bcb08235a6f4341332e4ca60482a4ba1a03b3e65008fc5da76b70bf1690db4eae29c5f1badd03c5ccf2a55d705ddcd86d449511ceb7ec30bf12b1fa35b913f9f747a8afd1b130e94bff94effd01a91735ca1726acd0b197c4e5b03393697e126826fb6bbde8ecc1e08298516e2c9ed03ff3c1b7860f6de76d4cecd94c8119855ef5297ca67e9f3e7ff72b1e99785ca0a7e7720c5b36dc6d72cac9574c8cbbc2f801e23e56fd344b07f22154beba0f08ce8891e643ed995c94d9a69c9f1b5f499027a78572aeebd74d20cc39881c213ee770b1010e4bea718846977ae119f7a023ab58cca0ad752afe656bb3c17256a9f6e9bf19fdd5a38fc82bbe872c5539edb609ef4f79c203ebb140f2e583cb2ad15b4aa5b655016a8449277dbd477ef2c8d6c017db738b18deb4a427d1923ce3ff262735779a418f20a282df920147beabe421ee5319d0568",
	}, { // Vector 5
		"2718281828459045235360287471352631415926535897932384626433832795",
		1,
		"27a7479befa1d476489f308cd4cfa6e2a96e4bbe3208ff25287dd3819616e89cc78cf7f5e543445f8333d8fa7f56000005279fa5d8b5e4ad40e736ddb4d35412328063fd2aab53e5ea1e0a9f332500a5df9487d07a5c92cc512c8866c7e860ce93fdf166a24912b422976146ae20ce846bb7dc9ba94a767aaef20c0d61ad02655ea92dc4c4e41a8952c651d33174be51a10c421110e6d81588ede82103a252d8a750e8768defffed9122810aaeb99f9172af82b604dc4b8e51bcb08235a6f4341332e4ca60482a4ba1a03b3e65008fc5da76b70bf1690db4eae29c5f1badd03c5ccf2a55d705ddcd86d449511ceb7ec30bf12b1fa35b913f9f747a8afd1b130e94bff94effd01a91735ca1726acd0b197c4e5b03393697e126826fb6bbde8ecc1e08298516e2c9ed03ff3c1b7860f6de76d4cecd94c8119855ef5297ca67e9f3e7ff72b1e99785ca0a7e7720c5b36dc6d72cac9574c8cbbc2f801e23e56fd344b07f22154beba0f08ce8891e643ed995c94d9a69c9f1b5f499027a78572aeebd74d20cc39881c213ee770b1010e4bea718846977ae119f7a023ab58cca0ad752afe656bb3c17256a9f6e9bf19fdd5a38fc82bbe872c5539edb609ef4f79c203ebb140f2e583cb2ad15b4aa5b655016a8449277dbd477ef2c8d6c017db738b18deb4a427d1923ce3ff262735779a418f20a282df920147beabe421ee5319d0568",
		"264d3ca8512194fec312c8c9891f279fefdd608d0c027b60483a3fa811d65ee59d52d9e40ec5672d81532b38b6b089ce951f0f9c35590b8b978d175213f329bb1c2fd30f2f7f30492a61a532a79f51d36f5e31a7c9a12c286082ff7d2394d18f783e1a8e72c722caaaa52d8f065657d2631fd25bfd8e5baad6e527d763517501c68c5edc3cdd55435c532d7125c8614deed9adaa3acade5888b87bef641c4c994c8091b5bcd387f3963fb5bc37aa922fbfe3df4e5b915e6eb514717bdd2a74079a5073f5c4bfd46adf7d282e7a393a52579d11a028da4d9cd9c77124f9648ee383b1ac763930e7162a8d37f350b2f74b8472cf09902063c6b32e8c2d9290cefbd7346d1c779a0df50edcde4531da07b099c638e83a755944df2aef1aa31752fd323dcb710fb4bfbb9d22b925bc3577e1b8949e729a90bbafeacf7f7879e7b1147e28ba0bae940db795a61b15ecf4df8db07b824bb062802cc98a9545bb2aaeed77cb3fc6db15dcd7d80d7d5bc406c4970a3478ada8899b329198eb61c193fb6275aa8ca340344a75a862aebe92eee1ce032fd950b47d7704a3876923b4ad62844bf4a09c4dbe8b4397184b7471360c9564880aedddb9baa4af2e75394b08cd32ff479c57a07d3eab5d54de5f9738b8d27f27a9f0ab11799d7b7ffefb2704c95c6ad12c39f1e867a4b7b1d7818a4b753dfd2a89ccb45e001a03a867b187f225dd",
	}, { // Vector 10, XTS-AES-256 applied for a data unit of 512 bytes
		"27182818284590452353602874713526624977572470936999595749669676273141592653589793238462643383279502884197169399375105820974944592",
		0xff,
		"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedfe0e1e2e3e4e5e6e7e8e9eaebecedeeeff0f1f2f3f4f5f6f7f8f9fafbfcfdfeff000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedfe0e1e2e3e4e5e6e7e8e9eaebecedeeeff0f1f2f3f4f5f6f7f8f9fafbfcfdfeff",
		"1c3b3a102f770386e4836c99e370cf9bea00803f5e482357a4ae12d414a3e63b5d31e276f8fe4a8d66b317f9ac683f44680a86ac35adfc3345befecb4bb188fd5776926c49a3095eb108fd1098baec70aaa66999a72a82f27d848b21d4a741b0c5cd4d5fff9dac89aeba122961d03a757123e9870f8acf1000020887891429ca2a3e7a7d7df7b10355165c8b9a6d0a7de8b062c4500dc4cd120c0f7418dae3d0b5781c34803fa75421c790dfe1de1834f280d7667b327f6c8cd7557e12ac3a0f93ec05c52e0493ef31a12d3d9260f79a289d6a379bc70c50841473d1a8cc81ec583e9645e07b8d9670655ba5bbcfecc6dc3966380ad8fecb17b6ba02469a020a84e18e8f84252070c13e9f1f289be54fbc481457778f616015e1327a02b140f1505eb309326d68378f8374595c849d84f4c333ec4423885143cb47bd71c5edae9be69a2ffeceb1bec9de244fbe15992b11b77c040f12bd8f6a975a44a0f90c29a9abc3d4d893927284c58754cce294529f8614dcd2aba991925fedc4ae74ffac6e333b93eb4aff0479da9a410e4450e0dd7ae4c6e2910900575da401fc07059f645e8b7e9bfdef33943054ff84011493c27b3429eaedb4ed5376441a77ed43851ad77f16f541dfd269d50d6a5f14fb0aab1cbb4c1550be97f7ab4066193c4caa773dad38014bd2092fa755c824bb5e54c4f36ffda9fcea70b9c6e693e148c151",
	}, { // XTS-AES-128 applied for a data unit that is not a multiple of 16 bytes, but should be a complte byte
		"c46acc2e7e013cb71cdbf750cf76b000249fbf4fb6cd17607773c23ffa2c4330",
		94,
		"7e9c2289cba460e470222953439cdaa892a5433d4dab2a3f67",
		"9af624641d42b036377ef37b4a158f49e49f6ee308ad449ecf",
	}, {
		"56ffcc9bbbdf413f0fc0f888f44b7493bb1925a39b8adf02d9009bb16db0a887",
		144,
		"9a839cc14363bafcfc0cc93b14f8e769d35b94cc98267438e3",
		"fbd8dcfc4d662259f48b151728c3b37233a35127a77051ee9d",
	},
	{
		"7454a43b87b1cf0dec95032c22873be3cace3bb795568854c1a008c07c5813f3",
		108,
		"41088fa15195b2733fe824d2c1fdc8306080863945fb2a73cf",
		"f916d877f817ae390f42dc54723bda0ad3ba5f331a72d05ccb",
	},
}

func fromHex(s string) []byte {
	ret, err := hex.DecodeString(s)
	if err != nil {
		panic("xts: invalid hex in test")
	}
	return ret
}

func TestXTS(t *testing.T) {
	for i, test := range xtsTestVectors {
		c, err := NewXTS(aes.NewCipher, fromHex(test.key))
		if err != nil {
			t.Errorf("#%d: failed to create cipher: %s", i, err)
			continue
		}
		plaintext := fromHex(test.plaintext)
		ciphertext := make([]byte, len(plaintext))

		c.Encrypt(ciphertext, plaintext, test.sector)
		expectedCiphertext := fromHex(test.ciphertext)
		if !bytes.Equal(ciphertext, expectedCiphertext) {
			t.Errorf("#%d: encrypted failed, got: %x, want: %x", i, ciphertext, expectedCiphertext)
			continue
		}

		decrypted := make([]byte, len(ciphertext))
		c.Decrypt(decrypted, ciphertext, test.sector)
		if !bytes.Equal(decrypted, plaintext) {
			t.Errorf("#%d: decryption failed, got: %x, want: %x", i, decrypted, plaintext)
		}
	}
}

func TestShorterCiphertext(t *testing.T) {
	c, err := NewXTS(aes.NewCipher, make([]byte, 32))
	if err != nil {
		t.Fatalf("NewCipher failed: %s", err)
	}

	plaintext := make([]byte, 32)
	encrypted := make([]byte, 48)
	decrypted := make([]byte, 48)

	c.Encrypt(encrypted, plaintext, 0)
	c.Decrypt(decrypted, encrypted[:len(plaintext)], 0)

	if !bytes.Equal(plaintext, decrypted[:len(plaintext)]) {
		t.Errorf("En/Decryption is not inverse")
	}
}
