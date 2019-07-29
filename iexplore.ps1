$website = $Args[0]
$BrowserNavFlag = 2048 #navOpenInNewTab

$ie = (New-Object -ComObject "Shell.Application").Windows() |
    Where-Object { $_.Name -eq "Internet Explorer" }
if (-not $ie) {
    $ie = (New-Object -com "InternetExplorer.Application" )
    $BrowserNavFlag = 64 #navHyperlink
}

echo $BrowserNavFlag
$ie.Navigate2("$website", $BrowserNavFlag)
$ie.Visible = $true