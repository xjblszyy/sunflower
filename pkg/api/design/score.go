package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("Score", func() {
	Description("成绩系统")

	Error("internal_server_error", ErrorResult)
	Error("bad_request", ErrorResult)

	HTTP(func() {
		Path("/score")
		Response("bad_request", StatusBadRequest)
		Response("internal_server_error", StatusInternalServerError)
	})

	Method("Upload", func() {
		Meta("swagger:summary", "上传学生成绩文件")
		Description("上传学生成绩文件")

		Payload(func() {
			Attribute("content_type", String, "Content-Type header, must define value for multipart boundary.", func() {
				Default("multipart/form-data; boundary=goa")
				Pattern("multipart/[^;]+; boundary=.+")
				Example("multipart/form-data; boundary=goa")
			})
		})

		Result(func() {
			Attribute("errcode", Int, "错误码", func() {
				Minimum(0)
				Maximum(999999)
				Example(0)
			})
			Attribute("errmsg", String, "错误消息", func() {
				Example("")
			})
			Attribute("data", SuccessResult)
			Required("errcode", "errmsg")
		})

		HTTP(func() {
			POST("/upload")
			Header("content_type:Content-Type")
			SkipRequestBodyEncodeDecode()
			Response(StatusOK)
		})
	})

	Method("ScoreList", func() {
		Description("成绩列表")
		Meta("swagger:summary", "成绩列表")

		Payload(func() {
			Field(1, "cursor", Int, "cursor of page", func() {
				Example(0)
			})
			Field(2, "limit", Int, "limit of items", func() {
				Example(20)
				Maximum(100)
			})
			Field(3, "sortField", String, "排序字段 默认createdAt", func() {
				Enum("class", "name", "score", "subject", "createdAt")
			})
			Field(4, "sortOrder", String, "排序方式 默认desc", func() {
				Enum("asc", "desc")
			})
			Field(5, "name", String, "姓名")
			Field(6, "class", String, "班级")
			Field(7, "scores", Int, "分数")
			Field(8, "subject", String, "科目")
			Required("cursor", "limit")
		})

		Result(func() {
			Field(1, "errcode", Int, "错误码", func() {
				Minimum(0)
				Maximum(999999)
				Example(0)
			})
			Field(2, "errmsg", String, "错误消息", func() {
				Example("")
			})
			Field(3, "data", ArrayOf(GradeResult), "结果")
			Field(4, "nextCursor", Int, "下一页游标")
			Field(5, "total", Int, "总记录数")
			Required("errcode", "errmsg")
		})

		HTTP(func() {
			GET("/")
			Params(func() {
				Param("sortField")
				Param("cursor")
				Param("limit")
				Param("sortOrder")
				Param("name")
				Param("class")
				Param("scores")
				Param("subject")
			})
			Response(StatusOK)
		})

		GRPC(func() {
			Response(CodeOK)
		})
	})

	Method("ScoreDetail", func() {
		Description("成绩详情")
		Meta("swagger:summary", "成绩详情")

		Payload(func() {
			Field(1, "id", Int, "id", func() {
				Example(1)
			})

			Required("id")
		})

		Result(func() {
			Field(1, "errcode", Int, "错误码", func() {
				Minimum(0)
				Maximum(999999)
				Example(0)
			})
			Field(2, "errmsg", String, "错误消息", func() {
				Example("")
			})
			Field(3, "data", CollectionOf(GradeResult, func() {
				View("default")
			}))
			Field(4, "nextCursor", Int, "下一页游标", func() {
				Example(100)
			})
			Field(5, "total", Int, "总记录数")
			Required("errcode", "errmsg")
		})

		Result(func() {
			Field(1, "errcode", Int, "错误码", func() {
				Minimum(0)
				Maximum(999999)
				Example(0)
			})
			Field(2, "errmsg", String, "错误消息", func() {
				Example("")
			})
			Field(3, "data", GradeResult, "结果")
			Required("errcode", "errmsg")
		})

		HTTP(func() {
			GET("/{id}")
			Params(func() {
				Param("id")
			})
			Response(StatusOK)
		})

		GRPC(func() {
			Response(CodeOK)
		})
	})
})
