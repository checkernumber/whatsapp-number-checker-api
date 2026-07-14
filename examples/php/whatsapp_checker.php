<?php
// WhatsApp Number Checker API — bulk verification example (PHP + cURL).
// Workflow: submit a file of E.164 numbers -> poll status -> download results.
// Docs: https://docs.checknumber.ai/whatsapp-bulk-checker

const BASE_URL = "https://api.checknumber.ai";
const TASK_TYPE = "ws"; // ws | ws_active | ws_avatar
$apiKey = getenv("CHECKNUMBER_API_KEY") ?: "YOUR_API_KEY";

function apiPost(string $url, array $fields, string $apiKey): array {
    $ch = curl_init($url);
    curl_setopt_array($ch, [
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_POST => true,
        CURLOPT_HTTPHEADER => ["X-API-Key: {$apiKey}"],
        CURLOPT_POSTFIELDS => $fields,
    ]);
    $resp = curl_exec($ch);
    $code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    if ($code !== 200) {
        throw new Exception("HTTP {$code}: {$resp}");
    }
    return json_decode($resp, true);
}

$submit = apiPost(BASE_URL . "/v1/tasks", [
    "file" => new CURLFile("numbers.txt", "text/plain", "numbers.txt"),
    "task_type" => TASK_TYPE,
], $apiKey);
$taskId = $submit["task_id"];
echo "task_id: {$taskId}\n";

do {
    $task = apiPost(BASE_URL . "/v1/gettasks", ["task_id" => $taskId], $apiKey);
    echo "status={$task['status']} success={$task['success']}/{$task['total']}\n";
    if ($task["status"] === "failed") {
        throw new Exception("task failed");
    }
    if ($task["status"] !== "exported") {
        sleep(5);
    }
} while ($task["status"] !== "exported");

if (!empty($task["result_url"])) {
    file_put_contents("results.xlsx", fopen($task["result_url"], "r"));
    echo "saved to: results.xlsx\n";
}
