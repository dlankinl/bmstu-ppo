<template>
  <div class="profile-container">
    <template v-if="entrepreneur.full_name">
      <div class="profile-header">
        <h1>Ваш профиль</h1>
        <div class="rating">
          <i class="pi pi-star"></i>
          {{ (5 * entrepreneur.rating).toFixed(1) }}
        </div>
      </div>
      <div class="profile-info">
        <h2>{{ entrepreneur.full_name }}</h2>
        <p><i class="pi pi-calendar"></i> {{ formatBirthday(entrepreneur.birthday) }}</p>
        <p><i class="pi pi-map-marker"></i> {{ entrepreneur.city }}</p>
        <p><i class="pi pi-venus-mars"></i> {{ formatGender(entrepreneur.gender) }}</p>
      </div>
    </template>
    <template v-else>
      <h1 class="no-info">Ваш профиль не заполнен. Обратитесь к администратору.</h1>
    </template>
    <div v-if="entrepreneur.full_name" class="settings">
      <ButtonGroup class="button-group">
        <RouterLink :to="`/entrepreneurs/${entrepreneur.id}/contacts`">
          <Button label="Средства связи" icon="pi pi-address-book" class="button"></Button>
        </RouterLink>
        <RouterLink :to="`/entrepreneurs/${entrepreneur.id}/skills`">
          <Button label="Навыки" icon="pi pi-bolt" class="button"></Button>
        </RouterLink>
        <RouterLink :to="`/entrepreneurs/${entrepreneur.id}/companies`">
          <Button label="Компании" icon="pi pi-building" class="button"></Button>
        </RouterLink>
        <RouterLink :to="`/entrepreneurs/${entrepreneur.id}/reviews`">
          <Button label="Отзывы" icon="pi pi-star" class="button"></Button>
        </RouterLink>
        <RouterLink :to="`/profile/reviews`">
          <Button label="Мои отзывы" icon="pi pi-star" class="button"></Button>
        </RouterLink>
      </ButtonGroup>
    </div>
    <div v-if="role == 'admin'" class="admin-settings">
      <ButtonGroup class="button-group">
        <RouterLink :to="`/activity-fields`">
          <Button label="Все сферы деятельности" icon="pi pi-address-book" class="button"></Button>
        </RouterLink>
        <RouterLink :to="`/skills`">
          <Button label="Все навыки" icon="pi pi-bolt" class="button"></Button>
        </RouterLink>
      </ButtonGroup>
    </div>
  </div>
</template>

<style scoped>
.profile-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background-color: #f8f8f8;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.profile-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.rating {
  background-color: #fff;
  padding: 5px 10px;
  border-radius: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 5px;
}

.profile-info {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.profile-info h2 {
  margin-top: 0;
}

.profile-info p {
  margin-bottom: 10px;
}

.profile-info p i {
  margin-right: 5px;
}

.no-info {
  text-align: center;
  color: #888;
}

.settings,
.admin-settings {
  margin-bottom: 20px;
}

.button-group {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 10px;
}

.button {
  flex: 1 0 calc(33.33% - 10px);
  max-width: 200px;
}

@media (max-width: 768px) {
  .button {
    flex: 0 0 calc(50% - 5px);
  }
}
</style>

<script>
import EntrepreneurService from '../services/entrepreneur.service'
import ButtonGroup from 'primevue/buttongroup';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Utils from '../services/auth-header';

export default {
  name: 'ProfilePage',
  components: {
    Button,
    Column,
    ButtonGroup
  },
  data() {
    return {
      entrepreneur: {},
      isAuthValue: null,
      role: 'guest'
    }
  },
  created() {
    if (Utils.isAuth()) {
      this.role = Utils.getUserRoleJWT();
    }
    this.fetchEntrepreneurDetails()
  },
  methods: {
    fetchEntrepreneurDetails() {
      const entrepreneurId = Utils.getUserIdJWT();
      EntrepreneurService.getEntrepreneurDetails(entrepreneurId)
        .then(response => {
          const entrepreneur = response.data.data.entrepreneur
          if (entrepreneur.full_name) {
              this.fetchRating(entrepreneur)
          } else {
              this.entrepreneur = entrepreneur
          }
        })
        .catch(error => {
          console.error('Error fetching entrepreneur details:', error)
        })
    },
    formatBirthday(birthday) {
      const date = new Date(birthday)
      return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
    },
    formatGender(gender) {
      return gender === 'm' ? 'мужской' : 'женский'
    },
    isAuth() {
      this.isAuthValue = Utils.isAuth();
    },
    fetchRating(entrepreneur) {
      EntrepreneurService.getEntrepreneurRating(entrepreneur.id)
        .then(response => {
          entrepreneur.rating = response.data.data.rating
          this.entrepreneur = entrepreneur
        })
        .catch(error => {
          console.error('Error fetching entrepreneur rating:', error)
        })
    },
  }
}
</script>