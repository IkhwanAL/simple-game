import { connectWebSocket } from "./ws.js"

const statsEl = document.getElementById("stats")
let lastUpdate = 0
export let latestSnapshot = null

connectWebSocket((snapshot) => {
  latestSnapshot = snapshot

  const now = performance.now();

  if ((now - lastUpdate) > 250) {
    statsEl.textContent = `
      Tick: ${snapshot.tick} |
      AvgEnergy: ${snapshot.avgEnergy} |
      Alive: ${snapshot.agents.length} |
      Born: ${snapshot.bornCount} |
      Dead: ${snapshot.deathCount}
    `
    lastUpdate = now
  }
})


