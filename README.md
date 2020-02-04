# Browser Selector

This tool can be configured as the default browser and switch to the correct browser depending on the FQDN. This is supposed to help in environments where you need to open certain tools in a specific browser.

## Build

To build this tool you require the tool goversioninfo.

```
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go generate
go build
```

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

## Installation

The attached Registry file can be used to configure the tool as the default browser. The path to the binary has to be updated depending on your installation.

```
; PATH NEEDS TO BE EDITED HERE
[HKEY_CURRENT_USER\Software\Classes\BrowserselectorURL\shell\open\command]
@="\"C:\\tools\\browserselector\\browserselector.exe\" -- \"%1\""
```

For Windows 10 the configuration changed. New .reg file is still work in progress. If you want to figure it out yourself the way Firefox registers itself seems to be documented here:

`https://github.com/mozilla/gecko-dev/blob/master/browser/installer/windows/nsis/shared.nsh#L1456`

