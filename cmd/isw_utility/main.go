package main

import (
	"context"
	"isw_utility/internal/repository"
	"isw_utility/internal/service"
	"isw_utility/internal/ui"
	"log/slog"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Create repositories
	iswRepo := repository.NewISWRepository()
	// Using ISWRepository for both sensors and control as it implements both
	
	// Create telemetry service
	telemetryService := service.NewTelemetryService(iswRepo, iswRepo, 1*time.Second)

	// Create application and window
	myApp := app.New()
	myApp.Settings().SetTheme(ui.NewArchitectTheme())
	
	myWindow := myApp.NewWindow("ISW UTILITY | TERMINAL ARCHITECT")
	myWindow.Resize(fyne.NewSize(1024, 600))

	// Create Dashboard
	dashboard := ui.NewDashboard(myWindow, telemetryService)
	myWindow.SetContent(dashboard.BuildLayout())

	// Context for background tasks
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Start telemetry loop
	go telemetryService.Start(ctx)

	// Update UI from telemetry updates
	go func() {
		for state := range telemetryService.Updates() {
			dashboard.Update(state)
		}
	}()

	myWindow.ShowAndRun()
}
