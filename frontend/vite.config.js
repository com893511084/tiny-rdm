import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Icons from 'unplugin-icons/vite'
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'
import { defineConfig } from 'vite'
import path from 'path' // 引入 path 模块


// 在Windows 中的路径是 D:\D:\%E9%A1%B9%E7%9B%AE\tiny-rdm\frontend\src\utils\i18n.js 有问题
// const rootPath = new URL('.', import.meta.url).pathname
// 使用 path.resolve() 方法来构建绝对路径
const rootPath = path.resolve(__dirname, './')

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        AutoImport({
            imports: [
                {
                    'naive-ui': ['useDialog', 'useMessage', 'useNotification', 'useLoadingBar'],
                },
            ],
        }),
        Components({
            resolvers: [NaiveUiResolver()],
        }),
        Icons(),
    ],
    resolve: {
        alias: {
            '@': rootPath + '/src',
            stores: rootPath + '/src/stores',
            wailsjs: rootPath + '/wailsjs',
        },
    },
})
