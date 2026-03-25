export type LineType = "type" | "instant" | "pause";

export interface TerminalLine {
  text: string;
  type: LineType;
  delay: number; // ms to wait before showing this line
  color?: "green" | "yellow" | "cyan" | "red" | "muted" | "white";
  append?: boolean; // append to previous line instead of new line
}

export const terminalScript: TerminalLine[] = [
  // User types the command
  { text: "$ press", type: "type", delay: 0, color: "white" },

  // Menu appears
  { text: "", type: "instant", delay: 400 },
  { text: "  pressctl — Common Actions", type: "instant", delay: 100, color: "cyan" },
  { text: "", type: "instant", delay: 100 },
  { text: "  1.  Provision a new server", type: "instant", delay: 80, color: "yellow" },
  { text: "  2.  Create a WordPress site", type: "instant", delay: 80, color: "yellow" },
  { text: "  3.  Delete a site", type: "instant", delay: 80, color: "yellow" },
  { text: "  4.  Add a domain to a site", type: "instant", delay: 80, color: "yellow" },
  { text: "  5.  Issue / renew SSL for a domain", type: "instant", delay: 80, color: "yellow" },
  { text: "  6.  Check server health", type: "instant", delay: 80, color: "yellow" },
  { text: "  7.  List servers", type: "instant", delay: 80, color: "yellow" },
  { text: "  8.  List sites", type: "instant", delay: 80, color: "yellow" },
  { text: "", type: "instant", delay: 100 },
  { text: "? Enter number: ", type: "instant", delay: 200, color: "green" },

  // User picks option 1
  { text: "1", type: "type", delay: 900, append: true },

  // Provisioning output
  { text: "", type: "instant", delay: 300 },
  { text: "Checking SSH connectivity...", type: "instant", delay: 100 },
  { text: "\u2713 SSH connectivity check passed", type: "instant", delay: 600, color: "green" },
  { text: "", type: "instant", delay: 200 },
  { text: "Checking for port conflicts...", type: "instant", delay: 100 },
  { text: "\u2713 No port conflicts detected", type: "instant", delay: 500, color: "green" },
  { text: "", type: "instant", delay: 300 },
  {
    text: "\u2550\u2550\u2550 Starting provisioning: myserver \u2550\u2550\u2550",
    type: "instant",
    delay: 200,
    color: "cyan",
  },
  { text: "", type: "instant", delay: 100 },
  { text: "  \u2713 bootstrap     [done]", type: "instant", delay: 800, color: "green" },
  { text: "  \u2713 database      [done]", type: "instant", delay: 600, color: "green" },
  { text: "  \u2713 nginx         [done]", type: "instant", delay: 500, color: "green" },
  { text: "  \u2713 php           [done]", type: "instant", delay: 700, color: "green" },
  { text: "  \u2713 security      [done]", type: "instant", delay: 400, color: "green" },
  { text: "", type: "instant", delay: 200 },
  {
    text: "\u2713 Server 'myserver' provisioned successfully!",
    type: "instant",
    delay: 300,
    color: "green",
  },
  { text: "", type: "instant", delay: 100 },
  { text: "Next steps:", type: "instant", delay: 200 },
  { text: "  Create a WordPress site: press site create", type: "instant", delay: 100, color: "muted" },
];
