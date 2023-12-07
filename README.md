# process-logger
Tracker for process CPU/Memory etc.

How to use
Build for you platform
```
    go build process-tracker.go
    GOOS=windows GOARCH=386 go build process-tracker.go
```

Execute with process you want to track
```
    ./process-tracker openvpn
    ./process-tracker "Outlook.exe"
```
