package setup

import (
	"fmt"
	"path"
	"path/filepath"
	"time"
)

type TaskSetupCollection struct {
	SetupFolder string
	Tasks       map[string]*TaskSetup
}

func CreateTaskSetupCollection(setupFolder string) (taskSetupCollection *TaskSetupCollection, err error) {
	taskSetupCollection = &TaskSetupCollection{
		SetupFolder: setupFolder,
		Tasks:       make(map[string]*TaskSetup),
	}
	if err = isValidPath(setupFolder, true); err != nil {
		err = fmt.Errorf("path not found %v", err)
		return
	}
	taskSetupCollection.Update()
	return
}

func (coll *TaskSetupCollection) removeUnexistentFiles() (changed bool) {
	unexistentFiles := make([]string, 0)
	for file := range coll.Tasks {
		if isValidPath(file, false) != nil {
			unexistentFiles = append(unexistentFiles, file)
		}
	}
	for _, file := range unexistentFiles {
		delete(coll.Tasks, file)
		changed = true
	}
	return
}

func (coll *TaskSetupCollection) updateFile(filename string) (changed bool) {
	changed = false
	taskSetup, err := ReadTaskSetup(filename)
	if err != nil {
		return
	}
	if existent, ok := coll.Tasks[filename]; ok {
		if existent.taskSetupHash == taskSetup.taskSetupHash {
			// Same content
			return
		}
	}
	coll.Tasks[filename] = taskSetup
	changed = true
	return
}

func (coll *TaskSetupCollection) Update() (changed bool, err error) {
	changed = coll.removeUnexistentFiles()

	var files []string
	if files, err = filepath.Glob(path.Join(coll.SetupFolder, "*.yaml")); err != nil {
		return
	}

	for _, file := range files {
		if coll.updateFile(file) {
			changed = true
		}
	}
	return
}

func (coll *TaskSetupCollection) GetNextTaskToRun() (nextTask *TaskSetup) {
	firstTaskTime := time.Now()
	nextTask = nil
	for _, task := range coll.Tasks {
		if task.isRunning && task.WaitTermination {
			continue
		}
		taskNextRun := task.GetNextRun()
		if taskNextRun.Before(firstTaskTime) {
			nextTask = task
			firstTaskTime = taskNextRun
		}
	}
	return
}
