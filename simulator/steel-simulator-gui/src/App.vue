<template>
  <div id="app">
    <HelloWorld msg="Welcome to steel-simulator-gui"/>
    <form @submit.prevent="greet">
      <InputText type="text" v-model="text"/>
      <Button type="submit" label="Submit"/>
      <h3>{{message}}</h3>
    </form>
    <Toast/>
  </div>
</template>

<script>
import HelloWorld from './components/HelloWorld.vue'

export default {
    data() {
        return {
            message: null,
            text: null
        }
    },
    methods: {
        greet() {
            this.$toast.add({severity: 'info', summary: 'Hello '  + this.text});
            this.message = 'Hello ' + this.text;
        }
    },
    components: {
      HelloWorld
    }
}
</script>

<style scoped>
#app {
  font-family: Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>