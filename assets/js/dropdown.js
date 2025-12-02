
function toggleDetail() {
  const div = htmx.find("#agent-detail")
  htmx.toggleClass(div, "hidden")

  const svg = document.querySelector("#chevron-toggle svg")

  if (svg) {
    svg.classList.toggle("rotate")
  }
}
