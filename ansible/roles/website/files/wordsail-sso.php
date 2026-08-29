<?php

/**
 * Plugin Name: WordSail WP Admin SSO
 * Description: Signs panel users into wp-admin via single-use SSO tokens
 * Author:      WordSail
 * Version:     1.2.0
 */

defined('ABSPATH') || exit;

define('WORDSAIL_WP_SSO_SERVER', 'app.wordsail.com');
define('WORDSAIL_WP_SSO_ENDPOINT', '/api/hosting/v1/wp-sso/');

add_action('login_init', 'wordsail_wp_sso_login_init');
add_filter('login_message', 'wordsail_wp_sso_login_message');

/**
 * Intercept wp-login.php requests carrying an SSO token.
 */
function wordsail_wp_sso_login_init(): void
{
    if (empty($_GET['token'])) { // phpcs:ignore WordPress.Security.NonceVerification.Recommended
        return;
    }

    // Already authenticated — don't burn the token, just go to wp-admin.
    if (is_user_logged_in()) {
        nocache_headers();
        wp_safe_redirect(admin_url());
        exit;
    }

    $token = sanitize_text_field(wp_unslash((string) $_GET['token'])); // phpcs:ignore WordPress.Security.NonceVerification.Recommended

    if (1 !== preg_match('/^[a-zA-Z0-9]{64}$/', $token) || ! wordsail_wp_sso_verify($token)) {
        wordsail_wp_sso_fail();
    }

    $user = wordsail_wp_sso_admin_user();

    if (! $user instanceof WP_User) {
        wordsail_wp_sso_fail();
    }

    wp_set_current_user($user->ID);
    wp_set_auth_cookie($user->ID);

    do_action('wp_login', $user->user_login, $user);

    nocache_headers();
    wp_safe_redirect(admin_url());
    exit;
}

/**
 * Verify the token against the control panel and confirm it was
 * minted for this site (the panel returns the site's domain, which
 * must match the host being logged into).
 */
function wordsail_wp_sso_verify(string $token): bool
{
    $response = wp_remote_get(
        'https://' . wordsail_wp_sso_server() . WORDSAIL_WP_SSO_ENDPOINT . rawurlencode($token),
        array(
            'timeout'     => 10,
            'redirection' => 0, // Never follow redirects with the token in the URL.
            'sslverify'   => true,
            'user-agent'  => 'WordSail Hosting API: ' . wp_parse_url(home_url(), PHP_URL_HOST),
        )
    );

    if (is_wp_error($response) || 200 !== wp_remote_retrieve_response_code($response)) {
        return false;
    }

    $data = json_decode(wp_remote_retrieve_body($response), true);

    if (! is_array($data) || empty($data['domain']) || ! is_string($data['domain'])) {
        return false;
    }

    $request_host = wordsail_wp_sso_request_host();

    return '' !== $request_host && hash_equals(strtolower($data['domain']), $request_host);
}

/**
 * The current request host, lowercased and stripped of any port,
 * so "example.com:443" still matches the panel's "example.com".
 */
function wordsail_wp_sso_request_host(): string
{
    $host = isset($_SERVER['HTTP_HOST']) ? strtolower(trim((string) wp_unslash($_SERVER['HTTP_HOST']))) : ''; // phpcs:ignore WordPress.Security.ValidatedSanitizedInput

    if ('' === $host) {
        return '';
    }

    return (string) preg_replace('/:\d+$/', '', $host);
}

/**
 * The panel host to verify against. Staging panels pass x_src_url,
 * restricted to *.wordsail.com so tokens can't be "verified" elsewhere.
 */
function wordsail_wp_sso_server(): string
{
    if (! empty($_GET['x_src_url'])) { // phpcs:ignore WordPress.Security.NonceVerification.Recommended
        $host = wp_parse_url(
            filter_var(wp_unslash((string) $_GET['x_src_url']), FILTER_SANITIZE_URL), // phpcs:ignore WordPress.Security.ValidatedSanitizedInput
            PHP_URL_HOST
        );

        if (is_string($host) && 1 === preg_match('/^[a-z0-9-]+(\.[a-z0-9-]+)*\.wordsail\.com$/i', $host)) {
            return strtolower($host);
        }
    }

    return WORDSAIL_WP_SSO_SERVER;
}

/**
 * The account to sign in as: the super admin on multisite,
 * otherwise the oldest administrator.
 */
function wordsail_wp_sso_admin_user(): ?WP_User
{
    if (is_multisite()) {
        foreach (get_super_admins() as $login) {
            $user = get_user_by('login', $login);
            if ($user instanceof WP_User) {
                return $user;
            }
        }
    }

    $admins = get_users(
        array(
            'role'        => 'administrator',
            'orderby'     => 'ID',
            'order'       => 'ASC',
            'number'      => 1,
            'count_total' => false,
        )
    );

    return $admins[0] ?? null;
}

/**
 * Bounce back to a clean login form with an error notice
 * instead of leaving the spent token in the address bar.
 */
function wordsail_wp_sso_fail(): void
{
    nocache_headers();
    wp_safe_redirect(add_query_arg('wordsail_sso_failed', '1', wp_login_url()));
    exit;
}

/**
 * Explain the failure on the login form.
 *
 * @param string $message Existing login message HTML.
 */
function wordsail_wp_sso_login_message(string $message): string
{
    if (! empty($_GET['wordsail_sso_failed'])) { // phpcs:ignore WordPress.Security.NonceVerification.Recommended
        $message .= '<div id="login_error">' . esc_html__('Your sign-on link has expired or was already used. Please open WP admin from the control panel again.', 'wordsail') . '</div>';
    }

    return $message;
}
