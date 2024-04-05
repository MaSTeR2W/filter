package filter

import "strconv"

func sl_of_int_to_sl_of_str[T int | int8 | int16 | int32 | int64](s []T) []string {
	var str_sl = make([]string, 0, len(s))
	switch t_s := any(s).(type) {
	case []int:
		for _, el := range t_s {
			str_sl = append(str_sl, strconv.FormatInt(int64(el), 10))
		}
	case []int8:
		for _, el := range t_s {
			str_sl = append(str_sl, strconv.FormatInt(int64(el), 10))
		}
	case []int16:
		for _, el := range t_s {
			str_sl = append(str_sl, strconv.FormatInt(int64(el), 10))
		}
	case []int32:
		for _, el := range t_s {
			str_sl = append(str_sl, strconv.FormatInt(int64(el), 10))
		}
	case []int64:
		for _, el := range t_s {
			str_sl = append(str_sl, strconv.FormatInt(el, 10))
		}
	}
	return str_sl
}
