# Filesystem

File calls live on the sandbox handle, next to lifecycle. Daytona, Modal, and Blaxel
implement the same methods. Open the handle with `Create`, Daytona or Blaxel `Get`, or
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

Calls on one handle run one at a time, including lifecycle and process. A
second handle for the same sandbox does not wait. Abort or deadline rejects the
TypeScript `await` (and ends the Go wait), but the handle stays busy until that
provider call settles so the next call does not overlap it.

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
An existing destination file is replaced.

## Timeouts

Leave the timeout unset to use the client deadline. Set it to zero when this
call should add no Kit deadline. TypeScript records the timeout in
milliseconds. Go records it as a `time.Duration`.

TypeScript Daytona `readFile`, `writeFile`, `copyFromLocal`, and `copyToLocal`
pass that duration to Daytona as a whole number of seconds. `2000` is two
seconds. `1500` is 1.5 seconds, so those four calls reject it. The error field
is `timeout`. An omitted timeout, and an explicit zero, on those four calls
are passed as `0`. Daytona adds no transfer limit. When the timeout is
omitted, the abort signal carries the client deadline.

`makeDirectory`, `listDirectory`, `stat`, and `remove` use the abort signal.
`1500` ends the call at 1.5 seconds. Every Go file call uses the context as
the deadline, so `1500 * time.Millisecond` ends the call at 1.5 seconds. Modal
and Blaxel file calls use the context or the abort signal. `1500` milliseconds
ends the call at 1.5 seconds. Blaxel rejects `parents: false` before the
provider call (field `parents`).

```ts
await sandbox.filesystem.readFile("/tmp/note.txt", { timeout: 2_000 });
await sandbox.filesystem.listDirectory("/tmp", { timeout: 1_500 });
```

