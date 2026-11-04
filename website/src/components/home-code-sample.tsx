"use client";

import { useState, type ReactNode } from "react";
import { cn } from "@/lib/cn";

const samples = {
  TypeScript: {
    filename: "create.ts",
    code: `import { daytona } from "@sandbox-kit/daytona";
import { Client } from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
});

const sandbox = await client.create();
const result = await sandbox.process.run({
  argv: ["echo", "hello from kit"],
});

console.log(result.exitCode, result.stdout);`,
  },
  Go: {
    filename: "create.go",
    code: `import (
	"context"
	"os"

	"github.com/sandbox-kit/kit/sdks/go/providers/daytona"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{
	Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{
		ApiKey: &sandbox.APIKeyCredentials{
			Key: os.Getenv("DAYTONA_API_KEY"),
		},
	}},
})
defer client.Close(context.Background())

instance, err := client.Create(context.Background(), nil)
result, err := instance.Process().Run(
	context.Background(),
	&sandbox.ExecOptions{
		Argv: []string{"echo", "hello from kit"},
	},
)`,
  },
} as const;

type Lang = keyof typeof samples;

function highlight(code: string, lang: Lang) {
  const keyword =
    lang === "TypeScript"
      ? /\b(import|from|await|using|const|new|console)\b/g
      : /\b(import|defer|err|nil|NewClient|New|Create|Run|Background|Getenv)\b/g;

  const stringRe = /("(?:\\.|[^"\\])*"|`(?:\\.|[^`\\])*`)/g;
  const commentRe = /(\/\/.*$)/gm;

  const parts: { text: string; kind: "plain" | "string" | "comment" }[] = [];
  let cursor = 0;
  const marks: { start: number; end: number; kind: "string" | "comment" }[] =
    [];

  for (const match of code.matchAll(commentRe)) {
    if (match.index == null) continue;
    marks.push({
      start: match.index,
      end: match.index + match[0].length,
      kind: "comment",
    });
  }
  for (const match of code.matchAll(stringRe)) {
    if (match.index == null) continue;
    const overlaps = marks.some(
      (m) => match.index! < m.end && match.index! + match[0].length > m.start,
    );
    if (!overlaps) {
      marks.push({
        start: match.index,
        end: match.index + match[0].length,
        kind: "string",
      });
    }
  }
  marks.sort((a, b) => a.start - b.start);

  for (const mark of marks) {
    if (mark.start > cursor) {
      parts.push({ text: code.slice(cursor, mark.start), kind: "plain" });
    }
    parts.push({ text: code.slice(mark.start, mark.end), kind: mark.kind });
    cursor = mark.end;
  }
  if (cursor < code.length)
    parts.push({ text: code.slice(cursor), kind: "plain" });

  return parts.map((part, i) => {
    if (part.kind === "string") {
      return (
        <span key={i} className="text-[#2f6f5e]">
          {part.text}
        </span>
      );
    }
    if (part.kind === "comment") {
      return (
        <span key={i} className="text-fd-muted-foreground/70">
          {part.text}
        </span>
      );
    }
    const nodes: ReactNode[] = [];
    let last = 0;
    for (const match of part.text.matchAll(keyword)) {
      if (match.index == null) continue;
      if (match.index > last) nodes.push(part.text.slice(last, match.index));
      nodes.push(
        <span key={`${i}-${match.index}`} className="text-[#8a4b2e]">
          {match[0]}
        </span>,
      );
      last = match.index + match[0].length;
    }
    if (last < part.text.length) nodes.push(part.text.slice(last));
    return <span key={i}>{nodes}</span>;
  });
}

export function HomeCodeSample() {
  const [lang, setLang] = useState<Lang>("TypeScript");
  const sample = samples[lang];
  const lines = sample.code.split("\n");

  return (
    <div className="home-code flex max-h-full w-full flex-col overflow-hidden border border-fd-border bg-[#fbfbf6] dark:bg-fd-card">
      <div className="flex shrink-0 items-center gap-1 border-b border-fd-border px-3 py-2">
        {(Object.keys(samples) as Lang[]).map((item) => (
          <button
            key={item}
            type="button"
            onClick={() => setLang(item)}
            className={cn(
              "rounded px-2.5 py-1 text-xs font-medium transition-colors",
              lang === item
                ? "bg-fd-foreground text-fd-background"
                : "text-fd-muted-foreground hover:text-fd-foreground",
            )}
          >
            {item}
          </button>
        ))}
        <span className="ms-auto font-mono text-[11px] text-fd-muted-foreground">
          {sample.filename}
        </span>
      </div>
      <pre className="overflow-auto p-4 font-mono text-[12px] leading-[1.55] text-fd-foreground sm:text-[13px]">
        <code>
          {lines.map((line, index) => (
            <div key={index} className="flex gap-4">
              <span className="w-5 shrink-0 select-none text-right text-fd-muted-foreground/50">
                {index + 1}
              </span>
              <span className="min-w-0 whitespace-pre">
                {highlight(line || " ", lang)}
              </span>
            </div>
          ))}
        </code>
      </pre>
    </div>
  );
}
