package redfox

func conditionalDefault(condition bool, defaultValue interface{}) interface{} {
	if !condition {
		return nil
	}

	return defaultValue
}
