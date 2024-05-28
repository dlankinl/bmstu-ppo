<template>
  <div v-if="role=='admin'">
    <div class="card flex flex-column md:flex-row gap-3">
      <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-user"></i>
        </InputGroupAddon>
        <InputText v-model="name" placeholder="Название" :class="{ 'p-invalid': !name }" />
      </InputGroup>

      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-map"></i>
        </InputGroupAddon>
        <InputText v-model="description" placeholder="Описание" :class="{ 'p-invalid': !description }"/>
      </InputGroup>
      <Button @click="createSkill">Создать</Button>
    </div>
  </div>
  <div v-else>
  </div>  
</template>
  
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