const terminals = [
  {
    title: "cat /etc/nginx/sites-available/example.com/example.com",
    lines: [
      { text: "server {", color: "text-zinc-500" },
      { text: "    listen 80;", color: "text-zinc-500" },
      { text: "    server_name example.com;", color: "text-terminal-green" },
      { text: "    root /sites/example.com/files;", color: "text-terminal-green" },
      { text: "", color: "" },
      { text: "    # modular per-site includes", color: "text-zinc-600" },
      {
        text: "    include sites-available/example.com/server/*;",
        color: "text-white",
      },
      { text: "", color: "" },
      { text: "    location ~ \\.php$ {", color: "text-zinc-500" },
      {
        text: "        fastcgi_pass unix:/run/php/php8.3-examplecom.sock;",
        color: "text-terminal-green",
      },
      {
        text: "        include sites-available/example.com/location/*;",
        color: "text-white",
      },
      { text: "    }", color: "text-zinc-500" },
      { text: "}", color: "text-zinc-500" },
    ],
  },
  {
    title: "cat /etc/php/8.3/fpm/pool.d/examplecom.conf",
    lines: [
      { text: "[examplecom]", color: "text-white" },
      { text: "user = examplecom", color: "text-terminal-green" },
      { text: "group = examplecom", color: "text-terminal-green" },
      { text: "", color: "" },
      {
        text: "listen = /run/php/php8.3-examplecom.sock",
        color: "text-terminal-green",
      },
      { text: "listen.owner = examplecom", color: "text-terminal-green" },
      { text: "listen.group = www-data", color: "text-zinc-500" },
      { text: "", color: "" },
      { text: "pm = dynamic", color: "text-white" },
      { text: "pm.max_children = 5", color: "text-zinc-500" },
      { text: "pm.start_servers = 1", color: "text-zinc-500" },
      { text: "pm.max_requests = 500", color: "text-zinc-500" },
    ],
  },
  {
    title: "tree /sites/example.com/",
    lines: [
      { text: "/sites/example.com/", color: "text-terminal-green" },
      { text: "├── files/              ← document root", color: "text-white" },
      { text: "│   ├── wp-config.php", color: "text-zinc-500" },
      { text: "│   ├── wp-content/", color: "text-zinc-500" },
      { text: "│   └── index.php", color: "text-zinc-500" },
      { text: "├── logs/", color: "text-white" },
      { text: "│   ├── access.log", color: "text-zinc-500" },
      { text: "│   └── error.log", color: "text-zinc-500" },
      { text: "├── .local/bin/", color: "text-zinc-500" },
      { text: "│   └── php → /usr/bin/php8.3", color: "text-zinc-500" },
      {
        text: "└── .pressctl-pool.conf  ← custom overrides",
        color: "text-white",
      },
    ],
  },
  {
    title: 'mysql -e "SHOW GRANTS FOR examplecom@localhost"',
    lines: [
      {
        text: "GRANT ALL PRIVILEGES ON `examplecom`.*",
        color: "text-white",
      },
      {
        text: "  TO 'examplecom'@'localhost'",
        color: "text-terminal-green",
      },
      { text: "", color: "" },
      { text: "# one database, one user, localhost only", color: "text-zinc-600" },
      { text: "# credentials generated at site creation", color: "text-zinc-600" },
    ],
  },
];

function Terminal({ title, lines }: (typeof terminals)[number]) {
  return (
    <div className="rounded-xl border border-zinc-800 bg-zinc-900 overflow-hidden shadow-2xl shadow-black/50">
      <div className="flex items-center gap-2 bg-zinc-800/80 px-4 py-3">
        <div className="flex gap-1.5">
          <div className="w-3 h-3 rounded-full bg-red-500/80" />
          <div className="w-3 h-3 rounded-full bg-yellow-500/80" />
          <div className="w-3 h-3 rounded-full bg-green-500/80" />
        </div>
        <span className="ml-2 text-xs text-zinc-500 font-mono truncate">
          bash
        </span>
      </div>
      <div className="p-5 font-mono text-xs sm:text-sm leading-relaxed overflow-x-auto">
        <div className="mb-3">
          <span className="text-terminal-green">$</span>{" "}
          <span className="text-zinc-300">{title}</span>
        </div>
        {lines.map((line, i) => (
          <div key={i} className={`${line.color} whitespace-pre`}>
            {line.text || "\u00A0"}
          </div>
        ))}
      </div>
    </div>
  );
}

export default function Architecture() {
  return (
    <section className="px-6 py-24 border-t border-zinc-800/50">
      <div className="max-w-5xl mx-auto">
        <h2 className="text-2xl sm:text-3xl font-bold text-white text-center mb-4">
          What&apos;s on the server
        </h2>
        <p className="text-zinc-500 text-center mb-16 max-w-xl mx-auto">
          Every site gets its own Linux user, PHP-FPM pool, Nginx config, and
          database. No shared resources. Here&apos;s exactly what that looks like.
        </p>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {terminals.map((t) => (
            <Terminal key={t.title} {...t} />
          ))}
        </div>
      </div>
    </section>
  );
}
