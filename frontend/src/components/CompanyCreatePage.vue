<!-- <template>
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
        <i class="pi pi-map"></i>
      </InputGroupAddon>
      <InputText v-model="city" placeholder="Город" :class="{ 'p-invalid': !city }"/>
    </InputGroup>

    <InputGroup>
      <InputGroupAddon>
        <i class="pi pi-venus"></i>
      </InputGroupAddon>
      <Dropdown v-model="selectedField" placeholder="Сфера деятельности" :options="fields" optionLabel="name" :class="{ 'p-invalid': !selectedField }" />
    </InputGroup>

    <Button @click="createCompany">Создать</Button>
  </div>
</template> -->
  


<template>
  <div class="create-company-container">
    <div class="card">
      <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
      <h2 class="form-title">Создать компанию</h2>
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
            <i class="pi pi-map"></i>
          </InputGroupAddon>
          <InputText v-model="city" placeholder="Город" :class="{ 'p-invalid': !city }" />
        </InputGroup>
      </div>

      <div class="form-group">
        <InputGroup>
          <InputGroupAddon>
            <i class="pi pi-venus"></i>
          </InputGroupAddon>
          <Dropdown
            v-model="selectedField"
            placeholder="Сфера деятельности"
            :options="fields"
            optionLabel="name"
            :class="{ 'p-invalid': !selectedField }"
          />
        </InputGroup>
      </div>

      <div class="form-group">
        <Button @click="createCompany" class="create-button">Создать</Button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.create-company-container {
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
  .create-company-container {
    height: auto;
    padding: 20px;
  }
}
</style>


<script>
import CompaniesService from '../services/companies.service';
import ActivityFieldsService from '../services/activity-fields.service';
import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';
import InputText from 'primevue/inputtext';
import Dropdown from 'primevue/dropdown';
import Button from 'primevue/button';
import CompanyModel from '../models/CompanyModel.js'
import Message from 'primevue/message';
import Utils from '../services/auth-header';
import { ref } from 'vue';
  
export default {
  name: 'UpdateCompanyPage',
  components: {
    Button,
    InputGroup,
    InputGroupAddon,
    InputText,
    Dropdown,
    Message
  },
  setup() {
    const name = ref('');
    const city = ref('');
    const selectedField = ref(null);
    const message = ref(null);
    const fields = ref(null);

    return {
      name,
      city,
      selectedField,
      message,
      fields
    };
  },
  created() {
    this.getFields();
  },
  methods: {
    createCompany() {
      const ownerId = Utils.getUserIdJWT();
      const company = new CompanyModel(0, this.name, this.selectedField.id, this.city, ownerId);

      CompaniesService.createCompany(company)
        .then(response => {
          this.message = { severity: 'success', content: 'Компания успешно добавлена' }
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Произошла ошибка при добавлении компании: ${error.response.data.error}` }
          console.error(error)
        })
    },
    getFields() {
      ActivityFieldsService.getFields()
        .then(response => {
          this.fields = response.data.data.activity_fields;
        })
        .catch(error => {
          console.error("Ошибка при получении всех сфер деятельности:", error)
        })
    }
  }
}

</script>