package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type copy struct {
	path          string
	hasFrontmatter bool
}

type cmdCopy struct {
	path          string
	hasFrontmatter bool
	isTOML         bool
}

func main() {
	repoRoot := filepath.Join("..")

	failed := false

	// --- Skill copies ---
	skillCanonical := readAndStripFrontmatter(filepath.Join(repoRoot, "SKILL.md"))

	skillCopies := []copy{
		{path: filepath.Join(repoRoot, ".cursor/rules/devsync.mdc"), hasFrontmatter: true},
		{path: filepath.Join(repoRoot, ".windsurf/rules/devsync.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".clinerules/devsync.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".github/copilot-instructions.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".kiro/steering/devsync.md"), hasFrontmatter: true},
		{path: filepath.Join(repoRoot, ".agents/rules/devsync.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".agents/skills/devsync/SKILL.md"), hasFrontmatter: true},
		{path: filepath.Join(repoRoot, ".opencode/skills/devsync/SKILL.md"), hasFrontmatter: true},
	}

	fmt.Println("--- Skill copies ---")
	for _, c := range skillCopies {
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

		if body != skillCanonical {
			fmt.Fprintf(os.Stderr, "DRIFT: %s does not match SKILL.md body\n", c.path)
			failed = true
		} else {
			fmt.Printf("OK: %s\n", c.path)
		}
	}

	// --- Command copies: devsync ---
	cmdCanonical := readAndStripFrontmatter(filepath.Join(repoRoot, "commands/devsync.md"))

	cmdCopies := []cmdCopy{
		{path: filepath.Join(repoRoot, ".cursor/commands/devsync.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".windsurf/workflows/devsync.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".cline/skills/devsync/SKILL.md"), hasFrontmatter: true},
		{path: filepath.Join(repoRoot, ".claude/skills/devsync/SKILL.md"), hasFrontmatter: true},
		{path: filepath.Join(repoRoot, ".github/prompts/devsync.prompt.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".gemini/commands/devsync.toml"), hasFrontmatter: false, isTOML: true},
		{path: filepath.Join(repoRoot, ".kiro/steering/devsync-cmd.md"), hasFrontmatter: true},
		{path: filepath.Join(repoRoot, ".opencode/commands/devsync.md"), hasFrontmatter: true},
	}

	fmt.Println("\n--- Command copies: devsync ---")
	checkCmdCopies(cmdCopies, cmdCanonical, "commands/devsync.md", &failed)

	// --- Command copies: devsync-commit ---
	cmdCommitCanonical := readAndStripFrontmatter(filepath.Join(repoRoot, "commands/devsync-commit.md"))

	cmdCommitCopies := []cmdCopy{
		{path: filepath.Join(repoRoot, ".cursor/commands/devsync-commit.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".windsurf/workflows/devsync-commit.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".cline/skills/devsync-commit/SKILL.md"), hasFrontmatter: true},
		{path: filepath.Join(repoRoot, ".claude/skills/devsync-commit/SKILL.md"), hasFrontmatter: true},
		{path: filepath.Join(repoRoot, ".github/prompts/devsync-commit.prompt.md"), hasFrontmatter: false},
		{path: filepath.Join(repoRoot, ".gemini/commands/devsync-commit.toml"), hasFrontmatter: false, isTOML: true},
		{path: filepath.Join(repoRoot, ".kiro/steering/devsync-commit-cmd.md"), hasFrontmatter: true},
		{path: filepath.Join(repoRoot, ".opencode/commands/devsync-commit.md"), hasFrontmatter: true},
	}

	fmt.Println("\n--- Command copies: devsync-commit ---")
	checkCmdCopies(cmdCommitCopies, cmdCommitCanonical, "commands/devsync-commit.md", &failed)

	if failed {
		fmt.Fprintln(os.Stderr, "\nOne or more copies drifted. Update them manually.")
		os.Exit(1)
	}
	fmt.Println("\nAll copies match their canonical sources.")
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

func readTOMLPrompt(path string) string {
	raw, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: cannot read %s: %v\n", path, err)
		os.Exit(1)
	}

	content := string(raw)

	start := strings.Index(content, `prompt = """`)
	if start == -1 {
		fmt.Fprintf(os.Stderr, "ERROR: cannot find prompt field in %s\n", path)
		os.Exit(1)
	}

	start += len(`prompt = """`)
	end := strings.Index(content[start:], `"""`)
	if end == -1 {
		fmt.Fprintf(os.Stderr, "ERROR: unclosed prompt field in %s\n", path)
		os.Exit(1)
	}

	body := content[start : start+end]
	body = strings.ReplaceAll(body, "{{args}}", "$1")
	return strings.TrimSpace(body)
}

func checkCmdCopies(copies []cmdCopy, canonical string, label string, failed *bool) {
	for _, c := range copies {
		var body string
		if c.isTOML {
			body = readTOMLPrompt(c.path)
		} else if c.hasFrontmatter {
			body = readAndStripFrontmatter(c.path)
		} else {
			raw, err := os.ReadFile(c.path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: cannot read %s: %v\n", c.path, err)
				*failed = true
				continue
			}
			body = strings.TrimSpace(string(raw))
		}

		if body != canonical {
			fmt.Fprintf(os.Stderr, "DRIFT: %s does not match %s body\n", c.path, label)
			*failed = true
		} else {
			fmt.Printf("OK: %s\n", c.path)
		}
	}
}
