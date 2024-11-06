package mcd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/gobwas/glob"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"io"
	"os"
	"path"
	"sort"
	"strings"
)

type PackageJSONWorkspaces struct {
	Workspaces []string `json:"workspaces"`
}

type PnpmWorkspacesPackages struct {
	Packages []string `json:"packages"`
}

type MonocdrcInclude struct {
	Include []string `json:"include"`
}

type WalkFunc func(path string)

func Walk(root string, absolutePath string, excludeGlobs []glob.Glob, negativeGlobs []glob.Glob, positiveGlobs []glob.Glob, candidates []Candidate) []Candidate {
	isNegativeMatched := false
	isPositiveMatched := false

	p := strings.Replace(absolutePath, root+"/", "", 1)

	for _, g := range negativeGlobs {
		isNegativeMatched = g.Match(p)

		if isNegativeMatched {
			break
		}
	}

	if !isNegativeMatched {
		for _, g := range positiveGlobs {
			isPositiveMatched = g.Match(p)

			if isPositiveMatched {
				break
			}
		}
	}

	if !isNegativeMatched && isPositiveMatched {
		candidates = append(candidates, Candidate{
			name: p,
			path: absolutePath,
		})
	}

	isExcludeMatched := false

	for _, g := range excludeGlobs {
		isExcludeMatched = g.Match(p)

		if isExcludeMatched {
			break
		}
	}

	if !isExcludeMatched {
		entries, _ := os.ReadDir(absolutePath)

		for _, entry := range entries {
			if entry.IsDir() {
				nextPath := path.Join(absolutePath, entry.Name())
				candidates = Walk(root, nextPath, excludeGlobs, negativeGlobs, positiveGlobs, candidates)
			}
		}
	}

	return candidates
}

