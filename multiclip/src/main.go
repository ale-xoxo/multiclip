package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"fyne.io/systray"
	"golang.design/x/clipboard"
)

type ClipboardManager struct {
	history []string
	mu      sync.RWMutex
	maxSize int
}

func NewClipboardManager() *ClipboardManager {
	return &ClipboardManager{
		history: make([]string, 0, 5),
		maxSize: 5,
	}
}

func (cm *ClipboardManager) AddItem(text string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Clean up text - remove excessive whitespace
	text = strings.TrimSpace(text)
	if text == "" || len(text) > 200 {
		return // Skip empty or very long items
	}

	// Check if item already exists
	for _, item := range cm.history {
		if item == text {
			return // Skip duplicates
		}
	}

	// Add to front of history
	cm.history = append([]string{text}, cm.history...)
	
	// Keep only last 5 items
	if len(cm.history) > cm.maxSize {
		cm.history = cm.history[:cm.maxSize]
	}
}

func (cm *ClipboardManager) GetHistory() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	result := make([]string, len(cm.history))
	copy(result, cm.history)
	return result
}

var (
	clipManager = NewClipboardManager()
	ctx, cancel = context.WithCancel(context.Background())
	menuItems   = make([]*systray.MenuItem, 0, 5)
	quitItem    *systray.MenuItem
)

func main() {
	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
		systray.Quit()
	}()

	log.Println("Starting MultiClip...")
	systray.Run(onReady, onExit)
}

func onReady() {
	log.Println("Initializing MultiClip...")
	
	// Set up system tray with Apple-like minimalist icon
	systray.SetTemplateIcon(getClipboardIcon(), getClipboardIcon())
	systray.SetTitle("📋")
	systray.SetTooltip("MultiClip - Clipboard Manager")

	// Initialize clipboard
	err := clipboard.Init()
	if err != nil {
		log.Printf("Failed to initialize clipboard: %v", err)
		// Don't return - still show the menu even if clipboard fails
	} else {
		// Start clipboard monitoring only if init succeeded
		go monitorClipboard(ctx)
		log.Println("Clipboard monitoring started")
	}

	// Create initial menu
	setupMenu()

	// Start menu update routine
	go menuUpdateRoutine(ctx)

	log.Println("MultiClip started successfully")
}

func monitorClipboard(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Clipboard monitoring recovered from panic: %v", r)
		}
	}()

	ch := clipboard.Watch(ctx, clipboard.FmtText)
	
	for {
		select {
		case data := <-ch:
			if data != nil {
				text := string(data)
				if text != "" {
					clipManager.AddItem(text)
					log.Printf("New clipboard item: %.50s...", text)
				}
			}
		case <-ctx.Done():
			log.Println("Clipboard monitoring stopped")
			return
		}
	}
}

func menuUpdateRoutine(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			updateMenuItems()
		case <-ctx.Done():
			log.Println("Menu update routine stopped")
			return
		}
	}
}

func setupMenu() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Menu setup recovered from panic: %v", r)
		}
	}()

	// Clear existing menu items (but not the quit item)
	menuItems = menuItems[:0]

	// Add clipboard history items
	history := clipManager.GetHistory()
	if len(history) == 0 {
		mEmpty := systray.AddMenuItem("No clipboard history", "Clipboard is empty")
		mEmpty.Disable()
		menuItems = append(menuItems, mEmpty)
	} else {
		for i, item := range history {
			// Truncate long items for menu display
			displayText := item
			if len(displayText) > 50 {
				displayText = displayText[:47] + "..."
			}
			
			menuItem := systray.AddMenuItem(displayText, fmt.Sprintf("Copy item %d", i+1))
			menuItems = append(menuItems, menuItem)
			
			// Handle click to copy back to clipboard
			go handleMenuItemClick(item, menuItem)
		}
	}

	// Add separator and quit option (only once)
	if quitItem == nil {
		systray.AddSeparator()
		quitItem = systray.AddMenuItem("Quit MultiClip", "Exit the application")
		
		go func() {
			for {
				select {
				case <-quitItem.ClickedCh:
					log.Println("User requested quit")
					cancel()
					systray.Quit()
					return
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func updateMenuItems() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Menu update recovered from panic: %v", r)
		}
	}()

	history := clipManager.GetHistory()
	
	// Remove old items (except quit)
	for _, item := range menuItems {
		item.Hide()
	}
	menuItems = menuItems[:0]

	// Add current history
	if len(history) == 0 {
		mEmpty := systray.AddMenuItemCheckbox("No clipboard history", "Clipboard is empty", false)
		mEmpty.Disable()
		menuItems = append(menuItems, mEmpty)
	} else {
		for i, item := range history {
			displayText := item
			if len(displayText) > 50 {
				displayText = displayText[:47] + "..."
			}
			
			menuItem := systray.AddMenuItemCheckbox(displayText, fmt.Sprintf("Copy item %d", i+1), false)
			menuItems = append(menuItems, menuItem)
			
			go handleMenuItemClick(item, menuItem)
		}
	}
}

func handleMenuItemClick(text string, menuItem *systray.MenuItem) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Menu item click handler recovered from panic: %v", r)
		}
	}()

	select {
	case <-menuItem.ClickedCh:
		clipboard.Write(clipboard.FmtText, []byte(text))
		log.Printf("Copied back to clipboard: %.50s...", text)
	case <-ctx.Done():
		return
	}
}

func onExit() {
	log.Println("MultiClip shutting down gracefully...")
	cancel() // Cancel all contexts
	time.Sleep(100 * time.Millisecond) // Give goroutines time to cleanup
}

// Simple clipboard icon data (monochrome for template icon)
func getClipboardIcon() []byte {
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x91, 0x68, 0x36, 0x00, 0x00, 0x00,
		0x1C, 0x49, 0x44, 0x41, 0x54, 0x28, 0x91, 0x63, 0x78, 0xCF, 0x80, 0x01,
		0x86, 0x40, 0x92, 0x20, 0x89, 0x81, 0x24, 0x76, 0x92, 0x58, 0x48, 0x62,
		0x27, 0x89, 0x85, 0x24, 0x00, 0x00, 0x0E, 0x30, 0x01, 0x2F, 0xED, 0x83,
		0x34, 0x8C, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42,
		0x60, 0x82,
	}
}