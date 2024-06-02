<template>
  <div class="pt-4" v-if="visitorId!=entId">
    <div class="card flex flex-column justify-content-center gap-3">
      <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
      <FloatLabel>
        <Textarea v-model="pros" rows="5" cols="30"></Textarea>
        <label>Понравилось</label>
      </FloatLabel>

      <FloatLabel>
        <Textarea v-model="cons" rows="5" cols="30"></Textarea>
        <label>Не понравилось</label>
      </FloatLabel>

      <FloatLabel>
        <Textarea v-model="description" rows="5" cols="30"></Textarea>
        <label>Описание</label>
      </FloatLabel>

      Оценка: <Rating v-model="rating" :cancel="false" />
      <Button @click="createReview" :disabled="!pros || !cons || !rating">Создать</Button>
    </div>
  </div>
  <div v-else>
    <h2>Вы не можете добавить себе отзыв.</h2>
  </div>  
</template>
  
<script>
import ReviewsService from '../services/reviews.service';
import Button from 'primevue/button';
import ReviewModel from '../models/ReviewModel.js'
import Message from 'primevue/message';
import Utils from '../services/auth-header';
import Rating from 'primevue/rating';
import Textarea from 'primevue/textarea';
import FloatLabel from 'primevue/floatlabel';
import { ref } from 'vue';
  
export default {
  name: 'CreateReviewPage',
  components: {
    Button,
    Message,
    Rating,
    Textarea,
    FloatLabel
  },
  setup() {
    const pros = ref('');
    const cons = ref('');
    const description = ref('');
    const rating = ref(null);
    const message = ref(null);
    const visitorId = null;
    const entId = null;

    return {
      pros,
      cons,
      description,
      rating,
      message,
      visitorId,
      entId
    };
  },
  created() {
    if (Utils.isAuth()) {
      this.role = Utils.getUserRoleJWT();
      this.visitorId = Utils.getUserIdJWT();
    }
    this.entId = this.$route.params.id;
  },
  methods: {
    createReview() {
      const rev = new ReviewModel(0, this.visitorId, this.$route.params.id, this.pros, this.cons, this.description, this.rating);
      ReviewsService.createReview(rev)
        .then(response => {
          this.message = { severity: 'success', content: 'Отзыв добавлен.' };
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Произошла ошибка при добавлении отзыва: ${error.response.data.error}` };
          console.error("Ошибка при добавлении отзыва:", error);
        })
    }
  }
}

</script>