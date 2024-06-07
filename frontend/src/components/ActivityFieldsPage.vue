<!-- <template>
  <div class="card">
    <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
    <template v-if="entId==visitorId">
      <RouterLink :to="`/activity-fields/create`"><Button label="Добавить" icon="pi pi-plus-circle"></Button></RouterLink>
    </template>
    <DataView :value="fields" paginator :rows="rows" :totalRecords="totalPages*rows" @page="onPageChange">
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
                        <i class="pi pi-tag"></i>
                        {{ item.cost }}
                        <RouterLink :to="`/activity-fields/${item.id}/update`"><Button icon="pi pi-wrench" class="p-button-rounded p-button-secondary p-button-text"></Button></RouterLink>
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
                        {{ item.description }}
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
</template> -->



<template>
  <div class="card">
    <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
    <div class="add-button-container">
      <template v-if="entId == visitorId">
        <RouterLink :to="`/activity-fields/create`">
          <Button label="Добавить" icon="pi pi-plus-circle" class="add-button"></Button>
        </RouterLink>
      </template>
    </div>
    <DataView :value="fields" paginator :rows="rows" :totalRecords="totalPages * rows" @page="onPageChange">
      <template #list="slotProps">
        <div class="grid grid-nogutter">
          <div
            v-for="(item, index) in slotProps.items"
            :key="index"
            class="col-12 relative activity-field-card"
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
                      <template v-if="item == null"></template>
                      <template v-else="item.name">
                        {{ item.name }}
                        <i class="pi pi-tag"></i>
                        {{ item.cost }}
                        <RouterLink :to="`/activity-fields/${item.id}/update`">
                          <Button
                            icon="pi pi-wrench"
                            class="p-button-rounded p-button-secondary p-button-text activity-field-button"
                          ></Button>
                        </RouterLink>
                      </template>
                    </div>
                  </div>
                  <div class="surface-100 p-1" style="border-radius: 30px"></div>
                </div>
                <div class="flex flex-column md:align-items-end gap-5">
                  <span class="text-xl font-semibold text-900">
                    <template v-if="item == null"></template>
                    <template v-else>
                      <i class="pi pi-book"></i>
                      {{ item.description }}
                    </template>
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </DataView>
  </div>
</template>

<style scoped>
.card {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background-color: #f8f8f8;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.add-button-container {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 20px;
}

.add-button {
  background-color: #2196f3;
  color: #fff;
}

.activity-field-card {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: box-shadow 0.3s ease;
}

.activity-field-card:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.activity-field-button {
  margin-left: 5px;
}

@media (max-width: 768px) {
  .activity-field-button {
    margin-top: 5px;
    margin-left: 0;
  }
}
</style>

  
<script>
import ActivityFieldsService from '../services/activity-fields.service';
import Utils from '../services/auth-header';
import Button from 'primevue/button';
import DataView from 'primevue/dataview';
import Message from 'primevue/message';

export default {
  name: 'ActivityFieldsPage',
  components: {
    Button,
    DataView,
    Message
  },
  data() {
    return {
      first: 0,
      rows: 3,
      totalPages: 0,
      fields: [],
      currentPage: 1,
      message: null
    }
  },
  created() {
    if (Utils.isAuth()) {
      this.role = Utils.getUserRoleJWT();
      if (this.role !== 'admin') {
        this.$router.push('/404');
      }
    }
    this.fetchFields();
  },
  methods: {
    fetchFields() {
      ActivityFieldsService.getPaginatedFields(this.currentPage)
        .then(response => {
          this.fields = [...Array((this.currentPage - 1) * 3).fill(null), ...response.data.data.activity_fields];
          this.totalPages = response.data.data.num_pages;
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Ошибка при получении сфер деятельности: ${error.response.data.error}` };
          console.error("Ошибка при получении сфер деятельности:", error)
        })
    },
    onPageChange(event) {
      this.currentPage = event.page + 1;
      this.fetchFields();
    },
  }
}
</script>