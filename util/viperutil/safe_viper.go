package viperutil

import (
    "errors"
    "io"
    "sync"
    "time"

    "github.com/spf13/viper"
)


// SafeViper 包装viper+读写锁，无需批量封装接口
//  1.不直接提供viper, 会破坏并发安全
//  2. 直接访问viper, viper.Sub() 会让并发失效, 所以不提供直接接口
//  3. 下面包装了viper常用接口, 特殊接口使用Read,Write访问
type SafeViper struct {
    mu sync.RWMutex
    v  *viper.Viper
}

func NewSafeViper() *SafeViper {
    return &SafeViper{
        v: viper.New(),
    }
}

// R 读锁块，读操作包裹
func (s *SafeViper) Read(fn func(v *viper.Viper)) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    fn(s.v)
}

// W 写锁块，修改/重载/Set 包裹
func (s *SafeViper) Write(fn func(v *viper.Viper)) {
    s.mu.Lock()
    defer s.mu.Unlock()
    fn(s.v)
}

//-----------------------------------
//   读
//-----------------------------------

func (s *SafeViper) Get(key string) any{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.Get(key)
}

func (s *SafeViper) GetString(key string) string{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetString(key)
}

func (s *SafeViper) GetBool(key string) bool{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetBool(key)
}

func (s *SafeViper) GetInt(key string) int{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetInt(key)
}

func (s *SafeViper) GetInt32(key string) int32{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetInt32(key)
}

func (s *SafeViper) GetInt64(key string) int64{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetInt64(key)
}

func (s *SafeViper) GetUint(key string) uint{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetUint(key)
}

func (s *SafeViper) GetUint16(key string) uint16{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetUint16(key)
}

func (s *SafeViper) GetUint32(key string) uint32{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetUint32(key)
}

func (s *SafeViper) GetUint64(key string) uint64{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetUint64(key)
}

func (s *SafeViper) GetFloat64(key string) float64{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetFloat64(key)
}

func (s *SafeViper) GetDuration(key string) time.Duration{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetDuration(key)
}

func (s *SafeViper) GetTime(key string) time.Time{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetTime(key)
}


func (s *SafeViper) GetIntSlice(key string) []int{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetIntSlice(key)
}

func (s *SafeViper) GetSizeInBytes(key string) uint{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetSizeInBytes(key)
}

func (s *SafeViper) GetStringMap(key string) map[string]any{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetStringMap(key)
}

func (s *SafeViper) GetStringSlice(key string) []string{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetStringSlice(key)
}

func (s *SafeViper) GetStringMapString(key string) map[string]string{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetStringMapString(key)
}

func (s *SafeViper) GetStringMapStringSlice(key string) map[string][]string{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.GetStringMapStringSlice(key)
}

func (s *SafeViper) Unmarshal(rawVal any, opts ...viper.DecoderConfigOption) error{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.Unmarshal(rawVal, opts...)
}

func (s *SafeViper) AllSettings() map[string]any{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.AllSettings()
}

func (s *SafeViper) IsSet(key string) bool{
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.v.IsSet(key)
}

func (s *SafeViper) DebugTo(w io.Writer){
    s.mu.RLock()
    defer s.mu.RUnlock()
    s.v.DebugTo(w)
}


func (s *SafeViper) SubUnmarshal(key string, rawVal any, opts ...viper.DecoderConfigOption) error{
    s.mu.RLock()
    defer s.mu.RUnlock()
    sub := s.v.Sub(key)
    if sub==nil{
        return errors.New("no sub found")
    }
    return sub.Unmarshal(rawVal, opts...)
}


//-----------------------------------
//   更新
//-----------------------------------

func (s *SafeViper) Set(key string, val any){
    s.mu.Lock()
    defer s.mu.Unlock()
    s.v.Set(key, val)
}

func (s *SafeViper) SetDefault(key string, val any){
    s.mu.Lock()
    defer s.mu.Unlock()
    s.v.SetDefault(key, val)
}

func (s *SafeViper) MergeConfigMap(cfg map[string]any) error{
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.v.MergeConfigMap(cfg)
}

func (s *SafeViper) MergeConfig(in io.Reader) error{
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.v.MergeConfig(in)
}
