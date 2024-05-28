<template>
  <div v-if="role=='admin'">
    <div class="card flex flex-column md:flex-row gap-3">
      <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-user"></i>
        </InputGroupAddon>
        <InputText v-model="name" placeholder="Название" :class="{ 'p-invalid': !name }" />
      </InputGroup>

      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-book"></i>
        </InputGroupAddon>
        <InputText v-model="description" placeholder="Описание" :class="{ 'p-invalid': !description }"/>
      </InputGroup>

      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-book"></i>
        </InputGroupAddon>
        <InputNumber v-model="cost" :minFractionDigits="0" :maxFractionDigits="2" placeholder="Вес" :class="{ 'p-invalid': !cost }"/>
      </InputGroup>
      <Button @click="createActivityField">Создать</Button>
    </div>
  </div>
  <div v-else>
  </div>  
</template>
  
<script>
import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import Button from 'primevue/button';
import ActivityFieldModel from '../models/ActivityFieldModel.js'
import Message from 'primevue/message';
import Utils from '../services/auth-header';
import { ref } from 'vue';
import ActivityFieldsService from '../services/activity-fields.service';
  
export default {
  name: 'CreateActivityFieldPage',
  components: {
    Button,
    InputGroup,
    InputGroupAddon,
    InputText,
    InputNumber,
    Message
  },
  setup() {
    const name = ref('');
    const description = ref('');
    const cost = ref(null);
    const message = ref(null);
    var role = null;

    return {
      name,
      description,
      cost,
      message,
      role
    };
  },
  created() {
    if (Utils.isAuth()) {
      this.role = Utils.getUserRoleJWT();
      if (this.role !== 'admin') {
        this.$router.push('/404');
      }
    }
  },
  methods: {
    createActivityField() {
      const field = new ActivityFieldModel(0, this.name, this.description, this.cost);

      ActivityFieldsService.createField(field)
        .then(response => {
          this.message = { severity: 'success', content: 'Сфера деятельности успешно добавлена' }
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Произошла ошибка при добавлении сферы деятельности: ${error.response.data.error}` }
          console.error(error)
        })
    }
  }
}

</script>