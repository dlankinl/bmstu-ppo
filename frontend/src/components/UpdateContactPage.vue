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

  <Button @click="updateContact">Обновить</Button>
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
  name: 'UpdateContactPage',
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
  created() {
    this.fillInfo()
  },
  methods: {
    fillInfo() {
      ContactsService.getById(this.$route.params.id)
        .then(response => {
          this.contact = response.data.data.contact;
          this.name = this.contact.name;
          this.value = this.contact.value;
        })
        .catch(error => {
          console.error(error)
          this.message = { severity: 'error', content: `Произошла ошибка при получении информации о средстве связи: ${error.response.data.error}` };
        })
    },
    updateContact() {
      const contact = new ContactModel(this.$route.params.id, this.name, this.value);
      console.log(contact);

      ContactsService.updateContact(this.$route.params.id, contact)
        .then(response => {
          this.message = { severity: 'success', content: 'Средство связи успешно обновлено' }
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Произошла ошибка при обновлении средства связи: ${error.response.data.error}` }
          console.error(error)
        })
    }
  }
}

</script>