interface TerminalWindowProps {
  title?: string;
  children: React.ReactNode;
}

export default function TerminalWindow({
  title = "press — bash",
  children,
}: TerminalWindowProps) {
  return (
    <div className="rounded-xl border border-zinc-800 bg-zinc-900 overflow-hidden shadow-2xl shadow-black/50">
      {/* Title bar */}
      <div className="flex items-center gap-2 bg-zinc-800/80 px-4 py-3">
        <div className="flex gap-1.5">
          <div className="w-3 h-3 rounded-full bg-red-500/80" />
          <div className="w-3 h-3 rounded-full bg-yellow-500/80" />
          <div className="w-3 h-3 rounded-full bg-green-500/80" />
        </div>
        <span className="ml-2 text-xs text-zinc-500 font-mono">{title}</span>
      </div>

      {/* Terminal body */}
      <div className="p-5 font-mono text-sm leading-relaxed min-h-[300px] max-h-[520px] overflow-y-auto">
        {children}
      </div>
    </div>
  );
}
