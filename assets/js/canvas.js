import { connectWebSocket } from "./ws.js"

/**
 * @type {HTMLCanvasElement}
 */
const canvas = document.getElementById("worldCanvas")
const context = canvas.getContext("2d")

const cellSize = 16

const agents = {}

let latestSnapshot = null

connectWebSocket((snapshot) => {
  latestSnapshot = snapshot
})


function renderWorld() {
  if (!latestSnapshot) {
    requestAnimationFrame(renderWorld)
    return
  }

  let snapshot = latestSnapshot

  context.fillStyle = "#1f2937"
  context.fillRect(0, 0, canvas.width, canvas.height)

  // Render Food
  context.fillStyle = "#4ade80" // green-400
  for (let index = 0; index < snapshot.foods.length; index++) {
    const [x, y] = snapshot.foods[index];
    context.fillRect(x * cellSize, y * cellSize, cellSize, cellSize)
  }

  // Render Obstacles
  context.fillStyle = "#9ca3af" // stone-400
  for (let index = 0; index < snapshot.obstacles.length; index++) {
    const [x, y] = snapshot.obstacles[index];
    context.fillRect(x * cellSize, y * cellSize, cellSize, cellSize)
  }

  for (let index = 0; index < snapshot.agents.length; index++) {
    const { id, x, y, isDead } = snapshot.agents[index];

    if (isDead) {
      delete agents[id]
      continue
    }

    const prev = agents[id] || { x, y }

    // An Act to Move Toward Target But Slowly Instead of Teleport
    const smoothX = prev.x + (x - prev.x) * 0.25 // Linear interpolation
    const smoothY = prev.y + (y - prev.y) * 0.25 // Linear Interpolation 

    agents[id] = { x: smoothX, y: smoothY }

    context.fillStyle = "#ef4444"
    context.fillRect(smoothX * cellSize, smoothY * cellSize, cellSize, cellSize)
  }

  requestAnimationFrame(renderWorld)
}

renderWorld()
