export function refreshToken() {
  return fetch("http://localhost:8080/token", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
  })
    .then((response) => response.json())
    .then((data) => {
      if (data.error) {
        throw new Error(data.error);
      }
      return data.token;
    });
}
