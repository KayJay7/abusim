<template>
  <div>
    <ProgressBar v-if="agentsSettings.autoRefreshInterval != null" :value="countdown" :showValue="false" style="height: .5em; margin-bottom: 1em;" />
    <DataView :value="agents" :layout="layout" :paginator="true" :rows="6">
      <template #header>
        <DataViewLayoutOptions v-model="layout" />
      </template>

      <template #list="slotProps">
        <div class="p-col-12">
          <div class="agent-list-item p-grid p-ai-center vertical-container">
            <h1 class="p-col agent-list-item-title">
            {{slotProps.data.name}}
            </h1>
            <TreeTable :value="slotProps.data.memoryTree" class="p-treetable-sm treetable-very-sm p-col-4">
              <template #header>
                Memory
              </template>
              <template #empty>
                <span class="empty-text">Nothing in memory</span>
              </template>
              <Column field="name" header="Name" :expander="true"></Column>
              <Column field="value" header="Value"></Column>
            </TreeTable>
            <TreeTable :value="slotProps.data.poolTree" class="p-treetable-sm treetable-very-sm p-col-4">
              <template #header>
                Pool
              </template>
              <template #empty>
                <span class="empty-text">Nothing in pool</span>
              </template>
              <Column field="index" header="Index" :expander="true"></Column>
              <Column field="resource" header="Resource"></Column>
              <Column field="value" header="Value"></Column>
            </TreeTable>
            <div class="p-inputgroup p-col-3" style="margin-bottom: .5em;">
              <span class="p-inputgroup-addon">
                <i class="pi pi-play"></i>
              </span>
              <InputText placeholder="Input" v-model="slotProps.data.input" />
              <Button icon="pi pi-send" :disabled="slotProps.data.input == '' || ! slotProps.data.input" @click="sendInput(slotProps.data.name, slotProps.data.input)"/>
            </div>
          </div>
        </div>
      </template>

      <template #grid="slotProps">
        <div class="p-col-12 p-md-4">
          <div class="agent-grid-item p-shadow-6">
            <h1 class="agent-grid-item-title">{{slotProps.data.name}}</h1>
            <div class="p-inputgroup" style="margin-bottom: .5em;">
              <Button :icon="slotProps.data.paused ? 'pi pi-play' : 'pi pi-pause'" @click="togglePause(slotProps.data.name)"/>
              <Button icon="pi pi-step-forward" :disabled="! slotProps.data.paused" @click="stepForward(slotProps.data.name)"/>
              <InputText placeholder="Input" v-model="slotProps.data.input"/>
              <Button icon="pi pi-send" :disabled="slotProps.data.input == '' || ! slotProps.data.input" @click="sendInput(slotProps.data.name, slotProps.data.input)"/>
              <Dropdown v-model="slotProps.data.verbosity" :options="verbosityOptions" optionLabel="name" style="flex: 0 0 8em" @change="updateAgentDebugStatus(slotProps.data.name)"/>
            </div>
            <TreeTable :value="slotProps.data.memoryTree" class="p-treetable-sm treetable-very-sm" style="margin-bottom: 0.5em">
              <template #header>
                Memory
              </template>
              <template #empty>
                <span class="empty-text">Nothing in memory</span>
              </template>
              <Column field="name" header="Name" :expander="true"></Column>
              <Column field="value" header="Value"></Column>
            </TreeTable>
            <TreeTable :value="slotProps.data.poolTree" class="p-treetable-sm treetable-very-sm">
              <template #header>
                Pool
              </template>
              <template #empty>
                <span class="empty-text">Nothing in pool</span>
              </template>
              <Column field="index" header="Index" :expander="true"></Column>
              <Column field="resource" header="Resource"></Column>
              <Column field="value" header="Value"></Column>
            </TreeTable>
          </div>
        </div>
      </template>
    </DataView>
  </div>
</template>

<script>
import { ref, watch, onMounted, onUnmounted } from "vue";
import { useToast } from 'primevue/usetoast';

import { getAgentMemory, decorateAgentMemory, decorateAgentPool, postAgentInput, getAgentDebugStatus, postAgentDebugStatusChange, postAgentDebugStep } from '@/functions/coordinatorService'

