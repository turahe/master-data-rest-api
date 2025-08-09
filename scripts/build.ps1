# PowerShell build script for Windows
param(
    [string]$Version = "1.0.0",
    [string]$BuildDir = "bin",
    [switch]$Clean
)

Write-Host "Master Data REST API - Build Script" -ForegroundColor Green
Write-Host "====================================" -ForegroundColor Green

# Clean build directory if requested
if ($Clean) {
    Write-Host "Cleaning build directory..." -ForegroundColor Yellow
    if (Test-Path $BuildDir) {
        Remove-Item -Path "$BuildDir/*" -Force -Recurse -ErrorAction SilentlyContinue
    }
}

# Create build directory if it doesn't exist
if (!(Test-Path $BuildDir)) {
    New-Item -ItemType Directory -Path $BuildDir -Force | Out-Null
}

# Set build information
$BuildTime = Get-Date -Format "yyyy-MM-dd_HH:mm:ss"
$GitCommit = try { git rev-parse --short HEAD 2>$null } catch { "unknown" }

Write-Host "Building version: $Version" -ForegroundColor Cyan
Write-Host "Build time: $BuildTime" -ForegroundColor Cyan
Write-Host "Git commit: $GitCommit" -ForegroundColor Cyan

# Build flags
$LDFlags = "-X main.version=$Version -X main.buildTime=$BuildTime -X main.gitCommit=$GitCommit"

# Build Windows executable
Write-Host "`nBuilding Windows executable..." -ForegroundColor Yellow
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags $LDFlags -o "$BuildDir/master-data-api.exe" main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Windows build completed successfully!" -ForegroundColor Green
} else {
    Write-Host "Windows build failed!" -ForegroundColor Red
    exit 1
}

# Build Linux executable
Write-Host "`nBuilding Linux executable..." -ForegroundColor Yellow
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags $LDFlags -o "$BuildDir/master-data-api" main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Linux build completed successfully!" -ForegroundColor Green
} else {
    Write-Host "Linux build failed!" -ForegroundColor Red
    exit 1
}

# Build macOS executable
Write-Host "`nBuilding macOS executable..." -ForegroundColor Yellow
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -ldflags $LDFlags -o "$BuildDir/master-data-api-darwin" main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "macOS build completed successfully!" -ForegroundColor Green
} else {
    Write-Host "macOS build failed!" -ForegroundColor Red
    exit 1
}

# Reset environment variables
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue

Write-Host "`nAll builds completed successfully!" -ForegroundColor Green
Write-Host "`nBuilt files:" -ForegroundColor Cyan
Get-ChildItem -Path $BuildDir | ForEach-Object {
    $size = [math]::Round($_.Length / 1MB, 2)
    Write-Host "  File: $($_.Name) (Size: $size MB)" -ForegroundColor White
}

Write-Host "`nTo run the application:" -ForegroundColor Yellow
Write-Host "  Windows: .\$BuildDir\master-data-api.exe --help" -ForegroundColor White
Write-Host "  Linux:   ./$BuildDir/master-data-api --help" -ForegroundColor White
Write-Host "  macOS:   ./$BuildDir/master-data-api-darwin --help" -ForegroundColor White