import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import "bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";
import 'primevue/resources/primevue.min.css'  
import PrimeVue from 'primevue/config';
import { FontAwesomeIcon } from './plugins/font-awesome'
import 'primevue/resources/themes/lara-light-green/theme.css';
import 'primeicons/primeicons.css';

const app = createApp(App)
  .use(router)
  .use(store)
  .component("font-awesome-icon", FontAwesomeIcon)
  
app.use(PrimeVue);
  
app.mount("#app");

// app.use(PrimeVue, {
//     unstyled: true,
//     pt: Lara
// });