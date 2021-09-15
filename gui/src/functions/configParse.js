const yaml = require('js-yaml')

export function configParse(configSourceCode) {
  try {
    const doc = yaml.load(configSourceCode);
    if (doc['version']) {
      return doc
    } 
    return null
  } catch (e) {
    return null
  }
}

export function getConfigTree(configDoc) {
  var tree = [{
    key: "agents",
    label: "Agents",
    data: "Agents data",
    icon: "pi pi-users"
  },{
    key: "prototypes",
    label: "Prototypes",
    data: "Prototypes data",
    icon: "pi pi-tags"
  }]
  var agents = []
  for (const key of Object.keys(configDoc['agents'])) {
    agents.push({
      key: `agents-${key}`,
      label: key,
      data: `${key} data`,
      icon: "pi pi-user"
    })
  }
  tree[0].children = agents
  var prototypes = []
  for (const [key, value] of Object.entries(configDoc['prototypes'])) {
    var properties = []
    for (const [propkey, propvalue] of Object.entries(value)) {
      switch (propkey) {
        case 'memorycontroller':
          properties.push({
            key: `prototypes-${key}-memorycontroller`,
            label: 'Memory controller',
            data: `${key} memorycontroller data`,
            icon: "pi pi-clone",
            children: [{
              key: `prototypes-${key}-memorycontroller-value`,
              label: propvalue,
              data: `${key} memorycontroller data value`,
              icon: "pi pi-clone",
            }]
          })
          break;
        case 'tick':
          properties.push({
            key: `prototypes-${key}-tick`,
            label: 'Tick',
            data: `${key} tick data`,
            icon: "pi pi-clock",
            children: [{
              key: `prototypes-${key}-tick-value`,
              label: propvalue,
              data: `${key} tick data value`,
              icon: "pi pi-clock",
            }]
          })
          break;
        case 'memory':
          properties.push({
            key: `prototypes-${key}-memory`,
            label: 'Memory',
            data: `${key} memory data`,
            icon: "pi pi-flag",
            children: propvalue.map((el, i) => {
              return {
                key: `prototypes-${key}-memory-${i}`,
                label: el,
                data: `${key} memory data value`,
                icon: "pi pi-flag",
              }
            })
          })
          break;
        case 'rules':
          properties.push({
            key: `prototypes-${key}-rules`,
            label: 'Rules',
            data: `${key} rules data`,
            icon: "pi pi-list",
            children: propvalue.map((el, i) => {
              return {
                key: `prototypes-${key}-rules-${i}`,
                label: el,
                data: `${key} rules data value`,
                icon: "pi pi-list",
              }
            })
          })
          break;
        default:
          break;
      }
    }
    prototypes.push({
      key: `prototypes-${key}`,
      label: key,
      data: `${key} data`,
      icon: "pi pi-tag",
      children: properties
    })
  }
  tree[1].children = prototypes
  return tree
}

export function decorateAgentTree(configTree, agentData) {
  agentData.memorycontroller = agentData.memorycontroller ?? ''
  agentData.tick = agentData.tick ?? ''
  agentData.memory = agentData.memory ?? []
  agentData.rules = agentData.rules ?? []
  agentData.endpoints = agentData.endpoints ?? []
  var children = [
    {
      key: `agents-${agentData.name}-memorycontroller`,
      label: 'Memory controller',
      data: `${agentData.name} memorycontroller data`,
      icon: "pi pi-clone",
      children: [{
        key: `agents-${agentData.name}-memorycontroller-value`,
        label: agentData.memorycontroller,
        data: `${agentData.name} memorycontroller data value`,
        icon: "pi pi-clone",
      }]
    },
    {
      key: `agents-${agentData.name}-tick`,
      label: 'Tick',
      data: `${agentData.name} tick data`,
      icon: "pi pi-clock",
      children: [{
        key: `agents-${agentData.name}-tick-value`,
        label: agentData.tick,
        data: `${agentData.name} tick data value`,
        icon: "pi pi-clock",
      }]
    },
    {
      key: `agents-${agentData.name}-memory`,
      label: 'Memory',
      data: `${agentData.name} memory data`,
      icon: "pi pi-flag",
      children: agentData.memory.map((el, i) => {
        return {
          key: `agents-${agentData.name}-memory-${i}`,
          label: el,
          data: `${agentData.name} memory data value`,
          icon: "pi pi-flag",
        }
      })
    },
    {
      key: `agents-${agentData.name}-rules`,
      label: 'Rules',
      data: `${agentData.name} rules data`,
      icon: "pi pi-list",
      children: agentData.rules.map((el, i) => {
        return {
          key: `agents-${agentData.name}-rules-${i}`,
          label: el,
          data: `${agentData.name} rules data value`,
          icon: "pi pi-list",
        }
      })
    },
    {
      key: `agents-${agentData.name}-endpoints`,
      label: 'Endpoints',
      data: `${agentData.name} endpoints data`,
      icon: "pi pi-link",
      children: agentData.endpoints.map((el, i) => {
        return {
          key: `agents-${agentData.name}-endpoints-${i}`,
          label: el,
          data: `${agentData.name} endpoints data value`,
          icon: "pi pi-link",
        }
      })
    }
  ]

  configTree[0].children.forEach((agent, index, agentChildren) => {
    if (agent.label == agentData.name) {
      agentChildren[index].children = children
    }
  });
}