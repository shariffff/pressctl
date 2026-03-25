"use client";

import TerminalWindow from "@/components/ui/TerminalWindow";
import { useTypingAnimation } from "@/hooks/useTypingAnimation";
import { terminalScript } from "@/lib/terminalScript";

const colorMap: Record<string, string> = {
  green: "text-terminal-green",
  yellow: "text-terminal-yellow",
  cyan: "text-terminal-cyan",
  red: "text-terminal-red",
  muted: "text-zinc-500",
  white: "text-white",
};

export default function TerminalDemo() {
  const { lines, showCursor } = useTypingAnimation(terminalScript, 5000);

  return (
    <section id="demo" className="px-6 py-24 max-w-3xl mx-auto">
      <h2 className="text-2xl sm:text-3xl font-mono font-bold text-white text-center mb-4">
        See it in action
      </h2>
      <p className="text-zinc-500 text-center mb-10 max-w-lg mx-auto">
        This is what it looks like to provision a fresh server with pressctl.
      </p>

      <TerminalWindow>
        {lines.map((line, i) => (
          <div key={i} className={`${line.color ? colorMap[line.color] : "text-zinc-300"} whitespace-pre`}>
            {line.text}
            {i === lines.length - 1 && showCursor && (
              <span className="animate-blink text-terminal-green">
                &#9608;
              </span>
            )}
            {line.text === "" && "\u00A0"}
          </div>
        ))}
      </TerminalWindow>
    </section>
  );
}
