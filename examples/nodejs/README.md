# Node.js example

Requires Node.js 18+ (built-in `fetch` / `FormData`). No dependencies.

```bash
export CHECKNUMBER_API_KEY="YOUR_API_KEY"
node whatsapp_checker.js   # reads numbers.txt (one E.164 number per line)
```

Submits `numbers.txt` to `POST /v1/tasks` (`task_type=ws`), polls `POST /v1/gettasks`, and downloads the result file. Full docs: https://docs.checknumber.ai/whatsapp-bulk-checker
