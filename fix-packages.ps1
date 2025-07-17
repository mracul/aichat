# For each .go file (excluding .bak), set the package name to the directory name (or 'main' for project root)
Get-ChildItem -Recurse -Filter *.go | Where-Object { $_.FullName -notmatch '\\.bak\\\\' } | ForEach-Object {
    $file = $_.FullName
    $dir = Split-Path $file -Parent
    $base = Split-Path $dir -Leaf

    # If in project root, use 'main'
    if ($dir -eq (Get-Location).Path) {
        $package = 'main'
    } else {
        $package = $base
    }

    $lines = Get-Content $file
    # If the first line is not a package declaration, insert one
    if ($lines.Count -eq 0 -or $lines[0] -notmatch '^\\s*package\\s+') {
        $lines = @("package $package") + $lines
    } else {
        $lines[0] = "package $package"
    }
    Set-Content $file $lines
}
Write-Host "All package declarations have been updated."
