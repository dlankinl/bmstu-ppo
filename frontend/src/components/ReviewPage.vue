<template>
  <div class="card">
    <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
    <template v-if="entId!=visitorId && role!='admin'">
      <RouterLink :to="`/entrepreneurs/${entId}/reviews/create`"><Button label="Добавить" icon="pi pi-plus-circle"></Button></RouterLink>
    </template>
    <DataView :value="reviews" paginator :rows="rows" :totalRecords="totalPages*rows" @page="onPageChange">
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
                      <template v-else="item">
                        Автор: <a :href="`/entrepreneurs/${item.reviewer_id}`">{{ item.author.full_name }}</a>
                        <div class="flex flex-column">
                          <p>Понравилось: {{ item.pros }}</p>
                          <p>Не понравилось: {{ item.cons }}</p>
                          <p v-if="item.description">Описание: {{ item.description }}</p>
                          <Rating v-model="item.rating" readonly :cancel="false" />
                        </div>
                        <template v-if="role=='admin'">
                          <Button icon="pi pi-trash" class="p-button-rounded p-button-secondary p-button-text" @click="deleteReview(item.id)"></Button>
                        </template>
                      </template>
                    </div>
                  </div>
                  <div
                    class="surface-100 p-1"
                    style="border-radius: 30px"
                  ></div>
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
import ReviewsService from '../services/reviews.service'
import Utils from '../services/auth-header';
import Button from 'primevue/button';
import DataView from 'primevue/dataview';
import Message from 'primevue/message';
import EntrepreneurService from '../services/entrepreneur.service';
import Rating from 'primevue/rating';

export default {
  name: 'ReviewsPage',
  components: {
    Button,
    DataView,
    Message,
    Rating
  },
  data() {
    return {
      first: 0,
      rows: 3,
      totalPages: 0,
      reviews: [],
      currentPage: 1,
      role: 'guest',
      message: null
    }
  },
  created() {
    this.entId = this.$route.params.id;
    if (Utils.isAuth()) {
      this.visitorId = Utils.getUserIdJWT();
      this.role = Utils.getUserRoleJWT();
    }
    this.fetchReviews();
  },
  methods: {
    fetchReviews() {
      ReviewsService.getEntrepreneurReviews(this.$route.params.id, this.currentPage)
        .then(response => {
          this.reviews = [...Array((this.currentPage - 1) * 3).fill(null), ...response.data.data.reviews];
          this.totalPages = response.data.data.num_pages;
          Promise.all(
            this.reviews
              .filter(rev => rev !== null)
              .map(rev =>
                EntrepreneurService.getEntrepreneurDetails(rev.reviewer_id).then(response => {
                  rev.author = response.data.data.entrepreneur;
                })
              )
            )
        })
        .catch(error => {
          console.error('Ошибка получения списка отзывов: ', error)
        })
    },
    onPageChange(event) {
      this.currentPage = event.page + 1;
      this.fetchReviews();
    },
    deleteReview(id) {
      ReviewsService.deleteReview(id)
        .then(response => {
          this.message = { severity: 'success', content: 'Отзыв удален.' };
          this.fetchReviews();
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Произошла ошибка при удалении отзыва: ${error.response.data.error}` };
          console.error("Ошибка при удалении отзыва:", error);
        })
    }
  }
}
</script>