```go
data, err := instance.Filesystem().ReadFile(ctx, "/tmp/note.txt", &sandbox.FileReadOptions{
    Timeout: sandbox.Duration(1500 * time.Millisecond),
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

[Go Daytona](../examples/go/daytona/filesystem/read-daytona-file/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/filesystem/read-daytona-file/README.md) ·
[Go Modal](../examples/go/modal/filesystem/read-modal-file/README.md) ·
[TypeScript Modal](../examples/typescript/modal/filesystem/read-modal-file/README.md) ·
[Go Blaxel](../examples/go/blaxel/filesystem/read-blaxel-file/README.md) ·
[TypeScript Blaxel](../examples/typescript/blaxel/filesystem/read-blaxel-file/README.md).

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

[Go Daytona](../examples/go/daytona/filesystem/write-daytona-file/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/filesystem/write-daytona-file/README.md) ·
[Go Modal](../examples/go/modal/filesystem/write-modal-file/README.md) ·
[TypeScript Modal](../examples/typescript/modal/filesystem/write-modal-file/README.md) ·
[Go Blaxel](../examples/go/blaxel/filesystem/write-blaxel-file/README.md) ·
[TypeScript Blaxel](../examples/typescript/blaxel/filesystem/write-blaxel-file/README.md).

## Copy a file

`CopyFromLocal` reads a file on your machine and writes it in the sandbox.
`CopyToLocal` does the reverse. `CopyFromLocal` creates missing parent
directories of the sandbox path.

A missing local file is not found. The error field is `local_path`. A local
directory is an invalid argument with the same field.

`CopyToLocal` creates a missing local parent directory. A local destination
that is already a directory is an invalid argument, and the field is
`local_path`.

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
[Go Daytona](../examples/go/daytona/filesystem/copy-file-into-daytona/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/filesystem/copy-file-into-daytona/README.md) ·
[Go Modal](../examples/go/modal/filesystem/copy-file-into-modal/README.md) ·
[TypeScript Modal](../examples/typescript/modal/filesystem/copy-file-into-modal/README.md) ·
[Go Blaxel](../examples/go/blaxel/filesystem/copy-file-into-blaxel/README.md) ·
[TypeScript Blaxel](../examples/typescript/blaxel/filesystem/copy-file-into-blaxel/README.md).

Out of the sandbox:
[Go Daytona](../examples/go/daytona/filesystem/copy-file-from-daytona/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/filesystem/copy-file-from-daytona/README.md) ·
[Go Modal](../examples/go/modal/filesystem/copy-file-from-modal/README.md) ·
[TypeScript Modal](../examples/typescript/modal/filesystem/copy-file-from-modal/README.md) ·
[Go Blaxel](../examples/go/blaxel/filesystem/copy-file-from-blaxel/README.md) ·
[TypeScript Blaxel](../examples/typescript/blaxel/filesystem/copy-file-from-blaxel/README.md).

## Create a directory

`parents` defaults to true, so the call creates missing parents. Modal accepts
a directory that already exists. Daytona returns the provider error when it
already exists. Blaxel creates missing parents when `parents` is omitted or
true.

Set `parents` to false when the parent directory must already exist. Modal
then requires that parent, and returns an error if the new directory is
already there. Daytona and Blaxel reject `false` before the call. The error
field is `parents`.

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

[Go Daytona](../examples/go/daytona/filesystem/make-daytona-directory/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/filesystem/make-daytona-directory/README.md) ·
[Go Modal](../examples/go/modal/filesystem/make-modal-directory/README.md) ·
[TypeScript Modal](../examples/typescript/modal/filesystem/make-modal-directory/README.md) ·
[Go Blaxel](../examples/go/blaxel/filesystem/make-blaxel-directory/README.md) ·
[TypeScript Blaxel](../examples/typescript/blaxel/filesystem/make-blaxel-directory/README.md).

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

[Go Daytona](../examples/go/daytona/filesystem/list-daytona-directory/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/filesystem/list-daytona-directory/README.md) ·
[Go Modal](../examples/go/modal/filesystem/list-modal-directory/README.md) ·
[TypeScript Modal](../examples/typescript/modal/filesystem/list-modal-directory/README.md) ·
[Go Blaxel](../examples/go/blaxel/filesystem/list-blaxel-directory/README.md) ·
[TypeScript Blaxel](../examples/typescript/blaxel/filesystem/list-blaxel-directory/README.md).

## Read metadata

`Stat` returns one path. A symlink result describes the symlink. Blaxel has no
dedicated stat API: TypeScript uses `fs.ls` (a directory listing, or a file
when Blaxel returns "Directory not found" for a file path); Go calls the
sandbox filesystem GET through a rebuilt lower client because
`SandboxInstanceFSService` has no public Get/Stat and `LS` always treats the
response as a directory. Modal can
report a symlink. Daytona reports a file or a directory. `mode` is an octal
string when the provider supplies one, such as `0644`. `permissions`, `owner`,
and `group` are set when the provider supplies them. `modifiedAt` is set when
the provider supplies a time.

```ts
const info = await sandbox.filesystem.stat("/tmp/note.txt");
console.log(info.type, info.size);
```

```go
info, err := instance.Filesystem().Stat(ctx, "/tmp/note.txt", nil)
```

[Go Daytona](../examples/go/daytona/filesystem/stat-daytona-file/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/filesystem/stat-daytona-file/README.md) ·
[Go Modal](../examples/go/modal/filesystem/stat-modal-file/README.md) ·
[TypeScript Modal](../examples/typescript/modal/filesystem/stat-modal-file/README.md) ·
[Go Blaxel](../examples/go/blaxel/filesystem/stat-blaxel-file/README.md) ·
[TypeScript Blaxel](../examples/typescript/blaxel/filesystem/stat-blaxel-file/README.md).

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

[Go Daytona](../examples/go/daytona/filesystem/remove-daytona-file/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/filesystem/remove-daytona-file/README.md) ·
[Go Modal](../examples/go/modal/filesystem/remove-modal-file/README.md) ·
[TypeScript Modal](../examples/typescript/modal/filesystem/remove-modal-file/README.md) ·
[Go Blaxel](../examples/go/blaxel/filesystem/remove-blaxel-file/README.md) ·
[TypeScript Blaxel](../examples/typescript/blaxel/filesystem/remove-blaxel-file/README.md).
