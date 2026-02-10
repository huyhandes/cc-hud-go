#!/usr/bin/env python3
"""Environment variable loader for docs-seeker skill.

Respects order: os.environ > skill/.env > skills/.env > .claude/.env
"""

import os
from pathlib import Path
from typing import Dict


def parse_env_file(content: str) -> Dict[str, str]:
    """Parse .env file content into key-value pairs."""
    env = {}
    for line in content.splitlines():
        line = line.strip()
        if not line or line.startswith("#"):
            continue

        if "=" not in line:
            continue

        key, _, value = line.partition("=")
        key = key.strip()
        value = value.strip()

        # Validate key
        if not key or not key[0].isalpha() and key[0] != "_":
            continue

        # Remove quotes if present
        if len(value) >= 2:
            if (value.startswith('"') and value.endswith('"')) or (
                value.startswith("'") and value.endswith("'")
            ):
                value = value[1:-1]

        env[key] = value

    return env


def load_env() -> Dict[str, str]:
    """Load environment variables from .env files in priority order.

    Priority: os.environ > skill/.env > skills/.env > .claude/.env
    """
    script_dir = Path(__file__).resolve().parent
    skill_dir = script_dir.parent.parent
    skills_dir = skill_dir.parent
    claude_dir = skills_dir.parent

    env_paths = [
        claude_dir / ".env",  # Lowest priority
        skills_dir / ".env",
        skill_dir / ".env",  # Highest priority (file)
    ]

    merged_env: Dict[str, str] = {}

    # Load .env files in order (lowest to highest priority)
    for env_path in env_paths:
        if env_path.exists():
            try:
                content = env_path.read_text(encoding="utf-8")
                parsed = parse_env_file(content)
                merged_env.update(parsed)
            except Exception:
                pass  # Silently skip unreadable files

    # os.environ has highest priority
    merged_env.update(os.environ)

    return merged_env


def get_env(key: str, default: str = "") -> str:
    """Get environment variable with fallback."""
    env = load_env()
    return env.get(key, default)


if __name__ == "__main__":
    # Test: print all loaded env vars
    import json

    env = load_env()
    print(json.dumps({k: v for k, v in env.items() if not k.startswith("_")}, indent=2))
