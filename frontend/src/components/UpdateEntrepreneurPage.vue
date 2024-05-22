<template>
  <div class="card flex flex-column md:flex-row gap-3">
  <InputGroup>
    <InputGroupAddon>
      <i class="pi pi-user"></i>
    </InputGroupAddon>
    <InputText v-model="fullName" placeholder="ФИО" :class="{ 'p-invalid': !fullName }" />
  </InputGroup>

  <InputGroup>
    <InputGroupAddon>
      <i class="pi pi-calendar"></i>
    </InputGroupAddon>
    <Calendar v-model="selectedDate" placeholder="Дата рождения" :class="{ 'p-invalid': !selectedDate }"/>
  </InputGroup>

  <InputGroup>
    <InputGroupAddon>
      <i class="pi pi-map"></i>
    </InputGroupAddon>
    <InputText v-model="city" placeholder="Город" :class="{ 'p-invalid': !city }"/>
  </InputGroup>

  <InputGroup>
    <InputGroupAddon>
      <i class="pi pi-venus"></i>
    </InputGroupAddon>
    <Dropdown v-model="selectedGender" placeholder="Пол" :options="genders" :class="{ 'p-invalid': !selectedGender }" />
  </InputGroup>

  <Button @click="updateEntrepreneur">Обновить</Button>
</div>
</template>

<script>
import EntrepreneurService from '../services/entrepreneur.service';
import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';
import InputText from 'primevue/inputtext';
import Dropdown from 'primevue/dropdown';
import Calendar from 'primevue/calendar';
import Button from 'primevue/button';
import UserModel from '../models/UserModel.js'
import { ref } from 'vue';

export default {
  name: 'UpdateEntrepreneurPage',
  components: {
    Button,
    InputGroup,
    InputGroupAddon,
    InputText,
    Dropdown,
    Calendar
  },
  setup() {
    const fullName = ref('');
    const city = ref('');
    const selectedGender = ref(null);
    const selectedDate = ref(null);
    const genders = ref(['мужской', 'женский']);
    const username = "";

    return {
      fullName,
      city,
      selectedGender,
      selectedDate,
      genders
    };
  },
  created() {
    this.fillInfo();
  },
  methods: {
    fillInfo() {
      EntrepreneurService.getEntrepreneurDetails(this.$route.params.id)
        .then(response => {
          const entrepreneur = response.data.data.entrepreneur
          this.fullName = entrepreneur.full_name;
          this.city = entrepreneur.city;
          this.selectedGender = entrepreneur.gender === "m" ? "мужской" : "женский";
          this.selectedDate = new Date(entrepreneur.birthday);
          this.username = entrepreneur.username;
        })
        .catch(error => {
          console.error(error)
        })
    },
    updateEntrepreneur() {
      this.selectedGender = this.selectedGender === "мужской" ? "m" : "w";
      const user = new UserModel(this.$route.params.id, this.username, this.fullName, this.city, this.selectedGender, this.selectedDate)

      EntrepreneurService.updateEntrepreneur(this.$route.params.id, user)
      .then(response => {
        console.log(response.status)
      })
      .catch(error => {
        console.error(error)
      })
    }
  }
}

</script>