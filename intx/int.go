package intx

type Int interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr
}

// Min 返回两个值中的最小值
func Min[T Int](x, y T) T {
	if x < y {
		return x
	}

	return y
}

// Max 返回两个值中的最大值
func Max[T Int](x, y T) T {
	if x > y {
		return x
	}

	return y
}

// Abs 返回输入值的绝对值
func Abs[T Int](input T) T {
	if input < 0 {
		return -input
	}

	return input
}
