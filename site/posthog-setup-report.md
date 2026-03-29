# PostHog post-wizard report

The wizard has completed a deep integration of PostHog analytics into the pressctl marketing site. The integration uses `posthog-js` initialized via Next.js 15.3+'s `instrumentation-client.ts` hook for automatic client-side initialization. Because the site uses `output: 'export'` (static site generation) in production, the PostHog host is set directly rather than through a `/ingest` proxy (which requires Next.js rewrites that aren't available in static exports). Five components were instrumented with custom event tracking to measure CTA performance, install intent, and GitHub engagement.

| Event | Description | File |
|---|---|---|
| `command_copied` | User copies a command snippet (install command or code block) using the copy button | `site/src/components/ui/CopyButton.tsx` |
| `github_link_clicked` | User clicks the View on GitHub CTA button in the Hero section | `site/src/components/Hero.tsx` |
| `docs_link_clicked` | User clicks the Read the docs link in the Hero section | `site/src/components/Hero.tsx` |
| `github_link_clicked` | User clicks the View on GitHub CTA button in the HowItWorks section | `site/src/components/HowItWorks.tsx` |
| `docs_link_clicked` | User clicks the Read the docs link in the HowItWorks section | `site/src/components/HowItWorks.tsx` |
| `github_link_clicked` | User clicks the GitHub link in the Footer | `site/src/components/Footer.tsx` |

## Next steps

We've built some insights and a dashboard for you to keep an eye on user behavior, based on the events we just instrumented:

- **Dashboard**: [Analytics basics](https://us.posthog.com/project/360767/dashboard/1409137)
- **Insight**: [GitHub link clicks over time](https://us.posthog.com/project/360767/insights/QvPtUlqh)
- **Insight**: [Command copies over time](https://us.posthog.com/project/360767/insights/BGcC8Agf)
- **Insight**: [CTA clicks by source](https://us.posthog.com/project/360767/insights/ggNGaGo2)
- **Insight**: [Visitor to GitHub CTA conversion funnel](https://us.posthog.com/project/360767/insights/HNea8YcG)
- **Insight**: [Visitor to install command copied funnel](https://us.posthog.com/project/360767/insights/y7ogk4FQ)

### Agent skill

We've left an agent skill folder in your project at `.claude/skills/integration-nextjs-app-router/`. You can use this context for further agent development when using Claude Code. This will help ensure the model provides the most up-to-date approaches for integrating PostHog.
