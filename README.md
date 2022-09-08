
# GM-Standards SM2/SM3/SM4/SM9/ZUC for Go

![Travis (.org)](https://img.shields.io/travis/emmansun/gmsm?label=arm64)
[![Github CI](https://github.com/emmansun/gmsm/actions/workflows/ci.yml/badge.svg)](https://github.com/emmansun/gmsm/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/emmansun/gmsm/branch/main/graph/badge.svg?token=Otdi8m8sFj)](https://codecov.io/gh/emmansun/gmsm)
[![Go Report Card](https://goreportcard.com/badge/github.com/emmansun/gmsm)](https://goreportcard.com/report/github.com/emmansun/gmsm)
[![Documentation](https://godoc.org/github.com/emmansun/gmsm?status.svg)](https://godoc.org/github.com/emmansun/gmsm)
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/emmansun/gmsm)
[![Release](https://img.shields.io/github/release/emmansun/gmsm/all.svg)](https://github.com/emmansun/gmsm/releases)

## Packages
* **SM2** - This is a SM2 sm2p256v1 implementation whose performance is similar like golang native NIST P256 under **amd64** and **arm64**, for implementation detail, please refer [SM2实现细节](https://github.com/emmansun/gmsm/wiki/SM2%E6%80%A7%E8%83%BD%E4%BC%98%E5%8C%96). It supports ShangMi sm2 digital signature, public key encryption algorithm and also key exchange.

* **SM3** - This is also a SM3 implementation whose performance is similar like golang native SHA 256 with SIMD under **amd64**, for implementation detail, please refer [SM3性能优化](https://github.com/emmansun/gmsm/wiki/SM3%E6%80%A7%E8%83%BD%E4%BC%98%E5%8C%96). It also provides A64 cryptographic instructions SM3 POC without test.

* **SM4** - For SM4 implementation, SIMD & AES-NI are used under **amd64** and **arm64**, for detail please refer [SM4性能优化](https://github.com/emmansun/gmsm/wiki/SM4%E6%80%A7%E8%83%BD%E4%BC%98%E5%8C%96), it supports CBC/CFB/OFB/CTR/GCM/CCM/XTS modes. It also provides A64 cryptographic instructions SM4 POC without test.

* **SM9** - For SM9 implementation, please reference [sm9/bn256 README.md](https://github.com/emmansun/gmsm/tree/main/sm9/bn256).

* **ZUC** - For ZUC implementation, SIMD, AES-NI and CLMUL are used under **amd64** and **arm64**, for detail please refer [Efficient Software Implementations of ZUC](https://github.com/emmansun/gmsm/wiki/Efficient-Software-Implementations-of-ZUC)

* **CIPHER** - CCM/XTS cipher modes.

* **SMX509** - a fork of golang X509 that supports ShangMi.

* **PKCS8** - a fork of [youmark/pkcs8](https://github.com/youmark/pkcs8) that supports ShangMi.

* **ECDH** - a similar implementation of golang ECDH that supports SM2 ECDH & SM2MQV without usage of **big.Int**, a replacement of SM2 key exchange. For detail, pleaes refer [is my code constant time?](https://github.com/emmansun/gmsm/wiki/is-my-code-constant-time%3F)

## Some Related Projects
* **PKCS12** - https://github.com/emmansun/go-pkcs12
* **PKCS7** - https://github.com/emmansun/pkcs7
* **TLCP** - https://github.com/Trisia/gotlcp