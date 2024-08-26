package nf

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"
)

func ProcessDirectory(dir string) ([]*Module, error) {
	var modules []*Module
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".nf" {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				module, err := BuildModule(path)
				if err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf("%s: %v", path, err))
					mu.Unlock()
					return
				}
				mu.Lock()
				modules = append(modules, module)
				mu.Unlock()
			}(path)
		}
		return nil
	})

	wg.Wait()

	if err != nil {
		return nil, err
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("encountered %d errors: %v", len(errors), errors)
	}

	return modules, nil
}
