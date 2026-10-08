// Code generated from YAML validation specifications. DO NOT EDIT.
package sandbox

import (
	fmt "fmt"
	v10 "github.com/go-playground/validator/v10"
	math "math"
	strings "strings"
)

func newCreateOptionsValidator() *v10.Validate {
	v := v10.New(v10.WithRequiredStructEnabled())
	_ = v.RegisterValidation("finite", func(fl v10.FieldLevel) bool {
		value := fl.Field().Float()
		return !math.IsNaN(value) && !math.IsInf(value, 0)
	})
	_ = v.RegisterValidation("nonblank", func(fl v10.FieldLevel) bool { return strings.TrimSpace(fl.Field().String()) != "" })
	v.RegisterStructValidation(func(sl v10.StructLevel) {
		x := sl.Current().Interface().(AutomaticAction)
		_ = x
		count0 := 0
		if x.After != nil {
			count0++
		}
		if (x.GetMode() == 2) && (count0 == 0) {
			sl.ReportError(x, "after", "after", "requires_any", "")
		}
		count1 := 0
		if x.After != nil {
			count1++
		}
		if (x.GetMode() != 2) && (count1 > 0) {
			sl.ReportError(x, "after", "after", "forbids", "")
		}
	}, AutomaticAction{})
	v.RegisterStructValidation(func(sl v10.StructLevel) {
		x := sl.Current().Interface().(NetworkConfig)
		_ = x
		count0 := 0
		if x.OutboundCIDRs != nil {
			count0++
		}
		if x.OutboundDomains != nil {
			count0++
		}
		if (x.GetEgress() == 3) && (count0 == 0) {
			sl.ReportError(x, "outbound_cidrs,outbound_domains", "outbound_cidrs,outbound_domains", "requires_any", "")
		}
		count1 := 0
		if x.OutboundCIDRs != nil {
			count1++
		}
		if x.OutboundDomains != nil {
			count1++
		}
		if x.GetPrivateNetwork() {
			count1++
		}
		if x.LinkedSandbox != nil {
			count1++
		}
		if x.InboundCIDRs != nil {
			count1++
		}
		if len(x.Ports) > 0 {
			count1++
		}
		if (x.GetEgress() == 2) && (count1 > 0) {
			sl.ReportError(x, "outbound_cidrs,outbound_domains,private_network,linked_sandbox,inbound_cidrs,ports", "outbound_cidrs,outbound_domains,private_network,linked_sandbox,inbound_cidrs,ports", "forbids", "")
		}
		count2 := 0
		if x.OutboundCIDRs != nil {
			count2++
		}
		if x.OutboundDomains != nil {
			count2++
		}
		if (x.GetEgress() == 1) && (count2 > 0) {
			sl.ReportError(x, "outbound_cidrs,outbound_domains", "outbound_cidrs,outbound_domains", "forbids", "")
		}
	}, NetworkConfig{})
	v.RegisterStructValidation(func(sl v10.StructLevel) {
		x := sl.Current().Interface().(ReadinessProbe)
		_ = x
		count0 := 0
		if x.TCPPort != nil {
			count0++
		}
		if x.Command != nil {
			count0++
		}
		if (true) && (count0 != 1) {
			sl.ReportError(x, "tcp_port,command", "tcp_port,command", "exactly_one", "")
		}
	}, ReadinessProbe{})
	v.RegisterStructValidation(func(sl v10.StructLevel) {
		x := sl.Current().Interface().(Resources)
		_ = x
		if (true) && (x.GetCPULimitCores() > 0 && x.GetCPULimitCores() < x.GetCPUCores()) {
			sl.ReportError(x, "cpu_limit_cores", "cpu_limit_cores", "limit", "")
		}
		if (true) && (x.GetMemoryLimitMiB() > 0 && x.GetMemoryLimitMiB() < x.GetMemoryMiB()) {
			sl.ReportError(x, "memory_limit_mib", "memory_limit_mib", "limit", "")
		}
	}, Resources{})
	v.RegisterStructValidation(func(sl v10.StructLevel) {
		x := sl.Current().Interface().(SandboxSource)
		_ = x
		count0 := 0
		if x.ProviderDefault != nil {
			count0++
		}
		if x.Image != nil {
			count0++
		}
		if x.Snapshot != nil {
			count0++
		}
		if x.WarmPool != nil {
			count0++
		}
		if (true) && (count0 > 1) {
			sl.ReportError(x, "provider_default,image,snapshot,warm_pool", "provider_default,image,snapshot,warm_pool", "at_most_one", "")
		}
	}, SandboxSource{})
	return v
}

// ValidateCreateOptions applies the generated shared rules without filling defaults.
func ValidateCreateOptions(request *CreateOptions) error {
	if request == nil {
		return nil
	}
	if err := newCreateOptionsValidator().Struct(request); err != nil {
		return fmt.Errorf("sandbox-kit: invalid configuration: %w", err)
	}
	return nil
}
