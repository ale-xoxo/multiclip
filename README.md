# MultiClip - macOS Menu Bar Clipboard Manager

A simple, minimalistic macOS menu bar application that tracks your last 5 clipboard items.

## Features

- 📋 **Menu Bar Integration**: Clean, Apple-like icon in the menu bar
- 🔄 **Real-time Monitoring**: Automatically tracks clipboard changes
- 📝 **Last 5 Items**: Shows your most recent clipboard entries
- 🖱️ **Click to Copy**: Click any item to copy it back to clipboard
- 🎨 **Minimalist Design**: Follows Apple's design guidelines

## Requirements

- macOS 10.12+ (Sierra or later)
- No additional dependencies

## Installation

1. Download the `multiclip` executable
2. Place it in your Applications folder or preferred location
3. Run the application
4. The clipboard icon will appear in your menu bar

## Usage

1. **Start the app**: Double-click `multiclip` or run from terminal
2. **View clipboard history**: Click the 📋 icon in the menu bar
3. **Copy previous items**: Click any item in the dropdown to copy it back
4. **Quit**: Select "Quit MultiClip" from the dropdown menu

## Technical Details

- **Language**: Go
- **Menu Bar Library**: fyne.io/systray
- **Clipboard Library**: golang.design/x/clipboard
- **Icon**: Template icon that adapts to light/dark mode

## Building from Source

```bash
# Clone and build
git clone <repo-url>
cd multiclip
go mod tidy
CGO_ENABLED=1 go build -ldflags="-s -w" -o build/multiclip src/main.go
```

## Architecture

- **ClipboardManager**: Thread-safe clipboard history management
- **Menu Integration**: Dynamic menu updates every 2 seconds
- **Template Icon**: Automatically adapts to macOS appearance settings

Built with ❤️ for macOS productivity.