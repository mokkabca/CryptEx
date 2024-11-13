# CryptEx

CryptEx is a Go-based ransomware tool for encrypting files and exfiltrating data through the Telegram Bot API. This project demonstrates the potential risks of ransomware attacks and serves as an educational tool for cybersecurity professionals.

---

## Installation

To get started with CryptEx, follow these steps:

1. Clone the repository:
    ```bash
    git clone https://github.com/mokkabca/CryptEx.git
    ```

2. Initialize Go modules:
    ```bash
    go mod init cryptex
    ```

3. Install the required Telegram Bot API package:
    ```bash
    go get github.com/go-telegram-bot-api/telegram-bot-api
    ```

---

## Usage

### Encrypting Files

To encrypt all files in the current directory, run the following command:

```bash
go run CryptEx.go
```

### Decrypting Files

To decrypt the files, use the following command:

```bash
go run DeCryptEx.go
```

### Building the Executable
To convert the project into a .exe file, use the following command:

```bash
go build -o CryptEx.exe -ldflags "-s -w -H windowsgui" CryptEx.go
go build -o DeCryptEx.exe -ldflags "-s -w -H windowsgui" DeCryptEx.go
```

Compressing the Executable
To reduce the size of the .exe file, you can use UPX compression:
```bash
upx --best --lzma CryptEx.exe
upx --best --lzma DeCryptEx.exe
```
