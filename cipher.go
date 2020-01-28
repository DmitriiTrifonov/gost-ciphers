package main

type Cipher interface {
	SetKey()
	Encrypt(data [0x10]byte) []byte
	Decrypt(data [0x10]byte) []byte
	ResetKey()
}
