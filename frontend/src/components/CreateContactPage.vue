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
    <InputText v-model="value" placeholder="Значение" :class="{ 'p-invalid': !value }"/>
  </InputGroup>

  <Button @click="createContact">Создать</Button>
</div>
</template>
  
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

    return {
      name,
      value
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