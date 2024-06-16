<!-- <template>
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
</template> -->

<template>
  <div class="pt-4" v-if="visitorId != entId">
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

<style scoped>
.pt-4 {
  padding-top: 1rem;
}

.card {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 1rem;
  padding: 1rem;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.float-label {
  position: relative;
  margin-bottom: 1rem;
}

.float-label textarea {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.float-label label {
  position: absolute;
  top: 0.5rem;
  left: 0.5rem;
  background-color: #fff;
  padding: 0 0.25rem;
  color: #333;
  font-size: 0.8rem;
  pointer-events: none;
  transition: top 0.3s, left 0.3s, font-size 0.3s;
}

.float-label textarea:focus + label,
.float-label textarea:not(:placeholder-shown) + label {
  top: -0.5rem;
  left: 0;
  font-size: 0.7rem;
}

.rating {
  margin-bottom: 1rem;
}

.button {
  width: 100%;
}

@media (max-width: 768px) {
  .card {
    padding: 0.5rem;
  }
}
</style>
  
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