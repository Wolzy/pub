package goftp

const (
	DataSize = 1024
)

type Server struct {
	Protocol   string `json:"Protocol"`
	IP         string `json:"IP"`
	Port       string `json:"Port"`
	UserName   string `json:"UserName"`
	PassWord   string `json:"PassWord"`
	RemotePath string `json:"RemotePath"` // 远程临时目录（上传用）
	TargetPath string `json:"TargetPath"` // 远程目标目录（上传用，下载时为目标目录）
	RmtTmpPath string `json:"RmtTmpPath"` // 远程临时目录（上传用，向下兼容）
	BackupPath string `json:"BackupPath"` // 备份目录（上传后将本地文件转移到的目标目录，下载后将远程文件转移到的目标目录）
	LocalPath  string `json:"LocalPath"`  // 本地目标目录（下载为目标目录，上传为源目录）
	TempPath   string `json:"TempPath"`   // 本地临时目录（下载为临时目录，上传不用）
}
