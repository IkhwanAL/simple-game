import { latestSnapshot } from "./client.js"

/**
 * @type {HTMLCanvasElement}
 */
const canvas = document.getElementById("worldCanvas")
const context = canvas.getContext("2d")

const cellSize = 16

const agents = {}

function renderLoop() {
  if (!latestSnapshot) return requestAnimationFrame(renderLoop)

  renderWorld(latestSnapshot)
  requestAnimationFrame(renderLoop)
}

requestAnimationFrame(renderLoop)

export function renderWorld(snapshot) {
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

    let prev = agents[id]
    // To Prevent Inherit Previous State
    if (!prev) {
      prev = { x, y, opacity: 0 }
      agents[id] = prev
    }

    // An Act to Move Toward Target But Slowly Instead of Teleport
    const smoothX = prev.x + (x - prev.x) * 0.25 // Linear interpolation
    const smoothY = prev.y + (y - prev.y) * 0.25 // Linear Interpolation 

    const targetOpacity = isDead ? 0 : 1;
    const smoothOpacity = prev.opacity + (targetOpacity - prev.opacity) * 0.1

    agents[id] = { x: smoothX, y: smoothY, opacity: smoothOpacity }

    context.globalAlpha = smoothOpacity
    context.fillStyle = "#ef4444"
    context.fillRect(smoothX * cellSize, smoothY * cellSize, cellSize, cellSize)
    context.globalAlpha = 1

    if (agents[id].opacity < 0.05 && isDead) {
      delete agents[id]
    }
  }
}

