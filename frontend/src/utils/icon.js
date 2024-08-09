const iconType = [
    {
        key: "宠物",
        val: "el-icon-cat",
    },
    { key: "交通出行", val: "el-icon-car" },
    { key: "充值缴费", val: "el-icon-shopping" },
    { key: "医疗健康", val: "" },
    { key: "商业服务", val: "" },
    { key: "家居家装", val: "el-icon-house" },
    { key: "数码电器", val: "el-icon-shopping" },
    { key: "文化休闲", val: "el-icon-coffee" },
    { key: "日用百货", val: "el-icon-shopping" },
    { key: "服饰装扮", val: "" },
    { key: "爱车养车", val: "el-icon-car" },
    { key: "生活服务", val: "el-icon-coffee" },
    { key: "美容美发", val: "" },
    { key: "转账红包", val: "" },
    { key: "酒店旅游", val: "el-icon-house" },
    { key: "餐饮美食", val: "el-icon-eat" },
];

// 转换icon (%s )
export function getIconClass(typeName) {
    var iconStr = "";
    iconType.forEach((v) => {
        if (v.key == typeName) {
            iconStr = v.val;
        }
    });
    if (iconStr === "" || iconStr.length === 0) {
        iconStr = "el-icon-shopping";
    }
    // console.debug("getIconClass2", iconStr);
    return iconStr;
}
