package pkg

type StringArray []string

func FilterString(_arr []string, fn func(t string) bool) []string {
	arr := []string{}
	for _, i := range _arr {
		if fn(i) {
			arr = append(arr, i)
		}
	}
	return arr
}

func (_arr *StringArray) Filter(fn func(t string) bool) *StringArray {
	arr := new(StringArray)
	for _, v := range *_arr {
		if fn(v) {
			*arr = append(*arr, v)
		}
	}
	return arr
}
