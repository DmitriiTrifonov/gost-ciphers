package main

type Cipher interface {
	SetKey(data []byte)
	SetSubKeys()
	Encrypt(data []byte) []byte
	Decrypt(data []byte) []byte
	ResetSubKeys()
}
