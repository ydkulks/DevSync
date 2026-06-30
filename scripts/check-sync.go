package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	repoRoot := filepath.Join("..")

	skillPath := filepath.Join(repoRoot, "SKILL.md")
	canonical := readAndStripFrontmatter(skillPath)

	type copy struct {
		path          string
		hasFrontmatter bool
	}

	copies := []copy{
		{path: filepath.Join(repoRoot, ".cursor/rules/devsync.mdc"), hasFrontmatter: true},
		{path: filepath.Join(repoRoot, ".windsurf/rules/devsync.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".clinerules/devsync.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".github/copilot-instructions.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".kiro/steering/devsync.md"), hasFrontmatter: true},
		{path: filepath.Join(repoRoot, ".agents/rules/devsync.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".zed/devsync.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".opencode/skills/devsync/SKILL.md"), hasFrontmatter: true},
	}

	failed := false
	for _, c := range copies {
		var body string
		if c.hasFrontmatter {
			body = readAndStripFrontmatter(c.path)
		} else {
			raw, err := os.ReadFile(c.path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: cannot read %s: %v\n", c.path, err)
				failed = true
				continue
			}
			body = strings.TrimSpace(string(raw))
		}

		if body != canonical {
			fmt.Fprintf(os.Stderr, "DRIFT: %s does not match SKILL.md body\n", c.path)
			failed = true
		} else {
			fmt.Printf("OK: %s\n", c.path)
		}
	}

	if failed {
		fmt.Fprintln(os.Stderr, "\nOne or more copies drifted from SKILL.md body. Update them manually.")
		os.Exit(1)
	}
	fmt.Println("\nAll copies match SKILL.md body.")
}

func readAndStripFrontmatter(path string) string {
	raw, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: cannot read %s: %v\n", path, err)
		os.Exit(1)
	}

	content := string(raw)

	if strings.HasPrefix(strings.TrimLeft(content, " \t\r\n"), "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) == 3 {
			content = parts[2]
		}
	}

	return strings.TrimSpace(content)
}
