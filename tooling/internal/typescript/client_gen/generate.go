// Package clientgen emits the TypeScript client runtime from compiled API and
// lifecycle declarations.
package clientgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/typescript/declaration_gen"
	lifecyclegen "github.com/sandbox-kit/kit/tooling/internal/typescript/lifecycle_gen"
	typesgen "github.com/sandbox-kit/kit/tooling/internal/typescript/types_gen"
	validationgen "github.com/sandbox-kit/kit/tooling/internal/typescript/validation_gen"
	"google.golang.org/protobuf/compiler/protogen"
)

func Generate(plugin *protogen.Plugin, file *protogen.File, program *model.Program, profile spec.LanguageProfile) error {
	if err := declarationgen.ValidateModule(program.API.Modules["client"], profile); err != nil {
		return fmt.Errorf("TypeScript client declarations: %w", err)
	}
	if err := spec.ValidateOperations(program.Generation.Operations); err != nil {
		return fmt.Errorf("TypeScript client operations: %w", err)
	}
	g := declarationgen.GeneratedFile(plugin, "client"+profile.Outputs["client"], string(file.Desc.Path()))
	e := declarationgen.NewEmitter(g)
	validationTypes, validationValues, err := validationgen.RequiredSymbols(program.Schema, "CreateOptions")
	if err != nil {
		return err
	}
	configTypes, configValues, err := validationgen.RequiredSymbols(program.Schema, program.Client.ConfigMessage)
	if err != nil {
		return err
	}
	typeImports := map[string]bool{"ArchiveOptions": true, "CreateOptions": true, "CreateResult": true, "DeleteOptions": true, "DurationMs": true, "ErrorInfo": true, "ExitResult": true, "LifecycleCapability": true, "PauseOptions": true, "SandboxInfo": true, "StartOptions": true, "StopOptions": true, "TerminateOptions": true}
	for _, message := range validationTypes {
		typeImports[string(message.Name())] = true
	}
	_ = configTypes // Config declarations are emitted into this same output file.
	valueImports := map[string]bool{"ErrorKind": true, "SandboxState": true}
	for _, value := range validationValues {
		valueImports[value] = true
	}
	for _, value := range configValues {
		valueImports[value] = true
	}
	e.Line("import { ", strings.Join(sortedImportNames(valueImports, ""), ", "), ", ", strings.Join(sortedImportNames(typeImports, "type "), ", "), ` } from "./sandbox.gen.js";`)
	e.Blank()
	inputMessages, err := typesgen.InputMessages(program)
	if err != nil {
		return err
	}
	if err := typesgen.EmitDeclarations(g, file, false, inputMessages); err != nil {
		return err
	}
	e.Block(`
export interface OperationOptions {
  readonly signal?: AbortSignal;
  readonly timeoutMs?: number;
}

export interface Backend {
	readonly name: string;`)
	lifecyclegen.EmitServiceBackendMethods(g, program.Clients[file.Desc.Path()].ServiceOperations)
	e.Block(`
  close(options: OperationOptions): Promise<void>;
  closeSandbox?(id: string, options: OperationOptions): Promise<void>;
  capabilities?(info: SandboxInfo): readonly LifecycleCapability[];
}

export interface Provider {
  readonly name: string;
  initialize(config: Config): Backend;
}

export interface ClientConfig extends Config {
  readonly provider: Provider;
}

export class SandboxError extends Error {
  readonly kind: ErrorInfo["kind"];
  readonly provider: string;
  readonly operation: string;
  readonly field: string;
  readonly rule: string | undefined;
  readonly statusCode: number | undefined;
  readonly providerCode: string | undefined;
  readonly providerSource: string | undefined;

  constructor(info: ErrorInfo, options?: ErrorOptions) {
    super(info.message, options);
    this.name = "SandboxError";
    this.kind = info.kind;
    this.provider = info.provider;
    this.operation = info.operation;
    this.field = info.field;
    this.rule = info.rule;
    this.statusCode = info.statusCode;
    this.providerCode = info.providerCode;
    this.providerSource = info.providerSource;
  }
}

export function sandboxError(info: ErrorInfo, cause?: unknown): SandboxError {
  return new SandboxError(info, cause === undefined ? undefined : { cause });
}
`)
	if err := validationgen.Emit(g, program.Schema, program.Validation, "CreateOptions", true); err != nil {
		return err
	}
	if err := validationgen.Emit(g, program.Schema, program.Client.Validation, program.Client.ConfigMessage, false); err != nil {
		return err
	}
	e.Block(runtimeHelpers)
	g.P(`interface ClientLifecycle { begin(operation: string): () => void; isClosed(): boolean; }`)
	g.P(`const sandboxConstructorToken: unique symbol = Symbol("sandbox-kit.Sandbox");`)
	g.P(`export class Sandbox implements AsyncDisposable {`)
	g.P(`  readonly #backend: Backend; readonly #owner: ClientLifecycle; readonly #timeoutMs: number | undefined; #info: SandboxInfo; #closed = false; #tail: Promise<void> = Promise.resolve();`)
	g.P(`  constructor(token: typeof sandboxConstructorToken, backend: Backend, owner: ClientLifecycle, info: SandboxInfo, timeoutMs?: number) { if (token !== sandboxConstructorToken) throw new TypeError("Sandbox handles are created by Client"); this.#backend = backend; this.#owner = owner; this.#info = clone(info); this.#timeoutMs = timeoutMs; }`)
	g.P(`  get id(): string { return this.#info.id; } get providerName(): string { return this.#info.provider; } info(): SandboxInfo { return clone(this.#info); }`)
	g.P(`  capabilities(): readonly LifecycleCapability[] { if (this.#closed || this.#owner.isClosed() || this.#info.terminal) return []; return this.#backend.capabilities?.(clone(this.#info)) ?? []; }`)
	g.P(`  supports(capability: LifecycleCapability): boolean { return this.capabilities().includes(capability); }`)
	if err := lifecyclegen.EmitHandleMethods(g, program.Lifecycle); err != nil {
		return err
	}
	g.P(`  async close(options: OperationOptions = {}): Promise<void> { if (this.#closed) return; await this.#serialized("close", options, this.#timeoutMs, async (signal) => { const close = this.#backend.closeSandbox; if (close) await awaitWithSignal(() => close.call(this.#backend, this.id, operationOptions(options, signal)), signal); this.#closed = true; }); }`)
	g.P(`  async [Symbol.asyncDispose](): Promise<void> { await this.close(); }`)
	updateOptions, err := lifecyclegen.UpdateOptionsUnion(program.Lifecycle)
	if err != nil {
		return err
	}
	if err := e.IndentedTemplate("sandbox update helper", sandboxUpdateTemplate, struct{ Options string }{updateOptions}, "  "); err != nil {
		return err
	}
	g.P(`  async #serialized<T>(operation: string, options: OperationOptions, timeoutMs: number | undefined, body: (signal: AbortSignal | undefined) => Promise<T>): Promise<T> { this.#ensureOpen(operation); const releaseClient = this.#owner.begin(operation); const previous = this.#tail; let unlock!: () => void; this.#tail = new Promise<void>((resolve) => { unlock = resolve; }); let acquired = false; let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, timeoutMs); await awaitWithSignal(() => previous, scoped.signal); acquired = true; this.#ensureOpen(operation); return await body(scoped.signal); } catch (error) { throw withContext(error, this.providerName, operation); } finally { scoped?.dispose(); if (acquired) unlock(); else void previous.then(unlock, unlock); releaseClient(); } }`)
	g.P(`  #ensureOpen(operation: string): void { if (this.#closed || this.#owner.isClosed()) throw sandboxError({ kind: ErrorKind.FailedPrecondition, provider: this.providerName, operation, field: "sandbox", message: "sandbox-kit: sandbox handle is closed" }); if (operation !== "close" && this.#info.terminal) throw sandboxError({ kind: ErrorKind.FailedPrecondition, provider: this.providerName, operation, field: "sandbox", message: "sandbox-kit: sandbox lifecycle is terminal" }); }`)
	g.P(`}`)
	g.P()
	g.P(`export class Client implements AsyncDisposable {`)
	g.P(`  readonly #backend: Backend; readonly #providerName: string; readonly #timeoutMs: number | undefined; readonly #lifecycle: ClientLifecycle; #active = 0; #closing = false; #closed = false; #closePromise: Promise<void> | undefined; #drain: Promise<void> | undefined; #resolveDrain: (() => void) | undefined;`)
	g.P(`  constructor(config: ClientConfig) { if (!config.provider?.name?.trim()) throw sandboxError({ kind: ErrorKind.InvalidArgument, provider: "", operation: "initialize", field: "provider", message: "sandbox-kit: provider is required" }); validateConfig(config, config.provider.name, "initialize"); this.#providerName = config.provider.name; this.#timeoutMs = config.timeout; this.#lifecycle = { begin: (operation) => this.#begin(operation), isClosed: () => this.#closing || this.#closed }; try { this.#backend = config.provider.initialize(cloneConfig(config)); } catch (error) { throw withContext(error, config.provider.name, "initialize"); } if (!this.#backend) throw sandboxError({ kind: ErrorKind.InvalidResponse, provider: config.provider.name, operation: "initialize", field: "provider", message: "sandbox-kit: provider returned no backend" }); if (this.#backend.name !== this.#providerName) { void this.#backend.close({}).catch(() => {}); throw sandboxError({ kind: ErrorKind.InvalidResponse, provider: this.#providerName, operation: "initialize", field: "provider", message: "sandbox-kit: initialized backend identity differs from selected provider" }); } }`)
	g.P(`  get providerName(): string { return this.#providerName; }`)
	g.P(`  async attach(id: string, options: OperationOptions = {}): Promise<Sandbox> { const release = this.#begin("attach"); let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, this.#timeoutMs); if (!id.trim()) throw sandboxError({ kind: ErrorKind.InvalidArgument, provider: this.providerName, operation: "attach", field: "id", message: "sandbox-kit: sandbox ID must not be blank" }); if (!this.#backend.attach) unsupported(this.providerName, "attach"); const info = await awaitWithSignal(() => this.#backend.attach!(id, operationOptions(options, scoped?.signal)), scoped.signal); return new Sandbox(sandboxConstructorToken, this.#backend, this.#lifecycle, { ...info, id, provider: this.providerName }, this.#timeoutMs); } catch (error) { throw withContext(error, this.providerName, "attach"); } finally { scoped?.dispose(); release(); } }`)
	g.P(`  async get(id: string, options: OperationOptions = {}): Promise<Sandbox> { const release = this.#begin("get"); let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, this.#timeoutMs); if (!id.trim()) throw sandboxError({ kind: ErrorKind.InvalidArgument, provider: this.providerName, operation: "get", field: "id", message: "sandbox-kit: sandbox ID must not be blank" }); if (!this.#backend.refresh) unsupported(this.providerName, "get"); const info = await awaitWithSignal(() => this.#backend.refresh!(id, operationOptions(options, scoped?.signal)), scoped.signal); return new Sandbox(sandboxConstructorToken, this.#backend, this.#lifecycle, { ...info, id, provider: this.providerName }, this.#timeoutMs); } catch (error) { throw withContext(error, this.providerName, "get"); } finally { scoped?.dispose(); release(); } }`)
	g.P(`  async create(request: CreateOptions = {}, options: OperationOptions = {}): Promise<Sandbox> { const release = this.#begin("create"); const timeoutMs = options.timeoutMs ?? request.provisioning?.timeout ?? this.#timeoutMs; let scoped: ReturnType<typeof operationSignal> | undefined; try { validateCreateOptions(request, this.providerName, "create"); scoped = operationSignal(options, timeoutMs); const create = this.#backend.create; if (!create) unsupported(this.providerName, "create"); const result = await awaitWithSignal(() => create.call(this.#backend, clone(request), operationOptions({ ...(timeoutMs === undefined ? {} : { timeoutMs }) }, scoped?.signal)), scoped.signal); if (!result.sandbox?.id) throw sandboxError({ kind: ErrorKind.InvalidResponse, provider: this.providerName, operation: "create", field: "id", message: "sandbox-kit: provider returned no sandbox identity" }); if (result.sandbox.provider && result.sandbox.provider !== this.providerName) throw sandboxError({ kind: ErrorKind.InvalidResponse, provider: this.providerName, operation: "create", field: "provider", message: "sandbox-kit: response provider differs from selected provider" }); return new Sandbox(sandboxConstructorToken, this.#backend, this.#lifecycle, { ...result.sandbox, provider: this.providerName }, this.#timeoutMs); } catch (error) { throw withContext(error, this.providerName, "create"); } finally { scoped?.dispose(); release(); } }`)
	g.P(`  close(options: OperationOptions = {}): Promise<void> { if (this.#closePromise) return this.#closePromise; this.#closing = true; this.#closePromise = this.#close(options); return this.#closePromise; }`)
	g.P(`  async #close(options: OperationOptions): Promise<void> { let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, this.#timeoutMs); if (this.#active > 0) { this.#drain ??= new Promise<void>((resolve) => { this.#resolveDrain = resolve; }); await awaitWithSignal(() => this.#drain!, scoped.signal); } await awaitWithSignal(() => this.#backend.close(operationOptions(options, scoped?.signal)), scoped.signal); this.#closed = true; } catch (error) { this.#closePromise = undefined; this.#closing = false; throw withContext(error, this.providerName, "close"); } finally { scoped?.dispose(); } }`)
	g.P(`  async [Symbol.asyncDispose](): Promise<void> { await this.close(); }`)
	g.P(`  #begin(operation: string): () => void { if (this.#closing || this.#closed) throw sandboxError({ kind: ErrorKind.FailedPrecondition, provider: this.providerName, operation, field: "client", message: "sandbox-kit: client is closed" }); this.#active++; let released = false; return () => { if (released) return; released = true; this.#active--; if (this.#active === 0) { this.#resolveDrain?.(); this.#resolveDrain = undefined; this.#drain = undefined; } }; }`)
	g.P(`}`)
	g.P()
	g.P(`function cloneConfig(config: ClientConfig): Config { const { provider: _provider, ...shared } = config; return clone(shared); }`)
	g.P(`function unsupported(provider: string, operation: string): never { throw sandboxError({ kind: ErrorKind.Unsupported, provider, operation, field: "", message: "sandbox-kit: provider does not support " + operation }); }`)
	g.P(`function withContext(error: unknown, provider: string, operation: string): SandboxError { if (error instanceof SandboxError) { if (error.provider === provider && error.operation === operation) return error; return sandboxError({ kind: error.kind, provider: error.provider || provider, operation: error.operation || operation, field: error.field, message: error.message, ...(error.rule === undefined ? {} : { rule: error.rule }), ...(error.statusCode === undefined ? {} : { statusCode: error.statusCode }), ...(error.providerCode === undefined ? {} : { providerCode: error.providerCode }), ...(error.providerSource === undefined ? {} : { providerSource: error.providerSource }) }, error.cause); } if (error instanceof DOMException && error.name === "AbortError") return sandboxError({ kind: ErrorKind.Canceled, provider, operation, field: "", message: error.message }, error); if (error instanceof DOMException && error.name === "TimeoutError") return sandboxError({ kind: ErrorKind.Timeout, provider, operation, field: "", message: error.message }, error); return sandboxError({ kind: ErrorKind.Unknown, provider, operation, field: "", message: error instanceof Error ? error.message : "sandbox-kit: provider operation failed" }, error); }`)
	return nil
}

func sortedImportNames(values map[string]bool, prefix string) []string {
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, prefix+name)
	}
	sort.Strings(names)
	return names
}
