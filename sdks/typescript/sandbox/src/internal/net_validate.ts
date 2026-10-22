/** Shared network / env name validators used by provider create checks. */

export function validEnvName(value: string): boolean {
  return /^[A-Za-z_][A-Za-z0-9_]*$/.test(value);
}

export function validDomain(value: string): boolean {
  const domain = value.startsWith("*.") ? value.slice(2) : value;
  return (
    domain.length <= 253 &&
    /^(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)*[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.?$/i.test(
      domain,
    )
  );
}

export function validCidr(value: string): boolean {
  const separator = value.lastIndexOf("/");
  if (separator <= 0) return false;
  const address = value.slice(0, separator);
  const prefix = Number(value.slice(separator + 1));
  if (!Number.isInteger(prefix)) return false;
  return address.includes(":")
    ? prefix >= 0 && prefix <= 128 && validIpv6(address)
    : prefix >= 0 && prefix <= 32 && validIpv4(address);
}

function validIpv4(value: string): boolean {
  const parts = value.split(".");
  return (
    parts.length === 4 &&
    parts.every(
      (part) => /^(0|[1-9][0-9]{0,2})$/.test(part) && Number(part) <= 255,
    )
  );
}

function validIpv6(value: string): boolean {
  if ((value.match(/::/g)?.length ?? 0) > 1) return false;
  const halves = value.split("::");
  if (halves.length > 2) return false;
  const count = (part: string): number | undefined => {
    if (part === "") return 0;
    const groups = part.split(":");
    let total = 0;
    for (const group of groups) {
      if (/^[0-9a-f]{1,4}$/i.test(group)) total++;
      else if (validIpv4(group)) total += 2;
      else return undefined;
    }
    return total;
  };
  const left = count(halves[0] ?? "");
  const right = count(halves[1] ?? "");
  return (
    left !== undefined &&
    right !== undefined &&
    (halves.length === 2 ? left + right < 8 : left === 8)
  );
}
