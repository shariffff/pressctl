"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import type { TerminalLine } from "@/lib/terminalScript";

interface DisplayLine {
  text: string;
  color?: string;
}

export function useTypingAnimation(script: TerminalLine[], loopPause = 4000) {
  const [lines, setLines] = useState<DisplayLine[]>([]);
  const [isRunning, setIsRunning] = useState(true);
  const [showCursor, setShowCursor] = useState(true);
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const abortRef = useRef(false);

  const sleep = useCallback(
    (ms: number) =>
      new Promise<void>((resolve) => {
        timeoutRef.current = setTimeout(resolve, ms);
      }),
    []
  );

  const runScript = useCallback(async () => {
    abortRef.current = false;
    setLines([]);
    setIsRunning(true);
    setShowCursor(true);

    for (const step of script) {
      if (abortRef.current) return;

      // Wait before this step
      if (step.delay > 0) {
        await sleep(step.delay);
      }
      if (abortRef.current) return;

      if (step.type === "pause") {
        await sleep(step.delay);
        continue;
      }

      if (step.type === "type") {
        // Type character by character
        for (let i = 0; i <= step.text.length; i++) {
          if (abortRef.current) return;

          const partial = step.text.slice(0, i);

          setLines((prev) => {
            if (step.append && prev.length > 0) {
              const updated = [...prev];
              const last = updated[updated.length - 1];
              updated[updated.length - 1] = {
                ...last,
                text: last.text + step.text.charAt(i - 1),
              };
              return i === 0 ? prev : updated;
            }
            if (i === 0) {
              return [...prev, { text: "", color: step.color }];
            }
            const updated = [...prev];
            updated[updated.length - 1] = { text: partial, color: step.color };
            return updated;
          });

          if (i < step.text.length) {
            await sleep(60);
          }
        }
      } else {
        // Instant
        if (step.append) {
          setLines((prev) => {
            if (prev.length === 0) return [{ text: step.text, color: step.color }];
            const updated = [...prev];
            const last = updated[updated.length - 1];
            updated[updated.length - 1] = {
              ...last,
              text: last.text + step.text,
            };
            return updated;
          });
        } else {
          setLines((prev) => [...prev, { text: step.text, color: step.color }]);
        }
      }
    }

    // Script done — pause with cursor off, then loop
    setShowCursor(false);
    setIsRunning(false);
    await sleep(loopPause);
    if (!abortRef.current) {
      runScript();
    }
  }, [script, sleep, loopPause]);

  useEffect(() => {
    runScript();
    return () => {
      abortRef.current = true;
      if (timeoutRef.current) clearTimeout(timeoutRef.current);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return { lines, isRunning, showCursor };
}
