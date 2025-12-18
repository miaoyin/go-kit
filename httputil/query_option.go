package httputil

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"
)

// url parse

func NewQueryOption(q url.Values) *QueryOption {
	return &QueryOption{q: q}
}

type QueryOption struct {
	q url.Values
	err error
}

func (o *QueryOption) Has(name string) bool {
	return len(o.q[name]) > 0
}

func (o *QueryOption) String(name string) string {
	vs := o.q[name]
	if len(vs) == 0 {
		return ""
	}
	delete(o.q, name) // enable detection of unknown parameters
	return vs[len(vs)-1]
}

func (o *QueryOption) Strings(name string) []string {
	vs := o.q[name]
	delete(o.q, name)
	return vs
}

func (o *QueryOption) Int(name string) int {
	s := o.String(name)
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	}
	if o.err == nil {
		o.err = fmt.Errorf("url: invalid %s number: %s", name, err)
	}
	return 0
}

func (o *QueryOption) Duration(name string) time.Duration {
	s := o.String(name)
	if s == "" {
		return 0
	}
	// try plain number first
	if i, err := strconv.Atoi(s); err == nil {
		if i <= 0 {
			// disable timeouts
			return -1
		}
		return time.Duration(i) * time.Second
	}
	dur, err := time.ParseDuration(s)
	if err == nil {
		return dur
	}
	if o.err == nil {
		o.err = fmt.Errorf("url: invalid %s duration: %w", name, err)
	}
	return 0
}

func (o *QueryOption) Bool(name string) bool {
	switch s := o.String(name); s {
	case "true", "1":
		return true
	case "false", "0", "":
		return false
	default:
		if o.err == nil {
			o.err = fmt.Errorf("url: invalid %s boolean: expected true/false/1/0 or an empty string, got %q", name, s)
		}
		return false
	}
}

func (o *QueryOption) Remaining() []string {
	if len(o.q) == 0 {
		return nil
	}
	keys := make([]string, 0, len(o.q))
	for k := range o.q {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (o *QueryOption) Err() error {
	return o.err
}

