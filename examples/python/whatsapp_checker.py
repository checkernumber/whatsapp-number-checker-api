#!/usr/bin/env python3
"""WhatsApp Number Checker API — bulk verification example.

Workflow: submit a file of E.164 numbers -> poll status -> download results.
Docs: https://docs.checknumber.ai/whatsapp-bulk-checker
"""
import os
import time
import requests

BASE_URL = os.environ.get("API_BASE_URL", "https://api.checknumber.ai")
API_KEY = os.environ.get("CHECKNUMBER_API_KEY", "YOUR_API_KEY")
TASK_TYPE = "ws"  # ws | ws_active | ws_avatar


def submit_task(file_path: str) -> dict:
    with open(file_path, "rb") as f:
        resp = requests.post(
            f"{BASE_URL}/v1/tasks",
            headers={"X-API-Key": API_KEY},
            files={"file": (os.path.basename(file_path), f, "text/plain")},
            data={"task_type": TASK_TYPE},
            timeout=30,
        )
    resp.raise_for_status()
    return resp.json()


def get_task(task_id: str) -> dict:
    resp = requests.post(
        f"{BASE_URL}/v1/gettasks",
        headers={"X-API-Key": API_KEY},
        data={"task_id": task_id},
        timeout=30,
    )
    resp.raise_for_status()
    return resp.json()


def poll(task_id: str, interval: int = 5) -> dict:
    while True:
        task = get_task(task_id)
        print(f"status={task['status']} success={task['success']}/{task['total']}")
        if task["status"] == "exported":
            return task
        if task["status"] == "failed":
            raise RuntimeError("task failed")
        time.sleep(interval)


def download(result_url: str, out_path: str = "results.xlsx") -> str:
    resp = requests.get(result_url, stream=True, timeout=300)
    resp.raise_for_status()
    with open(out_path, "wb") as f:
        for chunk in resp.iter_content(chunk_size=8192):
            f.write(chunk)
    return out_path


if __name__ == "__main__":
    task = submit_task("numbers.txt")
    print("task_id:", task["task_id"])
    final = poll(task["task_id"])
    if final.get("result_url"):
        print("saved to:", download(final["result_url"]))
