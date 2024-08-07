package lockscreendetector

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func isScreenLocked() bool {
	var command, testval string

	switch os.Getenv("XDG_CURRENT_DESKTOP") {

	case "Unity":
		command, testval = "gdbus call -e -d com.canonical.Unity -o /com/canonical/Unity/Session -m com.canonical.Unity.Session.IsLocked", "true"

	case "KDE":
		command, testval =
			"qdbus org.kde.screensaver /ScreenSaver org.freedesktop.ScreenSaver.GetActive", "true"

	case "XFCE":
		command, testval =
			"xfconf-query -c xfce4-session -p /general/LockDialogIsVisible", "true"

	case "LXQt":
	case "MATE":
		command, testval =
			"mate-screensaver-command -q", "is active"

	case "Cinnamon":
		command, testval =
			"gdbus call --session --dest org.Cinnamon.ScreenSaver --object-path /org/Cinnamon/ScreenSaver --method org.Cinnamon.ScreenSaver.GetActive", "true"

	case "LXDE":
	case "Deepin":
		command, testval =
			"dbus-send --session --dest=com.deepin.ScreenSaver --type=method_call --print-reply /com/deepin/ScreenSaver com.deepin.ScreenSaver.GetStatus", "true"

	case "Gnome":
	default:
		command, testval =
			"gdbus call --session --dest org.gnome.ScreenSaver --object-path /org/gnome/ScreenSaver --method org.gnome.ScreenSaver.GetActive", "true"
	}

	{
		command := strings.Split(command, " ")
		cmdOutput, err := exec.Command(command[0], command[1:]...).Output()

		if err == nil {
			cmdOutput := string(cmdOutput)
			return strings.Contains(cmdOutput, testval)
		} else {
			log.Println("an error occurred while trying to check if screen was locked: ", err)
			return false
		}
	}
}
