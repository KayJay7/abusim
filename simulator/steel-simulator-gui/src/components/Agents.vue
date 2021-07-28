<template>
  <TabView class="tabview-custom" ref="tabview4">
    <TabPanel>
      <template #header>
        <i class="pi pi-sitemap sep-pi"></i>
        <span>Configuration</span>
      </template>
      <pre v-if="config != null" v-highlightjs="configSourceCode"><code class="yaml"></code></pre>
      <Message v-else severity="warn" :closable="false">No config loaded, please add one using the button below</Message>
    </TabPanel>
    <TabPanel>
      <template #header>
        <i class="pi pi-compass sep-pi"></i>
        <span>Explore</span>
      </template>
      <Tree v-if="config != null" :value="configTree"></Tree>
      <Message v-else severity="warn" :closable="false">No config loaded, please add one using the button below</Message>
    </TabPanel>
    <TabPanel>
      <template #header>
        <i class="pi pi-comments sep-pi"></i>
        <span>Interact</span>
      </template>
      <Interact v-if="config != null" :agents-list="agentsList" :refresh-rate="10" :refresh="refresh"/>
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
    'configSource',
    'refresh'
  ],
  components: {
    Interact
  },
  setup(props) {
    const toast = useToast()

    const config = ref(null)
    const configSourceCode = ref('')
    const configTree = ref([])
    const agentsList = ref([])

    watch(() => props.configSource, (current) => {
      if (current == '') {
        config.value = null
        configSourceCode.value = ''
        agentsList.value = []
        configTree.value = {}
        return
      }
      var configDoc = configParse(current)
      if (configDoc != null) {
        config.value = configDoc
        configSourceCode.value = current
        agentsList.value = Object.keys(configDoc['agents'])
        configTree.value = getConfigTree(configDoc)
        configTree.value[0].children.forEach((agentTree) => {
        getAgentConfig(agentTree.label)
        .then(agent => {
          decorateAgentTree(configTree.value, agent)
        })
        .catch(error => {
          toast.add({ severity: 'error', summary: 'API Error', detail: `There has been a problem with the API operation: ${error}` })
        })
      })
      } else {
        toast.add({ severity: 'error', summary: 'Invalid config', detail: 'The provided configuration is not a valid YAML file or is not semantically valid' })
      }
    });

    return {
      config,
      configSourceCode,
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