#!/usr/bin/env python3
"""llms.txt Analyzer Script.

Parses llms.txt content and categorizes URLs for optimal agent distribution.
"""

import json
import re
import sys
from typing import Dict, List, Any

from utils.env_loader import load_env

env = load_env()
DEBUG = env.get("DEBUG", "").lower() == "true"

# URL priority categories
PRIORITY_KEYWORDS = {
    "critical": [
        "getting-started",
        "quick-start",
        "quickstart",
        "introduction",
        "intro",
        "overview",
        "installation",
        "install",
        "setup",
        "basics",
        "core-concepts",
        "fundamentals",
    ],
    "supplementary": [
        "advanced",
        "internals",
        "migration",
        "migrate",
        "troubleshooting",
        "troubleshoot",
        "faq",
        "frequently-asked",
        "changelog",
        "contributing",
        "contribute",
    ],
    "important": [
        "guide",
        "tutorial",
        "example",
        "api-reference",
        "api",
        "reference",
        "configuration",
        "config",
        "routing",
        "route",
        "data-fetching",
        "authentication",
        "auth",
    ],
}


def categorize_url(url: str) -> str:
    """Categorize URL by priority."""
    url_lower = url.lower()

    for priority in ["critical", "supplementary", "important"]:
        keywords = PRIORITY_KEYWORDS[priority]
        for keyword in keywords:
            if keyword in url_lower:
                return priority

    return "important"


def parse_urls(content: str) -> List[str]:
    """Parse llms.txt content to extract URLs."""
    if not content:
        return []

    urls = []
    for line in content.splitlines():
        trimmed = line.strip()

        # Skip comments and empty lines
        if not trimmed or trimmed.startswith("#"):
            continue

        # Extract URLs
        match = re.search(r'https?://[^\s<>"]+', trimmed, re.I)
        if match:
            urls.append(match.group(0))

    return urls


def group_by_priority(urls: List[str]) -> Dict[str, List[str]]:
    """Group URLs by priority."""
    groups: Dict[str, List[str]] = {
        "critical": [],
        "important": [],
        "supplementary": [],
    }

    for url in urls:
        priority = categorize_url(url)
        groups[priority].append(url)

    return groups


def suggest_agent_distribution(url_count: int) -> Dict[str, Any]:
    """Suggest optimal agent distribution."""
    if url_count <= 3:
        return {
            "agentCount": 1,
            "strategy": "single",
            "urlsPerAgent": url_count,
            "description": "Single agent can handle all URLs",
        }
    elif url_count <= 10:
        agents = min((url_count + 1) // 2, 5)
        return {
            "agentCount": agents,
            "strategy": "parallel",
            "urlsPerAgent": (url_count + agents - 1) // agents,
            "description": f"Deploy {agents} agents in parallel",
        }
    elif url_count <= 20:
        return {
            "agentCount": 7,
            "strategy": "parallel",
            "urlsPerAgent": (url_count + 6) // 7,
            "description": "Deploy 7 agents with balanced workload",
        }
    else:
        return {
            "agentCount": 7,
            "strategy": "phased",
            "urlsPerAgent": (url_count + 6) // 7,
            "phases": 2,
            "description": "Use two-phase approach: critical first, then important",
        }


def analyze_llms_txt(content: str) -> Dict[str, Any]:
    """Analyze llms.txt content."""
    urls = parse_urls(content)
    grouped = group_by_priority(urls)
    distribution = suggest_agent_distribution(len(urls))

    return {
        "totalUrls": len(urls),
        "urls": urls,
        "grouped": grouped,
        "distribution": distribution,
        "summary": {
            "critical": len(grouped["critical"]),
            "important": len(grouped["important"]),
            "supplementary": len(grouped["supplementary"]),
        },
    }


def main():
    if len(sys.argv) < 2:
        print(
            "Usage: python analyze_llms_txt.py <content-file-or-stdin>", file=sys.stderr
        )
        print(
            "Or pipe content: cat llms.txt | python analyze_llms_txt.py -",
            file=sys.stderr,
        )
        sys.exit(1)

    if sys.argv[1] == "-":
        content = sys.stdin.read()
    else:
        file_path = sys.argv[1]
        try:
            with open(file_path, "r", encoding="utf-8") as f:
                content = f.read()
        except FileNotFoundError:
            print(f"Error: File not found: {file_path}", file=sys.stderr)
            sys.exit(1)

    result = analyze_llms_txt(content)
    print(json.dumps(result, indent=2))


if __name__ == "__main__":
    main()
