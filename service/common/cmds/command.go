package cmds

import (
	"fmt"
	"github.com/v2rayA/v2rayA/pkg/util/log"
	"os"
	"os/exec"
	"strings"
)

func IsCommandValid(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func ExecCommands(commands string, stopWhenError bool) error {
	lines := strings.Split(commands, "\n")
	var e error
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) <= 0 || strings.HasPrefix(line, "#") {
			continue
		}
		log.Alert("%v", line)
		// 将命令行追加到日志文件
		f, err := os.OpenFile("/tmp/v2.shell.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Trace("打开日志文件失败: %v", err)
		} else {
			_, err = f.WriteString(line + "\n")
			if err != nil {
				log.Trace("写入日志失败: %v", err)
			}
			f.Close()
		}
		// 执行命令
		out, err := exec.Command("sh", "-c", line).CombinedOutput()
		if err != nil {
			e = fmt.Errorf("ExecCommands: %v %v: %w", line, string(out), err)
			if stopWhenError {
				log.Trace("%v", e)
				return e
			}
		}
	}
	return e
}
