package rules

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type PathInfo struct {
	SrcPath   string
	IsSymLink bool
	FileName  string
}

type PathAction struct {
	DestPaths []string
	ToDelete  bool
	PathInfo  PathInfo
}

func (p PathInfo) matchesDataSize(size DataSize, fileSize int64) bool {
	if !size.Enable || p.IsSymLink {
		return false
	}
	if fileSize < int64(size.Min) {
		return false
	}
	if size.Max != 0 && fileSize > int64(size.Max) {
		return false
	}
	if size.Equal != 0 && fileSize != int64(size.Equal) {
		return false
	}
	return true
}

func addPath(act Action, p PathInfo, m map[string]*PathAction) {
	if _, ok := m[p.SrcPath]; !ok {
		f := PathAction{
			DestPaths: []string{},
			ToDelete:  false,
			PathInfo:  p,
		}
		m[p.SrcPath] = &f
	}

	if act.Delete && !m[p.SrcPath].ToDelete {
		m[p.SrcPath].ToDelete = act.Delete
	}
	if act.Copy != "" {
		m[p.SrcPath].DestPaths = append(m[p.SrcPath].DestPaths, act.Copy)
	}
	if act.Move != "" {
		m[p.SrcPath].DestPaths = append(m[p.SrcPath].DestPaths, act.Move)
		m[p.SrcPath].ToDelete = true
	}
}

func Planner(rules Rules, workingDir string) (map[string]*PathAction, error) {
	actions := map[string]*PathAction{}

	type compiledRule struct {
		rule Rule
		re   *regexp.Regexp
	}

	compiled := []compiledRule{}

	for _, rule := range rules {
		pattern := fmt.Sprintf("%s.*%s$", rule.Options.Name, rule.Options.Extension)

		compiled = append(compiled, compiledRule{
			rule: rule,
			re:   regexp.MustCompile(pattern),
		})
	}

	err := filepath.WalkDir(workingDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		f, err := os.Lstat(p)
		if err != nil {
			return err
		}

		info := PathInfo{
			SrcPath:   p,
			IsSymLink: f.Mode()&fs.ModeSymlink != 0,
			FileName:  f.Name(),
		}

		for _, cr := range compiled {
			if !cr.re.MatchString(d.Name()) {
				continue
			}

			if !info.matchesDataSize(cr.rule.Options.DataSize, f.Size()) {
				continue
			}

			addPath(cr.rule.Options.Action, info, actions)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return actions, nil
}
