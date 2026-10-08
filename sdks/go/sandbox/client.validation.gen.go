// Code generated from YAML validation specifications. DO NOT EDIT.
package sandbox

import (
	fmt "fmt"
	v10 "github.com/go-playground/validator/v10"
	math "math"
	strings "strings"
)

func newConfigValidator() *v10.Validate {
	v := v10.New(v10.WithRequiredStructEnabled())
	_ = v.RegisterValidation("finite", func(fl v10.FieldLevel) bool {
		value := fl.Field().Float()
		return !math.IsNaN(value) && !math.IsInf(value, 0)
	})
	_ = v.RegisterValidation("nonblank", func(fl v10.FieldLevel) bool { return strings.TrimSpace(fl.Field().String()) != "" })
	v.RegisterStructValidation(func(sl v10.StructLevel) {
		x := sl.Current().Interface().(AuthConfig)
		_ = x
		count0 := 0
		if x.APIKey != nil {
			count0++
		}
		if x.TokenPair != nil {
			count0++
		}
		if x.BearerToken != nil {
			count0++
		}
		if x.OAuthRefresh != nil {
			count0++
		}
		if (true) && (count0 != 1) {
			sl.ReportError(x, "api_key,token_pair,bearer_token,oauth_refresh", "api_key,token_pair,bearer_token,oauth_refresh", "exactly_one", "")
		}
	}, AuthConfig{})
	v.RegisterStructValidation(func(sl v10.StructLevel) {
		x := sl.Current().Interface().(OAuthCredentials)
		_ = x
		count0 := 0
		if x.ClientSecret != nil {
			count0++
		}
		if x.JWTKey != nil {
			count0++
		}
		if (true) && (count0 != 1) {
			sl.ReportError(x, "client_secret,jwt_key", "client_secret,jwt_key", "exactly_one", "")
		}
	}, OAuthCredentials{})
	return v
}

// ValidateConfig applies the generated shared rules without filling defaults.
func ValidateConfig(request *Config) error {
	if request == nil {
		return nil
	}
	if err := newConfigValidator().Struct(request); err != nil {
		return fmt.Errorf("sandbox-kit: invalid configuration: %w", err)
	}
	return nil
}
