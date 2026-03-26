import TerminalWindow from "@/components/ui/TerminalWindow";

const commands = [
  {
    title: "press server provision",
    lines: [
      { text: "$ press server provision myserver", style: "text-white" },
      { text: "", style: "" },
      { text: "Checking SSH connectivity...", style: "text-zinc-400" },
      { text: "\u2713 SSH connectivity check passed", style: "text-terminal-green" },
      { text: "\u2713 No port conflicts detected", style: "text-terminal-green" },
      { text: "", style: "" },
      { text: "\u2550\u2550\u2550 Starting provisioning: myserver \u2550\u2550\u2550", style: "text-terminal-cyan" },
      { text: "", style: "" },
      { text: "  \u2713 bootstrap     [done]", style: "text-terminal-green" },
      { text: "  \u2713 database      [done]", style: "text-terminal-green" },
      { text: "  \u2713 nginx         [done]", style: "text-terminal-green" },
      { text: "  \u2713 php           [done]", style: "text-terminal-green" },
      { text: "  \u2713 security      [done]", style: "text-terminal-green" },
      { text: "", style: "" },
      { text: "\u2713 Server 'myserver' provisioned successfully!", style: "text-terminal-green" },
    ],
  },
  {
    title: "press site create",
    lines: [
      { text: "$ press site create", style: "text-white" },
      { text: "", style: "" },
      { text: "? Select server: myserver (203.0.113.10)", style: "text-terminal-green" },
      { text: "? Domain: blog.example.com", style: "text-terminal-green" },
      { text: "? Admin email: admin@example.com", style: "text-terminal-green" },
      { text: "", style: "" },
      { text: "Creating WordPress site...", style: "text-zinc-400" },
      { text: "", style: "" },
      { text: "  \u2713 System user created", style: "text-terminal-green" },
      { text: "  \u2713 PHP-FPM pool configured", style: "text-terminal-green" },
      { text: "  \u2713 Nginx vhost active", style: "text-terminal-green" },
      { text: "  \u2713 Database ready", style: "text-terminal-green" },
      { text: "  \u2713 WordPress installed", style: "text-terminal-green" },
      { text: "", style: "" },
      { text: "\u2713 Site live at https://blog.example.com", style: "text-terminal-green" },
    ],
  },
  {
    title: "press domain ssl",
    lines: [
      { text: "$ press domain ssl", style: "text-white" },
      { text: "", style: "" },
      { text: "? Select server: myserver (203.0.113.10)", style: "text-terminal-green" },
      { text: "? Select site: blog.example.com", style: "text-terminal-green" },
      { text: "? Select domain: blog.example.com", style: "text-terminal-green" },
      { text: "", style: "" },
      { text: "Issuing SSL certificate...", style: "text-zinc-400" },
      { text: "", style: "" },
      { text: "  \u2713 HTTP challenge passed", style: "text-terminal-green" },
      { text: "  \u2713 Certificate issued (Let's Encrypt)", style: "text-terminal-green" },
      { text: "  \u2713 Nginx reloaded with HTTPS", style: "text-terminal-green" },
      { text: "  \u2713 Auto-renewal configured", style: "text-terminal-green" },
      { text: "", style: "" },
      { text: "\u2713 SSL active — expires 2026-06-24", style: "text-terminal-green" },
    ],
  },
];

export default function CommandShowcase() {
  return (
    <section className="px-6 py-24 max-w-5xl mx-auto">
      <h2 className="text-2xl sm:text-3xl font-bold text-white text-center mb-4">
        Everything happens in your terminal
      </h2>
      <p className="text-zinc-500 text-center mb-14 max-w-xl mx-auto">
        Provision servers, deploy sites, manage domains and SSL — all through
        simple, interactive commands.
      </p>

      <div className="space-y-8">
        {commands.map((cmd) => (
          <TerminalWindow key={cmd.title} title={cmd.title}>
            {cmd.lines.map((line, i) => (
              <div key={i} className={`${line.style} whitespace-pre`}>
                {line.text || "\u00A0"}
              </div>
            ))}
          </TerminalWindow>
        ))}
      </div>
    </section>
  );
}
