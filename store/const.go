package store

const casbinText = `

[request_definition]
r = sub,obj,act

[policy_definition]
p = sub,obj,act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act
`

const textExample = `
# 请求
[request_definition]
r = sub,obj,act
# sub ——> 想要访问资源的用户角色(Subject)——请求实体
# obj ——> 访问的资源(Object)
# act ——> 访问的方法(Action: get、post...)

# 策略(.csv文件p的格式，定义的每一行为policy rule;p,p2为policy rule的名字。)
[policy_definition]
p = sub,obj,act
# p2 = sub,act 表示sub对所有资源都能执行act

# 组定义
[role_definition]
g = _, _
# _,_表示用户，角色/用户组
# g2 = _,_,_ 表示用户, 角色/用户组, 域(也就是租户)

# 策略效果
[policy_effect]
e = some(where (p.eft == allow))
# 上面表示有任意一条 policy rule 满足, 则最终结果为 allow；p.eft它可以是allow或deny，它是可选的，默认是allow

# 匹配器
[matchers]
# m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act

#m = r.sub == p.sub && keyMatch(r.obj,p.obj) && (r.act==p.act || p.act == "*") || r.sub=="superTest"
# r.sub="xxx"表示实体为superTest的 deny
`
