package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/pkg/errors"
	"github.com/sirkon/gosrcfmt"

	parser2 "github.com/sirkon/go-imports-rename/internal/parser"
	"github.com/sirkon/go-imports-rename/internal/replacer"
)

type args struct {
	Root string   `arg:"--root" help:"root path to search go files in"`
	Save bool     `arg:"-s,--save" help:"save changes"`
	Rule RuleType `arg:"positional,required" help:"A rule to make import path changes"`
}

func (args) Description() string {
	return "A tool to change import paths based on either prefix switch or regular expressions"
}

func main() {
	var inputArgs args
	inputArgs.Root = "."
	argParse := arg.MustParse(&inputArgs)

	var rep replacer.Replacer
	switch v := inputArgs.Rule.Rule.(type) {
	case parser2.Prefix:
		rep = replacer.Prefix(v.From, v.To)
	case parser2.Add:
		var err error
		rep, err = replacer.Versioned(v.Import, v.Jump)
		if err != nil {
			argParse.Fail(err.Error())
		}
	case parser2.Regexp:
		var err error
		rep, err = replacer.Regexp(v.From, v.To)
		if err != nil {
			argParse.Fail(err.Error())
		}
	}

	logger := newLogger()

	var changesCounter int
	var actualChanges int
	var filesCounter int
	err := filepath.Walk(inputArgs.Root, func(path string, info os.FileInfo, err error) error {
		_, base := filepath.Split(path)
		if info.IsDir() {
			if strings.HasPrefix(base, ".") && len(base) > 1 {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		filesCounter++

		fset := token.NewFileSet()
		goFile, err := parser.ParseFile(fset, path, nil, parser.AllErrors|parser.ParseComments)
		if err != nil {
			logger.Error().Err(err).Msgf("failed to parse %s", path)
			return nil
		}

		var localChanges int
		for _, imp := range goFile.Imports {
			pathValue := strings.Trim(imp.Path.Value, `"`)
			rep := rep.Replace(pathValue)
			switch v := rep.(type) {
			case replacer.Replacement:
				if !inputArgs.Save {
					logger.Info().Msgf("%s: import %s => %s", path, pathValue, v.String())
				} else {
					imp.Path.Value = fmt.Sprintf(`"%s"`, v.String())
				}
				changesCounter++
				localChanges++
			case replacer.Nothing:
				continue
			default:
				logger.Fatal().Msgf("invalid variant case %T", v)
			}
		}

		if inputArgs.Save && localChanges > 0 {
			// create some temporary file in
			fullPath, err := getFullPath(inputArgs.Root, info.Name())
			if err != nil {
				logger.Error().Err(err).Msgf("failed to resolve absolute path of %s", path)
			}
			dir, base := filepath.Split(fullPath)
			file, err := ioutil.TempFile(dir, base)
			if err != nil {
				logger.Error().Err(err).Msgf("failed to update %s", path)
				return nil
			}
			formatted, err := gosrcfmt.AST(fset, goFile)
			if err != nil {
				logger.Error().Err(err).Msg("error when formatting a file")
				return nil
			}
			if _, err := io.Copy(file, bytes.NewBuffer(formatted)); err != nil {
				logger.Error().Err(err).Msgf("error when saving changes to %s", path)
				return nil
			}
			if err := file.Close(); err != nil {
				logger.Error().Err(err).Msgf("something went wrong for %s", path)
			}
			if err := os.Rename(file.Name(), path); err != nil {
				logger.Error().Err(err).Msgf("failed to update %s", path)
			}

			actualChanges += localChanges
		}

		return nil
	})

	var filesMention string
	switch filesCounter {
	case 0:
		logger.Warn().Msgf("no *.go files detected in %s", inputArgs.Root)
		return
	case 1:
		filesMention = "1 *.go file"
	default:
		filesMention = fmt.Sprintf("%d *.go files", filesCounter)
	}

	if changesCounter == 0 {
		logger.Info().Msgf("no changes detected in %d files", filesCounter)
	} else {
		if inputArgs.Save {
			if actualChanges < changesCounter {
				switch actualChanges {
				case 0:
					logger.Warn().Msgf("there were errors on saving changes, noting out of %d was commited in %s", changesCounter, filesMention)
				case 1:
					logger.Warn().Msgf("there were errors on saving changes, only %d out of %d was commited in %s", actualChanges, changesCounter, filesMention)
				default:
					logger.Warn().Msgf("there were errors on saving changes, only %d out of %d were commited in %s", actualChanges, changesCounter, filesMention)
				}
			} else {
				switch changesCounter {
				case 0:
					logger.Info().Msgf("no changes were detected in %s", filesMention)
				case 1:
					logger.Info().Msgf("%d change was detected and commited in %s", changesCounter, filesMention)
				default:
					logger.Info().Msgf("%d changes were detected and commited in %s", changesCounter, filesMention)
				}
			}
		} else {
			switch changesCounter {
			case 0:
				logger.Info().Msgf("no changes were detected in %s", changesCounter, filesMention)
			case 1:
				logger.Info().Msgf("%d change was detected in %s", changesCounter, filesMention)
			default:
				logger.Info().Msgf("%d changes were detected in %s", changesCounter, filesMention)
			}
		}
	}
	if err != nil {
		logger.Error().Err(err).Msgf("failed to scan %s directory tree", inputArgs.Root)
	}
}

func getFullPath(root string, name string) (string, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", errors.WithMessage(err, "full absolute path computation")
	}
	return filepath.Join(rootAbs, name), nil
}
