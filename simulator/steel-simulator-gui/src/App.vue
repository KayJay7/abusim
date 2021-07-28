<template>
  <div id="app">
    <div id="app-content">
      <Agents :configsource="configSourceCode" @invalid-config="showInvalidConfig" />
    </div>
    <SpeedDial :model="commands" :radius="140" direction="up-left" type="quarter-circle" :style="{ position: 'fixed', bottom: '25px', right: '25px'}" />
    <Toast/>
    <input id="config-file-input" type="file" style="display: none;" @change="uploadConfigFile" />
  </div>
</template>

<script>
import { ref } from 'vue';
import { useToast } from 'primevue/usetoast';

import Agents from '@/components/Agents.vue'

export default {
  components: {
    Agents
  },
  setup() {
    const toast = useToast();

    const configSourceCode = ref('')

    const uploadConfigFile = (evt) => {
      const reader = new FileReader()
      reader.onload = (evt) => {
        configSourceCode.value = evt.target.result
        toast.add({ severity: 'success', summary: 'Config load', detail: 'Config uploaded', life: 3000 });
      };
      reader.readAsText(evt.target.files[0])
    }

    const showInvalidConfig = () => {
      toast.add({ severity: 'error', summary: 'Invalid config', detail: 'The provided configuration is not a valid YAML file or is not semantically valid', life: 3000 });
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
                toast.add({ severity: 'info', summary: 'Refresh', detail: 'Data refreshed', life: 3000 });
            }
        },
        {
            label: 'Config reset',
            icon: 'pi pi-trash',
            command: () => {
                configSourceCode.value = ''
                toast.add({ severity: 'error', summary: 'Config reset', detail: 'Config removed', life: 3000 });
            }
        },
        {
            label: 'Settings',
            icon: 'pi pi-cog',
            command: () => {
                toast.add({ severity: 'info', summary: 'Settings', detail: 'Opened settings', life: 3000 });
            }
        }
    ]);

    return {
      configSourceCode,
      uploadConfigFile,
      showInvalidConfig,
      commands
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