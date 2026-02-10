# Repository Analysis (No llms.txt)

**Use when:** llms.txt not available on context7.com or official site

**Speed:** ⚡⚡⚡ Slower (5-10min)
**Token usage:** 🔴 High
**Accuracy:** 🔍 Code-based

## When to Use

- Library not on context7.com
- No llms.txt on official site
- Need to analyze code structure
- Documentation incomplete

## Workflow

```
1. Find repository
   → WebSearch: "[library] github repository"
   → Verify: Official, active, has docs/

2. Clone repository
   → Bash: git clone [repo-url] /tmp/docs-analysis
   → Optional: checkout specific version/tag

3. Generate code map with codemap
   → Bash: cd /tmp/docs-analysis && codemap . > codemap-output.txt
   → For dependency analysis: codemap --deps . > codemap-deps.txt
   → codemap creates AI-friendly codebase overview

4. Read code map
   → Read: /tmp/docs-analysis/codemap-output.txt
   → Extract: README, docs/, examples/, API files structure

5. Analyze structure
   → Identify: Documentation sections
   → Extract: Installation, usage, API, examples
   → Note: Code patterns, best practices

6. Present findings
   → Source: Repository analysis
   → Caveat: Based on code, not official docs
   → Include: Repository health (stars, activity)
```

## Example

**Obscure library without llms.txt:**
```bash
# 1. Find
WebSearch: "MyLibrary github repository"
# Found: https://github.com/org/mylibrary

# 2. Clone
git clone https://github.com/org/mylibrary /tmp/docs-analysis

# 3. Generate code map with codemap
cd /tmp/docs-analysis
codemap .                    # Basic tree view
codemap --deps .             # Dependency flow map (optional)

# 4. Read
Read: codemap output from terminal
# Or save to file: codemap . > codemap-output.txt

# 5. Extract documentation
- README.md: Installation, overview
- docs/: Usage guides, API reference
- examples/: Code samples
- src/: Implementation patterns

# 6. Present
Source: Repository analysis (no llms.txt)
Health: 1.2K stars, active
```

## Codemap Benefits

✅ Fast codebase overview
✅ Shows file structure with sizes
✅ Dependency flow visualization (--deps)
✅ Diff mode for changed files (--diff)
✅ No external dependencies

## Codemap Options

```bash
codemap .                    # Basic tree view
codemap --skyline .          # City skyline visualization
codemap --deps /path/to/proj # Dependency flow map
codemap --diff               # Files changed vs main
codemap --diff --ref develop # Files changed vs develop
```

## Alternative

If no GitHub repo exists:
→ Deploy multiple Researcher agents
→ Gather: Official site, blog posts, tutorials, Stack Overflow
→ Note: Quality varies, cross-reference sources
