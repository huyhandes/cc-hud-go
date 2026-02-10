#!/usr/bin/env python3
"""Documentation Fetcher Script.

Fetches documentation from context7.com with topic support and fallback chain.
"""

import json
import re
import ssl
import sys
import urllib.request
import urllib.error
from typing import Dict, Any, List, Optional
from urllib.parse import quote

from utils.env_loader import load_env
from detect_topic import detect_topic

env = load_env()
DEBUG = env.get("DEBUG", "").lower() == "true"
API_KEY = env.get("CONTEXT7_API_KEY", "")

# Known repo mappings
KNOWN_REPOS = {
    "next.js": "vercel/next.js",
    "nextjs": "vercel/next.js",
    "remix": "remix-run/remix",
    "astro": "withastro/astro",
    "shadcn": "shadcn-ui/ui",
    "shadcn/ui": "shadcn-ui/ui",
    "better-auth": "better-auth/better-auth",
    "temporal": "temporalio/temporal",
    "temporal-go": "temporalio/sdk-go",
    "temporal-python": "temporalio/sdk-python",
}


def https_get(url: str) -> Optional[str]:
    """Make HTTPS GET request."""
    headers = {}
    if API_KEY:
        headers["Authorization"] = f"Bearer {API_KEY}"

    req = urllib.request.Request(url, headers=headers)
    ctx = ssl.create_default_context()

    try:
        with urllib.request.urlopen(req, context=ctx, timeout=30) as resp:
            if resp.status == 200:
                return resp.read().decode("utf-8")
            return None
    except urllib.error.HTTPError as e:
        if e.code == 404:
            return None
        raise
    except Exception:
        return None


def build_context7_url(library: str, topic: Optional[str] = None) -> str:
    """Construct context7.com URL."""
    if "/" in library:
        org, repo = library.split("/", 1)
        base_path = f"{org}/{repo}"
    else:
        normalized = re.sub(r"[^a-z0-9-]", "", library.lower())
        base_path = f"websites/{normalized}"

    base_url = f"https://context7.com/{base_path}/llms.txt"

    if topic:
        return f"{base_url}?topic={quote(topic)}"
    return base_url


def get_url_variations(library: str, topic: Optional[str] = None) -> List[str]:
    """Get URL variations to try for a library."""
    urls = []

    normalized = library.lower()
    repo = KNOWN_REPOS.get(normalized, library)

    # Primary: Try with topic if available
    if topic:
        urls.append(build_context7_url(repo, topic))

    # Fallback: Try without topic
    urls.append(build_context7_url(repo))

    return urls


def fetch_docs(query: str) -> Dict[str, Any]:
    """Fetch documentation from context7.com."""
    topic_info = detect_topic(query)

    if DEBUG:
        print(f"[DEBUG] Topic detection result: {topic_info}", file=sys.stderr)

    urls: List[str] = []

    if topic_info and topic_info.get("isTopicSpecific"):
        urls = get_url_variations(topic_info["library"], topic_info["topic"])
        if DEBUG:
            print(f"[DEBUG] Topic-specific URLs: {urls}", file=sys.stderr)
    else:
        # Extract library from general query
        match = re.search(r"(?:documentation|docs|guide) (?:for )?(.+)", query, re.I)
        if match:
            library = match.group(1).strip()
            urls = get_url_variations(library)
            if DEBUG:
                print(f"[DEBUG] General library URLs: {urls}", file=sys.stderr)

    # Try each URL
    for url in urls:
        if DEBUG:
            print(f"[DEBUG] Trying URL: {url}", file=sys.stderr)

        try:
            content = https_get(url)
            if content:
                return {
                    "success": True,
                    "source": "context7.com",
                    "url": url,
                    "content": content,
                    "topicSpecific": "?topic=" in url,
                }
        except Exception as e:
            if DEBUG:
                print(f"[DEBUG] Failed to fetch {url}: {e}", file=sys.stderr)

    return {
        "success": False,
        "source": "context7.com",
        "error": "Documentation not found on context7.com",
        "urls": urls,
        "suggestion": "Try repository analysis or web search",
    }


def main():
    if len(sys.argv) < 2:
        print('Usage: python fetch_docs.py "<user query>"', file=sys.stderr)
        sys.exit(1)

    query = " ".join(sys.argv[1:])

    try:
        result = fetch_docs(query)
        print(json.dumps(result, indent=2))
        sys.exit(0 if result["success"] else 1)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
