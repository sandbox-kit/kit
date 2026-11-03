package modal_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/providers/modal"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func ExampleNew() {
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Auth: &sandbox.AuthConfig{Credentials: &sandbox.TokenPair{
			TokenPair: &sandbox.TokenPairCredentials{
				Id:     os.Getenv("MODAL_TOKEN_ID"),
				Secret: os.Getenv("MODAL_TOKEN_SECRET"),
			},
		}},
		Scope: &sandbox.Scope{
			AppName:     sandbox.Value(os.Getenv("MODAL_APP_NAME")),
			Environment: sandbox.Value("main"),
		},
		Timeout: sandbox.Duration(2 * time.Minute),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Source: &sandbox.SandboxSource{Kind: &sandbox.Image{
			Image: &sandbox.ImageSource{Reference: "alpine:3.21"},
		}},
		Lifetime: &sandbox.LifetimePolicy{
			MaximumLifetime: sandbox.Duration(5 * time.Minute),
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(instance.ID())
}

func ExampleNew_attach() {
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Scope: &sandbox.Scope{AppName: sandbox.Value(os.Getenv("MODAL_APP_NAME"))},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Attach(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(instance.ID())
}

func ExampleNew_terminate() {
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Scope: &sandbox.Scope{AppName: sandbox.Value(os.Getenv("MODAL_APP_NAME"))},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Attach(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = instance.Terminate(context.Background(), &sandbox.TerminateOptions{
		Wait: sandbox.Value(true),
	})
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_errors() {
	_, err := sandbox.NewClient(modal.New(), &sandbox.Config{
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
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Scope: &sandbox.Scope{AppName: sandbox.Value(os.Getenv("MODAL_APP_NAME"))},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Attach(context.Background(), os.Getenv("SANDBOX_ID"))
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
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Scope: &sandbox.Scope{AppName: sandbox.Value(os.Getenv("MODAL_APP_NAME"))},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Attach(context.Background(), os.Getenv("SANDBOX_ID"))
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
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Scope: &sandbox.Scope{AppName: sandbox.Value(os.Getenv("MODAL_APP_NAME"))},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Attach(context.Background(), os.Getenv("SANDBOX_ID"))
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
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Scope: &sandbox.Scope{AppName: sandbox.Value(os.Getenv("MODAL_APP_NAME"))},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Attach(context.Background(), os.Getenv("SANDBOX_ID"))
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
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Scope: &sandbox.Scope{AppName: sandbox.Value(os.Getenv("MODAL_APP_NAME"))},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Attach(context.Background(), os.Getenv("SANDBOX_ID"))
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
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Scope: &sandbox.Scope{AppName: sandbox.Value(os.Getenv("MODAL_APP_NAME"))},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Attach(context.Background(), os.Getenv("SANDBOX_ID"))
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
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Scope: &sandbox.Scope{AppName: sandbox.Value(os.Getenv("MODAL_APP_NAME"))},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Attach(context.Background(), os.Getenv("SANDBOX_ID"))
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
