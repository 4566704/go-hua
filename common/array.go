package common

func DeleteElement[S []T, T any](slice S, index int) S {
	if index < 0 || index >= len(slice) {
		return slice
	}
	// 计算要删除的元素后面的元素需要向前移动多少位置
	n := len(slice) - 1
	copy(slice[index:n], slice[index+1:n+1])
	return slice[:n]
}
