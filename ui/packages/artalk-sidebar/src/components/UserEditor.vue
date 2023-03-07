<script setup lang="ts">
import { artalk, bootParams } from '../global'
import type { UserData, UserDataForAdmin } from 'artalk/types/artalk-data'

const { t } = useI18n()

interface IUserEditData {
  name: string
  email: string
  link: string
  password: string
  badge_name: string
  badge_color: string
  is_admin: boolean
  site_names_raw: string
  receive_email: boolean
}

const props = defineProps<{
  user?: IUserEditData
}>()

const emit = defineEmits<{
  (evt: 'close'): void
  (evt: 'update', user: UserDataForAdmin): void
}>()

const user = ref<IUserEditData>()

const isLoading = ref(false)
const isCreateMode = ref(false)

const showFullDetails = ref(false)

onBeforeMount(() => {
  if (!props.user) {
    isCreateMode.value = true
    user.value = {
      name: '',
      email: '',
      link: '',
      password: '',
      badge_name: '',
      badge_color: '',
      is_admin: false,
      site_names_raw: '',
      receive_email: true,
    }
  } else {
    user.value = props.user
  }
})

function close() {
  emit('close')
}

function submit() {
  // 执行保存
  isLoading.value = true

  // 默认标签颜色
  if (user.value!.badge_name !== '' && user.value!.badge_color === '') {
    user.value!.badge_color = '#0083ff'
  }

  if (isCreateMode.value) {
    artalk!.ctx.getApi().user.userAdd(user.value!, user.value!.password)
      .then((respUser) => {
        emit('update', respUser)
      }).catch((e) => {
        alert('用户创建错误：'+e.msg)
      }).finally(() => {
        isLoading.value = false
      })
  } else {
    artalk!.ctx.getApi().user.userEdit(user.value!, user.value!.password)
      .then((respUser) => {
        emit('update', respUser)
      }).catch((e) => {
        alert('用户保存错误：'+e.msg)
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
      <div><span>{{ t('comments') }}</span>{{(user as any).comment_count}}</div>
      <div><span>{{ t('last') }} IP</span>{{(user as any).last_ip || '-'}}</div>
      <div><span>{{ t('last') }} UA</span>
        <template v-if="showFullDetails || !(user as any).last_ua">{{(user as any).last_ua || '-'}}</template>
        <template v-else><span @click="showFullDetails = true" style="cursor: pointer;color: var(--at-color-main)">{{ t('Show') }}</span></template>
      </div>
    </div>
    <form v-if="user" class="atk-form" @submit.prevent="submit()">
      <div class="atk-label required">{{ t('username') }}</div>
      <input v-model="user.name" type="text" placeholder="" autocomplete="off">
      <div class="atk-label required">{{ t('email')  }}</div>
      <input v-model="user.email" type="text" placeholder="" autocomplete="off">
      <div class="atk-label">{{ t('link') }}</div>
      <input v-model="user.link" type="text" placeholder="" autocomplete="off">
      <div class="atk-label">{{ t('badgeText') }}</div>
      <input v-model="user.badge_name" type="text" placeholder="" autocomplete="off">
      <div class="atk-label">{{ t('badgeColor') }} (Color Hex)</div>
      <input v-model="user.badge_color" type="text" placeholder="" autocomplete="off">
      <div class="atk-label required">{{ t('role') }}</div>
      <select v-model="user.is_admin">
        <option :value="false">{{ t('normal') }}</option>
        <option :value="true">{{ t('admin') }}</option>
      </select>
      <template v-if="user.is_admin">
        <div class="atk-label required">{{ t('password') }}</div>
        <input v-model="user.password" type="text" :placeholder="isCreateMode ? '' : `(${t('passwordEmptyHint')})`" autocomplete="off">
        <div class="atk-label">{{ t('siteAttached') }}</div>
        <input v-model="user.site_names_raw" type="text" :placeholder="`(${t('siteEmptyHint')})`" autocomplete="off">
      </template>
      <div class="atk-label required">{{ t('emailNotify') }}</div>
      <select v-model="user.receive_email">
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

  &:hover i {
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
