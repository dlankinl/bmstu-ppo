<template>
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
        <i class="pi pi-tag"></i>
      </InputGroupAddon>
      <InputNumber v-model="cost" :minFractionDigits="0" :maxFractionDigits="2" placeholder="Вес" :class="{ 'p-invalid': !cost }"/>
    </InputGroup>
    <Button @click="updateActivityField">Обновить</Button>
  </div>
</template>
  
<script>
import ActivityFieldsService from '../services/activity-fields.service';
import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import Button from 'primevue/button';
import ActivityFieldModel from '../models/ActivityFieldModel.js'
import Message from 'primevue/message';
import { ref } from 'vue';
  
export default {
  name: 'UpdateActivityFieldPage',
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

    return {
      name,
      description,
      cost,
      message
    };
  },
  created() {
    if (Utils.isAuth()) {
      this.role = Utils.getUserRoleJWT();
      if (this.role !== 'admin') {
        this.$router.push('/404');
      }
    }
    this.fillInfo();
  },
  methods: {
    fillInfo() {
      ActivityFieldsService.getActivityField(this.$route.params.id)
        .then(response => {
          this.activityField = response.data.data.activity_field;
          this.name = this.activityField.name;
          this.description = this.activityField.description;
          this.cost = this.activityField.cost;
        })
        .catch(error => {
          console.error(error)
        })
    },
    updateActivityField() {
      const field = new ActivityFieldModel(this.$route.params.id, this.name, this.description, this.cost);

      ActivityFieldsService.updateField(this.$route.params.id, field)
      .then(response => {
        this.message = { severity: 'success', content: 'Данные успешно обновлены' }
      })
      .catch(error => {
        this.message = { severity: 'error', content: `Произошла ошибка при обновлении данных: ${error.response.data.error}` }
        console.error(error)
      })
    }
  }
}

</script>