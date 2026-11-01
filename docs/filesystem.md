# Filesystem

File calls live on the sandbox handle, next to lifecycle. Daytona and Modal
implement the same methods. Open the handle with `Create`, Daytona `Get`, or
Modal `Attach`, then call `Filesystem` / `sandbox.filesystem`.

| You want to                   | Go              | TypeScript      |
| ----------------------------- | --------------- | --------------- |
| Read the bytes of one file    | `ReadFile`      | `readFile`      |
| Replace one file              | `WriteFile`     | `writeFile`     |
| Copy a file from your machine | `CopyFromLocal` | `copyFromLocal` |
| Copy a file onto your machine | `CopyToLocal`   | `copyToLocal`   |
| Create a directory            | `MakeDirectory` | `makeDirectory` |
| List one directory            | `ListDirectory` | `listDirectory` |
| Read metadata for one path    | `Stat`          | `stat`          |
| Delete a file or directory    | `Remove`        | `remove`        |

Calls on one handle run one at a time, including lifecycle calls. A second
handle for the same sandbox does not wait.

## Paths

A path inside the sandbox is absolute: `/tmp/note.txt`. A relative path, a
blank path, or a path that contains a NUL byte is rejected before the provider
is called. The error field is `path`.

A path on your machine, used by copy, may be relative to the process, such as
`./note.txt`. It must be non-blank and must not contain a NUL byte. A rejected
local path uses the field `local_path`.

## Bytes and disk

`ReadFile` and `WriteFile` move bytes. Encode text before you write it, and
decode text after you read it. `WriteFile` replaces an existing file and
creates missing parent directories.

`CopyFromLocal` and `CopyToLocal` use files on disk. The local file is the
first argument of `CopyFromLocal` and the second argument of `CopyToLocal`.
Copy creates missing parent directories of the destination and replaces a file
that is already there.

## Timeouts

Leave the timeout unset to use the client deadline. Set it to zero when this
call should add no Kit deadline. Daytona accepts a timeout only as a whole
number of seconds, so 1500 milliseconds is rejected and 2 seconds is passed
through. Modal filesystem methods take the context or abort signal as the
deadline.

```ts
const data = await sandbox.filesystem.readFile("/tmp/note.txt", {
  timeout: 2_000,
});
```

```go
data, err := instance.Filesystem().ReadFile(ctx, "/tmp/note.txt", &sandbox.FileReadOptions{
    Timeout: sandbox.Duration(2 * time.Second),
})
```

## Read a file

`ReadFile` returns the bytes of one file. A missing file comes back as not
found. A directory path fails at the provider.

```ts
const data = await sandbox.filesystem.readFile("/tmp/note.txt");
const text = new TextDecoder().decode(data);
```

```go
data, err := instance.Filesystem().ReadFile(ctx, "/tmp/note.txt", nil)
```

[Go Daytona](../examples/go/daytona/read-daytona-file/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/read-daytona-file/README.md) ·
[Go Modal](../examples/go/modal/read-modal-file/README.md) ·
[TypeScript Modal](../examples/typescript/modal/read-modal-file/README.md).

## Write a file

`WriteFile` replaces the file at that path. Parent directories are created
when they are missing.

```ts
const body = new TextEncoder().encode("hello\n");
await sandbox.filesystem.writeFile("/tmp/kit-files/note.txt", body);
```

```go
err := instance.Filesystem().WriteFile(ctx, "/tmp/kit-files/note.txt", []byte("hello\n"), nil)
```

[Go Daytona](../examples/go/daytona/write-daytona-file/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/write-daytona-file/README.md) ·
[Go Modal](../examples/go/modal/write-modal-file/README.md) ·
[TypeScript Modal](../examples/typescript/modal/write-modal-file/README.md).

## Copy a file

`CopyFromLocal` reads a file on your machine and writes it in the sandbox.
`CopyToLocal` does the reverse. An existing destination file is replaced.

```ts
await sandbox.filesystem.copyFromLocal("./note.txt", "/tmp/note.txt");
await sandbox.filesystem.copyToLocal("/tmp/note.txt", "./note.txt");
```

```go
files := instance.Filesystem()
if err := files.CopyFromLocal(ctx, "note.txt", "/tmp/note.txt", nil); err != nil {
    return err
}
if err := files.CopyToLocal(ctx, "/tmp/note.txt", "note.txt", nil); err != nil {
    return err
}
```

