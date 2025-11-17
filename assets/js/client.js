import { connectWebSocket } from "./ws.js"

const statsEl = document.getElementById("stats")

export let lastUpdate = 0
export let timeInterval = 500
export let prevSnapshot = null
export let nextSnapshot = null

connectWebSocket((snapshot) => {
  prevSnapshot = nextSnapshot

  nextSnapshot = snapshot

  const now = performance.now();
  lastUpdate = now

  statsEl.textContent = `
    Tick: ${snapshot.tick} |
    AvgEnergy: ${snapshot.avgEnergy} |
    Alive: ${snapshot.agents.length} |
    Born: ${snapshot.bornCount} |
    Dead: ${snapshot.deathCount}
  `
})


