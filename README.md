# WhatsApp Number Checker API — Bulk Verification | CheckNumber

The **CheckNumber WhatsApp Number Checker API** is a family of three asynchronous, file-based products for WhatsApp registration, activity, and profile signals. Upload E.164 phone numbers, poll the task, and download the exported result. This is the official CheckNumber example repository, with integrations in Python, Node.js, Go, Java, C#, PHP, JavaScript, and Shell.

- **Product page:** https://checknumber.ai/products/whatsapp
- **Base URL:** `https://api.checknumber.ai`
- **Auth:** `X-API-Key` request header
- **Full API docs:** https://docs.checknumber.ai/whatsapp-bulk-checker
- **Get an API key:** https://platform.checknumber.ai

> [!NOTE]
> This is a bulk, file-based, asynchronous API: submit a task, poll for status, download results. It is designed for verifying large lists of numbers, not single realtime lookups in a request/response cycle.

## Table of Contents

- [What is the CheckNumber WhatsApp Number Checker API?](#what-is-the-checknumber-whatsapp-number-checker-api)
- [How do I check if a phone number is registered on WhatsApp via API?](#how-do-i-check-if-a-phone-number-is-registered-on-whatsapp-via-api)
- [How do I verify WhatsApp numbers in bulk?](#how-do-i-verify-whatsapp-numbers-in-bulk)
- [What request parameters does the API take?](#what-request-parameters-does-the-api-take)
- [What does the response look like?](#what-does-the-response-look-like)
- [How do I check when a phone number was last active on WhatsApp?](#how-do-i-check-when-a-phone-number-was-last-active-on-whatsapp)
- [How do I get WhatsApp profile signals via API?](#how-do-i-get-whatsapp-profile-signals-via-api)
- [What does the WhatsApp number checker API cost?](#what-does-the-whatsapp-number-checker-api-cost)
- [WhatsApp number checker API vs alternatives](#whatsapp-number-checker-api-vs-alternatives)
- [Code examples](#code-examples)
- [Status codes](#status-codes)
- [FAQ](#faq)
- [Support & legal](#support--legal)

## What is the CheckNumber WhatsApp Number Checker API?

CheckNumber's WhatsApp API verifies lists of phone numbers through three task types: `ws` checks registration, `ws_active` adds documented activity and Business-account signals, and `ws_avatar` returns documented profile and inferred image signals. All three use the same bulk workflow: upload a text file to `POST /v1/tasks`, poll `POST /v1/gettasks`, and download the exported result. The API is intended for authorized list verification and enrichment, not messaging, contact discovery, or access to private WhatsApp content.

## How do I check if a phone number is registered on WhatsApp via API?

Send a `POST` request to `https://api.checknumber.ai/v1/tasks` with your `X-API-Key` header, a `file` of phone numbers, and `task_type=ws`. The API returns a `task_id`. Poll `https://api.checknumber.ai/v1/gettasks` with that ID until `status` becomes `exported`, then download the file at `result_url`.

```bash
# 1. Submit a task
curl --location 'https://api.checknumber.ai/v1/tasks' \
  --header 'X-API-Key: YOUR_API_KEY' \
  --form 'file=@"./numbers.txt"' \
  --form 'task_type="ws"'
```

Each line of `numbers.txt` is one phone number in E.164 format (e.g. `+14155550100`).

## How do I verify WhatsApp numbers in bulk?

The workflow is three steps:

1. **Submit** — upload the number file (`POST /v1/tasks`, `task_type=ws`). You receive a `task_id`.
2. **Poll** — call `POST /v1/gettasks` with the `task_id`. `status` moves `pending` → `processing` → `exported`.
3. **Download** — when `status` is `exported`, download the file at `result_url`.

```bash
# 2. Poll task status
curl --location 'https://api.checknumber.ai/v1/gettasks' \
  --header 'X-API-Key: YOUR_API_KEY' \
  --form 'task_id="YOUR_TASK_ID"'
```

See [`examples/`](./examples) for complete, runnable implementations of this full submit → poll → download loop in every supported language.

## What request parameters does the API take?

**Submit task** — `POST https://api.checknumber.ai/v1/tasks`

| Parameter   | Type   | Description                                                        |
| ----------- | ------ | ----------------------------------------------------------------- |
| `file`      | file   | Text file, one phone number per line in E.164 format (`+14155550100`) |
| `task_type` | string | Set to `ws` for the WhatsApp real-time checker                    |

Current product limits:

| `task_type` | Minimum rows | Maximum rows | Product documentation |
| ----------- | -----------: | -----------: | --- |
| `ws` | 500 | 10,000,000 | [WhatsApp registration checker](https://docs.checknumber.ai/whatsapp-bulk-checker) |
| `ws_active` | 1,000 | 10,000,000 | [WhatsApp activity checker](https://docs.checknumber.ai/whatsapp-activity-checker) |
| `ws_avatar` | 500 | 10,000,000 | [WhatsApp profile checker](https://docs.checknumber.ai/whatsapp-bulk-number-checker-avatar) |

**Check status** — `POST https://api.checknumber.ai/v1/gettasks`

| Parameter | Type   | Description                          |
| --------- | ------ | ------------------------------------ |
| `task_id` | string | Task ID returned from task creation  |

**Result fields** (in the downloaded file)

| Field      | Description                                | Example        |
| ---------- | ------------------------------------------ | -------------- |
| `number`   | Submitted phone number                     | `+14155550100` |
| `activated` | Whether WhatsApp registration was detected | `yes` / `no`   |

## What does the response look like?

`POST /v1/tasks` and `POST /v1/gettasks` return the same task object:

```json
{
  "created_at": "2024-10-19T18:24:56.450567423Z",
  "updated_at": "2024-10-19T18:53:43.141760071Z",
  "task_id": "cs9viu7i61pkfs4oavvg",
  "user_id": "your-user-id",
  "status": "exported",
  "total": 20000,
  "success": 20000,
  "failure": 0,
  "result_url": "https://example-link-to-results.zip"
}
```

| Field        | Description                                                        |
| ------------ | ----------------------------------------------------------------- |
| `task_id`    | Unique task identifier                                             |
| `user_id`    | Account identifier when exposed by the deployment                  |
| `status`     | `pending` (queued) → `processing` → `exported` (results ready)     |
| `total`      | Total phone numbers processed                                     |
| `success`    | Numbers successfully identified                                   |
| `failure`    | Numbers that failed processing                                    |
| `result_url` | Download URL for results, present when `status` is `exported`      |

The three WhatsApp checks below share this exact request/response flow — only the `task_type` value and the columns in the result file change.

## How do I check when a phone number was last active on WhatsApp?

Use `task_type=ws_active`. Instead of a simple yes/no, the result reports a recent-activity window for each number — useful for scoring how "fresh" a contact list is before outreach or re-engagement.

```bash
curl --location 'https://api.checknumber.ai/v1/tasks' \
  --header 'X-API-Key: YOUR_API_KEY' \
  --form 'file=@"./numbers.txt"' \
  --form 'task_type="ws_active"'
```

Result fields:

| Field    | Description                                   | Example        |
| -------- | --------------------------------------------- | -------------- |
| `number` | Phone number in E.164 format | `+14155550100` |
| `activated` | Whether the number is registered | `yes` / `no` |
| `activetime` | Activity value returned by the service | service value |
| `activedays` | Active-day value returned by the service | service value |
| `business` | Whether the account is a Business account | `yes` / `no` |

The retired `signature` column is not returned.

The status/poll/download workflow is identical to the real-time checker above — change only `task_type` in any of the [code examples](#code-examples).

## How do I get WhatsApp profile signals via API?

Use `task_type=ws_avatar`. The result enriches each number with profile signals (such as avatar availability plus inferred age and gender), for use cases like audience segmentation and profile enrichment.

```bash
curl --location 'https://api.checknumber.ai/v1/tasks' \
  --header 'X-API-Key: YOUR_API_KEY' \
  --form 'file=@"./numbers.txt"' \
  --form 'task_type="ws_avatar"'
```

Result fields:

| Field     | Description                                | Example        |
| --------- | ------------------------------------------ | -------------- |
| `number` | Phone number in E.164 format | `+14155550100` |
| `activated` | Whether the number is registered | `yes` / `no` |
| `age` | Inferred age signal | service value |
| `avatar` | Avatar signal returned by the service | service value |
| `category` | Profile category signal | service value |
| `gender` | Inferred gender signal | service value |
| `hair_color` | Inferred hair-color signal | service value |
| `skin_color` | Inferred skin-tone signal | service value |

The catalog normalizes the export to one canonical `number` column.

## What does the WhatsApp number checker API cost?

Pricing is per phone number checked, billed against your account balance. Each `task_type` (`ws`, `ws_active`, `ws_avatar`) has its own per-number rate and minimum batch size.

See current pricing at **https://checknumber.ai/pricing**, and query your remaining balance any time with `GET https://api.checknumber.ai/v1/balance`.

## WhatsApp number checker API vs alternatives

| Capability                     | CheckNumber WhatsApp API | Generic HLR lookup | Manual / mobile app check |
| ------------------------------ | ------------------------ | ------------------ | ------------------------- |
| Confirms WhatsApp registration | ✅ Yes                    | ❌ No (carrier only) | ✅ Yes                     |
| Bulk file upload               | ✅ Yes                    | Varies             | ❌ No                      |
| Asynchronous large batches     | ✅ Yes                    | Varies             | ❌ No                      |
| Structured result file (xlsx/csv) | ✅ Yes                 | Varies             | ❌ No                      |
| Extra signals (last-active, avatar) | ✅ Optional add-ons  | ❌ No               | Partial (manual)          |

*Comparison is against generic categories of alternatives, not any specific named vendor.*

## Code examples

Each folder is a self-contained, runnable implementation of the full submit → poll → download workflow:

| Language   | Example                                        |
| ---------- | ---------------------------------------------- |
| Python     | [`examples/python`](./examples/python)         |
| Node.js    | [`examples/nodejs`](./examples/nodejs)         |
| Go         | [`examples/go`](./examples/go)                 |
| Java       | [`examples/java`](./examples/java)             |
| C#         | [`examples/csharp`](./examples/csharp)         |
| PHP        | [`examples/php`](./examples/php)               |
| JavaScript (browser) | [`examples/javascript`](./examples/javascript) |
| Shell (curl) | [`examples/shell`](./examples/shell)         |

Set your API key as an environment variable before running:

```bash
export CHECKNUMBER_API_KEY="YOUR_API_KEY"
```

## Status codes

| Status | Billing  | Description                                              |
| ------ | -------- | ------------------------------------------------------- |
| `200` | — | Task status returned successfully |
| `200` / `202` | estimated charge applied | Task accepted; current and legacy deployments may use either success status |
| `400` | no completed charge | Invalid file, task type, or row count |
| `401` | none | Missing or invalid API key |
| `402` | none | Insufficient balance |
| `404` | none | Task not found |
| `413` | none | Uploaded file is too large |
| `500` | none | Internal error; retry with backoff |

## FAQ

**Is this a realtime single-number lookup?**
No. It is a bulk, asynchronous, file-based API. Submit a list, poll for completion, download results. This is optimized for verifying large lists efficiently.

**What number format is required?**
E.164, one per line — e.g. `+14155550100`. Include the country code with a leading `+`.

**How do I also get last-active dates or profile info?**
Use `task_type=ws_active` for activity and Business signals, or `task_type=ws_avatar` for the documented profile signals. The request/response shape is identical; only `task_type`, limits, and exported fields change.

**How do I check my balance?**
`GET https://api.checknumber.ai/v1/balance` with your `X-API-Key` header.

## Support & legal

- **Official product page:** https://checknumber.ai/products/whatsapp
- **Registration API documentation:** https://docs.checknumber.ai/whatsapp-bulk-checker
- **Activity API documentation:** https://docs.checknumber.ai/whatsapp-activity-checker
- **Profile API documentation:** https://docs.checknumber.ai/whatsapp-bulk-number-checker-avatar
- **Dashboard & API keys:** https://platform.checknumber.ai
- **Current pricing:** https://checknumber.ai/pricing
- **License:** [MIT](./LICENSE) (sample code)

Use this API only to verify numbers you are authorized to process, and in compliance with applicable privacy laws (GDPR, CCPA, etc.) and WhatsApp's terms. This is an unofficial checker and is not affiliated with or endorsed by WhatsApp/Meta.

---

*Last reviewed: 2026-08-18 · Maintained by CheckNumber. Canonical product page: https://checknumber.ai/products/whatsapp*
