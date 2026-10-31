package viperx

import (
	"errors"
	"io"
	"time"

	"github.com/miaoyin/go-kit/atomicx"
	"github.com/spf13/viper"
)

type viperSnapshot struct {
	vp        *viper.Viper
	version   uint64
	changedAt time.Time
}

type Updater interface {
	Swap(newVal *viper.Viper) (old *viper.Viper)
}

// ViperReader 包装viper+读写锁，无需批量封装接口
//
//	1.不直接提供viper, 会破坏并发安全
//	2. 直接访问viper, viper.Sub() 会让并发失效, 所以不提供直接接口
//	3. 下面包装了viper常用接口, 特殊接口使用Read,Write访问
type ViperReader struct {
	holder *atomicx.Holder[viper.Viper]
}

func NewViperReader(holder *atomicx.Holder[viper.Viper]) *ViperReader {
	return &ViperReader{holder: holder}
}

func (s *ViperReader) GetSnapshot() (*viper.Viper, error) {
	return s.holder.GetSnapshot()
}

//-----------------------------------
//   读
//-----------------------------------

func (s *ViperReader) Get(key string) any {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return nil
	}
	return vp.Get(key)
}

func (s *ViperReader) GetString(key string) string {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return ""
	}
	return vp.GetString(key)
}

func (s *ViperReader) GetBool(key string) bool {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return false
	}
	return vp.GetBool(key)
}

func (s *ViperReader) GetInt(key string) int {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return 0
	}
	return vp.GetInt(key)
}

func (s *ViperReader) GetInt32(key string) int32 {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return 0
	}
	return vp.GetInt32(key)
}

func (s *ViperReader) GetInt64(key string) int64 {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return 0
	}
	return vp.GetInt64(key)
}

func (s *ViperReader) GetUint(key string) uint {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return 0
	}
	return vp.GetUint(key)
}

func (s *ViperReader) GetUint16(key string) uint16 {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return 0
	}
	return vp.GetUint16(key)
}

func (s *ViperReader) GetUint32(key string) uint32 {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return 0
	}
	return vp.GetUint32(key)
}

func (s *ViperReader) GetUint64(key string) uint64 {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return 0
	}
	return vp.GetUint64(key)
}

func (s *ViperReader) GetFloat64(key string) float64 {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return 0
	}
	return vp.GetFloat64(key)
}

func (s *ViperReader) GetDuration(key string) time.Duration {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return 0
	}
	return vp.GetDuration(key)
}

func (s *ViperReader) GetTime(key string) time.Time {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return time.Time{}
	}
	return vp.GetTime(key)
}

func (s *ViperReader) GetIntSlice(key string) []int {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return nil
	}
	return vp.GetIntSlice(key)
}

func (s *ViperReader) GetSizeInBytes(key string) uint {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return 0
	}
	return vp.GetSizeInBytes(key)
}

func (s *ViperReader) GetStringMap(key string) map[string]any {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return nil
	}
	return vp.GetStringMap(key)
}

func (s *ViperReader) GetStringSlice(key string) []string {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return nil
	}
	return vp.GetStringSlice(key)
}

func (s *ViperReader) GetStringMapString(key string) map[string]string {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return nil
	}
	return vp.GetStringMapString(key)
}

func (s *ViperReader) GetStringMapStringSlice(key string) map[string][]string {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return nil
	}
	return vp.GetStringMapStringSlice(key)
}

func (s *ViperReader) Unmarshal(rawVal any, opts ...viper.DecoderConfigOption) error {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return err
	}
	return vp.Unmarshal(rawVal, opts...)
}

func (s *ViperReader) AllSettings() map[string]any {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return map[string]any{}
	}
	return vp.AllSettings()
}

func (s *ViperReader) IsSet(key string) bool {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return false
	}
	return vp.IsSet(key)
}

func (s *ViperReader) DebugTo(w io.Writer) {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return
	}
	vp.DebugTo(w)
}

func (s *ViperReader) SubUnmarshal(key string, rawVal any, opts ...viper.DecoderConfigOption) error {
	vp, err := s.holder.GetSnapshot()
	if err != nil {
		return err
	}
	sub := vp.Sub(key)
	if sub == nil {
		return errors.New("no sub found")
	}
	return sub.Unmarshal(rawVal, opts...)
}
