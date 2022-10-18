import { createApp } from 'vue'
// import App from './App.vue'
import SideApp from './SideApp.vue'

/* import the fontawesome core */
import { library } from '@fortawesome/fontawesome-svg-core'
/* import font awesome icon component */
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
/* import specific icons */
import { faPlay, faPause, faForwardStep, faRepeat, faUpload, faFileArrowDown } from '@fortawesome/free-solid-svg-icons'

import './assets/main.css'

/* add icons to the library */
library.add(
    faPlay, faPause, faForwardStep, faRepeat, faUpload, faFileArrowDown);

const app = createApp(SideApp);

// console.log(app.config);
// app.config.errorHandler = (err, instance, info) => {
//     console.log(err, instance, info);
// }

app.component('font-awesome-icon', FontAwesomeIcon);
app.mount('#app');
