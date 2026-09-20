package functions

func GetDefaultValue[T any](value *T) T { //nolint:ireturn
	if value == nil {
		var zero T

		return zero
	}

	return *value
}