export default {
  name: 'Interact',
  props: [
    'agentsList',
    'agentsSettings',
    'refresh'
  ],
  setup(props) {
    const toast = useToast()

    const agents = ref([])
    const layout = ref('grid')
    const countdown = ref(100)
    const interval = ref(null)
    const verbosityOptions = ref([
        {name: 'Fatal'},
        {name: 'Error'},
        {name: 'Warning'},
        {name: 'Info'},
        {name: 'Debug'}
    ])
    const updateRefreshInterval = (agentsSettings) => {
      stopRefreshInterval()
      if (agentsSettings.autoRefreshInterval != null) {
        interval.value = setInterval(() => {
          countdown.value -= 10
          if (countdown.value < 0) {
            refreshAgents()
            countdown.value = 100
          }
        }, agentsSettings.autoRefreshInterval / 10 * 1000);
      }
    }

    const stopRefreshInterval = () => {
      clearInterval(interval.value)
      countdown.value = 100
    }

    const loadAgents = () => {
      agents.value = []
      props.agentsList.forEach(agentName => {
        agents.value.push({
          name: agentName
        })
      })
      refreshAgents()
    }

    const refreshAgents = () => {
      agents.value.forEach((oldAgent, index, agentsValue) => {
        getAgentDebugStatus(oldAgent.name)
        .then(agent => {
          agentsValue[index].verbosity = {name: agent.status.verbosity}
          agentsValue[index].paused = agent.status.paused
        })
        .catch(error => {
          toast.add({ severity: 'error', summary: 'API Error', detail: `There has been a problem with the API operation: ${error}` })
        })
      })
      agents.value.forEach((oldAgent, index, agentsValue) => {
        getAgentMemory(oldAgent.name)
        .then(agent => {
          agentsValue[index].memory = agent.memory
          decorateAgentMemory(agentsValue[index])
          agentsValue[index].pool = agent.pool
          decorateAgentPool(agentsValue[index])
        })
        .catch(error => {
          toast.add({ severity: 'error', summary: 'API Error', detail: `There has been a problem with the API operation: ${error}` })
        })
      })
    }

    const sendInput = (agentName, input) => {
      postAgentInput(agentName, input)
      .then(() => {
        agents.value.filter(a => a.name == agentName)[0].input = ''
        if (props.agentsSettings.refreshOnInput) {
          toast.add({ severity: 'success', summary: 'Input', detail: `Input performed succesfully, now updating data`, life: 3000 })
          setTimeout(refreshAgents, 250)
        } else {
          toast.add({ severity: 'success', summary: 'Input', detail: `Input performed succesfully`, life: 3000 })
        }
      })
      .catch(error => {
        toast.add({ severity: 'error', summary: 'API Error', detail: `There has been a problem with the API operation: ${error}` })
      })
    }

    const togglePause = (agentName) => {
      let agent = agents.value.filter(a => a.name == agentName)[0]
      agent.paused = !agent.paused
      updateAgentDebugStatus(agentName)
    }

    const stepForward = (agentName) => {
      postAgentDebugStep(agentName)
      .then(() => {
         if (props.agentsSettings.refreshOnInput) {
          toast.add({ severity: 'success', summary: 'Debugger', detail: `Stepped succesfully, now updating data`, life: 3000 })
          setTimeout(refreshAgents, 250)
        } else {
          toast.add({ severity: 'success', summary: 'Debugger', detail: `Stepped succesfully`, life: 3000 })
        }
      })
      .catch(error => {
        toast.add({ severity: 'error', summary: 'API Error', detail: `There has been a problem with the API operation: ${error}` })
      })
    }

    const updateAgentDebugStatus = (agentName) => {
      let agent = agents.value.filter(a => a.name == agentName)[0]
      postAgentDebugStatusChange(agentName, agent.paused, agent.verbosity.name)
      .then(() => {
        toast.add({ severity: 'success', summary: 'Debugger', detail: `Debugger status updated succesfully`, life: 3000 })
      })
      .catch(error => {
        toast.add({ severity: 'error', summary: 'API Error', detail: `There has been a problem with the API operation: ${error}` })
      })
    }

    watch(() => props.agentsList, (current) => {
      console.log(current);
      if (current != []) {
        loadAgents()
        updateRefreshInterval(props.agentsSettings)
      } else {
        stopRefreshInterval()
      }
    })

    watch(() => props.agentsSettings, (current) => {
      updateRefreshInterval(current)
    })

    watch(() => props.refresh, () => {
      refreshAgents()
      countdown.value = 100
    })

    onMounted(() => {
      loadAgents()
      updateRefreshInterval(props.agentsSettings)
    })
    
    onUnmounted(() => {
      stopRefreshInterval()
    })
    
    return {
      agents,
      layout,
      interval,
      verbosityOptions,
      countdown,
      sendInput,
      togglePause,
      stepForward,
      updateAgentDebugStatus
    }
  }
}
</script>

<style lang="css" scoped>
::v-deep(.p-dataview-header) {
  text-align: right;
}

::v-deep(.agent-list-item) {
  background: #ffffff;
  padding: 1rem;
}
::v-deep(.agent-list-item-title) {
  text-align: center;
  font-size: 1.5em;
  margin: 0;
}

::v-deep(.agent-grid-item) {
  background: #ffffff;
  padding: 2rem;
  border-radius: 4px;
  margin: 1rem;
}

::v-deep(.agent-grid-item-title) {
  text-align: center;
  font-size: 1.5em;
  margin-top: 0;
}

::v-deep(.treetable-very-sm) {
  font-size: .85em;
}

::v-deep(.treetable-very-sm tr td) {
  padding: 0 0.5rem !important;
}

::v-deep(.empty-text) {
  margin: 0.5em;
  display: block;
  text-align: center;
}
</style>
