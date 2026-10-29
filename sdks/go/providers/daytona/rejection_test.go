package daytona

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"go.yaml.in/yaml/v3"
)

func TestSharedRejectionSpec(t *testing.T) {
	ran := 0
	for _, tc := range loadRejectionSpec(t) {
		if tc.Provider != "daytona" {
			continue
		}
		ran++
		t.Run(tc.Name, func(t *testing.T) {
			var err error
			switch tc.Operation {
			case "create":
				_, err = planCreate(rejectionCreate(t, tc.Create))
			default:
				t.Fatalf("operation %q", tc.Operation)
			}
			var detail *sandbox.Error
			if !errors.As(err, &detail) || detail.Kind != rejectionKind(t, tc.Kind) || detail.Field != tc.Field {
				t.Fatalf("error = %v detail = %#v", err, detail)
			}
		})
	}
	if ran == 0 {
		t.Fatal("no daytona cases")
	}
}

type rejectionCase struct {
	Name      string             `yaml:"name"`
	Provider  string             `yaml:"provider"`
	Operation string             `yaml:"operation"`
	Kind      string             `yaml:"kind"`
	Field     string             `yaml:"field"`
	Create    *rejectionCreateIn `yaml:"create"`
	Config    *rejectionConfigIn `yaml:"config"`
}

type rejectionCreateIn struct {
	Source        *rejectionSource        `yaml:"source"`
	Runtime       *rejectionRuntime       `yaml:"runtime"`
	Isolation     *rejectionIsolation     `yaml:"isolation"`
	Resources     *rejectionResources     `yaml:"resources"`
	Placement     *rejectionPlacement     `yaml:"placement"`
	Security      *rejectionSecurity      `yaml:"security"`
	Network       *rejectionNetwork       `yaml:"network"`
	Lifetime      *rejectionLifetime      `yaml:"lifetime"`
	Provisioning  *rejectionProvisioning  `yaml:"provisioning"`
	Observability *rejectionObservability `yaml:"observability"`
}

type rejectionSource struct {
	Image *struct {
		Reference string `yaml:"reference"`
	} `yaml:"image"`
}

type rejectionRuntime struct {
	Language *string `yaml:"language"`
	Version  *string `yaml:"version"`
}

type rejectionIsolation struct {
	Kind                 *string `yaml:"kind"`
	NestedVirtualization *bool   `yaml:"nested_virtualization"`
}

type rejectionResources struct {
	CPULimitCores  *float64 `yaml:"cpu_limit_cores"`
	MemoryLimitMiB *uint64  `yaml:"memory_limit_mib"`
	DiskMiB        *uint64  `yaml:"disk_mib"`
}

type rejectionPlacement struct {
	Regions []string `yaml:"regions"`
	Target  *string  `yaml:"target"`
}

type rejectionSecurity struct {
	PublicAccess     *bool `yaml:"public_access"`
	WorkloadIdentity *bool `yaml:"workload_identity"`
}

type rejectionNetwork struct {
	PrivateNetwork *bool   `yaml:"private_network"`
	LinkedSandbox  *string `yaml:"linked_sandbox"`
}

type rejectionLifetime struct {
	MaximumLifetimeMs *int64 `yaml:"maximum_lifetime_ms"`
	Ephemeral         *bool  `yaml:"ephemeral"`
	IdleTerminate     *struct {
		Mode string `yaml:"mode"`
	} `yaml:"idle_terminate"`
}

type rejectionProvisioning struct {
	TimeoutMs      *int64 `yaml:"timeout_ms"`
	QueueTimeoutMs *int64 `yaml:"queue_timeout_ms"`
}

type rejectionObservability struct {
	Verbose           *bool   `yaml:"verbose"`
	TelemetryEndpoint *string `yaml:"telemetry_endpoint"`
}

type rejectionConfigIn struct {
	Scope *struct {
		OrganizationID *string `yaml:"organization_id"`
		ProjectID      *string `yaml:"project_id"`
	} `yaml:"scope"`
}

type rejectionFile struct {
	Version int             `yaml:"version"`
	Cases   []rejectionCase `yaml:"cases"`
}

func loadRejectionSpec(t *testing.T) []rejectionCase {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("spec path")
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(file), "../../../../specs/rejections.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	var spec rejectionFile
	if err := dec.Decode(&spec); err != nil {
		t.Fatal(err)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		t.Fatalf("extra document: %v", err)
	}
	if spec.Version != 1 {
		t.Fatalf("version %d", spec.Version)
	}
	return spec.Cases
}

func rejectionKind(t *testing.T, kind string) sandbox.ErrorKind {
	t.Helper()
	switch kind {
	case "unsupported":
		return sandbox.ErrorKindUnsupported
	case "invalid_argument":
		return sandbox.ErrorKindInvalidArgument
	default:
		t.Fatalf("kind %q", kind)
		return 0
	}
}

