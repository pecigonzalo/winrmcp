package winrmcp

import (
	//	"errors"

	"time"

	"github.com/dylanmei/iso8601"
	"github.com/masterzen/winrm/winrm"
)

type Winrmcp struct {
	client *winrm.Client
	config *Config
}

type Config struct {
	Auth                  Auth
	Https                 bool
	Insecure              bool
	CACertBytes           []byte
	OperationTimeout      time.Duration
	MaxOperationsPerShell int
	MaxShell              int
}

type Auth struct {
	User     string
	Password string
}

func NewWinrmcp(addr string, config *Config) (*Winrmcp, error) {
	endpoint, err := parseEndpoint(addr, config.Https, config.Insecure, config.CACertBytes)
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = &Config{}
	}

	params := winrm.DefaultParameters()
	if config.OperationTimeout.Seconds() > 0 {
		params.Timeout = iso8601.FormatDuration(config.OperationTimeout)
	}
	client, err := winrm.NewClientWithParameters(
		endpoint, config.Auth.User, config.Auth.Password, params)
	return &Winrmcp{client, config}, err
}

// func (fs *Winrmcp) Copy(fromPath, toPath string) error {
// 	f, err := os.Open(fromPath)
// 	if err != nil {
// 		return errors.New(fmt.Sprintf("Couldn't read file %s: %v", fromPath, err))
// 	}

// 	defer f.Close()
// 	fi, err := f.Stat()
// 	if err != nil {
// 		return errors.New(fmt.Sprintf("Couldn't stat file %s: %v", fromPath, err))
// 	}

// 	if !fi.IsDir() {
// 		return fs.Write(toPath, f)
// 	} else {
// 		fw := &FileWalker{
// 			client:  fs.client,
// 			config:  fs.config,
// 			toDir:   toPath,
// 			fromDir: fromPath,
// 		}
// 		return filepath.Walk(fromPath, fw.copyFile)
// 	}
// }

// func (fs *Winrmcp) Write(toPath string, src io.Reader) error {
// 	return doCopy(fs.client, fs.config, src, winPath(toPath))
// }

// func (fs *Winrmcp) List(remotePath string) ([]FileItem, error) {
// 	return fetchList(fs.client, winPath(remotePath))
// }
