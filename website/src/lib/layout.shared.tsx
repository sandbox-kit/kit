import type { BaseLayoutProps } from "fumadocs-ui/layouts/shared";
import { appName, assetUrl, gitConfig } from "./shared";

export function baseOptions(): BaseLayoutProps {
  return {
    nav: {
      title: (
        <span className="inline-flex items-center gap-2 font-medium">
          <img
            src={assetUrl("logo.png")}
            alt=""
            width={22}
            height={22}
            className="h-[22px] w-[22px]"
          />
          {appName}
        </span>
      ),
    },
    githubUrl: `https://github.com/${gitConfig.user}/${gitConfig.repo}`,
  };
}
