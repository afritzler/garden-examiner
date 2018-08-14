package util

func StringValue(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
