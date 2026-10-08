package sandbox

func prepareCreateRequest(request *CreateOptions) (*CreateOptions, error) {
	if request == nil {
		request = &CreateOptions{}
	}
	copy, err := cloneData(request)
	if err != nil {
		return nil, err
	}
	if err := ValidateCreateOptions(copy); err != nil {
		return nil, err
	}
	return copy, nil
}
