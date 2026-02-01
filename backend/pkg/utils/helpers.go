package utils

// StringToPtr converts a string to *string, returning nil for empty strings
func StringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
