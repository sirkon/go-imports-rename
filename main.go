package main

import (
	"encoding"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/pkg/errors"

	"github.com/sirkon/go-imports-rename/internal/replacer"
)

var _ encoding.TextUnmarshaler = &RuleType{}

type RuleType struct {
	From string
	To   string
}

func (r *RuleType) UnmarshalText(rawText []byte) error {
	text := string(rawText)
	pos := strings.Index(text, "=>")
	if pos < 0 {
		return errors.Errorf("\033[1m=>\033[0m required in `From => To`, got `%s` instead", text)
	}
	from := strings.TrimSpace(text[:pos])
	to := strings.TrimSpace(text[pos+2:])
	if len(from) == 0 {
		return errors.Errorf("From rule must not be emptyin `%s`", text)
	}
	if len(to) == 0 {
		return errors.Errorf("To rule must not be empty in `%s`", text)
	}
	r.From = from
	r.To = to
	return nil
}

type args struct {
	Regexp bool     `arg:"-r,--regexp" help:"use regexp to replace import paths"`
	Root   string   `arg:"--root" help:"root path to search go files in"`
	Save   bool     `arg:"-s,--save" help:"save changes"`
	Rule   RuleType `arg:"positional,required" help:"From => To, where From and To must be either import prefix to switch or (Regexp, Replacement) couple"`
}

func (args) Description() string {
	return "A tool to change import paths based on either prefix switch or regular expressions"
}

func main() {
	var inputArgs args
	inputArgs.Root = "."
	argParse := arg.MustParse(&inputArgs)

	var rep replacer.Replacer
	if inputArgs.Regexp {
		var err error
		rep, err = replacer.RegexpReplace(inputArgs.Rule.From, inputArgs.Rule.To)
		if err != nil {
			argParse.Fail(err.Error())
		}
	} else {
		rep = replacer.PrefixReplace(inputArgs.Rule.From, inputArgs.Rule.To)
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

		var fset token.FileSet
		ast, err := parser.ParseFile(&fset, path, nil, parser.AllErrors|parser.ParseComments)
		if err != nil {
			logger.Error().Err(err).Msgf("failed to parse %s", path)
			return nil
		}

		var localChanges int
		for _, imp := range ast.Imports {
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
			fullRelPath := filepath.Join(inputArgs.Root, info.Name())
			fullPath, err := filepath.Abs(fullRelPath)
			if err != nil {
				logger.Error().Err(err).Msgf("failed to resolve absolute path of %s", path)
			}
			dir, _ := filepath.Split(fullPath)
			file, err := ioutil.TempFile(dir, path)
			if err != nil {
				logger.Error().Err(err).Msgf("failed to update %s", path)
				return nil
			}
			if err := printer.Fprint(file, &fset, ast); err != nil {
				logger.Error().Err(err).Msgf("error when saving changes to %s", path)
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
