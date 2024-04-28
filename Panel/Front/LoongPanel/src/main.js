import {createApp} from 'vue'
import App from './App.vue'
import FilesPage from './components/FilesPage.vue'
import Home from './components/HomePage.vue'
import TerminalPage from './components/TerminalPage.vue'
import {createRouter, createWebHistory} from "vue-router";


const router = createRouter({
    history: createWebHistory(),
    routes: [
        {path: '/', component: Home},
        {path: '/files', component: FilesPage},
        {path: '/terminal', component: TerminalPage},
    ],
})

createApp(App).use(router).mount('#app')
