// 自定义代码编辑器组件
//src\components\monacoEditor\index.vue
<template>
    <div class="full flex flex-column">
        <div class="toolbar flex align-center">
            <el-tooltip class="item" effect="dark" content="格式化" placement="top-start">
                <span class="format" @click.stop="handleFormat"> { } </span>
            </el-tooltip>
            <el-tooltip class="item" effect="dark" content="全屏" placement="top-start">
                <span class="screen" @click.stop="handleScreen">
                    <i class="el-icon-full-screen"></i>
                </span>
            </el-tooltip>
        </div>
        <div class="flex-1 mt10 main" style="height:100%;" ref="main"></div>
    </div>
</template>

<script>
import loadMonaco from '@/utils/loadMonaco';
// import * as monaco from "monaco-editor";
// import * as monaco from 'monaco-editor/esm/vs/editor/editor.api.js';
import beautify from 'js-beautify';

export default {
    model: {
        prop: 'value',
        event: 'onchange'
    },
    props: {
        value: {
            type: String,
            default: '<div id="container" style="height:100%;"></div>'
        },
        config: [Object]
    },
    data() {
        return {
            code: '',
            monacoEditor: null,
            fullscreen: false,
            configType: {
                html: {
                    indent_size: '2',
                    indent_char: ' ',
                    max_preserve_newlines: '-1',
                    preserve_newlines: false,
                    keep_array_indentation: false,
                    break_chained_methods: false,
                    indent_scripts: 'separate',
                    brace_style: 'end-expand',
                    space_before_conditional: true,
                    unescape_strings: false,
                    jslint_happy: false,
                    end_with_newline: true,
                    wrap_line_length: '110',
                    indent_inner_html: true,
                    comma_first: false,
                    e4x: true,
                    indent_empty_lines: true
                },
                css: {
                    indent_size: '2',
                    indent_char: ' ',
                    max_preserve_newlines: '-1',
                    preserve_newlines: false,
                    keep_array_indentation: false,
                    break_chained_methods: false,
                    indent_scripts: 'separate',
                    brace_style: 'end-expand',
                    space_before_conditional: true,
                    unescape_strings: false,
                    jslint_happy: false,
                    end_with_newline: true,
                    wrap_line_length: '110',
                    indent_inner_html: true,
                    comma_first: false,
                    e4x: true,
                    indent_empty_lines: true
                },
                js: {
                    indent_size: '2',
                    indent_char: ' ',
                    max_preserve_newlines: '-1',
                    preserve_newlines: false,
                    keep_array_indentation: false,
                    break_chained_methods: false,
                    indent_scripts: 'normal',
                    brace_style: 'end-expand',
                    space_before_conditional: true,
                    unescape_strings: false,
                    jslint_happy: true,
                    end_with_newline: true,
                    wrap_line_length: '110',
                    indent_inner_html: true,
                    comma_first: false,
                    e4x: true,
                    indent_empty_lines: true
                }
            }
        };
    },
    mounted() {
        this.init();
    },
    watch: {
        value(val) {
            if (this.monacoEditor.getValue() !== val) {
                this.monacoEditor.setValue(val);
            }
        }
    },
    methods: {
        init() {
            // this.monacoEditor = monaco.editor.create(this.$refs.main, {
            //     theme: 'vs-dark', // 主题
            //     value: this.value, // 默认显示的值
            // });
            loadMonaco(monaco => {
                this.monacoEditor = monaco.editor.create(this.$refs.main, {
                    theme: 'vs-dark', // 主题
                    value: this.value, // 默认显示的值
                    language: this.config.language || 'html',
                    folding: true, // 是否折叠
                    minimap: {
                        enabled: false // 是否启用预览图
                    }, // 预览图设置
                    foldingHighlight: true, // 折叠等高线
                    foldingStrategy: 'indentation', // 折叠方式  auto | indentation
                    showFoldingControls: 'always', // 是否一直显示折叠 always | mouseover
                    disableLayerHinting: true, // 等宽优化
                    emptySelectionClipboard: false, // 空选择剪切板
                    selectionClipboard: false, // 选择剪切板
                    automaticLayout: true, // 自动布局
                    codeLens: true, // 代码镜头
                    scrollBeyondLastLine: false, // 滚动完最后一行后再滚动一屏幕
                    colorDecorators: true, // 颜色装饰器
                    accessibilitySupport: 'on', // 辅助功能支持  "auto" | "off" | "on"
                    lineNumbers: 'on', // 行号 取值： "on" | "off" | "relative" | "interval" | function
                    lineNumbersMinChars: 5, // 行号最小字符   number
                    enableSplitViewResizing: false,
                    overviewRulerBorder: false, // 是否应围绕概览标尺绘制边框
                    renderLineHighlight: 'gutter', // 当前行突出显示方式
                    readOnly: false //是否只读  取值 true | false
                });
                this.monacoEditor.onDidChangeModelContent(e => {
                    let val = this.monacoEditor.getValue();
                    this.$emit('onchange', val);
                    this.$emit('getValue', val);
                });
            });
        },
        handleScreen() {
            let element = this.$refs.main;
            if (this.fullscreen) {
                if (document.exitFullscreen) {
                    document.exitFullscreen();
                } else if (document.webkitCancelFullScreen) {
                    document.webkitCancelFullScreen();
                } else if (document.mozCancelFullScreen) {
                    document.mozCancelFullScreen();
                } else if (document.msExitFullscreen) {
                    document.msExitFullscreen();
                }
            } else {
                if (element.requestFullscreen) {
                    element.requestFullscreen();
                } else if (element.webkitRequestFullScreen) {
                    element.webkitRequestFullScreen();
                } else if (element.mozRequestFullScreen) {
                    element.mozRequestFullScreen();
                } else if (element.msRequestFullscreen) {
                    element.msRequestFullscreen();
                }
            }
            this.fullscreen = !this.fullscreen;
        },
        handleFormat() {
            const code = this.monacoEditor.getValue();
            let type = undefined;
            let newCode = undefined;
            if (this.config.language === 'javascript') {
                type = this.configType.js;
                newCode = beautify.js(code, type);
            } else if (this.config.language === 'css') {
                type = this.configType.css;
                newCode = beautify.css(code, type);
            } else {
                type = this.configType.html;
                newCode = beautify.html(code, type);
            }
            this.monacoEditor.setValue(newCode);
        }
    },
    beforeDestroy() {
        // this.monacoEditor?.dispose();
    }
};
</script>

<style lang="scss" scoped>
.full {
    width: 100%;
    height: 100%;
    overflow: hidden;
    position: relative;
}

.screen,
.format {
    display: inline-block;
    width: 20px;
    height: 20px;
    line-height: 20px;
    cursor: pointer;
    font-size: 20px;
}

.format {
    float: left;
    font-size: 18px;
    margin: 7px 10px 0;
}

.el-icon-full-screen {
    float: left;
    margin-top: 5px;
}

.main {
    min-height: 0;

    &::v-deep {
        .monaco-editor {
            height: 100% !important;
        }
    }
}

.toolbar {
    position: absolute;
    right: 15px;
    top: 10px;
    color: #fff;
    z-index: 2;
}
</style>