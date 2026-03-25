import { Server, Globe, Lock, TerminalSquare } from "lucide-react";

const features = [
  {
    icon: Server,
    title: "Server Provisioning",
    description:
      "Full LEMP stack in one command — Nginx, PHP 8.3, MariaDB, Redis, Certbot, UFW, and Fail2ban on Ubuntu 24.04.",
  },
  {
    icon: Globe,
    title: "WordPress Deployment",
    description:
      "Isolated sites with dedicated Linux users, PHP-FPM pools, databases, and WP-CLI — no shared hosting footguns.",
  },
  {
    icon: Lock,
    title: "SSL & Domain Management",
    description:
      "Add domains and issue Let's Encrypt certificates interactively. No certbot flags to memorize.",
  },
  {
    icon: TerminalSquare,
    title: "Interactive Menu",
    description:
      "Run `press` with no arguments and pick from a numbered action menu. No subcommands to look up.",
  },
];

export default function Features() {
  return (
    <section className="px-6 py-24 max-w-5xl mx-auto">
      <h2 className="text-2xl sm:text-3xl font-mono font-bold text-white text-center mb-4">
        Everything you need to host WordPress
      </h2>
      <p className="text-zinc-500 text-center mb-14 max-w-xl mx-auto">
        Production-ready infrastructure from a single CLI. Powered by Ansible
        playbooks under the hood.
      </p>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-5">
        {features.map((feature) => (
          <div
            key={feature.title}
            className="group p-6 rounded-xl border border-zinc-800 bg-zinc-900/50 hover:bg-zinc-900 hover:border-zinc-700 transition-all"
          >
            <feature.icon className="w-8 h-8 text-terminal-green mb-4 group-hover:scale-110 transition-transform" />
            <h3 className="text-lg font-mono font-semibold text-white mb-2">
              {feature.title}
            </h3>
            <p className="text-sm text-zinc-400 leading-relaxed">
              {feature.description}
            </p>
          </div>
        ))}
      </div>
    </section>
  );
}
