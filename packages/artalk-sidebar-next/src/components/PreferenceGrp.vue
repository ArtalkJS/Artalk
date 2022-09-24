<script setup lang="ts">
const props = defineProps<{
  title: string
  subTitle: string
  level: number
}>()
</script>

<template>
  <div class="pf-grp" :class="`level-${level}`">
    <div class="pf-head">
      <div class="title">{{ props.title }}</div>
      <div v-if="!!props.subTitle" class="sub-title">{{ props.subTitle }}</div>
    </div>
    <div class="pf-body">
      <slot />
    </div>
  </div>
</template>

<style scoped lang="scss">
.pf-grp.level-1 {
  & > .pf-head {
    margin-top: 30px;
    margin-bottom: 20px;

    .title {
      font-size: 1.4em;
      font-weight: bold;
      padding-left: 12px;

      &::before {
        position: absolute;
        top: calc(50% - 25px / 2);
        left: -10px;
        content: '';
        height: 25px;
        width: 7px;
        background: #B7DCFF;
      }
    }
  }

  & > .pf-body {
  }
}

.pf-grp.level-2 {
  & > .pf-head {
    margin-top: 15px;
    margin-bottom: 20px;

    .title {
      position: relative;
      font-weight: bold;
      font-size: 1.2em;
    }
  }
}

.pf-grp.level-1, .pf-grp.level-2 {
  & > .pf-head > .sub-title {
    padding: 0 10px 0 10px;
    margin-left: 4px;
    margin-top: 15px;
    border-left: 2px solid #eee;
  }
}


.pf-grp.level-3 {
  margin-left: 25px;

  & > .pf-head {
    margin-top: 15px;
    margin-bottom: 20px;

    & > .title {
      position: relative;
      font-weight: bold;
      font-size: 0.9em;
    }
  }
}

.pf-head {
  & > .title {
    position: relative;
  }

  & > .sub-title {
    font-size: 14px;
    margin-top: 5px;
  }
}
</style>
