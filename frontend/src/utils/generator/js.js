import { exportDefault, titleCase } from '@/utils/index'
import ruleTrigger from './ruleTrigger'

const units = {
  KB: '1024',
  MB: '1024 / 1024',
  GB: '1024 / 1024 / 1024'
}
let confGlobal
const inheritAttrs = {
  noform: '',
  file: '',
  dialog: 'inheritAttrs: false,'
}


export function makeUpJs(conf, type, fileName) {
  confGlobal = conf = JSON.parse(JSON.stringify(conf))
  const importList = []
  const componentList = []
  const dataList = []
  const ruleList = []
  const optionsList = []
  const propsList = []
  const methodList = mixinMethod(type)
  const uploadVarList = []
  console.debug('makeUpJs', conf)

  conf.fields.forEach(el => {
    buildAttributes(el, dataList, ruleList, optionsList, methodList, propsList, uploadVarList, importList, componentList)
  })
  if (fileName) {
    fileName = fileName.substring(0, fileName.lastIndexOf("."))
  }
  //截取文件名的. 前缀
  const script = buildexport(
    conf,
    type,
    fileName,
    dataList.join('\n'),
    ruleList.join('\n'),
    optionsList.join('\n'),
    uploadVarList.join('\n'),
    propsList.join('\n'),
    methodList.join('\n'),
    importList.join('\n'),
    componentList.join('\n')
  )
  confGlobal = null
  return script
}

function buildAttributes(el, dataList, ruleList, optionsList, methodList, propsList, uploadVarList, importList, componentList) {
  buildData(el, dataList)
  buildRules(el, ruleList)

  buildImoprt(el, importList)
  buildComponents(el, componentList)

  if (el.options && el.options.length) {
    buildOptions(el, optionsList)
    if (el.dataType === 'dynamic') {
      const model = `${el.vModel}Options`
      const options = titleCase(model)
      buildOptionMethod(`get${options}`, model, methodList)
    }
  }

  if (el.props && el.props.props) {
    buildProps(el, propsList)
  }

  if (el.action && el.tag === 'el-upload') {
    uploadVarList.push(
      `${el.vModel}Action: '${el.action}',
      ${el.vModel}fileList: [],`
    )
    methodList.push(buildBeforeUpload(el))
    if (!el['auto-upload']) {
      methodList.push(buildSubmitUpload(el))
    }
  }

  if (el.children) {
    el.children.forEach(el2 => {
      buildAttributes(el2, dataList, ruleList, optionsList, methodList, propsList, uploadVarList)
    })
  }
}

function mixinMethod(type) {
  const list = []; const
    minxins = {
      file: confGlobal.formBtns ? {
        submitForm: `submitForm() {
        this.$refs['${confGlobal.formRef}'].validate(valid => {
          if(!valid) return
          // TODO 提交表单
        })
      },`,
        resetForm: `resetForm() {
        this.$refs['${confGlobal.formRef}'].resetFields()
      },`
      } : null,
      dialog: {
        onOpen: 'onOpen() {},',
        onClose: `onClose() {
        this.$refs['${confGlobal.formRef}'].resetFields()
      },`,
        close: `close() {
        this.$emit('update:visible', false)
      },`,
        handleConfirm: `handleConfirm() {
        this.$refs['${confGlobal.formRef}'].validate(valid => {
          if(!valid) return
          this.close()
        })
      },`
      }
    }

  const methods = minxins[type]
  if (methods) {
    Object.keys(methods).forEach(key => {
      list.push(methods[key])
    })
  }

  return list
}
// imports
function buildImoprt(conf, importList) {
  console.debug('buildImoprt', conf)
  if (conf.__config__.import === undefined) return

  importList.push(`${conf.__config__.import}`)
  console.debug('buildImoprt-importList', importList)
}
//components
function buildComponents(conf, componentList) {
  console.debug('buildComponents', conf)
  let component = ""
  if (conf.__config__.tag === undefined) return
  if (conf.__config__.tagComponent != undefined) {
    component = conf.__config__.tagComponent
  } else {
    component = conf.__config__.tag
  }
  if (!component.startsWith('el-')) {
    componentList.push(`${component},`)
  }
  console.debug('buildComponents-componentList', componentList)
}
function buildData(conf, dataList) {
  if (conf.vModel === undefined) return
  let defaultValue
  if (typeof (conf.defaultValue) === 'string' && !conf.multiple) {
    defaultValue = `'${conf.defaultValue}'`
  } else {
    defaultValue = `${JSON.stringify(conf.defaultValue)}`
  }
  dataList.push(`${conf.vModel}: ${defaultValue},`)
}

