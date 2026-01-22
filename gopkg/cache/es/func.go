package es

// QueryMap 查询条件数据结构
type QueryMap map[string]interface{}

// TextToKeyword Text类型子字段查询关键字
func TextToKeyword(field string) string {
	return field + ".keyword"
}

// TermQuery term查询条件封装
func TermQuery(field string, value interface{}) QueryMap {
	return QueryMap{
		"term": QueryMap{
			field: value,
		},
	}
}

// TermsQuery terms查询条件封装
func TermsQuery(field string, value []interface{}) QueryMap {
	return QueryMap{
		"terms": QueryMap{
			field: value,
		},
	}
}

// MustQueryMap bool中的must查询条件
type MustQueryMap []QueryMap

// ShouldQueryMap bool中的should查询条件
type ShouldQueryMap []QueryMap

// MustNotQueryMap bool中的must_not查询条件
type MustNotQueryMap []QueryMap

// BoolQueryParam BoolQuery方法函数
type BoolQueryParam struct {
	Must    MustQueryMap
	Should  ShouldQueryMap
	MustNot MustNotQueryMap
}

// BoolQuery bool查询条件封装
// func BoolQuery(must, should, mustNot []QueryMap) QueryMap {
func BoolQuery(param BoolQueryParam) QueryMap {

	// bool查询条件
	boolQueryMap := QueryMap{}

	// must条件
	if len(param.Must) > 0 {
		boolQueryMap["must"] = param.Must
	}

	// should条件
	if len(param.Should) > 0 {
		boolQueryMap["should"] = param.Should
	}

	// mustNot条件
	if len(param.MustNot) > 0 {
		boolQueryMap["must_not"] = param.MustNot
	}

	return QueryMap{
		"bool": boolQueryMap,
	}
}

// ConstantScoreParam ConstantScoreQuery参数
type ConstantScoreParam struct {
	Boost  *int
	Filter QueryMap
}

// ConstantScoreQuery constant_score查询条件封装
func ConstantScoreQuery(param ConstantScoreParam) QueryMap {

	query := QueryMap{}

	// 分数条件
	if param.Boost != nil {
		query["boost"] = *param.Boost
	}

	// 过滤条件
	if len(param.Filter) != 0 {
		query["filter"] = param.Filter
	}

	return QueryMap{
		"constant_score": query,
	}
}

// FilterQuery filter查询条件封装
func FilterQuery(value QueryMap) QueryMap {
	return QueryMap{
		"filter": value,
	}
}

// MatchQueryParam matchQuery参数
type MatchQueryParam struct {
	Field string
	Query interface{}
	Boost *int
}

// MatchQuery match查询条件封装
func MatchQuery(param MatchQueryParam) QueryMap {

	// 字段匹配参数
	fieldQueryMap := QueryMap{
		"query": param.Query,
	}

	// 过滤无效评分参数
	if param.Boost != nil {
		fieldQueryMap["boost"] = *param.Boost
	}

	// 拼装参数
	return QueryMap{
		"match": QueryMap{
			param.Field: fieldQueryMap,
		},
	}
}

// MatchPhraseQueryParam matchPhraseQuery参数
type MatchPhraseQueryParam struct {
	Field string
	Query interface{}
	Boost *int
}

// MatchPhraseQuery match_phrase查询条件封装
func MatchPhraseQuery(param MatchPhraseQueryParam) QueryMap {

	// 字段匹配参数
	fieldQueryMap := QueryMap{
		"query": param.Query,
	}

	// 过滤无效评分参数
	if param.Boost != nil {
		fieldQueryMap["boost"] = *param.Boost
	}

	// 拼装参数
	return QueryMap{
		"match_phrase": QueryMap{
			param.Field: fieldQueryMap,
		},
	}
}

// Query query查询条件封装
func Query(value QueryMap) QueryMap {
	return QueryMap{
		"query": value,
	}
}

// scoreMode 指定了该如何去合并每个文档生成的评分
type scoreMode string

const (

	// ScoreModeMultiply 函数结果相乘（默认）
	ScoreModeMultiply scoreMode = "multiply"

	// ScoreModeSum 函数结果相加
	ScoreModeSum scoreMode = "sum"

	// ScoreModeAvg 函数结果的平均值
	ScoreModeAvg scoreMode = "avg"

	// ScoreModeFirst 使用首个函数的结果做为最终结果
	ScoreModeFirst scoreMode = "first"

	// ScoreModeMax 函数结果的最大值
	ScoreModeMax scoreMode = "max"

	// ScoreModeMin	函数结果的最小值
	ScoreModeMin scoreMode = "min"
)

