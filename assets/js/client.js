import { connectWebSocket } from "./ws.js"

const worldEl = document.getElementById("world");
const statsEl = document.getElementById("stats");

let gridWidth = 0;
let gridHeight = 0;
let cells = [];
let agentEls = {};

connectWebSocket((snapshot) => {
  renderWorldAndStat(snapshot);
});

function renderWorldAndStat(snapshot) {
  if (cells.length === 0) {
    gridWidth = snapshot.width
    gridHeight = snapshot.height
    worldEl.style.gridTemplateColumns = `repeat(${gridWidth}, 16px)`
    createGrid(gridWidth, gridHeight)
  }

  statsEl.textContent = `Tick: ${snapshot.tick} | Agents: ${snapshot.agents?.length || 0} | Food: ${snapshot.foods?.length || 0}`;

  for (let y = 0; y < gridHeight; y++) {
    for (let x = 0; x < gridWidth; x++) {
      cells[y][x].className = "cell bg-gray-800"
    }
  }

  for (let i = 0; i < snapshot.foods.length; i++) {
    const [x, y] = snapshot.foods[i];
    cells[y][x].className = "cell bg-green-400"
  }


  for (let i = 0; i < snapshot.obstacles.length; i++) {
    const [x, y] = snapshot.obstacles[i];
    cells[y][x].className = "cell bg-stone-400"
  }

  for (let i = 0; i < snapshot.agents.length; i++) {
    const agent = snapshot.agents[i];
    let agentElement = agentEls[agent.id]

    if (!agentElement) {
      let el = document.createElement("div")
      el.className = "agent w-4 h-4 bg-red-500 absolute"
      worldEl.appendChild(el)
      agentEls[agent.id] = el
      agentElement = el
    }
    agentElement.style.transform = `translate(${agent.x * 16}px, ${agent.y * 16}px)`;
  }
}

function createGrid(gridWidth, gridHeight) {
  for (let y = 0; y < gridHeight; y++) {
    let row = []
    for (let x = 0; x < gridWidth; x++) {
      const div = document.createElement("div")
      div.className = "cell bg-gray-800"
      worldEl.appendChild(div)
      row.push(div)
    }
    cells.push(row)
  }
}
