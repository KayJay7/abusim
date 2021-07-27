<template>
  <div id="app">
    <SpeedDial :model="commands" :radius="140" direction="up-left" type="quarter-circle" :style="{ position: 'fixed', bottom: '25px', right: '25px'}" />
    <Toast/>
  </div>
</template>

<script>
import { ref } from 'vue';
import { useToast } from "primevue/usetoast";

export default {
  setup() {
    const toast = useToast();

    const commands = ref([
        {
            label: 'Config load',
            icon: 'pi pi-upload',
            command: () => {
                toast.add({ severity: 'success', summary: 'Config load', detail: 'Config loaded', life: 3000 });
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

    return { commands }
  }
}
</script>

<style scoped>
#app {
  font-family: Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  width: 100vw;
  height: 100vh;
}
</style>