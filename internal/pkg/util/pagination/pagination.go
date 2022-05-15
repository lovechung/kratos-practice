package pagination

func GetOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}

func GetPage(page, pageSize *int32) (int, int) {
	var p, ps int
	if page == nil {
		p = 1
	} else {
		p = int(*page)
	}
	if pageSize == nil {
		ps = 5
	} else {
		ps = int(*pageSize)
	}
	return p, ps
}
