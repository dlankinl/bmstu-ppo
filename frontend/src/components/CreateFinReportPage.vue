<template>
  <div class="create-fin-report-container">
    <!-- <div class="card flex flex-column md:flex-row gap-3"> -->
    <div class="card">
      <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
      <h2 class="form-title">Добавить отчет</h2>
      <div class="form-group">
        <InputGroup>
          <InputGroupAddon>
            <i class="pi pi-money-bill"></i>
          </InputGroupAddon>
          <InputNumber v-model="revenue" :minFractionDigits="0" :maxFractionDigits="2" placeholder="Выручка" :class="{ 'p-invalid': !revenue }"/>
        </InputGroup>
      </div>

      <div class="form-group">
        <InputGroup>
          <InputGroupAddon>
            <i class="pi pi-money-bill"></i>
          </InputGroupAddon>
          <InputNumber v-model="costs" :minFractionDigits="0" :maxFractionDigits="2" placeholder="Расходы" :class="{ 'p-invalid': !costs }"/>
        </InputGroup>
      </div>

      <div class="form-group">
        <InputGroup>
          <InputGroupAddon>
            <i class="pi pi-calendar"></i>
          </InputGroupAddon>
          <InputNumber v-model="year" :useGrouping="false" :min="1900" :max="`${new Date().getFullYear()}`" placeholder="Год" :class="{ 'p-invalid': !year }"/>
        </InputGroup>
      </div>

      <div class="form-group">
        <InputGroup>
          <InputGroupAddon>
            <i class="pi pi-calendar-clock"></i>
          </InputGroupAddon>
          <InputNumber v-model="quarter" :min="1" :max="4" placeholder="Квартал" :class="{ 'p-invalid': !quarter }"/>
        </InputGroup>
      </div>
      <div class="form-group">
        <Button @click="createFinReport">Создать</Button>
      </div>
    </div>
  </div>
</template>
  
<style scoped>
.create-fin-report-container {
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
  name: 'CreateFinReportPage',
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
    if (Utils.isAuth()) {
      this.role = Utils.getUserRoleJWT();
    }
  },
  methods: {
    createFinReport() {
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