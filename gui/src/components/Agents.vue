<template>
  <TabView class="tabview-custom" ref="tabview4">
    <TabPanel>
      <template #header>
        <i class="pi pi-sitemap sep-pi"></i>
        <span>Configuration</span>
      </template>
      <TabView v-if="config" >
        <TabPanel v-for="code in configSources" :key="code.filename" :header="code.filename">
          <pre v-highlightjs="code.content"><code class="yaml"></code></pre>
        </TabPanel>
      </TabView>
      <Message v-else severity="warn" :closable="false">No config loaded, please add one using the button below</Message>
    </TabPanel>
    <TabPanel>
      <template #header>
        <i class="pi pi-compass sep-pi"></i>
        <span>Explore</span>
      </template>
      <Tree v-if="config" :value="configTree"></Tree>
      <Message v-else severity="warn" :closable="false">No config loaded, please add one using the button below</Message>
    </TabPanel>
    <TabPanel>
      <template #header>
        <i class="pi pi-comments sep-pi"></i>
        <span>Interact</span>
      </template>
      <Interact v-if="config" :agents-list="agentsList" :refresh="refresh" :agents-settings="agentsSettings"/>
      <Message v-else severity="warn" :closable="false">No config loaded, please add one using the button below</Message>
    </TabPanel>
  </TabView>
</template>

<script>
import { ref, watch } from 'vue';
import { useToast } from 'primevue/usetoast';

import { configParse, getConfigTree, decorateAgentTree } from '@/functions/configParse'
import { getAgentConfig } from '@/functions/coordinatorService'

import Interact from '@/components/Interact.vue'

export default {
  name: 'Agents',
  props: [
    'configSources',
    'refresh',
    'agentsSettings'
  ],
  components: {
    Interact
  },
  setup(props) {
    const toast = useToast()

    const config = ref(false)
    const configTree = ref([])
    const agentsList = ref([])

    watch(() => props.configSources, (current) => {
      if (current == []) {
        config.value = false
        agentsList.value = []
        configTree.value = {}
        return
      }
      let configs = [], agentsLists = [], configTrees = []
      current.forEach(configFile => {
        var configDoc = configParse(configFile.content)
        if (configDoc) {
          configs.push(configDoc)
          if (configDoc.agents) {
            agentsLists.push(Object.keys(configDoc['agents']))
            let configDocTree = getConfigTree(configDoc)
            configDocTree[0].children.forEach((agentTree) => {
              getAgentConfig(agentTree.label)
              .then(agent => {
                decorateAgentTree(configDocTree, agent)
              })
              .catch(error => {
                toast.add({ severity: 'error', summary: 'API Error', detail: `There has been a problem with the API operation: ${error}` })
              })
            })
            configTrees.push(configDocTree)
          } else {
            toast.add({ severity: 'warn', summary: 'Empty config', detail: `The provided configuration (${configFile.filename}) does not contain any agent. Maybe avoid loading the file next time?` })
          }
        } else {
          toast.add({ severity: 'error', summary: 'Invalid config', detail: `The provided configuration (${configFile.filename}) is not a valid YAML file or is not semantically valid` })
        }
      })
      agentsList.value = agentsLists.flat().sort()
      configTree.value = [{
        key: "agents",
        label: "Agents",
        data: "Agents data",
        icon: "pi pi-users",
        children: configTrees.map(el => el[0].children).flat().sort((a, b) => a.label > b.label)
      },{
        key: "prototypes",
        label: "Prototypes",
        data: "Prototypes data",
        icon: "pi pi-tags",
        children: configTrees.map(el => el[1].children).flat().sort((a, b) => a.label > b.label)
      }]
      config.value = true
    })

    return {
      config,
      configTree,
      agentsList
    }
  }
}
</script>

<style scoped>
.sep-pi {
  margin-right: .5em;
}
</style>