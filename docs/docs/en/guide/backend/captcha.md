# Captcha

Artalk comes with a built-in image captcha feature, allowing you to configure operation frequency limits that activate the captcha when exceeded.

Additionally, you can integrate various verification service providers to gain diversified verification functionality.

Currently, Artalk supports: Turnstile, reCAPTCHA, hCaptcha, and Geetest.

You can modify these settings in the [Dashboard](../frontend/sidebar.md#Settings) or via the [configuration file](./config.md#captcha) or [environment variables](../env.md#captcha).

## Configuration File

The full `captcha` configuration is as follows:

```yaml
# Captcha
captcha:
  enabled: true # Master switch
  always: false # Always require captcha
  captcha_type: image # Type of captcha
  action_limit: 3 # Number of actions required to activate captcha
  action_reset: 60 # Timeout to reset action counter (unit: s, set to -1 to never reset)
  # Turnstile
  # (https://www.cloudflare.com/products/turnstile/)
  turnstile:
    site_key: ''
    secret_key: ''
  # reCAPTCHA
  # (https://www.google.com/recaptcha/about/)
  recaptcha:
    site_key: ''
    secret_key: ''
  # hCaptcha (https://www.hcaptcha.com/)
  hcaptcha:
    site_key: ''
    secret_key: ''
  # Geetest (https://www.geetest.com)
  geetest:
    captcha_id: ''
    captcha_key: ''
```

- **always**: When set to `true`, captcha is always required.
- **captcha_type**: Type of captcha, options include: `image`, `turnstile`, `recaptcha`, `hcaptcha`, `geetest`.
- **action_limit**: The number of actions required to activate the captcha.
- **action_reset**: When the time exceeds this value, the action counter resets. Unit is seconds; set to `-1` to never reset.

Note: When `always` is enabled, `action_limit` and `action_reset` configurations are ignored.

## Configuration Examples

### Example 1

Within a 60-second time frame, if the number of actions exceeds 3, captcha will be required:

```yaml
captcha:
  action_limit: 3
  action_reset: 60
```

The counter will automatically reset after 60 seconds, allowing for 3 more actions without requiring a captcha.

### Example 2

Regardless of the time frame, if the number of actions from an IP address exceeds 5, captcha will be required:

```yaml
captcha:
  action_limit: 5
  action_reset: -1
```

### Example 3

Always require captcha, regardless of the number of actions:

```yaml
captcha:
  always: true
```

## Definition of Actions

Each "comment, vote, image upload, password verification" by an IP address counts as an "action."

## Turnstile

[Turnstile](https://www.cloudflare.com/zh-cn/products/turnstile/) is a verification service from Cloudflare. You can obtain the `site_key` and `secret_key` from the CF dashboard, then fill in these keys in the Artalk settings and change `captcha_type` to `turnstile`.

Illustrations:

<img src="/images/captcha/cf-turnstile-1.png" width="400px">

<img src="/images/captcha/cf-turnstile-2.png" width="400px">

Corresponding configuration file:

```yaml
captcha:
  # Omit other configurations...
  captcha_type: turnstile
  turnstile:
    site_key: ''
    secret_key: ''
```

## reCAPTCHA

[reCAPTCHA](https://developers.google.com/recaptcha) is a verification service from Google, and Artalk supports reCAPTCHA v3. You can obtain the `site_key` and `secret_key` from the Google developer dashboard, then fill in these keys in the Artalk settings and change `captcha_type` to `recaptcha`.

Corresponding configuration file:

```yaml
captcha:
  # Omit other configurations...
  captcha_type: recaptcha
  recaptcha:
    site_key: ''
    secret_key: ''
```

Note: Accessing Google APIs from within China may require configuring a system proxy.

Google provides test keys for reCAPTCHA: [see here](https://developers.google.com/recaptcha/docs/faq?hl=en#id-like-to-run-automated-tests-with-recaptcha.-what-should-i-do).

## hCaptcha

[hCaptcha](https://www.hcaptcha.com/) is a verification service. You can obtain the `site_key` and `secret_key` from its official website, then fill in these keys in the Artalk settings and change `captcha_type` to `hcaptcha`.

Corresponding configuration file:

```yaml
captcha:
  # Omit other configurations...
  captcha_type: hcaptcha
  hcaptcha:
    site_key: ''
    secret_key: ''
```

hCaptcha provides test keys: [see here](https://docs.hcaptcha.com/#integration-testing-test-keys).

## Geetest

Artalk supports integrating [Geetest](https://www.geetest.com/adaptive-captcha) for advanced behavior verification.

First, register an account on the Geetest website to obtain the `captcha_id` and `captcha_key`, then modify the configuration in the Artalk settings and change `captcha_type` to `geetest`.

Corresponding configuration file:

```yaml
captcha:
  # Omit other configurations...
  captcha_type: geetest
  geetest:
    captcha_id: ''
    captcha_key: ''
```
