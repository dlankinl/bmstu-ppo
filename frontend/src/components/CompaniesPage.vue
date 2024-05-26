<template>
  <div class="card">
    <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
    <template v-if="entId==visitorId">
      <RouterLink :to="`/companies/create`"><Button label="Добавить" icon="pi pi-plus-circle"></Button></RouterLink>
    </template>
    <DataView :value="companies" paginator :rows="rows" :totalRecords="totalPages*rows" @page="onPageChange">
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
                        {{ item.name }} <Chip :label="`${item.activity_field.name}`" />
                        <template v-if="entId==visitorId">
                          <Button icon="pi pi-trash" class="p-button-rounded p-button-secondary p-button-text" @click="deleteCompany(item)"></Button>
                          <RouterLink :to="`/companies/${item.id}/update`"><Button icon="pi pi-wrench" class="p-button-rounded p-button-secondary p-button-text"></Button></RouterLink>
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
                        <i class="pi pi-map-marker"></i>
                      {{ item.city }}
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

<style scoped>
.rating-container {
  position: absolute;
  top: 10px; /* Adjust the top position as needed */
  right: 10px; /* Adjust the right position as needed */
  background-color: #fff; /* Set the background color */
  padding: 5px 10px; /* Add some padding */
  border-radius: 5px; /* Add border-radius for a rounded corner */
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); /* Add a box-shadow for a subtle elevation */
  display: flex;
  align-items: center;
  gap: 5px; /* Add some gap between the icon and rating value */
}
</style>

  
<script>
import CompaniesService from '../services/companies.service';
import ActivityFieldsService from '../services/activity-fields.service';
import Utils from '../services/auth-header';
import Button from 'primevue/button';
import DataView from 'primevue/dataview';
import Chip from 'primevue/chip';
import Message from 'primevue/message';

export default {
  name: 'CompaniesPage',
  components: {
    Button,
    DataView,
    Chip,
    Message
  },
  data() {
    return {
      first: 0,
      rows: 3,
      totalPages: 0,
      companies: [],
      currentPage: 1,
      visitorId: null,
      entId: null,
      message: null
    }
  },
  created() {
    this.visitorId = Utils.getUserIdJWT();
    this.fetchCompanies();
  },
  methods: {
    fetchCompanies() {
      this.entId = this.$route.params.id;
      CompaniesService.getEntrepreneursCompanies(this.entId, this.currentPage)
        .then(response => {
          this.companies = [...Array((this.currentPage - 1) * 3).fill(null), ...response.data.data.companies];
          this.totalPages = response.data.data.num_pages;
          const uniqueActivityFieldIds = [...new Set(this.companies.filter(company => company !== null).map(company => company.activity_field_id))];
          var uniqueActivityFieldIdsMap = uniqueActivityFieldIds.map(id => ({ id, value: {} }));
          Promise.all(
            uniqueActivityFieldIdsMap
              .map(field =>
                ActivityFieldsService.getActivityField(field.id).then(response => {
                  field.value = response.data.data.activity_field;
                })
              )
            ).then(() => {
              this.companies
                .filter(company => company !== null)
                .forEach(company => {
                  company.activity_field = uniqueActivityFieldIdsMap.find(field => field.id === company.activity_field_id).value;
                });
            })
        })
        .catch(error => {
          console.error('Ошибка получения списка компаний предпринимателя: ', error)
        })
    },
    onPageChange(event) {
      this.currentPage = event.page + 1;
      this.fetchCompanies();
    },
    deleteCompany(company) {
      CompaniesService.deleteCompany(company.id)
        .then(response => {
          this.message = { severity: 'success', content: 'Компания удалена.' };
          this.fetchCompanies();
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Произошла ошибка при удалении компании: ${error.response.data.error}` };
          console.error("Ошибка при удалении компании:", error);
        })
    }
  }
}
</script>