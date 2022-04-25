package setup

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/guionardo/gs-service-scheduler/logging"
	"github.com/robfig/cron/v3"
)

type TaskSetup struct {
	Name                string   `yaml:"name"`
	Schedule            string   `yaml:"schedule"`
	ExecutionPath       string   `yaml:"execution_path"`
	ExecutablePath      string   `yaml:"executable_path"`
	Args                []string `yaml:"args"`
	CaptureOutputToFile string   `yaml:"capture_output_to_file"`
	WaitTermination     bool     `yaml:"wait_termination"`

	schedule      cron.Schedule `yaml:"-"`
	lastRun       time.Time     `yaml:"-"`
	isRunning     bool          `yaml:"-"`
	taskSetupFile string        `yaml:"-"`
	taskSetupHash string        `yaml:"-"`
	lock          sync.RWMutex  `yaml:"-"`
}

func (taskSetup *TaskSetup) ToString() string {
	return fmt.Sprintf("%s: ", taskSetup.Name)
}

func ReadTaskSetup(filename string) (taskSetup *TaskSetup, err error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(content, &taskSetup); err != nil {
		return
	}
	if len(taskSetup.Name) == 0 {
		err = fmt.Errorf("invalid file %s : Empty 'name'", filename)
		return
	}
	parser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	schedule, err := parser.Parse(taskSetup.Schedule)
	if err != nil {
		return nil, err
	}
	taskSetup.schedule = schedule
	taskSetup.lastRun = time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
	taskSetup.isRunning = false

	hash := md5.Sum(content)
	taskSetup.taskSetupHash = hex.EncodeToString(hash[:])
	taskSetup.taskSetupFile = filename
	return
}

func (taskSetup *TaskSetup) GetNextRun() time.Time {
	return taskSetup.schedule.Next(taskSetup.lastRun)
}

func isValidPath(pathName string, isDir bool) error {
	stat, err := os.Stat(pathName)
	if err == nil {
		if stat.IsDir() != isDir {
			if isDir {
				err = fmt.Errorf("expected path %s", pathName)
			} else {
				err = fmt.Errorf("expected file %s", pathName)
			}
		}
	} else {
		if isDir {
			err = fmt.Errorf("path not found %s", pathName)
		} else {
			err = fmt.Errorf("file not found %s", pathName)
		}
	}
	return err
}

func (taskSetup *TaskSetup) ValidatePaths() error {
	executableDir, executableFile := path.Split(taskSetup.ExecutablePath)
	if len(executableDir) == 0 {
		if len(taskSetup.ExecutionPath) > 0 {
			executableDir = taskSetup.ExecutionPath
		} else {
			return fmt.Errorf("expected existent execution path")
		}
	}
	if err := isValidPath(executableDir, true); err != nil {
		return fmt.Errorf("expected executable path exists: %s - %v", executableDir, err)
	}

	taskSetup.ExecutablePath = path.Join(executableDir, executableFile)
	if err := isValidPath(taskSetup.ExecutablePath, false); err != nil {
		return fmt.Errorf("expected executable path exists: %s - %v", taskSetup.ExecutablePath, err)
	}

	if len(taskSetup.ExecutionPath) == 0 {
		taskSetup.ExecutionPath = executableDir
	}
	if err := isValidPath(taskSetup.ExecutionPath, true); err != nil {
		return fmt.Errorf("expected execution path exists: %s", taskSetup.ExecutionPath)
	}
	return nil
}

func (taskSetup *TaskSetup) Run() {

	if taskSetup.GetNextRun().After(time.Now()) {
		return
	}
	if taskSetup.isRunning && taskSetup.WaitTermination {
		logging.InfoF("%s is running previous schedule", taskSetup)
		return
	}
	if err := taskSetup.ValidatePaths(); err != nil {
		logging.ErrorF("%s in invalid state", taskSetup)
		return
	}
	if taskSetup.WaitTermination {
		taskSetup.lock.Lock()
		defer taskSetup.lock.Unlock()
	}
	

	// out,err:=exec.Command()
	//TODO: Implementar execução da task

}

func StartProccess(pwd string, command string, args ...string) (p *os.Process, err error) {
	if args[0], err = exec.LookPath(args[0]); err == nil {
		procAttr:=os.ProcAttr{
			Files: []*os.File{
				os.Stdin,
				os.Stdout,
				os.Stderr,
			},
			Dir: pwd,			
		}
		
		p, err := os.StartProcess(args[0], args, &procAttr)
		if err == nil {
			return p, nil
		}
	}
	return nil, err
}
