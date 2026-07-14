#!/usr/bin/env bash
# WhatsApp Number Checker API — bulk verification example (curl + jq).
# Workflow: submit a file of E.164 numbers -> poll status -> download results.
# Docs: https://docs.checknumber.ai/whatsapp-bulk-checker
set -euo pipefail

BASE_URL="https://api.checknumber.ai"
API_KEY="${CHECKNUMBER_API_KEY:-YOUR_API_KEY}"
TASK_TYPE="ws"            # ws | ws_active | ws_avatar
INPUT_FILE="${1:-numbers.txt}"

# 1. Submit task
submit=$(curl -sS --location "${BASE_URL}/v1/tasks" \
  --header "X-API-Key: ${API_KEY}" \
  --form "file=@\"${INPUT_FILE}\"" \
  --form "task_type=\"${TASK_TYPE}\"")
task_id=$(echo "$submit" | jq -r '.task_id')
echo "task_id: ${task_id}"

# 2. Poll status
while :; do
  task=$(curl -sS --location "${BASE_URL}/v1/gettasks" \
    --header "X-API-Key: ${API_KEY}" \
    --form "task_id=\"${task_id}\"")
  status=$(echo "$task" | jq -r '.status')
  echo "status=${status} success=$(echo "$task" | jq -r '.success')/$(echo "$task" | jq -r '.total')"
  [ "$status" = "exported" ] && break
  [ "$status" = "failed" ] && { echo "task failed"; exit 1; }
  sleep 5
done

# 3. Download results
result_url=$(echo "$task" | jq -r '.result_url')
if [ "$result_url" != "null" ] && [ -n "$result_url" ]; then
  curl -sS -L "$result_url" -o results.xlsx
  echo "saved to: results.xlsx"
fi
