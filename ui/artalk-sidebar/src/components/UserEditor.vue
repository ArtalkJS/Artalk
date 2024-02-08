<script setup lang="ts">
import { artalk, bootParams } from '../global'
import type { ArtalkType } from 'artalk'

const { t } = useI18n()

interface IUserEditData {
  id: number
  name: string
  email: string
  link: string
  password?: string
  badge_name: string
  badge_color: string
  is_admin: boolean
  receive_email: boolean
}

const props = defineProps<{
  user?: IUserEditData
}>()

const emit = defineEmits<{
  (evt: 'close'): void
  (evt: 'update', user: ArtalkType.UserDataForAdmin): void
}>()

const editUser = ref<IUserEditData>()

const isLoading = ref(false)
const isCreateMode = ref(false)

const showFullDetails = ref(false)

onBeforeMount(() => {
  if (!props.user) {
    isCreateMode.value = true
    editUser.value = {
      id: 0,
      name: '',
      email: '',
      link: '',
      password: '',
      badge_name: '',
      badge_color: '',
      is_admin: false,
      receive_email: true,
    }
  } else {
    editUser.value = {
      ...props.user,
      password: '',
    }
  }
})

function close() {
  emit('close')
}

function submit() {
  // 执行保存
  isLoading.value = true

  // 默认标签颜色
  if (editUser.value!.badge_name !== '' && editUser.value!.badge_color === '') {
    editUser.value!.badge_color = '#0083ff'
  }

  if (isCreateMode.value) {
    artalk!.ctx.getApi().users.createUser({
      ...editUser.value!
    })
      .then((res) => {
        emit('update', res.data)
      }).catch((e: ArtalkType.FetchError) => {
        alert('用户创建错误：'+e.message)
      }).finally(() => {
        isLoading.value = false
      })
  } else {
    const user = editUser.value!
    artalk!.ctx.getApi().users.updateUser(user.id, {
      ...editUser.value!
    })
      .then((res) => {
        emit('update', res.data)
      }).catch((e: ArtalkType.FetchError) => {
        alert('用户保存错误：'+e.message)
      }).finally(() => {
        isLoading.value = false
      })
  }
}
</script>

<template>
  <div class="user-editor-layer">
    <div class="header">
      <div class="title">{{ isCreateMode ? t('userCreate') : t('userEdit') }}</div>
      <div v-if="!isCreateMode" class="close-btn" @click="close()">
        <i class="atk-icon atk-icon-close"></i>
      </div>
    </div>
    <div v-if="!isCreateMode" class="user-log">
      <div><span>{{ t('comments') }}</span>{{(editUser as any).comment_count}}</div>
      <div><span>{{ t('last') }} IP</span>{{(editUser as any).last_ip || '-'}}</div>
      <div><span>{{ t('last') }} UA</span>
        <template v-if="showFullDetails || !(editUser as any).last_ua">{{(editUser as any).last_ua || '-'}}</template>
        <template v-else><span @click="showFullDetails = true" style="cursor: pointer;color: var(--at-color-main)">{{ t('Show') }}</span></template>
      </div>
    </div>
    <form v-if="editUser" class="atk-form" @submit.prevent="submit()">
      <div class="atk-label required">{{ t('username') }}</div>
      <input v-model="editUser.name" type="text" placeholder="" autocomplete="off">
      <div class="atk-label required">{{ t('email')  }}</div>
      <input v-model="editUser.email" type="text" placeholder="" autocomplete="off">
      <div class="atk-label">{{ t('link') }}</div>
      <input v-model="editUser.link" type="text" placeholder="" autocomplete="off">
      <div class="atk-label">{{ t('badgeText') }}</div>
      <input v-model="editUser.badge_name" type="text" placeholder="" autocomplete="off">
      <div class="atk-label">{{ t('badgeColor') }} (Color Hex)</div>
      <input v-model="editUser.badge_color" type="text" placeholder="" autocomplete="off">
      <div class="atk-label required">{{ t('role') }}</div>
      <select v-model="editUser.is_admin">
        <option :value="false">{{ t('normal') }}</option>
        <option :value="true">{{ t('admin') }}</option>
      </select>
      <template v-if="editUser.is_admin">
        <div class="atk-label required">{{ t('password') }}</div>
        <input v-model="editUser.password" type="text" :placeholder="isCreateMode ? '' : `(${t('passwordEmptyHint')})`" autocomplete="off">
      </template>
      <div class="atk-label required">{{ t('emailNotify') }}</div>
      <select v-model="editUser.receive_email">
        <option :value="true">{{ t('enabled') }}</option>
        <option :value="false">{{ t('disabled') }}</option>
      </select>
      <button type="submit" class="atk-btn">{{ t('save') }}</button>
    </form>
    <LoadingLayer v-if="isLoading" />
  </div>
</template>

<style scoped lang="scss">
.user-editor-layer {
  z-index: 201;
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: var(--at-color-bg);
  overflow-y: auto;
}

.header {
  position: sticky;
  top: 0;
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 20px 20px 10px 20px;
  background: var(--at-color-bg-transl);

  .title {
    margin-left: 10px;
    flex: 1;
    font-size: 19px;
  }
}

.close-btn {
  width: 50px;
  height: 50px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;

  &:hover i::after {
    background-color: var(--at-color-red);
  }
}

.user-log {
  margin: 0 30px;
  padding: 15px 25px;
  border-radius: 2px;
  font-size: 13px;
  background: var(--at-color-bg-grey);
  line-height: 1.8;

  & > div > span:first-child {
    display: inline-block;
    width: 5em;
  }
}
</style>
