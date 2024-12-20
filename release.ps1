go generate
go build -ldflags "-s -w"
# メッセージボックスを利用するためのアセンブリの読み込み(必須)
Add-Type -Assembly System.Windows.Forms
# コピー対象ファイル
$target = "hachinoheToGTFS.exe"
# サーバーパス
$path = "\\contents-storage\horus\tools\路線バス_ダイヤ取込\八戸市営_GTFSコンバーター\検証版\" + $target
# ローカルファイルのファイルバージョン
$v1 = (Get-ItemProperty $target).VersionInfo.ProductVersion
if( Test-Path $path ){
  # サーバーファイルのファイルバージョン
  $v2 = (Get-ItemProperty $path).VersionInfo.ProductVersion
  if ( !( $v1 -eq $v2 ) ) {
    $result = [System.Windows.Forms.MessageBox]::Show("検証版にリリースしますか？",$v1,"YesNo","Question","Button2")
    if ($result -eq "Yes") {
      $command = "copy " + $target + " " + $path
      powershell -Command $command
    }
  }
} else {
  $result = [System.Windows.Forms.MessageBox]::Show("検証版にリリースしますか？",$v1,"YesNo","Question","Button2")
  if ($result -eq "Yes") {
    $command = "copy " + $target + " " + $path
    powershell -Command $command
  }
}