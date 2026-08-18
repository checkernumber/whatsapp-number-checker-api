# Python example

From the repository root:

```bash
python3 -m pip install requests
export CHECKNUMBER_API_KEY="YOUR_API_KEY"
cd examples
python3 python/whatsapp_checker.py
```

Submits `numbers.txt` to `POST /v1/tasks` (`task_type=ws`), polls `POST /v1/gettasks`, and downloads the result file. Full docs: https://docs.checknumber.ai/whatsapp-bulk-checker
