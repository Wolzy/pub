package goftp

import (
	_ "errors"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// FTP 文件对象
type FileSource struct {
	entry *ftp.Entry // ftp库的entry对象
	path  string     // 文件的全路径+文件名
}

// EntryHandler 遍历ftp目录时的文件handler
type EntryHandler func(e *ftp.Entry, currentPath string) error

// FTP文件信息
type FtpFile struct {
	FileName string //FTP文件名
	Path     string //FTP文件的全路径+文件名
	Type     int    //FTP文件类型，文件:0, 文件夹:1
	Size     int    //FTP文件大小
}

type FtpClient struct {
	Protocol string
	Ftp      *ftp.ServerConn
	Ssh      *ssh.Client
	Sftp     *sftp.Client
	fileChan chan interface{}
	wg       sync.WaitGroup
}

func NewFtpClient() *FtpClient {
	return &FtpClient{}
}

func (this *FtpClient) Init(protocol, ip, port, user, pass string) error {
	var err error
	this.Protocol = protocol
	if "FTP" == protocol {
		this.Ftp, err = ftp.Connect(ip + ":" + port)
		if err != nil {
			return err
		}
		err = this.Ftp.Login(user, pass)
		if err != nil {
			return err
		}

		this.fileChan = make(chan interface{}, DataSize)
		this.wg = sync.WaitGroup{}
	} else if "SFTP" == protocol {
		var (
			auth         []ssh.AuthMethod
			addr         string
			clientConfig *ssh.ClientConfig
		)
		// get auth method
		auth = make([]ssh.AuthMethod, 0)
		auth = append(auth, ssh.Password(pass))

		clientConfig = &ssh.ClientConfig{
			User:            user,
			Auth:            auth,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         30 * time.Second,
		}
		addr = fmt.Sprintf("%s:%s", ip, port)
		// fmt.Println("Prepare to connnect ", protocol, ip, port)

		if this.Ssh, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
			return err
		}
		// fmt.Println("Estalished ssh channel, prepare sftp client...")
		// open an SFTP session over an existing ssh connection.
		this.Sftp, err = sftp.NewClient(this.Ssh)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (this *FtpClient) Finish() error {
	if nil != this.Ftp {
		// err := this.Ftp.Logout()
		err := this.Ftp.Quit()
		if err != nil {
			return err
		}
	}

	if nil != this.Sftp {
		err := this.Sftp.Close()
		if err != nil {
			return err
		}
	}

	if nil != this.Ssh {
		err := this.Ssh.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// 函调函数
func (this *FtpClient) Handler(e *ftp.Entry, currentPath string) error {

	stru := FtpFile{}
	stru.FileName = e.Name
	stru.Path = currentPath + "//" + e.Name
	stru.Type = int(e.Type)
	stru.Size = int(e.Size)
	select {
	case this.fileChan <- stru:
		//global.Log.Debug("fileChandata: %v", stru)
	default:
		log.Println("fileChan data chan is full")
		time.Sleep(time.Second)
		break
	}
	return nil
}

// M1: 遍历ftp目录，获取文件
func (this *FtpClient) Walk(rootDir string) error {
	entries, err := this.Ftp.List(rootDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		switch entry.Type {
		case ftp.EntryTypeFile:
			//正在上传的文件，先不进行下载
			if entry.Size != 0 {
				stru := FtpFile{}
				stru.FileName = entry.Name
				stru.Path = rootDir + "//" + entry.Name
				stru.Type = int(entry.Type)
				stru.Size = int(entry.Size)
				if len(this.fileChan) > (DataSize - 1) {
					log.Println("管道大小超限，完成本地扫描:", len(this.fileChan))
					return nil
				}
				select {
				case this.fileChan <- stru:
					//global.Log.Debug("fileChan: %v", stru.FileName)
				default:
					log.Println("fileChan data chan is full")
					time.Sleep(time.Second)
					return nil
				}
			}
		case ftp.EntryTypeFolder:
			continue
			// err = this.walkthrough(fmt.Sprintf("%s/%s", rootDir, entry.Name))
		default:
		}
	}
	return nil
}

// M2: 遍历ftp目录，回调获取文件
func (this *FtpClient) walkCall(rootDir string, handler EntryHandler) error {
	entries, err := this.Ftp.List(rootDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		switch entry.Type {
		case ftp.EntryTypeFile:
			//正在上传的文件，先不进行下载
			if entry.Size != 0 {
				err = handler(entry, rootDir)
			}
		case ftp.EntryTypeFolder:
			// 子目录不做处理
			// err = this.walkCall(fmt.Sprintf("%s/%s", rootDir, entry.Name), handler)
		default:
		}
	}
	return nil
}

// M3: 遍历ftp目录，获取文件
func (this *FtpClient) listfiles(rootDir string) error {
	err := this.walkCall(rootDir, func(entry *ftp.Entry, currentPath string) error {
		stru := FtpFile{}
		stru.FileName = entry.Name
		stru.Path = currentPath + "//" + entry.Name
		stru.Type = int(entry.Type)
		stru.Size = int(entry.Size)
		if len(this.fileChan) > DataSize {
			return nil
		}
		select {
		case this.fileChan <- stru:
			fmt.Println("fileChan:", stru.FileName)
		default:
			fmt.Println("fileChan data chan is full")
			time.Sleep(time.Second)
			break
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 函调函数-上传
func (this *FtpClient) HandlerUpload(e *ftp.Entry, currentPath string) error {

	stru := FtpFile{}
	stru.FileName = e.Name
	stru.Path = currentPath + "//" + e.Name
	stru.Type = int(e.Type)
	stru.Size = int(e.Size)
	select {
	case this.fileChan <- stru:
		//global.Log.Debug("fileChandata: %v", stru)
	default:
		log.Println("fileChan data chan is full")
		time.Sleep(time.Second)
		break
	}
	return nil
}

// 获取文件列表
func (this *FtpClient) GetFile(path string) ([]string, error) {
	fileList := make([]string, 0, 0)

	if nil != this.Ftp {
		entries, _ := this.Ftp.List(path)
		for _, entry := range entries {
			switch entry.Type {
			case ftp.EntryTypeFile:
				//正在上传的文件，先不进行下载(0字节)
				if entry.Size == 0 {
					continue
				}
				fileList = append(fileList, entry.Name)
			case ftp.EntryTypeFolder:
				// 子目录不做处理
				// err = this.walkCall(fmt.Sprintf("%s/%s", rootDir, entry.Name), handler)
			default:
			}
		}
	} else if nil != this.Sftp {
		w := this.Sftp.Walk(path)
		for w.Step() {
			if w.Err() != nil {
				break
			}
			if true == w.Stat().IsDir() {
				continue
			}
			fileList = append(fileList, w.Stat().Name())
		}
	}

	return fileList, nil
}

// 文件下载
func (this *FtpClient) Download(local, remote, file string) error {
	localfile := path.Join(local, file)
	remotefile := path.Join(remote, file)

	if "FTP" == this.Protocol {
		r, err := this.Ftp.Retr(remotefile)
		defer r.Close()

		fp, _ := os.Create(localfile)
		defer fp.Close()

		if err != nil {
			return err
		} else {
			for {
				buf := make([]byte, 1024)
				n, _ := r.Read(buf)
				if n == 0 {
					break
				}
				fp.Write(buf[:n])
			}
		}
		r.Close()
	} else if "SFTP" == this.Protocol {
		r, err := this.Sftp.Open(remotefile)
		if err != nil {
			log.Println("Rmtfile:", remotefile, "open error, ", err)
			return err
		}
		defer r.Close()

		fp, _ := os.Create(localfile)
		defer fp.Close()

		_, err = r.WriteTo(fp)
		if err != nil {
			log.Println("Rmtfile:", remotefile, "down error, ", err)
			return err
		}
		fp.Close()
		r.Close()
	}

	return nil
}

func (this *FtpClient) sftpUpload(local, remote, filename string) error {
	var (
		err    error
		ftpcli = this.Sftp
	)
	// 用来测试的本地文件路径 和 远程机器上的文件夹
	var localPath = path.Join(local, filename)
	srcFile, err := os.Open(localPath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	dstFile, err := ftpcli.Create(path.Join(remote, filename))
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	buf := make([]byte, 1048576)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[:n])
	}

	return nil
}

func (this *FtpClient) Upload(local, remote, file string) error {
	localfile := path.Join(local, file)
	remotefile := path.Join(remote, file)
	fp, err := os.Open(localfile)
	if err != nil {
		return err
	}

	if "FTP" == this.Protocol {
		err = this.Ftp.Stor(remotefile, fp)
		if err != nil {
			return err
		}
	} else if "SFTP" == this.Protocol {
		err = this.sftpUpload(local, remote, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *FtpClient) Rename(tmp, dst, file string) error {
	tmpfile := path.Join(tmp, file)
	dstfile := path.Join(dst, file)
	var err error

	if "FTP" == this.Protocol {
		err = this.Ftp.Rename(tmpfile, dstfile)
		if err != nil {
			return err
		}
	} else if "SFTP" == this.Protocol {
		err = this.Sftp.Rename(tmpfile, dstfile)
		if err != nil {
			return err
		}
	}

	return nil
}
