# 数据同步API参考实现

[数据同步API定义(v1版本)](https://github.com/idaaser/syncspecv1/blob/master/README.md)

开发者可以依此为基础, 通过自定义如下接口来从其他上游获取数据
- [AuthnStore](server/authn_store.go): 定义了如何颁发access_token, 以及如何校验access_token
- [ContactStore](server/contact_store.go): 定义了如何拉取用户、部门数据