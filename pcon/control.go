package main

import (
	"os"
	"runtime"
	"strings"
	"time"

	installer "github.com/weareplanet/ifcv5-main/installer/handler"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/pkg/errors"
)

var (
	cmdDelay = 500 * time.Millisecond
)

func (a *App) controlHandler(cmd string, _ string, _ uint64) error {

	var err error

	log.Debug().Msgf("%T received control command '%s'", a, cmd)

	switch cmd {

	case "update":

		switch runtime.GOOS {

		case "windows", "linux":

			err = a.update()
			if err == nil {
				a.restart()
			} else {
				log.Error().Msgf("%T update failed, err=%s", a, err)
			}

		default:

			err = errors.Errorf("update not supported on platform '%s'", runtime.GOOS)

		}

	case "restart":

		a.restart()

	default:
		err = errors.Errorf("unknown control command '%s'", cmd)
	}

	return err
}

func (a *App) update() error {

	target, err := os.Executable()
	if err != nil {
		return err
	}

	// normalise windows path
	target = strings.Replace(target, "\\", "/", -1)

	currentVersion := installer.GetVersion(target)
	currentHash := installer.GetHash(target)
	driver := installer.GetFileName(target)

	log.Debug().Msgf("%T binary '%s', current version '%s', current hash '%s'", a, driver, currentVersion, currentHash)
	if update, statusCode, err := installer.Install(target, driver, currentVersion, currentHash, nil, nil); err != nil {
		if statusCode == 404 {
			return errors.Errorf("driver '%s' not available", driver)
		}
		return errors.Errorf("install '%s' failed, err=%s", driver, err)
	} else if update {
		version := installer.GetVersion(target)
		hash := installer.GetHash(target)
		log.Debug().Msgf("%T downloaded version '%s', hash '%s'", a, version, hash)
	} else {
		return errors.Errorf("the latest version '%s' is already running", currentVersion)
	}

	return nil
}

func (a *App) restart() {
	go func() {
		time.Sleep(cmdDelay)
		a.Close()
		// todo: save handling -> close log
		os.Exit(0)
	}()
}
