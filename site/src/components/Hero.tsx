"use client";

import CopyButton from "@/components/ui/CopyButton";
import GitHubIcon from "@/components/ui/GitHubIcon";
import posthog from "posthog-js";

const INSTALL_CMD =
  "curl -fsSL https://raw.githubusercontent.com/shariffff/pressctl/main/install.sh | bash";

export default function Hero() {
  return (
    <section className="relative px-6 pt-28 pb-24">
      {/* Subtle grid background */}
      <div
        className="absolute inset-0 opacity-[0.03]"
        style={{
          backgroundImage:
            "linear-gradient(rgba(255,255,255,0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.1) 1px, transparent 1px)",
          backgroundSize: "64px 64px",
        }}
      />

      <div className="relative z-10 max-w-6xl mx-auto grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-16 items-center">
        {/* Text */}
        <div className="animate-fade-in-up">
          <h1 className="text-4xl sm:text-5xl font-bold tracking-tight text-white mb-6">
            WordPress hosting infrastructure{" "}
            <span className="text-accent">you actually control.</span>
          </h1>

          <p className="text-lg text-zinc-400 mb-8 leading-relaxed max-w-lg">
            Run it from your terminal, script it in your pipeline, or hand it to
            your AI agent. No DevOps knowledge required.
          </p>

          {/* Install command */}
          <div className="flex items-start gap-2 mb-8 px-4 py-3 rounded-lg bg-zinc-900 border border-zinc-800 font-mono text-sm text-left">
            <span className="text-accent select-none mt-0.5">$</span>
            <code className="text-zinc-300 flex-1 break-all whitespace-pre-wrap">
              {INSTALL_CMD}
            </code>
            <CopyButton text={INSTALL_CMD} />
          </div>

          {/* CTAs */}
          <div className="flex flex-wrap items-center gap-4">
            <a
              href="https://github.com/shariffff/pressctl"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-6 py-3 text-sm font-medium text-white border border-zinc-700 rounded-lg hover:bg-zinc-800 hover:border-zinc-600 transition-all"
              onClick={() => posthog.capture("github_link_clicked", { source: "hero" })}
            >
              <GitHubIcon className="w-4 h-4" />
              View on GitHub
            </a>
            <a
              href="/docs"
              className="inline-flex items-center gap-2 px-6 py-3 text-sm font-medium text-zinc-400 hover:text-white transition-colors"
              onClick={() => posthog.capture("docs_link_clicked", { source: "hero" })}
            >
              Read the docs
            </a>
          </div>
        </div>

        {/* Video */}
        <div className="animate-fade-in-up delay-200">
          <div className="rounded-xl border border-zinc-800 overflow-hidden shadow-2xl shadow-black/50">
            <div
              className="relative w-full"
              style={{ paddingBottom: "56.25%" }}
            >
              <iframe
                className="absolute inset-0 w-full h-full"
                src="https://www.youtube.com/embed/oZmqxB_eMMY"
                title="pressctl demo"
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                allowFullScreen
              />
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
