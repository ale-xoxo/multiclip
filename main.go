package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"golang.design/x/clipboard"
)

const maxClips = 5

type ClipboardManager struct {
	clips       []string
	lastContent string
	dataFile    string
}

func NewClipboardManager() *ClipboardManager {
	homeDir, _ := os.UserHomeDir()
	dataFile := filepath.Join(homeDir, ".multiclip.json")

	cm := &ClipboardManager{
		clips:    make([]string, 0, maxClips),
		dataFile: dataFile,
	}

	cm.loadClips()
	return cm
}

func (cm *ClipboardManager) loadClips() {
	data, err := os.ReadFile(cm.dataFile)
	if err != nil {
		return // File doesn't exist, start with empty clips
	}

	json.Unmarshal(data, &cm.clips)
}

func (cm *ClipboardManager) saveClips() {
	data, err := json.Marshal(cm.clips)
	if err != nil {
		return
	}

	os.WriteFile(cm.dataFile, data, 0644)
}

func (cm *ClipboardManager) addClip(content string) {
	content = strings.TrimSpace(content)
	if content == "" || content == cm.lastContent {
		return
	}

	// Remove if already exists
	for i, clip := range cm.clips {
		if clip == content {
			cm.clips = append(cm.clips[:i], cm.clips[i+1:]...)
			break
		}
	}

	// Add to front
	cm.clips = append([]string{content}, cm.clips...)

	// Keep only last 5
	if len(cm.clips) > maxClips {
		cm.clips = cm.clips[:maxClips]
	}

	cm.lastContent = content
	cm.saveClips()
}

func (cm *ClipboardManager) getClips() []string {
	return cm.clips
}

func (cm *ClipboardManager) setClip(content string) {
	clipboard.Write(clipboard.FmtText, []byte(content))
	cm.lastContent = content
}

var clipManager *ClipboardManager
var menuItems []*systray.MenuItem

func main() {
	err := clipboard.Init()
	if err != nil {
		log.Fatal("Failed to initialize clipboard:", err)
	}

	clipManager = NewClipboardManager()

	// Start clipboard monitoring in background
	go monitorClipboard()

	// Start system tray
	systray.Run(onReady, onExit)
}

func monitorClipboard() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		data := clipboard.Read(clipboard.FmtText)
		if len(data) > 0 {
			content := string(data)
			clipManager.addClip(content)
		}
	}
}

func onReady() {
	systray.SetIcon(getIcon())
	systray.SetTitle("MultiClip")
	systray.SetTooltip("Clipboard Manager - Last 5 copied items")

	// Create menu items
	updateMenu()

	// Refresh menu every 2 seconds
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			updateMenu()
		}
	}()
}

func updateMenu() {
	// Hide existing menu items
	for _, item := range menuItems {
		item.Hide()
	}
	menuItems = nil

	clips := clipManager.getClips()

	if len(clips) == 0 {
		mEmpty := systray.AddMenuItem("No clips yet", "No clipboard history")
		mEmpty.Disable()
		menuItems = append(menuItems, mEmpty)
	} else {
		for i, clip := range clips {
			// Truncate long clips for display
			displayText := clip
			if len(displayText) > 50 {
				displayText = displayText[:47] + "..."
			}

			// Replace newlines with spaces for menu display
			displayText = strings.ReplaceAll(displayText, "\n", " ")
			displayText = strings.ReplaceAll(displayText, "\r", " ")

			title := fmt.Sprintf("%d. %s", i+1, displayText)
			mClip := systray.AddMenuItem(title, "Click to copy to clipboard")
			menuItems = append(menuItems, mClip)

			// Handle click in a goroutine with captured clip content
			go func(content string, item *systray.MenuItem) {
				<-item.ClickedCh
				clipManager.setClip(content)
			}(clip, mClip)
		}
	}

	systray.AddSeparator()

	mClear := systray.AddMenuItem("Clear History", "Clear all clipboard history")
	menuItems = append(menuItems, mClear)
	go func() {
		<-mClear.ClickedCh
		clipManager.clips = []string{}
		clipManager.saveClips()
	}()

	mQuit := systray.AddMenuItem("Quit", "Quit MultiClip")
	menuItems = append(menuItems, mQuit)
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// Cleanup
}

func getIcon() []byte {
	// Simple 16x16 clipboard icon in ICO format (base64 encoded)
	// This is a minimal clipboard icon
	iconData := []byte{
		0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x10, 0x10, 0x00, 0x00, 0x01, 0x00, 0x20, 0x00, 0x68, 0x04,
		0x00, 0x00, 0x16, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x20, 0x00,
		0x00, 0x00, 0x01, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	// Add the actual icon pixels (simplified black clipboard shape)
	pixels := make([]byte, 16*16*4) // 16x16 RGBA

	// Draw a simple clipboard shape
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			idx := (y*16 + x) * 4

			// Create clipboard outline
			if (x >= 2 && x <= 13 && (y == 2 || y == 14)) ||
				(y >= 2 && y <= 14 && (x == 2 || x == 13)) ||
				(x >= 5 && x <= 10 && y >= 0 && y <= 3) {
				pixels[idx] = 0x00   // B
				pixels[idx+1] = 0x00 // G
				pixels[idx+2] = 0x00 // R
				pixels[idx+3] = 0xFF // A
			} else {
				pixels[idx] = 0xFF   // B
				pixels[idx+1] = 0xFF // G
				pixels[idx+2] = 0xFF // R
				pixels[idx+3] = 0x00 // A (transparent)
			}
		}
	}

	return append(iconData, pixels...)
}
