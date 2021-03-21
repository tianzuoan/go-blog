package convert

import "strconv"

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

func (s StrTo) Int() (int, error) {
	return strconv.Atoi(s.String())
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

func (s StrTo) Uint32() (uint32, error) {
	value, err := strconv.Atoi(s.String())
	return uint32(value), err
}

func (s StrTo) MustUint32() uint32 {
	value, _ := s.Uint32()
	return value
}
