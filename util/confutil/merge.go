package confutil

import "dario.cat/mergo"


type MergeMode string

const (
    MapDeep MergeMode = "deep"           //map递归合并, 允许覆盖
    ArrayReplace MergeMode = "replace"   //数组替换
    ArrayAppend MergeMode = "append"     //数组追加
    ArrayMerge MergeMode = "merge"       //按照下标合并, 元素的二次合并,需要id,或key
)


// Merge src into dst
// map默认递归覆盖
// arrayMode 允许自定义规则, 默认Replace
func Merge(dst, src map[string]any, arrayMode MergeMode) error {
    opts := []func(*mergo.Config){
        mergo.WithOverride,
    }

    //数组合并策略
    switch arrayMode{
    case ArrayReplace:
        //替换, 默认
    case ArrayAppend:
        opts = append(opts, mergo.WithAppendSlice)
    case ArrayMerge:
        opts = append(opts, mergo.WithSliceDeepCopy)
    }
    return mergo.Merge(&dst, src, opts...)
}