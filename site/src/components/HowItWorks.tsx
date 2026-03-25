const steps = [
  {
    number: "1",
    title: "Install",
    description: "One curl command copies the press binary and Ansible playbooks to your machine.",
    code: "curl -fsSL https://...install.sh | bash\npress init",
  },
  {
    number: "2",
    title: "Add your server",
    description: "Point pressctl at any fresh Ubuntu 24.04 VPS — DigitalOcean, Hetzner, AWS, anything with root SSH.",
    code: "press server provision\n# Enter IP, SSH key path",
  },
  {
    number: "3",
    title: "Deploy sites",
    description: "Create WordPress sites, add domains, and issue SSL certificates — all from the same menu.",
    code: "press site create\npress domain add\npress domain ssl",
  },
];

export default function HowItWorks() {
  return (
    <section className="px-6 py-24 max-w-5xl mx-auto">
      <h2 className="text-2xl sm:text-3xl font-mono font-bold text-white text-center mb-4">
        Three steps. That&apos;s it.
      </h2>
      <p className="text-zinc-500 text-center mb-16 max-w-lg mx-auto">
        From a blank VPS to a live WordPress site in minutes.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-8 relative">
        {/* Connecting line (desktop only) */}
        <div className="hidden md:block absolute top-10 left-[16.5%] right-[16.5%] h-px border-t border-dashed border-zinc-700" />

        {steps.map((step) => (
          <div key={step.number} className="relative text-center">
            {/* Step number */}
            <div className="inline-flex items-center justify-center w-10 h-10 rounded-full border border-terminal-green/30 bg-terminal-green/10 text-terminal-green font-mono font-bold text-lg mb-5 relative z-10">
              {step.number}
            </div>

            <h3 className="text-lg font-mono font-semibold text-white mb-2">
              {step.title}
            </h3>
            <p className="text-sm text-zinc-400 mb-4 leading-relaxed">
              {step.description}
            </p>

            {/* Code snippet */}
            <div className="text-left px-4 py-3 rounded-lg bg-zinc-900 border border-zinc-800 font-mono text-xs">
              {step.code.split("\n").map((line, i) => (
                <div key={i} className={line.startsWith("#") ? "text-zinc-600" : "text-zinc-400"}>
                  {!line.startsWith("#") && (
                    <span className="text-terminal-green select-none">$ </span>
                  )}
                  {line.startsWith("#") ? (
                    <span className="text-zinc-600">{line}</span>
                  ) : (
                    line
                  )}
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}
