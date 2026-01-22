package gorms

type Paging[T any] struct {
	Total int `json:"-"` // 总数
	List  []T `json:"-"` // 数据列表
}

type Page struct {
	PageIndex int `json:"page_index" form:"page_index"` // 页码
	PageSize  int `json:"page_size" form:"page_size"`   // 数组列表长度
}

type FindByPageFunc[T any] func(offset int, limit int) ([]*T, int64, error)

func FindByPage[T any](findByPage FindByPageFunc[T], page Page) ([]*T, int, error) {
	ts, total, err := findByPage(ComputeOffsetLimit(page))
	return ts, int(total), err
}

func ComputeOffsetLimit(page Page) (int, int) {
	if page.PageIndex < 1 {
		page.PageIndex = 1
	}
	if page.PageSize < 1 {
		page.PageSize = 100
	}
	return (page.PageIndex - 1) * page.PageSize, page.PageSize
}

func PaginationQuery[T any](findByPage FindByPageFunc[T], page Page) (*Paging[*T], error) {
	list, total, err := FindByPage(findByPage, page)
	if err != nil {
		return &Paging[*T]{}, err
	}

	return &Paging[*T]{
		Total: total,
		List:  list,
	}, nil
}
