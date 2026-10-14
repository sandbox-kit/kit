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
	operations := program.Generation.Operations
	for _, name := range []string{"initialize", "create", "close"} {
		if len(operations[name]) == 0 {
			return fmt.Errorf("client operation %s has no compiled instructions", name)
		}
	}
	g := declarationgen.GeneratedFile(plugin, "client"+profile.Outputs["client"], string(file.Desc.Path()))
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
	g.P("import { ", strings.Join(sortedImportNames(valueImports, ""), ", "), ", ", strings.Join(sortedImportNames(typeImports, "type "), ", "), ` } from "./sandbox.gen.js";`)
	g.P()
	inputMessages, err := typesgen.InputMessages(program)
	if err != nil {
		return err
	}
	if err := typesgen.EmitDeclarations(g, file, false, inputMessages); err != nil {
		return err
	}
	g.P(`export interface OperationOptions { readonly signal?: AbortSignal; readonly timeoutMs?: number; }`)
	g.P(`export interface Backend {`)
	g.P(`  readonly name: string;`)
	g.P(`  create(request: CreateOptions, options: OperationOptions): Promise<CreateResult>;`)
	g.P(`  close(options: OperationOptions): Promise<void>;`)
	g.P(`  attach?(id: string, options: OperationOptions): Promise<SandboxInfo>;`)
	lifecyclegen.EmitBackendMethods(g, program.Lifecycle)
	g.P(`  closeSandbox?(id: string, options: OperationOptions): Promise<void>;`)
	g.P(`  capabilities?(info: SandboxInfo): readonly LifecycleCapability[];`)
	g.P(`}`)
	g.P(`export interface Provider { readonly name: string; initialize(config: Config): Backend; }`)
	g.P(`export interface ClientConfig extends Config { readonly provider: Provider; }`)
	g.P()
	g.P(`export class SandboxError extends Error {`)
	g.P(`  readonly kind: ErrorInfo["kind"]; readonly provider: string; readonly operation: string; readonly field: string; readonly statusCode: number | undefined; readonly providerCode: string | undefined; readonly providerSource: string | undefined;`)
	g.P(`  constructor(info: ErrorInfo, options?: ErrorOptions) { super(info.message, options); this.name = "SandboxError"; this.kind = info.kind; this.provider = info.provider; this.operation = info.operation; this.field = info.field; this.statusCode = info.statusCode; this.providerCode = info.providerCode; this.providerSource = info.providerSource; }`)
	g.P(`}`)
	g.P(`export function sandboxError(info: ErrorInfo, cause?: unknown): SandboxError { return new SandboxError(info, cause === undefined ? undefined : { cause }); }`)
	g.P()
	if err := validationgen.Emit(g, program.Schema, program.Validation, "CreateOptions", true); err != nil {
		return err
	}
	if err := validationgen.Emit(g, program.Schema, program.Client.Validation, program.Client.ConfigMessage, false); err != nil {
		return err
	}
	g.P(`function clone<T>(value: T): T { return structuredClone(value); }`)
	g.P(`function operationSignal(options: OperationOptions, inheritedTimeoutMs?: number): { signal?: AbortSignal; dispose(): void } {`)
	g.P(`  const timeoutMs = options.timeoutMs ?? inheritedTimeoutMs; if (timeoutMs !== undefined && (!Number.isFinite(timeoutMs) || timeoutMs < 0)) throw sandboxError({ kind: ErrorKind.InvalidArgument, provider: "", operation: "", field: "timeoutMs", message: "sandbox-kit: timeoutMs must be a finite non-negative number" });`)
	g.P(`  if ((timeoutMs === undefined || timeoutMs === 0) && options.signal === undefined) return { dispose() {} };`)
	g.P(`  const controller = new AbortController(); const abort = () => controller.abort(options.signal?.reason); if (options.signal?.aborted) abort(); else options.signal?.addEventListener("abort", abort, { once: true });`)
	g.P(`  const timer = timeoutMs !== undefined && timeoutMs > 0 ? setTimeout(() => controller.abort(new DOMException("The operation timed out", "TimeoutError")), timeoutMs) : undefined;`)
	g.P(`  return { signal: controller.signal, dispose() { if (timer !== undefined) clearTimeout(timer); options.signal?.removeEventListener("abort", abort); } };`)
	g.P(`}`)
	g.P(`function operationOptions(options: OperationOptions, signal: AbortSignal | undefined): OperationOptions { return { ...(options.timeoutMs === undefined ? {} : { timeoutMs: options.timeoutMs }), ...(signal === undefined ? {} : { signal }) }; }`)
	g.P(`function abortReason(signal: AbortSignal): DOMException { const reason = signal.reason; return reason instanceof DOMException && (reason.name === "AbortError" || reason.name === "TimeoutError") ? reason : new DOMException("The operation was aborted", "AbortError"); }`)
	g.P(`async function awaitWithSignal<T>(operation: () => Promise<T>, signal?: AbortSignal): Promise<T> { if (signal?.aborted) throw abortReason(signal); const pending = operation(); if (signal === undefined) return pending; return new Promise<T>((resolve, reject) => { const aborted = () => reject(abortReason(signal)); const cleanup = () => signal.removeEventListener("abort", aborted); signal.addEventListener("abort", aborted, { once: true }); pending.then((value) => { cleanup(); resolve(value); }, (error: unknown) => { cleanup(); reject(error); }); }); }`)
	g.P()
	g.P(`export class Sandbox {`)
	g.P(`  readonly #backend: Backend; readonly #timeoutMs: number | undefined; #info: SandboxInfo; #closed = false;`)
	g.P(`  constructor(backend: Backend, info: SandboxInfo, timeoutMs?: number) { this.#backend = backend; this.#info = clone(info); this.#timeoutMs = timeoutMs; }`)
	g.P(`  get id(): string { return this.#info.id; } get providerName(): string { return this.#info.provider; } info(): SandboxInfo { return clone(this.#info); }`)
	g.P(`  capabilities(): readonly LifecycleCapability[] { return this.#backend.capabilities?.(clone(this.#info)) ?? []; }`)
	g.P(`  supports(capability: LifecycleCapability): boolean { return this.capabilities().includes(capability); }`)
	if err := lifecyclegen.EmitHandleMethods(g, program.Lifecycle); err != nil {
		return err
	}
	g.P(`  async close(options: OperationOptions = {}): Promise<void> { if (this.#closed) return; let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, this.#timeoutMs); const close = this.#backend.closeSandbox; if (close) await awaitWithSignal(() => close.call(this.#backend, this.id, operationOptions(options, scoped?.signal)), scoped.signal); this.#closed = true; } catch (error) { throw withContext(error, this.providerName, "close"); } finally { scoped?.dispose(); } }`)
	g.P(`  async [Symbol.asyncDispose](): Promise<void> { await this.close(); }`)
	updateOptions, err := lifecyclegen.UpdateOptionsUnion(program.Lifecycle)
	if err != nil {
		return err
	}
	g.P(`  async #update<T extends `, updateOptions, `>(operation: string, method: ((id: string, request: T, options: OperationOptions) => Promise<SandboxInfo>) | undefined, request: T, options: OperationOptions): Promise<SandboxInfo> { this.#ensureOpen(operation); if (!method) unsupported(this.providerName, operation); let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, request.timeout ?? this.#timeoutMs); const info = await awaitWithSignal(() => method.call(this.#backend, this.id, clone(request), operationOptions(options, scoped?.signal)), scoped.signal); this.#info = clone({ ...info, provider: this.providerName }); return this.info(); } catch (error) { throw withContext(error, this.providerName, operation); } finally { scoped?.dispose(); } }`)
	g.P(`  #ensureOpen(operation: string): void { if (this.#closed) throw sandboxError({ kind: ErrorKind.FailedPrecondition, provider: this.providerName, operation, field: "sandbox", message: "sandbox-kit: sandbox handle is closed" }); }`)
	g.P(`}`)
	g.P()
	g.P(`export class Client implements AsyncDisposable {`)
	g.P(`  readonly #backend: Backend; readonly #timeoutMs: number | undefined; #active = 0; #closing = false; #closed = false; #closePromise: Promise<void> | undefined; #drain: Promise<void> | undefined; #resolveDrain: (() => void) | undefined;`)
	g.P(`  constructor(config: ClientConfig) { if (!config.provider?.name?.trim()) throw sandboxError({ kind: ErrorKind.InvalidArgument, provider: "", operation: "initialize", field: "provider", message: "sandbox-kit: provider is required" }); validateConfig(config, config.provider.name, "initialize"); this.#timeoutMs = config.timeout; try { this.#backend = config.provider.initialize(cloneConfig(config)); } catch (error) { throw withContext(error, config.provider.name, "initialize"); } if (!this.#backend || this.#backend.name !== config.provider.name) throw sandboxError({ kind: ErrorKind.InvalidResponse, provider: config.provider.name, operation: "initialize", field: "provider", message: "sandbox-kit: initialized backend identity differs from selected provider" }); }`)
	g.P(`  get providerName(): string { return this.#backend.name; }`)
	g.P(`  async attach(id: string, options: OperationOptions = {}): Promise<Sandbox> { const release = this.#begin("attach"); let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, this.#timeoutMs); if (!id.trim()) throw sandboxError({ kind: ErrorKind.InvalidArgument, provider: this.providerName, operation: "attach", field: "id", message: "sandbox-kit: sandbox ID must not be blank" }); if (!this.#backend.attach) unsupported(this.providerName, "attach"); const info = await awaitWithSignal(() => this.#backend.attach!(id, operationOptions(options, scoped?.signal)), scoped.signal); return new Sandbox(this.#backend, { ...info, id, provider: this.providerName }, this.#timeoutMs); } catch (error) { throw withContext(error, this.providerName, "attach"); } finally { scoped?.dispose(); release(); } }`)
	g.P(`  async get(id: string, options: OperationOptions = {}): Promise<Sandbox> { const release = this.#begin("get"); let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, this.#timeoutMs); if (!id.trim()) throw sandboxError({ kind: ErrorKind.InvalidArgument, provider: this.providerName, operation: "get", field: "id", message: "sandbox-kit: sandbox ID must not be blank" }); if (!this.#backend.refresh) unsupported(this.providerName, "get"); const info = await awaitWithSignal(() => this.#backend.refresh!(id, operationOptions(options, scoped?.signal)), scoped.signal); return new Sandbox(this.#backend, { ...info, id, provider: this.providerName }, this.#timeoutMs); } catch (error) { throw withContext(error, this.providerName, "get"); } finally { scoped?.dispose(); release(); } }`)
	g.P(`  async create(request: CreateOptions = {}, options: OperationOptions = {}): Promise<Sandbox> { const release = this.#begin("create"); const timeoutMs = options.timeoutMs ?? request.provisioning?.timeout ?? this.#timeoutMs; let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, timeoutMs); validateCreateOptions(request, this.providerName, "create"); const result = await awaitWithSignal(() => this.#backend.create(clone(request), operationOptions({ ...(timeoutMs === undefined ? {} : { timeoutMs }) }, scoped?.signal)), scoped.signal); if (!result.sandbox?.id) throw sandboxError({ kind: ErrorKind.InvalidResponse, provider: this.providerName, operation: "create", field: "id", message: "sandbox-kit: provider returned no sandbox identity" }); if (result.sandbox.provider && result.sandbox.provider !== this.providerName) throw sandboxError({ kind: ErrorKind.InvalidResponse, provider: this.providerName, operation: "create", field: "provider", message: "sandbox-kit: response provider differs from selected provider" }); return new Sandbox(this.#backend, { ...result.sandbox, provider: this.providerName }, this.#timeoutMs); } catch (error) { throw withContext(error, this.providerName, "create"); } finally { scoped?.dispose(); release(); } }`)
	g.P(`  close(options: OperationOptions = {}): Promise<void> { if (this.#closePromise) return this.#closePromise; this.#closing = true; this.#closePromise = this.#close(options); return this.#closePromise; }`)
	g.P(`  async #close(options: OperationOptions): Promise<void> { let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, this.#timeoutMs); if (this.#active > 0) { this.#drain ??= new Promise<void>((resolve) => { this.#resolveDrain = resolve; }); await awaitWithSignal(() => this.#drain!, scoped.signal); } await awaitWithSignal(() => this.#backend.close(operationOptions(options, scoped?.signal)), scoped.signal); this.#closed = true; } catch (error) { this.#closePromise = undefined; this.#closing = false; throw withContext(error, this.providerName, "close"); } finally { scoped?.dispose(); } }`)
	g.P(`  async [Symbol.asyncDispose](): Promise<void> { await this.close(); }`)
	g.P(`  #begin(operation: string): () => void { if (this.#closing || this.#closed) throw sandboxError({ kind: ErrorKind.FailedPrecondition, provider: this.providerName, operation, field: "client", message: "sandbox-kit: client is closed" }); this.#active++; let released = false; return () => { if (released) return; released = true; this.#active--; if (this.#active === 0) { this.#resolveDrain?.(); this.#resolveDrain = undefined; this.#drain = undefined; } }; }`)
	g.P(`}`)
	g.P()
	g.P(`function cloneConfig(config: ClientConfig): Config { const { provider: _provider, ...shared } = config; return clone(shared); }`)
	g.P(`function unsupported(provider: string, operation: string): never { throw sandboxError({ kind: ErrorKind.Unsupported, provider, operation, field: "", message: "sandbox-kit: provider does not support " + operation }); }`)
	g.P(`function withContext(error: unknown, provider: string, operation: string): SandboxError { if (error instanceof SandboxError) { if (error.provider === provider && error.operation === operation) return error; return sandboxError({ kind: error.kind, provider: error.provider || provider, operation: error.operation || operation, field: error.field, message: error.message, ...(error.statusCode === undefined ? {} : { statusCode: error.statusCode }), ...(error.providerCode === undefined ? {} : { providerCode: error.providerCode }), ...(error.providerSource === undefined ? {} : { providerSource: error.providerSource }) }, error.cause); } if (error instanceof DOMException && error.name === "AbortError") return sandboxError({ kind: ErrorKind.Canceled, provider, operation, field: "", message: error.message }, error); if (error instanceof DOMException && error.name === "TimeoutError") return sandboxError({ kind: ErrorKind.Timeout, provider, operation, field: "", message: error.message }, error); return sandboxError({ kind: ErrorKind.Unknown, provider, operation, field: "", message: error instanceof Error ? error.message : "sandbox-kit: provider operation failed" }, error); }`)
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
