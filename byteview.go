package cache

// ByteView 存储真实的缓存值 --- byte 类型是为了能够支持任意的数据类型的存储
type ByteView struct {
	b []byte
}

// Len 所占的内存大小
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 拷贝 b ，防止缓存值被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String 以字符串形式返回数据，并在必要时进行拷贝。
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
