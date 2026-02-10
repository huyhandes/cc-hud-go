# General Library Documentation Search

**Use when:** User asks about entire library/framework

**Speed:** ⚡⚡ Moderate (30-60s)
**Token usage:** 🟡 Medium
**Accuracy:** 📚 Comprehensive

## Trigger Patterns

- "Documentation for [LIBRARY]"
- "[LIBRARY] getting started"
- "How to use [LIBRARY]"
- "[LIBRARY] API reference"

## Workflow (Script-First)

```bash
# STEP 1: Execute detect_topic.py script
python scripts/detect_topic.py "<user query>"
# Returns: {"isTopicSpecific": false} for general queries

# STEP 2: Execute fetch_docs.py script (handles URL construction)
python scripts/fetch_docs.py "<user query>"
# Script constructs context7.com URL automatically
# Script handles GitHub/website URL patterns
# Returns: llms.txt content with 5-20+ URLs

# STEP 3: Execute analyze_llms_txt.py script
cat llms.txt | python scripts/analyze_llms_txt.py -
# Groups URLs: critical, important, supplementary
# Recommends: agent distribution strategy
# Returns: {totalUrls, grouped, distribution}

# STEP 4: Deploy agents based on script recommendation
# - 1-3 URLs: Single agent or direct WebFetch
# - 4-10 URLs: Deploy 3-5 Explorer agents
# - 11+ URLs: Deploy 7 agents or phased approach

# STEP 5: Aggregate and present
# Synthesize findings: installation, concepts, API, examples
```

## Examples

**Temporal Go SDK:**
```bash
# Execute scripts (no manual URL construction)
python scripts/detect_topic.py "Documentation for Temporal Go SDK"
# {"isTopicSpecific": false}

python scripts/fetch_docs.py "Documentation for Temporal Go SDK"
# Script fetches: context7.com/temporalio/sdk-go/llms.txt
# Returns: llms.txt with 8 URLs

python scripts/analyze_llms_txt.py < llms.txt
# {totalUrls: 8, distribution: "3-agents", grouped: {...}}

# Deploy 3 Explorer agents as recommended:
# Agent 1: Getting started, installation, setup
# Agent 2: Core concepts, workflows, activities
# Agent 3: Configuration, API reference

# Aggregate and present comprehensive report
```

## Agent Distribution

**1-3 URLs:** Single agent
**4-10 URLs:** 3-5 agents (2-3 URLs each)
**11-20 URLs:** 7 agents (balanced)
**21+ URLs:** Two-phase (critical first, then important)

## Known Libraries

- Temporal: `temporalio/temporal`
- Temporal Go SDK: `temporalio/sdk-go`
- Temporal Python SDK: `temporalio/sdk-python`
- Next.js: `vercel/next.js`
- Astro: `withastro/astro`
- Remix: `remix-run/remix`
- shadcn/ui: `shadcn-ui/ui`

## Fallback

Scripts handle fallback automatically:
1. `fetch_docs.py` tries context7.com
2. If 404, script suggests WebSearch for llms.txt
3. If still unavailable: [Repository Analysis](./repo-analysis.md)