// boostMode 可以用来控制函数与查询评分_score合并后的结果
type boostMode string

const (

	// BoostModeMultiply 评分_score与函数值的乘积(默认)
	BoostModeMultiply boostMode = "multiply"

	// BoostModeReplace	评分_score会被忽略，仅使用函数值
	BoostModeReplace boostMode = "replace"

	// BoostModeSum	评分_score与函数值之和
	BoostModeSum boostMode = "sum"

	// BoostModeAvg	评分_score与函数值的平均值
	BoostModeAvg boostMode = "avg"

	// BoostModeMax	评分_score与函数值间的最大值
	BoostModeMax boostMode = "max"

	// BoostModeMin	评分_score与函数值间的最小值
	BoostModeMin boostMode = "min"
)

// FunctionScoreParam 评分函数参数
type FunctionScoreParam struct {

	// 过滤结果
	Query QueryMap

	// 评分函数
	Functions []QueryMap

	// 脚本评分函数允许计算自定义查询的评分，脚本表达式需使用文档中的数值字段。
	// 查询的分数将与脚本评分的结果相乘，如果不想使用这种方式，可通过设置"boost_mode":"replace"来禁止。
	ScriptScore QueryMap

	// 可以限制函数的最大效果，但是不会对最终的评分_score产生直接的影响。
	MaxBoost *int

	// 指定了该如何去合并每个文档生成的评分
	ScoreMode scoreMode

	// 以用来控制函数与查询评分_score合并后的结果
	BoostMode boostMode

	// 以设置为期望分数的阈值，能够排出不符合特定分数阈值的文档。
	MinScore *int
}

// FunctionScoreQuery 评分函数
func FunctionScoreQuery(param FunctionScoreParam) QueryMap {

	functionScoreQueryMap := QueryMap{}

	// 查询条件
	if len(param.Query) != 0 {
		functionScoreQueryMap["query"] = param.Query
	}

	// 评分方法
	if len(param.Functions) != 0 {
		functionScoreQueryMap["functions"] = param.Functions
	}

	// 评分脚本
	if len(param.ScriptScore) != 0 {
		functionScoreQueryMap["script_score"] = param.Functions
	}

	// max_boost
	if param.MaxBoost != nil {
		functionScoreQueryMap["max_boost"] = *param.MaxBoost
	}

	// score_mode
	if param.ScoreMode != "" {
		functionScoreQueryMap["score_mode"] = param.ScoreMode
	}

	// boost_mode
	if param.BoostMode != "" {
		functionScoreQueryMap["boost_mode"] = param.BoostMode
	}

	// min_score
	if param.MinScore != nil {
		functionScoreQueryMap["min_score"] = *param.MinScore
	}

	return QueryMap{
		"function_score": functionScoreQueryMap,
	}
}

// FunctionsQueryMap 评分过滤参数
type FunctionsQueryMap []QueryMap

// FunctionsQueryParam FunctionsQuery参数
type FunctionsQueryParam struct {
	Filter QueryMap
	Weight *int
}

// FunctionsQuery 评分参数
func FunctionsQuery(param FunctionsQueryParam) QueryMap {

	functionsQuery := QueryMap{}

	// 过滤条件
	if len(param.Filter) != 0 {
		functionsQuery["filter"] = param.Filter
	}

	// 权重
	if param.Weight != nil {
		functionsQuery["weight"] = *param.Weight
	}

	return functionsQuery
}

// RangeQueryParam RangeQuery方法参数
type RangeQueryParam struct {
	Field string
	Gt    *int
	Gte   *int
	Lt    *int
	Lte   *int
	Boost *int
}

// RangeQuery 范围查询条件
func RangeQuery(param RangeQueryParam) QueryMap {

	fieldQuery := QueryMap{}

	// 大于
	if param.Gt != nil {
		fieldQuery["gt"] = *param.Gt
	}

	// 大于等于
	if param.Gte != nil {
		fieldQuery["gte"] = *param.Gte
	}

	// 小于
	if param.Lt != nil {
		fieldQuery["lt"] = *param.Lt
	}

	// 小于等于
	if param.Lte != nil {
		fieldQuery["lte"] = *param.Lte
	}

	// 分数条件
	if param.Boost != nil {
		fieldQuery["boost"] = *param.Boost
	}

	return QueryMap{
		"range": QueryMap{
			param.Field: fieldQuery,
		},
	}
}

// ExistsQuery 字段存在查询
func ExistsQuery(filed string) QueryMap {
	return QueryMap{
		"exists": QueryMap{
			"field": filed,
		},
	}
}
