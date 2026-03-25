import type { Metadata } from "next";
import CopyButton from "@/components/ui/CopyButton";
import Footer from "@/components/Footer";

export const metadata: Metadata = {
  title: "Documentation — pressctl",
  description:
    "Complete guide to installing, configuring, and using pressctl to manage WordPress hosting infrastructure from the command line.",
};

function CodeBlock({ code, copyable = true }: { code: string; copyable?: boolean }) {
  return (
    <div className="relative group mb-4">
      <pre className="bg-zinc-900 border border-zinc-800 rounded-lg px-4 py-3 overflow-x-auto font-mono text-sm">
        <code className="text-zinc-300 whitespace-pre-wrap break-all">{code}</code>
      </pre>
      {copyable && (
        <div className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
          <CopyButton text={code} />
        </div>
      )}
    </div>
  );
}

function SidebarLink({ href, children }: { href: string; children: React.ReactNode }) {
  return (
    <a
      href={href}
      className="block text-sm text-zinc-500 hover:text-white transition-colors py-1"
    >
      {children}
    </a>
  );
}

export default function DocsPage() {
  return (
    <>
      <div className="pt-14">
        <div className="max-w-6xl mx-auto px-6 py-16 flex gap-12">
          {/* Sidebar */}
          <aside className="hidden lg:block w-56 shrink-0 sticky top-20 self-start">
            <nav className="space-y-1">
              <p className="text-xs font-semibold text-zinc-600 uppercase tracking-wider mb-3">
                Getting Started
              </p>
              <SidebarLink href="#requirements">Requirements</SidebarLink>
              <SidebarLink href="#install">Installation</SidebarLink>
              <SidebarLink href="#initialize">Initialize</SidebarLink>

              <p className="text-xs font-semibold text-zinc-600 uppercase tracking-wider mt-6 mb-3">
                Server
              </p>
              <SidebarLink href="#provision">Provision a Server</SidebarLink>
              <SidebarLink href="#server-list">List Servers</SidebarLink>
              <SidebarLink href="#server-health">Health Check</SidebarLink>

              <p className="text-xs font-semibold text-zinc-600 uppercase tracking-wider mt-6 mb-3">
                Sites
              </p>
              <SidebarLink href="#site-create">Create a Site</SidebarLink>
              <SidebarLink href="#site-list">List Sites</SidebarLink>
              <SidebarLink href="#site-delete">Delete a Site</SidebarLink>

              <p className="text-xs font-semibold text-zinc-600 uppercase tracking-wider mt-6 mb-3">
                Domains &amp; SSL
              </p>
              <SidebarLink href="#domain-add">Add a Domain</SidebarLink>
              <SidebarLink href="#domain-ssl">Enable SSL</SidebarLink>

              <p className="text-xs font-semibold text-zinc-600 uppercase tracking-wider mt-6 mb-3">
                Reference
              </p>
              <SidebarLink href="#config">Configuration</SidebarLink>
              <SidebarLink href="#stack">Technology Stack</SidebarLink>
              <SidebarLink href="#commands">All Commands</SidebarLink>
            </nav>
          </aside>

          {/* Main content */}
          <main className="flex-1 min-w-0 docs-prose">
            <h1 className="text-3xl sm:text-4xl font-bold text-white mb-2">
              Documentation
            </h1>
            <p className="text-zinc-500 text-lg mb-10">
              Everything you need to go from a blank VPS to a live WordPress site.
            </p>

            {/* ── Requirements ── */}
            <h2 id="requirements">Requirements</h2>
            <p>Before installing pressctl, make sure you have:</p>
            <ol>
              <li><strong>Linux or macOS</strong> on your local machine</li>
              <li><strong>Ansible 2.14+</strong> installed (<code>pip install ansible</code> or <code>brew install ansible</code>)</li>
              <li><strong>curl</strong> installed (ships with most systems)</li>
              <li>
                A <strong>fresh Ubuntu 24.04 server</strong> with root SSH access —
                DigitalOcean, Hetzner, AWS, Linode, or any VPS provider
              </li>
            </ol>
            <p>Check if curl is available:</p>
            <CodeBlock code="curl --version" />
            <p>
              If not installed — on Ubuntu/Debian: <code>sudo apt update &amp;&amp; sudo apt install curl -y</code>.
              On macOS: <code>brew install curl</code>.
            </p>

            {/* ── Installation ── */}
            <h2 id="install">Installation</h2>
            <p>Install pressctl with a single command:</p>
            <CodeBlock
              code="curl -fsSL https://raw.githubusercontent.com/pressctl/cli/main/install.sh | bash"
            />
            <p><strong>What this does:</strong></p>
            <ol>
              <li>Downloads the <code>press</code> binary from the pressctl GitHub repository</li>
              <li>Copies the Ansible playbooks to <code>~/.pressctl/ansible/</code></li>
              <li>Makes the <code>press</code> command available in your terminal</li>
            </ol>
            <p>
              No repository cloning, no manual build steps. One command and you&apos;re ready.
            </p>

            {/* ── Initialize ── */}
            <h2 id="initialize">Initialize</h2>
            <p>After installation, run the init command to create your configuration:</p>
            <CodeBlock code="press init" />
            <p><strong>What this does:</strong></p>
            <ul>
              <li>Creates <code>~/.pressctl/pressctl.yaml</code> — the config file that tracks all your servers, sites, and domains</li>
              <li>Prompts for your <strong>SSH public key path</strong> (used to access provisioned servers)</li>
              <li>Prompts for your <strong>email address</strong> (used for Let&apos;s Encrypt SSL certificates)</li>
            </ul>
            <p>After this step, pressctl is fully configured and ready to provision servers.</p>

            {/* ── Provision ── */}
            <h2 id="provision">Provision a Server</h2>
            <CodeBlock code="press server provision" />
            <p>
              This is the main command that sets up everything your server needs to host
              WordPress. You can also just run <code>press</code> and choose option <strong>1</strong> from the
              quick actions menu.
            </p>
            <p>pressctl will ask you a few questions:</p>
            <ol>
              <li><strong>Server name</strong> — a label for your reference (e.g. <code>production</code>)</li>
              <li><strong>IP address</strong> — your server&apos;s public IP</li>
              <li><strong>SSH user</strong> — defaults to <code>root</code> (press Enter to accept)</li>
              <li><strong>SSH port</strong> — defaults to <code>22</code></li>
              <li><strong>SSH private key</strong> — defaults to <code>~/.ssh/id_rsa</code></li>
            </ol>

            <p><strong>Pre-flight checks:</strong></p>
            <ul>
              <li><strong>SSH connectivity</strong> — verifies pressctl can reach your server</li>
              <li><strong>Port conflicts</strong> — checks that ports 80, 443, and 3306 are free. If Apache, Caddy, or MySQL are already installed, pressctl warns you before proceeding</li>
            </ul>

            <p><strong>What gets installed:</strong></p>
            <ul>
              <li>Nginx from the official repository</li>
              <li>PHP 8.3 with optimized FPM pools</li>
              <li>MariaDB with secure defaults</li>
              <li>Redis for object caching</li>
              <li>Let&apos;s Encrypt via Certbot</li>
              <li>UFW firewall (ports 22, 80, 443 only)</li>
              <li>Fail2ban for brute-force protection</li>
              <li>SSH hardening</li>
            </ul>
            <p>
              Provisioning takes about 5–10 minutes. When complete, you&apos;ll see a success
              message with the generated MySQL password (also saved in your config file).
            </p>

            <blockquote>
              <strong>Important:</strong> At this stage, pressctl only prepares the server.
              No WordPress site exists yet — site-specific features like isolated Linux
              users and per-site PHP-FPM pools happen when you create a site.
            </blockquote>

            {/* ── Server List ── */}
            <h2 id="server-list">List Servers</h2>
            <CodeBlock code="press server list" />
            <p>Shows all servers in your pressctl inventory with their status, IP, and site count.</p>
            <CodeBlock code="press server list --json" />
            <p>Use <code>--json</code> for machine-readable output, useful for scripts and automation.</p>

            {/* ── Health Check ── */}
            <h2 id="server-health">Server Health Check</h2>
            <CodeBlock code="press server health-check" />
            <p>
              Tests SSH connectivity to a server. Select from the list or pass the name
              directly: <code>press server health-check myserver</code>.
            </p>

            {/* ── Create Site ── */}
            <h2 id="site-create">Create a WordPress Site</h2>
            <CodeBlock code="press site create" />
            <p>
              Creates a new WordPress site on a provisioned server. pressctl will
              prompt you for:
            </p>
            <ol>
              <li><strong>Target server</strong> — choose which provisioned server to use</li>
              <li><strong>Primary domain</strong> — e.g. <code>yoursite.com</code></li>
              <li><strong>WordPress admin username</strong></li>
              <li><strong>WordPress admin email</strong></li>
              <li><strong>Password</strong> — let pressctl generate a secure one, or set your own</li>
            </ol>

            <p><strong>What happens for each site:</strong></p>
            <ul>
              <li>A <strong>dedicated Linux user</strong> is created for the site</li>
              <li>A <strong>separate PHP-FPM pool</strong> runs under that user</li>
              <li>An <strong>isolated database</strong> is created with its own credentials</li>
              <li>File permissions are locked down — sites can&apos;t read each other&apos;s files</li>
              <li>Nginx is configured with the site&apos;s domain</li>
              <li>WordPress is installed via WP-CLI</li>
              <li>A cron job is set up for <code>wp-cron</code></li>
            </ul>

            <blockquote>
              Make sure to save the admin credentials shown after site creation.
              They&apos;re stored in your local config file at <code>~/.pressctl/pressctl.yaml</code>.
            </blockquote>

            {/* ── Site List ── */}
            <h2 id="site-list">List Sites</h2>
            <CodeBlock code="press site list" />
            <p>Displays all sites across all servers with their domains and status.</p>

            {/* ── Delete Site ── */}
            <h2 id="site-delete">Delete a Site</h2>
            <CodeBlock code="press site delete" />
            <p>
              Removes a WordPress site and all its associated resources (database, Linux
              user, Nginx config, files, PHP-FPM pool). You&apos;ll be asked to confirm
              before anything is deleted.
            </p>

            {/* ── Domain Add ── */}
            <h2 id="domain-add">Add a Domain</h2>
            <CodeBlock code="press domain add" />
            <p>
              Attaches a domain to an existing WordPress site. pressctl will prompt for
              the server, the site, and the new domain name.
            </p>
            <blockquote>
              Make sure the domain&apos;s DNS A record points to your server&apos;s IP address
              before adding it. If DNS isn&apos;t configured, SSL certificate issuance will
              fail in the next step.
            </blockquote>

            {/* ── SSL ── */}
            <h2 id="domain-ssl">Enable SSL</h2>
            <CodeBlock code="press domain ssl" />
            <p>
              Issues a free SSL certificate for your domain using Let&apos;s Encrypt,
              making your site accessible over HTTPS. The certificate auto-renews via
              Certbot&apos;s systemd timer.
            </p>

            {/* ── Config ── */}
            <h2 id="config">Configuration</h2>
            <CodeBlock code="press config show" />
            <p>
              Displays your current pressctl configuration — servers, sites, domains,
              and global settings. Read-only, safe to run anytime.
            </p>
            <CodeBlock code="press config validate" />
            <p>Checks your configuration for errors and missing required values.</p>
            <CodeBlock code="press config edit" />
            <p>
              Opens the config file in your default editor. The config lives
              at <code>~/.pressctl/pressctl.yaml</code>.
            </p>

            {/* ── Stack ── */}
            <h2 id="stack">Technology Stack</h2>
            <p>Every provisioned server runs this exact stack:</p>
            <div className="overflow-x-auto mb-4">
              <table className="w-full text-sm border border-zinc-800 rounded-lg overflow-hidden">
                <thead>
                  <tr className="bg-zinc-900 text-left">
                    <th className="px-4 py-2 text-zinc-400 font-medium">Component</th>
                    <th className="px-4 py-2 text-zinc-400 font-medium">Details</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-zinc-800">
                  <tr><td className="px-4 py-2 text-zinc-300">OS</td><td className="px-4 py-2 text-zinc-400">Ubuntu 24.04 LTS</td></tr>
                  <tr><td className="px-4 py-2 text-zinc-300">Web Server</td><td className="px-4 py-2 text-zinc-400">Nginx (official repo)</td></tr>
                  <tr><td className="px-4 py-2 text-zinc-300">PHP</td><td className="px-4 py-2 text-zinc-400">8.3 (ondrej/php PPA)</td></tr>
                  <tr><td className="px-4 py-2 text-zinc-300">Database</td><td className="px-4 py-2 text-zinc-400">MariaDB</td></tr>
                  <tr><td className="px-4 py-2 text-zinc-300">Cache</td><td className="px-4 py-2 text-zinc-400">Redis</td></tr>
                  <tr><td className="px-4 py-2 text-zinc-300">SSL</td><td className="px-4 py-2 text-zinc-400">Let&apos;s Encrypt (Certbot)</td></tr>
                  <tr><td className="px-4 py-2 text-zinc-300">Firewall</td><td className="px-4 py-2 text-zinc-400">UFW (ports 22, 80, 443)</td></tr>
                  <tr><td className="px-4 py-2 text-zinc-300">Security</td><td className="px-4 py-2 text-zinc-400">Fail2ban + SSH hardening</td></tr>
                </tbody>
              </table>
            </div>

            {/* ── All Commands ── */}
            <h2 id="commands">All Commands</h2>
            <p>Quick reference for every <code>press</code> command:</p>
            <div className="overflow-x-auto mb-4">
              <table className="w-full text-sm border border-zinc-800 rounded-lg overflow-hidden">
                <thead>
                  <tr className="bg-zinc-900 text-left">
                    <th className="px-4 py-2 text-zinc-400 font-medium">Command</th>
                    <th className="px-4 py-2 text-zinc-400 font-medium">Description</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-zinc-800">
                  <tr><td className="px-4 py-2"><code>press</code></td><td className="px-4 py-2 text-zinc-400">Quick actions menu</td></tr>
                  <tr><td className="px-4 py-2"><code>press init</code></td><td className="px-4 py-2 text-zinc-400">Initialize configuration</td></tr>
                  <tr><td className="px-4 py-2"><code>press server provision</code></td><td className="px-4 py-2 text-zinc-400">Add and provision a new server</td></tr>
                  <tr><td className="px-4 py-2"><code>press server list</code></td><td className="px-4 py-2 text-zinc-400">List all servers</td></tr>
                  <tr><td className="px-4 py-2"><code>press server health-check</code></td><td className="px-4 py-2 text-zinc-400">Test server connectivity</td></tr>
                  <tr><td className="px-4 py-2"><code>press server update</code></td><td className="px-4 py-2 text-zinc-400">Update server configuration</td></tr>
                  <tr><td className="px-4 py-2"><code>press server remove</code></td><td className="px-4 py-2 text-zinc-400">Remove from inventory</td></tr>
                  <tr><td className="px-4 py-2"><code>press site create</code></td><td className="px-4 py-2 text-zinc-400">Create a WordPress site</td></tr>
                  <tr><td className="px-4 py-2"><code>press site list</code></td><td className="px-4 py-2 text-zinc-400">List all sites</td></tr>
                  <tr><td className="px-4 py-2"><code>press site delete</code></td><td className="px-4 py-2 text-zinc-400">Delete a site</td></tr>
                  <tr><td className="px-4 py-2"><code>press domain add</code></td><td className="px-4 py-2 text-zinc-400">Add domain to a site</td></tr>
                  <tr><td className="px-4 py-2"><code>press domain remove</code></td><td className="px-4 py-2 text-zinc-400">Remove a domain</td></tr>
                  <tr><td className="px-4 py-2"><code>press domain ssl</code></td><td className="px-4 py-2 text-zinc-400">Issue SSL certificate</td></tr>
                  <tr><td className="px-4 py-2"><code>press config show</code></td><td className="px-4 py-2 text-zinc-400">Show current configuration</td></tr>
                  <tr><td className="px-4 py-2"><code>press config validate</code></td><td className="px-4 py-2 text-zinc-400">Validate configuration</td></tr>
                  <tr><td className="px-4 py-2"><code>press config edit</code></td><td className="px-4 py-2 text-zinc-400">Edit config in your editor</td></tr>
                  <tr><td className="px-4 py-2"><code>press version</code></td><td className="px-4 py-2 text-zinc-400">Show version info</td></tr>
                </tbody>
              </table>
            </div>

            <p className="mt-10 text-zinc-600 text-sm">
              Need help? <a href="https://github.com/shariffff/wpsh/issues">Open an issue on GitHub</a>.
            </p>
          </main>
        </div>
      </div>
      <Footer />
    </>
  );
}
