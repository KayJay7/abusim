const api = 'http://localhost:4000'

export function ping() {
  return fetch(`${api}/`)
  .then(response => {
    if (response.status == 418) {
      return 'pong'
    }
    return 'pang'
  })
  .catch(() => {
    return 'pang'
  })
}

export function getAgentConfig(agentName) {
  return fetch(`${api}/config/${agentName}`)
  .then(response => response.json().then(data => ({response, json: data})))
  .then(result => {
    if (! result.response.ok) {
      return Promise.reject(`${result.response.status} ${result.response.statusText}: ${result.json.error}`)
    }
    return Promise.resolve(result.json)
  })
}

export function getAgentMemory(agentName) {
  return fetch(`${api}/memory/${agentName}`)
  .then(response => response.json().then(data => ({response, json: data})))
  .then(result => {
    if (! result.response.ok) {
      return Promise.reject(`${result.response.status} ${result.response.statusText}: ${result.json.error}`)
    }
    return Promise.resolve(result.json)
  })
}

export function decorateAgentMemory(agentObj) {
  agentObj.memoryTree = []
  if (! agentObj.memory) {
    return
  }
  if (agentObj.memory.bool != {}) {
    var bools = []
    for (const [name, value] of Object.entries(agentObj.memory.bool)) {
      bools.push({
        key: `${agentObj.name}-memory-bool-${name}`,
        data: {
          name,
          value,
        },
      })
    }
    agentObj.memoryTree.push({
      key: `${agentObj.name}-memory-bool`,
      data: {
        name: "Bool",
        value: "",
      },
      "children": bools
    })
  }
}

export function postAgentInput(agentName, input) {
  return fetch(`${api}/memory/${agentName}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      actions: input
    })
  })
  .then(response => response.json().then(data => ({response, json: data})))
  .then(result => {
    if (! result.response.ok) {
      return Promise.reject(`${result.response.status} ${result.response.statusText}: ${result.json.error}`)
    }
    return Promise.resolve(result.json)
  })
}
