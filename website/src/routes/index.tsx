import { createFileRoute, Link } from "@tanstack/react-router";
import { HomeLayout } from "fumadocs-ui/layouts/home";
import { HomeCodeSample } from "@/components/home-code-sample";
import { baseOptions } from "@/lib/layout.shared";
import { assetUrl } from "@/lib/shared";

export const Route = createFileRoute("/")({
  component: Home,
});

function Home() {
  return (
    <HomeLayout
      {...baseOptions()}
      className="home-shell h-dvh max-h-dvh overflow-hidden"
      nav={{
        ...baseOptions().nav,
        transparentMode: "always",
      }}
    >
      <div className="home-atmosphere relative flex min-h-0 flex-1 flex-col overflow-hidden">
        <div
          aria-hidden
          className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_at_12%_18%,rgba(17,17,16,0.05),transparent_42%),radial-gradient(ellipse_at_88%_78%,rgba(47,111,94,0.07),transparent_38%)]"
        />
        <div
          aria-hidden
          className="pointer-events-none absolute inset-0 opacity-[0.35] [background-image:linear-gradient(to_right,rgba(17,17,16,0.045)_1px,transparent_1px),linear-gradient(to_bottom,rgba(17,17,16,0.045)_1px,transparent_1px)] [background-size:48px_48px] dark:opacity-[0.18]"
        />

        <div className="relative mx-auto grid min-h-0 w-full max-w-[1400px] flex-1 grid-cols-1 items-stretch gap-3 px-4 py-3 max-[699px]:grid-rows-[auto_minmax(0,1fr)] min-[700px]:grid-cols-[minmax(0,0.88fr)_minmax(0,1.12fr)] min-[700px]:gap-8 min-[700px]:px-10 min-[700px]:py-8 sm:gap-5 sm:px-6 sm:py-5">
          <section className="home-copy flex min-h-0 flex-col justify-center py-1 min-[700px]:py-2">
            <img
              src={assetUrl("logo.png")}
              alt=""
              width={56}
              height={56}
              className="h-9 w-9 sm:h-12 sm:w-12 min-[700px]:h-14 min-[700px]:w-14"
            />
            <h1 className="mt-3 max-w-[14ch] text-[clamp(1.85rem,5vw,4.25rem)] font-medium leading-[0.95] tracking-tight text-fd-foreground sm:mt-5 min-[700px]:mt-6">
              Sandbox Kit
            </h1>
            <p className="mt-3 max-w-md text-[clamp(0.85rem,1.4vw,1.125rem)] leading-relaxed text-fd-muted-foreground sm:mt-4 min-[700px]:mt-5">
              Lightweight TypeScript and Go SDK for cloud sandboxes. Plug in a
              provider—one client for create, files, and commands.
            </p>
            <div className="mt-4 flex flex-wrap items-center gap-3 sm:mt-6 min-[700px]:mt-8">
              <Link
                to="/docs/$"
                params={{ _splat: "" }}
                className="rounded-md bg-fd-primary px-4 py-2.5 text-sm font-medium text-fd-primary-foreground transition-opacity hover:opacity-90"
              >
                Quickstart
              </Link>
              <Link
                to="/docs/$"
                params={{ _splat: "concepts" }}
                className="rounded-md border border-fd-border bg-transparent px-4 py-2.5 text-sm font-medium text-fd-foreground transition-colors hover:bg-fd-accent"
              >
                Concepts
              </Link>
            </div>
            <p className="mt-5 hidden font-mono text-[11px] tracking-wide text-fd-muted-foreground uppercase min-[700px]:mt-8 min-[700px]:block">
              Daytona · Modal · your provider next
            </p>
          </section>

          <section className="home-sample flex min-h-0 min-w-0 flex-col justify-center min-[700px]:ps-2">
            <p className="mb-2 shrink-0 font-mono text-[11px] tracking-wide text-fd-muted-foreground uppercase">
              Same call shape on every provider
            </p>
            <div className="min-h-0 w-full shrink">
              <HomeCodeSample />
            </div>
          </section>
        </div>
      </div>
    </HomeLayout>
  );
}
