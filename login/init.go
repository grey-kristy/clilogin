package login

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/markbates/goth"
)

func InitGoogleAuth(stop chan CallbackResponse) error {
	err := InitProviders()
	if err != nil {
		return err
	}

	url, err := GetAuthURL()
	if err != nil {
		return err
	}

	RunHTTPServer(stop)
	return OpenBrowser(url)
}

func GetAuthURL() (string, error) {
	provider, err := goth.GetProvider("google")
	if err != nil {
		return "", err
	}
	sess, err := provider.BeginAuth("state1")
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	return url, err
}

// OpenBrowser opens up the provided URL in a browser
func OpenBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "openbsd":
		fallthrough
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		r := strings.NewReplacer("&", "^&")
		cmd = exec.Command("cmd", "/c", "start", r.Replace(url))
	}
	if cmd != nil {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			return fmt.Errorf("failed to open browser: " + err.Error())
		}
		err = cmd.Wait()
		if err != nil {
			return fmt.Errorf("failed to wait for open browser command to finish: " + err.Error())
		}
		return nil
	} else {
		return errors.New("unsupported platform")
	}
}
