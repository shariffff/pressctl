"use client";

import Link from "next/link";
import GitHubIcon from "@/components/ui/GitHubIcon";
import posthog from "posthog-js";

export default function Footer() {
  return (
    <footer className="px-6 py-12 border-t border-zinc-800/50">
      <div className="max-w-5xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <Link href="/" className="font-mono font-bold text-white text-sm hover:opacity-80 transition-opacity">
            press<span className="text-accent">ctl</span>
          </Link>
          <span className="text-zinc-600 text-xs">MIT License</span>
        </div>

        <div className="flex items-center gap-6">
          <Link
            href="/docs"
            className="text-sm text-zinc-500 hover:text-white transition-colors"
          >
            Documentation
          </Link>
          <a
            href="https://github.com/shariffff/pressctl"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 text-sm text-zinc-500 hover:text-white transition-colors"
            onClick={() => posthog.capture("github_link_clicked", { source: "footer" })}
          >
            <GitHubIcon className="w-4 h-4" />
            GitHub
          </a>
        </div>
      </div>
    </footer>
  );
}
