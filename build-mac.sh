# 测试有没有警告、问题 
# go mod vendor 
# go test -mod=vendor

# node 18+
# go 1.18+


current_dir=$(cd `dirname $0` && pwd)
dir_path=${current_dir}/frontend
echo "current_dir dir: "${current_dir}
echo "rebuild dir: "${dir_path}
# 构建
rm -rf ${dir_path}/package-lock.json
rm -rf ${dir_path}/node_modules
# rm -rf ${dir_path}/wailsjs
npm install --prefix ${dir_path}
npm run build  --prefix ${dir_path}
wails build