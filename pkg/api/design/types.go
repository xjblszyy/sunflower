package design

import . "goa.design/goa/v3/dsl"

var (
	ExampleJwt  = "eyJhbGciOiJIUz..."
	ExampleUUID = "91cc3eb9-ddc0-4cf7-a62b-c85df1a9166f"
)

var SuccessResult = ResultType("SuccessResult", func() {
	Description("成功信息")
	ContentType("application/json")
	TypeName("SuccessResult")

	Attributes(func() {
		Field(1, "ok", Boolean, "success", func() {
			Example(true)
		})
		Required("ok")
	})

	View("default", func() {
		Attribute("ok")
	})
})

var GradeResult = ResultType("GradeResult", func() {
	Description("分数信息")
	ContentType("application/json")
	TypeName("GradeResult")

	Attributes(func() {
		Field(1, "id", String, "ID")
		Field(2, "class", String, "班级")
		Field(3, "name", String, "姓名")
		Field(4, "score", Int, "得分")
		Field(5, "subject", String, "科目")
		Field(6, "createdAt", String, "创建时间")
		Field(7, "updatedAt", String, "更新时间")
		Required("id", "class", "name", "score", "subject", "createdAt", "updatedAt")
	})

	View("default", func() {
		Attribute("id", String, "ID")
		Attribute("class", String, "班级")
		Attribute("name", String, "姓名")
		Attribute("score", Int, "得分")
		Attribute("subject", String, "科目")
		Attribute("createdAt", String, "创建时间")
		Attribute("updatedAt", String, "更新时间")
	})
})
