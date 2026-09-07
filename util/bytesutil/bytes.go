package bytesutil


//CopyByte 复制
func CopyByte(b []byte) []byte {
    if b == nil {
        return nil
    }
    cp := make([]byte, len(b))
    copy(cp, b)
    return cp
}

//AppendByte 追加一个字节, 不修改原始内容
func AppendByte(orig []byte, b byte) []byte {
    n := len(orig)
    out := make([]byte, n+1)
    copy(out, orig)
    out[n] = b
    return out
}

//TruncateLastByte 去掉最后一个字节, 不修改原始内容
func TruncateLastByte(orig []byte) []byte {
    n := len(orig)
    if n == 0 {
        return nil
    }
    out := make([]byte, n-1)
    copy(out, orig)
    return out
}

