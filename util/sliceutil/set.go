package sliceutil

//Compare 比较
func Compare[T comparable](a, b []T) map[T]int {
    m := make(map[T]int)
    for _, v := range a {
        m[v]++
    }
    for _, v := range b {
        m[v]--
    }
    return m
}

//Equal 是否相等
func Equal[T comparable ](a, b []T) bool{
    if len(a) != len(b) {
        return false
    }
    m := make(map[T]int)
    for _, v := range a {
        m[v]++
    }
    for _, v := range b {
        m[v]--
        if m[v] < 0 {
            return false
        }
    }
    for _, count := range m {
        if count != 0 {
            return false
        }
    }
    return true
}

//Difference 差集
func Difference[T comparable](a, b []T) []T{
    m := Compare[T](a, b)
    diff := make([]T, 0, len(m))
    for v, count := range m {
        if count != 0 {
            diff = append(diff, v)
        }
    }
    return diff
}

//Intersection 交集
func Intersection[T comparable](a, b []T) []T{
    m := Compare[T](a, b)
    equals := make([]T, 0, len(m))
    for v, count := range m {
        if count == 0 {
            equals = append(equals, v)
        }
    }
    return nil
}