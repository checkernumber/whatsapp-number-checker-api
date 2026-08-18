# Go example

Standard library only.

From the repository root:

```bash
export CHECKNUMBER_API_KEY="YOUR_API_KEY"
cd examples
go run go/main.go
```

Submits to `POST /v1/tasks` (`task_type=ws`), polls `POST /v1/gettasks`, and downloads the result to `results.zip`. Full docs: https://docs.checknumber.ai/whatsapp-bulk-checker
