const ws = new WebSocket("http://localhost:8080/ws")

ws.onmessage = async (e) => {
  // const world = JSON.parse(e.data);
  // console.log(world)

  const html = await fetch("/world-fragment", {
    method: "POST",
    body: e.data
  }).then(r => r.text())

  document.getElementById("world").outerHTML = html
};
