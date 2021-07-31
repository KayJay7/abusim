<template>
  <Dialog header="Settings" v-model:visible="dialogVisible" :style="{width: '50vw'}" :modal="true">
    <div class="p-fluid">
      <div class="p-field p-grid">
        <label for="settings-autorefresh" class="p-col-12 p-mb-3 p-md-3 p-mb-md-0">AutoRefresh</label>
        <div class="p-col-12 p-md-9">
          <SelectButton id="settings-autorefresh" v-model="autoRefresh" :options="onOffOptions" />
        </div>
      </div>
      <div class="p-field p-grid">
        <label for="settings-autorefresh-interval" class="p-col-12 p-mb-3 p-md-3 p-mb-md-0">AutoRefresh Interval</label>
        <div class="p-col-12 p-md-9">
          <InputNumber id="settings-autorefresh-interval" v-model="autoRefreshInterval" suffix=" seconds" />
        </div>
      </div>
      <div class="p-field p-grid">
        <label for="settings-autorefresh-oninput" class="p-col-12 p-mb-3 p-md-3 p-mb-md-0">Refresh on input</label>
        <div class="p-col-12 p-md-9">
          <SelectButton id="settings-autorefresh" v-model="refreshOnInput" :options="onOffOptions" />
        </div>
      </div>
    </div>
    <template #footer>
        <Button label="Cancel" @click="close" class="p-button-secondary"/>
        <Button label="Save" icon="pi pi-save" @click="save" autofocus />
    </template>
  </Dialog>
</template>

<script>
import { ref, watch } from 'vue';

export default {
  name: 'Settings',
  props: [
    'visible'
  ],
  emits: [
    'update',
    'close'
  ],
  setup(props, { emit }) {
    const dialogVisible = ref(false)
    const autoRefresh = ref('On')
    const refreshOnInput = ref('On')
    const onOffOptions = ref(['On', 'Off'])
    const autoRefreshInterval = ref(30)

    const close = () => {
      dialogVisible.value = false
      emit('close')
    }

    const save = () => {
      emit('update', {
        autoRefresh: autoRefresh.value == 'On',
        autoRefreshInterval: autoRefreshInterval.value,
        refreshOnInput: refreshOnInput.value == 'On',
      })
      close()
    }

    watch(() => props.visible, (current) => {
      dialogVisible.value = current
    })

    return {
      dialogVisible,
      autoRefresh,
      refreshOnInput,
      onOffOptions,
      autoRefreshInterval,
      close,
      save
    }
  }
}
</script>

<style scoped>
</style>