function buildRules(conf, ruleList) {
  if (conf.vModel === undefined) return
  const rules = []
  if (ruleTrigger[conf.tag]) {
    if (conf.required) {
      const type = Array.isArray(conf.defaultValue) ? 'type: \'array\',' : ''
      let message = Array.isArray(conf.defaultValue) ? `请至少选择一个${conf.vModel}` : conf.placeholder
      if (message === undefined) message = `${conf.label}不能为空`
      rules.push(`{ required: true, ${type} message: '${message}', trigger: '${ruleTrigger[conf.tag]}' }`)
    }
    if (conf.regList && Array.isArray(conf.regList)) {
      conf.regList.forEach(item => {
        if (item.pattern) {
          rules.push(`{ pattern: ${eval(item.pattern)}, message: '${item.message}', trigger: '${ruleTrigger[conf.tag]}' }`)
        }
      })
    }
    ruleList.push(`${conf.vModel}: [${rules.join(',')}],`)
  }
}

function buildOptions(conf, optionsList) {
  if (conf.vModel === undefined) return
  if (conf.dataType === 'dynamic') { conf.options = [] }
  const str = `${conf.vModel}Options: ${JSON.stringify(conf.options)},`
  optionsList.push(str)
}

function buildProps(conf, propsList) {
  if (conf.dataType === 'dynamic') {
    conf.valueKey !== 'value' && (conf.props.props.value = conf.valueKey)
    conf.labelKey !== 'label' && (conf.props.props.label = conf.labelKey)
    conf.childrenKey !== 'children' && (conf.props.props.children = conf.childrenKey)
  }
  const str = `${conf.vModel}Props: ${JSON.stringify(conf.props.props)},`
  propsList.push(str)
}

function buildBeforeUpload(conf) {
  const unitNum = units[conf.sizeUnit]; let rightSizeCode = ''; let acceptCode = ''; const
    returnList = []
  if (conf.fileSize) {
    rightSizeCode = `let isRightSize = file.size / ${unitNum} < ${conf.fileSize}
    if(!isRightSize){
      this.$message.error('文件大小超过 ${conf.fileSize}${conf.sizeUnit}')
    }`
    returnList.push('isRightSize')
  }
  if (conf.accept) {
    acceptCode = `let isAccept = new RegExp('${conf.accept}').test(file.type)
    if(!isAccept){
      this.$message.error('应该选择${conf.accept}类型的文件')
    }`
    returnList.push('isAccept')
  }
  const str = `${conf.vModel}BeforeUpload(file) {
    ${rightSizeCode}
    ${acceptCode}
    return ${returnList.join('&&')}
  },`
  return returnList.length ? str : ''
}

function buildSubmitUpload(conf) {
  const str = `submitUpload() {
    this.$refs['${conf.vModel}'].submit()
  },`
  return str
}

function buildOptionMethod(methodName, model, methodList) {
  const str = `${methodName}() {
    // TODO 发起请求获取数据
    this.${model}
  },`
  methodList.push(str)
}

function buildexport(conf, type, fileName, data, rules, selectOptions, uploadVar, props, methods, importList, components) {
  const str = `
  ${importList}

  ${exportDefault}{
  ${inheritAttrs[type]}
  name: '${fileName}',
  components: {
     ${components}
  },
  props: [],
  data () {
    return {
      ${conf.formModel}: {
        ${data}
      },
      ${conf.formRules}: {
        ${rules}
      },
      ${uploadVar}
      ${selectOptions}
      ${props}
    }
  },
  computed: {},
  watch: {},
  created () {},
  mounted () {},
  methods: {
    ${methods}
  }
}`
  return str
}
