import { cn } from "@/lib/cn";
import { LogoMark } from "@/components/logo-mark";

/** Nav / sidebar wordmark — white/black cube badge + lowercase name. */
export function BrandMark({
  className,
  logoClassName,
  showLogo = true,
}: {
  className?: string;
  logoClassName?: string;
  showLogo?: boolean;
}) {
  return (
    <span
      className={cn(
        "inline-flex flex-row flex-nowrap items-center gap-2.5 text-[0.95rem] font-medium leading-none tracking-[-0.03em] text-fd-foreground",
        className,
      )}
    >
      {showLogo ? (
        <LogoMark
          className={cn(
            "size-[1.2em] shrink-0 -translate-y-px ring-1 ring-black/10 dark:ring-white/15",
            logoClassName,
          )}
        />
      ) : null}
      <span aria-label="sandbox-kit" className="whitespace-nowrap">
        sandbox-kit
      </span>
    </span>
  );
}
