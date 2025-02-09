package main

import (
    "Tai/internal/application"
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
)

func main() {
    if runtime.GOOS == "windows" {
        cmd := exec.Command("chcp", "65001")
        cmd.Run()
    }
    os.Stdout.WriteString("\xEF\xBB\xBF")

    homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal("Failed to get home directory:", err)
    }
    // Ensure the directory exists (creates all parent directories if necessary)
    err = os.MkdirAll(filepath.Join(homeDir, "TaiInnovation/TaiCli/logs/"), os.ModePerm)
    if err != nil {
        fmt.Println("Error creating directory:", err)
        return
    }
    file, err := os.OpenFile(filepath.Join(homeDir, "TaiInnovation/TaiCli/logs/app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
    if err != nil {
        log.Fatal("Failed to open log file:", err)
    }
    log.SetOutput(file)
    application.Run()
}
