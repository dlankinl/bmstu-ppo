<template>
  <div class="card">
    <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
    <template v-if="role=='admin'">
      <RouterLink :to="`/skills/create`"><Button label="Добавить" icon="pi pi-plus-circle"></Button></RouterLink>
    </template>
    <DataView :value="skills" paginator :rows="rows" :totalRecords="totalPages*rows" @page="onPageChange">
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
                      <template v-else="item.name">
                        {{ item.name }}
                        <template v-if="role=='admin'">
                          <Button icon="pi pi-trash" class="p-button-rounded p-button-secondary p-button-text" @click="deleteSkill(item.id)"></Button>
                          <RouterLink :to="`/skills/${item.id}/update`"><Button icon="pi pi-wrench" class="p-button-rounded p-button-secondary p-button-text"></Button></RouterLink>
                        </template>
                        <template v-if="role!='guest'">
                          <Button icon="pi pi-plus-circle" class="p-button-rounded p-button-secondary p-button-text" @click="addSkill(item.id)"></Button>
                        </template>
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
                      <template v-if="item==null"></template>
                      <template v-else>
                        <i class="pi pi-book"></i>
                      {{ item.description }}
                      </template>
                    </span
                  >
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
import UserSkillsService from '../services/user-skills.service';
import SkillsService from '../services/skills.service'
import Utils from '../services/auth-header';
import Button from 'primevue/button';
import DataView from 'primevue/dataview';
import Chip from 'primevue/chip';
import Message from 'primevue/message';
import UserSkillModel from '../models/UserSkillModel';

export default {
  name: 'SkillsPage',
  components: {
    Button,
    DataView,
    Chip,
    Message
  },
  data() {
    return {
      first: 0,
      rows: 3,
      totalPages: 0,
      skills: [],
      currentPage: 1,
      role: 'guest',
      message: null
    }
  },
  created() {
    if (Utils.isAuth()) {
      this.visitorId = Utils.getUserIdJWT();
      this.role = Utils.getUserRoleJWT();
    }
    this.fetchSkills();
  },
  methods: {
    fetchSkills() {
      SkillsService.getSkills(this.currentPage)
        .then(response => {
          this.skills = [...Array((this.currentPage - 1) * 3).fill(null), ...response.data.data.skills];
          this.totalPages = response.data.data.num_pages;
        })
        .catch(error => {
          console.error('Ошибка получения списка навыков: ', error)
        })
    },
    onPageChange(event) {
      this.currentPage = event.page + 1;
      this.fetchSkills();
    },
    deleteSkill(skillId) {
      SkillsService.deleteSkill(skillId)
        .then(response => {
          this.message = { severity: 'success', content: 'Навык удален.' };
          this.fetchSkills();
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Произошла ошибка при удалении навыка: ${error.response.data.error}` };
          console.error("Ошибка при удалении навыка:", error);
        })
    },
    addSkill(skillId) {
      const userSkill = new UserSkillModel(this.visitorId, skillId);
      UserSkillsService.createUserSkill(userSkill)
        .then(response => {
          this.message = { severity: 'success', content: 'Навык добавлен.' };
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Произошла ошибка при добавлении навыка: ${error.response.data.error}` };
          console.error("Ошибка при добавлении навыка:", error);
        })
    }
  }
}
</script>