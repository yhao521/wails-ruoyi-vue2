package jwt

import (
	"fmt"
	"mySparkler/pkg/cache"
	"mySparkler/pkg/utils"
	"strconv"
)

func CacheUserId(userId int) {
	fmt.Println("CacheUserId: " + strconv.Itoa(userId))
	cache.Cache("").Put("userId", strconv.Itoa(userId), 0)
}

func CacheCleanUserId() {
	cache.Cache("").Del("userId")
}
func CacheCleanAll() {
	cache.Cache("").Clear()
}

func CacheGetUserId() int {
	userId, _ := cache.Cache("").Get("userId")
	fmt.Println("CacheGetUserId: " + userId)
	i := utils.GetInterfaceToInt(userId)
	// i, err := strconv.Atoi(userId)
	// if err != nil {
	// 	// 转换错误的处理逻辑
	// 	return 0
	// }
	return i
}

func CacheRoleId(roleId int) {
	fmt.Println("CacheRoleId: " + strconv.Itoa(roleId))
	cache.Cache("").Put("roleId", strconv.Itoa(roleId), 0)
}

func CacheGetRoleId() string {
	roleId, _ := cache.Cache("").Get("roleId")
	fmt.Println("CacheGetRoleId: " + roleId)
	return roleId
}
