<template>
  <div id="app">
    <div id="app-content">
      <Agents :config-sources="configSources" :refresh="refreshClick" :agents-settings="agentsSettings"/>
    </div>
    <SpeedDial :model="commands" :radius="140" direction="up-left" type="quarter-circle" :disabled="!online" :style="{ position: 'fixed', bottom: '25px', right: '25px'}" />
    <Settings :visible="settingsVisible" @update="updateSettings" @close="closeSettings"/>
    <Toast/>
    <input id="config-file-input" type="file" multiple style="display: none;" @change="uploadConfigFiles" />
  </div>
</template>

<script>
import { ref, onMounted, nextTick } from 'vue';
import { useToast } from 'primevue/usetoast';

import { ping } from '@/functions/coordinatorService'

import Agents from '@/components/Agents.vue'
import Settings from '@/components/Settings.vue'

export default {
  components: {
    Agents,
    Settings
  },
  setup() {
    const toast = useToast()

    const configSources = ref([])
    const online = ref(false)
    const refreshClick = ref(false)
    const agentsSettings = ref({
      autoRefresh: true,
      autoRefreshInterval: 30,
      refreshOnInput: true,
    })
    const settingsVisible = ref(false)

    const checkOnline = () => {
      ping().then(res => {
        if (res != 'pong') {
          toast.add({ severity: 'error', summary: 'Coordinator offline', detail: 'The coordinator is offline, please start the simulator and reload the page' })
          online.value = false
        } else {
          online.value = true
        }
      })
    }

    const uploadConfigFiles = (evt) => {
      let promises = []
      evt.target.files.forEach(file => {
        promises.push(new Promise(resolve => {
          const reader = new FileReader()
          reader.onload = (evt) => {
            resolve({
              filename: file.name,
              content: evt.target.result,
            })
          }
          reader.readAsText(file)
        }))
      })
      Promise.all(promises).then(values => {
        configSources.value = values
        toast.add({ severity: 'success', summary: 'Config load', detail: 'Config uploaded', life: 3000 })
      })
    }

    const updateSettings = (settings) => {
      agentsSettings.value = {
        autoRefresh: settings.autoRefresh,
        autoRefreshInterval: settings.autoRefresh ? settings.autoRefreshInterval : null,
        refreshOnInput: settings.refreshOnInput,
      }
    }

    const closeSettings = () => {
      settingsVisible.value = false
    }

    const commands = ref([
      {
        label: 'Config load',
        icon: 'pi pi-upload',
        command: () => {
          document.querySelector('#config-file-input').click()
        }
      },
      {
        label: 'Refresh',
        icon: 'pi pi-refresh',
        command: () => {
          toast.add({ severity: 'info', summary: 'Refresh', detail: 'Data manually refreshed', life: 3000 })
          refreshClick.value = ! refreshClick.value
        }
      },
      {
        label: 'Config reset',
        icon: 'pi pi-trash',
        command: () => {
          configSources.value = []
          toast.add({ severity: 'error', summary: 'Config reset', detail: 'Config removed', life: 3000 })
        }
      },
      {
        label: 'Settings',
        icon: 'pi pi-cog',
        command: () => {
          settingsVisible.value = true
          // toast.add({ severity: 'info', summary: 'Settings', detail: 'Opened settings', life: 3000 })
        }
      }
    ])

    onMounted(() => {
      nextTick(() => {
        checkOnline()
      })
    })

    return {
      configSources,
      uploadConfigFiles,
      commands,
      online,
      refreshClick,
      agentsSettings,
      settingsVisible,
      updateSettings,
      closeSettings
    }
  }
}
</script>

<style scoped>
#app {
  font-family: Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  width: 100%;
}

#app-content {
  width: 100%;
  height: 100%;
  padding: 25px;
}
</style>