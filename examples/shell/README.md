# Shell (curl) example

Requires `curl` and `jq`.

From the repository root:

```bash
export CHECKNUMBER_API_KEY="YOUR_API_KEY"
cd examples
./shell/whatsapp_checker.sh numbers.txt
```

Submits to `POST /v1/tasks` (`task_type=ws`), polls `POST /v1/gettasks`, downloads results. Full docs: https://docs.checknumber.ai/whatsapp-bulk-checker
