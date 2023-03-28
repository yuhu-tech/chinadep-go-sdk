package ptrutil

func GetStrPtr(s string) *string {
	return &s
}

func GetIntPtr(i int) *int {
	return &i
}

func GetInt64Ptr(i int64) *int64 {
	return &i
}
