
//代码编辑器
import monacoEditor from '@/components/template/monacoEditor';
// 导航栏+表格组件
import yhContainer from '@/components/template/Container';
// 全局注册富文本组件
import tinymce from '@/components/template/tinymce/index.vue'
// 表格组件模板
import yhTableView from '@/components/template/yhTableView/index.vue'
// md渲染/解析模板
import yhMarked from '@/components/template/yhMarked/index.vue'


export function setTemplate(vue) {
    vue.component('monacoEditor', monacoEditor);
    vue.component('yhContainer', yhContainer);
    vue.component('tinymce', tinymce)
    vue.component('yhTableView', yhTableView)
    vue.component('yhMarked', yhMarked)

}