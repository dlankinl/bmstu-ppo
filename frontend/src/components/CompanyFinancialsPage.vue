<template>
  <div class="card">
    <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
    <InputGroup>
      <InputNumber v-model="startYear" :useGrouping="false" :min="1900" :max="`${new Date().getFullYear()}`" placeholder="Начальный год" :class="{ 'p-invalid': !startYear }"/>
    </InputGroup>

    <InputGroup>
      <InputNumber v-model="startQuarter" :min="1" :max="4" placeholder="Начальный квартал" :class="{ 'p-invalid': !startQuarter }"/>
    </InputGroup>

    <InputGroup>
      <InputNumber v-model="endYear" :useGrouping="false" :min="1900" :max="`${new Date().getFullYear()}`" placeholder="Конечный год" :class="{ 'p-invalid': !endYear }"/>
    </InputGroup>

    <InputGroup>
      <InputNumber v-model="endQuarter" :min="1" :max="4" placeholder="Конечный квартал" :class="{ 'p-invalid': !endQuarter }"/>
    </InputGroup>

    <Button @click="fetchReport">Получить отчёт</Button>
    
    <template v-if="financialsAvailable">
      <p>Выручка: {{ revenue }}</p>
      <p>Расходы: {{ costs }}</p>
      <p>Прибыль: {{ profit }}</p>
    </template>
  </div>
</template>

  
<script>
import FinReportService from '../services/fin-report.service';
import Utils from '../services/auth-header';
import Button from 'primevue/button';
import DataView from 'primevue/dataview';
import Message from 'primevue/message';
import InputGroup from 'primevue/inputgroup';
import InputNumber from 'primevue/inputnumber';
import {ref} from 'vue';

export default {
  name: 'CompanyFinancialsPage',
  components: {
    Button,
    DataView,
    Message,
    InputNumber,
    InputGroup
  },
  data() {
    return {
      entId: null,
      message: null,
      financials: null,
      startYear: ref(null),
      startQuarter: ref(null),
      endYear: ref(null),
      endQuarter: ref(null),
      financialsAvailable: false,
    }
  },
  created() {
    if (Utils.isAuth()) {
      this.visitorId = Utils.getUserIdJWT();
    }
    this.fetchReport();
  },
  methods: {
    fetchReport() {
      this.compId = this.$route.params.id;
      if (this.startYear != null && this.startQuarter != null && this.endYear != null & this.endQuarter != null) {
        FinReportService.getCompanyReport(this.compId, this.startYear, this.endYear, this.startQuarter, this.endQuarter)
        .then(response => {
          this.report = response.data.data;
          this.revenue = this.report.revenue;
          this.costs = this.report.costs;
          this.profit = this.report.profit;
          this.financialsAvailable = true;
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Произошла ошибка при получении отчета компании: ${error.response.data.error}` }
          console.error('Ошибка получения отчета компании: ', error)
        })
      }
    } 
  }
}
</script>