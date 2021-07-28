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

export function getAgentMemory(agentName) {
  return fetch(`${api}/memory/${agentName}`)
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.json();
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
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.json();
  })
}
