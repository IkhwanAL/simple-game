export let socket = null

export function connectWebSocket(onMessage) {
  console.log("Start Opening Websocket")
  socket = new WebSocket("ws://127.0.0.1:8080/listen")

  socket.onopen = () => {
    console.log("✅ Connected to TinyWorlds WebSocket");
  };

  socket.onmessage = (event) => {
    const snapshot = JSON.parse(event.data);
    onMessage(snapshot);
  };

  setInterval(() => {
    const imAlive = { "Type": "ping" }
    if (!socket) {
      return
    }
    socket.send(JSON.stringify(imAlive))
  }, 1000 * 5)

  socket.onclose = (err) => {
    console.dir(err);

    console.log("❌ Disconnected from WebSocket");

    socket = null

    setTimeout(() => connectWebSocket(onMessage), 500);
  };
}

