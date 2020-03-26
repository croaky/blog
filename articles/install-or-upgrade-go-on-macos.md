# Install or Upgrade Go on macOS

Go to <https://golang.org/dl/>.
Copy the latest stable version number into a temporary shell variable:

```
gover="1.14"
```

Run:

```
if ! go version | grep -Fq "$gover"; then
  sudo rm -rf /usr/local/go
  curl "https://dl.google.com/go/go$gover.darwin-amd64.tar.gz" | sudo tar xz -C /usr/local
fi
```

With a 200 Mbps connection,
it completes in under 10 seconds.
