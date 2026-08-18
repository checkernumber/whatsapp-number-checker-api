# PHP example

Requires PHP 7.4+ with the cURL extension.

From the repository root:

```bash
export CHECKNUMBER_API_KEY="YOUR_API_KEY"
cd examples
php php/whatsapp_checker.php
```

Submits to `POST /v1/tasks` (`task_type=ws`), polls `POST /v1/gettasks`, downloads results. Full docs: https://docs.checknumber.ai/whatsapp-bulk-checker
