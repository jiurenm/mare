package intx

// Pagination 用于计算分页的起始位置和结束位置
func Pagination[T Int](page, pageSize, length T) (start, end T) {
	if length <= 0 {
		return 0, 0
	}

	if pageSize <= 0 {
		return 0, 0
	}

	if page <= 0 {
		return 0, 0
	}

	if page == 1 && pageSize > length {
		return 0, length
	}

	// 计算总页数
	totalPage := (length + pageSize - 1) / pageSize
	if page > totalPage {
		return 0, 0
	}

	// 计算起始位置和结束位置
	start = (page - 1) * pageSize
	end = start + pageSize

	if end > length {
		end = length
	}

	return start, end
}
