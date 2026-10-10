package sandbox

func prepareCreateRequest(request *CreateOptions) (*CreateOptions, error) {
	if request == nil {
		request = &CreateOptions{}
	}
	copy, err := request.Clone()
	if err != nil {
		return nil, NewError(ErrorInfo{
			Kind:      ErrorKindInvalidArgument,
			Operation: "create",
			Field:     CreateFieldProviderOptions,
			Message:   err.Error(),
		}, err)
	}
	if err := ValidateCreateOptions(copy); err != nil {
		return nil, err
	}
	return copy, nil
}
