package tunnel

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Tunnel struct {
	proc         *exec.Cmd
	FrpcPath     string
	RemoteHost   string
	RemotePort   int
	LocalHost    string
	LocalPort    int
	ShareToken   string
	stdoutReader io.ReadCloser
	stderrReader io.ReadCloser
}

func (t *Tunnel) String() string {
	result := "Tunnel(\n"
	//追加
	result += fmt.Sprintf("FrpcPath: %s,\n ", t.FrpcPath)
	result += fmt.Sprintf("RemoteHost: %s,\n ", t.RemoteHost)
	result += fmt.Sprintf("RemotePort: %d,\n ", t.RemotePort)
	result += fmt.Sprintf("LocalHost: %s,\n ", t.LocalHost)
	result += fmt.Sprintf("LocalPort: %d,\n", t.LocalPort)
	result += fmt.Sprintf("ShareToken: %s,\n", t.ShareToken)
	result += ")\n"
	return fmt.Sprintf(result)
}

func (t *Tunnel) Start() (string, error) {
	cmdArgs := []string{
		"http",
		"-n",
		t.ShareToken,
		"-l",
		strconv.Itoa(t.LocalPort),
		"-i",
		t.LocalHost,
		"--uc",
		"--sd",
		"random",
		"--ue",
		"--server_addr",
		fmt.Sprintf("%s:%d", t.RemoteHost, t.RemotePort),
		"--disable_log_color",
	}
	t.proc = exec.Command(t.FrpcPath, cmdArgs...)
	log.Printf("cmd:%s\n", t.proc.String())

	// 创建一个管道来读取命令的标准输出
	stdout, err := t.proc.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return "", nil
	}

	// 启动命令
	if err := t.proc.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return "", nil
	}

	url, readErr := t.readURLFromTunnelStream(stdout)
	if readErr != nil {
		return "", readErr
	}
	log.Println("开启代理完成")
	return url, nil
}

// readURLFromTunnelStream reads the URL from the tunnel's stdout stream with timeout handling.
func (t *Tunnel) readURLFromTunnelStream(r io.Reader) (string, error) {
	log.Println("Reading from stream...")

	// Compile regex pattern once outside the loop for efficiency.
	re := regexp.MustCompile(`start proxy success: (.+)`)
	reader := bufio.NewReader(r)
	var url string
	var err error

	// Setup a single-use timer for timeout.
	timeout := time.After(30 * time.Second)

	// Read lines from the stream with timeout.
	for {
		select {
		case <-timeout:
			log.Println("Timeout occurred while reading from stream.")
			err = errors.New("read timeout")
			goto exit
		default:
			line, readErr := reader.ReadString('\n')
			if readErr != nil {
				if readErr != io.EOF { // Ignore EOF which can be a normal termination signal.
					err = readErr
				}
				goto exit
			}
			line = strings.TrimSpace(line) // Remove leading/trailing whitespaces.
			if line == "" {
				continue
			}
			log.Println("Read line:", line)
			if strings.Contains(line, "start proxy success") {
				matches := re.FindStringSubmatch(line)
				if len(matches) == 2 {
					url = matches[1]
					goto exit
				}
			} else if strings.Contains(line, "login to server failed") {
				err = errors.New("login to server failed")
				goto exit
			}
		}
	}

exit:
	log.Println("Read operation completed.")
	return url, err
}
