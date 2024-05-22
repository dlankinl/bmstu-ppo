<template>
  <div class="card flex flex-column md:flex-row gap-3">
    <InputGroup>
      <InputGroupAddon>
        <i class="pi pi-user"></i>
      </InputGroupAddon>
      <InputText placeholder="ФИО" />
    </InputGroup>

    <InputGroup>
      <InputGroupAddon>
        <i class="pi pi-calendar"></i>
      </InputGroupAddon>
      <Calendar v-model="selectedDate" placeholder="Дата рождения" />
    </InputGroup>

    <InputGroup>
      <InputGroupAddon>
        <i class="pi pi-map"></i>
      </InputGroupAddon>
      <InputText placeholder="Город" />
    </InputGroup>

    <InputGroup>
      <InputGroupAddon>
        <i class="pi pi-venus"></i>
      </InputGroupAddon>
      <Dropdown v-model="selectedGender" placeholder="Пол" :options="genders" />
    </InputGroup>

    <Button></Button>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const selectedGender = ref();
const selectedDate = ref();
const genders = ref(['мужской', 'женский']);
</script>
  
<script>
import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';
import InputText from 'primevue/inputtext';
import Calendar from 'primevue/calendar';
import Dropdown from 'primevue/dropdown';
import Button from 'primevue/button';


export default {
  name: 'UpdateUser',
  components: {
    InputGroup,
    InputGroupAddon,
    Calendar,
    Dropdown,
    InputText
  },
  // data() {
  //   return {
  //     genders: ['мужской', 'женский'],
  //   };
  // },
};

</script>
  