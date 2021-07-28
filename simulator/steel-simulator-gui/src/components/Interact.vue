<template>
  <div>
    <ProgressBar :value="countdown" :showValue="false" style="height: .5em; margin-bottom: 1em;" />
    <DataView :value="agents" :layout="layout" :paginator="true" :rows="9">
      <template #header>
        <DataViewLayoutOptions v-model="layout" />
      </template>

      <template #list="slotProps">
        <div class="p-col-12">
          <div class="agent-list-item p-grid p-ai-center vertical-container">
            <h1 class="p-col agent-list-item-title">
            {{slotProps.data.name}}
            </h1>
            <TreeTable :value="slotProps.data.memoryTree" class="p-treetable-sm treetable-very-sm p-col-5">
                <Column field="name" header="Name" :expander="true"></Column>
                <Column field="value" header="Value"></Column>
            </TreeTable>
            <div class="p-inputgroup p-col-6" style="margin-bottom: .5em;">
              <span class="p-inputgroup-addon">
                <i class="pi pi-play"></i>
              </span>
              <InputText placeholder="Input" v-model="slotProps.data.input"/>
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
              <span class="p-inputgroup-addon">
                <i class="pi pi-play"></i>
              </span>
              <InputText placeholder="Input" v-model="slotProps.data.input"/>
              <Button icon="pi pi-send" :disabled="slotProps.data.input == '' || ! slotProps.data.input" @click="sendInput(slotProps.data.name, slotProps.data.input)"/>
            </div>
            <TreeTable :value="slotProps.data.memoryTree" class="p-treetable-sm treetable-very-sm">
                <Column field="name" header="Name" :expander="true"></Column>
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

import { getAgentMemory, decorateAgentMemory, postAgentInput } from '@/functions/coordinatorService'

export default {
  name: 'Interact',
  props: [
    'agentsList',
    'refreshRate'
  ],
  setup(props) {
    const toast = useToast()

    const agents = ref([])
    const layout = ref('grid')
    const countdown = ref(100)
    const interval = ref(null)

    const updateRefreshInterval = (refreshRate) => {
      if (refreshRate) {
        interval.value = setInterval(() => {
          countdown.value -= 10
          if (countdown.value < 0) {
            refreshAgents()
            countdown.value = 100
          }
        }, refreshRate / 10 * 1000);
      } else {
        stopRefreshInterval()
      }
    }

    const stopRefreshInterval = () => {
      if (interval.value) {
        clearInterval(interval.value)
      }
      interval.value = null
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
        getAgentMemory(oldAgent.name)
        .then(agent => {
          agentsValue[index].memory = agent.memory
          decorateAgentMemory(agentsValue[index])
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
        toast.add({ severity: 'success', summary: 'Input', detail: `Input performed succesfully` })
      })
      .catch(error => {
        toast.add({ severity: 'error', summary: 'API Error', detail: `There has been a problem with the API operation: ${error}` })
      })
    }

    watch(() => props.agentsList, (current) => {
      console.log(current);
      if (current != []) {
        loadAgents()
        updateRefreshInterval(props.refreshRate)
      } else {
        stopRefreshInterval()
      }
    })

    watch(() => props.refreshRate, (current) => {
      updateRefreshInterval(current)
    })

    onMounted(() => {
      loadAgents()
      updateRefreshInterval(props.refreshRate)
    })
    
    onUnmounted(() => {
      stopRefreshInterval()
    })
    
    return {
      agents,
      layout,
      interval,
      countdown,
      sendInput
    }
  }
}
</script>

<style lang="css" scoped>
.p-dataview >>> .p-dataview-header {
  text-align: right;
}

.p-dataview >>> .agent-list-item {
  background: #ffffff;
  padding: 1rem;
}
.p-dataview >>> .agent-list-item-title {
  text-align: center;
  font-size: 1.5em;
  margin: 0;
}

.p-dataview >>> .agent-grid-item {
  background: #ffffff;
  padding: 2rem;
  border-radius: 4px;
  margin: 1rem;
}

.p-dataview >>> .agent-grid-item-title {
  text-align: center;
  font-size: 1.5em;
  margin-top: 0;
}

.p-dataview >>> .treetable-very-sm {
  font-size: .85em;
}

.p-dataview >>> .treetable-very-sm tr td{
  padding: 0 0.5rem !important;
}
</style>
