<template>
  <div class="card">
    <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
    <template v-if="entId==visitorId">
      <RouterLink :to="`/contacts/create`"><Button label="Добавить" icon="pi pi-plus-circle"></Button></RouterLink>
    </template>
    <DataView :value="contacts" :rows="rows" :totalRecords="totalPages*rows">
      <template #list="slotProps">
        <div class="grid grid-nogutter">
          <div
            v-for="(item, index) in slotProps.items"
            :key="index"
            class="col-12 relative"
          >
            <div
              class="flex flex-column sm:flex-row sm:align-items-center p-4 gap-3"
              :class="{ 'border-top-1 surface-border': index !== 0 }"
            >
              <div
                class="flex flex-column md:flex-row justify-content-between md:align-items-center flex-1 gap-4"
              >
                <div
                  class="flex flex-row md:flex-column justify-content-between align-items-start gap-2"
                >
                  <div>
                    <div class="text-lg font-medium text-900 mt-2">
                      <template v-if="item==null"></template>
                      <template v-else="item.name">
                        {{ item.name }}
                        <template v-if="entId==visitorId">
                          <Button icon="pi pi-trash" class="p-button-rounded p-button-secondary p-button-text" @click="deleteContact(item)"></Button>
                          <RouterLink :to="`/contacts/${item.id}/update`"><Button icon="pi pi-wrench" class="p-button-rounded p-button-secondary p-button-text"></Button></RouterLink>
                        </template>
                      </template>
                    </div>
                  </div>
                  <div
                    class="surface-100 p-1"
                    style="border-radius: 30px"
                  ></div>
                </div>
                <div class="flex flex-column md:align-items-end gap-5">
                  <span class="text-xl font-semibold text-900"
                    >
                      <template v-if="item==null"></template>
                      <template v-else>
                        <i class="pi pi-book"></i>
                      {{ item.value }}
                      </template>
                    </span
                  >
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </DataView>
  </div>
</template>
  
  
<script>
import ContactsService from '../services/contacts.service';
import Utils from '../services/auth-header';
import Button from 'primevue/button';
import DataView from 'primevue/dataview';
import Message from 'primevue/message';

export default {
  name: 'ContactsPage',
  components: {
    Button,
    DataView,
    Message
  },
  data() {
    return {
      first: 0,
      rows: 5,
      totalPages: 0,
      contacts: [],
      visitorId: null,
      entId: null,
      message: null
    }
  },
  created() {
    if (Utils.isAuth()) {
      this.visitorId = Utils.getUserIdJWT();
    }
    this.fetchContacts();
  },
  methods: {
    fetchContacts() {
      this.entId = this.$route.params.id;
      ContactsService.getByOwnerId(this.entId)
        .then(response => {
          this.contacts = [...response.data.data.contacts];
        })
        .catch(error => {
          console.error('Ошибка получения списка средств связи предпринимателя: ', error)
        })
    },
    deleteContact(contact) {
      ContactsService.deleteContact(contact.id)
        .then(response => {
          this.message = { severity: 'success', content: 'Средство связи удалено.' };
          this.fetchContacts();
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Произошла ошибка при удалении средства связи: ${error.response.data.error}` };
          console.error("Ошибка при удалении средства связи:", error);
        })
    }
  }
}
</script>