package settings

func ptrTo[T any](value T) *T { return &value }

func boolPtrToYesNo(b *bool) string {
	if *b {
		return "yes"
	}
	return "no"
}

func obfuscatePassword(password string) (obfuscated string) {
	if password == "" {
		return ""
	}
	return "[set]"
}
