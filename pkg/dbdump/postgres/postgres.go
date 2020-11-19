package postgres

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Dump provides dump execution arguments.
type Dump struct {
	Host     string
	Username string
	Password string
	Name     string
	Schema   string
	Opts     string
}

func getHostPort(h string) (string, string) {
	host, port := "localhost", "5432"
	data := strings.Split(h, ":")
	host = data[0]
	if len(data) > 1 {
		port = data[1]
	}

	return host, port
}

func (d Dump) Exec() error {

	// Print the version number fo rht ecommand line tools
	cmd := exec.Command("pg_dump", "--version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)
	if err := cmd.Run(); err != nil {
		return err
	}

	flags := []string{}
	if d.Name != "" {
		flags = append(flags, "-d", d.Name)
	}

	host, port := getHostPort(d.Host)
	if host != "" {
		flags = append(flags, "-h", host)
	}
	if port != "" {
		flags = append(flags, "-p", port)
	}

	if d.Username != "" {
		flags = append(flags, "-U", d.Username)
	}

	if d.Opts != "" {
		flags = append(flags, d.Opts)
	}

	if d.Name != "" {
		flags = append(flags, d.Name)
	}

	// gzip > dump.sql.gz
	flags = append(flags, "|", "gzip", ">", "dump.sql.gz")

	envs := []string{}
	if d.Password != "" {
		envs = append(envs, fmt.Sprintf("PGPASSWORD=%s", d.Password))
	}

	cmd = exec.Command("pg_dump", flags...)
	cmd.Env = envs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)
	return cmd.Run()
}

// trace prints the command to the stdout.
func trace(cmd *exec.Cmd) {
	fmt.Printf("$ %s\n", strings.Join(cmd.Args, " "))
}