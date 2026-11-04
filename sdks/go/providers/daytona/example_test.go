package daytona_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/providers/daytona"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func ExampleNew() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{
		Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{
			ApiKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")},
		}},
		Region:  sandbox.Value("us"),
		Timeout: sandbox.Duration(2 * time.Minute),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(instance.ProviderName(), instance.ID())
}

func ExampleNew_image() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Source: &sandbox.SandboxSource{Kind: &sandbox.Image{
			Image: &sandbox.ImageSource{Reference: "python:3.11"},
		}},
		Resources: &sandbox.Resources{
			CpuCores:  sandbox.Value(2.0),
			MemoryMib: sandbox.Value(uint64(4096)),
			DiskMib:   sandbox.Value(uint64(8192)),
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(instance.ID())
}

func ExampleNew_get() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = instance.Refresh(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_stop() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = instance.Stop(context.Background(), &sandbox.StopOptions{Force: sandbox.Value(false)})
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_info() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	info := instance.Info()
	fmt.Println(info.GetId(), instance.ProviderName())
}

func ExampleNew_supports() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if instance.Supports(sandbox.LifecycleCapabilityPause) {
		err = instance.Pause(context.Background(), nil)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func ExampleNew_start() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = instance.Start(context.Background(), &sandbox.StartOptions{
		Timeout: sandbox.Duration(time.Minute),
	})
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_pause() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = instance.Pause(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_archive() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = instance.Archive(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_delete() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = instance.Delete(context.Background(), &sandbox.DeleteOptions{})
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_close() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Close leaves the Daytona sandbox running.
	err = instance.Close(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_errors() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		var detail *sandbox.Error
		if errors.As(err, &detail) {
			fmt.Println(detail.Kind, detail.Field)
		}
		return
	}
	defer client.Close(context.Background())
}

func ExampleNew_process() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
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
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
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
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
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
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
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
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
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
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
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
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
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
