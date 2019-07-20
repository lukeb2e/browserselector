# Browser Selector

This tool can be configured as the default browser and switch to the correct browser depending on the FQDN. This is supposed to help in environment where you need to open certain tools in a specific browser.

## Config

You need to configure a `config.yml` file directly next to the browserselector.

```yaml
domain:
  - browser: "/bin/firefox" # path to browser
    regex: ".*"             # regex to match domain
    priority: 999           # priority - rules will be evaluated from lowest to highest
  - browser: "/bin/chrome"
    regex: ".*web\\.de"
    priority: 10
debug: true # optional
```
