# Define all import path replacements as @{ 'old' = 'new' }
$replacements = @{
    'aichat/app' = 'aichat/types'
    'aichat/errors' = 'aichat/errors'
    'aichat/models' = 'aichat/models'
    'aichat/navigation' = 'aichat/navigation'
    'aichat/state' = 'aichat/state'
    'aichat/documentation' = 'aichat/documentation'
    'aichat/components/chatwindow' = 'aichat/components/chatwindow'
    'aichat/components/chat/mvc' = 'aichat/components/chat/mvc'
    'aichat/components/chat/interfaces' = 'aichat/components/chat/interfaces'
    'aichat/components/chat' = 'aichat/components/chat'
    'aichat/components/common' = 'aichat/components/common'
    'aichat/components/flows' = 'aichat/components/flows'
    'aichat/components/input' = 'aichat/components/input'
    'aichat/components/menus' = 'aichat/components/menus'
    'aichat/components/modals/dialogs' = 'aichat/components/modals/dialogs'
    'aichat/components/modals' = 'aichat/components/modals'
    'aichat/components/sidebar/tabs' = 'aichat/components/sidebar/tabs'
    'aichat/components/sidebar' = 'aichat/components/sidebar'
    'aichat/services/ai/providers' = 'aichat/services/ai/providers'
    'aichat/services/ai/types' = 'aichat/services/ai/types'
    'aichat/services/ai' = 'aichat/services/ai'
    'aichat/services/cache' = 'aichat/services/cache'
    'aichat/services/config' = 'aichat/services/config'
    'aichat/services/storage/repositories' = 'aichat/services/storage/repositories'
    'aichat/services/storage' = 'aichat/services/storage'
    'aichat/types/flows' = 'aichat/types/flows'
    'aichat/types/modals' = 'aichat/types/modals'
    'aichat/types/render' = 'aichat/types/render'
    'aichat/types' = 'aichat/types'
    'aichat/views/menu' = 'aichat/views/menu'
    'aichat/views' = 'aichat/views'
}

# Get all .go files except those in .bak/
$files = Get-ChildItem -Recurse -Filter *.go | Where-Object { $_.FullName -notmatch '\\.bak\\' }

foreach ($file in $files) {
    $content = Get-Content $file.FullName -Raw
    foreach ($old in $replacements.Keys) {
        $new = $replacements[$old]
        $content = $content -replace [regex]::Escape($old), $new
    }
    Set-Content $file.FullName $content
}
Write-Host "All import paths updated."
