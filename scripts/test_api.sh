#!/usr/bin/env bash
set -euo pipefail

HOST=${HOST:-localhost}
PORT=${PORT:-8080}
BASE="http://${HOST}:${PORT}"

json() { jq -c .; }

echo "==> Testing Song APIs"
# embedding omitted to let service generate 512-dim vector
curl -sS -X POST "$BASE/songs" -H "Content-Type: application/json" --data @<(printf %s '{"song": {"name": "song_test"}}') | json
curl -sS "$BASE/songs/song_test" | json
curl -sS "$BASE/songs" | json
curl -sS -X PUT "$BASE/songs/song_test" -H "Content-Type: application/json" --data @<(printf %s '{"song": {"name": "song_test"}}') | json

echo "==> Testing User APIs"
curl -sS -X POST "$BASE/users" -H "Content-Type: application/json" --data @<(printf %s '{"user": {"id": "user_test", "name": "Tester"}}') | json
curl -sS "$BASE/users/user_test" | json
curl -sS "$BASE/users" | json

echo "==> Testing Like/Unlike APIs"
curl -sS -X POST "$BASE/users/user_test/like/song_test" | json
curl -sS "$BASE/users/user_test/liked_songs" | json
curl -sS -X DELETE "$BASE/users/user_test/unlike/song_test" | json
curl -sS "$BASE/users/user_test/liked_songs" | json

echo "==> Cleanup"
curl -sS -X DELETE "$BASE/users/user_test" | json
curl -sS -X DELETE "$BASE/songs/song_test" | json

echo "All API checks completed."
