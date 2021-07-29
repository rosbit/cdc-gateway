// global conf
// ENV:
//   CONF_FILE      --- 配置文件名
//   TZ             --- 时区名称"Asia/Shanghai"
//
// YAML
// ---
// listen-host: ""
// listen-port: 7080
// apps:
//   - name: app-name
//     service-url: "http://183.62.66.60:9080/cdcserver/api/v2"
//     rsa-private-key: "private-key"
//     aes-key: "aes-key"
//     bank-public-key: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5Ec7viMyQC5SShRz1jP0IQRLLVGDQ4f1rgZwtxT4ZOgnWxUoAHquj2yIrgFjNpWVnt/1dJGtXWkpp2UN3jMI5ubjVQkL0OFD+8r0IFXYAARsCLAwVLF0LE487KvVRaQC7A7rPlFfBtE/v++KajzMuDauNlIASYobcFKYdZ89vIfE/xMg/44QJqQ2XBkoMnJ7ul0kMdh4YWOQnO0qqvXD2eK3KPaXMRtxieGsVBgsvtETprw98bTl9tPUBUrneyirrccS8/Z6raV6nioyx2RzrMld8YnjlnV2YTJpNAlG+y/wLoKY55Rkjcvg9wSe8qbI/VtYVQfQz8gfeUzFQTKKCwIDAQAB"
//     uids:
//        - N002467299
//        - N002467305
//     account:
//        no: 755915709610102
//        bank-no: 75
//   #-
// common-endpoints:
//   health-check: "/health"
//   make-request: "/make-request/:app/:api"
//   parse-response: "/parse-response/:app"
//   get-trans-info: "/get-trans-info/:app/:cdMark"
//
// Rosbit Xu

package gwconf

import (
	"gopkg.in/yaml.v2"
	"path"
	"fmt"
	"os"
	"time"
)

type AppParams struct {
	Name   string `yaml:"name"`
	ServiceURL string `yaml:"service-url"`
	RSAPrivateKey string `yaml:"rsa-private-key"`
	AesKey string `yaml:"aes-key"`
	BankPublicKey string `yaml:"bank-public-key"`
	UIDs []string `yaml:"uids"`
	Account struct {
		No string `yaml:"no"`
		BankNo string `yaml:"bank-no"`
	} `yaml:"account"`
}

type ServiceConfT struct {
	ListenHost     string `yaml:"listen-host"`
	ListenPort     int    `yaml:"listen-port"`
	Apps []AppParams `yaml:"apps"`
	CommonEndpoints struct {
		HealthCheck string `yaml:"health-check"`
		MakeRequest string `yaml:"make-request"`
		ParseResponse string `yaml:"parse-response"`
		GetTransInfo string `yaml:"get-trans-info"`
	} `yaml:"common-endpoints"`
}

var (
	ServiceConf ServiceConfT
	Loc = time.FixedZone("UTC+8", 8*60*60)
)


func getEnv(name string, result *string, must bool) error {
	s := os.Getenv(name)
	if s == "" {
		if must {
			return fmt.Errorf("env \"%s\" not set", name)
		}
	}
	*result = s
	return nil
}

func CheckGlobalConf() error {
	var p string
	getEnv("TZ", &p, false)
	if p != "" {
		if loc, err := time.LoadLocation(p); err == nil {
			Loc = loc
		}
	}

	var confFile string
	if err := getEnv("CONF_FILE", &confFile, true); err != nil {
		return err
	}

	fp, err := os.Open(confFile)
	if err != nil {
		return err
	}
	defer fp.Close()

	dec := yaml.NewDecoder(fp)
	if err := dec.Decode(&ServiceConf); err != nil {
		return err
	}

	if err = checkMust(confFile); err != nil {
		return err
	}

	return nil
}

func DumpConf() {
	fmt.Printf("conf: %v\n", ServiceConf)
	fmt.Printf("TZ time location: %v\n", Loc)
}

func checkMust(confFile string) error {
	// confRoot := path.Dir(confFile)

	if ServiceConf.ListenPort <= 0 {
		return fmt.Errorf("listen-port expected in conf")
	}

	apps := ServiceConf.Apps
	if len(apps) == 0 {
		return fmt.Errorf("apps expected in conf")
	}
	for i, _ := range apps {
		appConf := &apps[i]
		if len(appConf.Name) == 0 {
			return fmt.Errorf("apps[%d]/name expected in conf", i)
		}
		if len(appConf.ServiceURL) == 0 {
			return fmt.Errorf("apps[%d]/service-url expected in conf", i)
		}
		if len(appConf.RSAPrivateKey) == 0 {
			return fmt.Errorf("apps[%d]/rsa-private-key expected in conf", i)
		}
		if len(appConf.AesKey) == 0 {
			return fmt.Errorf("apps[%d]/aes-key expected in conf", i)
		}
		if len(appConf.BankPublicKey) == 0 {
			return fmt.Errorf("apps[%d]/bank-public-key expected in conf", i)
		}
		if len(appConf.UIDs) == 0 {
			return fmt.Errorf("apps[%d]/uids array expected in conf", i)
		}
		acct := &appConf.Account
		if len(acct.No) == 0 {
			return fmt.Errorf("apps[%d]/account/no expected in conf", i)
		}
		if len(acct.BankNo) == 0 {
			return fmt.Errorf("apps[%d]/account/bank-no expected in conf", i)
		}
	}

	ce := &ServiceConf.CommonEndpoints
	if len(ce.HealthCheck) == 0 {
		return fmt.Errorf("common-endpoints/health-check expected in conf")
	}
	if len(ce.MakeRequest) == 0 {
		return fmt.Errorf("common-endpoints/make-request expected in conf")
	}
	if len(ce.ParseResponse) == 0 {
		return fmt.Errorf("common-endpoints/parse-response expected in conf")
	}
	if len(ce.GetTransInfo) == 0 {
		return fmt.Errorf("common-endpoints/get-trans-info expected in conf")
	}

	return nil
}

func checkDir(path, prompt string) error {
	if fi, err := os.Stat(path); err != nil {
		return err
	} else if !fi.IsDir() {
		return fmt.Errorf("%s %s is not a directory", prompt, path)
	}
	return nil
}

func checkFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	return nil
}

func getExecWD() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return path.Dir(exePath), nil
}

func toAbsPath(absRoot, filePath string) string {
	if path.IsAbs(filePath) {
		return filePath
	}
	return path.Join(absRoot, filePath)
}

