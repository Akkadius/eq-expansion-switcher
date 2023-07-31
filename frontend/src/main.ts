import {createApp} from 'vue'
import App from './App.vue'

// EQ Assets - These should be moved into the assets themselves
import '@/components/eq-ui/styles/eq-ui.css'

import 'bootstrap/dist/css/bootstrap-grid.css'

import '@/assets/global.css'

const app = createApp(App)
app.mount('#app')