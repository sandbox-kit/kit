import { cn } from "@/lib/cn";

/** White square badge with a black isometric cube. */
export function LogoMark({
  className,
  title,
}: {
  className?: string;
  title?: string;
}) {
  return (
    <svg
      viewBox="0 0 32 32"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={cn(
        "block size-7 shrink-0 overflow-hidden rounded-[0.45rem]",
        className,
      )}
      aria-hidden={title ? undefined : true}
      role={title ? "img" : undefined}
    >
      {title ? <title>{title}</title> : null}
      <rect width="32" height="32" rx="5.5" fill="#ffffff" />
      {/* Cube sized like the original mark — room around it on the white badge */}
      <path fill="#111110" d="M16 6.2 25.4 11.5 16 16.8 6.6 11.5 16 6.2Z" />
      <path fill="#111110" d="M6.6 12.35 15.35 17.35v9.4L6.6 21.75V12.35Z" />
      <path fill="#111110" d="M25.4 12.35 16.65 17.35v9.4L25.4 21.75V12.35Z" />
      <path
        stroke="#ffffff"
        strokeWidth="1"
        strokeLinecap="round"
        d="M16 16.8V26.75M16 16.8 6.6 12.35M16 16.8 25.4 12.35"
      />
    </svg>
  );
}
