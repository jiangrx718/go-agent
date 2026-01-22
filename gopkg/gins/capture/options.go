package gins

import (
	"context"
	"strings"
)

type Options struct {
	FilterMethods       []string                                     // 需要过滤的Method
	FilterPaths         []string                                     // 需要过滤的Path
	GetCustomFieldFuncs []func(ctx context.Context) (string, string) // 获取自定义字段函数
	MethodFilterType    FilterType                                   // 需要过滤的Method的操作类型
	PathFilterType      FilterType                                   // 需要过滤的Path的操作类型

	filterMethodMap map[string]bool
	filterPathMap   map[string]bool
}

func (o *Options) format() {
	if o.filterMethodMap == nil {
		o.filterMethodMap = make(map[string]bool)
	}

	if o.filterPathMap == nil {
		o.filterPathMap = make(map[string]bool)
	}

	for _, method := range o.FilterMethods {
		o.filterMethodMap[strings.ToUpper(method)] = true
	}

	for _, path := range o.FilterPaths {
		o.filterPathMap[path] = true
	}
}

func (o *Options) filterMethodExist(method string) bool {
	return o.filterMethodMap[method]
}

func (o *Options) filterPathExist(path string) bool {
	return o.filterPathMap[path]
}

func (o *Options) GetCustomField(ctx context.Context) CustomFields {
	customFields := make(CustomFields)
	for _, getCustomFieldFunc := range o.GetCustomFieldFuncs {
		key, value := getCustomFieldFunc(ctx)
		customFields[key] = value
	}

	return customFields
}
