CGO_ENABLED=1 go build -ldflags="-s -w" -o r11-bin
scp -i ~/key-store/id_ed25519 r11-bin pojos@10.19.142.254:

