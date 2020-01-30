package magma

type Magma struct {
	key [0x20]byte
	subKeys [0x20][]byte
}

const t32  = 4294967296 // 2^32

var pi = [8][16]byte{
	{1, 7, 14, 13, 0, 5, 8, 3, 4, 15, 10, 6, 9, 12, 11, 2},
	{8, 14, 2, 5, 6, 9, 1, 12, 15, 4, 11, 0, 13, 10, 3, 7},
	{5, 13, 15, 6, 9, 2, 12, 10, 11, 7, 8, 1, 4, 3, 14, 0},
	{7, 15, 5, 10, 8, 1, 6, 13, 0, 9, 3, 14, 11, 4, 2, 12},
	{12, 8, 2, 1, 13, 4, 15, 6, 7, 0, 10, 5, 3, 14, 9, 11},
	{11, 3, 5, 8, 2, 15, 10, 13, 14, 1, 7, 4, 12, 9, 6, 0},
	{6, 8, 2, 3, 9, 10, 5, 12, 1, 14, 4, 7, 11, 13, 0, 15},
	{12, 4, 6, 2, 10, 5, 11, 9, 14, 8, 13, 7, 0, 3, 15, 1},
}

func xor(a []byte, b []byte) []byte {
	c := make([]byte, 4)
	for i := 0; i < 4; i++ {
		c[i] = a[i] ^ b[i]
	}
	return c
}

func convertToUInt32(a []byte) uint32 {
	var r uint32
	for i := 0; i < 3; i++ {
		r |= uint32(a[i])
		r <<= 8
	}
	r |= uint32(a[3])
	return r
}

func add32(a, b uint32) uint32 {
	return uint32(int(a + b) % t32)
}

func convertToArray(a uint32) []byte {
	arr := make([]byte, 4)
	arr[3] = byte(a)
	arr[2] = byte(a >> 8)
	arr[1] = byte(a >> 16)
	arr[0] = byte(a >> 24)
	return arr
}

func x32(a []byte, b []byte) []byte {
	c := make([]byte, 4)
	var internal int
	for i := 3; i >= 0; i-- {
		internal = int(a[i]) + int(b[i]) + (internal >> 8)
		c[i] = byte(internal & 0xFF)
	}
	return c
}

// Splits bytes in array by two 4 bits numbers and changes value from pi table
func t(input []byte) []byte {
	out := make([]byte, 4)
	var fbp, sbp byte
	for i := 0; i < 4; i++ {
		fbp = (input[i] & 0xF0) >> 4
		sbp = input[i] & 0x0F
		fbp = pi[i * 2][fbp]
		sbp = pi[i * 2 + 1][sbp]
		out[i] = (fbp << 4) | sbp
	}
	return out
}

func (m *Magma) SetSubKeys() {
	m.subKeys[0] = m.key[:4]
	m.subKeys[1] = m.key[4:8]
	m.subKeys[2] = m.key[8:12]
	m.subKeys[3] = m.key[12:16]
	m.subKeys[4] = m.key[16:20]
	m.subKeys[5] = m.key[20:24]
	m.subKeys[6] = m.key[24:28]
	m.subKeys[7] = m.key[28:]
	m.subKeys[8] = m.key[:4]
	m.subKeys[9] = m.key[4:8]
	m.subKeys[10] = m.key[8:12]
	m.subKeys[11] = m.key[12:16]
	m.subKeys[12] = m.key[16:20]
	m.subKeys[13] = m.key[20:24]
	m.subKeys[14] = m.key[24:28]
	m.subKeys[15] = m.key[28:]
	m.subKeys[16] = m.key[:4]
	m.subKeys[17] = m.key[4:8]
	m.subKeys[18] = m.key[8:12]
	m.subKeys[19] = m.key[12:16]
	m.subKeys[20] = m.key[16:20]
	m.subKeys[21] = m.key[20:24]
	m.subKeys[22] = m.key[24:28]
	m.subKeys[23] = m.key[28:]
	m.subKeys[24] = m.key[28:]
	m.subKeys[25] = m.key[24:28]
	m.subKeys[26] = m.key[20:24]
	m.subKeys[27] = m.key[16:20]
	m.subKeys[28] = m.key[12:16]
	m.subKeys[29] = m.key[8:12]
	m.subKeys[30] = m.key[4:8]
	m.subKeys[31] = m.key[:4]
}

func gSwap(block, key []byte) []byte {
	out32 := x32(block, key)
	outT := t(out32)
	n := convertToUInt32(outT)
	n = (n << 11) | (n >> 21)
	return convertToArray(n)
}

func gIter(block, key []byte) []byte {
	rh := make([]byte, 4)
	lh := make([]byte, 4)
	G := make([]byte, 4)
	out := make([]byte, 8)

	for i := 0; i < 4; i++ {
		rh[i] = block[4 + i]
		lh[i] = block[i]
	}

	G = gSwap(key, rh)
	G = xor(G,lh)

	for i := 0; i < 4; i++ {
		lh[i] = rh[i]
		rh[i] = G[i]
	}

	for i := 0; i < 4; i++ {
		out[i] = lh[i]
		out[4 + i] = rh[i]
	}
	return out
}

func gFinal(block, key []byte) []byte {
	rh := make([]byte, 4)
	lh := make([]byte, 4)
	G := make([]byte, 4)
	out := make([]byte, 8)

	for i := 0; i < 4; i++ {
		rh[i] = block[4 + i]
		lh[i] = block[i]
	}

	G = gSwap(key, rh)
	G = xor(G,lh)

	for i := 0; i < 4; i++ {
		lh[i] = G[i]
	}

	for i := 0; i < 4; i++ {
		out[i] = lh[i]
		out[4 + i] = rh[i]
	}
	return out
}

func (m *Magma) Encrypt(data []byte) []byte {
	out := make([]byte, 8)
	out = gIter(data, m.subKeys[0])

	for i := 1; i < 31; i++ {
		out = gIter(out, m.subKeys[i])
	}

	out = gFinal(out, m.subKeys[31])

	return out
}

func (m *Magma) Decrypt(data []byte) []byte {
	out := make([]byte, 8)
	out = gIter(data, m.subKeys[31])

	for i := 30; i > 0; i-- {
		out = gIter(out, m.subKeys[i])
	}

	out = gFinal(out, m.subKeys[0])

	return out
}

func (m *Magma) SetKey(data []byte) {
	var arr [0x20]byte
	copy(arr[:], data[:0x20])
	m.key = arr
}