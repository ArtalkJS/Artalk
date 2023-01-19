<script setup lang="ts">
import { useNavStore } from '../stores/nav'
import { artalk, bootParams } from '../global'
import Pagination from '../components/Pagination.vue'
import type { UserDataForAdmin } from 'artalk/types/artalk-data'
import { storeToRefs } from 'pinia'

const nav = useNavStore()
const router = useRouter()
const { curtTab } = storeToRefs(nav)
const users = ref<UserDataForAdmin[]>([])

const pageSize = ref(30)
const pageTotal = ref(0)
const pagination = ref<InstanceType<typeof Pagination>>()
const curtType = ref('all')

const addingUser = ref(false)
const editingUser = ref<UserDataForAdmin|undefined>()

onMounted(() => {
  nav.updateTabs({
    'all': '全部',
    'admin': '管理员',
    'create': '新增'
  }, 'all')

  watch(curtTab, (tab) => {
    if (tab === 'create') {
      if (editingUser.value !== undefined) {
        editingUser.value = undefined
        nextTick(() => {
          addingUser.value = true
        })
      } else {
        addingUser.value = true
      }
    } else {
      addingUser.value = false
      editingUser.value = undefined

      curtType.value = tab
      pagination.value?.reset()
      reqUsers(0)
    }
  })

  reqUsers(0)
})

function reqUsers(offset: number) {
  nav.setPageLoading(true)
  artalk?.ctx.getApi().user.userList(offset, pageSize.value, curtType.value as any)
    .then(got => {
      pageTotal.value = got.total
      users.value = got.users
      nav.scrollPageToTop()
    }).finally(() => {
      nav.setPageLoading(false)
    })
}

function onChangePage(offset: number) {
  reqUsers(offset)
}

function editUser(user: UserDataForAdmin) {
  if (user.is_in_conf) {
    alert('暂不支持在线编辑配置文件中的用户，请手动修改配置文件')
    return
  }

  addingUser.value = false
  editingUser.value = user
}

function updateUser(user: UserDataForAdmin) {
  const index = users.value.findIndex(u => u.id === user.id)
  if (index != -1) {
    // 修改用户
    const orgUser = users.value[index]
    Object.keys(users).forEach(key => {
      ;(orgUser as any)[key] = (users as any)[key]
    })
  } else {
    // 创建用户
    pagination.value!.reset()
    reqUsers(0)
  }

  closeEditUser()
}

function closeEditUser() {
  addingUser.value = false
  editingUser.value = undefined
}

function delUser(user: UserDataForAdmin) {
  if (window.confirm(`该操作将删除 用户："${user.name}" 邮箱："${user.email}" 所有评论，包括其评论下面他人的回复评论，是否继续？`)) {
    artalk!.ctx.getApi().user.userDel(user.id)
      .then(() => {
        const index = users.value.findIndex(u => u.id === user.id)
        users.value.splice(index, 1)

        if (user.is_in_conf) {
          alert('用户已从数据库删除，请手动编辑配置文件并删除用户')
        }
      })
      .catch((e) => {
        alert('删除失败：'+e.msg)
      })
  }
}
</script>

<template>
  <div class="user-list-wrap">
    <div class="user-list">
      <div v-for="(user) in users" class="user-item">
        <div class="user-main">
          <div class="title">
            {{ user.name }}
            <span class="badge-grp">
              <span v-if="user.badge_name" class="badge" :style="{ backgroundColor: user.badge_color }">{{user.badge_name}}</span>
              <span v-else-if="user.is_admin" class="badge admin" title="该用户具有管理员权限">管理员</span>
              <span v-if="user.is_in_conf" class="badge in-conf" title="该用户存在于配置文件中">配置文件</span>
            </span>
          </div>
          <div class="sub">{{ user.email }}</div>
        </div>
        <div class="user-actions">
          <span @click="editUser(user)">编辑</span>
          <span>评论 ({{user.comment_count}})</span>
          <span @click="delUser(user)">删除</span>
        </div>
      </div>
    </div>
    <Pagination
      ref="pagination"
      :pageSize="pageSize"
      :total="pageTotal"
      :disabled="nav.isPageLoading"
      @change="onChangePage"
    />

    <UserEditor
      v-if="addingUser || editingUser !== undefined"
      :user="editingUser"
      @update="updateUser"
      @close="closeEditUser"
    />
  </div>
</template>

<style scoped lang="scss">
.user-list-wrap {
  .user-list {
    .user-item {
      padding: 15px 30px;

      &:not(:last-child) {
        border-bottom: 1px solid var(--at-color-border);
      }

      .user-main {
        .title {
          display: flex;
          flex-direction: row;
          align-items: center;
          color: var(--at-color-deep);
          font-size: 20px;

          .badge-grp {
            display: flex;
            margin-left: 5px;
          }

          .badge {
            margin-left: 6px;
            font-size: 13px;
            color: var(--at-color-meta);
            background: var(--at-color-bg-grey);
            padding: 0 6px;
            line-height: 17px;
            border-radius: 2px;
            color: #fff;

            &.admin {
              background: #ff6c00;
            }

            &.in-conf {
              background: #297bff;
            }
          }
        }

        .sub {
          color: var(--at-color-sub);
          font-size: 15px;
          margin-top: 5px;
        }
      }

      .user-actions {
        margin-top: 10px;

        & > span {
          cursor: pointer;
          color: var(--at-color-meta);
          font-size: 13px;

          &:not(:last-child) {
            margin-right: 16px;
          }

          &:hover {
            color: var(--at-color-deep);
          }
        }
      }
    }
  }
}
</style>
