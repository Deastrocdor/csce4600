package builtins

import (
    "os"
	  "fmt"
)

func MakeDirectory(dirPath string, mode os.FileMode) error {
    if dirPath == "" {
        return fmt.Errorf("mkdir: missing operand")
    }
    // defaults 0777 
    if mode == 0 {
        mode = 0777
    }

    return os.MkdirAll(dirPath, mode)
}
