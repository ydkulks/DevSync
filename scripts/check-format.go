package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type provider struct {
	name        string
	path        string
	frontmatter bool
	required    map[string]string
}

func main() {
	repoRoot := filepath.Join("..")

	providers := []provider{
		{
			name:        "Cursor",
			path:        filepath.Join(repoRoot, ".cursor/rules/devsync.mdc"),
			frontmatter: true,
			required: map[string]string{
				"description": "description of the rule",
				"alwaysApply": "must be 'true'",
			},
		},
		{
			name:        "Windsurf",
			path:        filepath.Join(repoRoot, ".windsurf/rules/devsync.md"),
			frontmatter: false,
		},
		{
			name:        "Cline",
			path:        filepath.Join(repoRoot, ".clinerules/devsync.md"),
			frontmatter: false,
		},
		{
			name:        "GitHub Copilot",
			path:        filepath.Join(repoRoot, ".github/copilot-instructions.md"),
			frontmatter: false,
		},
		{
			name:        "Kiro",
			path:        filepath.Join(repoRoot, ".kiro/steering/devsync.md"),
			frontmatter: true,
			required: map[string]string{
				"title":     "title of the steering rule",
				"inclusion": "inclusion mode (e.g. 'always')",
			},
		},
		{
			name:        "Claude Code",
			path:        filepath.Join(repoRoot, ".agents/rules/devsync.md"),
			frontmatter: false,
		},
		{
			name:        "OpenCode",
			path:        filepath.Join(repoRoot, ".opencode/skills/devsync/SKILL.md"),
			frontmatter: true,
			required: map[string]string{
				"name":              "name of the skill",
				"description":       "description of the skill",
				"license":           "license identifier",
				"compatibility":     "compatible tool",
				"metadata.audience": "target audience",
				"metadata.workflow": "workflow context",
			},
		},
		{
			name:        "Zed",
			path:        filepath.Join(repoRoot, ".zed/devsync.md"),
			frontmatter: false,
		},
	}

	failed := false

	for _, p := range providers {
		fmt.Printf("[%s]\n", p.name)

		info, err := os.Stat(p.path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  LOCATION FAIL: file not found at %s\n", p.path)
			failed = true
			fmt.Println()
			continue
		}
		if info.IsDir() {
			fmt.Fprintf(os.Stderr, "  LOCATION FAIL: path is a directory: %s\n", p.path)
			failed = true
			fmt.Println()
			continue
		}
		fmt.Printf("  LOCATION OK: %s\n", p.path)

		raw, err := os.ReadFile(p.path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  FORMAT FAIL: cannot read: %v\n", err)
			failed = true
			fmt.Println()
			continue
		}
		content := string(raw)

		if p.frontmatter {
			trimmed := strings.TrimLeft(content, " \t\r\n")
			if !strings.HasPrefix(trimmed, "---") {
				fmt.Fprintf(os.Stderr, "  FORMAT FAIL: expected YAML frontmatter (starting with ---)\n")
				failed = true
				fmt.Println()
				continue
			}

			parts := strings.SplitN(content, "---", 3)
			if len(parts) < 3 {
				fmt.Fprintf(os.Stderr, "  FORMAT FAIL: unclosed frontmatter (missing closing ---)\n")
				failed = true
				fmt.Println()
				continue
			}

			frontmatterText := strings.TrimSpace(parts[1])
			fmFields := parseFrontmatterFields(frontmatterText)

			for field, desc := range p.required {
				if val, ok := fmFields[field]; !ok || val == "" {
					fmt.Fprintf(os.Stderr, "  METADATA FAIL: missing required field '%s' (%s)\n", field, desc)
					failed = true
				}
			}
			fmt.Printf("  FORMAT OK: valid YAML frontmatter\n")
		} else {
			trimmed := strings.TrimLeft(content, " \t\r\n")
			if strings.HasPrefix(trimmed, "---") {
				fmt.Fprintf(os.Stderr, "  FORMAT FAIL: unexpected frontmatter (provider does not use it)\n")
				failed = true
				fmt.Println()
				continue
			}
			fmt.Printf("  FORMAT OK: no frontmatter (plain markdown)\n")
		}

		if !p.frontmatter {
			fmt.Printf("  METADATA: none required\n")
		}
		fmt.Println()
	}

	if failed {
		fmt.Fprintln(os.Stderr, "One or more providers failed format, location, or metadata checks.")
		os.Exit(1)
	}
	fmt.Println("All providers pass format, location, and metadata checks.")
}

func parseFrontmatterFields(text string) map[string]string {
	fields := make(map[string]string)
	lines := strings.Split(text, "\n")
	var prefix string

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		rawKey := parts[0]
		value := strings.TrimSpace(parts[1])

		indented := strings.HasPrefix(rawKey, " ") || strings.HasPrefix(rawKey, "\t")
		key := strings.TrimSpace(rawKey)

		if indented && prefix != "" {
			key = prefix + "." + key
		} else {
			prefix = key
		}

		fields[key] = value
	}
	return fields
}
