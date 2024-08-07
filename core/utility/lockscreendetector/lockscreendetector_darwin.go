package lockscreendetector

func isScreenLocked() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	output, _ := exec.CommandContext(ctx, "ioreg", "ioreg", "-n", "Root", "-d1").Output()
	return strings.Contains(string(output), "CGSSessionScreenIsLocked")
}
