package str

import (
	"encoding/json"
	"strconv"
	"unsafe"
)

// S String type conversion
type S string

func NewWithByte(b []byte) S {
	return *(*S)(unsafe.Pointer(&b))
}

func (a S) String() string {
	return string(a)
}

// Bytes Convert to []byte
func (a S) Bytes() []byte {
	return *(*[]byte)(unsafe.Pointer(&a))
}

// Bool Convert to bool
func (a S) Bool() (bool, error) {
	b, err := strconv.ParseBool(a.String())
	if err != nil {
		return false, err
	}
	return b, nil
}

// DefaultBool Convert to bool, use default value if error occurs
func (a S) DefaultBool(defaultVal bool) bool {
	b, err := a.Bool()
	if err != nil {
		return defaultVal
	}
	return b
}

// Int64 Convert to int64
func (a S) Int64() (int64, error) {
	i, err := strconv.ParseInt(a.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// DefaultInt64 Convert to int64, use default value if error occurs
func (a S) DefaultInt64(defaultVal int64) int64 {
	i, err := a.Int64()
	if err != nil {
		return defaultVal
	}
	return i
}

// Int Convert to int
func (a S) Int() (int, error) {
	i, err := a.Int64()
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// DefaultInt Convert to int, use default value if error occurs
func (a S) DefaultInt(defaultVal int) int {
	i, err := a.Int()
	if err != nil {
		return defaultVal
	}
	return i
}

// Uint64 Convert to uint64
func (a S) Uint64() (uint64, error) {
	i, err := strconv.ParseUint(a.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// DefaultUint64 Convert to uint64, use default value if error occurs
func (a S) DefaultUint64(defaultVal uint64) uint64 {
	i, err := a.Uint64()
	if err != nil {
		return defaultVal
	}
	return i
}

// Uint Convert to uint
func (a S) Uint() (uint, error) {
	i, err := a.Uint64()
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}

// DefaultUint Convert to uint, use default value if error occurs
func (a S) DefaultUint(defaultVal uint) uint {
	i, err := a.Uint()
	if err != nil {
		return defaultVal
	}
	return uint(i)
}

// Float64 Convert to float64
func (a S) Float64() (float64, error) {
	f, err := strconv.ParseFloat(a.String(), 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

// DefaultFloat64 Convert to float64, use default value if error occurs
func (a S) DefaultFloat64(defaultVal float64) float64 {
	f, err := a.Float64()
	if err != nil {
		return defaultVal
	}
	return f
}

// Float32 Convert to float32
func (a S) Float32() (float32, error) {
	f, err := a.Float64()
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

// DefaultFloat32 Convert to float32, use default value if error occurs
func (a S) DefaultFloat32(defaultVal float32) float32 {
	f, err := a.Float32()
	if err != nil {
		return defaultVal
	}
	return f
}

// ToJSON Convert to JSON
func (a S) ToJSON(v interface{}) error {
	return json.Unmarshal(a.Bytes(), v)
}