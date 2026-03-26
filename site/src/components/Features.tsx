const features = [
  {
    title: "Run it yourself",
    description:
      "Interactive CLI with numbered menus. No subcommands to memorize. Run press and pick an action.",
  },
  {
    title: "Script it",
    description:
      "Every action has flags for non-interactive mode. Drop it into CI/CD, cron jobs, or shell scripts.",
  },
  {
    title: "Hand it to your AI agent",
    description:
      "Structured commands and predictable output make pressctl a natural fit for AI agents like Claude Code or Cursor.",
  },
  {
    title: "Production-grade stack",
    description:
      "Nginx, PHP 8.3, MariaDB, Redis, Let's Encrypt, UFW, Fail2ban — on Ubuntu 24.04. Battle-tested Ansible under the hood.",
  },
];

export default function Features() {
  return (
    <section className="px-6 py-24 border-t border-zinc-800/50">
      <div className="max-w-5xl mx-auto">
        <h2 className="text-2xl sm:text-3xl font-bold text-white text-center mb-4">
          Your workflow. Your way.
        </h2>
        <p className="text-zinc-500 text-center mb-16 max-w-xl mx-auto">
          Terminal, automation, or AI — pressctl fits however you work.
        </p>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-px bg-zinc-800/50 rounded-xl overflow-hidden border border-zinc-800">
          {features.map((feature) => (
            <div
              key={feature.title}
              className="p-8 bg-zinc-950"
            >
              <h3 className="text-lg font-semibold text-white mb-2">
                {feature.title}
              </h3>
              <p className="text-sm text-zinc-400 leading-relaxed">
                {feature.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
