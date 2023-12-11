package builtins

import (
    "os"
	  "fmt"
)

func RemoveDirectory(dirPath string) error {
    if dirPath == "" {
        return fmt.Errorf("rmdir: missing operand")
    }
    //removes everything
    return os.RemoveAll(dirPath)
}
