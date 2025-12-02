/**
 * @type {WebSocket}
 */
export let socket = null

let pingInterval

export function connectWebSocket(onMessage) {
  console.log("Start Opening Websocket")
  socket = new WebSocket("ws://127.0.0.1:8080/listen")

  socket.onopen = () => {
    console.log("✅ Connected to TinyWorlds WebSocket");

    pingInterval = setInterval(() => {
      const imAlive = { "Type": "ping" }

      if (socket.readyState == WebSocket.OPEN) {
        socket.send(JSON.stringify(imAlive))
      }

    }, 1000 * 5)
  };

  socket.onmessage = (event) => {
    const snapshot = JSON.parse(event.data);
    onMessage(snapshot);
  };

  socket.onerror = (ev) => {
    console.log("On Erro Trigger", ev.err)
  }

  socket.onclose = (err) => {
    console.dir(err);

    console.log("❌ Disconnected from WebSocket");

    clearInterval(pingInterval)
    socket = null

    setTimeout(() => connectWebSocket(onMessage), 500);
  };
}
