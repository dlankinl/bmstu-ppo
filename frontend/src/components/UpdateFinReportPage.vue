<template>
    <div class="card flex flex-column md:flex-row gap-3">
      <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-money-bill"></i>
        </InputGroupAddon>
        <InputNumber v-model="revenue" :minFractionDigits="0" :maxFractionDigits="2" placeholder="Выручка" :class="{ 'p-invalid': !revenue }"/>
      </InputGroup>
  
      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-money-bill"></i>
        </InputGroupAddon>
        <InputNumber v-model="costs" :minFractionDigits="0" :maxFractionDigits="2" placeholder="Расходы" :class="{ 'p-invalid': !costs }"/>
      </InputGroup>
  
      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-calendar"></i>
        </InputGroupAddon>
        <InputNumber v-model="year" :min="1900" :max="`${new Date().getFullYear()}`" placeholder="Год" :class="{ 'p-invalid': !year }"/>
      </InputGroup>
  
      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-calendar-clock"></i>
        </InputGroupAddon>
        <InputNumber v-model="quarter" :min="1" :max="4" placeholder="Квартал" :class="{ 'p-invalid': !quarter }"/>
      </InputGroup>
      <Button @click="updateFinReport">Создать</Button>
    </div>
  </template>
    
  <script>
  import InputGroup from 'primevue/inputgroup';
  import InputGroupAddon from 'primevue/inputgroupaddon';
  import InputNumber from 'primevue/inputnumber';
  import Button from 'primevue/button';
  import FinReportModel from '../models/FinReportModel.js'
  import Message from 'primevue/message';
  import Utils from '../services/auth-header';
  import FinReportService from '../services/fin-report.service.js';
  import { ref } from 'vue';
    
  export default {
    name: 'UpdateFinReportPage',
    components: {
      Button,
      InputGroup,
      InputGroupAddon,
      InputNumber,
      Message
    },
    setup() {
      const revenue = ref(null);
      const costs = ref(null);
      const year = ref(null);
      const quarter = ref(null);
      const message = ref(null);
      var role = 'guest';
      var companyId = null;
  
      return {
        revenue,
        costs,
        year,
        quarter,
        message,
        role
      };
    },
    created() {
      this.fillInfo();
      if (Utils.isAuth()) {
        this.role = Utils.getUserRoleJWT();
      }
    },
    methods: {
      fillInfo() {
        FinReportService.getFinReport(this.$route.params.id)
        .then(response => {
          const report = response.data.data.financial_report;
          this.companyId = report.company_id;
          this.revenue = report.revenue;
          this.costs = report.costs;
          this.year = report.year;
          this.quarter = report.quarter;
        })
        .catch(error => {
          console.error(error)
        })
      },
      updateFinReport() {
        this.companyId = this.$route.params.id;
        const report = new FinReportModel(0, this.companyId, this.revenue, this.costs, this.year, this.quarter);
  
        FinReportService.createFinReport(report)
          .then(response => {
            this.message = { severity: 'success', content: 'Финансовый отчет успешно добавлен' }
          })
          .catch(error => {
            this.message = { severity: 'error', content: `Произошла ошибка при добавлении финансового отчета: ${error.response.data.error}` }
            console.error(error)
          })
      }
    }
  }
  
  </script>