
export function connectWebSocket(onMessage) {
  console.log("Start Opening Websocket")
  const socket = new WebSocket("ws://127.0.0.1:8080/listen")

  socket.onopen = () => {
    console.log("✅ Connected to TinyWorlds WebSocket");
  };

  socket.onmessage = (event) => {
    const snapshot = JSON.parse(event.data);
    onMessage(snapshot);
  };

  socket.onclose = (err) => {
    console.dir(err);

    console.log("❌ Disconnected from WebSocket");

    setTimeout(() => connectWebSocket(onMessage), 500);
  };
}

