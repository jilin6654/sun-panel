<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { NEllipsis } from 'naive-ui'
import { ItemIcon } from '@/components/common'
import { PanelPanelConfigStyleEnum } from '@/enums'
import { changeServerState, getServerState } from '@/api/panel/itemIcon'

interface Prop {
  itemInfo?: Panel.ItemInfo
  size?: number // 默认70
  forceBackground?: string // 强制背景色
  iconTextColor?: string
  iconTextInfoHideDescription: boolean
  iconTextIconHideTitle: boolean
  style: PanelPanelConfigStyleEnum
}

const props = withDefaults(defineProps<Prop>(), {
  size: 70,
})

let timer: NodeJS.Timer
const isServer = ref<boolean>(false)
const isLoadingServerState = ref<boolean>(false)
const serverState = ref<number>(0)
const url = ref<string>('')
const mac = ref<string>('')
const defaultBackground = '#2a2a2a6b'

const calculateLuminance = (color: string) => {
  const hex = color.replace(/^#/, '')
  const r = parseInt(hex.substring(0, 2), 16)
  const g = parseInt(hex.substring(2, 4), 16)
  const b = parseInt(hex.substring(4, 6), 16)
  return (0.299 * r + 0.587 * g + 0.114 * b) / 255
}

const textColor = computed(() => {
  const luminance = calculateLuminance(props.itemInfo?.icon?.backgroundColor || defaultBackground)
  return luminance > 0.5 ? 'black' : 'white'
})

onMounted(() => {
  url.value = props.itemInfo?.lanUrl || ''
  mac.value = props.itemInfo?.description || ''
  isServer.value = props.itemInfo?.title.includes('服务器') || false
  if (!isServer.value)
    return
  getData()
})

onUnmounted(() => {
  clearTimeout(timer)
})

const serverIsRunning = computed(() => serverState.value === 1)

async function getData() {
  timer = setTimeout(async () => {
    try {
      if (!url.value)
        return
      const { data, code } = await getServerState<number>(url.value)
      if (code === 0)
        serverState.value = data
      await getData()
    }
    catch (error) {

    }
  }, 3000)
}

async function handleChangeState() {
  console.log('handle', url.value, mac.value, isServer.value)
  if (!url.value || !mac.value || isLoadingServerState.value || !isServer.value)
    return
  isLoadingServerState.value = true
  try {
    const { data, code } = await changeServerState<number>(url.value, mac.value, serverState.value)
    if (code === 0)
      serverState.value = data
  }
  catch (error) {

  }
  finally {
    isLoadingServerState.value = false
  }
}
</script>

<template>
  <div class="app-icon w-full" @click.stop="handleChangeState()">
    <!-- 详情图标 -->
    <div
      v-if="style === PanelPanelConfigStyleEnum.info"
      class="app-icon-info w-full rounded-2xl  transition-all duration-200 hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)] flex"
      :style="{ background: itemInfo?.icon?.backgroundColor || defaultBackground }"
    >
      <!-- 图标 -->
      <div class="app-icon-info-icon w-[70px] h-[70px]">
        <div class="w-[70px] h-full flex items-center justify-center ">
          <ItemIcon :item-icon="itemInfo?.icon" force-background="transparent" :size="50" class="overflow-hidden rounded-xl" />
        </div>
        <span v-if="isServer" class="ml-[2px] top-[-80px] server-status" :class="[serverIsRunning ? 'before:bg-green-600 after:bg-green-600' : 'before:bg-zinc-600 after:bg-zinc-600']" />
      </div>

      <!-- 文字 -->
      <!-- 如果为纯白色，将自动根据背景的明暗计算字体的黑白色 -->
      <div class="text-white flex items-center" :style="{ color: (iconTextColor === '#ffffff') ? textColor : iconTextColor, maxWidth: 'calc(100% - 80px)' }">
        <div class="app-icon-info-text-box w-full">
          <div class="app-icon-info-text-box-title font-semibold w-full">
            <NEllipsis>
              {{ itemInfo?.title }}
            </NEllipsis>
          </div>
          <div v-if="!iconTextInfoHideDescription" class="app-icon-info-text-box-description">
            <NEllipsis :line-clamp="2" class="text-xs">
              {{ itemInfo?.description }}
            </NEllipsis>
          </div>
        </div>
      </div>
    </div>

    <!-- 极简(小)图标（APP） -->
    <div v-if="style === PanelPanelConfigStyleEnum.icon" class="app-icon-small">
      <div
        class="app-icon-small-icon overflow-hidden rounded-2xl sunpanel w-[70px] h-[70px] mx-auto rounded-2xl transition-all duration-200 hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)]"
        :title="itemInfo?.description"
      >
        <ItemIcon :item-icon="itemInfo?.icon" />
      </div>
      <span v-if="isServer" class="ml-[2px] top-[-80px] server-status" :class="[serverIsRunning ? 'before:bg-green-600 after:bg-green-600' : 'before:bg-zinc-600 after:bg-zinc-600']" />
      <div
        v-if="!iconTextIconHideTitle"
        class="app-icon-small-title text-center app-icon-text-shadow cursor-pointer mt-[2px]"
        :class="[isServer ? 'top-[-22px]' : '']"
        :style="{ color: iconTextColor, position: 'relative' }"
      >
        <span>{{ itemInfo?.title }}</span>
      </div>
    </div>
  </div>
</template>

<style>
.server-status {
  position: relative;
  width: 20px;
  height: 20px;
}

.server-status::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 6px;
  height: 6px;
  /* background-color: #00ff00; */
  border-radius: 50%;
  transform: translate(-50%, -50%);
}

.server-status::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 6px;
  height: 6px;
  /* background-color: #00ff00; */
  border-radius: 50%;
  transform: translate(-50%, -50%);
  animation: pulse 1.5s ease-out infinite;
}

@keyframes pulse {
  0% {
    transform: translate(-50%, -50%) scale(1);
    opacity: 1;
  }
  100% {
    transform: translate(-50%, -50%) scale(4);
    opacity: 0;
  }
}
</style>
