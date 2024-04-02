package slicex

// MergeSlices 将多个相同类型的切片合并为一个切片
func MergeSlices[T any](slices ...[]T) []T {
	totalLen := 0
	for _, s := range slices {
		totalLen += len(s)
	}

	// 创建合并后的切片
	merged := make([]T, totalLen)
	i := 0

	// 将各个切片依次拷贝到合并后的切片中
	for _, s := range slices {
		i += copy(merged[i:], s)
	}

	return merged
}

// PadSlice 填充切片，将切片填充到指定长度
func PadSlice[T any](slice []T, length int, defaultValue ...T) []T {
	if len(slice) >= length {
		return slice
	}

	paddedSlice := make([]T, length)

	if len(defaultValue) > 0 {
		for i := 0; i < len(paddedSlice); i++ {
			paddedSlice[i] = defaultValue[0]
		}
	}

	copy(paddedSlice, slice)

	return paddedSlice
}
