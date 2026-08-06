async function sign(secret, value) {
  const key = await crypto.subtle.importKey(
    "raw",
    new TextEncoder().encode(secret),
    { name: "HMAC", hash: "SHA-256" },
    false,
    ["sign"],
  );
  const result = await crypto.subtle.sign(
    "HMAC",
    key,
    new TextEncoder().encode(value),
  );
  return Array.from(new Uint8Array(result), (byte) =>
    byte.toString(16).padStart(2, "0"),
  ).join("");
}

export default {
  async email(message, env) {
    const timestamp = Math.floor(Date.now() / 1000).toString();
    const path = "/api/report";
    const signature = await sign(
      env.REPORT_HMAC_SECRET,
      `${timestamp}\nPOST\n${path}\n${message.to}`,
    );

    const response = await fetch(
      `https://mail.sunls.de${path}?to=${encodeURIComponent(message.to)}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/octet-stream",
          "X-Tmail-Timestamp": timestamp,
          "X-Tmail-Signature": signature,
        },
        body: message.raw,
      },
    );
    const result = await response.json().catch(() => null);
    if (!response.ok || result?.code !== 0) {
      throw new Error(result?.message || `Report failed: ${response.status}`);
    }
  },
};
