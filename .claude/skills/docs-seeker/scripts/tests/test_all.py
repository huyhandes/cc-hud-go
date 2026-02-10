#!/usr/bin/env python3
"""Test suite for docs-seeker scripts."""

import sys
from pathlib import Path

# Add parent directory to path for imports
sys.path.insert(0, str(Path(__file__).parent.parent))

from detect_topic import detect_topic, normalize_topic, normalize_library
from analyze_llms_txt import (
    parse_urls,
    categorize_url,
    suggest_agent_distribution,
    analyze_llms_txt,
)


def test_detect_topic():
    """Test topic detection."""
    print("Testing detect_topic...")

    # Topic-specific queries
    result = detect_topic("How do I use date picker in shadcn?")
    assert result is not None, "Should detect topic-specific query"
    assert result["isTopicSpecific"] is True
    print(f"  ✓ Topic query: {result}")

    # General queries
    result = detect_topic("Documentation for Next.js")
    assert result is None, "Should return None for general query"
    print("  ✓ General query: None (as expected)")

    result = detect_topic("Temporal getting started")
    assert result is None, "Should return None for getting started query"
    print("  ✓ Getting started query: None (as expected)")

    print("  All detect_topic tests passed!\n")


def test_normalize():
    """Test normalization functions."""
    print("Testing normalization...")

    assert normalize_topic("Date Picker") == "date"
    assert normalize_topic("CACHING!") == "caching"
    print("  ✓ normalize_topic works")

    assert normalize_library("Next.js") == "next.js"
    assert normalize_library("shadcn/ui") == "shadcn/ui"
    print("  ✓ normalize_library works")

    print("  All normalization tests passed!\n")


def test_parse_urls():
    """Test URL parsing from llms.txt."""
    print("Testing parse_urls...")

    content = """# llms.txt
https://example.com/getting-started
https://example.com/api-reference
# Comment line
https://example.com/advanced-usage
"""
    urls = parse_urls(content)
    assert len(urls) == 3, f"Expected 3 URLs, got {len(urls)}"
    print(f"  ✓ Parsed {len(urls)} URLs")

    print("  All parse_urls tests passed!\n")


def test_categorize_url():
    """Test URL categorization."""
    print("Testing categorize_url...")

    assert categorize_url("https://docs.com/getting-started") == "critical"
    assert categorize_url("https://docs.com/installation") == "critical"
    assert categorize_url("https://docs.com/api-reference") == "important"
    assert categorize_url("https://docs.com/advanced-usage") == "supplementary"
    assert categorize_url("https://docs.com/random-page") == "important"
    print("  ✓ All categorizations correct")

    print("  All categorize_url tests passed!\n")


def test_suggest_agent_distribution():
    """Test agent distribution suggestions."""
    print("Testing suggest_agent_distribution...")

    result = suggest_agent_distribution(2)
    assert result["agentCount"] == 1, "2 URLs should use 1 agent"
    print(f"  ✓ 2 URLs: {result['strategy']}")

    result = suggest_agent_distribution(8)
    assert result["agentCount"] <= 5, "8 URLs should use 3-5 agents"
    print(f"  ✓ 8 URLs: {result['strategy']} ({result['agentCount']} agents)")

    result = suggest_agent_distribution(15)
    assert result["agentCount"] == 7, "15 URLs should use 7 agents"
    print(f"  ✓ 15 URLs: {result['strategy']}")

    result = suggest_agent_distribution(25)
    assert result["strategy"] == "phased", "25 URLs should use phased approach"
    print(f"  ✓ 25 URLs: {result['strategy']}")

    print("  All suggest_agent_distribution tests passed!\n")


def test_analyze_llms_txt():
    """Test full analysis."""
    print("Testing analyze_llms_txt...")

    content = """# Example llms.txt
https://docs.com/getting-started
https://docs.com/installation
https://docs.com/api-reference
https://docs.com/configuration
https://docs.com/advanced-usage
https://docs.com/migration-guide
"""
    result = analyze_llms_txt(content)
    assert result["totalUrls"] == 6
    assert len(result["grouped"]["critical"]) >= 2
    print(f"  ✓ Analyzed {result['totalUrls']} URLs")
    print(f"  ✓ Distribution: {result['distribution']['strategy']}")

    print("  All analyze_llms_txt tests passed!\n")


def main():
    """Run all tests."""
    print("=" * 50)
    print("docs-seeker Python Scripts Test Suite")
    print("=" * 50 + "\n")

    try:
        test_detect_topic()
        test_normalize()
        test_parse_urls()
        test_categorize_url()
        test_suggest_agent_distribution()
        test_analyze_llms_txt()

        print("=" * 50)
        print("All tests passed!")
        print("=" * 50)
        return 0
    except AssertionError as e:
        print(f"\n❌ Test failed: {e}")
        return 1
    except Exception as e:
        print(f"\n❌ Error: {e}")
        return 1


if __name__ == "__main__":
    sys.exit(main())
