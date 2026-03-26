import CopyButton from "@/components/ui/CopyButton";
import GitHubIcon from "@/components/ui/GitHubIcon";

const INSTALL_CMD =
  "curl -fsSL https://raw.githubusercontent.com/pressctl/cli/main/install.sh | bash";

export default function HowItWorks() {
  return (
    <section className="px-6 py-24 border-t border-zinc-800/50">
      <div className="max-w-3xl mx-auto text-center">
        <h2 className="text-2xl sm:text-3xl font-bold text-white mb-4">
          Get started in seconds
        </h2>
        <p className="text-zinc-500 mb-10 max-w-lg mx-auto">
          Install pressctl, point it at any Ubuntu 24.04 VPS, and you&apos;re hosting WordPress.
        </p>

        {/* Install command */}
        <div className="flex items-start gap-2 max-w-xl mx-auto mb-8 px-4 py-3 rounded-lg bg-zinc-900 border border-zinc-800 font-mono text-sm text-left">
          <span className="text-accent select-none mt-0.5">$</span>
          <code className="text-zinc-300 flex-1 break-all whitespace-pre-wrap">
            {INSTALL_CMD}
          </code>
          <CopyButton text={INSTALL_CMD} />
        </div>

        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <a
            href="https://github.com/shariffff/pressctl"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-6 py-3 text-sm font-medium text-white bg-accent/10 border border-accent/20 rounded-lg hover:bg-accent/20 transition-all"
          >
            <GitHubIcon className="w-4 h-4" />
            View on GitHub
          </a>
          <a
            href="/docs"
            className="inline-flex items-center gap-2 px-6 py-3 text-sm font-medium text-zinc-400 hover:text-white transition-colors"
          >
            Read the docs
          </a>
        </div>
      </div>
    </section>
  );
}
