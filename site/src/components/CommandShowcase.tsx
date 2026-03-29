import TerminalWindow from "@/components/ui/TerminalWindow";

const lines = [
  { text: "➜  ~ press", style: "text-white" },
  { text: "", style: "" },
  { text: "  pressctl  github.com/shariffff/pressctl", style: "text-zinc-500" },
  { text: "", style: "" },
  { text: "  ▶ 1  Provision a new server", style: "text-terminal-cyan" },
  { text: "    2  Create a WordPress site", style: "text-zinc-300" },
  { text: "    3  Delete a site", style: "text-zinc-300" },
  { text: "    4  Add a domain to a site", style: "text-zinc-300" },
  { text: "    5  Issue / renew SSL for a domain", style: "text-zinc-300" },
  { text: "    6  Check server health", style: "text-zinc-300" },
  { text: "    7  List servers", style: "text-zinc-300" },
  { text: "    8  List sites", style: "text-zinc-300" },
  { text: "", style: "" },
  {
    text: "  ↑↓ navigate  enter or number  type to filter  M more  Q quit",
    style: "text-zinc-600",
  },
];

export default function CommandShowcase() {
  return (
    <section className="px-6 py-24 max-w-3xl mx-auto">
      <h2 className="text-2xl sm:text-3xl font-bold text-white text-center mb-4">
        Everything happens in your terminal
      </h2>
      <p className="text-zinc-500 text-center mb-14 max-w-xl mx-auto">
        Provision servers, deploy sites, manage domains and SSL — all through
        one interactive menu.
      </p>

      <TerminalWindow title="press — bash">
        {lines.map((line, i) => (
          <div key={i} className={`${line.style} whitespace-pre`}>
            {line.text || "\u00A0"}
          </div>
        ))}
      </TerminalWindow>
    </section>
  );
}