func getCandidates() []Candidate {
	candidates := make([]Candidate, 0, 32)

	workingDir, err := os.Getwd()

	if err != nil {
		fmt.Fprintf(os.Stderr, "COULD NOT FIND CURRENT WORKING DIRECTORY\n")
		return candidates
	}

	currentDirectory := workingDir
	workspacePath := currentDirectory
	workspaceGlobs := make([]string, 0, 32)

	workspaceGlobs = append(workspaceGlobs)

	for {
		packageJSONPath := path.Join(currentDirectory, "./package.json")
		packageJSONExists, _ := fileExists(packageJSONPath)

		if packageJSONExists {
			file, fileOpenError := os.Open(packageJSONPath)

			if fileOpenError != nil {
				fmt.Fprintf(os.Stderr, "could not open file %s \n", packageJSONPath)
				return candidates
			}

			reader := bufio.NewReader(file)
			packageJSONContents, fileReadError := io.ReadAll(reader)
			file.Close()

			if fileReadError != nil {
				fmt.Fprintf(os.Stderr, "could not read file %s \n", packageJSONPath)
				return candidates
			}

			var packageJSONData map[string]interface{}
			json.Unmarshal(packageJSONContents, &packageJSONData)

			if _, ok := packageJSONData["workspaces"]; ok {
				packageJSONWorkspaces := PackageJSONWorkspaces{}
				json.Unmarshal(packageJSONContents, &packageJSONWorkspaces)

				workspacePath = currentDirectory
				workspaceGlobs = packageJSONWorkspaces.Workspaces
			}
		}

		pnpmWorkspaceYAMLPath := path.Join(currentDirectory, "./pnpm-workspace.yaml")
		pnpmWorkspaceYAMLExists, _ := fileExists(pnpmWorkspaceYAMLPath)

		if pnpmWorkspaceYAMLExists {
			file, fileOpenError := os.Open(pnpmWorkspaceYAMLPath)

			if fileOpenError != nil {
				fmt.Fprintf(os.Stderr, "could not open file %s \n", pnpmWorkspaceYAMLPath)
				return candidates
			}

			reader := bufio.NewReader(file)
			pnpmWorkspaceYAMLPathContents, fileReadError := io.ReadAll(reader)
			file.Close()

			if fileReadError != nil {
				fmt.Fprintf(os.Stderr, "could not read file %s \n", packageJSONPath)
				return candidates
			}

			pnpmWorkspaceYAMLData := PnpmWorkspacesPackages{}
			yaml.Unmarshal(pnpmWorkspaceYAMLPathContents, &pnpmWorkspaceYAMLData)

			workspacePath = currentDirectory
			workspaceGlobs = pnpmWorkspaceYAMLData.Packages
		}

		monocdrcJSONPath := path.Join(currentDirectory, "./.monocdrc.json")
		monocdrcJSONExists, _ := fileExists(monocdrcJSONPath)

		if monocdrcJSONExists {
			file, fileOpenError := os.Open(monocdrcJSONPath)

			if fileOpenError != nil {
				fmt.Fprintf(os.Stderr, "could not open file %s \n", monocdrcJSONPath)
				return candidates
			}

			reader := bufio.NewReader(file)
			monocdrcJSONPathContents, fileReadError := io.ReadAll(reader)
			file.Close()

			if fileReadError != nil {
				fmt.Fprintf(os.Stderr, "could not read file %s \n", monocdrcJSONPath)
				return candidates
			}

			monocdrcJSONData := MonocdrcInclude{}
			json.Unmarshal(monocdrcJSONPathContents, &monocdrcJSONData)

			workspacePath = currentDirectory
			workspaceGlobs = append(workspaceGlobs, monocdrcJSONData.Include...)
			break
		}

		nextDirectoryPath := path.Join(currentDirectory, "../")

		if nextDirectoryPath == currentDirectory {
			break
		}

		currentDirectory = nextDirectoryPath
	}

	negativeGlobs := make([]glob.Glob, 0, len(workspaceGlobs))
	positiveGlobs := make([]glob.Glob, 0, len(workspaceGlobs))

	excludeGlobs := make([]glob.Glob, 0, len(workspaceGlobs))

	excludeGlobs = append(
		excludeGlobs,
		glob.MustCompile("node_modules", '/'),
		glob.MustCompile("**/node_modules", '/'),
		glob.MustCompile(".turbo", '/'),
		glob.MustCompile("**/.turbo", '/'),
		glob.MustCompile(".next", '/'),
		glob.MustCompile("**/.next", '/'),
		glob.MustCompile(".vercel", '/'),
		glob.MustCompile("**/.vercel", '/'),
	)

	for _, g := range workspaceGlobs {
		firstCharacter := g[0]

		if firstCharacter == '!' {
			negativeGlobs = append(negativeGlobs, glob.MustCompile(g[1:], '/'))
		} else {
			positiveGlobs = append(positiveGlobs, glob.MustCompile(g, '/'))
		}
	}

	candidates = append(candidates, Candidate{
		name: ".",
		path: workingDir,
	})

	if workingDir != workspacePath {
		candidates = append(candidates, Candidate{
			name: "/",
			path: workspacePath,
		})
	}

	candidates = Walk(workspacePath, workspacePath, excludeGlobs, negativeGlobs, positiveGlobs, candidates)

	return candidates
}

func getFilteredCandidates(candidates *[]Candidate, searchText string) *[]FilteredCandidate {
	filteredCandidates := make([]FilteredCandidate, 0, 32)

	if searchText == "" {
		for _, c := range *candidates {
			filteredCandidates = append(filteredCandidates, FilteredCandidate{
				candidate: c,
				rank:      1,
			})
		}
	} else {
		lowerSearchText := strings.ToLower(searchText)

		for _, c := range *candidates {
			rank := fuzzy.RankMatch(
				lowerSearchText,
				strings.ToLower(c.name),
			)

			if rank >= 0 {
				filteredCandidates = append(filteredCandidates, FilteredCandidate{
					candidate: c,
					rank:      rank,
				})
			}
		}

		sort.Slice(filteredCandidates, func(i, j int) bool {
			return filteredCandidates[i].rank > filteredCandidates[j].rank
		})
	}

	return &filteredCandidates
}
