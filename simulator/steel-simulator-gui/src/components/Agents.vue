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
      <Message v-if="config != null" severity="success" :closable="false">Explore here</Message>
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

import { configParse } from '@/functions/configParse'

export default {
  name: 'Agents',
  props: [
    'configsource'
  ],
  emits: [
    'invalid-config'
  ],
  setup(props, { emit }) {
    const configSourceCode = ref('')
    const config = ref(null)

    watch(() => props.configsource, (current) => {
      if (current == '') {
        config.value = null
        configSourceCode.value = ''
        return
      }
      var configDoc = configParse(current)
      console.log(configDoc);
      if (configDoc != null) {
        config.value = configDoc
        configSourceCode.value = current
      } else {
        emit('invalid-config')
      }
    });

    return {
      configSourceCode,
      config
    }
  }
}
</script>

<style scoped>
.sep-pi {
  margin-right: .5em;
}
</style>