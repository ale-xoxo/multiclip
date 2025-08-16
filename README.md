# MultiClip - macOS Clipboard Manager

A lightweight macOS clipboard manager built with Go that keeps track of your last 5 copied items and provides easy access through the menu bar.

## Features

- 📋 Tracks last 5 clipboard items automatically
- 🔝 Menu bar icon for quick access
- 💾 Persistent storage (survives app restart)
- 🖱️ Click any item to copy it back to clipboard
- 🧹 Clear history option
- ⚡ Lightweight and fast

## Installation

1. Clone or download this repository
2. Run the build script:
   ```bash
   chmod +x build.sh
   ./build.sh
   ```
3. Run the application:
   ```bash
   ./multiclip
   ```

## Usage

1. Launch MultiClip - a clipboard icon will appear in your menu bar
2. Copy any text (Cmd+C) - it will automatically be saved
3. Click the menu bar icon to see your clipboard history
4. Click any item in the dropdown to copy it back to your clipboard
5. Use "Clear History" to remove all saved clips

## System Requirements

- macOS 10.12 or later
- Go 1.20+ (for building from source)

## Technical Details

- Built with Go using `systray` for menu bar integration
- Uses `golang.design/x/clipboard` for clipboard monitoring
- Stores clipboard history in `~/.multiclip.json`
- Monitors clipboard every 500ms for changes
- Menu updates every 2 seconds

## Auto-start Setup

To run MultiClip automatically at login:
1. Open System Preferences → Users & Groups
2. Select your user and go to "Login Items"
3. Click "+" and add the `multiclip` executable

## Building

```bash
go mod tidy
go build -o multiclip main.go
```

## License

MIT License - see LICENSE file for details