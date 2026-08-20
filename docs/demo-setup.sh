#!/bin/sh
# Recreate the fixture tree used by docs/demo.tape. Scene 2 really deletes
# tests/test_main.py, so re-run this before every recording.
set -eu

d=/tmp/fzt-demo
rm -rf "$d"
mkdir -p "$d/docs" "$d/src/utils" "$d/tests"

printf '# fzt demo\n\nA tiny project tree for the recording.\n' > "$d/README.md"
printf '# Guide\n\nSee src/main.py for the entry point.\n' > "$d/docs/guide.md"
printf 'from utils.helper import greet\n\n\ndef main():\n    print(greet("fzt"))\n\n\nif __name__ == "__main__":\n    main()\n' > "$d/src/main.py"
printf 'def greet(name: str) -> str:\n    """Return a friendly greeting."""\n    return f"hello, {name}!"\n' > "$d/src/utils/helper.py"
printf 'DEBUG = True\nPORT = 8080\n' > "$d/src/utils/config.py"
printf 'from src.utils.helper import greet\n\n\ndef test_greet():\n    assert greet("fzt") == "hello, fzt!"\n' > "$d/tests/test_main.py"
