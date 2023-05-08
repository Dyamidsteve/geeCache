package main

type ByteView struct {
	b []byte
}

func (b ByteView) Len() int {
	return len(b.b)
}

// ByteSlice returns a copy of the data as a byte slice.
func (b ByteView) ByteSlice() []byte {
	buf := make([]byte, b.Len())
	copy(buf, b.b)
	return buf
}

func (v ByteView) String() string {
	return string(v.b)
}
