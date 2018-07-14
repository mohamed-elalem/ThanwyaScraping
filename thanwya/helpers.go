package thanwya

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func branch(section int) string {
	var branchText string

	switch section {
	case 1:
		branchText = "علوم"
	case 2:
		branchText = "رياضة"
	case 5:
		branchText = "ادبى"
	default:
		branchText = "غير محدد"
	}
	return branchText
}
