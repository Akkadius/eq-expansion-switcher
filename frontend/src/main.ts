import {createApp} from 'vue'
import App from './App.vue'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap/dist/css/bootstrap-grid.css'

import '@/assets/global.css'

// EQ Assets - These should be moved into the assets themselves
import '@/components/eq-ui/styles/eq-ui.css'

import 'css.gg/icons/icons.css'

const app = createApp(App)

app.mount('#app')