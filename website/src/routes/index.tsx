import { createFileRoute, Link } from "@tanstack/react-router";
import { HomeLayout } from "fumadocs-ui/layouts/home";
import { HomeCodeSample } from "@/components/home-code-sample";
import { baseOptions } from "@/lib/layout.shared";
import { siteSocialHead } from "@/lib/seo";

export const Route = createFileRoute("/")({
  head: () => siteSocialHead(),
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

        <div className="relative mx-auto grid min-h-0 w-full max-w-[1400px] flex-1 grid-cols-1 items-stretch gap-3 px-4 py-3 max-[699px]:grid-rows-[auto_minmax(0,1fr)] min-[700px]:grid-cols-[minmax(0,0.92fr)_minmax(0,1.08fr)] min-[700px]:gap-10 min-[700px]:px-10 min-[700px]:py-8 sm:gap-5 sm:px-6 sm:py-5">
          <section className="home-copy flex min-h-0 flex-col justify-center py-1 min-[700px]:py-2">
            <h1 className="home-brand text-[clamp(1.75rem,4.6vw,3.75rem)] font-medium leading-none tracking-[-0.04em] text-fd-foreground">
              sandbox-kit
            </h1>

            <p className="home-lede mt-4 max-w-[34rem] text-[clamp(0.9rem,1.35vw,1.125rem)] leading-relaxed text-fd-muted-foreground sm:mt-5">
              SDK for cloud sandboxes in TypeScript or Go. Plug in a provider;
              the client talks to it for create, files, and commands.
            </p>

            <div className="home-cta mt-5 flex flex-wrap items-center gap-3 sm:mt-7">
              <Link
                to="/docs/$"
                params={{ _splat: "" }}
                className="rounded-md bg-fd-primary px-4 py-2.5 text-sm font-medium text-fd-primary-foreground transition-[opacity,transform] duration-200 hover:opacity-90 active:scale-[0.98]"
              >
                Quickstart
              </Link>
              <Link
                to="/docs/$"
                params={{ _splat: "concepts" }}
                className="rounded-md border border-fd-border bg-transparent px-4 py-2.5 text-sm font-medium text-fd-foreground transition-colors duration-200 hover:bg-fd-accent"
              >
                Concepts
              </Link>
            </div>

            <p className="home-meta mt-6 hidden font-mono text-[11px] tracking-[0.14em] text-fd-muted-foreground uppercase min-[700px]:mt-9 min-[700px]:block">
              One client · many providers
            </p>
          </section>

          <section className="home-sample flex min-h-0 min-w-0 flex-col justify-center min-[700px]:ps-2">
            <p className="mb-2 shrink-0 font-mono text-[11px] tracking-[0.14em] text-fd-muted-foreground uppercase">
              Create, then run
            </p>
            <div className="home-code min-h-0 w-full shrink">
              <HomeCodeSample />
            </div>
          </section>
        </div>
      </div>
    </HomeLayout>
  );
}
