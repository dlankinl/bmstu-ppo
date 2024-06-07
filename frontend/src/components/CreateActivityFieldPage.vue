<!-- <template>
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
</template> -->
  

<template>
  <div v-if="role=='admin'">
    <div class="create-field-container">
      <div class="card">
        <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>

        <h2 class="form-title">Создать сферу деятельности</h2>
        <div class="form-group">
          <InputGroup>
            <InputGroupAddon>
              <i class="pi pi-user"></i>
            </InputGroupAddon>
            <InputText v-model="name" placeholder="Название" :class="{ 'p-invalid': !name }" />
          </InputGroup>
        </div>

        <div class="form-group">
          <InputGroup>
            <InputGroupAddon>
              <i class="pi pi-book"></i>
            </InputGroupAddon>
            <InputText v-model="description" placeholder="Описание" :class="{ 'p-invalid': !description }"/>
          </InputGroup>
        </div>

        <div class="form-group">
          <InputGroup>
            <InputGroupAddon>
              <i class="pi pi-book"></i>
            </InputGroupAddon>
            <InputNumber v-model="cost" :minFractionDigits="0" :maxFractionDigits="2" placeholder="Вес" :class="{ 'p-invalid': !cost }"/>
          </InputGroup>
        </div>

        <div class="form-group">
          <Button @click="createActivityField" class="create-button" :disabled="!name || !description || !cost">Создать</Button>
        </div>
      </div>
    </div>
  </div>
  <div v-else>
  </div>  
</template>

<style scoped>
.create-field-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.card {
  max-width: 500px;
  width: 100%;
  padding: 20px;
  background-color: #f8f8f8;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.form-title {
  text-align: center;
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.create-button {
  width: 100%;
}

@media (max-width: 768px) {
  .create-field-container {
    height: auto;
    padding: 20px;
  }
}
</style>

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