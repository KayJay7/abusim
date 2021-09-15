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
  if (Object.keys(agentObj.memory.bool).length != 0) {
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
  if (Object.keys(agentObj.memory.integer).length != 0) {
    var integers = []
    for (const [name, value] of Object.entries(agentObj.memory.integer)) {
      integers.push({
        key: `${agentObj.name}-memory-integer-${name}`,
        data: {
          name,
          value,
        },
      })
    }
    agentObj.memoryTree.push({
      key: `${agentObj.name}-memory-integer`,
      data: {
        name: "Integer",
        value: "",
      },
      "children": integers
    })
  }
  if (Object.keys(agentObj.memory.float).length != 0) {
    var floats = []
    for (const [name, value] of Object.entries(agentObj.memory.float)) {
      floats.push({
        key: `${agentObj.name}-memory-float-${name}`,
        data: {
          name,
          value,
        },
      })
    }
    agentObj.memoryTree.push({
      key: `${agentObj.name}-memory-float`,
      data: {
        name: "Float",
        value: "",
      },
      "children": floats
    })
  }
  if (Object.keys(agentObj.memory.text).length != 0) {
    var texts = []
    for (const [name, value] of Object.entries(agentObj.memory.text)) {
      texts.push({
        key: `${agentObj.name}-memory-text-${name}`,
        data: {
          name,
          value,
        },
      })
    }
    agentObj.memoryTree.push({
      key: `${agentObj.name}-memory-text`,
      data: {
        name: "Text",
        value: "",
      },
      "children": texts
    })
  }
  if (Object.keys(agentObj.memory.time).length != 0) {
    var times = []
    for (const [name, value] of Object.entries(agentObj.memory.time)) {
      times.push({
        key: `${agentObj.name}-memory-time-${name}`,
        data: {
          name,
          value,
        },
      })
    }
    agentObj.memoryTree.push({
      key: `${agentObj.name}-memory-time`,
      data: {
        name: "Time",
        value: "",
      },
      "children": times
    })
  }
}

export function decorateAgentPool(agentObj) {
  agentObj.poolTree = []
  if (! agentObj.pool) {
    return
  }
  agentObj.pool.forEach((action, i) => {
    var resources = []
    action.forEach((res, j) => {
      resources.push({
        key: `${agentObj.name}-pool-${i}-${j}`,
        data: {
          index: `${i}`,
          resource: res.resource,
          value: res.value,
        },
      })
    })
    agentObj.poolTree.push({
      key: `${agentObj.name}-pool-${i}`,
      data: {
        index: `${i}`,
        resource: "",
        value: "",
      },
      "children": resources
    })
  })
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

export function getAgentDebugStatus(agentName) {
  return fetch(`${api}/debug/${agentName}`)
  .then(response => response.json().then(data => ({response, json: data})))
  .then(result => {
    if (! result.response.ok) {
      return Promise.reject(`${result.response.status} ${result.response.statusText}: ${result.json.error}`)
    }
    return Promise.resolve(result.json)
  })
}

export function postAgentDebugStatusChange(agentName, paused, verbosity) {
  return fetch(`${api}/debug/${agentName}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      paused,
      verbosity
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

export function postAgentDebugStep(agentName) {
  return fetch(`${api}/debug/${agentName}/step`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: ''
  })
  .then(response => response.json().then(data => ({response, json: data})))
  .then(result => {
    if (! result.response.ok) {
      return Promise.reject(`${result.response.status} ${result.response.statusText}: ${result.json.error}`)
    }
    return Promise.resolve(result.json)
  })
}
