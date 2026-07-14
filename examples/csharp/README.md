# C# example

Requires .NET 6+.

```bash
export CHECKNUMBER_API_KEY="YOUR_API_KEY"
dotnet run   # reads numbers.txt (one E.164 number per line)
```

Submits to `POST /v1/tasks` (`task_type=ws`), polls `POST /v1/gettasks`, downloads results. Full docs: https://docs.checknumber.ai/whatsapp-bulk-checker
