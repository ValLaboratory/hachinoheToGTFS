go generate
go build -ldflags "-s -w"
# ���b�Z�[�W�{�b�N�X�𗘗p���邽�߂̃A�Z���u���̓ǂݍ���(�K�{)
Add-Type -Assembly System.Windows.Forms
# �R�s�[�Ώۃt�@�C��
$target = "hachinoheToGTFS.exe"
# �T�[�o�[�p�X
$path = "\\contents-storage\horus\tools\�H���o�X_�_�C���捞\���ˎs�c_GTFS�R���o�[�^�[\���ؔ�\" + $target
# ���[�J���t�@�C���̃t�@�C���o�[�W����
$v1 = (Get-ItemProperty $target).VersionInfo.FileVersion
if( Test-Path $path ){
  # �T�[�o�[�t�@�C���̃t�@�C���o�[�W����
  $v2 = (Get-ItemProperty $path).VersionInfo.FileVersion
  if ( !( $v1 -eq $v2 ) ) {
    $result = [System.Windows.Forms.MessageBox]::Show("���ؔłɃ����[�X���܂����H",$v1,"YesNo","Question","Button2")
    if ($result -eq "Yes") {
      powershell -Command "./release.bat"
    }
  }
} else {
  $result = [System.Windows.Forms.MessageBox]::Show("���ؔłɃ����[�X���܂����H",$v1,"YesNo","Question","Button2")
  if ($result -eq "Yes") {
    powershell -Command "./release.bat"
  }
}