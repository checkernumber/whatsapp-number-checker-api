// WhatsApp Number Checker API — bulk verification example (Java 11+).
// Workflow: submit a file of E.164 numbers -> poll status -> download results.
// Docs: https://docs.checknumber.ai/whatsapp-bulk-checker
import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.file.Files;
import java.nio.file.Path;

public class WhatsAppChecker {
    static final String BASE_URL = "https://api.checknumber.ai";
    static final String TASK_TYPE = "ws"; // ws | ws_active | ws_avatar
    static final String API_KEY =
        System.getenv().getOrDefault("CHECKNUMBER_API_KEY", "YOUR_API_KEY");
    static final HttpClient CLIENT = HttpClient.newHttpClient();

    // Minimal multipart builder (kept dependency-free for the example).
    static HttpRequest.BodyPublisher multipart(String boundary, byte[] fileBytes,
                                               String fileName, String taskType) throws IOException {
        var sb = new StringBuilder();
        sb.append("--").append(boundary).append("\r\n")
          .append("Content-Disposition: form-data; name=\"file\"; filename=\"").append(fileName).append("\"\r\n")
          .append("Content-Type: text/plain\r\n\r\n");
        var head = sb.toString().getBytes();
        var tail = ("\r\n--" + boundary + "\r\n"
                + "Content-Disposition: form-data; name=\"task_type\"\r\n\r\n"
                + taskType + "\r\n--" + boundary + "--\r\n").getBytes();
        var out = new java.io.ByteArrayOutputStream();
        out.write(head); out.write(fileBytes); out.write(tail);
        return HttpRequest.BodyPublishers.ofByteArray(out.toByteArray());
    }

    static String field(String json, String key) {
        int i = json.indexOf("\"" + key + "\"");
        if (i < 0) return null;
        int c = json.indexOf(':', i) + 1;
        while (json.charAt(c) == ' ' || json.charAt(c) == '"') c++;
        int e = c;
        while (e < json.length() && json.charAt(e) != '"' && json.charAt(e) != ',' && json.charAt(e) != '}') e++;
        return json.substring(c, e);
    }

    public static void main(String[] args) throws Exception {
        String boundary = "----checknumber" + System.currentTimeMillis();
        byte[] fileBytes = Files.readAllBytes(Path.of("numbers.txt"));
        var submitReq = HttpRequest.newBuilder(URI.create(BASE_URL + "/v1/tasks"))
            .header("X-API-Key", API_KEY)
            .header("Content-Type", "multipart/form-data; boundary=" + boundary)
            .POST(multipart(boundary, fileBytes, "numbers.txt", TASK_TYPE))
            .build();
        String submit = CLIENT.send(submitReq, HttpResponse.BodyHandlers.ofString()).body();
        String taskId = field(submit, "task_id");
        System.out.println("task_id: " + taskId);

        String status;
        String task;
        do {
            var getReq = HttpRequest.newBuilder(URI.create(BASE_URL + "/v1/gettasks"))
                .header("X-API-Key", API_KEY)
                .header("Content-Type", "application/x-www-form-urlencoded")
                .POST(HttpRequest.BodyPublishers.ofString("task_id=" + taskId))
                .build();
            task = CLIENT.send(getReq, HttpResponse.BodyHandlers.ofString()).body();
            status = field(task, "status");
            System.out.println("status=" + status);
            if ("failed".equals(status)) throw new RuntimeException("task failed");
            if (!"exported".equals(status)) Thread.sleep(5000);
        } while (!"exported".equals(status));

        System.out.println("result_url: " + field(task, "result_url"));
    }
}
