# Taste

## Infrastructure & Deployment
- Prefers supporting multiple deployment stacks side by side and letting the user choose their stack (PHP-FPM + Nginx, FrankenPHP Docker, FrankenPHP native). Confidence: 0.9
- Prefers focusing on the native FrankenPHP (deb package) setup first; Docker is deferred to later. Confidence: 0.7
- Wants per-site access logs and per-site resource usage summaries so it is easy to see which site consumes the most resources. Confidence: 0.8
- Prefers configuring PHP via PHP_INI_SCAN_DIR / conf.d-style .ini files rather than editing the main php.ini directly. Confidence: 0.6
- Manages server provisioning through Ansible roles (e.g., nginx role with a preset Cloudflare certificate). Confidence: 0.6
- Prefers pinning external app/tool versions in role defaults with sha256 checksums so downloads are reproducible. Confidence: 0.8
- Prefers writing tasks with `ansible.builtin` modules and deploying config via `*.j2` templates that carry an `{{ ansible_managed }}` header. Confidence: 0.7
- Prefers keeping secrets outside the web root with restrictive file modes (e.g. 0600), using `no_log: true` on tasks that handle them. Confidence: 0.8
- Prefers exposing extra web tools via nginx `^~` prefix locations (to win over regex handlers) and explicitly denying direct access to implementation files to prevent an auth bypass. Confidence: 0.7
- Prefers validating provisioning work with yamllint, ansible-lint, and `ansible-playbook --syntax-check` before finishing. Confidence: 0.6
- Prefers tool/internal names, URLs, and template filenames to accurately reflect the software actually deployed (e.g. "filemanager" rather than a "elfinder" codename when Tiny File Manager is the real software). Confidence: 0.6
