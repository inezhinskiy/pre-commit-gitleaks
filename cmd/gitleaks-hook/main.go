package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	autoinstall := gitConfigBool("hooks.gitleaks.autoinstall")

	if _, err := exec.LookPath("gitleaks"); err != nil {
		if !autoinstall {
			fail("gitleaks не знайдено в PATH.\n" +
				"  Встанови вручну: https://github.com/gitleaks/gitleaks#installing\n" +
				"  або увімкни автоматичне встановлення:\n" +
				"    git config hooks.gitleaks.autoinstall true")
		}
		fmt.Println("[pre-commit:gitleaks] gitleaks не знайдено — встановлюю автоматично...")
		if err := installGitleaks(); err != nil {
			fail(fmt.Sprintf("не вдалося встановити gitleaks: %v", err))
		}
	}

	cmd := exec.Command("gitleaks", "protect", "--staged", "--redact", "--no-banner")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fail("gitleaks виявив секрети у staged-змінах. Коміт відхилено.")
	}

	fmt.Println("[pre-commit:gitleaks] секретів не знайдено, коміт дозволено.")
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, "\n[pre-commit:gitleaks] "+msg)
	os.Exit(1)
}

func gitConfigBool(key string) bool {
	out, err := exec.Command("git", "config", "--get", key).Output()
	if err != nil {
		return false
	}
	return string(out) == "true\n"
}

func installGitleaks() error {
	switch runtime.GOOS {
	case "darwin":
		if _, err := exec.LookPath("brew"); err == nil {
			return runShell("brew install gitleaks")
		}
	case "linux":
		if _, err := exec.LookPath("apt-get"); err == nil {
			return runShell("sudo apt-get update && sudo apt-get install -y gitleaks")
		}
	case "windows":
		if _, err := exec.LookPath("choco"); err == nil {
			return runShell("choco install gitleaks -y")
		}
		if _, err := exec.LookPath("scoop"); err == nil {
			return runShell("scoop install gitleaks")
		}
	}
	if _, err := exec.LookPath("go"); err == nil {
		return runShell("go install github.com/zricethezav/gitleaks/v8@latest")
	}
	return fmt.Errorf("не знайдено ні пакетного менеджера, ні Go — встанови gitleaks вручну")
}

func runShell(command string) error {
	shell, flag := "sh", "-c"
	if runtime.GOOS == "windows" {
		shell, flag = "cmd", "/C"
	}
	cmd := exec.Command(shell, flag, command)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}
