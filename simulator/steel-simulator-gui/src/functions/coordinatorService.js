const api = 'http://localhost:4000'

export function getAgentConfig(agentName) {
  return fetch(`${api}/config/${agentName}`)
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.json();
  })
}