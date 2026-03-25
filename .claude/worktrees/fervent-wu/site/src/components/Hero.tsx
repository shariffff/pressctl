import CopyButton from "@/components/ui/CopyButton";
import GitHubIcon from "@/components/ui/GitHubIcon";
import { ArrowDown } from "lucide-react";

const INSTALL_CMD =
  "curl -fsSL https://raw.githubusercontent.com/pressctl/cli/main/install.sh | bash";

export default function Hero() {
  return (
    <section className="relative flex flex-col items-center justify-center min-h-screen px-6 pt-14 text-center">
      {/* Subtle grid background */}
      <div
        className="absolute inset-0 opacity-[0.03]"
        style={{
          backgroundImage:
            "linear-gradient(rgba(255,255,255,0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.1) 1px, transparent 1px)",
          backgroundSize: "64px 64px",
        }}
      />

      <div className="relative z-10 max-w-3xl mx-auto animate-fade-in-up">
        {/* Badge */}
        <div className="inline-flex items-center gap-2 px-3 py-1 mb-8 text-xs font-mono text-terminal-green border border-terminal-green/20 rounded-full bg-terminal-green/5">
          <span className="w-1.5 h-1.5 rounded-full bg-terminal-green animate-pulse" />
          Open Source CLI Tool
        </div>

        {/* Headline */}
        <h1 className="text-4xl sm:text-5xl md:text-6xl font-bold tracking-tight text-white mb-6">
          WordPress infrastructure
          <br />
          <span className="text-terminal-green">from your terminal</span>
        </h1>

        {/* Subheadline */}
        <p className="text-lg sm:text-xl text-zinc-400 max-w-2xl mx-auto mb-10 leading-relaxed">
          One command to provision a server. One command to deploy a site.
          <br className="hidden sm:block" />
          No YAML editing. No Ansible knowledge required.
        </p>

        {/* Install command — wraps naturally, no horizontal scroll */}
        <div className="flex items-start gap-2 max-w-xl mx-auto mb-10 px-4 py-3 rounded-lg bg-zinc-900 border border-zinc-800 font-mono text-sm text-left">
          <span className="text-terminal-green select-none mt-0.5">$</span>
          <code className="text-zinc-300 flex-1 break-all whitespace-pre-wrap">
            {INSTALL_CMD}
          </code>
          <CopyButton text={INSTALL_CMD} />
        </div>

        {/* CTAs */}
        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <a
            href="https://github.com/shariffff/wpsh"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-6 py-3 text-sm font-medium text-white border border-zinc-700 rounded-lg hover:bg-zinc-800 hover:border-zinc-600 transition-all"
          >
            <GitHubIcon className="w-4 h-4" />
            View on GitHub
          </a>
          <a
            href="/docs"
            className="inline-flex items-center gap-2 px-6 py-3 text-sm font-medium text-zinc-400 hover:text-white transition-colors"
          >
            Read the docs
            <ArrowDown className="w-4 h-4" />
          </a>
        </div>
      </div>
    </section>
  );
}
