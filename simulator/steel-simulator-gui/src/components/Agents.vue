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
      <Message v-if="config != null" severity="success" :closable="false">Interact here</Message>
      <Message v-else severity="warn" :closable="false">No config loaded, please add one using the button below</Message>
    </TabPanel>
  </TabView>
</template>

<script>
import { ref, watch } from 'vue';

import { configParse, getConfigTree, decorateAgentTree } from '@/functions/configParse'
import { getAgentConfig } from '@/functions/coordinatorService'

export default {
  name: 'Agents',
  props: [
    'configsource'
  ],
  emits: [
    'invalid-config'
  ],
  setup(props, { emit }) {
    const config = ref(null)
    const configSourceCode = ref('')
    const configTree = ref([])

    watch(() => props.configsource, (current) => {
      if (current == '') {
        config.value = null
        configSourceCode.value = ''
        configTree.value = {}
        return
      }
      var configDoc = configParse(current)
      console.log(configDoc);
      if (configDoc != null) {
        config.value = configDoc
        configSourceCode.value = current
        configTree.value = getConfigTree(configDoc)
        configTree.value[0].children.forEach((agentTree) => {
        getAgentConfig(agentTree.label)
        .then(agent => {
          decorateAgentTree(configTree.value, agent)
        })
        .catch(error => {
          console.error('There has been a problem with your fetch operation:', error);
        })
      })
      } else {
        emit('invalid-config')
      }
    });

    return {
      config,
      configSourceCode,
      configTree
    }
  }
}
</script>

<style scoped>
.sep-pi {
  margin-right: .5em;
}
</style>