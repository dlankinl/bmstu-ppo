<template>
  <div class="company-details">
    <div class="header">
      <h1>Информация о компании</h1>
      <div class="company-name">{{ company.name }}</div>
    </div>
    <div class="info-card">
      <p><span class="label">Сфера деятельности:</span> {{ company.activity_field.name }}</p>
      <p><span class="label">Город:</span> {{ company.city }}</p>
    </div>
    <div class="button-group">
      <template v-if="entId == visitorId">
        <ButtonGroup>
          <RouterLink :to="`/companies/${company.id}/update`">
            <Button icon="pi pi-cog" label="Обновить информацию" class="button"></Button>
          </RouterLink>
          <RouterLink :to="`/companies/${company.id}/financials/create`">
            <Button icon="pi pi-money-bill" label="Добавить отчет" class="button"></Button>
          </RouterLink>
        </ButtonGroup>
      </template>
    </div>
    <CompanyFinancialsPage />
  </div>
</template>

<style scoped>
.company-details {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background-color: #f8f8f8;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.company-name {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.info-card {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.info-card p {
  margin-bottom: 10px;
}

.info-card .label {
  font-weight: bold;
}

.button-group {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-bottom: 20px;
}

.button {
  flex: 1;
}

@media (max-width: 768px) {
  .button-group {
    flex-wrap: wrap;
  }

  .button {
    flex: 0 0 calc(50% - 5px);
  }
}
</style>

  
<script>
import CompaniesService from '../services/companies.service';
import ActivityFieldsService from '../services/activity-fields.service';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Buttongroup from 'primevue/buttongroup';
import CompanyFinancialsPage from './CompanyFinancialsPage.vue';
import Utils from '../services/auth-header';

export default {
  name: 'EntrepreneurPage',
  components: {
    Button,
    DataTable,
    Column,
    Buttongroup,
    CompanyFinancialsPage
  },
  data() {
    return {
      company: {},
      financials: {"revenue": 0.0, "profit": 0.0, "costs": 0.0},
      currentPage: 1,
      rows: 3,
      totalPages: 0,
      isAuthValue: null,
      visitorId: null,
      entId: null,
    }
  },
  created() {
    if (Utils.isAuth()) {
      this.visitorId = Utils.getUserIdJWT();
    }
    this.fetchCompanyDetails()
  },
  methods: {
    fetchCompanyDetails() {
      this.entId = this.$route.params.id;
      console.log(this.entId, this.visitorId);
      const companyId = this.$route.params.id
      CompaniesService.getCompanyDetails(companyId)
        .then(response => {
          this.company = response.data.data.company;
          ActivityFieldsService.getActivityField(this.company.activity_field_id)
            .then((response) => {
              this.company.activity_field = response.data.data.activity_field;
            })
            .catch(error => {
              console.error("Ошибка получения сферы деятельности компании:", error)
            })
        })
        .catch(error => {
          console.error('Ошибка получения информации о компании:', error)
        })
    },
    isAuth() {
      this.isAuthValue = Utils.isAuth();
    },
  }
}
</script>