Into the sandbox:
[Go Daytona](../examples/go/daytona/copy-file-into-daytona/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/copy-file-into-daytona/README.md) ·
[Go Modal](../examples/go/modal/copy-file-into-modal/README.md) ·
[TypeScript Modal](../examples/typescript/modal/copy-file-into-modal/README.md).

Out of the sandbox:
[Go Daytona](../examples/go/daytona/copy-file-from-daytona/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/copy-file-from-daytona/README.md) ·
[Go Modal](../examples/go/modal/copy-file-from-modal/README.md) ·
[TypeScript Modal](../examples/typescript/modal/copy-file-from-modal/README.md).

## Create a directory

`parents` defaults to true, so the call creates missing parents. Modal accepts
a directory that already exists. Daytona returns the provider error when it
already exists.

Set `parents` to false when the parent directory must already exist. Modal
then requires that parent, and returns an error if the new directory is
already there. Daytona rejects `false` before the call. The error field is
`parents`.

```ts
await sandbox.filesystem.makeDirectory("/tmp/kit-files");
await sandbox.filesystem.makeDirectory("/tmp/kit-files/nested", {
  parents: false,
});
```

```go
err := instance.Filesystem().MakeDirectory(ctx, "/tmp/kit-files", nil)
err = instance.Filesystem().MakeDirectory(ctx, "/tmp/kit-files/nested", &sandbox.MakeDirectoryOptions{
	Parents: sandbox.Value(false),
})
```

[Go Daytona](../examples/go/daytona/make-daytona-directory/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/make-daytona-directory/README.md) ·
[Go Modal](../examples/go/modal/make-modal-directory/README.md) ·
[TypeScript Modal](../examples/typescript/modal/make-modal-directory/README.md).

## List a directory

`ListDirectory` returns the entries in that directory. It does not walk into
subdirectories. A file path fails at the provider.

```ts
const entries = await sandbox.filesystem.listDirectory("/tmp/kit-files");
for (const entry of entries) console.log(entry.name, entry.size);
```

```go
entries, err := instance.Filesystem().ListDirectory(ctx, "/tmp/kit-files", nil)
```

`size` is a byte count: `bigint` in TypeScript and `int64` in Go.

[Go Daytona](../examples/go/daytona/list-daytona-directory/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/list-daytona-directory/README.md) ·
[Go Modal](../examples/go/modal/list-modal-directory/README.md) ·
[TypeScript Modal](../examples/typescript/modal/list-modal-directory/README.md).

## Read metadata

`Stat` returns one path. A symlink result describes the symlink. Modal can
report a symlink. Daytona reports a file or a directory. `mode` is an octal
string when the provider supplies one, such as `0644`. `modifiedAt` is set
when the provider supplies a time.

```ts
const info = await sandbox.filesystem.stat("/tmp/note.txt");
console.log(info.type, info.size);
```

```go
info, err := instance.Filesystem().Stat(ctx, "/tmp/note.txt", nil)
```

[Go Daytona](../examples/go/daytona/stat-daytona-file/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/stat-daytona-file/README.md) ·
[Go Modal](../examples/go/modal/stat-modal-file/README.md) ·
[TypeScript Modal](../examples/typescript/modal/stat-modal-file/README.md).

## Delete a file or directory

`Remove` with no `recursive` flag deletes a file or an empty directory. Set
`recursive` to delete a directory and its contents. A directory that still
has entries fails when the flag is omitted. A missing path comes back as not
found.

```ts
await sandbox.filesystem.remove("/tmp/note.txt");
await sandbox.filesystem.remove("/tmp/kit-files", { recursive: true });
```

```go
files := instance.Filesystem()
if err := files.Remove(ctx, "/tmp/note.txt", nil); err != nil {
    return err
}
if err := files.Remove(ctx, "/tmp/kit-files", &sandbox.RemoveOptions{
    Recursive: sandbox.Value(true),
}); err != nil {
    return err
}
```

[Go Daytona](../examples/go/daytona/remove-daytona-file/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/remove-daytona-file/README.md) ·
[Go Modal](../examples/go/modal/remove-modal-file/README.md) ·
[TypeScript Modal](../examples/typescript/modal/remove-modal-file/README.md).
