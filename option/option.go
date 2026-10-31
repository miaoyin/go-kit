package option

// Option 选项函数
type Option[T any] func(*T)

// Apply 应用函数
func Apply[T any](target *T, opts ...Option[T]) {
    for _, opt := range opts {
        opt(target)
    }
}
