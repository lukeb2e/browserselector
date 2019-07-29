# Browser Selector

This tool can be configured as the default browser and switch to the correct browser depending on the FQDN. This is supposed to help in environments where you need to open certain tools in a specific browser.

## Config

You need to configure a `config.yml` file directly next to the browserselector.

```yaml
browser: # define possible browsers
  iexplore: # call iexplorer with helper script to open new tab instead of new window
    exe: "C:/Windows/System32/WindowsPowerShell/v1.0/powershell.exe"
    script: "C:/tools/browserselector/iexplore.ps1"
  firefox:
    exe: "C:/Program Files/Mozilla Firefox/firefox.exe"

domain:
   - browser: "firefox" # name of browser
     regex: ".*"        # regex to match domain
     priority: 999      # priority - rules will be evaluated from lowest to highest
   - browser: "iexplore"
     regex: ".*\\.corpintra\\.net"
     priority: 10

debug: true # optional
```

