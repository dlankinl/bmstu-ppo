<template>
  <div class="create-contact-container">
    <div class="card">
      <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
      <h2 class="form-title">Создать контакт</h2>
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
          <InputText v-model="value" placeholder="Значение" :class="{ 'p-invalid': !value }" />
        </InputGroup>
      </div>

      <div class="form-group">
        <Button @click="createContact" class="create-button" :disabled="!name || !value">Создать</Button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.create-contact-container {
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
  .create-contact-container {
    height: auto;
    padding: 20px;
  }
}
</style>

  
<script>
import ContactsService from '../services/contacts.service';
import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import ContactModel from '../models/ContactModel.js'
import Message from 'primevue/message';
import Utils from '../services/auth-header';
import { ref } from 'vue';
  
export default {
  name: 'CreateContactPage',
  components: {
    Button,
    InputGroup,
    InputGroupAddon,
    InputText,
    Message
  },
  setup() {
    const name = ref('');
    const value = ref('');
    var message = ref('');

    return {
      name,
      value,
      message
    };
  },
  methods: {
    createContact() {
      const ownerId = Utils.getUserIdJWT();
      const contact = new ContactModel(0, this.name, this.value);

      ContactsService.createContact(contact)
        .then(response => {
          this.message = { severity: 'success', content: 'Средство связи успешно добавлено' }
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Произошла ошибка при добавлении средства связи: ${error.response.data.error}` }
          console.error(error)
        })
    }
  }
}

</script>