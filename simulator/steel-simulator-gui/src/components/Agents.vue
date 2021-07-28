<template>
  <TabView class="tabview-custom" ref="tabview4">
    <TabPanel>
      <template #header>
        <i class="pi pi-sitemap sep-pi"></i>
        <span>Configuration</span>
      </template>
      <pre v-if="configSourceCode != ''" v-highlightjs="configSourceCode"><code class="yaml"></code></pre>
      <Message v-else severity="warn" :closable="false">No config loaded, please add one using the button below</Message>
    </TabPanel>
    <TabPanel>
      <template #header>
        <i class="pi pi-compass sep-pi"></i>
        <span>Explore</span>
      </template>
      <Message v-if="configSourceCode != ''" severity="success" :closable="false">Explore here</Message>
      <Message v-else severity="warn" :closable="false">No config loaded, please add one using the button below</Message>
    </TabPanel>
    <TabPanel>
      <template #header>
        <i class="pi pi-comments sep-pi"></i>
        <span>Interact</span>
      </template>
      <Message v-if="configSourceCode != ''" severity="success" :closable="false">Interact here</Message>
      <Message v-else severity="warn" :closable="false">No config loaded, please add one using the button below</Message>
    </TabPanel>
  </TabView>
</template>

<script>
import { ref, watch } from 'vue';

export default {
  name: 'Agents',
  props: [
    'configsource'
  ],
  setup(props) {
    const configSourceCode = ref('')

    watch(() => props.configsource, (current) => {
      configSourceCode.value = current
    });

    return {
      configSourceCode
    }
  }
}
</script>

<style scoped>
.sep-pi {
  margin-right: .5em;
}
</style>