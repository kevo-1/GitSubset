$ErrorActionPreference = "Stop"

$repo = "kevo-1/GitSubset"

# GitHub API to get the latest release
Write-Host "Fetching latest release from $repo..."
$release = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest"

# Find the Windows binary asset
$asset = $release.assets | Where-Object { $_.name -match "windows_amd64|windows_x86_64" -and $_.name -match "\.exe$" }

if (-not $asset) {
    Write-Error "Could not find a Windows .exe binary in the latest release. Are you sure Goreleaser successfully built it?"
    exit 1
}

$downloadUrl = $asset.browser_download_url
$installDir = "$env:LOCALAPPDATA\GitSubset"

# Check for existing installation to clean up
if (Test-Path $installDir) {
    Write-Host "Cleaning up old installation..."
    Remove-Item -Path "$installDir\*" -Recurse -Force -ErrorAction SilentlyContinue
} else {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

$exePath = Join-Path $installDir "gitsubset.exe"

Write-Host "Downloading $downloadUrl..."
Invoke-WebRequest -Uri $downloadUrl -OutFile $exePath

if (-not (Test-Path $exePath)) {
    Write-Error "gitsubset.exe was not found after downloading! Please check the release package."
    exit 1
}

# Update User PATH if needed
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
$escapedInstallDir = [regex]::Escape($installDir)

if ($userPath -notmatch "(^|;)$escapedInstallDir(;|$)") {
    Write-Host "Adding $installDir to your User PATH..."
    $newPath = "$userPath;$installDir"
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
    
    # Update the current session's PATH so the user can use it immediately without restarting
    $env:PATH = "$env:PATH;$installDir"
    
    Write-Host ""
    Write-Host "Installation complete! The tool 'gitsubset' is now available."
    Write-Host "Note: You might need to restart your terminal for changes to fully take effect."
} else {
    Write-Host ""
    Write-Host "Installation complete! 'gitsubset' is already in your PATH."
}

Write-Host ""
Write-Host "Run 'gitsubset' to get started."
