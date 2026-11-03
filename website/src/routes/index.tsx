import { createFileRoute, Link } from "@tanstack/react-router";
import { HomeLayout } from "fumadocs-ui/layouts/home";
import { baseOptions } from "@/lib/layout.shared";
import { assetUrl } from "@/lib/shared";

export const Route = createFileRoute("/")({
  component: Home,
});

function Home() {
  return (
    <HomeLayout {...baseOptions()}>
      <main className="mx-auto flex w-full max-w-3xl flex-1 flex-col px-6 pt-16 pb-24 md:pt-24">
        <img
          src={assetUrl("logo.png")}
          alt=""
          width={96}
          height={96}
          className="h-24 w-24"
        />
        <h1 className="mt-8 max-w-xl text-4xl font-medium tracking-tight text-balance md:text-5xl">
          One interface for your sandbox providers.
        </h1>
        <p className="mt-5 max-w-xl text-lg leading-relaxed text-fd-muted-foreground">
          Sandbox Kit is the client an agent uses with a cloud sandbox. Pass in
          a provider. The same client opens the sandbox where that agent writes
          files, runs commands, and finishes its work.
        </p>
        <p className="mt-8 text-sm leading-relaxed text-fd-muted-foreground">
          <Link
            to="/docs/$"
            params={{ _splat: "" }}
            className="text-fd-foreground underline decoration-fd-border underline-offset-4 hover:decoration-fd-foreground"
          >
            Read the guide
          </Link>
          . Open the page for the call you need. TypeScript and Go samples stay
          on the language you pick.
        </p>
      </main>
    </HomeLayout>
  );
}
