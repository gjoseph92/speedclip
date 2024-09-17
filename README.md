# speedclip

`speedclip` is a small utility for cropping [Speedscope profiles](https://github.com/jlfwong/speedscope/wiki/Importing-from-custom-sources#speedscopes-file-format).

Installation:
```
go install github.com/gjoseph92/speedclip@latest
```

Example usage:
```
speedclip -s 1m35s -e 2m3s profile.json > profile-clipped.json
speedscope <(speedclip -s 2.32s -e 9.9s profile.json)
```
