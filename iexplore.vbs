' VB Script Document ' http://stackoverflow.com/questions/31441974
' Opens blank page in a new MSIE tab and navigates to the URL if supplied. 
' Usage: ie [URL]
'        if saved as "IE.VBS" somewhere under PATH system environment variable.
'
option explicit
On Error Goto 0
Dim strMyUrl, intWExist, BrowserNavFlag, iRetNav2, objArgs, WshShell, IE, oIE, sTitle

Set objArgs = WScript.Arguments
If objArgs.Count > 0 Then
    strMyUrl = objArgs( 0)
Else
    strMyUrl = "about:blank"
End If

Set WshShell = WScript.CreateObject( "WScript.Shell")
Set IE  = Nothing
Set oIE = Nothing

intWExist = FindIE( strMyUrl, oIE) ' look for MSIE window pointer 

set IE = oIE
iRetNav2 = 0
Select Case intWExist
' Case 3
'   ''' 3 = MSIE window found, URL match, window title match
'   ''' (not implemented yet)
' Case 2
'   ''' 2 = MSIE window found, URL match
Case 1, 2, 3
    ''' 1 = MSIE window found, no URL match
    BrowserNavFlag = navOpenInNewTab ' 2048
    iRetNav2 = IE.Navigate2( strMyUrl, CLng( BrowserNavFlag), "_blank")
Case Else
    ''' 0 = MSIE window not found
    Set IE = CreateObject( "InternetExplorer.Application")
    BrowserNavFlag = 1
    iRetNav2 = IE.Navigate( strMyUrl)
End Select

IE.Visible = True

While IE.Busy
    Wscript.Sleep 100
Wend
'While IE.Document.ReadyState <> "complete" '(obsolete?) Or IE.ReadyState <> 4
'    Wscript.Sleep 100
'Wend

sTitle = ""
intWExist = FindIE( strMyUrl, oIE) ' look for MSIE window title
' AppActivate method could fail (no error) if MSIE window runs minimized
'                               off topic for current question (31441974)
If Not sTitle = "" Then WshShell.AppActivate sTitle

Private Function FindIE( ByVal sUrl, ByRef oObj)
' parameters
' sUrl (input)  string
' oObj (output) object
' returns 
' 0 = any MSIE window not found - or found but not accessible   
' 1 = a MSIE window found
' 2 = 1 and address line match
' 3 = 2 and title match (not implemented yet)
    Dim ww, tpnm, tptitle, tpfulln, tpUrl, tpUrlUnencoded
    Dim errNo, errStr, intLoop, intLoopLimit
    Dim iFound : iFound = 0
    Dim shApp    : Set shApp = CreateObject( "Shell.Application")
    With shApp
        For Each ww In .windows
            tpfulln = ww.FullName
            If Instr( 1, Lcase( tpfulln), "iexplore.exe", vbTextCompare) <> 0 _ 
        and Instr( 1, UCase( tpfulln), "SCODEF:", vbTextCompare) = 0 _ 
        and Instr( 1, UCase( tpfulln), "CREDAT:", vbTextCompare) = 0 Then
                If iFound > 0 Then
                Else
                    Set oObj = ww
                End If
                tptitle = "x x x" : tpUrl = "" : tpUrlUnencoded = ""
                intLoopLimit = 10 ' to look for attributes max. intLoopLimit/10 seconds
                intLoop = 0
                While intLoop < intLoopLimit
                    intLoop = intLoop + 1
                    On Error Resume Next
                    tpnm = typename( ww.document)
                    errNo = Err.Number
                    If errNo <> 0 Then
                        'error if  page not response (yet)' 
                        errStr = "Error # " & CStr( errNo) _
                & " """ & Err.Description & """ 0x" & Hex( errNo)  
                        Wscript.Sleep 100
                    Else
                        iFound = 1
                        intLoopLimit = intLoop  ' end the loop and preserve loop counter
                        tptitle = ww.document.title
                        tpUrl = ww.document.URL
                        tpUrlUnencoded = ww.document.URLUnencoded
                        errStr = tpnm
            sTitle = tptitle
                    End If
                    On Error Goto 0
                Wend
                If Instr( 1, Lcase( tpnm), "htmldocument", vbTextCompare) <> 0 then
                    If Instr( 1, Lcase( tpUrl), Lcase( sUrl), vbTextCompare) <> 0 Then
                        Set oObj = ww
                        iFound = 2
                        ' looking for all matching MSIE URLs 
                        ' this may take considerable time amount
                        ' to speed up script running, uncomment next line "exit for"
                        exit for
                    Else
                    End If 
                End If
            Else
                ' a program reports the same shell.application property as "iexplore.exe"
                ' i.e. "explorer.exe", "HTML preview" in some editors etc.
            End If
        Next
    End With
    Set shApp = Nothing
    FindIE = iFound
End Function
' 
' http://msdn.microsoft.com/en-us/library/aa768360(v=vs.85).aspx
' BrowserNavConstants Enumerated Type
' Contains values used by the IWebBrowser2::Navigate
'                         and IWebBrowser2::Navigate2 methods.
' typedef enum BrowserNavConstants {
Const navOpenInNewWindow      = &h01, _
      navNoHistory            = &h02, _
      navNoReadFromCache      = &h04, _
      navNoWriteToCache       = &h08, _
      navAllowAutosearch      = &h10, _
      navBrowserBar           = &h20, _
      navHyperlink            = &h40, _
      navEnforceRestricted    = &h80, _
      navNewWindowsManaged    = &h0100, _
      navUntrustedForDownload = &h0200, _
      navTrustedForActiveX    = &h0400, _
      navOpenInNewTab         = &h0800, _
      navOpenInBackgroundTab  = &h1000, _
      navKeepWordWheelText    = &h2000, _
      navVirtualTab           = &h4000, _
      navBlockRedirectsXDomain= &h8000, _
      navOpenNewForegroundTab = &h010000
' } BrowserNavConstants;
' 