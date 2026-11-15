package blaxel_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/providers/blaxel"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func ExampleNew() {
	region := os.Getenv("BL_REGION")
	if region == "" {
		region = "us-pdx-1"
	}
	config := &sandbox.Config{
		Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{
			ApiKey: &sandbox.APIKeyCredentials{Key: os.Getenv("BL_API_KEY")},
		}},
		Region:  sandbox.Value(region),
		Timeout: sandbox.Duration(2 * time.Minute),
	}
	if workspace := strings.TrimSpace(os.Getenv("BL_WORKSPACE")); workspace != "" {
		config.Scope = &sandbox.Scope{Workspace: sandbox.Value(workspace)}
	}
	client, err := sandbox.NewClient(blaxel.New(), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Name: sandbox.Value("kit-example"),
		Source: &sandbox.SandboxSource{Kind: &sandbox.Image{
			Image: &sandbox.ImageSource{Reference: "blaxel/base-image:latest"},
		}},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(instance.ProviderName(), instance.ID())
}

func ExampleNew_image() {
	region := os.Getenv("BL_REGION")
	if region == "" {
		region = "us-pdx-1"
	}
	config := &sandbox.Config{
		Region: sandbox.Value(region),
	}
	if workspace := strings.TrimSpace(os.Getenv("BL_WORKSPACE")); workspace != "" {
		config.Scope = &sandbox.Scope{Workspace: sandbox.Value(workspace)}
	}
	client, err := sandbox.NewClient(blaxel.New(), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Name: sandbox.Value("kit-example-memory"),
		Source: &sandbox.SandboxSource{Kind: &sandbox.Image{
			Image: &sandbox.ImageSource{Reference: "blaxel/base-image:latest"},
		}},
		Resources: &sandbox.Resources{
			MemoryMib: sandbox.Value(uint64(4096)),
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(instance.ID())
}

func ExampleNew_get() {
	config := &sandbox.Config{}
	if workspace := strings.TrimSpace(os.Getenv("BL_WORKSPACE")); workspace != "" {
		config.Scope = &sandbox.Scope{Workspace: sandbox.Value(workspace)}
	}
	client, err := sandbox.NewClient(blaxel.New(), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("BL_SANDBOX_NAME"))
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = instance.Refresh(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_delete() {
	config := &sandbox.Config{}
	if workspace := strings.TrimSpace(os.Getenv("BL_WORKSPACE")); workspace != "" {
		config.Scope = &sandbox.Scope{Workspace: sandbox.Value(workspace)}
	}
	client, err := sandbox.NewClient(blaxel.New(), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("BL_SANDBOX_NAME"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = instance.Delete(context.Background(), &sandbox.DeleteOptions{Wait: sandbox.Value(true)})
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_close() {
	region := os.Getenv("BL_REGION")
	if region == "" {
		region = "us-pdx-1"
	}
	config := &sandbox.Config{
		Region: sandbox.Value(region),
	}
	if workspace := strings.TrimSpace(os.Getenv("BL_WORKSPACE")); workspace != "" {
		config.Scope = &sandbox.Scope{Workspace: sandbox.Value(workspace)}
	}
	client, err := sandbox.NewClient(blaxel.New(), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Name: sandbox.Value("kit-example-close"),
		Source: &sandbox.SandboxSource{Kind: &sandbox.Image{
			Image: &sandbox.ImageSource{Reference: "blaxel/base-image:latest"},
		}},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	// Close leaves the Blaxel sandbox running.
	err = instance.Close(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_errors() {
	_, err := sandbox.NewClient(blaxel.New(), &sandbox.Config{
		Endpoint: sandbox.Value("https://example.com"),
	})
	var detail *sandbox.Error
	if errors.As(err, &detail) {
		fmt.Println(detail.Field)
	}
	// Output:
	// endpoint
}

func ExampleNew_process() {
	config := &sandbox.Config{}
	if workspace := strings.TrimSpace(os.Getenv("BL_WORKSPACE")); workspace != "" {
		config.Scope = &sandbox.Scope{Workspace: sandbox.Value(workspace)}
	}
	client, err := sandbox.NewClient(blaxel.New(), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("BL_SANDBOX_NAME"))
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := instance.Process().Run(context.Background(), &sandbox.ExecOptions{
		Argv: []string{"echo", "hello from kit"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result.GetExitCode(), string(result.GetStdout()))
}

func ExampleNew_processShell() {
	config := &sandbox.Config{}
	if workspace := strings.TrimSpace(os.Getenv("BL_WORKSPACE")); workspace != "" {
		config.Scope = &sandbox.Scope{Workspace: sandbox.Value(workspace)}
	}
	client, err := sandbox.NewClient(blaxel.New(), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("BL_SANDBOX_NAME"))
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := instance.Process().Run(context.Background(), &sandbox.ExecOptions{
		Shell: sandbox.Value(`echo "hello from shell"`),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result.GetExitCode())
}

func ExampleNew_processEnv() {
	config := &sandbox.Config{}
	if workspace := strings.TrimSpace(os.Getenv("BL_WORKSPACE")); workspace != "" {
		config.Scope = &sandbox.Scope{Workspace: sandbox.Value(workspace)}
	}
	client, err := sandbox.NewClient(blaxel.New(), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("BL_SANDBOX_NAME"))
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := instance.Process().Run(context.Background(), &sandbox.ExecOptions{
		Argv:        []string{"printenv", "KIT_DEMO"},
		Environment: map[string]string{"KIT_DEMO": "from-kit"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(result.GetStdout()))
}

func ExampleNew_filesystem() {
	config := &sandbox.Config{}
	if workspace := strings.TrimSpace(os.Getenv("BL_WORKSPACE")); workspace != "" {
		config.Scope = &sandbox.Scope{Workspace: sandbox.Value(workspace)}
	}
	client, err := sandbox.NewClient(blaxel.New(), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("BL_SANDBOX_NAME"))
	if err != nil {
		fmt.Println(err)
		return
	}
	files := instance.Filesystem()
	if err := files.WriteFile(context.Background(), "/tmp/note.txt", []byte("hello\n"), nil); err != nil {
		fmt.Println(err)
		return
	}
	data, err := files.ReadFile(context.Background(), "/tmp/note.txt", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
}

func ExampleNew_makeDirectory() {
	config := &sandbox.Config{}
	if workspace := strings.TrimSpace(os.Getenv("BL_WORKSPACE")); workspace != "" {
		config.Scope = &sandbox.Scope{Workspace: sandbox.Value(workspace)}
	}
	client, err := sandbox.NewClient(blaxel.New(), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("BL_SANDBOX_NAME"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = instance.Filesystem().MakeDirectory(context.Background(), "/tmp/kit-dir", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_listDirectory() {
	config := &sandbox.Config{}
	if workspace := strings.TrimSpace(os.Getenv("BL_WORKSPACE")); workspace != "" {
		config.Scope = &sandbox.Scope{Workspace: sandbox.Value(workspace)}
	}
	client, err := sandbox.NewClient(blaxel.New(), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("BL_SANDBOX_NAME"))
	if err != nil {
		fmt.Println(err)
		return
	}
	entries, err := instance.Filesystem().ListDirectory(context.Background(), "/tmp", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(entries))
}

func ExampleNew_stat() {
	config := &sandbox.Config{}
	if workspace := strings.TrimSpace(os.Getenv("BL_WORKSPACE")); workspace != "" {
		config.Scope = &sandbox.Scope{Workspace: sandbox.Value(workspace)}
	}
	client, err := sandbox.NewClient(blaxel.New(), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("BL_SANDBOX_NAME"))
	if err != nil {
		fmt.Println(err)
		return
	}
	info, err := instance.Filesystem().Stat(context.Background(), "/tmp/note.txt", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(info.GetPath(), info.GetSize())
}
