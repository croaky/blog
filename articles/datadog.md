# Datadog

Staging:

```
staging config:set DISABLE_DATADOG_AGENT=true
```

Production:

```
production config:set \
  DD_AGENT_MAJOR_VERSION=7 \
  DD_API_KEY=<REPLACE> \
  DD_DYNO_HOST=true \
  DD_ENV=production \
  DD_LOG_LEVEL=ERROR\
  DD_SITE=datadoghq.com
```
