package dir

import (
	"io/fs"
	"os"
	"path/filepath"
)

func Clean(p string) error {
	filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 跳过根目录本身
		if path == p {
			return nil
		}

		// 如果是文件，直接删除
		if !d.IsDir() {
			return os.Remove(path)
		}

		// 如果是目录，删除整个子目录
		return os.RemoveAll(path)
	})

	return nil
}
