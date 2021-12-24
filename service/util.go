package service

func maxInt64(i, j int64) int64 {
	if i > j {
		return i
	}
	return j
}

func minInt64(i, j int64) int64 {
	if i < j {
		return i
	}
	return j
}
