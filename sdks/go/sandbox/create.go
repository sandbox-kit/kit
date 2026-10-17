package sandbox

// prepareCreateRequest copies the request and validates it before the provider sees it.
// A nil request becomes an empty request, which means provider defaults.
// The caller can keep using the original value after this returns.
func prepareCreateRequest(request *CreateOptions) (*CreateOptions, error) {
	if request == nil {
		request = &CreateOptions{}
	}
	copy := CloneCreateOptions(request)
	if err := ValidateCreateOptions(copy); err != nil {
		return nil, err
	}
	return copy, nil
}
