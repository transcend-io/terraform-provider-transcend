package transcend

func toStringList(origs []interface{}) []string {
	vals := make([]string, len(origs))
	for i, orig := range origs {
		vals[i] = orig.(string)
	}

	return vals
}
