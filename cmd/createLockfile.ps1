$LolDir = 'C:\Program Files\Tencent\League'
$lockfile = Join-Path -Path $LolDir -ChildPath 'LeagueClient\lockfile'

function Check-If-Admin {
    $currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    return $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Set-Lcu-Info {
    $cmdline = Get-WmiObject -Class Win32_Process -Filter "name='LeagueClientUx.exe'" | Select-Object -Expand CommandLine
    if($cmdline.length -gt 1){
        if($cmdline -match 'app-port=(\d*)'){
            $port = $Matches[1]
        }
        if($cmdline -match 'remoting-auth-token=([\w-]*)'){
            $passwd = $Matches[1]
        }
        return $port+':'+$passwd
    }
}

If (-not (Check-If-Admin)){
    Write-Host "please run this script as admin, or it will not work"
    Write-Host -NoNewLine 'Press any key to continue...';
    $null = $Host.UI.RawUI.ReadKey('NoEcho,IncludeKeyDown');
    Exit(1)
}else{
    Set-Lcu-Info
}

# if(Test-Path -Path $lockfile){
#     Write-Host((Get-Item $lockfile).length/1KB)
# }

function Chech-If-Running {
    $process = Get-Process -Name 'LeagueClientUx' -ErrorAction SilentlyContinue
    if($null -ne $process){return "True"}else{Write-Host "False"}
}