func rejectionCreate(t *testing.T, in *rejectionCreateIn) *sandbox.CreateOptions {
	t.Helper()
	if in == nil {
		t.Fatal("create request is missing")
	}
	request := &sandbox.CreateOptions{}
	if in.Source != nil && in.Source.Image != nil {
		request.Source = &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: in.Source.Image.Reference}}}
	}
	if in.Runtime != nil {
		request.Runtime = &sandbox.RuntimeConfig{}
		if in.Runtime.Language != nil {
			request.Runtime.Language = sandbox.Value(*in.Runtime.Language)
		}
		if in.Runtime.Version != nil {
			request.Runtime.Version = sandbox.Value(*in.Runtime.Version)
		}
	}
	if in.Isolation != nil {
		request.Isolation = &sandbox.IsolationConfig{}
		if in.Isolation.Kind != nil {
			switch *in.Isolation.Kind {
			case "container":
				request.Isolation.Kind = sandbox.IsolationKindContainer
			default:
				t.Fatalf("isolation kind %q", *in.Isolation.Kind)
			}
		}
		if in.Isolation.NestedVirtualization != nil {
			request.Isolation.NestedVirtualization = sandbox.Value(*in.Isolation.NestedVirtualization)
		}
	}
	if in.Resources != nil {
		request.Resources = &sandbox.Resources{}
		if in.Resources.CPULimitCores != nil {
			request.Resources.CpuLimitCores = sandbox.Value(*in.Resources.CPULimitCores)
		}
		if in.Resources.MemoryLimitMiB != nil {
			request.Resources.MemoryLimitMib = sandbox.Value(*in.Resources.MemoryLimitMiB)
		}
		if in.Resources.DiskMiB != nil {
			request.Resources.DiskMib = sandbox.Value(*in.Resources.DiskMiB)
		}
	}
	if in.Placement != nil {
		request.Placement = &sandbox.Placement{Regions: in.Placement.Regions}
		if in.Placement.Target != nil {
			request.Placement.Target = sandbox.Value(*in.Placement.Target)
		}
	}
	if in.Security != nil {
		request.Security = &sandbox.SecurityConfig{}
		if in.Security.PublicAccess != nil {
			request.Security.PublicAccess = sandbox.Value(*in.Security.PublicAccess)
		}
		if in.Security.WorkloadIdentity != nil {
			request.Security.WorkloadIdentity = sandbox.Value(*in.Security.WorkloadIdentity)
		}
	}
	if in.Network != nil {
		request.Network = &sandbox.NetworkConfig{}
		if in.Network.PrivateNetwork != nil {
			request.Network.PrivateNetwork = sandbox.Value(*in.Network.PrivateNetwork)
		}
		if in.Network.LinkedSandbox != nil {
			request.Network.LinkedSandbox = sandbox.Value(*in.Network.LinkedSandbox)
		}
	}
	if in.Lifetime != nil {
		request.Lifetime = &sandbox.LifetimePolicy{}
		if in.Lifetime.MaximumLifetimeMs != nil {
			request.Lifetime.MaximumLifetime = sandbox.Duration(time.Duration(*in.Lifetime.MaximumLifetimeMs) * time.Millisecond)
		}
		if in.Lifetime.Ephemeral != nil {
			request.Lifetime.Ephemeral = sandbox.Value(*in.Lifetime.Ephemeral)
		}
		if in.Lifetime.IdleTerminate != nil {
			switch in.Lifetime.IdleTerminate.Mode {
			case "disabled":
				request.Lifetime.IdleTerminate = &sandbox.AutomaticAction{Mode: sandbox.PolicyModeDisabled}
			default:
				t.Fatalf("idle terminate mode %q", in.Lifetime.IdleTerminate.Mode)
			}
		}
	}
	if in.Provisioning != nil {
		request.Provisioning = &sandbox.ProvisioningOptions{}
		if in.Provisioning.TimeoutMs != nil {
			request.Provisioning.Timeout = sandbox.Duration(time.Duration(*in.Provisioning.TimeoutMs) * time.Millisecond)
		}
		if in.Provisioning.QueueTimeoutMs != nil {
			request.Provisioning.QueueTimeout = sandbox.Duration(time.Duration(*in.Provisioning.QueueTimeoutMs) * time.Millisecond)
		}
	}
	if in.Observability != nil {
		request.Observability = &sandbox.ObservabilityConfig{}
		if in.Observability.Verbose != nil {
			request.Observability.Verbose = sandbox.Value(*in.Observability.Verbose)
		}
		if in.Observability.TelemetryEndpoint != nil {
			request.Observability.TelemetryEndpoint = sandbox.Value(*in.Observability.TelemetryEndpoint)
		}
	}
	return request
}
