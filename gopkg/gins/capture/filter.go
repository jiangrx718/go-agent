package gins

type filterResult struct {
	isFilterMethodHitAlready bool // 是否已命中过滤Methods
	isFilterPathHitAlready   bool // 是否已命中过滤Paths
}

func newFilterResult(options Options, request *Request) *filterResult {
	isFilterMethodHitAlready := isIntercept(request.Method, options.MethodFilterType, options.filterMethodExist)
	if isFilterMethodHitAlready == nil {
		return nil
	}

	isFilterPathHitAlready := isIntercept(request.Path, options.PathFilterType, options.filterPathExist)
	if isFilterPathHitAlready == nil {
		return nil
	}

	return &filterResult{
		isFilterMethodHitAlready: *isFilterMethodHitAlready,
		isFilterPathHitAlready:   *isFilterPathHitAlready,
	}
}

func (c *filterResult) GetIsFilterMethodHitAlready() bool {
	return c.isFilterMethodHitAlready
}

func (c *filterResult) GetIsFilterPathHitAlready() bool {
	return c.isFilterPathHitAlready
}

func isIntercept(keyWord string, filterType FilterType, filterExistFunc func(keyWord string) bool) *bool {
	var interceptBool bool
	if !filterExistFunc(keyWord) {
		return &interceptBool
	}

	if filterType == Intercept {
		return nil
	}

	interceptBool = true
	return &interceptBool
}
