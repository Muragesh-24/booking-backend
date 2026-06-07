import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  stages: [
    { duration: "30s", target: 10 },  // slowly go to 10 users
    { duration: "1m", target: 10 },   // stay at 10 users
    { duration: "30s", target: 0 },   // come back to 0
  ],
  thresholds: {
    http_req_failed: ["rate<0.01"],        // less than 1% requests should fail
    http_req_duration: ["p(95)<500"],      // 95% requests should finish under 500ms
  },
};

export default function () {
  const res = http.get("https://api.kannadda.muragesh.tech/admin/");

  check(res, {
    "status is 200": (r) => r.status === 200,
    "response time < 500ms": (r) => r.timings.duration < 500,
  });

  sleep(1);
}