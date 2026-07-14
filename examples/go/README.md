# Go example

Standard library only.

```bash
export CHECKNUMBER_API_KEY="YOUR_API_KEY"
go run main.go   # reads numbers.txt (one E.164 number per line)
```

Submits to `POST /v1/tasks` (`task_type=ws`), polls `POST /v1/gettasks`, returns the result URL. Full docs: https://docs.checknumber.ai/whatsapp-bulk-checker
