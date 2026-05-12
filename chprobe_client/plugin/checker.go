package plugin

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ricky97gr/chprobe/chprobe_common/typed"
	"github.com/ricky97gr/chprobe/chprobe_common/utils"
)

var pluginBaseDir = "/opt/chprobe/plugins"

func getPluginDir(pluginID string) string {
	return filepath.Join(pluginBaseDir, pluginID)
}

func getPluginZipPath(pluginID string) string {
	return filepath.Join(getPluginDir(pluginID), pluginID+".zip")
}

func CheckPluginExistsAndMD5(pluginID string, expectedMD5 string) (bool, error) {
	pluginDir := getPluginDir(pluginID)
	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	var foundExec bool
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), pluginID) {
			foundExec = true
			break
		}
	}

	if !foundExec {
		utils.Logger.Infof("plugin %s executable not found in dir\n", pluginID)
		return false, nil
	}

	var calculatedMD5 string
	file, err := os.Open(filepath.Join(pluginDir, pluginID))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false, err
	}
	calculatedMD5 = hex.EncodeToString(hash.Sum(nil))

	if calculatedMD5 == expectedMD5 {
		utils.Logger.Infof("plugin %s md5 check passed: %s\n", pluginID, calculatedMD5)
		return true, nil
	}

	utils.Logger.Warnf("plugin %s md5 mismatch, expected: %s, actual: %s\n", pluginID, expectedMD5, calculatedMD5)
	return false, nil
}

func unzipPlugin(zipPath string, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			continue
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		srcFile, err := f.Open()
		if err != nil {
			dstFile.Close()
			return err
		}

		_, err = io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func DownloadPlugin(task *typed.Task) error {
	targetDir := getPluginDir(task.PluginID)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	zipPath := getPluginZipPath(task.PluginID)
	utils.Logger.Infof("start downloading plugin %s zip from %s to %s\n", task.PluginID, task.DownloadURL, zipPath)

	resp, err := http.Get(task.DownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	outFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	utils.Logger.Infof("plugin %s zip downloaded, start unzip\n", task.PluginID)
	if err := unzipPlugin(zipPath, targetDir); err != nil {
		return err
	}

	execPath := filepath.Join(targetDir, task.PluginID)
	os.Chmod(execPath, 0755)

	utils.Logger.Infof("plugin %s unzipped successfully\n", task.PluginID)
	return nil
}

func EnsurePluginReady(task *typed.Task) error {
	ok, err := CheckPluginExistsAndMD5(task.PluginID, task.Md5)
	if err != nil {
		return err
	}
	if !ok {
		return DownloadPlugin(task)
	}
	return nil
}
