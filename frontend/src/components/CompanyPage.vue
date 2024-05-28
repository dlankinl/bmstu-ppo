<template>
  <div>
    <h1>Информация о компании</h1>
    <div>
      <h2>Название: {{ company.name }}</h2>
      <p>Сфера деятельности: {{ company.activity_field.name }}</p>
      <p>Город: {{ company.city }}</p>
    </div>
    <RouterLink :to="`/companies/${company.id}/update`">
    <Button 
      icon="pi pi-cog"
      label="Обновить информацию"
      class="flex-auto md:flex-initial white-space-nowrap"
    >
    </Button>
    </RouterLink>
  </div>
</template>
  
  <script>
  import FinReportService from '../services/financials.service';
  import CompaniesService from '../services/companies.service';
  import ActivityFieldsService from '../services/activity-fields.service';
  import DataTable from 'primevue/datatable';
  import Column from 'primevue/column';
  import Button from 'primevue/button';
  // import isAuth from '../services/auth-header';
  import Utils from '../services/auth-header';
  
  export default {
    name: 'EntrepreneurPage',
    components: {
      Button,
      DataTable,
      Column,
    },
    data() {
      return {
        company: {},
        financials: {"revenue": 0.0, "profit": 0.0, "costs": 0.0},
        currentPage: 1,
        rows: 3,
        totalPages: 0,
        isAuthValue: null
      }
    },
    created() {
      this.isAuth()
      this.fetchCompanyDetails()
      // this.fetchFinancials()
    },
    methods: {
      fetchCompanyDetails() {
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
      onPageChange(event) {
        this.currentPage = event.page + 1;
        // this.fetchCompanies();
      },
    }
  }
  </script>