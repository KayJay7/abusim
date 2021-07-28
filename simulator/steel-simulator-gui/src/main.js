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
import DataView from 'primevue/dataview';
import DataViewLayoutOptions from 'primevue/dataviewlayoutoptions';
import TreeTable from 'primevue/treetable';
import Column from 'primevue/column';
import ProgressBar from 'primevue/progressbar';

import 'primevue/resources/themes/saga-blue/theme.css';
import 'primevue/resources/primevue.min.css';
import 'primeicons/primeicons.css';

import VueHighlightJS from 'vue3-highlightjs'
import 'highlight.js/styles/github.css'

import 'primeflex/primeflex.css';

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
app.component('DataView', DataView);
app.component('DataViewLayoutOptions', DataViewLayoutOptions);
app.component('TreeTable', TreeTable);
app.component('Column', Column);
app.component('ProgressBar', ProgressBar);

app.mount('#app')
