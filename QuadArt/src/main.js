import { createApp } from 'vue'
// import App from './App.vue'
import SideApp from './SideApp.vue'

import './assets/main.css'

const app = createApp(SideApp);
// createApp(App).mount('#app')

// console.log(app.config);
// app.config.errorHandler = (err, instance, info) => {
//     console.log(err, instance, info);
// }
app.mount('#app');
