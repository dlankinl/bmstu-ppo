<template>
  <div class="card">
    <DataView :value="entrepreneurs" paginator :rows="rows" :totalRecords="totalPages*rows" @page="onPageChange">
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
                      <template v-if="item==null">
                      </template>
                      <template v-else-if="item.full_name">
                        {{ item.full_name }}
                      </template>
                      <template v-else>
                        {{ item.username }} (карточка не заполнена) 
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
                    <template v-if="item==null">
                    </template>
                    <template v-else-if="item.full_name">
                      <i class="pi pi-map-marker"></i>
                      {{ item.city }}
                    </template>
                    </span
                  >
                  <div v-if="item==null"></div>
                  <div v-else-if="item.id" class="flex flex-row-reverse md:flex-row gap-2">
                    <RouterLink :to="`/entrepreneurs/${item.id}`">
                    <Button 
                      icon="pi pi-id-card"
                      label="Подробнее"
                      class="flex-auto md:flex-initial white-space-nowrap"
                    >
                    </Button>
                    </RouterLink>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="item==null"></div>
            <div v-else class="rating-container">
              <i class="pi pi-star"></i>
              {{ (5 * item.rating).toFixed(1) }}
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
import EntrepreneurService from '../services/entrepreneur.service';
import Button from 'primevue/button';
import DataView from 'primevue/dataview';
import { ref } from 'vue';

export default {
  name: 'EntrepreneursPage',
  components: {
    Button,
    DataView
  },
  // setup() {
  //   const entrepreneurs = ref();

  //   return entrepreneurs
  // },
  data() {
    return {
      first: 0,
      rows: 3,
      totalPages: 0,
      entrepreneurs: [],
      currentPage: 1,
    }
  },
  created() {
    this.fetchEntrepreneurs()
  },
  methods: {
    fetchEntrepreneurs() {
      EntrepreneurService.getEntrepreneurs(this.currentPage)
        .then(response => {
          this.entrepreneurs = [...Array((this.currentPage - 1) * 3).fill(null), ...response.data.data.users]
          this.totalPages = response.data.data.num_pages
          Promise.all(
              this.entrepreneurs
                .filter(entrepreneur => entrepreneur !== null)
                .map(entrepreneur =>
                  EntrepreneurService.getEntrepreneurRating(entrepreneur.id).then(response => {
                    entrepreneur.rating = response.data.data.rating
                  })
                )
            ).then(() => {})
        })
        .catch(error => {
          console.error(error)
        })
    },
    formatGender(gender) {
      return gender === 'm' ? 'мужской' : 'женский'
    },
    formatBirthday(birthday) {
      const date = new Date(birthday)
      return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
    },
    onPageChange(event) {
      this.currentPage = event.page + 1;
      this.fetchEntrepreneurs();
    },
  }
}
</script>