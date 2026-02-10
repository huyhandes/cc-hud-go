#!/usr/bin/env python3
"""Topic Detection Script.

Analyzes user queries to extract library name and topic keywords.
Returns null for general queries, topic info for specific queries.
"""

import json
import re
import sys
from typing import Optional, Dict, Any

from utils.env_loader import load_env

env = load_env()
DEBUG = env.get("DEBUG", "").lower() == "true"

# Topic-specific query patterns
TOPIC_PATTERNS = [
    # "How do I use X in Y?"
    re.compile(
        r"how (?:do i|to|can i) (?:use|implement|add|setup|configure) (?:the )?(.+?) (?:in|with|for) (.+)",
        re.I,
    ),
    # "Y X strategies/patterns"
    re.compile(
        r"(.+?) (.+?) (?:strategies|patterns|techniques|methods|approaches)", re.I
    ),
    # "X Y documentation" or "Y X docs"
    re.compile(r"(.+?) (.+?) (?:documentation|docs|guide|tutorial)", re.I),
    # "Using X with Y"
    re.compile(r"using (.+?) (?:with|in|for) (.+)", re.I),
    # "Y X guide/implementation/setup"
    re.compile(r"(.+?) (.+?) (?:guide|implementation|setup|configuration)", re.I),
    # "Implement X in Y"
    re.compile(r"implement(?:ing)? (.+?) (?:in|with|for|using) (.+)", re.I),
]

# General library query patterns (non-topic specific)
GENERAL_PATTERNS = [
    re.compile(r"(?:documentation|docs) for (.+)", re.I),
    re.compile(r"(.+?) (?:getting started|quick ?start|introduction)", re.I),
    re.compile(r"(?:how to use|learn) (.+)", re.I),
    re.compile(r"(.+?) (?:api reference|overview|basics)", re.I),
]


def normalize_topic(topic: str) -> str:
    """Normalize topic keyword."""
    normalized = topic.lower().strip()
    normalized = re.sub(r"[^a-z0-9\s-]", "", normalized)
    normalized = re.sub(r"\s+", "-", normalized)
    return normalized.split("-")[0][:20]


def normalize_library(library: str) -> str:
    """Normalize library name."""
    normalized = library.lower().strip()
    normalized = re.sub(r"[^a-z0-9\s\-/.]", "", normalized)
    return re.sub(r"\s+", "-", normalized)


def detect_topic(query: str) -> Optional[Dict[str, Any]]:
    """Detect if query is topic-specific or general.

    Returns topic info dict or None for general query.
    """
    if not query or not isinstance(query, str):
        return None

    trimmed = query.strip()

    # Check general patterns first
    for pattern in GENERAL_PATTERNS:
        if pattern.search(trimmed):
            if DEBUG:
                print("[DEBUG] Matched general pattern, no topic", file=sys.stderr)
            return None

    # Check topic-specific patterns
    for i, pattern in enumerate(TOPIC_PATTERNS):
        match = pattern.search(trimmed)
        if match:
            term1, term2 = match.group(1), match.group(2)

            # Pattern 1: "Y X strategies" → term1 is library, term2 is topic
            if i == 1:
                topic = normalize_topic(term2)
                library = normalize_library(term1)
            else:
                # For other patterns, term1 is topic, term2 is library
                topic = normalize_topic(term1)
                library = normalize_library(term2)

            if DEBUG:
                print("[DEBUG] Matched topic pattern", file=sys.stderr)
                print(f"[DEBUG] Topic: {topic}", file=sys.stderr)
                print(f"[DEBUG] Library: {library}", file=sys.stderr)

            return {
                "query": trimmed,
                "topic": topic,
                "library": library,
                "isTopicSpecific": True,
            }

    if DEBUG:
        print("[DEBUG] No pattern matched, treating as general", file=sys.stderr)
    return None


def main():
    if len(sys.argv) < 2:
        print('Usage: python detect_topic.py "<user query>"', file=sys.stderr)
        sys.exit(1)

    query = " ".join(sys.argv[1:])
    result = detect_topic(query)

    if result:
        print(json.dumps(result, indent=2))
    else:
        print(json.dumps({"isTopicSpecific": False}, indent=2))


if __name__ == "__main__":
    main()
