package purge

import (
	"io/fs"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

func buildDirTreeMap(root string) (map[string]map[string]bool, error) {
    dirContents := make(map[string]map[string]bool)

    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if path == root {
            dirContents[root] = make(map[string]bool)
            return nil
        }

        parent := filepath.Dir(path)
        if _, exists := dirContents[parent]; !exists {
            dirContents[parent] = make(map[string]bool)
        }
        dirContents[parent][path] = d.IsDir()

        if d.IsDir() && !dirExists(path, dirContents) {
            dirContents[path] = make(map[string]bool)
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    return dirContents, nil
}

func dirExists(path string, dirContents map[string]map[string]bool) bool {
    _, ok := dirContents[path]
    return ok
}

func simulateDeletions(changes []Change, dirContents map[string]map[string]bool) {
    for _, c := range changes {
        if c.Type == DeleteFile || c.Type == RemoveDir {
            targetPath, _ := filepath.Abs(c.Target)
            parent := filepath.Dir(targetPath)

            if children, ok := dirContents[parent]; ok {
                delete(children, targetPath)
            }

            if c.Type == RemoveDir {
                delete(dirContents, targetPath)
            }
        }
    }
}

func detectEmptyDirs(dirContents map[string]map[string]bool) []Change {
    var emptyDirs []Change
    processed := make(map[string]bool)

    var dirs []string
    for dir := range dirContents {
        dirs = append(dirs, dir)
    }

    sort.Slice(dirs, func(i, j int) bool {
        return strings.Count(dirs[i], string(os.PathSeparator)) >
            strings.Count(dirs[j], string(os.PathSeparator))
    })

    for _, dir := range dirs {
        if processed[dir] || dir == "." {
            continue
        }

        children := dirContents[dir]
        if len(children) == 0 {
            emptyDirs = append(emptyDirs, Change{
                Type:   RemoveDir,
                Target: dir,
            })
            processed[dir] = true

            parent := filepath.Dir(dir)
            if parentChildren, ok := dirContents[parent]; ok {
                delete(parentChildren, dir)
            }
        }
    }

    return emptyDirs
}

func findEmptyDirs(root string, changes []Change) ([]Change, error) {
    root, err := filepath.Abs(root)
    if err != nil {
        return nil, err
    }

    dirContents, err := buildDirTreeMap(root)
    if err != nil {
        return nil, err
    }

    simulateDeletions(changes, dirContents)
    emptyDirs := detectEmptyDirs(dirContents)

    return emptyDirs, nil
}