import { nextSnapshot, prevSnapshot, lastUpdate, timeInterval } from "./client.js"

/**
 * @type {HTMLCanvasElement}
 */
const canvas = document.getElementById("worldCanvas")
const context = canvas.getContext("2d")

const cellSize = 16

function renderLoop() {
  if (!prevSnapshot || !nextSnapshot) return requestAnimationFrame(renderLoop)

  const elapsed = performance.now() - lastUpdate
  let interpolation = Math.min(elapsed / timeInterval, 1)

  renderWorld(prevSnapshot, nextSnapshot, interpolation)
  requestAnimationFrame(renderLoop)
}

requestAnimationFrame(renderLoop)

export function renderWorld(prevSnapshot, nextSnapshot, interpolation) {
  context.fillStyle = "#1f2937"
  context.fillRect(0, 0, canvas.width, canvas.height)

  const snapshot = nextSnapshot

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
    const agent = snapshot.agents[index];

    let prevAgent = prevSnapshot.agents.find((item) => item.id === agent.id)

    let startX = prevAgent ? prevAgent.x : agent.x
    let startY = prevAgent ? prevAgent.y : agent.y

    // An Act to Move Toward Target But Slowly Instead of Teleport
    const smoothX = startX + (agent.x - startX) * interpolation // Linear interpolation
    const smoothY = startY + (agent.y - startY) * interpolation // Linear Interpolation 

    const targetOpacity = agent.isDead ? 0 : 1;

    const startOpacityValue = prevAgent ? prevAgent.opacity : targetOpacity

    const smoothOpacity = startOpacityValue + (targetOpacity - startOpacityValue) * interpolation

    context.globalAlpha = smoothOpacity
    context.fillStyle = "#ef4444"
    context.fillRect(smoothX * cellSize, smoothY * cellSize, cellSize, cellSize)
    context.globalAlpha = 1
  }
}

