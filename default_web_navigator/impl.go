package default_web_navigator

import "os/exec"

const (
	windows = "windows"
	linux   = "linux"
	darwin  = "darwin"
)

func OpenURL(os string, url string) *exec.Cmd {
	switch os {
	default:
		return nil
	case windows:
		return exec.Command("cmd", "/C", "start "+url)
	case darwin:
		return exec.Command("open", url)
	}
}
