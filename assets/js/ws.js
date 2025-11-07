
export function connectWebSocket(onMessage) {
  const socket = new WebSocket("ws://localhost:8080/listen")

  socket.onopen = () => {
    console.log("✅ Connected to TinyWorlds WebSocket");
  };

  socket.onmessage = (event) => {
    const snapshot = JSON.parse(event.data);
    onMessage(snapshot);
  };

  socket.onclose = (err) => {
    console.log(err)
    console.log("❌ Disconnected from WebSocket");


    // Timeout Increase Exponential
    setTimeout(() => connectWebSocket(onMessage), 1000);
  };
}

