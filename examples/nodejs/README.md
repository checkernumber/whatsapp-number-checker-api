# Node.js example

Requires Node.js 18+ (built-in `fetch` / `FormData`). No dependencies.

From the repository root:

```bash
export CHECKNUMBER_API_KEY="YOUR_API_KEY"
cd examples
node nodejs/whatsapp_checker.js
```

Submits `numbers.txt` to `POST /v1/tasks` (`task_type=ws`), polls `POST /v1/gettasks`, and downloads the result file. Full docs: https://docs.checknumber.ai/whatsapp-bulk-checker
