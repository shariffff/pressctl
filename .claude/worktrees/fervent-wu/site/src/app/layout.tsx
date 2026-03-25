import type { Metadata } from "next";
import localFont from "next/font/local";
import { Geist_Mono } from "next/font/google";
import Link from "next/link";
import GitHubIcon from "@/components/ui/GitHubIcon";
import "./globals.css";

const monaSans = localFont({
  src: "../../public/fonts/MonaSans-Variable.woff2",
  variable: "--font-sans",
  display: "swap",
});

const geistMono = Geist_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
});

export const metadata: Metadata = {
  title: "pressctl — WordPress infrastructure from your terminal",
  description:
    "Provision servers, deploy WordPress sites, and manage SSL — all from an interactive CLI powered by Ansible.",
  openGraph: {
    title: "pressctl — WordPress infrastructure from your terminal",
    description:
      "Provision servers, deploy WordPress sites, and manage SSL — all from an interactive CLI powered by Ansible.",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "pressctl — WordPress infrastructure from your terminal",
    description:
      "Provision servers, deploy WordPress sites, and manage SSL — all from an interactive CLI powered by Ansible.",
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className={`${monaSans.variable} ${geistMono.variable}`}>
      <body className="font-sans">
        {/* Top nav */}
        <nav className="fixed top-0 left-0 right-0 z-50 border-b border-zinc-800/50 bg-zinc-950/80 backdrop-blur-md">
          <div className="max-w-5xl mx-auto flex items-center justify-between px-6 h-14">
            <Link href="/" className="font-mono font-bold text-white text-sm hover:opacity-80 transition-opacity">
              press<span className="text-terminal-green">ctl</span>
            </Link>
            <div className="flex items-center gap-6">
              <Link
                href="/docs"
                className="text-sm text-zinc-400 hover:text-white transition-colors"
              >
                Docs
              </Link>
              <a
                href="https://github.com/shariffff/wpsh"
                target="_blank"
                rel="noopener noreferrer"
                className="text-zinc-400 hover:text-white transition-colors"
              >
                <GitHubIcon className="w-5 h-5" />
              </a>
            </div>
          </div>
        </nav>

        {children}
      </body>
    </html>
  );
}
