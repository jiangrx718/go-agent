package gins

// 过滤类型
type FilterType int

const (
	Intercept FilterType = 0 // 拦截
	Mark      FilterType = 1 // 标记
)
