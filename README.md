# Alertmanager Discord

# Environment configuration variables:

DISCORD_WEBHOOK - webhook, where to post alerts. For more details see: https://support.discordapp.com/hc/en-us/articles/228383668-Intro-to-Webhooks
DISCORD_NAME - bot name at Discord. Default AlertManager.


# Example Prometheus Alertmanager config:

```
global:
  # The smarthost and SMTP sender used for mail notifications.
  smtp_smarthost: 'localhost:25'
  smtp_from: 'alertmanager@example.org'
  smtp_require_tls: false

# The root route on which each incoming alert enters.
route:
  group_by: ['alertname', 'cluster', 'service']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  receiver: discord_webhook

receivers:
- name: 'discord_webhook'
  webhook_configs:
  - url: 'http://localhost:9095'
```
