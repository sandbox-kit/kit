import {
  ErrorKind,
  CreateFieldEnvironment,
  validEnvName,
} from "@sandbox-kit/sandbox";
import type { CreateOptions } from "@sandbox-kit/sandbox";
import { providerError } from "#errors.js";

// validateProviderCreate rejects creation settings this provider cannot honor.
// Rules are inventoried in specs/providers/blaxel.yaml checks when present.
export function validateProviderCreate(request: CreateOptions): void {
  for (const value of Object.keys(request.environment ?? {})) {
    if (!validEnvName(value)) {
      providerError({
        kind: ErrorKind.InvalidArgument,
        operation: "create",
        field: CreateFieldEnvironment,
        message: "environment is outside supported values",
      });
    }
  }
}
