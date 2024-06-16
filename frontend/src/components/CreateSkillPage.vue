<template>
  <div v-if="role === 'admin'">
    <div class="create-skill-container">
      <div class="card">
        <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>

        <h2 class="form-title">Создать навык</h2>
        <div class="form-group">
          <InputGroup>
            <InputGroupAddon>
              <i class="pi pi-user"></i>
            </InputGroupAddon>
            <InputText v-model="name" placeholder="Название" :class="{ 'p-invalid': !name }" />
          </InputGroup>
        </div>

        <div class="form-group">
          <InputGroup>
            <InputGroupAddon>
              <i class="pi pi-map"></i>
            </InputGroupAddon>
            <InputText v-model="description" placeholder="Описание" :class="{ 'p-invalid': !description }"/>
          </InputGroup>
        </div>

        <div class="form-group">
          <Button @click="createSkill" class="create-button" :disabled="!name || !description">Создать</Button>
        </div>
      </div>
    </div>
  </div>
  <div v-else>
  </div>
</template>


<style scoped>
.create-skill-container {
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
  .create-skill-container {
    height: auto;
    padding: 20px;
  }
}
</style>


<script>
import SkillsService from '../services/skills.service';
import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import SkillModel from '../models/SkillModel.js'
import Message from 'primevue/message';
import Utils from '../services/auth-header';
import { ref } from 'vue';
  
export default {
  name: 'CreateSkillPage',
  components: {
    Button,
    InputGroup,
    InputGroupAddon,
    InputText,
    Message
  },
  setup() {
    const name = ref('');
    const description = ref('');
    const message = ref(null);
    var role = null;

    return {
      name,
      description,
      message,
      role
    };
  },
  created() {
    if (Utils.isAuth()) {
      this.role = Utils.getUserRoleJWT();
      if (this.role !== 'admin') {
        this.$router.push('/404');
      }
    }
  },
  methods: {
    createSkill() {
      const skill = new SkillModel(0, this.name, this.description);

      SkillsService.createSkill(skill)
        .then(response => {
          this.message = { severity: 'success', content: 'Навык успешно добавлена' }
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Произошла ошибка при добавлении навыка: ${error.response.data.error}` }
          console.error(error)
        })
    }
  }
}

</script>