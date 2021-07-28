import { createApp } from 'vue'
import App from './App.vue';

import PrimeVue from 'primevue/config';
import Toast from 'primevue/toast';
import ToastService from 'primevue/toastservice';
import SpeedDial from 'primevue/speeddial';
import TabView from 'primevue/tabview';
import TabPanel from 'primevue/tabpanel';
import Message from 'primevue/message';
import Tree from 'primevue/tree';

import 'primevue/resources/themes/saga-blue/theme.css';
import 'primevue/resources/primevue.min.css';
import 'primeicons/primeicons.css';

import VueHighlightJS from 'vue3-highlightjs'
import 'highlight.js/styles/github.css'

const app = createApp(App);

app.use(PrimeVue);
app.use(ToastService);

app.use(VueHighlightJS )

app.component('SpeedDial', SpeedDial);
app.component('Toast', Toast);
app.component('TabView', TabView);
app.component('TabPanel', TabPanel);
app.component('Message', Message);
app.component('Tree', Tree);

app.mount('#app')
