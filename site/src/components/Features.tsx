const features = [
  "Isolated Linux user and group per site",
  "Dedicated PHP-FPM pool and socket per site",
  "Separate Nginx config with modular includes",
  "Per-site database with single-database privileges",
  "UFW firewall — ports 22, 80, 443 only",
  "Fail2ban brute-force protection",
  "Let's Encrypt auto SSL certificates",
  "System cron with file-locked WP-Cron per site",
  "Redis object cache",
  "FastCGI page cache per site",
  "Interactive CLI or fully scriptable flags",
  "AI-agent friendly — structured output, predictable commands",
];

export default function Features() {
  return (
    <section className="px-6 py-24 border-t border-zinc-800/50">
      <div className="max-w-3xl mx-auto">
        <h2 className="text-2xl sm:text-3xl font-bold text-white text-center mb-4">
          What you get out of the box
        </h2>
        <p className="text-zinc-500 text-center mb-14 max-w-xl mx-auto">
          Production-grade WordPress on Ubuntu 24.04. Nginx, PHP 8.3, MariaDB,
          Redis, and battle-tested Ansible under the hood.
        </p>

        <ul className="grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-4">
          {features.map((feature) => (
            <li key={feature} className="flex items-start gap-3">
              <span className="text-accent font-mono flex-shrink-0">—</span>
              <span className="text-sm text-zinc-300">{feature}</span>
            </li>
          ))}
        </ul>
      </div>
    </section>
  );
}
