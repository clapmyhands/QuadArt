import { createApp } from 'vue'
// import App from './App.vue'
import App from './App.vue'

const app = createApp(App);


/* import the fontawesome core */
import { library } from '@fortawesome/fontawesome-svg-core'
/* import font awesome icon component */
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
/* import specific icons */
import { faPlay, faPause, faForwardStep, faRepeat, faUpload, faFileArrowDown } from '@fortawesome/free-solid-svg-icons'

/* add icons to the library */
library.add(
    faPlay, faPause, faForwardStep, faRepeat, faUpload, faFileArrowDown);

app.component('font-awesome-icon', FontAwesomeIcon);

import './assets/main.css'


// console.log(app.config);
// app.config.errorHandler = (err, instance, info) => {
//     console.log(err, instance, info);
// }

app.mount('#app');
