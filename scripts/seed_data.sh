#!/usr/bin/env bash
set -euo pipefail

HOST=${HOST:-localhost}
PORT=${PORT:-8080}
BASE="http://${HOST}:${PORT}"

json() { jq -c .; }

echo "==> Seeding songs"
declare -a SONGS=(
  '{"song": {"name": "yesterday"}}'
  '{"song": {"name": "hey_jude"}}'
  '{"song": {"name": "imagine"}}'
  '{"song": {"name": "bohemian_rhapsody"}}'
  '{"song": {"name": "hotel_california"}}'
  '{"song": {"name": "stairway_to_heaven"}}'
  '{"song": {"name": "smells_like_teen_spirit"}}'
  '{"song": {"name": "billie_jean"}}'
  '{"song": {"name": "like_a_rolling_stone"}}'
  '{"song": {"name": "purple_rain"}}'
  '{"song": {"name": "what_a_wonderful_world"}}'
  '{"song": {"name": "sweet_child_o_mine"}}'
  '{"song": {"name": "thriller"}}'
  '{"song": {"name": "wonderwall"}}'
  '{"song": {"name": "hallelujah"}}'
  '{"song": {"name": "losing_my_religion"}}'
  '{"song": {"name": "blackbird"}}'
  '{"song": {"name": "here_comes_the_sun"}}'
  '{"song": {"name": "no_woman_no_cry"}}'
  '{"song": {"name": "let_it_be"}}'
  '{"song": {"name": "another_brick_in_the_wall"}}'
  '{"song": {"name": "paint_it_black"}}'
)

for payload in "${SONGS[@]}"; do
  curl -sS -X POST "$BASE/songs" -H "Content-Type: application/json" --data @<(printf %s "$payload") | json
done

echo "==> Seeding users"
declare -a USERS=(
  '{"user": {"id": "u_alice", "name": "Alice"}}'
  '{"user": {"id": "u_bob", "name": "Bob"}}'
  '{"user": {"id": "u_carla", "name": "Carla"}}'
)

for payload in "${USERS[@]}"; do
  curl -sS -X POST "$BASE/users" -H "Content-Type: application/json" --data @<(printf %s "$payload") | json
done

echo "==> Seeding user likes"
# Like 2-3 songs for each user from the seeded catalog
like_songs() {
  local uid="$1"; shift
  for s in "$@"; do
    curl -sS -X POST "$BASE/users/${uid}/like/${s}" | json
  done
}

like_songs u_alice imagine here_comes_the_sun hallelujah
like_songs u_bob hotel_california stairway_to_heaven sweet_child_o_mine
like_songs u_carla wonderwall paint_it_black let_it_be

echo "==> Verifying creations"
curl -sS "$BASE/songs" | json
# show one song embedding length (should be ~512)
curl -sS "$BASE/songs/imagine" | jq -c '.song | {name, emb_len: (.embedding|length)}'
curl -sS "$BASE/users/u_alice/liked_songs" | json
curl -sS "$BASE/users/u_bob/liked_songs" | json
curl -sS "$BASE/users/u_carla/liked_songs" | json
curl -sS "$BASE/users" | json

echo "Seed complete."


