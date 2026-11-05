"use client";

import type { ReactNode } from "react";
import { cn } from "@/lib/cn";

/**
 * Compact language/provider-aware note. Content swaps with the page language
 * tabs and provider select — this box has no switchers of its own.
 */
export function Note({
  children,
  className,
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <aside
      className={cn(
        "my-4 rounded-xl border border-fd-border bg-fd-secondary/40 px-4 py-3 text-sm text-fd-muted-foreground",
        "**:first:mt-0 **:last:mb-0",
        "[&_p]:my-0 [&_code]:text-[0.9em]",
        className,
      )}
    >
      {children}
    </aside>
  );
}
