package paging

type Page struct {
	PageIndex int `json:"page_index" form:"page_index"` // 页码
	PageSize  int `json:"page_size" form:"page_size"`   // 数组列表长度
}

type Paging[T any] struct {
	Total int `json:"total"` // 总数
	List  []T `json:"list"`  // 数据列表
}

func NewPaging[T any](total int, list []T) *Paging[T] {
	return &Paging[T]{
		Total: total,
		List:  list,
	}
}
